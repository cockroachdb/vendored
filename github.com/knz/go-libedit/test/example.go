package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	libedit "github.com/knz/go-libedit"
)

type example struct{}

func (_ example) GetLeftPrompt() string {
	return "hello> "
}

func (_ example) GetRightPrompt() string {
	return "(-:"
}

func (_ example) GetCompletions(word string) []string {
	if strings.HasPrefix(word, "he") {
		return []string{"hello!"}
	}
	return nil
}

func main() {
	el, err := libedit.Init("example")
	if err != nil {
		log.Fatal(err)
	}
	defer el.Close()

	// RebindControlKeys ensures that Ctrl+C, Ctrl+Z, Ctrl+R and Tab are
	// properly bound even if the user's .editrc has used bind -e or
	// bind -v to load a predefined keymap.
	el.RebindControlKeys()

	el.UseHistory(-1, true)
	el.SetCompleter(example{})
	el.SetLeftPrompt(example{})
	el.SetRightPrompt(example{})
	for {
		s, err := el.GetLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			if err == libedit.ErrInterrupted {
				fmt.Printf("interrupted! (%q)\n", s)
				continue
			}
			log.Fatal(err)
		}
		fmt.Printf("echo %q\n", s)
		if err := el.AddHistory(s); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("goodbye!")
}
