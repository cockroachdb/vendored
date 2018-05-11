// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package byteutil

import (
	"strings"
	"testing"
)

const printable = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"

func IsDigitUsingCompare(b byte) bool {
	return '0' <= b && b <= '9'
}

func IsDigitUsingSubtraction(b byte) bool {
	return b-'0' <= 9
}

func IsHexDigitUsingCompare(b byte) bool {
	return '0' <= b && b <= '9' || 'a' <= b && b <= 'f' || 'A' <= b && b <= 'F'
}

func IsLetterUsingCompare(b byte) bool {
	b |= 0x20 // to lower case
	return 'a' <= b && b <= 'z'
}

func IsLowercaseLetterUsingCompare(b byte) bool {
	return 'a' <= b && b <= 'z'
}

func IsUppercaseLetterUsingCompare(b byte) bool {
	return 'A' <= b && b <= 'Z'
}

func IsAlphaNumUsingCompare(b byte) bool {
	return '0' <= b && b <= '9' || 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z'
}

func BenchmarkIsDigitTable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 256; j++ {
			IsDigit(byte(j))
		}
	}
}

func BenchmarkIsDigitCompare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 256; j++ {
			IsDigitUsingCompare(byte(j))
		}
	}
}

func BenchmarkIsDigitSubtraction(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 256; j++ {
			IsDigitUsingSubtraction(byte(j))
		}
	}
}

func BenchmarkIsHexDigitTable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 256; j++ {
			IsHexDigit(byte(j))
		}
	}
}

func BenchmarkIsHexDigitCompare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 256; j++ {
			IsHexDigitUsingCompare(byte(j))
		}
	}
}

func BenchmarkIsLetterTable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 256; j++ {
			IsLetter(byte(j))
		}
	}
}

func BenchmarkIsLetterCompare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 256; j++ {
			IsLetterUsingCompare(byte(j))
		}
	}
}

func BenchmarkIsLowercaseLetterTable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 256; j++ {
			IsLowercaseLetter(byte(j))
		}
	}
}

func BenchmarkIsLowercaseLetterCompare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 256; j++ {
			IsLowercaseLetterUsingCompare(byte(j))
		}
	}
}

func BenchmarkIsUppercaseLetterTable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 256; j++ {
			IsUppercaseLetter(byte(j))
		}
	}
}

func BenchmarkIsUppercaseLetterCompare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 256; j++ {
			IsUppercaseLetterUsingCompare(byte(j))
		}
	}
}

func BenchmarkIsAlphaNumTable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 256; j++ {
			IsAlphaNum(byte(j))
		}
	}
}

func BenchmarkIsAlphaNumCompare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 256; j++ {
			IsAlphaNumUsingCompare(byte(j))
		}
	}
}

func BenchmarkToLowerTable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ToLower(printable)
	}
}

func BenchmarkToLowerStringsToLower(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strings.ToLower(printable)
	}
}

func BenchmarkToUpperTable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ToUpper(printable)
	}
}

func BenchmarkToUpperStringsToUpper(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strings.ToUpper(printable)
	}
}

func BenchmarkIndexAny(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IndexAny(printable, "abcdefghijklmnopqrstuvwxyz")
		IndexAny(printable, " abcdefghijklmnopqrstuvwxyz")
		IndexAny(printable, " ")
	}
}

func BenchmarkIndexAnyTable(b *testing.B) {
	var t1, t2, t3 [256]bool
	for _, b := range "abcdefghijklmnopqrstuvwxyz" {
		t1[b] = true
	}
	for _, b := range " abcdefghijklmnopqrstuvwxyz" {
		t2[b] = true
	}
	for _, b := range " " {
		t3[b] = true
	}
	for i := 0; i < b.N; i++ {
		IndexAnyTable(printable, &t1)
		IndexAnyTable(printable, &t2)
		IndexAnyTable(printable, &t3)
	}
}

func BenchmarkStringsIndexAny(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strings.IndexAny(printable, "abcdefghijklmnopqrstuvwxyz")
		strings.IndexAny(printable, " abcdefghijklmnopqrstuvwxyz")
		strings.IndexAny(printable, " ")
	}
}
