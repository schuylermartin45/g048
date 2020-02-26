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
	// TODO improve
	t.screen.Fill(' ', tcell.StyleDefault.Background(tcell.ColorBlack))
	// Screen constants
	const (
		valuePad    = 4
		boardWidth  = valuePad * model.BoardSize
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
		// Center the tile's value on the tile
		rPad := fmt.Sprintf(fmt.Sprintf("%%%dv", valuePad), tile)
		valueStr := fmt.Sprintf(fmt.Sprintf("%%-%dv", valuePad), rPad)

		// Place value as string on the board
		x := xBoard + (pos.Col * len(valueStr))
		t.drawStr(x, y, valueStr, whiteText)
		if isEOL {
			y++
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
