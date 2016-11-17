package datefmt

import (
	"fmt"
	"testing"
	"time"
)

func TestStrftime(t *testing.T) {
	timestamp, err := Strftime("%T", time.Now())
	if err != nil {
		t.Errorf("failed to format current time: %s", err)
	} else {
		fmt.Println(timestamp)
	}
}

func TestStrptime(t *testing.T) {
	d, err := Strptime("%F-%T", "2015-10-5-16:34:55")
	if err != nil {
		t.Errorf("failed to parse time: %s", err)
	} else {
		fmt.Println(d)
	}
}
