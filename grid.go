package conway

import (
	"errors"
	"math/rand"
	"regexp"
	"strings"
)

// Cell is an (x, y) coordinate.
type Cell struct {
	X int
	Y int
}

// Grid is a 2-D grid of Cells.
type Grid map[Cell]struct{}

/*
 * GRID CREATION
 */

// ErrEmptyGrid is returned when creating a Grid that would have no Cells.
var ErrEmptyGrid = errors.New("grid must have at least one cell")

// FromSlice constructs a new Grid from a slice of slices of ints.
// Each nonzero value will be converted into a Cell of its index ([y][x]).
func FromSlice(rows [][]int) (Grid, error) {
	grid := make(Grid)
	for y, row := range rows {
		for x, val := range row {
			if val != 0 {
				grid.Add(Cell{x, y})
			}
		}
	}
	return requireNonEmpty(grid)
}

// GridSplitRegex is the pattern used to split a Grid string into rows.
var GridSplitRegex = regexp.MustCompile("[\n;]")

// LiveCellInputChar represents a live Cell in a Grid string.
var LiveCellInputChar = 'x'

var trimChars = []string{"\n", "\r", "\t"}

// FromString constructs a new Grid from a multiline string.
// Each line represents a row, and each occurrence of the rune 'x' is
// converted to a Cell in that position in the Grid.
func FromString(s string) (Grid, error) {
	grid := make(Grid)
	s = trimAny(s, trimChars)
	rows := GridSplitRegex.Split(s, -1)
	for y, row := range rows {
		row = trimAny(row, trimChars)
		for x, char := range row {
			if char == LiveCellInputChar {
				grid.Add(Cell{x, y})
			}
		}
	}
	return requireNonEmpty(grid)
}

func trimAny(s string, cutsets []string) string {
	for _, cutset := range cutsets {
		s = strings.Trim(s, cutset)
	}
	return s
}

// RandomGrid creates a Grid of random Cells in the given dimensions.
// `p` is the probability of a Cell being generated, from 0 to 1.
func RandomGrid(width, height int, p float64) (Grid, error) {
	grid := make(Grid)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if p > rand.Float64() {
				grid.Add(Cell{x, y})
			}
		}
	}
	return requireNonEmpty(grid)
}

func requireNonEmpty(grid Grid) (Grid, error) {
	if len(grid) == 0 {
		return nil, ErrEmptyGrid
	}
	return grid, nil
}

// Add adds a Cell to a Grid.
func (g Grid) Add(cell Cell) {
	g[cell] = struct{}{}
}

// AddMany adds one or more Cells to a Grid.
func (g Grid) AddMany(cells ...Cell) {
	for _, cell := range cells {
		g.Add(cell)
	}
}

// Remove removes a Cell from a Grid.
func (g Grid) Remove(cell Cell) {
	delete(g, cell)
}

/*
 * GAME EXECUTION
 */

// Next creates a new Grid by applying GoL rules.
// It returns the new Grid and an ok value. The ok value will be true if
// the Grid still has live Cells, or false if the Grid is empty.
func (g Grid) Next() (Grid, bool) {
	grid := make(Grid)
	for cell := range g.withNeighbors() {
		if g.cellSurvives(cell) {
			grid.Add(cell)
		} else {
			grid.Remove(cell)
		}
	}
	ok := len(grid) > 0
	return grid, ok
}

func (g Grid) withNeighbors() Grid {
	grid := make(Grid)
	for cell := range g {
		grid.Add(cell)
		grid.AddMany(cell.neighbors()...)
	}
	return grid
}

func (cell Cell) neighbors() []Cell {
	cells := make([]Cell, 8)
	i := 0
	for y := cell.Y - 1; y <= cell.Y+1; y++ {
		for x := cell.X - 1; x <= cell.X+1; x++ {
			c := Cell{x, y}
			if cell == c {
				continue
			}
			cells[i] = c
			i++
		}
	}
	return cells
}

func (g Grid) cellSurvives(cell Cell) bool {
	switch g.liveNeighbors(cell) {
	case 3:
		return true
	case 2:
		_, ok := g[cell]
		return ok
	default:
		return false
	}
}

func (g Grid) liveNeighbors(cell Cell) int {
	n := 0
	for _, c := range cell.neighbors() {
		_, ok := g[c]
		if ok {
			n++
		}
	}
	return n
}

/*
 * GRID VISUALIZATION
 */

// LiveCellRepr is the string used to represent a live Cell.
const LiveCellRepr = `■`

// DeadCellRepr is the string used to represent a dead Cell.
const DeadCellRepr = ` `

// Show returns a human-readable string representation of a Grid.
func (g Grid) Show() string {
	str := ""
	min, max := g.xyBounds()
	for y := min.Y; y <= max.Y; y++ {
		for x := min.X; x <= max.X; x++ {
			_, ok := g[Cell{x, y}]
			if ok {
				str += LiveCellRepr
			} else {
				str += DeadCellRepr
			}
		}
		str += "\n"
	}
	return str
}

func (g Grid) xyBounds() (min, max Cell) {
	for cell := range g {
		min.X = minimum(cell.X, min.X)
		min.Y = minimum(cell.Y, min.Y)
		max.X = maximum(cell.X, max.X)
		max.Y = maximum(cell.Y, max.Y)
	}
	return
}

func (g Grid) maxXY() (max Cell) {
	for cell := range g {
		max.X = maximum(max.X, cell.X)
		max.Y = maximum(max.Y, cell.Y)
	}
	return
}

func minimum(n0, n1 int) int {
	if n0 < n1 {
		return n0
	}
	return n1
}

func maximum(n0, n1 int) int {
	if n0 > n1 {
		return n0
	}
	return n1
}
