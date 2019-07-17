package windowed

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/ajwerner/tdigest"
)

// Windowed is a TDigest that provides a single write path but
// provides mechanisms to read at various different timescales.The
// Windowed structure uses hierarchical windows to provide high
// precision estimates at approximate timescales while requiring
// sublinear number of buckets.
//
// The structure of the hiearchy is configurable to control
// size-accuracy tradeoffs for a given trailing trailing timescale.
//
// TODO(ajwerner): Figure out how clients can ergonomically express their timescale configuration.
//
// Imagine we wanted to know about the trailing 1s and 1m. One thing
// that we could do is keep a ring-buffer of the last 1m of tdigests
// in 1s intervals such that there are 59 "closed" digests which only
// contain a "merged" buffer and no write buffer as well as two "read"
// digests and an "open" digest. This will contain a single concurrent
// write buffer that gets merged into each of the "read" digests as
// well as the "open" digest.
//
// This already offers some issues, it is not actually the trailing 1s
// and 1m you'll be reading but rather it's the last 1-2s depending on
// when the last tick occurred and similarly for the minute it's the
// last 1m-1m1s. For the minute this certainly isn't a problem and
// even for the trailing 1s this is probably okay (furthermore, we're
// probably not interested in the trailing 1s, instead we're more
// likely to be interested in the last ~10s.
//
// We can extend this tradeoff further by allowing further window
// size. For example, imagine we keep this 1s buffer as described
// above, but then we also keep a next layer which represent 2s
// intervals, then we only need to keep 4 of them to get to a trailing
// 10s buffer with a 2s window size. We can then layer these things up
// to use 5 more to get the trailing 1m with a 10s window size. This
// allows us to get a reasonable window size on a 5 minute trailing
// period of 1m using just 15 tdigests instead of the 300 we'd need if
// we kept all 1s ring buffers.
//
//
//______________________________________________________________________________
//
// In the beginning all of the trailing levels start at the front of their cycles.
//
//      0  |  |  |  |  V  |  |  |  |  X  |  |  |  |  XV |  |  |  |  XX |  |  |  | XXV |  |  |  |XXX
//  o|0s]
//  0|  ( 1]
//  1|  (    2|    4|    6|    8]
//  2|  (                           10]                           20]                          30]
//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//  2|  (10]20]30]40]
//  3|  (            50]           100]           150]           200]           250]
//      0  |  |  |  |  L  |  |  |  |  C  |  |  |  |  CL |  |  |  |  CC |  |  |  | CCL |  |  |  |CCC
//
//_______________________________________________________________________________
//
// Over time buckets slide over until tick events happen.
//
//      0  |  |  |  |  V  |  |  |  |  X  |  |  |  |  XV |  |  |  |  XX |  |  |  | XXV |  |  |  |XXX
//  0| -1s]<
//  1|    (2s]
//           (  <4s]  <6s]  <8s] <10s]<
//  2|             (                         <14s]                        <24s]<
//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//  2|    (14]24]34]44]54]<
//  3|                   (              ]              ]              ]              ]             ]<
//      0  |  |  |  |  L  |  |  |  |  C  |  |  |  |  CL |  |  |  |  CC |  |  |  | CCL |  |  |  |CCC
//_______________________________________________________________________________
//
//
//      (1s]
//         (2s]
//            (   4s]   6s]   8s]  10s]
//                                    (                          20s]                          30s]
//      ___________________________________________________________________________________________
//      0  |  |  |  |  V  |  |  |  |  X  |  |  |  |  V  |  |  |  |  D  |  |  |  |  V  |  |  |  |  D
//
//
//  1|  (   2s]   4s]   6s]   8s]  10s]
//  2|                                (                          20s]                         30s]
//  3|..___________________________________________________________________________________________
//      0  |  |  |  |  V  |  |  |  |  X  |  |  |  |  XV |  |  |  |  XX |  |  |  | XXV |  |  |  |XXX
//
//
//
//  We only keep these around so that we can read at this level with a 2s window
//
//  1|  (   2s]   4s]   6s]   8s]
//  2|  (                             ]                          20s]
//  3|..___________________________________________________________________________________________
//      0  |  |  |  |  V  |  |  |  |  X  |  |  |  |  XV |  |  |  |  XX |  |  |  | XXV |  |  |  |XXX
//
//  2|  (20]  ]  ]  ]
//  3|  (           ]              ]              ]              ]
//      ___________________________________________________________________________________________
//      0  |  |  |  |  L  |  |  |  |  C  |  |  |  |  CL |  |  |  |  CC |  |  |  | CCL |  |  |  |CCC
//
//-------------------------------------------------------------------------------
//
//  1|  (   2s]   4s]   6s]   8s]
//  2|  (                             ]                          20s]
//  3|..___________________________________________________________________________________________
//      0  |  |  |  |  V  |  |  |  |  X  |  |  |  |  XV |  |  |  |  XX |  |  |  | XXV |  |  |  |XXX
//
//  2|  (20]  ]  ]  ]
//  3|  (           ]              ]              ]              ]
//      ___________________________________________________________________________________________
//      0  |  |  |  |  L  |  |  |  |  C  |  |  |  |  CL |  |  |  |  CC |  |  |  | CCL |  |  |  |CCC
//
//-------------------------------------------------------------------------------
type TDigest struct {
	tickInterval time.Duration

	mu struct {
		sync.RWMutex
		spare    *tdigest.TDigest
		lastTick time.Time
		ticks    int

		// We need to have levels
		open   *tdigest.Concurrent
		levels []level
	}

	// We now want some number of levels where each level has a tick period
	// a next tick var (or last, same difference). It also has a slice of
	// tdigest structs.
	//
	// Then, upon each tick, we add the open to what we need to and then
	// we tick the other levels as needed.

	// this is hand-wavy but maybe will work?

}

type digestRingBuf struct {
	head    int32
	len     int32
	digests []*tdigest.TDigest
}

func (rb *digestRingBuf) back() *tdigest.TDigest {
	return rb.at(rb.len - 1)
}

func (rb *digestRingBuf) at(idx int32) *tdigest.TDigest {
	return rb.digests[(rb.head+idx)%int32(len(rb.digests))]
}

func (rb *digestRingBuf) pushFront(td *tdigest.TDigest) {
	if rb.full() {
		panic("cannot push onto a full digest")
	}
	if rb.head--; rb.head < 0 {
		rb.head += int32(len(rb.digests))
	}
	rb.digests[rb.head] = td
	rb.len++
}

func (rb *digestRingBuf) popBack() *tdigest.TDigest {
	ret := rb.back()
	rb.len--
	return ret
}

func (rb *digestRingBuf) full() bool {
	return rb.len == int32(len(rb.digests))
}

func (rb *digestRingBuf) forEach(f func(int32, *tdigest.TDigest)) {
	for i := int32(0); i < rb.len; i++ {
		f(i, rb.at(i))
	}
}

type level struct {
	period int
	digestRingBuf
}

const size = 128

// NewTDigest returns a new TDigest TDigest.
func NewTDigest() *TDigest {
	// TODO(ajwerner): add configuration.
	w := &TDigest{
		tickInterval: time.Second,
	}
	// TODO(ajwerner): fix buf where the last one has only 5 buckets
	w.mu.levels = []level{
		{
			// (0-1)-(1-2)s
			period:        1,
			digestRingBuf: makeDigests(1, size),
		},
		{
			// (0-2)-(2-4), (0-2)-(4-6), (0-2)-(6-8), (0-2)-(8-10)s
			period:        2,
			digestRingBuf: makeDigests(4, size),
		},
		{
			period:        10,
			digestRingBuf: makeDigests(5, size),
		},
		{
			period:        60,
			digestRingBuf: makeDigests(3, size),
		},
		{
			period:        120,
			digestRingBuf: makeDigests(10, size),
		},
	}
	w.mu.open = tdigest.NewConcurrent(tdigest.Compression(size), tdigest.BufferFactor(10))
	w.mu.spare = tdigest.New(tdigest.Compression(size), tdigest.BufferFactor(2))
	return w
}

func makeDigests(n int, size int) digestRingBuf {
	ret := digestRingBuf{
		digests: make([]*tdigest.TDigest, 0, n),
	}
	for i := 0; i < n; i++ {
		ret.digests = append(ret.digests,
			tdigest.New(tdigest.Compression(float64(size)), tdigest.BufferFactor(1)))
		ret.len++
	}
	return ret

}

func (w *TDigest) AddAt(t time.Time, mean, count float64) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	if t.Sub(w.mu.lastTick) > w.tickInterval {
		w.tickAtRLocked(t)
	}
	w.mu.open.Add(mean, count)
}

func (w *TDigest) tickAtRLocked(t time.Time) {
	w.mu.RUnlock()
	defer w.mu.RLock()
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.mu.lastTick.IsZero() {
		w.mu.lastTick = t
		return
	}
	ticksNeeded := int(t.Sub(w.mu.lastTick) / w.tickInterval)
	// TODO(ajwerner): optimize when many ticks are needed.
	if ticksNeeded <= 0 {
		return
	}
	for i := 0; i < ticksNeeded; i++ {
		w.tickLocked()
	}
}

func (w *TDigest) tickLocked() {
	// A tick means moving the current open interval down to the next level
	// It may also mean merging all of the current bottom level into a new
	// digest for the next level which may need to happen recursively.
	w.mu.ticks++
	w.mu.lastTick = w.mu.lastTick.Add(w.tickInterval)
	// Take the merged buf from the top and write it into the "spare"
	closed := w.mu.spare
	w.mu.spare = nil
	w.mu.open.AddTo(closed)
	w.mu.open.Clear()
	for i := range w.mu.levels {
		l := &w.mu.levels[i]
		tail := l.popBack()
		l.forEach(func(_ int32, td *tdigest.TDigest) {
			closed.AddTo(td)
		})
		l.pushFront(closed)
		var next *level
		if i+1 < len(w.mu.levels) {
			next = &w.mu.levels[i+1]
		}
		tickNext := next != nil && w.mu.ticks%next.period == 0
		if tickNext {
			closed.AddTo(tail)
		}
		closed = tail
		if !tickNext {
			break
		}
	}
	closed.Clear()
	w.mu.spare = closed
}

func (w *TDigest) String() string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.stringRLocked(w.mu.lastTick)
}

func (w *TDigest) stringRLocked(now time.Time) string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "Windowed{lastTick: %v\n", w.mu.lastTick)
	curDur := now.Sub(w.mu.lastTick)
	tc := w.mu.open.TotalCount()
	fmt.Fprintf(&buf, "\tnow (0-%v): %v %v\n", now.Sub(w.mu.lastTick), w.mu.open.String(), tc)
	for i := range w.mu.levels {
		l := &w.mu.levels[i]
		offset := curDur + w.tickInterval*time.Duration(w.mu.ticks%(l.period))
		l.forEach(func(j int32, td *tdigest.TDigest) {
			dur := offset + w.tickInterval*time.Duration(j+1)*time.Duration(l.period)
			fmt.Fprintf(&buf, "\t%v,%v: (%v-%v): %v\n", i, j, offset, dur, td)
		})
	}
	return buf.String()
}

// Reader enables reading from a TDigest.
//
// Reader allocates memory lazily but will do so only once if it continues
// to see the same *TDigest.
//
// Reader is safe for concurrent use.
type Reader struct {
	mu struct {
		sync.Mutex
		w   *TDigest
		buf *tdigest.TDigest
	}
}

func (r *Reader) KSScore(
	a, b time.Duration, w *TDigest, steps int, other *Reader,
) (deltaMax, ks float64) {
	step := 1.0 / float64(steps)
	r.Read(a, w, func(_ time.Duration, ar tdigest.Reader) {
		other.Read(b, w, func(_ time.Duration, br tdigest.Reader) {
			q := 0.0
			for i := 0; i <= steps; i++ {
				av := ar.ValueAt(q)
				bq := br.QuantileOf(av)
				if math.Abs(q-bq) > math.Abs(deltaMax) {
					deltaMax = q - bq
				}
				if q += step; q > 1 { // deal with float precsion.
					q = 1
				}
			}
			ac := ar.TotalCount()
			bc := br.TotalCount()
			ks = math.Sqrt((ac + bc) / (ac * bc))
		})
	})
	return deltaMax, ks
}

func (r *Reader) Read(
	trailing time.Duration, w *TDigest, f func(last time.Duration, r tdigest.Reader),
) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.mu.w != w {
		r.mu.w = w
		r.mu.buf = tdigest.New(
			tdigest.Compression(size),
			tdigest.BufferFactor(len(w.mu.levels)),
		)
	}

	r.mu.buf.Clear()
	w.mu.RLock()
	// we want to find the bucket which contains this time window
	// and then fill it in below.
	// TODO(ajwerner): optimize the reader behavior to just accumulate
	// the indexes and do a single merge pass into the buffer.

	// First we work our way up the levels until we find the
	// level that contains this time period
	//
	// Then we work our way down to fill in the remainder
	curDur := time.Duration(0) // now.Sub(w.mu.lastTick)
	var i int
	var last time.Duration
	for i = 0; i < len(w.mu.levels); i++ {
		l := &w.mu.levels[i]
		bucketDur := w.tickInterval * time.Duration(l.period)
		levelDur := bucketDur * time.Duration(len(l.digests)+1)
		if levelDur <= trailing && i+1 < len(w.mu.levels) {
			continue
		}
		offset := w.tickInterval * time.Duration(w.mu.ticks%l.period)
		idx := int((trailing - offset) / bucketDur)
		if (trailing-offset)%bucketDur == 0 {
			idx--
		}
		trailing = offset
		l.at(int32(idx)).AddTo(r.mu.buf)
		last = offset + time.Duration(idx+1)*bucketDur
		break
	}
	for ; i >= 0 && trailing > curDur; i-- {
		l := &w.mu.levels[i]
		bucketDur := w.tickInterval * time.Duration(l.period)
		idx := int(trailing / bucketDur)
		if idx == len(l.digests) {
			idx--
		}
		trailing = curDur + w.tickInterval*time.Duration(w.mu.ticks%(l.period))
		l.at(int32(idx)).AddTo(r.mu.buf)
	}
	w.mu.open.Read(func(d tdigest.Reader) {
		d.AddTo(r.mu.buf)
	})
	w.mu.RUnlock()
	f(last, r.mu.buf)
}
