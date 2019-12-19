package ilog10

// NumInt32DecimalDigits returns the number of decimal digits in n
// (excluding the negative sign, if any). 0 is considered to have one
// digit.
func NumInt32DecimalDigits(n int32) uint {
	if n == 0 {
		return 1
	}
	u := uint32(n)
	if n < 0 {
		u = -u
	}
	return 1 + FastUint32Log10(u)
}

// NumUint32DecimalDigits returns the number of decimal digits in n
// (excluding the negative sign, if any). 0 is considered to have one
// digit.
func NumUint32DecimalDigits(n uint32) uint {
	if n == 0 {
		return 1
	}
	return 1 + FastUint32Log10(n)
}

// NumInt64DecimalDigits returns the number of decimal digits in n
// (excluding the negative sign, if any). 0 is considered to have one
// digit.
func NumInt64DecimalDigits(n int64) uint {
	if n == 0 {
		return 1
	}
	u := uint64(n)
	if n < 0 {
		u = -u
	}
	return 1 + FastUint64Log10(u)
}

// NumUint64DecimalDigits returns the number of decimal digits in n
// (excluding the negative sign, if any). 0 is considered to have one
// digit.
func NumUint64DecimalDigits(n uint64) uint {
	if n == 0 {
		return 1
	}
	return 1 + FastUint64Log10(n)
}
