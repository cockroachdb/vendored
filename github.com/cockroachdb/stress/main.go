// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The stress utility is intended for catching of episodic failures.
// It runs a given process in parallel in a loop and collects any failures.
// Usage:
// 	$ stress ./fmt.test -test.run=TestSometing -test.cpu=10
// You can also specify a number of parallel processes with -p flag;
// instruct the utility to not kill hanged processes for gdb attach;
// or specify the failure output you are looking for (if you want to
// ignore some other episodic failures).
package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"runtime"
	"sync"
	"syscall"
	"time"
)

var (
	flags        = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flagP        = flags.Int("p", runtime.NumCPU(), "run `N` processes in parallel")
	flagTimeout  = flags.Duration("timeout", 0, "timeout each process after `duration`")
	flagKill     = flags.Bool("kill", true, "kill timed out processes if true, otherwise just print pid (to attach with gdb)")
	flagFailure  = flags.String("failure", "", "fail only if output matches `regexp`")
	flagIgnore   = flags.String("ignore", "", "ignore failure if output matches `regexp`")
	flagMaxTime  = flags.Duration("maxtime", 0, "maximum time to run")
	flagMaxRuns  = flags.Int("maxruns", 0, "maximum number of runs")
	flagMaxFails = flags.Int("maxfails", 1, "maximum number of failures")
	flagStdErr   = flags.Bool("stderr", true, "output failures to STDERR instead of to a temp file")
)

func roundToSeconds(d time.Duration) time.Duration {
	return time.Duration(d.Seconds()+0.5) * time.Second
}

func run() error {
	if err := flags.Parse(os.Args[1:]); err != nil {
		return err
	}
	if *flagP <= 0 || *flagTimeout < 0 || len(flags.Args()) == 0 {
		var b bytes.Buffer
		flags.SetOutput(&b)
		flags.Usage()
		return errors.New(b.String())
	}
	var failureRe, ignoreRe *regexp.Regexp
	if *flagFailure != "" {
		var err error
		if failureRe, err = regexp.Compile(*flagFailure); err != nil {
			return fmt.Errorf("bad failure regexp: %s", err)
		}
	}
	if *flagIgnore != "" {
		var err error
		if ignoreRe, err = regexp.Compile(*flagIgnore); err != nil {
			return fmt.Errorf("bad ignore regexp: %s", err)
		}
	}

	c := make(chan os.Signal)
	defer close(c)
	signal.Notify(c, os.Interrupt)
	// TODO(tamird): put this behind a !windows build tag.
	signal.Notify(c, syscall.SIGHUP, syscall.SIGTERM)
	defer signal.Stop(c)
	var wg sync.WaitGroup
	defer wg.Wait()
	ctx, cancel := func(ctx context.Context) (context.Context, context.CancelFunc) {
		if *flagMaxTime > 0 {
			return context.WithTimeout(ctx, *flagMaxTime)
		}
		return context.WithCancel(ctx)
	}(context.Background())
	defer cancel()
	go func() {
		for range c {
			cancel()
		}
	}()

	startTime := time.Now()

	res := make(chan []byte)
	wg.Add(*flagP)
	for i := 0; i < *flagP; i++ {
		go func(ctx context.Context) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case res <- func(ctx context.Context) []byte {
					var cmd *exec.Cmd
					if *flagTimeout > 0 {
						if *flagKill {
							var cancel context.CancelFunc
							ctx, cancel = context.WithTimeout(ctx, *flagTimeout)
							defer cancel()
						} else {
							defer time.AfterFunc(*flagTimeout, func() {
								fmt.Printf("process %v timed out\n", cmd.Process.Pid)
							}).Stop()
						}
					}
					cmd = exec.CommandContext(ctx, flags.Args()[0], flags.Args()[1:]...)
					out, err := cmd.CombinedOutput()
					if err != nil && (failureRe == nil || failureRe.Match(out)) && (ignoreRe == nil || !ignoreRe.Match(out)) {
						out = append(out, fmt.Sprintf("\n\nERROR: %v\n", err)...)
					} else {
						out = []byte{}
					}
					return out
				}(ctx):
				}
			}
		}(ctx)
	}
	runs, fails := 0, 0
	ticker := time.NewTicker(5 * time.Second).C
	for {
		select {
		case out := <-res:
			runs++
			if *flagMaxRuns > 0 && runs >= *flagMaxRuns {
				cancel()
			}
			if len(out) == 0 {
				continue
			}
			fails++
			if *flagMaxFails > 0 && fails >= *flagMaxFails {
				cancel()
			}
			if *flagStdErr {
				fmt.Fprintf(os.Stderr, "\n%s\n", out)
			} else {
				f, err := ioutil.TempFile("", "go-stress")
				if err != nil {
					return fmt.Errorf("failed to create temp file: %v", err)
				}
				if _, err := f.Write(out); err != nil {
					return fmt.Errorf("failed to write temp file: %v", err)
				}
				if err := f.Close(); err != nil {
					return fmt.Errorf("failed to close temp file: %v", err)
				}
				if len(out) > 2<<10 {
					out = out[:2<<10]
				}
				fmt.Printf("\n%s\n%s\n", f.Name(), out)
			}
		case <-ticker:
			fmt.Printf("%v runs so far, %v failures, over %s\n",
				runs, fails, roundToSeconds(time.Since(startTime)))
		case <-ctx.Done():

			fmt.Printf("%v runs completed, %v failures, over %s\n",
				runs, fails, roundToSeconds(time.Since(startTime)))

			switch err := ctx.Err(); err {
			// A context timeout in this case is indicative of no failures
			// being detected in the allotted duration.
			case context.DeadlineExceeded:
				return nil
			case context.Canceled:
				if *flagMaxRuns > 0 && runs >= *flagMaxRuns {
					return nil
				}
				return err
			default:
				return fmt.Errorf("unexpected context error: %v", err)
			}
		}
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Println("FAIL")
		os.Exit(1)
	} else {
		fmt.Println("SUCCESS")
	}
}
