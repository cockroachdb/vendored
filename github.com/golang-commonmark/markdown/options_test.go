package markdown

import (
	"testing"
)

func TestQuote(t *testing.T) {
	quotes := "‘’“”"
	md := New(Quotes(quotes))
	if string(md.Quotes[:]) != quotes {
		t.Errorf("expected %q, got %q", quotes, md.Quotes)
	}
}
