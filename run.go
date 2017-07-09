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
const DefaultRunConfig = RunConfig{
	Prompt:   false,
	MaxTurns: 0,
	Delay:    time.Millisecond * 500,
	Outfile:  os.Stdout,
}

// Run runs a Game of Life simulation.
func Run(grid Grid, opts RunConfig) {
	if opts.MaxTurns > 0 {
		for i := 0; i < opts.MaxTurns; i++ {
			Turn(grid, opts)
		}
	} else {
		for {
			Turn(grid, opts)
		}
	}
}

func Turn(grid Grid, opts RunConfig) {
	fmt.Fprintln(opts.Outfile, grid.Show())
	time.Sleep(delay)
	grid = grid.Next()
}
