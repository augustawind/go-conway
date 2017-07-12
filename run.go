package conway

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/buger/goterm"
	util "github.com/dustinrohde/go-conway/util/go-conway"
)

// RunConfig holds settings for running the simulation.
type RunConfig struct {
	OutFile      io.Writer
	Delay        time.Duration
	MaxTurns     int
	ClearScreen  bool // TODO
	Interactive  bool
	Spinner      bool // TODO
	KeepCentered bool // TODO
	FixedSize    int  // TODO
}

// DefaultRunConfig returns the default run settings.
func DefaultRunConfig() RunConfig {
	return RunConfig{
		OutFile:      os.Stdout,
		Delay:        500 * time.Millisecond,
		MaxTurns:     0,
		ClearScreen:  true,
		Interactive:  false,
		Spinner:      false,
		KeepCentered: false,
	}
}

var spinner = util.Spinner{
	Seq:   []string{"-", "/", "|", "\\"},
	Delay: 50 * time.Millisecond,
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
	if opts.ClearScreen {
		goterm.Clear()
		goterm.Flush()
	}
	fmt.Fprintln(opts.OutFile, grid.Show())

	done := make(chan bool)
	if opts.Spinner {
		defer spinner.Done()
		go spinner.Spin(done)
	}

	if opts.Interactive {
		util.WaitForInput(done)
	} else {
		time.Sleep(opts.Delay)
	}

	return grid.Next()
}
