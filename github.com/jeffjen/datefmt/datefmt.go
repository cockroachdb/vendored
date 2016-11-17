package datefmt

/*
#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

void ptime(const char* src, const char* fmt, struct tm* tm) {
	strptime(src, fmt, tm);
}

char* ftime(const char* fmt, const struct tm* tm) {
	const size_t size = 64;
	char *timestamp;
	if (timestamp = malloc(size * sizeof(char))) {
		strftime(timestamp, size, fmt, tm);
		return timestamp;
	} else {
		errno = EINVAL;
		return NULL;
	}
}
*/
import "C"

import (
	"fmt"
	"time"
	"unsafe"
)

const ()

func Strptime(format, timestamp string) (t time.Time, err error) {
	c_timestamp := C.CString(timestamp)
	c_format := C.CString(format)
	defer func() {
		C.free(unsafe.Pointer(c_timestamp))
		C.free(unsafe.Pointer(c_format))
	}()
	var c_time C.struct_tm
	if _, trr := C.ptime(c_timestamp, c_format, &c_time); trr != nil {
		err = fmt.Errorf("%s - %s:%s", trr, timestamp, format)
		return
	}
	t = time.Date(int(c_time.tm_year)+1900,
		time.Month(c_time.tm_mon+1),
		int(c_time.tm_mday),
		int(c_time.tm_hour),
		int(c_time.tm_min),
		int(c_time.tm_sec),
		0,
		time.FixedZone("", int(c_time.tm_gmtoff)),
	)
	return
}

func Strftime(format string, t time.Time) (timestamp string, err error) {
	c_format := C.CString(format)
	defer func() { C.free(unsafe.Pointer(c_format)) }()

	tz, offset := t.Zone()
	c_tz := C.CString(tz)
	defer func() { C.free(unsafe.Pointer(c_tz)) }()

	c_time := C.struct_tm{
		tm_year:   C.int(t.Year() - 1900),
		tm_mon:    C.int(t.Month() - 1),
		tm_mday:   C.int(t.Day()),
		tm_hour:   C.int(t.Hour()),
		tm_min:    C.int(t.Minute()),
		tm_sec:    C.int(t.Second()),
		tm_gmtoff: C.long(offset),
		tm_zone:   c_tz,
	}

	c_timestamp, trr := C.ftime(c_format, &c_time)
	defer func() { C.free(unsafe.Pointer(c_timestamp)) }()

	timestamp = C.GoString(c_timestamp)
	if trr == nil {
		timestamp = C.GoString(c_timestamp)
	} else {
		err = fmt.Errorf("%s - %s", trr, t)
	}
	return
}
