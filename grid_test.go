package main

import (
	"testing"

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
	require := require.New(t)
	grid := FromSlice([][]int{
		{1, 0, 0, 1, 0},
		{0, 1, 0, 0, 0},
		{1, 0, 0, 1, 0},
		{1, 1, 0, 1, 1},
		{0, 0, 0, 1, 1},
	})
	// 0 live neighbors dies
	require.False(grid.cellSurvives(Cell{3, 0}))
	// 1 live neighbor dies
	require.False(grid.cellSurvives(Cell{0, 0}))
	// 2 live neighbors lives
	require.True(grid.cellSurvives(Cell{1, 1}))
	require.True(grid.cellSurvives(Cell{3, 2}))
	require.True(grid.cellSurvives(Cell{1, 3}))
	// 3 live neighbors lives
	require.True(grid.cellSurvives(Cell{0, 2}))
	require.True(grid.cellSurvives(Cell{0, 3}))
	require.True(grid.cellSurvives(Cell{3, 4}))
	require.True(grid.cellSurvives(Cell{4, 4}))
	// 4+ live neighbors dies
	require.False(grid.cellSurvives(Cell{3, 3}))
	require.False(grid.cellSurvives(Cell{4, 3}))

	// 0-2 live neighbors stays dead
	require.False(grid.cellSurvives(Cell{1, 0}))
	require.False(grid.cellSurvives(Cell{2, 0}))
	require.False(grid.cellSurvives(Cell{4, 0}))
	require.False(grid.cellSurvives(Cell{3, 1}))
	require.False(grid.cellSurvives(Cell{4, 1}))
	require.False(grid.cellSurvives(Cell{0, 4}))
	require.False(grid.cellSurvives(Cell{1, 4}))
	// 3 live neighbors is revived
	require.True(grid.cellSurvives(Cell{0, 1}))
	require.True(grid.cellSurvives(Cell{2, 1}))
	require.True(grid.cellSurvives(Cell{4, 2}))
	require.True(grid.cellSurvives(Cell{2, 4}))
	// 4+ live neighbors stays dead
	require.False(grid.cellSurvives(Cell{1, 2}))
	require.False(grid.cellSurvives(Cell{2, 2}))
	require.False(grid.cellSurvives(Cell{2, 3}))
}
