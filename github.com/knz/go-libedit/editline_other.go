// +build !darwin,!freebsd,!linux,!openbsd,!netbsd,!dragonfly

package libedit

import (
	"os"

	edit "github.com/knz/go-libedit/other"
)

type EditLine = edit.EditLine

func Init(x string) (EditLine, error) { return edit.Init(x) }
func InitFiles(a string, stdin, stdout, stderr *os.File) (EditLine, error) {
	return edit.InitFiles(a, stdin, stdout, stderr)
}
