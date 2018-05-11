// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"strings"

	"github.com/golang-commonmark/markdown/byteutil"
)

var output = flag.String("o", "", "Output file")

type node struct {
	b byte    // byte
	s int     // state
	c []*node // children
	f bool    // final state?
}

func main() {
	log.SetFlags(0)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-o output.go] input.txt package.function\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}
	if !strings.Contains(args[1], ".") {
		flag.Usage()
		os.Exit(1)
	}
	split := strings.SplitN(args[1], ".", 2)
	pkg, fun := split[0], split[1]

	f, err := os.Open(args[0])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	root := &node{}
	var nodes []*node
	nodes = append(nodes, root)
	s := 0 // state
	scan := bufio.NewScanner(f)

	for scan.Scan() {
		str := scan.Text()

		cur := root
		for i := 0; i < len(str); i++ {
			b := str[i]
			var n *node
			found := false

			for _, n = range cur.c {
				if n.b == b {
					found = true
					break
				}
			}

			if !found {
				s++
				n = &node{b: b, s: s}
				cur.c = append(cur.c, n)
				nodes = append(nodes, n)
			}

			cur = n
			if i == len(str)-1 {
				cur.f = true
			}
		}
	}

	buf := bytes.NewBuffer(nil)

	fmt.Fprintf(buf, `package %s
			import "github.com/golang-commonmark/markdown/byteutil"
	`, pkg)
	fmt.Fprintf(buf, `
	func %s(s string) int {
	st := 0
	length := -1

	for i := 0; i < len(s); i++ {
		b := s[i]

		switch st {`, fun)

	for _, n := range nodes {
		if len(n.c) > 0 {
			fmt.Fprintf(buf, "case %d:\n", n.s)
			hasLetter := false
			for _, c := range n.c {
				if byteutil.IsLetter(c.b) {
					hasLetter = true
				}
			}
			if hasLetter {
				fmt.Fprintf(buf, "switch byteutil.ByteToLower(b) {\n")
			} else {
				fmt.Fprintf(buf, "switch b {\n")
			}

			for _, c := range n.c {
				fmt.Fprintf(buf, "case '%c':\n", c.b)
				if c.f {
					fmt.Fprintln(buf, "length = i+1")
				}
				fmt.Fprintf(buf, "st = %d\n", c.s)
			}

			fmt.Fprintf(buf, `default:
			return length
		}

			`)

		}
	}

	fmt.Fprintln(buf, `
		}
	}

	return length
}
`)
	source, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	if *output == "" || *output == "-" {
		os.Stdout.Write(source)
	} else {
		out, err := os.OpenFile(*output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		out.Write(source)
	}
}
