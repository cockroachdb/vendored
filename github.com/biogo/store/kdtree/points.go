// Copyright ©2012 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kdtree

import (
	"math"
)

var (
	_ Interface  = Points{}
	_ Comparable = Point{}
)

// Randoms is the maximum number of random values to sample for calculation of median of
// random elements.
var Randoms = 100

// A Point represents a point in a k-d space that satisfies the Comparable interface.
type Point []float64

func (p Point) Compare(c Comparable, d Dim) float64 { q := c.(Point); return p[d] - q[d] }
func (p Point) Dims() int                           { return len(p) }
func (p Point) Distance(c Comparable) float64 {
	q := c.(Point)
	var sum float64
	for dim, c := range p {
		d := c - q[dim]
		sum += d * d
	}
	return sum
}
func (p Point) Extend(b *Bounding) *Bounding {
	if b == nil {
		b = &Bounding{append(Point(nil), p...), append(Point(nil), p...)}
	}
	min := b[0].(Point)
	max := b[1].(Point)
	for d, v := range p {
		min[d] = math.Min(min[d], v)
		max[d] = math.Max(max[d], v)
	}
	*b = Bounding{min, max}
	return b
}

// A Points is a collection of point values that satisfies the Interface.
type Points []Point

func (p Points) Bounds() *Bounding {
	if p.Len() == 0 {
		return nil
	}
	min := append(Point(nil), p[0]...)
	max := append(Point(nil), p[0]...)
	for _, e := range p[1:] {
		for d, v := range e {
			min[d] = math.Min(min[d], v)
			max[d] = math.Max(max[d], v)
		}
	}
	return &Bounding{min, max}
}
func (p Points) Index(i int) Comparable         { return p[i] }
func (p Points) Len() int                       { return len(p) }
func (p Points) Pivot(d Dim) int                { return Plane{Points: p, Dim: d}.Pivot() }
func (p Points) Slice(start, end int) Interface { return p[start:end] }

// A Plane is a wrapping type that allows a Points type be pivoted on a dimension.
type Plane struct {
	Dim
	Points
}

func (p Plane) Less(i, j int) bool              { return p.Points[i][p.Dim] < p.Points[j][p.Dim] }
func (p Plane) Pivot() int                      { return Partition(p, MedianOfRandoms(p, Randoms)) }
func (p Plane) Slice(start, end int) SortSlicer { p.Points = p.Points[start:end]; return p }
func (p Plane) Swap(i, j int) {
	p.Points[i], p.Points[j] = p.Points[j], p.Points[i]
}
