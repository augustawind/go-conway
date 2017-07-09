package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dustinrohde/go-conway"
	"github.com/urfave/cli"
)

func main() {
	app := initApp()
	app.Run(os.Args)
}

func initApp() *cli.App {
	app := cli.NewApp()
	app.Name = "conway"
	app.Usage = "Run Conway's Game of Life simulation."

	var grid conway.Grid
	config := conway.DefaultRunConfig()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "grid, g",
			Value: "",
			Usage: "Starting grid. If `FILE` starts with `@`, interpret it as" +
				" a file path. If `FILE` is `-`, read from stdout. If `FILE`" +
				" is absent or blank, use a demo starting grid.",
		},
		cli.StringFlag{
			Name:  "outfile, o",
			Value: "-",
			Usage: "File to write results to. If `-`, use stdout.",
		},
		cli.DurationFlag{
			Name:        "delay, d",
			Value:       config.Delay,
			Usage:       "Delay between turns; e.g. 500ms, 3s, 1m",
			Destination: &config.Delay,
		},
		cli.IntFlag{
			Name:        "turns, t",
			Value:       config.MaxTurns,
			Usage:       "Number of turns to run. If < 0, run forever.",
			Destination: &config.MaxTurns,
		},
		cli.BoolFlag{
			Name:        "interactive, i",
			Usage:       "Wait for input between each turn.",
			Destination: &config.Interactive,
		},
	}

	app.Action = func(c *cli.Context) error {
		if path := c.String("outfile"); path != "-" {
			// Write results to a file.
			file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			Guard(err)
			config.OutFile = file
			defer file.Close()
		}

		if path := c.String("grid"); path == "-" {
			// Read Grid from stdin.
			config.GridFile = os.Stdin
		} else if path[0] == '@' {
			// Read Grid from file.
			file, err := os.OpenFile(path[1:], os.O_RDONLY, 0600)
			Guard(err)
			config.GridFile = file
			defer file.Close()
		}
		gridBytes, err := ioutil.ReadAll(config.GridFile)
		Guard(err)
		grid = conway.FromString(string(gridBytes))

		return nil
	}

	return app
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
