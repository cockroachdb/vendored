package returncheck_test

import (
	"fmt"

	"github.com/kkaneda/returncheck/internal"
)

func ExampleRun() {
	if err := returncheck.Run([]string{"testdata/test.go"}, "testdata/test.go", "Error"); err != nil {
		fmt.Printf("failed to run return check: %s\n", err)
		return
	}
	// Output:
	// testdata/test.go:22:3	g()
	// testdata/test.go:23:3	h()
	// testdata/test.go:24:11	x, _ := h()
	// testdata/test.go:29:6	go g()
	// testdata/test.go:30:9	defer h()
	// failed to run return check: found an unchecked return value
}
