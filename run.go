package conway

import (
	"fmt"
	"os"
	"time"
)

// RunConfig holds settings for running the simulation.
type RunConfig struct {
	Prompt   bool
	MaxTurns int
	Delay    time.Duration
	Outfile  *os.File
}

// DefaultRunConfig returns the default run settings.
func DefaultRunConfig() RunConfig {
	return RunConfig{
		Prompt:   false,
		MaxTurns: 0,
		Delay:    time.Millisecond * 500,
		Outfile:  os.Stdout,
	}
}

// Run runs a Game of Life simulation.
func Run(grid Grid, opts RunConfig) {
	if opts.Outfile != os.Stdout {
		defer opts.Outfile.Close
	}
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

// RunDefault calls Run with the default settings.
func RunDefault(grid Grid) {
	Run(grid, DefaultRunConfig())
}

// NextTurn runs a single turn of a Game of Life simulation.
func NextTurn(grid Grid, opts RunConfig) Grid {
	fmt.Fprintln(opts.Outfile, grid.Show())
	time.Sleep(opts.Delay)
	return grid.Next()
}
