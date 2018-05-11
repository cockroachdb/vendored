package markdown

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/opennota/wd"
)

type example struct {
	Filename string
	Num      int `json:"example"`
	Markdown string
	HTML     string
	Section  string
}

func loadExamplesFromJSON(fn string) []example {
	f, err := os.Open(fn)
	if err != nil {
		panic(err)
	}

	var examples []example
	err = json.NewDecoder(f).Decode(&examples)
	if err != nil {
		panic(err)
	}

	return examples
}

func render(src string, options ...option) (_ string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	md := New(options...)
	return md.RenderToString([]byte(src)), nil
}

func TestCommonMark(t *testing.T) {
	examples := loadExamplesFromJSON("spec/commonmark-0.21.json")
	for _, ex := range examples {
		result, err := render(ex.Markdown, HTML(true), XHTMLOutput(true), Linkify(false), Typographer(false), LangPrefix("language-"))
		if err != nil {
			t.Errorf("#%d (%s): PANIC (%v)", ex.Num, ex.Section, err)
		} else if result != ex.HTML {
			d := wd.ColouredDiff(ex.HTML, result, false)
			d = wd.NumberLines(d)
			t.Errorf("#%d (%s):\n%s", ex.Num, ex.Section, d)
		}
	}
}

func TestRenderSpec(t *testing.T) {
	data, err := ioutil.ReadFile("spec/spec-0.21.txt")
	if err != nil {
		t.Fatal(err)
	}

	md := New(HTML(true), XHTMLOutput(true))
	md.RenderToString(data)
}
