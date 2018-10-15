// +build darwin
// +build cgo

package cpu

import (
	"errors"
	"fmt"
	"runtime/debug"
	"unsafe"
)

// #include <stdlib.h>
// #include "cpu_darwin.h"
import "C"

var timesStatSize = unsafe.Sizeof(C.times_stat{})

func perCPUTimes() ([]TimesStat, error) {
	var n C.int
	cTimesStats := C.per_cpu_times(&n)
	if cTimesStats == nil {
		return nil, errors.New("unable to collect per CPU times")
	}
	defer C.free(unsafe.Pointer(cTimesStats))

	out := make([]TimesStat, n)
	for i := range out {
		// We can't index into cTimesStats directly because it is a pointer, not a
		// slice.
		cTimesStat := (*C.times_stat)(unsafe.Pointer(uintptr(unsafe.Pointer(cTimesStats)) + uintptr(i)*timesStatSize))
		out[i] = TimesStat{
			CPU:    fmt.Sprintf("cpu%d", cTimesStat.cpu),
			User:   float64(cTimesStat.user),
			System: float64(cTimesStat.system),
			Idle:   float64(cTimesStat.idle),
			Nice:   float64(cTimesStat.nice),
		}
		fmt.Printf("cpu load per %#v\n", out[i])
	}
	debug.PrintStack()
	return out, nil
}

func allCPUTimes() (out []TimesStat, err error) {
	defer func() { fmt.Printf("cpu load total %#v\n", out) }()
	cTimesStat := C.all_cpu_times()
	return []TimesStat{{
		CPU:    fmt.Sprintf("cpu%d", cTimesStat.cpu),
		User:   float64(cTimesStat.user),
		System: float64(cTimesStat.system),
		Idle:   float64(cTimesStat.idle),
		Nice:   float64(cTimesStat.nice),
	}}, nil
}
