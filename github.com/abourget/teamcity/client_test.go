package teamcity

import "testing"

func TestClientGetBuildProperties(t *testing.T) {
	client := NewTestClient(newResponse(`{"property":[{"name": "build.counter", "value": "12"}], "count": 1}`), nil)

	props, err := client.GetBuildProperties("999999")

	if len(props) != 1 {
		t.Fatal("Expected to have 1 property, found", len(props))
	}

	if err != nil {
		t.Fatal("Expected no error, got", err)
	}
}
