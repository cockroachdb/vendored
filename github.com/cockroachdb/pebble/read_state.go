// Copyright 2019 The LevelDB-Go and Pebble Authors. All rights reserved. Use
// of this source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package pebble

import (
	"fmt"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// readState encapsulates the state needed for reading (the current version and
// list of memtables). Loading the readState is done without grabbing
// DB.mu. Instead, a separate DB.readState.RWMutex is used for
// synchronization. This mutex solely covers the current readState object which
// means it is rarely or ever contended.
//
// Note that various fancy lock-free mechanisms can be imagined for loading the
// readState, but benchmarking showed the ones considered to purely be
// pessimizations. The RWMutex version is a single atomic increment for the
// RLock and an atomic decrement for the RUnlock. It is difficult to do better
// than that without something like thread-local storage which isn't available
// in Go.
type readState struct {
	db        *DB
	refcnt    int32
	current   *version
	memtables flushableList

	// DEBUGGING.
	iterMap sync.Map
	doneC   chan struct{}
}

// ref adds a reference to the readState.
func (s *readState) ref() {
	atomic.AddInt32(&s.refcnt, 1)
}

// unref removes a reference to the readState. If this was the last reference,
// the reference the readState holds on the version is released. Requires DB.mu
// is NOT held as version.unref() will acquire it. See unrefLocked() if DB.mu
// is held by the caller.
func (s *readState) unref() {
	if atomic.AddInt32(&s.refcnt, -1) != 0 {
		return
	}
	s.current.Unref()
	for _, mem := range s.memtables {
		mem.readerUnref()
	}
	s.cleanupDebug()

	// The last reference to the readState was released. Check to see if there
	// are new obsolete tables to delete.
	s.db.maybeScheduleObsoleteTableDeletion()
}

// unrefLocked removes a reference to the readState. If this was the last
// reference, the reference the readState holds on the version is
// released. Requires DB.mu is held as version.unrefLocked() requires it. See
// unref() if DB.mu is NOT held by the caller.
func (s *readState) unrefLocked() {
	if atomic.AddInt32(&s.refcnt, -1) != 0 {
		return
	}
	s.current.UnrefLocked()
	for _, mem := range s.memtables {
		mem.readerUnref()
	}
	s.cleanupDebug()

	// NB: Unlike readState.unref(), we don't attempt to cleanup newly obsolete
	// tables as unrefLocked() is only called during DB shutdown to release the
	// current readState.
}

// loadReadState returns the current readState. The returned readState must be
// unreferenced when the caller is finished with it.
func (d *DB) loadReadState() *readState {
	d.readState.RLock()
	state := d.readState.val
	state.ref()
	d.readState.RUnlock()
	return state
}

// updateReadStateLocked creates a new readState from the current version and
// list of memtables. Requires DB.mu is held. If checker is not nil, it is called after installing
// the new readState
func (d *DB) updateReadStateLocked(checker func(*DB) error) {
	s := &readState{
		db:        d,
		refcnt:    1,
		current:   d.mu.versions.currentVersion(),
		memtables: d.mu.mem.queue,
	}
	s.current.Ref()
	for _, mem := range s.memtables {
		mem.readerRef()
	}

	d.readState.Lock()
	old := d.readState.val
	d.readState.val = s
	d.readState.Unlock()
	if checker != nil {
		if err := checker(d); err != nil {
			d.opts.Logger.Fatalf("checker failed with error: %s", err)
		}
	}
	if old != nil {
		old.unrefLocked()
	}
}

type iterDebug struct {
	iter      *Iterator
	allocated time.Time
	stack     []byte
}

func (d *iterDebug) String() string {
	iterStr := "(closed)"
	if d.iter.iter != nil {
		iterStr = d.iter.iter.String()
	}
	return fmt.Sprintf(
		"iter: %s, age: %s (allocated @ %s), stack:\n%s\n",
		iterStr,
		time.Since(d.allocated),
		d.allocated,
		string(d.stack),
	)
}

func newIterDebug(i *Iterator) iterDebug {
	return iterDebug{
		iter:      i,
		allocated: time.Now(),
		stack:     debug.Stack(),
	}
}

func (s *readState) initDebug() {
	// Spawn a goroutine that reports on the current state of the open refs on the
	// read state.
	s.doneC = make(chan struct{})
	go func() {
		t := time.NewTicker(30 * time.Second)
		for {
			select {
			case <-s.doneC:
				return
			case <-t.C:
				fmt.Println(s.debug())
			}
		}
	}()
}

func (s *readState) addDebug(i *Iterator) {
	if s.doneC == nil {
		s.initDebug()
	}
	s.iterMap.Store(i, newIterDebug(i))
}

func (s *readState) removeDebug(i *Iterator) {
	s.iterMap.Delete(i)
}

func (s *readState) cleanupDebug() {
	s.iterMap.Range(func(key, _ interface{}) bool {
		s.iterMap.Delete(key)
		return true
	})
	if s.doneC != nil {
		close(s.doneC)
	}
}

func (s *readState) debug() string {
	var sb strings.Builder
	sb.WriteString("--- READ STATE DEBUG---\n")

	// Sort our debug info from oldest to largest.
	var dbgs []iterDebug
	s.iterMap.Range(func(_, value interface{}) bool {
		dbgs = append(dbgs, value.(iterDebug))
		return true
	})
	sort.Slice(dbgs, func(i, j int) bool {
		return dbgs[i].allocated.Before(dbgs[j].allocated)
	})
	for _, dbg := range dbgs {
		sb.WriteString(dbg.String())
	}

	return sb.String()
}
