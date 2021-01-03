package main

import (
	"bytes"
	"fmt"
	"html/template"
	"strconv"
	"syscall/js"
)

var singleton *Board

var boardTemplate = `
<style>
#board {
	margin-left: auto;
	margin-right: auto;
}
#tableHeader {
	background-color: tomato;
	  color: black;
	  padding: 10px;
	  text-align: center;
}
#statusLine {
	width:100%;
	text-align: right;
}
.statusInner {
	display: inline-block;
}
.closedCell {
	background-color: darkseagreen;
	color: black;
	text-align: center;
	width: 30px;
	height: 30px;
}
.openedCell {
	background-color: lightseagreen;
	color: black;
	text-align: center;
	width: 30px;
	height: 30px;
}
.flaggedCell {
	background-color: lightpink;
	color: black;
	text-align: center;
	width: 30px;
	height: 30px;
}
</style>
<table id=board>
	<tr>
		<th id=tableHeader colspan={{.Cols}}>
		<select id=gameLevel onchange=resetGame()>
			<option value=0>Easy</option>
			<option value=1 selected>Medium</option>
			<option value=2>Hard</option>
		</select>
		</th>
	</tr>
	<tr>
		<td colspan={{.Cols}}>
			<div id=statusLine>
				<div class=statusInner sty> <p id=gameStatus style='font-weight: bold;'> Game in progress </p> </div>
				<div class=statusInner> <button id=restartGame onclick=resetGame() disabled> Play Again </button> </div>
			</div>
		</td>
	</tr>
	{{ range $row, $cells := .Cells }}
	<tr>
	{{ range $col, $cell := $cells}}
	<td class=closedCell id="cell{{$row}}-{{$col}}" onclick="openCell( {{$row}}, {{$col}} )" oncontextmenu="event.preventDefault(); flagCell( {{$row}}, {{$col}} )"></td>
	{{ end}}
	</tr>
	{{ end }}
</table>
`

// initGame initializes the board to specified level
func initGame(this js.Value, args []js.Value) interface{} {
	gameSection := js.Global().Get("document").Call("getElementById", "gameSection")
	lvlVal, err := strconv.Atoi(args[0].String())
	if err != nil {
		retError := fmt.Errorf("invalid level value: %w", err)
		gameSection.Set("innerHTML", retError.Error())
		return nil
	}
	singleton = NewBoard(Level(lvlVal))
	tmplt, err := template.New("table-template").Parse(boardTemplate)
	if err != nil {
		retError := fmt.Errorf("error parsing template: %w", err)
		gameSection.Set("innerHTML", retError.Error())
		return nil
	}
	var buf bytes.Buffer
	err = tmplt.Execute(&buf, singleton)
	if err != nil {
		retError := fmt.Errorf("error executing the template: %w", err)
		gameSection.Set("innerHTML", retError.Error())
		return nil
	}
	gameSection.Set("innerHTML", buf.String())
	return nil
}

// resetGame resets the game to user selected level.
func resetGame(this js.Value, args []js.Value) interface{} {
	lvl := js.Global().Get("document").Call("getElementById", "gameLevel").Get("value")
	initGame(js.Null(), []js.Value{lvl})
	// need to fetch gameLevel again as the board HTML is setup afresh on changing the level
	js.Global().Get("document").Call("getElementById", "gameLevel").Set("value", lvl)
	return nil
}

// openCell triggers action to open a closed cell
func openCell(this js.Value, args []js.Value) interface{} {
	rVal := args[0].Int()
	cVal := args[1].Int()
	singleton.OpenCell(rVal, cVal)
	document := js.Global().Get("document")
	switch singleton.State {
	case Won:
		document.Call("getElementById", "gameStatus").Set("innerHTML", "Yay, you won!")
		document.Call("getElementById", "restartGame").Set("disabled", false)
	case Lost:
		document.Call("getElementById", "gameStatus").Set("innerHTML", "Sorry, you lost!")
		document.Call("getElementById", "restartGame").Set("disabled", false)
	}
	drawBoard(singleton)
	return nil
}

// flagCell triggers action to flag a closed cell
func flagCell(this js.Value, args []js.Value) interface{} {
	rVal := args[0].Int()
	cVal := args[1].Int()
	singleton.FlagCell(rVal, cVal)
	drawBoard(singleton)
	return nil
}

// drawBoard redraws the board
func drawBoard(b *Board) {
	for row := 0; row < b.Rows; row++ {
		for col := 0; col < b.Cols; col++ {
			cell := b.Cells[row][col]
			cellName := fmt.Sprintf("cell%d-%d", row, col)
			cellTag := js.Global().Get("document").Call("getElementById", cellName)
			if cell.State == Opened {
				val := ""
				if cell.Val.IsMine() {
					val = "💣"
				}
				if cell.Val.IsNumber() {
					val = fmt.Sprint(cell.Val)
				}
				cellTag.Set("innerHTML", val)
				cellTag.Set("className", "openedCell")
			}
			if cell.State == Flagged {
				cellTag.Set("innerHTML", "🚩")
				cellTag.Set("className", "flaggedCell")
			}
			if cell.State == Closed {
				cellTag.Set("innerHTML", "")
				cellTag.Set("className", "closedCell")
			}
		}
	}
}

func registerCallbacks() {
	js.Global().Set("initGame", js.FuncOf(initGame))
	js.Global().Set("resetGame", js.FuncOf(resetGame))
	js.Global().Set("openCell", js.FuncOf(openCell))
	js.Global().Set("flagCell", js.FuncOf(flagCell))
}
