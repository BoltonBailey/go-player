package gogame

import (
	"fmt"
	//"time"
)

// We adopt the convention that Intersection{SIZE, SIZE} represents a pass
var PASS Intersection = Intersection{SIZE, SIZE}

type Game struct {
	// A slice of all boards so far in the game, starting with empty board
	BoardList []Board
	// Functions that are the players
	BlackPlayer func(Position) Intersection
	WhitePlayer func(Position) Intersection
}

// Makes the current position of the game.
func (game *Game) makeCurrentPosition() Position {
	// First we create a position to represent the current game position
	var currentPostion Position
	move := len(game.BoardList)
	// The board is the last element of the BoardList slice
	currentPostion.board = game.BoardList[move-1]
	// The player to move is black iff the boardlist has odd length
	currentPostion.blacksTurn = (move%2 == 1)
	// We must find all illegal intersections - occupied, suicide and ko
	// Loop through all intersections
	for i := 0; uint8(i) < SIZE; i++ {
		for j := 0; uint8(j) < SIZE; j++ {
			var intn Intersection = Intersection{uint8(i), uint8(j)}
			// If empty, create a copy of the board and make a move there
			if currentPostion.board.isEmpty(intn) {
				tempBoard := currentPostion.board
				if currentPostion.blacksTurn {
					tempBoard.playBlackStone(intn)
				} else {
					tempBoard.playWhiteStone(intn)
				}
				// If the intersection is now empty, suicide
				if tempBoard.isEmpty(intn) {
					currentPostion.setIllegal(intn)
					continue
				}
				// If the board state matches another, ko
				for moveIt := move - 1; moveIt >= 0; moveIt-- {
					if tempBoard == game.BoardList[moveIt] {
						currentPostion.setIllegal(intn)
						break
					}
				}
			} else {
				// Intersection nonempty, so illegal
				currentPostion.setIllegal(intn)
			}
		}
	}
	return currentPostion
}

func (game *Game) gameOver() bool {
	move := len(game.BoardList)
	if move <= 2 {
		return false
	}
	previousPass := game.BoardList[move-2] == game.BoardList[move-3]
	lastPass := game.BoardList[move-1] == game.BoardList[move-2]
	return lastPass && previousPass
}

// Prints out the record of the game, which has been played.
func (game *Game) PrintGame() {

	for i := 0; i < len(game.BoardList); i++ {
		fmt.Printf("Move number %d:\n", i)
		game.BoardList[i].PrintOut()
		fmt.Println()
	}
	move := len(game.BoardList)
	// The board is the last element of the BoardList slice
	fmt.Println("Chinese Scoring:")
	blackScore, whiteScore := game.BoardList[move-1].chineseScoring()
	fmt.Printf("Black's score is: %d\n", blackScore)
	fmt.Printf("White's score is: %d\n", whiteScore)
	fmt.Println()
}

// Plays a single turn, returns true if game ended
func (game *Game) playTurn() {
	currentPosition := game.makeCurrentPosition()
	// We now get the move from the player
	if currentPosition.blacksTurn {
		// Get blacks choice of move
		blacksMove := game.BlackPlayer(currentPosition)
		// Check legality of move
		if !currentPosition.isLegal(blacksMove) {
			panic("Illegal Move\n")
		}
		if blacksMove != PASS {
			currentPosition.board.playBlackStone(blacksMove)
		}
	} else {
		// Get whites choice of move
		whitesMove := game.WhitePlayer(currentPosition)
		// Check legality of move
		if !currentPosition.isLegal(whitesMove) {
			panic("Illegal Move\n")
		}
		if whitesMove != PASS {
			currentPosition.board.playWhiteStone(whitesMove)
		}
	}
	game.BoardList = append(game.BoardList, currentPosition.board)
}

// Activates the game,
// Keeps playing until two passes in a row
// Does not have ko or suicide rules, or komi
func (game *Game) PlayGame() (int, int) {
	for !game.gameOver() {
		game.playTurn()
		if len(game.BoardList) > 1000 {
			fmt.Println("Game ends on 1000 move rule")
			break
		}
	}
	// The game is over
	move := len(game.BoardList)
	// The board is the last element of the BoardList slice
	blackScore, whiteScore := game.BoardList[move-1].chineseScoring()
	return blackScore, whiteScore
}

func MakeGame(blackPlayer, whitePlayer func(Position) Intersection) Game {
	var game Game
	game.BoardList = make([]Board, 1, 1)
	game.BlackPlayer = blackPlayer
	game.WhitePlayer = whitePlayer
	return game
}

// Scores the came using chinese rules and prints out results
// Prints out scores, without komi, for black and white
// Prints out a board showing black and white area
// Returns black and white scores
func (board *Board) chineseScoring() (int, int) {

	// compute the area of each player
	var scoreBoard Board = *board
	blackScore := 0
	whiteScore := 0
	// Loop through all intersections
	for i := 0; uint8(i) < SIZE; i++ {
		for j := 0; uint8(j) < SIZE; j++ {
			var intn Intersection = Intersection{uint8(i), uint8(j)}
			// If empty, see if it can be filled
			if scoreBoard.isEmpty(intn) {
				if scoreBoard.isBlackTerritory(intn) {
					scoreBoard.fillSpaceBlack(intn)
				} else if scoreBoard.isWhiteTerritory(intn) {
					scoreBoard.fillSpaceWhite(intn)
				}
			}
			// If now nonempty, add to score
			if scoreBoard.isBlackStone(intn) {
				blackScore++
			}
			if scoreBoard.isWhiteStone(intn) {
				whiteScore++
			}
		}
	}
	return blackScore, whiteScore
}
