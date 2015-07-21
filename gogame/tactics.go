package gogame

// Returns true if the given intersection
// - is illegal
// - can be shown that pass is an equivalent or better move
// Sometimes
// Returns true if intn fills in an eye with no enemy stones in it
func (pos *Position) worseThanPass(intn Intersection) bool {
	if !pos.isLegal(intn) {
		return true
	}
	if pos.isEye(intn) {
		return true
	}
	return false
}

// Is this an eye for the person on move?
func (pos *Position) isEye(intn Intersection) bool {

	for _, adjIntn := range intn.adjacents() {
		if pos.blacksTurn {
			if !pos.board.isBlackStone(adjIntn) {
				return false
			}
		} else {
			if !pos.board.isWhiteStone(adjIntn) {
				return false
			}
		}
	}
	tempBoard := pos.board
	tempBoard.removeChain(intn.adjacents()[0])

	for _, adjIntn := range intn.adjacents() {
		if !tempBoard.isEmpty(adjIntn) {
			return false
		}
	}
	return true
}

// Returns true if the move would capture enemy stones
func (pos *Position) isCapture(intn Intersection) bool {
	var tempBoard Board = pos.board
	if pos.blacksTurn {
		tempBoard.playBlackStone(intn)
		if tempBoard.white != pos.board.white && pos.isLegal(intn) {
			return true
		}
	} else {
		tempBoard.playWhiteStone(intn)
		if tempBoard.black != pos.board.black && pos.isLegal(intn) {
			return true
		}
	}
	return false
}

// Possible states for single chains: Consider a black chain
// Dead:   black can not save the group, even on blacks turn.
// Danger: black can only save the group if it is blacks turn.
// Alive:  black can save the group if they choose.
// Unconditionally Alive:   white cannot possibly kill the group.

/*
func (board *Board) isUnconditionallyAlive(intn Intersection) bool {

}
*/
