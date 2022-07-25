package snowball

import (
	"fmt"

	"github.com/blevesearch/snowball/english"
	"github.com/blevesearch/snowball/french"
	"github.com/blevesearch/snowball/norwegian"
	"github.com/blevesearch/snowball/russian"
	"github.com/blevesearch/snowball/spanish"
	"github.com/blevesearch/snowball/swedish"
)

const (
	VERSION string = "v0.7.0"
)

// Stem a word in the specified language.
//
func Stem(word, language string, stemStopWords bool) (stemmed string, err error) {

	var f func(string, bool) string
	switch language {
	case "english":
		f = english.Stem
	case "spanish":
		f = spanish.Stem
	case "french":
		f = french.Stem
	case "russian":
		f = russian.Stem
	case "swedish":
		f = swedish.Stem
	case "norwegian":
		f = norwegian.Stem
	default:
		err = fmt.Errorf("Unknown language: %s", language)
		return
	}
	stemmed = f(word, stemStopWords)
	return

}
