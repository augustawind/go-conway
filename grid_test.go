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
