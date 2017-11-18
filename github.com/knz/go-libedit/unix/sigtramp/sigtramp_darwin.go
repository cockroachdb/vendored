// +build darwin

package sigtramp

import "unsafe"

func getptr() unsafe.Pointer

func Get() unsafe.Pointer {
	return getptr()
}
