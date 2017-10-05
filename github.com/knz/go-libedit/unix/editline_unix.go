package libedit_unix

import (
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"
	"unsafe"

	common "github.com/knz/go-libedit/common"
)

// #cgo openbsd netbsd freebsd dragonfly darwin LDFLAGS: -ledit
// #cgo openbsd netbsd freebsd dragonfly darwin CPPFLAGS: -Ishim
// #cgo linux LDFLAGS: -ltinfo
// #cgo linux CFLAGS: -Wno-unused-result
// #cgo linux CPPFLAGS: -Isrc -Isrc/c-libedit -Isrc/c-libedit/editline -Isrc/c-libedit/linux-build
//
// #include <stdlib.h>
// #include <stdio.h>
// #include <unistd.h>
// #include <limits.h>
//
// #include <histedit.h>
// #include "c_editline.h"
import "C"

type EditLine int
type LeftPromptGenerator = common.LeftPromptGenerator
type RightPromptGenerator = common.RightPromptGenerator
type CompletionGenerator = common.CompletionGenerator

type state struct {
	el              *C.EditLine
	h               *C.History
	cIn, cOut, cErr *C.FILE
	inf, outf, errf *os.File
	promptGenLeft   LeftPromptGenerator
	promptGenRight  RightPromptGenerator
	cPromptLeft     *C.char
	cPromptRight    *C.char
	completer       CompletionGenerator
	histFile        *C.char
	autoSaveHistory bool
}

var editors []state

var errUnknown = errors.New("unknown error")

func Init(appName string) (EditLine, error) {
	return InitFiles(appName, os.Stdin, os.Stdout, os.Stderr)
}

func InitFiles(appName string, inf, outf, errf *os.File) (e EditLine, err error) {
	var inFile, outFile, errFile *C.FILE
	defer func() {
		if err == nil {
			return
		}
		if inFile != nil {
			C.fclose(inFile)
		}
		if outFile != nil {
			C.fclose(outFile)
		}
		if errFile != nil {
			C.fclose(errFile)
		}
	}()

	inFile, err = C.fdopen(C.dup(C.int(inf.Fd())), C.go_libedit_mode_read)
	if err != nil {
		return -1, fmt.Errorf("fdopen(inf): %v", err)
	}
	outFile, err = C.fdopen(C.dup(C.int(outf.Fd())), C.go_libedit_mode_write)
	if err != nil {
		return -1, fmt.Errorf("fdopen(outf): %v", err)
	}
	errFile, err = C.fdopen(C.dup(C.int(errf.Fd())), C.go_libedit_mode_write)
	if err != nil {
		return -1, fmt.Errorf("fdopen(errf): %v", err)
	}
	cAppName := C.CString(appName)
	defer C.free(unsafe.Pointer(cAppName))
	el, err := C.go_libedit_init(cAppName, inFile, outFile, errFile)
	// If the settings file did not exist, ignore the error.
	if err == syscall.ENOENT {
		err = nil
	}
	if el == nil || err != nil {
		if err == nil {
			err = errUnknown
		}
		return -1, fmt.Errorf("el_init: %v", err)
	}

	st := state{
		el:  el,
		inf: inf, outf: outf, errf: errf,
		cIn: inFile, cOut: outFile, cErr: errFile,
	}
	editors = append(editors, st)
	return EditLine(len(editors) - 1), nil
}

func (el EditLine) RebindControlKeys() {
	st := &editors[el]
	C.go_libedit_rebind_ctrls(st.el)
}

func (el EditLine) Close() {
	st := &editors[el]
	if st.el == nil {
		// Already closed.
		return
	}
	C.el_end(st.el)
	if st.h != nil {
		C.history_end(st.h)
	}
	if st.cPromptLeft != nil {
		C.free(unsafe.Pointer(st.cPromptLeft))
	}
	if st.cPromptRight != nil {
		C.free(unsafe.Pointer(st.cPromptRight))
	}
	if st.histFile != nil {
		C.free(unsafe.Pointer(st.histFile))
	}
	C.fclose(st.cIn)
	C.fclose(st.cOut)
	C.fclose(st.cErr)
	*st = state{}
}

var errNoHistory = errors.New("history not configured")
var errNoFileConfigured = errors.New("no savefile configured")

func (el EditLine) SaveHistory() error {
	st := &editors[el]
	if st.h == nil {
		return errNoHistory
	}
	if st.histFile == nil {
		return errNoFileConfigured
	}
	_, err := C.go_libedit_write_history(st.h, st.histFile)
	if err != nil {
		return fmt.Errorf("write_history: %v", err)
	}
	return nil
}

func (el EditLine) AddHistory(line string) error {
	st := &editors[el]
	if st.h == nil {
		return errNoHistory
	}

	cLine := C.CString(line)
	defer C.free(unsafe.Pointer(cLine))

	_, err := C.go_libedit_add_history(st.h, cLine)
	if err != nil {
		return fmt.Errorf("add_history: %v", err)
	}

	if st.autoSaveHistory && st.histFile != nil {
		_, err := C.go_libedit_write_history(st.h, st.histFile)
		if err != nil {
			return fmt.Errorf("write_history: %v", err)
		}
	}
	return nil
}

func (el EditLine) LoadHistory(file string) error {
	st := &editors[el]
	if st.h == nil {
		return errNoHistory
	}

	histFile := C.CString(file)
	defer C.free(unsafe.Pointer(histFile))
	_, err := C.go_libedit_read_history(st.h, histFile)
	if err != nil && err != syscall.ENOENT {
		return fmt.Errorf("read_history: %v", err)
	}
	return nil
}

func (el EditLine) SetAutoSaveHistory(file string, autoSave bool) {
	st := &editors[el]
	if st.h == nil {
		return
	}
	var newHistFile *C.char
	if file != "" {
		newHistFile = C.CString(file)
	}
	if st.histFile != nil {
		C.free(unsafe.Pointer(st.histFile))
		st.histFile = nil
	}
	st.histFile = newHistFile
	st.autoSaveHistory = autoSave
}

func (el EditLine) UseHistory(maxEntries int, dedup bool) error {
	st := &editors[el]

	cDedup := 0
	if dedup {
		cDedup = 1
	}

	cMaxEntries := C.int(maxEntries)
	if maxEntries < 0 {
		cMaxEntries = C.INT_MAX
	}

	h, err := C.go_libedit_setup_history(st.el, cMaxEntries, C.int(cDedup))
	if err != nil {
		return fmt.Errorf("init_history: %v", err)
	}
	if st.h != nil {
		C.history_end(st.h)
	}
	st.h = h
	return nil
}

func (el EditLine) GetLine() (string, error) {
	st := &editors[el]

	var count C.int
	var interrupted C.int
	s, err := C.go_libedit_gets(st.el, &count, &interrupted)
	if interrupted > 0 {
		// Reveal the partial line.
		line, _ := el.GetLineInfo()
		C.el_reset(st.el)
		return line, common.ErrInterrupted
	}
	if err != nil {
		return "", err
	}
	if s == nil {
		return "", io.EOF
	}
	return C.GoStringN(s, count), nil
}

func (el EditLine) Stdin() *os.File {
	return editors[el].inf
}

func (el EditLine) Stdout() *os.File {
	return editors[el].outf
}

func (el EditLine) Stderr() *os.File {
	return editors[el].errf
}

func (el EditLine) SetCompleter(gen CompletionGenerator) {
	editors[el].completer = gen
}

func (el EditLine) SetLeftPrompt(gen LeftPromptGenerator) {
	st := &editors[el]
	st.promptGenLeft = gen
	if gen != nil {
		C.go_libedit_set_prompt(st.el, C.EL_PROMPT, C.go_libedit_prompt_left_ptr)
	} else {
		C.go_libedit_set_prompt(st.el, C.EL_PROMPT, nil)
	}
}

func (el EditLine) SetRightPrompt(gen RightPromptGenerator) {
	st := &editors[el]
	st.promptGenRight = gen
	if gen != nil {
		C.go_libedit_set_prompt(st.el, C.EL_RPROMPT, C.go_libedit_prompt_right_ptr)
	} else {
		C.go_libedit_set_prompt(st.el, C.EL_RPROMPT, nil)
	}
}

func (el EditLine) GetLineInfo() (string, int) {
	st := &editors[el]
	li := C.el_line(st.el)
	return C.GoStringN(li.buffer,
		C.int(uintptr(unsafe.Pointer(li.lastchar))-uintptr(unsafe.Pointer(li.buffer))),
	), int(uintptr(unsafe.Pointer(li.cursor)) - uintptr(unsafe.Pointer(li.buffer)))
}

//export go_libedit_getcompletions
func go_libedit_getcompletions(cI C.int, cWord *C.char) **C.char {
	if int(cI) < 0 || int(cI) >= len(editors) {
		return nil
	}
	st := &editors[int(cI)]
	if st.completer == nil {
		return nil
	}

	word := C.GoString(cWord)
	matches := st.completer.GetCompletions(word)
	if len(matches) == 0 {
		return nil
	}

	array := (**C.char)(C.malloc(C.size_t(C.sizeof_pchar * (len(matches) + 1))))
	for i, m := range matches {
		C.go_libedit_set_string_array(array, C.int(i), C.CString(m))
	}
	C.go_libedit_set_string_array(array, C.int(len(matches)), nil)
	return array
}

//export go_libedit_prompt_left
func go_libedit_prompt_left(el *C.EditLine) *C.char {
	i := C.go_libedit_get_clientdata(el)
	e := int(i)
	if e < 0 || e >= len(editors) {
		return C.go_libedit_emptycstring
	}
	st := &editors[e]
	if st.cPromptLeft != nil {
		C.free(unsafe.Pointer(st.cPromptLeft))
		st.cPromptLeft = nil
	}
	if st.promptGenLeft == nil {
		return C.go_libedit_emptycstring
	}
	st.cPromptLeft = C.CString(st.promptGenLeft.GetLeftPrompt())
	return st.cPromptLeft
}

//export go_libedit_prompt_right
func go_libedit_prompt_right(el *C.EditLine) *C.char {
	i := C.go_libedit_get_clientdata(el)
	e := int(i)
	if e < 0 || e >= len(editors) {
		return C.go_libedit_emptycstring
	}
	st := &editors[e]
	if st.cPromptRight != nil {
		C.free(unsafe.Pointer(st.cPromptRight))
		st.cPromptRight = nil
	}
	if st.promptGenRight == nil {
		return C.go_libedit_emptycstring
	}
	st.cPromptRight = C.CString(st.promptGenRight.GetRightPrompt())
	return st.cPromptRight
}
