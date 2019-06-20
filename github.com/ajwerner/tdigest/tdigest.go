// Package tdigest provides a concurrent, streaming quantiles estimation data
// structure for float64 data.
package tdigest

import (
	"math"
	"sort"
)

// Sketch is an
type Sketch interface {
	Reader
	Add(mean, count float64)
}

// Reader provides read access to a float64 valued distribution by
// quantile or by value.
type Reader interface {
	TotalCount() float64
	TotalSum() float64
	ValueAt(q float64) (v float64)
	QuantileOf(v float64) (q float64)
}

type centroid struct {
	mean, count float64
}

type TDigest struct {
	scale          scaleFunc
	compression    float64
	useWeightLimit bool

	centroids   []centroid
	numMerged   int
	unmergedIdx int
}

// New creates a new Concurrent.
func New(options ...Option) *TDigest {
	cfg := defaultConfig
	optionList(options).apply(&cfg)
	var td TDigest
	td.centroids = make([]centroid, cfg.bufferSize())
	td.compression = cfg.compression
	td.scale = cfg.scale
	td.useWeightLimit = cfg.useWeightLimit
	return &td
}

func (td *TDigest) ValueAt(q float64) (v float64) {
	td.compress()
	return valueAt(td.centroids[:td.numMerged], v)
}

// QuantileOf returns the estimated quantile at which this value falls in the
// distribution. If the v is smaller than any recorded value 0.0 will be
// returned and if v is larger than any recorded value 1.0 will be returned.
// An empty Concurrent will return 0.0 for all values.
func (td *TDigest) QuantileOf(v float64) (q float64) {
	td.compress()
	return quantileOf(td.centroids[:td.numMerged], v)
}

func (td *TDigest) TotalCount() (c float64) {
	td.compress()
	return totalCount(td.centroids[:td.numMerged])
}

func (td *TDigest) Add(mean, count float64) {
	if td.unmergedIdx == len(td.centroids) {
		td.compress()
	}
	td.centroids[td.unmergedIdx] = centroid{mean: mean, count: count}
	td.unmergedIdx++
}

func (td *TDigest) TotalSum() float64 {
	td.compress()
	return totalSum(td.centroids[:td.numMerged])
}

func (td *TDigest) compress() {
	td.numMerged = compress(td.centroids[:td.unmergedIdx], td.compression, td.scale, td.numMerged, td.useWeightLimit)
	td.unmergedIdx = td.numMerged
}

func (td *TDigest) Record(mean float64) { td.Add(mean, 1) }

func valueAt(merged []centroid, q float64) float64 {
	goal := q * merged[len(merged)-1].count
	i := sort.Search(len(merged), func(i int) bool {
		return merged[i].count >= goal
	})
	n := merged[i]
	k := 0.0
	if i > 0 {
		k = merged[i-1].count
		n.count -= merged[i-1].count
	}
	deltaK := goal - k - (n.count / 2)
	right := deltaK > 0

	// if before the first point or after the last point, return the current mean.
	if !right && i == 0 || right && (i+1) == len(merged) {
		return n.mean
	}
	var nl, nr centroid
	if right {
		nl = n
		nr = merged[i+1]
		nr.count -= merged[i].count
		k += nl.count / 2
	} else {
		nl = merged[i-1]
		if i > 1 {
			nl.count -= merged[i-2].count
		}
		nr = n
		k -= nr.count / 2
	}
	x := goal - k
	m := (nr.mean - nl.mean) / ((nl.count / 2) + (nr.count / 2))
	return m*x + nl.mean
}

func quantileOf(merged []centroid, v float64) float64 {
	i := sort.Search(len(merged), func(i int) bool {
		return merged[i].mean >= v
	})
	// Deal with the ends of the distribution.
	switch {
	case i == 0:
		return 0
	case i == len(merged):
		return 1
	case i+1 == len(merged) && v >= merged[i].mean:
		return 1
	}
	k := merged[i-1].count
	nr := merged[i]
	nr.count -= k
	nl := merged[i-1]
	if i > 1 {
		nl.count -= merged[i-2].count
	}
	delta := (nr.mean - nl.mean)
	cost := ((nl.count / 2) + (nr.count / 2))
	m := delta / cost
	return (k + ((v - nl.mean) / m)) / merged[len(merged)-1].count
}

func totalCount(merged []centroid) float64 {
	if len(merged) == 0 {
		return 0.0
	}
	return merged[len(merged)-1].count
}

func totalSum(merged []centroid) float64 {
	var countSoFar float64
	var sum float64
	for i := range merged {
		sum += (merged[i].count - countSoFar) * merged[i].mean
		countSoFar += merged[i].count
	}
	return sum
}

func compress(
	cl []centroid, compression float64, scale scaleFunc, numMerged int, useWeightLimit bool,
) (newNumMerged int) {
	if len(cl) == 0 {
		return 0
	}
	if numMerged == len(cl) {
		return numMerged
	}
	totalCount := 0.0
	for i := 0; i < numMerged; i++ {
		cl[i].count -= totalCount
		totalCount += cl[i].count
	}
	for i := numMerged; i < len(cl); i++ {
		totalCount += cl[i].count
	}
	sort.Sort(centroids(cl))
	normalizer := scale.normalizer(compression, totalCount)
	cur := 0
	countSoFar := 0.0
	var k1, scaleLimit float64 // for use with scale scaleLimitit
	if !useWeightLimit {
		k1 = scale.k(0, normalizer)
		scaleLimit = totalCount * scale.q(k1+1, normalizer)
	}
	for i := 1; i < len(cl); i++ {
		proposedCount := cl[cur].count + cl[i].count
		var shouldAdd bool
		if useWeightLimit {
			q0 := countSoFar / totalCount
			q2 := (countSoFar + proposedCount) / totalCount
			probDensity := math.Min(scale.max(q0, normalizer), scale.max(q2, normalizer))
			limit := totalCount * probDensity
			shouldAdd = proposedCount < limit
		} else /* useScaleLimit */ {
			shouldAdd = countSoFar+proposedCount <= scaleLimit

		}
		if shouldAdd {
			cl[cur].count += cl[i].count
			delta := cl[i].mean - cl[cur].mean
			if delta > 0 {
				weightedDelta := (delta * cl[i].count) / cl[cur].count
				cl[cur].mean += weightedDelta
			}
		} else {
			countSoFar += cl[cur].count
			cl[cur].count = countSoFar
			cur++
			cl[cur] = cl[i]
			if !useWeightLimit {
				k1 = scale.k(countSoFar/totalCount, normalizer)
				scaleLimit = totalCount * scale.q(k1+1, normalizer)
			}
		}
		if cur != i {
			cl[i] = centroid{}
		}
	}
	cl[cur].count += countSoFar
	return cur + 1
}

func decay(merged []centroid, factor float64) {
	const verySmall = 1e-9
	for i := range merged {
		if count := merged[i].count * factor; count > verySmall {
			merged[i].count = count
		}
	}
}

type centroids []centroid

var _ sort.Interface = centroids(nil)

func (c centroids) Len() int           { return len(c) }
func (c centroids) Less(i, j int) bool { return c[i].mean < c[j].mean }
func (c centroids) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
