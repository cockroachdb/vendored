package ilog10

import "math/bits"

const n = ^uint64(0)

var lookup64 = [64]uint64{
	// This initializer list is easier to read as follows:
	//     10000000000000000000, n, n, n, 1000000000000000000, n, n, 100000000000000000, n, n,
	//        10000000000000000, n, n, n,    1000000000000000, n, n,    100000000000000, n, n,
	//           10000000000000, n, n, n,       1000000000000, n, n,       100000000000, n, n,
	//              10000000000, n, n, n,          1000000000, n, n,          100000000, n, n,
	//                 10000000, n, n, n,             1000000, n, n,             100000, n, n,
	//                    10000, n, n, n,                1000, n, n,                100, n, n,
	//                       10, n, n, n
	10000000000000000000, n, n, n, 1000000000000000000, n, n, 100000000000000000, n, n,
	10000000000000000, n, n, n, 1000000000000000, n, n, 100000000000000, n, n,
	10000000000000, n, n, n, 1000000000000, n, n, 100000000000, n, n,
	10000000000, n, n, n, 1000000000, n, n, 100000000, n, n,
	10000000, n, n, n, 1000000, n, n, 100000, n, n,
	10000, n, n, n, 1000, n, n, 100, n, n,
	10, n, n, n,
}

// FastUint64Log10 computes the integer base-10 logarithm of v, that is,
// the number of decimal digits in v minus one.
// The function is not well-defiend for v == 0.
func FastUint64Log10(v uint64) uint {
	lz := uint(bits.LeadingZeros64(v)) & 0x3f // &63 to eliminate bounds checking
	g := uint(0)
	if v >= lookup64[lz] {
		g = 1
	}
	return (63-lz)*3/10 + g
}

// Note: in the following table we use 64-bit values otherwise the condition
// v >= lookup[clz(v)] will be true for the very first entry when v == 2^32-1.
// Trying to force this to be 32-bit by adding an additional condition below
// makes the code overall slower.
var lookup32 = [32]uint64{
	// This initializer list is easier to read as follows:
	//                              n, n,          1000000000, n, n,          100000000, n, n,
	//                 10000000, n, n, n,             1000000, n, n,             100000, n, n,
	//                    10000, n, n, n,                1000, n, n,                100, n, n,
	//                       10, n, n, l
	n, n, 1000000000, n, n, 100000000, n, n,
	10000000, n, n, n, 1000000, n, n, 100000, n, n,
	10000, n, n, n, 1000, n, n, 100, n, n,
	10, n, n, n,
}

// FastUint32Log10 computes the integer base-10 logarithm of v, that is,
// the number of decimal digits in v minus one.
// The function is not well-defiend for v == 0.
func FastUint32Log10(v uint32) uint {
	lz := uint(bits.LeadingZeros32(v)) & 0x1f // &31 to eliminate bounds checking
	g := uint(0)
	if uint64(v) >= lookup32[lz] {
		g = 1
	}
	return (31-lz)*3/10 + g
}
