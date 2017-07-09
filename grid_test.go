package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromSlice(t *testing.T) {
	require := require.New(t)
	actual := FromSlice([][]int{
		{1, 0, 0},
		{0, 0, 0},
		{0, 1, 1},
	})
	expected := make(Grid)
	expected[Cell{0, 0}] = struct{}{}
	expected[Cell{1, 2}] = struct{}{}
	expected[Cell{2, 2}] = struct{}{}
	require.Equal(expected, actual)
}

func TestCell_neighbors(t *testing.T) {
	require := require.New(t)
	cell := Cell{0, 2}
	actual := cell.neighbors()
	expected := []Cell{
		Cell{-1, 1}, Cell{0, 1}, Cell{1, 1},
		Cell{-1, 2}, Cell{1, 2},
		Cell{-1, 3}, Cell{0, 3}, Cell{1, 3},
	}
	require.Equal(expected, actual)
}

func TestGrid_liveNeighbors(t *testing.T) {
	require := require.New(t)
	grid := FromSlice([][]int{
		{1, 0, 0, 0},
		{1, 1, 1, 0},
		{0, 0, 1, 1},
		{0, 0, 0, 0},
	})
	require.Equal(2, grid.liveNeighbors(Cell{0, 0}))
	require.Equal(3, grid.liveNeighbors(Cell{2, 2}))
}

func TestGrid_cellSurvives(t *testing.T) {
	assert := assert.New(t)
	grid := FromSlice([][]int{
		{1, 0, 0, 1, 0},
		{0, 1, 0, 0, 0},
		{1, 0, 0, 1, 0},
		{1, 1, 0, 1, 1},
		{0, 0, 0, 1, 1},
	})
	// 0 live neighbors dies
	assert.False(grid.cellSurvives(Cell{3, 0}))
	// 1 live neighbor dies
	assert.False(grid.cellSurvives(Cell{0, 0}))
	// 2 live neighbors lives
	assert.True(grid.cellSurvives(Cell{1, 1}))
	assert.True(grid.cellSurvives(Cell{3, 2}))
	assert.True(grid.cellSurvives(Cell{1, 3}))
	// 3 live neighbors lives
	assert.True(grid.cellSurvives(Cell{0, 2}))
	assert.True(grid.cellSurvives(Cell{0, 3}))
	assert.True(grid.cellSurvives(Cell{3, 4}))
	assert.True(grid.cellSurvives(Cell{4, 4}))
	// 4+ live neighbors dies
	assert.False(grid.cellSurvives(Cell{3, 3}))
	assert.False(grid.cellSurvives(Cell{4, 3}))

	// 0-2 live neighbors stays dead
	assert.False(grid.cellSurvives(Cell{1, 0}))
	assert.False(grid.cellSurvives(Cell{2, 0}))
	assert.False(grid.cellSurvives(Cell{4, 0}))
	assert.False(grid.cellSurvives(Cell{3, 1}))
	assert.False(grid.cellSurvives(Cell{4, 1}))
	assert.False(grid.cellSurvives(Cell{0, 4}))
	assert.False(grid.cellSurvives(Cell{1, 4}))
	// 3 live neighbors is revived
	assert.True(grid.cellSurvives(Cell{0, 1}))
	assert.True(grid.cellSurvives(Cell{2, 1}))
	assert.True(grid.cellSurvives(Cell{4, 2}))
	assert.True(grid.cellSurvives(Cell{2, 4}))
	// 4+ live neighbors stays dead
	assert.False(grid.cellSurvives(Cell{1, 2}))
	assert.False(grid.cellSurvives(Cell{2, 2}))
	assert.False(grid.cellSurvives(Cell{2, 3}))
}

func TestGrid_withNeighbors(t *testing.T) {
	require := require.New(t)
	grid := FromSlice([][]int{
		{1, 0},
		{0, 1},
	})
	actual := grid.withNeighbors()
	expected := make(Grid)
	expected.AddMany(
		Cell{-1, -1}, Cell{0, -1}, Cell{1, -1},
		Cell{-1, 0}, Cell{0, 0}, Cell{1, 0}, Cell{2, 0},
		Cell{-1, 1}, Cell{0, 1}, Cell{1, 1}, Cell{2, 1},
		Cell{0, 2}, Cell{1, 2}, Cell{2, 2},
	)
	require.Equal(expected, actual)
}

func TestGrid_Next(t *testing.T) {
	require := require.New(t)
	type gridPair struct {
		start Grid
		next  Grid
	}
	pairs := []gridPair{
		{
			start: FromSlice([][]int{
				{0, 1, 0},
				{0, 1, 0},
				{0, 1, 0},
			}),
			next: FromSlice([][]int{
				{0, 0, 0},
				{1, 1, 1},
				{0, 0, 0},
			}),
		},
		{
			start: FromSlice([][]int{
				{0, 0, 0, 0},
				{0, 1, 1, 1},
				{1, 1, 1, 0},
				{0, 0, 0, 0},
			}),
			next: FromSlice([][]int{
				{0, 0, 1, 0},
				{1, 0, 0, 1},
				{1, 0, 0, 1},
				{0, 1, 0, 0},
			}),
		},
		{
			start: FromSlice([][]int{
				{0, 0, 0, 0},
				{0, 1, 1, 0},
				{0, 1, 1, 0},
				{0, 0, 0, 0},
			}),
			next: FromSlice([][]int{
				{0, 0, 0, 0},
				{0, 1, 1, 0},
				{0, 1, 1, 0},
				{0, 0, 0, 0},
			}),
		},
	}

	for _, pair := range pairs {
		require.Equal(pair.next, pair.start.Next())
	}
}

func TestGrid_maxXY(t *testing.T) {
	require := require.New(t)
	grid := FromSlice([][]int{
		{1, 0, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 1, 0},
		{0, 1, 1, 0},
	})
	maxX, maxY := grid.maxXY()
	require.Equal(2, maxX)
	require.Equal(3, maxY)
}

func TestGrid_Show(t *testing.T) {
	require := require.New(t)
	grid := FromSlice([][]int{
		{1, 0, 0},
		{0, 0, 0},
		{0, 1, 1},
	})
	actual := grid.Show()
	x := LiveCellRepr
	o := DeadCellRepr
	expected := strings.Join(
		[]string{
			x + o + o,
			o + o + o,
			o + x + x,
		},
		"\n",
	)
	require.Equal(strings.TrimSpace(expected), strings.TrimSpace(actual))
}
