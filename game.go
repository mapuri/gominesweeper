package main

import (
	"math/rand"
)

// board capture the runtime state of the game
type board struct {
	rows         int
	cols         int
	cells        [][]cell
	state        gameState
	openedCells  int
	flaggedCells int
}

// newBoard initializes and returns a game board.
func newBoard(lvl level) *board {
	b := &board{}
	b.init(lvl)
	return b
}

// init initializes a board based on user specified difficulty
func (b *board) init(lvl level) {
	rows := 6
	cols := 6
	numMines := 6
	switch lvl {
	case medium:
		rows = 10
		cols = 10
		numMines = 10
	case hard:
		rows = 20
		cols = 20
		numMines = 20
	}
	b.rows = rows
	b.cols = cols
	b.cells = make([][]cell, rows)
	for i := 0; i < rows; i++ {
		b.cells[i] = make([]cell, cols)
	}

	// randomly pick the cells with mines
	mineCells := make(map[int]struct{ row, col int })
	for len(mineCells) < numMines {
		key := rand.Intn(rows * cols)
		if _, ok := mineCells[key]; !ok {
			row := key / rows
			col := key % cols
			mineCells[key] = struct {
				row int
				col int
			}{row: row, col: col}
		}
	}
	for _, cell := range mineCells {
		row := cell.row
		col := cell.col
		b.cells[row][col].val = mine
	}

	// setup the numbered cells
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			if b.cells[row][col].val.isMine() {
				continue
			}
			neighs := []struct{ row, col int }{
				{row - 1, col - 1}, {row - 1, col}, {row - 1, col + 1},
				{row, col - 1}, {row, col + 1},
				{row + 1, col - 1}, {row + 1, col}, {row + 1, col + 1},
			}
			numMines := 0
			for _, neigh := range neighs {
				if neigh.row < 0 || neigh.row >= rows || neigh.col < 0 || neigh.col >= cols {
					continue
				}
				if b.cells[neigh.row][neigh.col].val.isMine() {
					numMines++
				}
			}
			if numMines > 0 {
				b.cells[row][col].val = value(numMines)
			}
		}
	}
}

// openCell marks the cell as open. When a mine is opened the game ends. When a clear cell is opened it may open more adjoimg cells.
// When a number cells is opend it just opens and reveals the value of that number
func (b *board) openCell(row, col int) {
	if b.cells[row][col].state == opened {
		// noop, if cell is already opened
		return
	}

	if b.cells[row][col].val.isMine() {
		b.state = lost
		return
	}
	defer func() {
		// check if we are done
		if b.openedCells+b.flaggedCells == b.rows*b.cols {
			b.state = won
		}
	}()
	if b.cells[row][col].val.isClear() {
		b.openIsland(row, col)

		return
	}
	b.cells[row][col].state = opened
	b.openedCells++
	return
}

// openIsland is called from openCell, when the cell being opened is clear. It performs a BFS to search all adjoinging cells that
// are either clear or numbered
func (b *board) openIsland(row, col int) {
	type cell struct{ row, col int }
	visited := make(map[cell]struct{})
	q := make([]cell, 0)
	q = append(q, cell{row, col})
	for len(q) > 0 {
		c := q[0]
		q = q[1:]
		b.cells[c.row][c.col].state = opened
		b.openedCells++
		visited[c] = struct{}{}
		if b.cells[c.row][c.col].val.isNumber() {
			// stop the search at numbered cells
			continue
		}
		neighs := []cell{
			{c.row - 1, c.col - 1}, {c.row - 1, c.col}, {c.row - 1, c.col + 1},
			{c.row, c.col - 1}, {c.row, c.col + 1},
			{c.row + 1, c.col - 1}, {c.row + 1, c.col}, {c.row + 1, c.col + 1},
		}
		for _, neigh := range neighs {
			if neigh.row < 0 || neigh.row >= b.rows || neigh.col < 0 || neigh.col >= b.cols {
				continue
			}
			if _, ok := visited[neigh]; ok {
				continue
			}
			if !b.cells[neigh.row][neigh.col].val.isMine() {
				// only add numbered and clear cells to queue
				q = append(q, neigh)
			}
		}
	}
}

// flagCell marks the cell as flagged.
func (b *board) flagCell(row, col int) {
	if b.cells[row][col].state == opened {
		// noop, if cell is already opened
		return
	}

	b.cells[row][col].state = flagged
}

// level is the game level as selected by the user
type level int

const (
	easy level = iota
	medium
	hard
)

// cell is a single cell on the board
type cell struct {
	val   value
	state state
}

// value is the value of cell initialized when board is setup
type value int

const (
	clear = 0
	mine  = 9
)

func (v value) isMine() bool {
	return v == mine
}

func (v value) isClear() bool {
	return v == clear
}

func (v value) isNumber() bool {
	return v != clear && v != mine
}

// state is the runtime state of the cell. It can change by user actions.
type state int

const (
	closed state = iota
	opened
	flagged
)

// gameState is the current state of the game
type gameState int

const (
	playing gameState = iota
	won
	lost
)
