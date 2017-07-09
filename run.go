package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

// RunConfig holds settings for running the simulation.
type RunConfig struct {
	Prompt   bool
	MaxTurns int
	Delay    time.Duration
	Outfile  io.Writer
}

// DefaultRunConfig defines the default run settings.
var DefaultRunConfig = RunConfig{
	Prompt:   false,
	MaxTurns: 0,
	Delay:    time.Millisecond * 500,
	Outfile:  os.Stdout,
}

// Run runs a Game of Life simulation.
func Run(grid Grid, opts RunConfig) {
	if opts.MaxTurns > 0 {
		for i := 0; i < opts.MaxTurns; i++ {
			grid = NextTurn(grid, opts)
		}
	} else {
		for {
			grid = NextTurn(grid, opts)
		}
	}
}

func NextTurn(grid Grid, opts RunConfig) Grid {
	fmt.Fprintln(opts.Outfile, grid.Show())
	time.Sleep(opts.Delay)
	return grid.Next()
}
