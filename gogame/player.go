package gogame

import (
	"fmt"
	"math/rand"
)

// A function that gets user input to return an intersection.
func HumanPlayer(pos Position) Intersection {

	pos.board.PrintOut()
	// Get first coordinate
	fmt.Println("Please enter coordinates, separated by space ")
	var i uint8
	var j uint8

	_, err := fmt.Scanf("%d %d", &i, &j)
	if err != nil {
		panic("Problem reading input")
	}

	if i == SIZE && j == SIZE {
		return PASS
	}
	// Throw errors when i, j is out of range
	if i >= SIZE || j >= SIZE {
		panic("Entered coordinates out of range")
	}
	// Create and return the intersection
	return Intersection{i, j}
}

// A function that ramdomly selects an intersection
func RandomPlayer(pos Position) Intersection {

	var chosenIntn Intersection
	for {
		if 0 == rand.Intn(50) {
			return PASS
		}
		i := rand.Intn(int(SIZE))
		j := rand.Intn(int(SIZE))
		chosenIntn = Intersection{uint8(i), uint8(j)}
		if pos.isLegal(chosenIntn) {
			return chosenIntn
		}
	}
}

// A function that scans right-left top-bottom and selects first legal move
func BadPlayer(pos Position) Intersection {

	// Will pass 1 out of 20 times
	if 0 == rand.Intn(20) {
		return PASS
	}
	// Loop through all intersections, find first empty one
	for i := 0; uint8(i) < SIZE; i++ {
		for j := 0; uint8(j) < SIZE; j++ {
			var intn Intersection = Intersection{uint8(i), uint8(j)}
			if pos.isLegal(intn) {
				return intn
			}
		}
	}
	return PASS
}

// A function that looks for an enemy stone, then tries to surround it
// If board is completely empty, will play in the middle
func SurroundPlayer(pos Position) Intersection {

	completelyEmpty := true
	// Loop through all intersections
	for i := 0; uint8(i) < SIZE; i++ {
		for j := 0; uint8(j) < SIZE; j++ {
			var intn Intersection = Intersection{uint8(i), uint8(j)}
			if !pos.board.isEmpty(intn) {
				completelyEmpty = false
			}
			// If empty, see if it can be filled
			if pos.blacksTurn {
				if pos.board.isWhiteStone(intn) {
					for _, adjIntn := range intn.adjacents() {
						if pos.isLegal(adjIntn) {
							return adjIntn
						}
					}
				}
			} else {
				if pos.board.isBlackStone(intn) {
					for _, adjIntn := range intn.adjacents() {
						if pos.isLegal(adjIntn) {
							return adjIntn
						}
					}
				}
			}
		}
	}
	if completelyEmpty {
		return RandomPlayer(pos)
	}
	return PASS
}

// A function that looks for a capture
// otherwise plays as SurroundPlayer
func CapturePlayer(pos Position) Intersection {
	board := pos.board
	for i := 0; uint8(i) < SIZE; i++ {
		for j := 0; uint8(j) < SIZE; j++ {
			var intn Intersection = Intersection{uint8(i), uint8(j)}
			// See if a play here by self will capture
			if board.isEmpty(intn) && pos.isCapture(intn) {
				return intn
			}
		}
	}
	return SurroundPlayer(pos)
}
