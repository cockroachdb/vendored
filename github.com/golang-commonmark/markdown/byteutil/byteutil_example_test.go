// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package byteutil

import (
	"fmt"
)

func ExampleIsDigit() {
	fmt.Println(IsDigit('0'))
	fmt.Println(IsDigit('9'))
	fmt.Println(IsDigit('a'))

	// Output:
	// true
	// true
	// false
}

func ExampleIsHexDigit() {
	fmt.Println(IsHexDigit('0'))
	fmt.Println(IsHexDigit('9'))
	fmt.Println(IsHexDigit('a'))
	fmt.Println(IsHexDigit('F'))
	fmt.Println(IsHexDigit('G'))

	// Output:
	// true
	// true
	// true
	// true
	// false
}

func ExampleIsLetter() {
	fmt.Println(IsLetter('a'))
	fmt.Println(IsLetter('Z'))
	fmt.Println(IsLetter('0'))

	// Output:
	// true
	// true
	// false
}

func ExampleIsLowercaseLetter() {
	fmt.Println(IsLowercaseLetter('a'))
	fmt.Println(IsLowercaseLetter('z'))
	fmt.Println(IsLowercaseLetter('A'))

	// Output:
	// true
	// true
	// false
}

func ExampleIsUppercaseLetter() {
	fmt.Println(IsUppercaseLetter('A'))
	fmt.Println(IsUppercaseLetter('Z'))
	fmt.Println(IsUppercaseLetter('a'))

	// Output:
	// true
	// true
	// false
}

func ExampleIsAlphaNum() {
	fmt.Println(IsAlphaNum('0'))
	fmt.Println(IsAlphaNum('9'))
	fmt.Println(IsAlphaNum('a'))
	fmt.Println(IsAlphaNum('Z'))
	fmt.Println(IsAlphaNum('.'))

	// Output:
	// true
	// true
	// true
	// true
	// false
}

func ExampleToLower() {
	fmt.Println(ToLower(".012345abcdefABCDEFабвгд"))

	// Output:
	// .012345abcdefabcdefабвгд
}

func ExampleToUpper() {
	fmt.Println(ToUpper(".012345abcdefABCDEFабвгд"))

	// Output:
	// .012345ABCDEFABCDEFабвгд
}

func ExampleByteToLower() {
	fmt.Printf("%c\n", ByteToLower('a'))
	fmt.Printf("%c\n", ByteToLower('A'))
	fmt.Printf("%c\n", ByteToLower('0'))
	fmt.Printf("%c\n", ByteToLower('.'))

	// Output:
	// a
	// a
	// 0
	// .
}

func ExampleByteToUpper() {
	fmt.Printf("%c\n", ByteToUpper('a'))
	fmt.Printf("%c\n", ByteToUpper('A'))
	fmt.Printf("%c\n", ByteToUpper('0'))
	fmt.Printf("%c\n", ByteToUpper('.'))

	// Output:
	// A
	// A
	// 0
	// .
}

func ExampleIndexAny() {
	fmt.Println(IndexAny("abcdefghijklmnopqrstuvwxyz", "zyx"))

	// Output:
	// 23
}

func ExampleIndexAnyTable() {
	var t [256]bool
	t['z'], t['y'], t['x'] = true, true, true
	fmt.Println(IndexAnyTable("abcdefghijklmnopqrstuvwxyz", &t))

	// Output:
	// 23
}

func ExampleUnhex() {
	fmt.Println(Unhex('0'))
	fmt.Println(Unhex('9'))
	fmt.Println(Unhex('a'))
	fmt.Println(Unhex('A'))
	fmt.Println(Unhex('f'))
	fmt.Println(Unhex('F'))

	// Output:
	// 0
	// 9
	// 10
	// 10
	// 15
	// 15
}
