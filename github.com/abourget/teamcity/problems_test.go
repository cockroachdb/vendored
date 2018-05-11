package teamcity

import (
	"encoding/json"
	"fmt"
	"testing"
)

var JSONresponse = []byte(`{
  "count": 1,
  "href": "https://teamcity-or.intel.com/app/rest/problemOccurrences?locator=build:(id:431935)",
  "problemOccurrence": [
    {
      "id": "problem:(id:135),build:(id:431935)",
      "type": "TC_FAILED_TESTS",
      "identity": "TC_FAILED_TESTS",
      "href": "/httpAuth/app/rest/problemOccurrences/problem:(id:135),build:(id:431935)"
    }
  ],
  "default": false
}`)

func TestResponseParsing(t *testing.T) {
	var v struct {
		Count             int64
		Default           bool
		HREF              string
		ProblemOccurrence ProblemOccurrence
	}

	err := json.Unmarshal(JSONresponse, &v)
	if err != nil {
		t.Errorf("json unmarshal: %s", err)
	}

	fmt.Printf("%#v", v)

}
