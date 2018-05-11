package markdown

import "testing"

func TestNextQuoteIndex(t *testing.T) {
	type testCase struct {
		in   string
		from int
		want int
	}
	testCases := []testCase{
		{"", 0, -1},
		{`"xxx"`, 0, 0},
		{`"xxx"`, 1, 4},
		{"'xxx'", 1, 4},
	}
	for _, tc := range testCases {
		got := nextQuoteIndex([]rune(tc.in), tc.from)
		if got != tc.want {
			t.Errorf("nextQuoteIndex(%q, %d) = %d, want %d", tc.in, tc.from, got, tc.want)
		}
	}
}
