/*
 * File:        g048.go
 *
 * Author:      Schuyler Martin <schuylermartin45@gmail.com>
 *
 * Description: Main execution point of the `g048` project.
 */
package main

import (
	"./model"
	"./view"
	"fmt"
	"os"
	"strings"
)

/***** Constants *****/

// USAGE message to display on bad input
const USAGE string = "Usage: g048 [help]"

/***** Functions *****/

/*
 Main entry point of the G048 project.
*/
func main() {
	// Handle user input
	argc := len(os.Args)
	if argc > 1 {
		if strings.ToLower(os.Args[1]) == "help" {
			fmt.Println("G048: A Go-implementation of 2048")
			fmt.Println("\nAbout")
			fmt.Println("  Author: Schuyler Martin")
			fmt.Println("  Date:   February 2020")
			fmt.Println("\n" + USAGE + "\n")
			fmt.Println("Controls")
			fmt.Println("  * W/[Up]:         Move up")
			fmt.Println("  * A/[Left]:       Move left")
			fmt.Println("  * S/[Down]:       Move right")
			fmt.Println("  * D/[Right]:      Move down")
			fmt.Println("  * [Esc]/[Ctrl-C]: Exit game")

			os.Exit(view.EXIT_SUCCESS)
		} else {
			fmt.Fprintf(os.Stderr, "%v\n", USAGE)
			os.Exit(view.ERROR_USAGE)
		}
	}

	// Initialize, run, and exit with the selected mode
	textGame := new(view.TextGame)
	playAgain := true
	for playAgain {
		textGame.InitGame(model.NewBoard())
		playAgain = textGame.RenderGame()
	}
	textGame.ExitGame()
}
