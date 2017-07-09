package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dustinrohde/go-conway"
)

func main() {
	grid, opts := ParseArgs(os.Args[1:])
	conway.Run(grid, opts)
}

func ParseArgs(args []string) (conway.Grid, conway.RunConfig) {
	config := conway.DefaultRunConfig()

	flag.BoolVar(&config.Prompt, "prompt", config.Prompt,
		"Wait for input between each turn.")
	flag.IntVar(&config.MaxTurns, "turns", config.MaxTurns,
		"Number of turns to run. If < 0, run forever.")
	flag.DurationVar(&config.Delay, "delay", config.Delay,
		"Delay between turns, e.g. 500ms, 3s, 1m")

	outfilePtr := flag.String("out", "", "File to write results to. Defaults to stdout.")

	flag.CommandLine.Parse(args)

	if *outfilePtr != "" {
		var err error
		config.Outfile, err = os.OpenFile(*outfilePtr, os.O_APPEND, os.ModeAppend)
		Guard(err)
	}

	if nargs := flag.NArg(); nargs != 1 {
		if nargs > 1 {
			Fail("too many arguments")
		} else {
			Fail("too few arguments")
		}
	}

	gridFile, err := os.Open(flag.Arg(0))
	Guard(err)

	gridBytes, err := ioutil.ReadAll(gridFile)
	Guard(err)

	grid := conway.FromString(string(gridBytes))

	return grid, config
}

func Guard(err error, fmtArgs ...interface{}) {
	if err != nil {
		if len(fmtArgs) > 0 {
			Fail(fmtArgs[0], fmtArgs[1:]...)
		} else {
			Fail(err)
		}
	}
}

func Fail(msg interface{}, fmtArgs ...interface{}) {
	fullMsg := fmt.Sprintf("conway: error: %s\n", msg)
	fmt.Printf(fullMsg, fmtArgs...)
	os.Exit(1)
}
