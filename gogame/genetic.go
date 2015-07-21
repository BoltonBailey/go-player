package gogame

import (
	//"fmt"
	//"io/ioutil"
	"math/rand"
	//"strconv"
)

// Uses a slice of bytes to create a Turing machine-like player
// Each machine has finite states, consisting of bytes
// The machine has two heads
// - A read head for the board, starting on the middle square
// - a read write board, that moves on the same board
//        The read/write cells can be incremented or decremented,
//        and are positive integers initialized to 0
// The machine starts on byte index 0
// Each state is either a command, consisting of one byte
//       executes inc/dec of counter
//       executes move up down left or right
//       that subsequently passes control to next byte
// Or a test, which consists of three bytes, and passes control to a
// different state.
func makeAutomatonPlayer(data []byte) func(Position) Intersection {

	return func(pos Position) Intersection {

		rwBoard := [SIZE][SIZE]int8{}
		rwFirst := (SIZE + 1) / 2
		rwSecond := (SIZE + 1) / 2

		var state uint8 = 0
		for i := 0; i < 100; i++ { // For loop limits number of steps
			if int(state) >= len(data) {
				return PASS
			}
			if rwFirst < 0 || rwFirst >= SIZE || rwSecond < 0 || rwSecond >= SIZE {
				rwFirst = (SIZE + 1) / 2
				rwSecond = (SIZE + 1) / 2
			}
			lit := [8]bool{} // To contain the bits of the byte of state
			//read in bits
			var bit uint8 = 0
			for ; bit < 8; bit++ {
				lit[bit] = (data[state] & (1 << bit)) != 0
			}
			if lit[0] { // 1 bit lit = command
				if lit[1] { // Change counter
					if lit[2] {
						rwBoard[rwFirst][rwSecond] += 2
					}
					if lit[3] {
						rwBoard[rwFirst][rwSecond] -= 1
					}
				}
				if lit[4] { // decide if head to move
					if lit[5] && lit[6] {
						rwFirst += 1
					} else if lit[5] && !lit[6] {
						rwFirst -= 1
					} else if !lit[5] && lit[6] {
						rwSecond += 1
					} else {
						rwSecond -= 1
					}
				}
				state += 1
			} else { // Go to command
				// Read the board
				onColor := (pos.board.isBlackStone(Intersection{rwFirst, rwSecond}) && pos.blacksTurn) ||
					(pos.board.isWhiteStone(Intersection{rwFirst, rwSecond}) && !pos.blacksTurn)
				empty := pos.board.isEmpty(Intersection{rwFirst, rwSecond})
				offColor := !empty && !onColor
				counter := rwBoard[rwFirst][rwSecond]
				if lit[6] {
					if empty && pos.isLegal(Intersection{rwFirst, rwSecond}) {
						return Intersection{rwFirst, rwSecond}
					} else {
						return RandomPlayer(pos)
					}
				} else if lit[2] {
					if onColor {
						state = data[state+1]
					} else {
						state = data[state+2]
					}
				} else if lit[3] {
					if offColor {
						state = data[state+1]
					} else {
						state = data[state+2]
					}
				} else if lit[4] {
					if empty {
						state = data[state+1]
					} else {
						state = data[state+2]
					}
				} else if lit[5] {
					if int(counter) > int(state%16) {
						state = data[state+1]
					} else {
						state = data[state+2]
					}
				}
			}
		}
		empty := pos.board.isEmpty(Intersection{rwFirst, rwSecond})
		if empty && pos.isLegal(Intersection{rwFirst, rwSecond}) {
			return Intersection{rwFirst, rwSecond}
		} else {
			return PASS
		}
	}
}

func RandomAutomatonPlayer() func(Position) Intersection {
	randBytes := []byte{}
	for len(randBytes) < 256 {
		randBytes = append(randBytes, byte(rand.Intn(256)))
	}
	return makeAutomatonPlayer(randBytes)
}
