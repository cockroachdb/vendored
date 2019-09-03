package tdigest

import (
	"fmt"
	"math"

	"github.com/ajwerner/tdigest/internal/scale"
)

type config struct {

	// compression factor controls the target maximum number of centroids in a
	// fully compressed TDigest.
	compression float64

	// bufferFactor is multiple of the compression size used to buffer unmerged
	// centroids. See config.bufferSize().
	bufferFactor int

	// scale controls how weight is apportioned to centroids.
	scale scale.Func

	useWeightLimit bool
}

func (cfg config) bufferSize() int {
	return int(math.Ceil(cfg.compression)) * (1 + cfg.bufferFactor)
}

// Option configures a TDigest.
type Option interface {
	apply(*config)
}

// BufferFactor configures the size of the buffer for uncompressed data.
// The default value is 5.
func BufferFactor(factor int) Option {
	return bufferFactorOption(factor)
}

func Compression(compression float64) Option {
	return compressionOption(compression)
}

func UseWeightLimit(useWeightLimit bool) Option {
	return weightLimitOption(useWeightLimit)
}

type bufferFactorOption int

func (o bufferFactorOption) apply(cfg *config) { cfg.bufferFactor = int(o) }
func (o bufferFactorOption) String() string {
	return fmt.Sprintf("bufferFactor=%d", o)
}

type compressionOption float64

func (o compressionOption) apply(cfg *config) { cfg.compression = float64(o) }
func (o compressionOption) String() string {
	return fmt.Sprintf("compression=%f", o)
}

type scaleOption struct{ scale.Func }

func (o scaleOption) apply(cfg *config) { cfg.scale = o.Func }
func (o scaleOption) String() string {
	return fmt.Sprintf("scale=%v", o.Func)
}

type weightLimitOption bool

func (o weightLimitOption) apply(cfg *config) { cfg.useWeightLimit = bool(o) }
func (o weightLimitOption) String() string {
	return fmt.Sprintf("weightLimit=%v", bool(o))
}

var defaultConfig = config{
	compression:    128,
	bufferFactor:   5,
	scale:          scaleOption{scale.K2{}},
	useWeightLimit: false,
}

type optionList []Option

func (l optionList) apply(cfg *config) {
	for _, o := range l {
		o.apply(cfg)
	}
}
