/*
 * File:        board.go
 *
 * Author:      Schuyler Martin <schuylermartin45@gmail.com>
 *
 * Description: Represents the 2048 board.
 */
package model

import (
	"fmt"
	"math/rand"
	"time"
)

/***** Constants *****/

const (
	// BoardSize is the default size of the board
	BoardSize = 4
)

/***** Types *****/

// Tile represents a single tile value on the board.
type Tile uint32

// Score represents the user's score
type Score uint32

// Grid represents tiles on the game board.
type Grid [BoardSize][BoardSize]Tile

// Coordinate is a convenience structure that stores a (row,col) pairing.
type Coordinate struct {
	Row int
	Col int
}

/*
 DrawTile is a callback that renders a single tile when called by
 `RenderBoard()`.

 @param pos   Row, column position in the board.
 @param isEOL Flag indicates if this is the last column drawn in a row.
 @param tile  Value of the tile at this position
*/
type DrawTile func(pos Coordinate, isEOL bool, tile Tile)

// Board is the primary structure that represents the game's state.
type Board struct {
	// Current tile layout
	grid Grid
	// Game's current score.
	score Score
	// Random number generator
	random *rand.Rand
}

// Helper lambda function that performs one iteration of the move process.
// This is dependent on the direction.
type moveBoard func()

/***** Functions *****/

/*
 NewBoard constructs a new board play with.

 @return New board object to start a game with
*/
func NewBoard() *Board {
	b := new(Board)
	b.score = 0
	// Set a new random generator per game. This ensures that we don't
	// constantly reconstruct the generator for every random value we need.
	b.random = rand.New(rand.NewSource(time.Now().UnixNano()))
	// Initialize the starting board configuration.
	b.initBoard()
	return b
}

/***** Internal Members *****/

/*
 Generates a new random tile and places it on the board
*/
func (b *Board) generateTile() {
	// Tiles that are added to the board start at either 2 or 4, with 2 having
	// a much higher probability to show up.
	tileValue := Tile(2)
	if b.random.Intn(4) == 0 {
		tileValue = 4
	}

	// Since the board is relatively small, iterate over the entire board and
	// record all possible positions.
	var possiblePositions []*Coordinate
	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			if b.grid[row][col] == 0 {
				possiblePositions = append(possiblePositions, &Coordinate{row, col})
			}
		}
	}

	// Only add a new tile if the board has an open position
	possibleSize := len(possiblePositions)
	if possibleSize > 0 {
		pos := possiblePositions[b.random.Intn(possibleSize)]
		b.grid[pos.Row][pos.Col] = tileValue
	}
}

/*
 Initializes the new board
*/
func (b *Board) initBoard() {
	// Two tiles are randomly placed on the empty board
	b.generateTile()
	b.generateTile()
}

/*
 Helper function that centralizes the move logic for all 4 moves, handling
 the accumulation of values and the current score.

 @param curPos  Current board position being examined
 @param nextPos Next board position in the direction of the move. This is the
                position that is being accumulated into.
*/
func (b *Board) calcMove(curPos Coordinate, nextPos Coordinate) {
	curValue := b.grid[curPos.Row][curPos.Col]
	nextValue := b.grid[nextPos.Row][nextPos.Col]
	// If the other value is 0, move the current value in
	if b.grid[nextPos.Row][nextPos.Col] == 0 {
		b.grid[nextPos.Row][nextPos.Col] = curValue
		b.grid[curPos.Row][curPos.Col] = 0
	} else if curValue == nextValue {
		// If the values are equal, accumulate
		b.grid[nextPos.Row][nextPos.Col] *= 2
		b.grid[curPos.Row][curPos.Col] = 0
		// Score increments with the value accumulated
		b.score += Score(nextValue)
	}
}

/*
 Helper function that de-dupes core move logic from directional iterations.

 To quote my alma mater, "Make Moves, Son!"

 @param move Helper lambda that iterates over the board in the desired
             direction.
*/
func (b *Board) makeMove(move moveBoard) {
	// Repeat the accumulation process until all positions move as far as they
	// can.
	for i := 1; i < BoardSize; i++ {
		move()
	}
	// Every move generates a tile, if possible
	b.generateTile()
}

/***** Members *****/

/*
 GetDisplayScore returns the score as a displayable string

 @return Score, as a displayable string.
*/
func (b *Board) GetDisplayScore() string {
	return fmt.Sprintf("Score: %10d", b.score)
}

/*
 RenderBoard iterates over the board, invoking a callback to render a tile
 at a given position.
*/
func (b *Board) RenderBoard(draw DrawTile) {
	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			isEOL := (col + 1) == BoardSize
			draw(Coordinate{row, col}, isEOL, b.grid[row][col])
		}
	}
}

/*
 IsEndGame determines if the game has ended.

 @return True if the game ended. False otherwise.
*/
func (b *Board) IsEndGame() bool {
	boundSize := BoardSize - 1
	// To end the game:
	//   1) The board must be filled.
	//   2) There are no 2 adjacent tiles with the same value.
	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			value := b.grid[row][col]
			// Board is not filled
			if value == 0 {
				return false
			}
			// Check surrounding positions for equivalent values
			if (row > 0) && (value == b.grid[row-1][col]) {
				return false
			}
			if (row < boundSize) && (value == b.grid[row+1][col]) {
				return false
			}
			if (col > 0) && (value == b.grid[row][col-1]) {
				return false
			}
			if (col < boundSize) && (value == b.grid[row][col+1]) {
				return false
			}
		}
	}
	return true
}

/*
 MoveLeft moves tiles to the left
*/
func (b *Board) MoveLeft() {
	b.makeMove(func() {
		for row := 0; row < BoardSize; row++ {
			for col := 1; col < BoardSize; col++ {
				b.calcMove(Coordinate{row, col}, Coordinate{row, col - 1})
			}
		}
	})
}

/*
 MoveRight moves tiles to the right
*/
func (b *Board) MoveRight() {
	b.makeMove(func() {
		for row := 0; row < BoardSize; row++ {
			for col := BoardSize - 2; col >= 0; col-- {
				b.calcMove(Coordinate{row, col}, Coordinate{row, col + 1})
			}
		}
	})
}

/*
 MoveUp moves tiles up
*/
func (b *Board) MoveUp() {
	b.makeMove(func() {
		for row := 1; row < BoardSize; row++ {
			for col := 0; col < BoardSize; col++ {
				b.calcMove(Coordinate{row, col}, Coordinate{row - 1, col})
			}
		}
	})
}

/*
 MoveDown moves tiles down
*/
func (b *Board) MoveDown() {
	b.makeMove(func() {
		for row := BoardSize - 2; row >= 0; row-- {
			for col := 0; col < BoardSize; col++ {
				b.calcMove(Coordinate{row, col}, Coordinate{row + 1, col})
			}
		}
	})
}
