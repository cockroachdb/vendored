package libedit_common

import "errors"

type LeftPromptGenerator interface {
	GetLeftPrompt() string
}

type RightPromptGenerator interface {
	GetRightPrompt() string
}

type CompletionGenerator interface {
	GetCompletions(word string) []string
}

var ErrInterrupted = errors.New("interrupted")
var ErrWidecharNotSupported = errors.New("cannot enable wide character support")
