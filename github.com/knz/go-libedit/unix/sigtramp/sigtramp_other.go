// +build !darwin !amd64

package sigtramp

import "unsafe"

func Get() unsafe.Pointer {
	return unsafe.Pointer(uintptr(0))
}
