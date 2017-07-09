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
		Guard(err, "%s", err)
	}

	if nargs := flag.NArg(); nargs != 1 {
		if nargs > 1 {
			Fail("too many arguments")
		} else {
			Fail("too few arguments")
		}
	}

	gridFile, err := os.Open(flag.Arg(0))
	Guard(err, err)

	gridBytes, err := ioutil.ReadAll(gridFile)
	Guard(err, err)

	grid := conway.FromString(string(gridBytes))

	return grid, config
}

func Guard(err error, msg interface{}, a ...interface{}) {
	if err != nil {
		Fail(msg, a...)
	}
}

func Fail(msg interface{}, a ...interface{}) {
	fullMsg := fmt.Sprintf("conway: error: %s\n", msg)
	fmt.Printf(fullMsg, a...)
	os.Exit(1)
}
