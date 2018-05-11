package markdown

import (
	"reflect"
	"testing"
)

func TestNormalizeAndIndex(t *testing.T) {
	type testCase struct {
		in      string
		out     string
		b, e, s []int
	}
	testCases := []testCase{
		{"abc", "abc", []int{0}, []int{3}, []int{0}},
		{" abc", " abc", []int{0}, []int{4}, []int{1}},
		{"    abc", "    abc", []int{0}, []int{7}, []int{4}},
		{"abc\n", "abc\n", []int{0}, []int{3}, []int{0}},
		{"abc\r", "abc\n", []int{0}, []int{3}, []int{0}},
		{"abc\r\n", "abc\n", []int{0}, []int{3}, []int{0}},
		{"abc\td", "abc d", []int{0}, []int{5}, []int{0}},
		{"abc\x00def", "abc\ufffddef", []int{0}, []int{9}, []int{0}},
		{"ab\tc", "ab  c", []int{0}, []int{5}, []int{0}},
		{"a\tabc", "a   abc", []int{0}, []int{7}, []int{0}},
		{"", "", nil, nil, nil},
		{"\tabc", "    abc", []int{0}, []int{7}, []int{4}},
		{"abc\n def\r\n\tghi\rj\tkl\x00\n", "abc\n def\n    ghi\nj   kl\ufffd\n", []int{0, 4, 9, 17}, []int{3, 8, 16, 26}, []int{0, 1, 4, 0}},
		{"абв\n где\r\n\tёжз\rи\tйк\x00\n", "абв\n где\n    ёжз\nи   йк\ufffd\n", []int{0, 7, 15, 26}, []int{6, 14, 25, 38}, []int{0, 1, 4, 0}},
	}
	for _, tc := range testCases {
		out, b, e, s := normalizeAndIndex([]byte(tc.in))
		if out != tc.out {
			t.Errorf("normalize(%q):\nstring = %q\n    want %q", tc.in, out, tc.out)
		}
		if !reflect.DeepEqual(b, tc.b) {
			t.Errorf("normalize(%q):\nbMarks = %#v\n    want %#v", tc.in, b, tc.b)
		}
		if !reflect.DeepEqual(e, tc.e) {
			t.Errorf("normalize(%q):\neMarks = %#v\n    want %#v", tc.in, e, tc.e)
		}
		if !reflect.DeepEqual(s, tc.s) {
			t.Errorf("normalize(%q):\ntShift = %#v\n    want %#v", tc.in, s, tc.s)
		}
	}
}
