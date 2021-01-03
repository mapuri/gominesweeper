package main

import (
	"math/rand"
	"time"
)

// Board capture the runtime state of the game
type Board struct {
	Rows  int
	Cols  int
	Cells [][]Cell
	State GameState

	minedCells   int
	openedCells  int
	flaggedCells int
}

// NewBoard initializes and returns a game board.
func NewBoard(lvl Level) *Board {
	b := &Board{}
	b.init(lvl)
	return b
}

// init initializes a board based on user specified difficulty
func (b *Board) init(lvl Level) {
	rows := 6
	cols := 6
	numMines := 6
	switch lvl {
	case Medium:
		rows = 10
		cols = 10
		numMines = 10
	case Hard:
		rows = 20
		cols = 20
		numMines = 20
	}
	b.Rows = rows
	b.Cols = cols
	b.minedCells = numMines
	b.Cells = make([][]Cell, rows)
	for i := 0; i < rows; i++ {
		b.Cells[i] = make([]Cell, cols)
	}

	// randomly pick the cells with mines
	rand.Seed(int64(time.Now().Nanosecond()))
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
		b.Cells[row][col].Val = Mine
	}

	// setup the numbered cells
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			b.Cells[row][col].Row = row
			b.Cells[row][col].Col = col
			if b.Cells[row][col].Val.IsMine() {
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
				if b.Cells[neigh.row][neigh.col].Val.IsMine() {
					numMines++
				}
			}
			if numMines > 0 {
				b.Cells[row][col].Val = Value(numMines)
			}
		}
	}
}

// OpenCell marks the cell as open. When a mine is opened the game ends. When a clear cell is opened it may open more adjoimg cells.
// When a number cells is opend it just opens and reveals the value of that number
func (b *Board) OpenCell(row, col int) {
	if b.Cells[row][col].State == Opened {
		// noop, if cell is already opened
		return
	}
	defer func() {
		if b.State == Lost || b.State == Won {
			// once the game completes, reveal the entire board
			for row := 0; row < b.Rows; row++ {
				for col := 0; col < b.Cols; col++ {
					b.Cells[row][col].State = Opened
				}
			}
		}
	}()

	if b.Cells[row][col].Val.IsMine() {
		b.State = Lost
		return
	}
	defer func() {
		// check if we are done
		if b.flaggedCells == b.minedCells && b.openedCells+b.flaggedCells == b.Rows*b.Cols {
			b.State = Won
		}
	}()
	if b.Cells[row][col].Val.IsClear() {
		b.openIsland(row, col)

		return
	}
	b.Cells[row][col].State = Opened
	b.openedCells++
	return
}

// openIsland is called from openCell, when the cell being opened is clear. It performs a BFS to search all adjoinging cells that
// are either clear or numbered
func (b *Board) openIsland(row, col int) {
	type cell struct{ row, col int }
	queued := make(map[cell]struct{})
	q := make([]cell, 0)
	c := cell{row, col}
	q = append(q, c)
	queued[c] = struct{}{}
	for len(q) > 0 {
		c := q[0]
		q = q[1:]
		b.Cells[c.row][c.col].State = Opened
		b.openedCells++
		if b.Cells[c.row][c.col].Val.IsNumber() {
			// stop the search at numbered cells
			continue
		}
		neighs := []cell{
			{c.row - 1, c.col - 1}, {c.row - 1, c.col}, {c.row - 1, c.col + 1},
			{c.row, c.col - 1}, {c.row, c.col + 1},
			{c.row + 1, c.col - 1}, {c.row + 1, c.col}, {c.row + 1, c.col + 1},
		}
		for _, neigh := range neighs {
			if neigh.row < 0 || neigh.row >= b.Rows || neigh.col < 0 || neigh.col >= b.Cols {
				continue
			}
			if _, ok := queued[neigh]; ok {
				continue
			}
			if !b.Cells[neigh.row][neigh.col].Val.IsMine() {
				// only add numbered and clear cells to queue
				q = append(q, neigh)
				queued[neigh] = struct{}{}
			}
		}
	}
}

// FlagCell marks the cell as flagged.
func (b *Board) FlagCell(row, col int) {
	if b.Cells[row][col].State == Opened {
		// noop, if cell is already opened
		return
	}

	if b.Cells[row][col].State == Closed {
		b.Cells[row][col].State = Flagged
		b.flaggedCells++
	} else {
		b.Cells[row][col].State = Closed
		b.flaggedCells--
	}
}

// Level is the game Level as selected by the user
type Level int

const (
	// Easy is a easy game level
	Easy Level = iota

	// Medium is a medium game level
	Medium

	// Hard is hard game level
	Hard
)

// Cell is a single Cell on the board
type Cell struct {
	Val      Value
	State    State
	Row, Col int // used to render the html table template
}

// Value is the Value of cell initialized when board is setup
type Value int

const (
	// Clear denotes a cell that is empty
	Clear = 0

	//Mine denotes a cell that contains a mine
	Mine = 9
)

// IsMine returns true if cell contains a mine
func (v Value) IsMine() bool {
	return v == Mine
}

// IsClear returns true if cell is empty
func (v Value) IsClear() bool {
	return v == Clear
}

// IsNumber returns true is cell is a empty cell that is surrounded by 1 or more mines
func (v Value) IsNumber() bool {
	return v != Clear && v != Mine
}

// State is the runtime State of the cell. It can change by user actions.
type State int

const (
	// Closed denotes a closed cell
	Closed State = iota

	// Opened denotes a cell that has been opened
	Opened

	// Flagged denotes a cell that has been flagged
	Flagged
)

// GameState is the current state of the game
type GameState int

const (
	// Playing denotes a game in progress
	Playing GameState = iota

	// Won denotes a game that has been won
	Won

	// Lost denotes a game that has been lost
	Lost
)
