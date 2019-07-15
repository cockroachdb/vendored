package tdigest

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/ajwerner/tdigest/internal/scale"
	"github.com/ajwerner/tdigest/internal/tdigest"
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
	scale          scale.Func
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

	centroids []tdigest.Centroid
}

// NewConcurrent creates a new Concurrent.
func NewConcurrent(options ...Option) *Concurrent {
	cfg := defaultConfig
	optionList(options).apply(&cfg)
	var td Concurrent
	td.mu.L = td.mu.RLocker()
	td.centroids = make([]tdigest.Centroid, cfg.bufferSize())
	td.compression = cfg.compression
	td.scale = cfg.scale
	td.useWeightLimit = cfg.useWeightLimit
	return &td
}

// Clear resets the data structure, clearing all recorded data.
func (td *Concurrent) Clear() {
	td.mu.Lock()
	defer td.mu.Unlock()
	atomic.StoreInt64(&td.unmergedIdx, 0)
	td.mu.numMerged = 0
}

// Read enables clients to perform a number of read operations on a snapshot
// of the data.
func (td *Concurrent) Read(f func(d Reader)) {
	td.compress()
	td.mu.RLock()
	defer td.mu.RUnlock()
	f((*readConcurrent)(td))
}

// ReadStale enables clients to perform a number of read operations on a
// snapshot of the data without forcing a compression of buffered data.
func (td *Concurrent) ReadStale(f func(d Reader)) {
	td.mu.RLock()
	defer td.mu.RUnlock()
	f((*readConcurrent)(td))
}

func (td *Concurrent) String() (s string) {
	td.ReadStale(func(r Reader) { s = readerString(r) })
	return s
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

// InnerMean returns the mean of the inner quantile range.
// It requires flushing the buffer then is an O(n) operation on the number
// of centroids.
func (td *Concurrent) InnerMean(inner float64) (c float64) {
	td.Read(func(r Reader) { c = r.InnerMean(inner) })
	return c
}

// TrimmedMean returns the mean of the inner quantile range from lo to hi.
// It requires flushing the buffer then is an O(n) operation.
func (td *Concurrent) TrimmedMean(lo, hi float64) (c float64) {
	td.Read(func(r Reader) { c = r.TrimmedMean(lo, hi) })
	return c
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
	td.centroids[td.getAddIndexRLocked()] = tdigest.Centroid{Mean: mean, Count: count}
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

func (td *Concurrent) addToRLocked(into Recorder) {
	addTo(into, td.centroids[:td.mu.numMerged])
}

// AddTo adds the currently recorded data into the provided Recorder.
func (td *Concurrent) AddTo(into Recorder) {
	td.Read(func(r Reader) { r.AddTo(into) })
}

func (td *Concurrent) getAddIndexRLocked() (r int) {
	compress := func() {
		td.mu.RUnlock()
		defer td.mu.RLock()
		td.compress()
		td.mu.Broadcast()
	}
	for {
		idx := int(atomic.AddInt64(&td.unmergedIdx, 1))
		idx--
		if idx < len(td.centroids) {
			return idx
		} else if idx == len(td.centroids) {
			compress()
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
	return tdigest.ValueAt(td.centroids[:td.mu.numMerged], q)
}

func (td *Concurrent) quantileOfRLocked(v float64) float64 {
	return tdigest.QuantileOf(td.centroids[:td.mu.numMerged], v)
}

func (td *Concurrent) totalCountRLocked() float64 {
	return tdigest.TotalCount(td.centroids[:td.mu.numMerged])
}

func (td *Concurrent) totalSumRLocked() float64 {
	return tdigest.TotalSum(td.centroids[:td.mu.numMerged])
}

func (td *Concurrent) innerMeanRLocked(inner float64) float64 {
	tails := (1 - inner) / 2
	return tdigest.TrimmedMean(td.centroids[:td.mu.numMerged], tails, 1-tails)
}

func (td *Concurrent) trimmedMeanRLocked(lo, hi float64) float64 {
	return tdigest.TrimmedMean(td.centroids[:td.mu.numMerged], lo, hi)
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
	td.mu.numMerged = tdigest.Compress(td.centroids[:idx], td.compression, td.scale, td.mu.numMerged, td.useWeightLimit)
	atomic.StoreInt64(&td.unmergedIdx, int64(td.mu.numMerged))
}

type readConcurrent Concurrent

var _ Reader = (*readConcurrent)(nil)

func (rtd *readConcurrent) AddTo(into Recorder) {
	(*Concurrent)(rtd).addToRLocked(into)
}

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

func (rtd *readConcurrent) InnerMean(inner float64) float64 {
	return (*Concurrent)(rtd).innerMeanRLocked(inner)
}

func (rtd *readConcurrent) TrimmedMean(lo, hi float64) float64 {
	return (*Concurrent)(rtd).trimmedMeanRLocked(lo, hi)
}
