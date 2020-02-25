/*
 * File:        display.go
 *
 * Author:      Schuyler Martin <schuylermartin45@gmail.com>
 *
 * Description: Display interface that defines how a display mode operates.
 */
package view

import (
	"../model"
)

/***** Constants *****/

// Error code constants
const (
	EXIT_SUCCESS      = 0
	ERROR_USAGE       = 1
	ERROR_SCREEN_INIT = 2
)

/***** Types *****/

// Action describes a user-caused event in the game.
type Action uint8

// Enumeration of actions
const (
	ActionIllegal Action = 0
	ActionLeft    Action = 1
	ActionRight   Action = 2
	ActionUp      Action = 3
	ActionDown    Action = 4
	ActionExit    Action = 5
)

// ExitFunc is a callback triggered on `ActionExit`. This breaks the game loop
type ExitFunc func()

// Display is an interface that describes the features of a way to render the
// game.
type Display interface {
	// Initializes the game.
	InitGame(b *model.Board)
	// Runs the primary gameplay loop, returning true to play again.
	RenderGame() bool
	// Callback for when the game terminates.
	ExitGame()
}

/***** Functions *****/

/*
 Action handler. Given an action, performs a board operation.

 @param board  Pointer to the board to modify.
 @param action Action to interpret
 @param onExit	Function to call on exit
*/
func ActionHandler(board *model.Board, action Action, onExit ExitFunc) {
	switch action {
	case ActionIllegal:
		return
	case ActionLeft:
		board.MoveLeft()
	case ActionRight:
		board.MoveRight()
	case ActionUp:
		board.MoveUp()
	case ActionDown:
		board.MoveDown()
	case ActionExit:
		onExit()
	}
}
