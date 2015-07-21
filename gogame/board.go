package gogame

import (
	"fmt"
	"strconv"
)

const SIZE uint8 = 9

// Struct to hold the coordinates of an intersection on the board
type Intersection struct {
	x uint8
	y uint8
}

// Returns a slice of Intersections (cap 4)
// This contains every adjacent intersection to the given intersection.
func (intn *Intersection) adjacents() []Intersection {
	// intn x and y are uints, so thet are always positive
	// We must only check that they are under the SIZE value
	adj := make([]Intersection, 0, 4)
	if intn.x+1 < SIZE {
		adj = append(adj, Intersection{intn.x + 1, intn.y})
	}
	if intn.y+1 < SIZE {
		adj = append(adj, Intersection{intn.x, intn.y + 1})
	}
	if intn.x-1 < SIZE {
		adj = append(adj, Intersection{intn.x - 1, intn.y})
	}
	if intn.y-1 < SIZE {
		adj = append(adj, Intersection{intn.x, intn.y - 1})
	}
	return adj
}

// Struct to hold the current position on the board
// Bitmaps hold locations of stones of a given color
type Board struct {
	black [SIZE]uint32
	white [SIZE]uint32
}

// Returns whether a given intersection (param i) holds a black stone
func (board *Board) isBlackStone(i Intersection) bool {
	return (board.black[i.x] & (1 << i.y)) != 0
}

// Returns whether a given intersection holds a white stone
func (board *Board) isWhiteStone(i Intersection) bool {
	return (board.white[i.x] & (1 << i.y)) != 0
}

// Returns whether a given intersection is empty
func (board *Board) isEmpty(i Intersection) bool {
	return !board.isBlackStone(i) && !board.isWhiteStone(i)
}

// Returns true if the given Intersections have stones of the same color
// False if different colors or either intersection is empty
func (board *Board) sameColor(i1, i2 Intersection) bool {
	if board.isBlackStone(i1) && board.isBlackStone(i2) {
		return true
	}
	if board.isWhiteStone(i1) && board.isWhiteStone(i2) {
		return true
	}
	return false
}

/**
 * Prints out a textual display of the board.
 * uses + to represent empty intersections
 * uses w to represent white stones
 * uses b to represent black stones
 */
func (board *Board) PrintOut() {
	var rowString string = "  "
	for i := uint8(0); i < SIZE; i++ {
		rowString += strconv.Itoa(int(i)%10) + " "
	}
	fmt.Println(rowString)
	for i := uint8(0); i < SIZE; i++ {
		var rowString string = strconv.Itoa(int(i)%10) + " "
		for j := uint8(0); j < SIZE; j++ {
			if board.isBlackStone(Intersection{i, j}) {
				rowString += "b "
			} else if board.isWhiteStone(Intersection{i, j}) {
				rowString += "w "
			} else {
				rowString += "+ "
			}
		}
		fmt.Println(rowString)
	}

}

/**
 * Puts a black stone in the given intersection on the board.
 * Does not check if space is occupied or for captures.
 */
func (board *Board) placeBlackStone(i Intersection) {
	board.white[i.x] &^= 1 << i.y
	board.black[i.x] |= 1 << i.y
}

/**
 * Puts a white stone in the given intersection on the board.
 * Does not check if space is occupied or for captures.
 */
func (board *Board) placeWhiteStone(i Intersection) {
	board.white[i.x] |= 1 << i.y
	board.black[i.x] &^= 1 << i.y
}

// Switches the color of the stone at the intersection
// Throws an error if intersection is empty
func (board *Board) switchStoneColor(i Intersection) {
	if board.isBlackStone(i) {
		board.placeWhiteStone(i)
	} else if board.isWhiteStone(i) {
		board.placeBlackStone(i)
	} else {
		panic("Tried to switch color for empty intersection")
	}
}

/**
 * Clears the given intersection on the board.
 * Does not check if space is occupied
 */
func (board *Board) clearIntersection(i Intersection) {
	board.white[i.x] &^= 1 << i.y
	board.black[i.x] &^= 1 << i.y
}

// Removes all stones in a chain, and returns the number of stones removed
func (board *Board) removeChain(intn Intersection) {
	// If intersection is empty, we are done
	if board.isEmpty(intn) {
		return
	}
	// List of adjacent stones to remove
	toRemove := make([]Intersection, 0, 4)
	for _, adjIntn := range intn.adjacents() {
		if board.sameColor(intn, adjIntn) {
			toRemove = append(toRemove, adjIntn)
		}
	}
	// Remove current stone
	board.clearIntersection(intn)
	// Loop through toRemove
	for _, sameColorIntn := range toRemove {
		board.removeChain(sameColorIntn)
	}
}

// Fills empty space with black stones
func (board *Board) fillSpaceBlack(intn Intersection) {
	// If intersection is not empty, we are done
	if !board.isEmpty(intn) {
		return
	}
	board.placeBlackStone(intn)
	// Fill next to intersection
	for _, adjIntn := range intn.adjacents() {
		board.fillSpaceBlack(adjIntn)
	}
}

// Fills empty space with White stones
func (board *Board) fillSpaceWhite(intn Intersection) {
	// If intersection is not empty, we are done
	if !board.isEmpty(intn) {
		return
	}
	board.placeWhiteStone(intn)
	// Fill next to intersection
	for _, adjIntn := range intn.adjacents() {
		board.fillSpaceWhite(adjIntn)
	}
}

// Enacts the play of a black stone at the given empty intersection
func (board *Board) playBlackStone(intn Intersection) {
	if !board.isEmpty(intn) {
		panic("Tried to play black stone in nonempty intersection")
	}
	board.placeBlackStone(intn)
	for _, adjIntn := range intn.adjacents() {
		if board.isWhiteStone(adjIntn) && !board.hasLiberty(adjIntn) {
			board.removeChain(adjIntn)
		}
	}
	if !board.hasLiberty(intn) {
		board.removeChain(intn)
	}
}

// Enacts the play of a white stone at the given empty intersection
func (board *Board) playWhiteStone(intn Intersection) {
	if !board.isEmpty(intn) {
		panic("Tried to play black stone in nonempty intersection")
	}
	board.placeWhiteStone(intn)
	for _, adjIntn := range intn.adjacents() {
		if board.isBlackStone(adjIntn) && !board.hasLiberty(adjIntn) {
			board.removeChain(adjIntn)
		}
	}
	if !board.hasLiberty(intn) {
		board.removeChain(intn)
	}
}

// Asks if the stone at the given intersection has a liberty.
// Requires that the Intersection intn contains a stone
func (board *Board) hasLiberty(intn Intersection) bool {

	// We first ensure that our intersection actually contains a stone
	if board.isEmpty(intn) {
		panic("Tried to find liberties for empty intersection")
	}
	// Is stone black?
	blackStone := board.isBlackStone(intn)
	// Create tempBoard
	var tempBoard Board = *board
	// Remove the chain
	tempBoard.removeChain(intn)
	// Fill the space, and compare with original board
	if blackStone {
		tempBoard.fillSpaceBlack(intn)
		return tempBoard.black != board.black
	} else {
		tempBoard.fillSpaceWhite(intn)
		return tempBoard.white != board.white
	}
}

// Asks if the empty space at the given intersection is black territory
// Requires that the Intersection intn be empty
func (board *Board) isBlackTerritory(intn Intersection) bool {
	// We first ensure that our intersection is empty
	if !board.isEmpty(intn) {
		panic("Tried to test black territory on nonempty intersection")
	}
	// Create tempBoard
	var tempBoard Board = *board
	// Fill the space with white, remove the white chain, and compare
	tempBoard.fillSpaceWhite(intn)
	tempBoard.removeChain(intn)
	return tempBoard.white == board.white
}

// Asks if the empty space at the given intersection is white territory
// Requires that the Intersection intn be empty
func (board *Board) isWhiteTerritory(intn Intersection) bool {
	// We first ensure that our intersection is empty
	if !board.isEmpty(intn) {
		panic("Tried to test white territory on nonempty intersection")
	}
	// Create tempBoard
	var tempBoard Board = *board
	// Fill the space with black, remove the white chain, and compare
	tempBoard.fillSpaceBlack(intn)
	tempBoard.removeChain(intn)
	return tempBoard.black == board.black
}

type Position struct {
	board      Board
	blacksTurn bool
	illegal    [SIZE]uint32
}

func (pos *Position) setIllegal(i Intersection) {
	pos.illegal[i.x] |= 1 << i.y
}

func (pos *Position) isLegal(i Intersection) bool {
	return i == PASS || pos.illegal[i.x]&(1<<i.y) == 0
}
