package main

import (
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dustinrohde/go-conway"
	"github.com/dustinrohde/go-conway/util/go-conway"
	"github.com/urfave/cli"
)

func main() {
	app := initApp()
	app.Run(os.Args)
}

func initApp() *cli.App {
	var grid conway.Grid
	var err error
	config := conway.DefaultRunConfig()

	app := cli.NewApp()
	app.Name = "conway"
	app.Usage = "Run Conway's Game of Life simulation."

	cli.HelpFlag = cli.BoolFlag{
		Name: "help",
		Usage: "show help",
	}

	app.Commands = []cli.Command{
		{
			Name:    "defined",
			Aliases: []string{"def"},
			Usage:   "start with a predefined grid",
			Flags: withCommonFlags(config,
				cli.StringFlag{
					Name:  "grid, g",
					Value: "-",
					Usage: strings.Join([]string{
						"Starting grid.",
						"\t\tIf `FILE` starts with '@', interpret it as a file path.",
						"\t\tIf `FILE` is '-', read from stdout.",
					}, "\n"),
				},
			),
			Action: func(c *cli.Context) error {
				arg := c.String("grid")
				if arg == "-" {
					// Read Grid from stdin.
					grid = readGrid(os.Stdin)
				} else if arg[0] == '@' {
					// Read Grid from file.
					file, err := os.Open(arg[1:])
					util.Guard(err)
					defer file.Close()
					grid = readGrid(file)
				} else {
					// Read Grid from argument.
					grid, err = conway.FromString(arg)
					util.Guard(err)
				}
				setOutFile(config, c.String("outfile"))
				conway.Run(grid, config)
				return nil
			},
		},
		{
			Name:    "random",
			Aliases: []string{"rand"},
			Usage:   "start with a randomly generated grid",
			Flags: withCommonFlags(config,
				cli.IntFlag{
					Name:  "width, w",
					Value: 9,
					Usage: "max `WIDTH` of the grid",
				},
				cli.IntFlag{
					Name:  "height, h",
					Value: 9,
					Usage: "max `HEIGHT` of the grid",
				},
				cli.Float64Flag{
					Name:  "probability, p",
					Value: 0.5,
					Usage: "probability of living cells, where 0 < `PROB` <= 1",
				},
			),
			Action: func(c *cli.Context) error {
				grid, err = conway.RandomGrid(
					c.Int("width"), c.Int("height"), c.Float64("probability"))
				util.Guard(err)
				setOutFile(config, c.String("outfile"))
				conway.Run(grid, config)
				return nil
			},
		},
	}

	return app
}

func withCommonFlags(config conway.RunConfig, flags ...cli.Flag) []cli.Flag {
	 commonFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "outfile, o",
			Value: "-",
			Usage: "File to write results to. If `FILE` is `-`, use stdout.",
		},
		cli.DurationFlag{
			Name:        "delay, d",
			Value:       config.Delay,
			Usage:       "`TIME` to pause between turns; e.g. 500ms, 3s, 1m",
			Destination: &config.Delay,
		},
		cli.IntFlag{
			Name:        "turns, t",
			Value:       config.MaxTurns,
			Usage:       "Max `TURNS` to run. If < 0, run indefinitely.",
			Destination: &config.MaxTurns,
		},
		cli.BoolFlag{
			Name:        "clear, c",
			Usage:       "Clear screen between each turn.",
			Destination: &config.ClearScreen,
		},
		cli.BoolFlag{
			Name:        "interactive, i",
			Usage:       "Wait for input between each turn.",
			Destination: &config.Interactive,
		},
		cli.BoolFlag{
			Name:        "spinner, s",
			Usage:       "Show an animated spinner between turns.",
			Destination: &config.Spinner,
		},
	}
	return append(flags, commonFlags...)
}


func setOutFile(config conway.RunConfig, path string) {
	if path != "-" {
		file, err := os.Create(path)
		util.Guard(err)
		config.OutFile = file
	}
}

func readGrid(r io.Reader) (grid conway.Grid) {
	gridBytes, err := ioutil.ReadAll(r)
	util.Guard(err)
	grid, err = conway.FromString(string(gridBytes))
	util.Guard(err)
	return
}
