package main

// LiveCellRepr is the string used to represent a live Cell.
const LiveCellRepr = `■`

// DeadCellRepr is the string used to represent a dead Cell.
const DeadCellRepr = ` `

// Cell is an (x, y) coordinate.
type Cell struct {
	X int
	Y int
}

// Grid is a 2-D grid of Cells.
type Grid map[Cell]struct{}

// FromSlice constructs a new Grid from a slice of slices of ints.
// Each nonzero value will be converted into a Cell of its index ([y][x]).
func FromSlice(rows [][]int) Grid {
	grid := make(Grid)
	for y, row := range rows {
		for x, val := range row {
			if val != 0 {
				grid.Add(Cell{x, y})
			}
		}
	}
	return grid
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

// Next creates a new Grid by applying GoL rules.
func (g Grid) Next() Grid {
	grid := make(Grid)
	for cell := range g.withNeighbors() {
		if g.cellSurvives(cell) {
			grid.Add(cell)
		} else {
			grid.Remove(cell)
		}
	}
	return grid
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

// Show returns a human-readable string representation of a Grid.
func (g Grid) Show() string {
	str := ""
	maxX, maxY := g.maxXY()
	for y := 0; y <= maxY; y++ {
		for x := 0; x <= maxX; x++ {
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

func (g Grid) maxXY() (maxX, maxY int) {
	for cell := range g {
		if cell.X > maxX {
			maxX = cell.X
		}
		if cell.Y > maxY {
			maxY = cell.Y
		}
	}
	return
}
