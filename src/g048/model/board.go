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

// Grid represents tiles on the game board.
type Grid [BoardSize][BoardSize]Tile

// Coordinate is a convenience structure that stores a (row,col) pairing.
type Coordinate struct {
	Row uint32
	Col uint32
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
	score uint32
	// Random number generator
	random *rand.Rand
}

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
	for row := uint32(0); row < BoardSize; row++ {
		for col := uint32(0); col < BoardSize; col++ {
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
	for row := uint32(0); row < BoardSize; row++ {
		for col := uint32(0); col < BoardSize; col++ {
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
	// To end the game:
	//   1) The board must be filled.
	//   2) There are no 2 adjacent tiles with the same value.
	// TODO implement check #2
	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			if b.grid[row][col] == 0 {
				return false
			}
		}
	}
	return false
}

/*
 MoveLeft moves tiles to the left
*/
func (b *Board) MoveLeft() {
	// Repeat the accumulation process until all positions move as far as they
	// can.
	for i := 1; i < BoardSize; i++ {
		for row := 0; row < BoardSize; row++ {
			for col := 1; col < BoardSize; col++ {
				value := b.grid[row][col]
				prevCol := col - 1
				// If the previous value is 0, move the next value in
				if b.grid[row][prevCol] == 0 {
					b.grid[row][prevCol] = value
					b.grid[row][col] = 0
				} else if value == b.grid[row][prevCol] {
					// If the values are equal, accumulate
					b.grid[row][prevCol] *= 2
					b.grid[row][col] = 0
				}
			}
		}
	}
	// Every move generates a tile, if possible
	b.generateTile()
}

/*
 MoveRight moves tiles to the right
*/
func (b *Board) MoveRight() {
	// TODO implement
}

/*
 MoveUp moves tiles up
*/
func (b *Board) MoveUp() {
	// TODO implement
}

/*
 MoveDown moves tiles down
*/
func (b *Board) MoveDown() {
	// TODO implement
}
