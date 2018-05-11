package markdown

import "testing"

func TestMatchAutolink(t *testing.T) {
	type testCase struct {
		in   string
		want string
	}
	testCases := []testCase{
		{"", ""},
		{"%#!", ""},
		{"<google.com>", ""},
		{"<http:>", ""},
		{"<http://google.com", ""},
		{"<http://google.com>", "http://google.com"},
		{"<http://url with spaces>", ""},
		{"<http://\x00>", ""},
		{"<http:\x00>", ""},
		{"<%#!://url>", ""},
		{"<waytoolongstringforschema://url>", ""},
		{"<ws:x>", "ws:x"},
		{"<xxx://url>", ""},
	}
	for _, tc := range testCases {
		got := matchAutolink(tc.in)
		if got != tc.want {
			t.Errorf("matchAutolink(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestMatchEmail(t *testing.T) {
	type testCase struct {
		in   string
		want string
	}
	testCases := []testCase{
		{"", ""},
		{"a@a.ru", ""},
		{"<a@a.ru>", "a@a.ru"},
		{"<@a.ru>", ""},
		{"<a@xxx>", ""},
		{"<bradfitz@golang.org>", "bradfitz@golang.org"},
		{"<r@golang.org", ""},
		{"<r@golang\x00org>", ""},
		{"<r@\x00golang.org>", ""},
		{"<r\x00@golang.org>", ""},
		{"<\x00r@golang.org>", ""},
	}
	for _, tc := range testCases {
		got := matchEmail(tc.in)
		if got != tc.want {
			t.Errorf("matchEmail(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
