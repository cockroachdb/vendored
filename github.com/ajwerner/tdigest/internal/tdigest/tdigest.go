// package tdigest contains the low-level algorithms to compress and query
// lists of centroids.
package tdigest

import (
	"math"
	"sort"

	"github.com/ajwerner/tdigest/internal/scale"
)

// Centroid represents a point in distribution.
type Centroid struct {
	Mean, Count float64
}

func TrimmedMean(merged []Centroid, lo, hi float64) float64 {
	if len(merged) == 0 {
		return 0
	}
	totalCount := merged[len(merged)-1].Count
	leftTailCount := lo * totalCount
	rightTailCount := hi * totalCount
	var countSeen float64
	var weightedMean float64
	for i, c := range merged {
		if i > 0 {
			countSeen = merged[i-1].Count
		}
		if c.Count < leftTailCount {
			continue
		}
		if countSeen > rightTailCount {
			break
		}
		if countSeen < leftTailCount {
			countSeen = leftTailCount
		}
		if c.Count > rightTailCount {
			c.Count = rightTailCount
		}
		weightedMean += c.Mean * (c.Count - countSeen)
	}
	includedCount := totalCount * (hi - lo)
	return weightedMean / includedCount
}

func ValueAt(merged []Centroid, q float64) float64 {
	if len(merged) == 0 {
		return 0
	}
	goal := q * merged[len(merged)-1].Count
	i := sort.Search(len(merged), func(i int) bool {
		return merged[i].Count >= goal
	})
	n := merged[i]
	k := 0.0
	if i > 0 {
		k = merged[i-1].Count
		n.Count -= merged[i-1].Count
	}
	deltaK := goal - k - (n.Count / 2)
	right := deltaK > 0

	// if before the first point or after the last point, return the current mean.
	if !right && i == 0 || right && (i+1) == len(merged) {
		return n.Mean
	}
	var nl, nr Centroid
	if right {
		nl = n
		nr = merged[i+1]
		nr.Count -= merged[i].Count
		k += nl.Count / 2
	} else {
		nl = merged[i-1]
		if i > 1 {
			nl.Count -= merged[i-2].Count
		}
		nr = n
		k -= nr.Count / 2
	}
	x := goal - k
	m := (nr.Mean - nl.Mean) / ((nl.Count / 2) + (nr.Count / 2))
	return m*x + nl.Mean
}

func QuantileOf(merged []Centroid, v float64) float64 {
	i := sort.Search(len(merged), func(i int) bool {
		return merged[i].Mean >= v
	})
	// Deal with the ends of the distribution.
	switch {
	case i == 0:
		return 0
	case i == len(merged):
		return 1
	case i+1 == len(merged) && v >= merged[i].Mean:
		return 1
	}
	k := merged[i-1].Count
	nr := merged[i]
	nr.Count -= k
	nl := merged[i-1]
	if i > 1 {
		nl.Count -= merged[i-2].Count
	}
	k -= nl.Count / 2
	delta := (nr.Mean - nl.Mean)
	cost := ((nl.Count / 2) + (nr.Count / 2))
	m := delta / cost
	return (k + ((v - nl.Mean) / m)) / merged[len(merged)-1].Count
}

func TotalCount(merged []Centroid) float64 {
	if len(merged) == 0 {
		return 0.0
	}
	return merged[len(merged)-1].Count
}

func TotalSum(merged []Centroid) float64 {
	var countSoFar float64
	var sum float64
	for i := range merged {
		sum += (merged[i].Count - countSoFar) * merged[i].Mean
		countSoFar = merged[i].Count
	}
	return sum
}

func Compress(
	cl []Centroid, compression float64, scale scale.Func, numMerged int, useWeightLimit bool,
) (newNumMerged int) {
	if len(cl) == 0 {
		return 0
	}
	if numMerged == len(cl) {
		return numMerged
	}
	totalCount := 0.0
	for i := 0; i < numMerged; i++ {
		cl[i].Count -= totalCount
		totalCount += cl[i].Count
	}
	for i := numMerged; i < len(cl); i++ {
		totalCount += cl[i].Count
	}
	sort.Sort(centroids(cl))
	normalizer := scale.Normalizer(compression, totalCount)
	cur := 0
	countSoFar := 0.0
	var k1, scaleLimit float64 // for use with scale scaleLimitit
	if !useWeightLimit {
		k1 = scale.K(0, normalizer)
		scaleLimit = totalCount * scale.Q(k1+1, normalizer)
	}
	for i := 1; i < len(cl); i++ {
		proposedCount := cl[cur].Count + cl[i].Count
		var shouldAdd bool
		if useWeightLimit {
			q0 := countSoFar / totalCount
			q2 := (countSoFar + proposedCount) / totalCount
			probDensity := math.Min(scale.Max(q0, normalizer), scale.Max(q2, normalizer))
			limit := totalCount * probDensity
			shouldAdd = proposedCount < limit
		} else /* useScaleLimit */ {
			shouldAdd = countSoFar+proposedCount <= scaleLimit

		}
		if shouldAdd {
			cl[cur].Count += cl[i].Count
			delta := cl[i].Mean - cl[cur].Mean
			if delta > 0 {
				weightedDelta := (delta * cl[i].Count) / cl[cur].Count
				cl[cur].Mean += weightedDelta
			}
		} else {
			countSoFar += cl[cur].Count
			cl[cur].Count = countSoFar
			cur++
			cl[cur] = cl[i]
			if !useWeightLimit {
				k1 = scale.K(countSoFar/totalCount, normalizer)
				scaleLimit = totalCount * scale.Q(k1+1, normalizer)
			}
		}
		if cur != i {
			cl[i] = Centroid{}
		}
	}
	cl[cur].Count += countSoFar
	return cur + 1
}

type centroids []Centroid

var _ sort.Interface = centroids(nil)

func (c centroids) Len() int           { return len(c) }
func (c centroids) Less(i, j int) bool { return c[i].Mean < c[j].Mean }
func (c centroids) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
