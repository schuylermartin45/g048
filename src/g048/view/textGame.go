/*
 * File:        textGame.go
 *
 * Author:      Schuyler Martin <schuylermartin45@gmail.com>
 *
 * Description: An advanced gameplay mode that runs in a text terminal.
 */
package view

import (
	"../model"
	"fmt"
	"github.com/gdamore/tcell"
	"os"
	// Ticking away, the moments that make up the dull day...
	"time"
)

/***** Types *****/

// TextGame renders G048 in an interactive text-based UI.
type TextGame struct {
	board  *model.Board
	screen tcell.Screen
}

/***** Internal Functions *****/

/*
 Returns a color for each tile
*/
func getTileColor(tile model.Tile) tcell.Style {
	tileColor := tcell.StyleDefault.Foreground(tcell.ColorBlack)
	switch tile {
	case 0:
		tileColor = tileColor.Background(tcell.ColorLightGrey)
	case 2:
		tileColor = tileColor.Background(tcell.ColorDarkGrey)
	case 4:
		tileColor = tileColor.Background(tcell.ColorSlateGrey)
	case 8:
		tileColor = tileColor.Background(tcell.ColorGrey)
	case 16:
		tileColor = tileColor.Background(tcell.ColorTeal)
	case 32:
		tileColor = tileColor.Background(tcell.ColorTan)
	case 64:
		tileColor = tileColor.Background(tcell.ColorOrange)
	case 128:
		tileColor = tileColor.Background(tcell.ColorDarkOrange)
	case 256:
		tileColor = tileColor.Background(tcell.ColorOrangeRed)
	case 512:
		tileColor = tileColor.Background(tcell.ColorRed)
	case 1024:
		tileColor = tileColor.Background(tcell.ColorDarkRed)
	case 2048:
		tileColor = tileColor.Background(tcell.ColorDarkViolet)
	case 4096:
		tileColor = tileColor.Background(tcell.ColorBlueViolet)
	// We don't have a lot of colors to work with so for now we'll make it boring
	// pass 4096.
	default:
		tileColor = tileColor.Background(tcell.ColorSlateGrey)
	}
	return tileColor
}

/***** Internal Members *****/

/*
 Draws a string.

 @param x         Left-top corner x position of the string
 @param y         Left-top corner y position of the string
 @param str       String to draw
 @param textColor Color to draw the text in
*/
func (t *TextGame) drawStr(x int, y int, str string, textColor tcell.Style) {
	sizeX, sizeY := t.screen.Size()
	if (x < 0) || (y < 0) || (y > sizeY) {
		return
	}
	for row := 0; row < len(str); row++ {
		screenX := x + row
		if screenX > sizeX {
			break
		}
		t.screen.SetContent(screenX, y, rune(str[row]), nil, textColor)
	}
}

/*
 Draws the game entire gameboard/screen
*/
func (t *TextGame) drawBoard() {
	t.screen.Fill(' ', tcell.StyleDefault.Background(tcell.ColorBlack))
	// Screen constants
	const (
		blankStr    = "        " // No border
		blockWidth  = len(blankStr)
		boardWidth  = blockWidth * model.BoardSize
		boardHeight = model.BoardSize
	)
	xScreen, yScreen := t.screen.Size()
	// Screen variables
	var (
		xBoard    = (xScreen / 2) - (boardWidth / 2)
		yBoard    = (yScreen / 2) - (boardHeight / 2)
		scoreStr  = t.board.GetDisplayScore()
		xScore    = (xScreen / 2) - (len(scoreStr) / 2)
		yScore    = yBoard - 2
		whiteText = tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
	)

	// Draw the score above the board
	t.drawStr(xScore, yScore, scoreStr, whiteText)

	// Draw the board
	y := yBoard
	t.board.RenderBoard(func(pos model.Coordinate, isEOL bool, tile model.Tile) {
		// Default to the blank row (don't render the 0 for blank tiles)
		valueStr := blankStr
		if tile != 0 {
			// TODO: fix for 3-length values
			// Center the tile's value on the tile
			tileStr := fmt.Sprintf("%v", tile)
			halfTile := (len(tileStr) / 2)
			// Single-digit numbers will int-divide to 0, so this keeps the pad
			// length consistent and prevents gaps in the other calculations.
			if halfTile == 0 {
				halfTile++
			}
			tilePadLen := (blockWidth / 2) - halfTile
			tilePad := ""
			for i := 0; i < tilePadLen; i++ {
				tilePad += " "
			}
			valueStr = tilePad + tileStr + tilePad
			// For odd numbers, add additional missing padding
			for len(valueStr) < blockWidth {
				valueStr += " "
			}
		}

		// Place value as string on the board
		x := xBoard + (pos.Col * len(valueStr))
		tileColor := getTileColor(tile)
		t.drawStr(x, y+0, blankStr, tileColor)
		t.drawStr(x, y+1, valueStr, tileColor)
		t.drawStr(x, y+2, blankStr, tileColor)
		// Increment to draw the next row
		if isEOL {
			y += 3
		}
	})
	t.screen.Sync()
}

/*
 Initializes the event listener
*/
func (t *TextGame) initEventListener() {
	for {
		event := t.screen.PollEvent()
		switch eventType := event.(type) {
		case *tcell.EventKey:
			var action Action = ActionIllegal
			switch eventType.Key() {
			// ASCII keys have to be handled separately
			case tcell.KeyRune:
				switch eventType.Rune() {
				case 'w':
					action = ActionUp
				case 'a':
					action = ActionLeft
				case 'd':
					action = ActionRight
				case 's':
					action = ActionDown
				}
			case tcell.KeyUp:
				action = ActionUp
			case tcell.KeyLeft:
				action = ActionLeft
			case tcell.KeyRight:
				action = ActionRight
			case tcell.KeyDown:
				action = ActionDown
			// Exit
			case tcell.KeyCtrlC:
				fallthrough
			case tcell.KeyEsc:
				action = ActionExit
			}
			if action != ActionIllegal {
				ActionHandler(t.board, action, func() {
					t.screen.Fini()
					os.Exit(EXIT_SUCCESS)
					return
				})
				// Re-render the board on a delta (the user has made a move).
				t.drawBoard()
			}
		default:
			continue
		}
	}
}

/***** Methods *****/

// InitGame initializes the game.
func (t *TextGame) InitGame(b *model.Board) {
	t.board = b

	// Init the screen on first game. Subsequent games do not re-initialized.
	if t.screen == nil {
		tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
		var error error
		t.screen, error = tcell.NewScreen()
		if error != nil {
			fmt.Fprintf(os.Stderr, "%v\n", error)
			os.Exit(ERROR_SCREEN_INIT)
		}
		if error = t.screen.Init(); error != nil {
			fmt.Fprintf(os.Stderr, "%v\n", error)
			os.Exit(ERROR_SCREEN_INIT)
		}
		// Kick off event listener thread.
		go t.initEventListener()
	}
}

// RenderGame runs the primary gameplay loop.
func (t *TextGame) RenderGame() bool {
	// Draw the initial board. Subsequent renders will come on a user's action.
	t.drawBoard()
	// This game only redraws when the user does something. So the main
	// thread just has to spin-lock until the game is over. In an effort not
	// to peg the CPU, we will make the thread sleep.
	for !t.board.IsEndGame() {
		time.Sleep(200 * time.Millisecond)
	}
	return true
}

// ExitGame is a callback triggered when the game terminates
func (t *TextGame) ExitGame() {
	// Clean up screen object
	t.screen.Fini()
}
