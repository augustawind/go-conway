package conway

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// RunConfig holds settings for running the simulation.
type RunConfig struct {
	GridFile    io.Reader
	OutFile     io.Writer
	Delay       time.Duration
	MaxTurns    int
	Interactive bool
}

var defaultGridStr = strings.Join(
	[]string{
		".....",
		"..x..",
		"...x.",
		".xxx.",
		".....",
	},
	"\n",
)

// DefaultRunConfig returns the default run settings.
func DefaultRunConfig() RunConfig {
	return RunConfig{
		GridFile:    strings.NewReader(defaultGridStr),
		OutFile:     os.Stdout,
		Delay:       time.Millisecond * 500,
		MaxTurns:    0,
		Interactive: false,
	}
}

// Run runs a Game of Life simulation.
func Run(grid Grid, opts RunConfig) {
	ok := true
	if opts.MaxTurns > 0 {
		for i := 0; i < opts.MaxTurns && ok; i++ {
			grid, ok = NextTurn(grid, opts)
		}
	} else {
		for ok {
			grid, ok = NextTurn(grid, opts)
		}
	}
}

// RunDefault calls Run with the default settings.
func RunDefault(grid Grid) {
	Run(grid, DefaultRunConfig())
}

// NextTurn runs a single turn of a Game of Life simulation.
func NextTurn(grid Grid, opts RunConfig) (Grid, bool) {
	fmt.Fprintln(opts.OutFile, grid.Show())
	time.Sleep(opts.Delay)
	return grid.Next()
}
