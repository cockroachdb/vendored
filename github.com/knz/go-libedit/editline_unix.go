// +build darwin freebsd linux openbsd netbsd dragonfly

package libedit

import (
	"os"

	edit "github.com/knz/go-libedit/unix"
)

type EditLine = edit.EditLine

func Init(x string, w bool) (EditLine, error) { return edit.Init(x, w) }
func InitFiles(a string, w bool, stdin, stdout, stderr *os.File) (EditLine, error) {
	return edit.InitFiles(a, w, stdin, stdout, stderr)
}
