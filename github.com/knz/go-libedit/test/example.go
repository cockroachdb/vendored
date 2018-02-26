package main

import (
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	libedit "github.com/knz/go-libedit"
)

type example struct{}

func (_ example) GetCompletions(word string) []string {
	if strings.HasPrefix(word, "he") {
		return []string{"hello!"}
	}
	return nil
}

func main() {
	// Open and immediately close a libedit instance to test that nonzero editor
	// IDs are tracked correctly.
	el, err := libedit.Init("example", true)
	if err != nil {
		log.Fatal(err)
	}
	el.Close()

	el, err = libedit.Init("example", true)
	if err != nil {
		log.Fatal(err)
	}
	defer el.Close()

	// RebindControlKeys ensures that Ctrl+C, Ctrl+Z, Ctrl+R and Tab are
	// properly bound even if the user's .editrc has used bind -e or
	// bind -v to load a predefined keymap.
	el.RebindControlKeys()

	el.UseHistory(-1, true)
	el.LoadHistory("hist")
	el.SetAutoSaveHistory("hist", true)
	el.SetCompleter(example{})
	el.SetLeftPrompt("hello> ")
	el.SetRightPrompt("(-:")
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
		time.Sleep(2 * time.Second)
		fmt.Println("ok")
		if err := el.AddHistory(s); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("goodbye!")
}
