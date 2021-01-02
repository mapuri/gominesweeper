package main

import (
	"bytes"
	"fmt"
	"html/template"
	"syscall/js"
)

var singleton *board

var boardTemplate = `
<style>
#tableHeader {
	background-color: tomato;
	  color: black;
	  padding: 10px;
	  text-align: center;
}
.cell {
	background-color: darkseagreen;
}
</style>
<table>
	<tr>
		<th id=tableHeader colspan={{.cols}}>
		<select id=gameLevel onchange=resetGame()>
			<option value=0>Easy</option>
			<option value=1 selected>Medium</option>
			<option value=2>Hard</option>
		</select>
		</th>
	</tr>
	<tr>
		<td>
			<p id=gameStatus></p>
		</td>
		<td>
			<button id=restartGame onclick=resetGame() disabled> Play Again </button>
		</td>
	</tr>
	{{ range .cells }}
	<tr>
	{{ range .}}
	<td class=cell id=cell{{.row}}-{{.col}} onclick=openCell({{.row}}, {{.col}}) oncontextmenu=flagCell({{.row}}, {{.col}})></td>
	{{ end}}
	</tr>
	{{ end }}
</table>
`

// initGame initializes the board to specified level
func initGame(lvl js.Value) {
	lvlVal := lvl.Int()
	singleton = newBoard(level(lvlVal))
	gameSelection := js.Global().Get("document").Call("getElementById", "gameSection")
	tmplt, err := template.New("table-template").Parse(boardTemplate)
	if err != nil {
		retError := fmt.Errorf("error parsing template: %w", err)
		gameSelection.Set("innerHtml", retError.Error())
		return
	}
	var buf bytes.Buffer
	err = tmplt.Execute(&buf, singleton)
	if err != nil {
		retError := fmt.Errorf("error executing the template: %w", err)
		gameSelection.Set("innerHtml", retError.Error())
		return
	}
	gameSelection.Set("innerHtml", buf.String())
}

// resetGame resets the game to user selected level.
func resetGame() {
	lvl := js.Global().Get("document").Call("getElementById", "gameLevel").Get("value")
	initGame(lvl)
}

// openCell triggers action to open a closed cell
func openCell(row, col js.Value) {
	rVal := row.Int()
	cVal := col.Int()
	singleton.openCell(rVal, cVal)
	document := js.Global().Get("document")
	switch singleton.state {
	case won:
		document.Call("getElementById", "gameStatus").Set("innerHTML", "Yay, you won!")
		document.Call("getElementById", "restartGame").Set("disabled", "")
	case lost:
		document.Call("getElementById", "gameStatus").Set("innerHTML", "Sorry, you lost!")
		document.Call("getElementById", "restartGame").Set("disabled", "")
	}
	drawBoard(singleton)
}

// flagCell triggers action to flag a closed cell
func flagCell(row, col js.Value) {
	rVal := row.Int()
	cVal := col.Int()
	singleton.openCell(rVal, cVal)
	drawBoard(singleton)
}

// drawBoard redraws the board
func drawBoard(b *board) {
	for row := 0; row < b.rows; row++ {
		for col := 0; col < b.cols; col++ {
			cell := b.cells[row][col]
			cellName := fmt.Sprintf("cell%d-%d", row, col)
			cellTag := js.Global().Get("document").Call("getElementById", cellName)
			if cell.state == opened {
				val := ""
				if cell.val.isMine() {
					val = "💣"
				}
				if cell.val.isNumber() {
					val = fmt.Sprint(cell.val)
				}
				cellTag.Set("innerHTML", val)
			}
			if cell.state == flagged {
				cellTag.Set("innerHTML", "🚩")
			}
		}
	}
}
