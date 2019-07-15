// package scale provides functionality to control the scaling of tdigests
package scale

import "math"

type Func interface {
	Normalizer(compression, totalCount float64) float64
	K(q, normalizer float64) float64
	Q(k, normalizer float64) float64
	Max(q, normalizer float64) float64
	String() string
}

type K2 struct{}

func (f K2) String() string { return "k2" }

func (f K2) Normalizer(compression, totalCount float64) float64 {
	return compression / f.z(compression, totalCount)
}

func (f K2) z(compression, totalCount float64) float64 {
	return 4*math.Log(totalCount/compression) + 24
}

func (f K2) K(q, normalizer float64) float64 {
	if q < 1e-15 {
		return 2 * f.K(1e-15, normalizer)
	} else if q > (1 - 1e-15) {
		return 2 * f.K(1-1e-15, normalizer)
	}
	return math.Log(q/(1-q)) * normalizer
}

func (f K2) Q(k, normalizer float64) (ret float64) {
	w := math.Exp(k / normalizer)
	return w / (1 + w)
}

func (f K2) Max(q, normalizer float64) float64 {
	return q * (1 - q) / normalizer
}
