package conway

import (
	"errors"
	"math/rand"
	"regexp"
	"strings"
)

// Grid is a 2-D grid of Cells.
type Grid map[Cell]CellState

// Cell is an (x, y) coordinate.
type Cell struct {
	X     int
	Y     int
}

// CellState represents the binary status of a Cell (alive/dead).
type CellState bool

// Alive and Dead are assigned to Cell.State to indicate whether the Cell is alive or dead.
const (
	Alive CellState = true
	Dead  CellState = false
)

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
			grid.Set(Cell{x, y}, val != 0)
		}
	}
	return requireNonEmpty(grid)
}

// GridSplitRegex is the pattern used to split a Grid string into rows.
var GridSplitRegex = regexp.MustCompile("[\n;]")

const (
	TokenLiveCell = 'x'
	TokenDeadCell = '.'
	TokenComment  = '#'
)

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
			switch char {
			case TokenLiveCell:
				grid.Set(Cell{x, y}, Alive)
			case TokenDeadCell:
				grid.Set(Cell{x, y}, Dead)
			case TokenComment:
				continue
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
	if width < 1 || height < 1 || p == 0.0 {
		return nil, ErrEmptyGrid
	}
	grid := make(Grid)
	for len(grid) == 0 {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				state := p > rand.Float64()
				grid.Set(Cell{x, y}, CellState(state))
			}
		}
	}
	return grid, nil
}

func requireNonEmpty(grid Grid) (Grid, error) {
	if len(grid) == 0 {
		return nil, ErrEmptyGrid
	}
	return grid, nil
}

// Set sets a Cell in a Grid to alive or dead.
func (g Grid) Set(cell Cell, state CellState) {
	g[cell] = state
}

// Add is like Set but it will only update a Cell if it already exists.
func (g Grid) Add(cell Cell, state CellState) {
	_, ok := g[cell]
	if !ok {
		g.Set(cell, state)
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
	nextGrid := make(Grid)
	for cell := range g.withNeighbors() {
		nextGrid.Set(cell, g.nextCell(cell))
	}
	ok := len(nextGrid) > 0
	return nextGrid, ok
}

func (g Grid) withNeighbors() Grid {
	grid := make(Grid)
	for cell := range g {
		grid.Add(cell, Dead)
		for _, c := range cell.neighbors() {
			grid.Add(c, Dead)
		}
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

func (g Grid) nextCell(cell Cell) CellState {
	switch g.liveNeighbors(cell) {
	case 3:
		return Alive
	case 2:
		return g[cell]
	default:
		return Dead
	}
}

func (g Grid) liveNeighbors(cell Cell) int {
	n := 0
	for _, c := range cell.neighbors() {
		if g[c] == Alive {
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
			if g[Cell{x, y}] == Alive {
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
