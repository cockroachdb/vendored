package libedit_common

import "errors"

type CompletionGenerator interface {
	GetCompletions(word string) []string
}

var ErrInterrupted = errors.New("interrupted")
var ErrWidecharNotSupported = errors.New("cannot enable wide character support")
