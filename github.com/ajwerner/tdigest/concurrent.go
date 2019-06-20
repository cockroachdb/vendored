// Package tdigest provides a concurrent, streaming quantiles estimation data
// structure for float64 data.
package tdigest

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
)

// Concurrent approximates a distribution of floating point numbers.
// All methods are safe to be called concurrently.
//
// Design
//
// The data structure is designed to maintain most of its state in a single
// slice
// The total in-memory size of a Concurrent is
//
//    (1+BufferFactor)*(int(Compression)+1)
//
// The data structure does not allocates memory after its construction.
type Concurrent struct {
	scale          scaleFunc
	compression    float64
	useWeightLimit bool
	// unmergedIdx is accessed with atomics.
	unmergedIdx int64
	mu          struct {
		sync.RWMutex
		sync.Cond

		// numMerged is the size of the prefix of centroid used for the merged
		// sorted data.
		numMerged int
	}

	centroids []centroid
}

// NewConcurrent creates a new Concurrent.
func NewConcurrent(options ...Option) *Concurrent {
	cfg := defaultConfig
	optionList(options).apply(&cfg)
	var td Concurrent
	td.mu.L = td.mu.RLocker()
	td.centroids = make(centroids, cfg.bufferSize())
	td.compression = cfg.compression
	td.scale = cfg.scale
	td.useWeightLimit = cfg.useWeightLimit
	return &td
}

// Read enables clients to perform a number of read operations on a snapshot
// of the data.
func (td *Concurrent) Read(f func(d Reader)) {
	td.compress()
	td.mu.RLock()
	defer td.mu.RUnlock()
	f((*readConcurrent)(td))
}

// Read enables clients to perform a number of read operations on a snapshot of
// the data without forcing a compression of buffered data.
func (td *Concurrent) ReadStale(f func(d Reader)) {
	td.mu.RLock()
	defer td.mu.RUnlock()
	f((*readConcurrent)(td))
}

func (td *Concurrent) String() (s string) {
	td.mu.RLock()
	defer td.mu.RUnlock()

	return fmt.Sprintf("(%.4f-[%.4f %.4f %.4f]-%.4f) numMerged: %v totalCount: %v",
		td.valueAtRLocked(0),
		td.valueAtRLocked(.25),
		td.valueAtRLocked(.5),
		td.valueAtRLocked(.75),
		td.valueAtRLocked(1),
		td.mu.numMerged,
		td.totalCountRLocked())
}

// TotalCount returns the total count that has been added to the Concurrent.
func (td *Concurrent) TotalCount() (c float64) {
	td.Read(func(r Reader) { c = r.TotalCount() })
	return c
}

// TotalSum returns the approximation of the weighted sum of all values
// recorded in to the sketch.
func (td *Concurrent) TotalSum() (s float64) {
	td.Read(func(r Reader) { s = r.TotalSum() })
	return s
}

// ValueAt returns the value of the quantile q.
// If q is not in [0, 1], ValueAt will panic.
// An empty Concurrent will return 0.
func (td *Concurrent) ValueAt(q float64) (v float64) {
	td.Read(func(r Reader) { v = r.ValueAt(q) })
	return v
}

// QuantileOf returns the estimated quantile at which this value falls in the
// distribution. If the v is smaller than any recorded value 0.0 will be
// returned and if v is larger than any recorded value 1.0 will be returned.
// An empty Concurrent will return 0.0 for all values.
func (td *Concurrent) QuantileOf(v float64) (q float64) {
	td.Read(func(r Reader) { q = r.QuantileOf(v) })
	return q
}

// Add adds the provided data to the Concurrent.
func (td *Concurrent) Add(mean, count float64) {
	td.mu.RLock()
	defer td.mu.RUnlock()
	td.centroids[td.getAddIndexRLocked()] = centroid{mean: mean, count: count}
}

// Record records adds a value with a count of 1.
func (td *Concurrent) Record(mean float64) { td.Add(mean, 1) }

// Decay decreases the weight of all tracked centroids by factor.
func (td *Concurrent) Decay(factor float64) {
	td.mu.Lock()
	defer td.mu.Unlock()
	td.compressLocked()
	decay(td.centroids[:td.mu.numMerged], factor)
}

// Merge combines other into td.
func (td *Concurrent) Merge(other *Concurrent) {
	td.mu.Lock()
	defer td.mu.Unlock()
	other.mu.Lock()
	defer other.mu.Unlock()
	other.compressLocked()
	totalCount := 0.0
	for i := range other.centroids {
		other.centroids[i].count -= totalCount
		totalCount += other.centroids[i].count
	}
	perm := rand.Perm(other.mu.numMerged)
	for _, i := range perm {
		td.centroids[td.getAddIndexLocked()] = other.centroids[i]
	}
}

func (td *Concurrent) getAddIndexRLocked() (r int) {
	for {
		idx := int(atomic.AddInt64(&td.unmergedIdx, 1))
		idx--
		if idx < len(td.centroids) {
			return idx
		} else if idx == len(td.centroids) {
			func() {
				td.mu.RUnlock()
				defer td.mu.RLock()
				td.compress()
				td.mu.Broadcast()
			}()
		} else {
			td.mu.Wait()
		}
	}
}

func (td *Concurrent) getAddIndexLocked() int {
	for {
		idx := int(atomic.AddInt64(&td.unmergedIdx, 1)) - 1
		if idx < len(td.centroids) {
			return idx
		}
		td.compressLocked()
	}
}

func (td *Concurrent) valueAtRLocked(q float64) float64 {
	if q < 0 || q > 1 {
		panic(fmt.Errorf("invalid quantile %v", q))
	}
	if td.mu.numMerged == 0 {
		return 0
	}
	return valueAt(td.centroids[:td.mu.numMerged], q)
}

func (td *Concurrent) quantileOfRLocked(v float64) float64 {
	if td.mu.numMerged == 0 {
		return 0
	}
	return quantileOf(td.centroids[:td.mu.numMerged], v)
}

func (td *Concurrent) totalCountRLocked() float64 {
	return totalCount(td.centroids[:td.mu.numMerged])
}

func (td *Concurrent) totalSumRLocked() float64 {
	return totalSum(td.centroids[:td.mu.numMerged])
}

func (td *Concurrent) compress() {
	td.mu.Lock()
	defer td.mu.Unlock()
	td.compressLocked()
}

func (td *Concurrent) compressLocked() {
	idx := int(atomic.LoadInt64(&td.unmergedIdx))
	if idx > len(td.centroids) {
		idx = len(td.centroids)
	}
	td.mu.numMerged = compress(td.centroids[:idx], td.compression, td.scale, td.mu.numMerged, td.useWeightLimit)
	atomic.StoreInt64(&td.unmergedIdx, int64(td.mu.numMerged))
}

type readConcurrent Concurrent

var _ Reader = (*readConcurrent)(nil)

func (rtd *readConcurrent) ValueAt(q float64) (v float64) {
	return (*Concurrent)(rtd).valueAtRLocked(q)
}

func (rtd *readConcurrent) QuantileOf(v float64) (q float64) {
	return (*Concurrent)(rtd).quantileOfRLocked(v)
}

func (rtd *readConcurrent) TotalCount() (c float64) {
	return (*Concurrent)(rtd).totalCountRLocked()
}

func (rtd *readConcurrent) TotalSum() float64 {
	return (*Concurrent)(rtd).totalSumRLocked()
}
