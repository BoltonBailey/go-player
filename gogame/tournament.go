package gogame

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
)

//Markers
const FILES_PRESERVED int = 1

const NUM_FILES int = 10

const PRINT bool = true

// Creates clean files in godata subdirectory
func MakeDatafiles() {

	for i := 0; i < NUM_FILES; i++ {
		// Get the zero padding right
		numberString := strconv.Itoa(i)
		for len(numberString) < len(strconv.Itoa(NUM_FILES)) {
			numberString = "0" + numberString
		}
		filename := "godata/datafile_" + numberString
		// Make the file, 0664 for read/write permission
		err := ioutil.WriteFile(filename, []byte{}, 0644)
		if err != nil {
			panic(err)
		}
	}

}

// Read the ith file.
// Takes as argument an int corresponding to one of the data files
// Returns a slice of bytes being the data within that file
func readDatafile(i int) []byte {
	if i < 0 || i >= NUM_FILES {
		panic("Illegal argument: file number out of range.\n")
	}

	numberString := strconv.Itoa(i)
	for len(numberString) < len(strconv.Itoa(NUM_FILES)) {
		numberString = "0" + numberString
	}

	filename := "godata/datafile_" + numberString

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return data
}

func printData(i int) {
	file := readDatafile(i)
	for _, datum := range file {
		fmt.Printf("%d\n", datum)
	}
}

// Write data to the ith file
// Takes as argument an int corresponding to one of the datafiles
// Takes as argument a slice of bytes
// Writes the slice of bytes into the file specified by the int
func writeDatafile(i int, data []byte) {
	// Get the zero padding right
	numberString := strconv.Itoa(i)
	for len(numberString) < len(strconv.Itoa(NUM_FILES)) {
		numberString = "0" + numberString
	}
	filename := "godata/datafile_" + numberString
	// Write to the file
	err := ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		panic(err)
	}
}

// Makes a player-type function using a the data from the ith file
// The function works by analyzing each legal move.
// The analysis is based on the board position, and on turn.
// The the analysis returns a int64 for each legal move.
// The Player will return the largest rating
func PlayerMaker(i int) func(Position) Intersection {
	return DataPlayerMaker(readDatafile(i))
}

// Helper function for Playermaker
func DataPlayerMaker(data []byte) func(Position) Intersection {

	return func(pos Position) Intersection {

		var maxAnalysis int64 = -1
		var bestIntn Intersection = PASS
		// Loop through all intersections
		for i := 0; uint8(i) < SIZE; i++ {
			for j := 0; uint8(j) < SIZE; j++ {
				var intn Intersection = Intersection{uint8(i), uint8(j)}
				if pos.worseThanPass(intn) {
					continue
				}
				// Create the board. If white to play, switch colors
				board := pos.board
				if !pos.blacksTurn {
					board.black, board.white = board.white, board.black
				}
				// Analyze this move
				currentAnalysis := analyzer(data, board, intn)
				// If analysis is better, update the move
				if currentAnalysis > maxAnalysis {
					maxAnalysis = currentAnalysis
					bestIntn = intn
				}
			}
		}
		return bestIntn
	}
}

/// TODO: Document
type Template struct {
	score       int64
	decidedMask [7]uint8
	blackMask   [7]uint8
	whiteMask   [7]uint8
}

func (tplt *Template) setDecided(i, j uint8) {
	tplt.decidedMask[i] |= (1 << j)
}

func (tplt *Template) setBlack(i, j uint8) {
	tplt.blackMask[i] |= (1 << j)
}

func (tplt *Template) setWhite(i, j uint8) {
	tplt.whiteMask[i] |= (1 << j)
}

func (tplt *Template) isDecided(i, j uint8) bool {
	return (tplt.decidedMask[i] & (1 << j)) != 0
}

func (tplt *Template) isBlack(i, j uint8) bool {
	return (tplt.blackMask[i] & (1 << j)) != 0
}

func (tplt *Template) isWhite(i, j uint8) bool {
	return (tplt.whiteMask[i] & (1 << j)) != 0
}

// Makes list of templates from the data
func makeTemplateList(data []byte) []Template {

	templateList := []Template{}
	currentTemplate := Template{}

	for _, readByte := range data {
		// Split up byte
		occupancy := int64(readByte / 64)
		blackOccupies := occupancy%2 == 1
		whiteOccupies := occupancy/2 == 1
		xDelta := uint8((readByte / 8) % 8)
		yDelta := uint8(readByte % 8)

		// Check if byte specifies score
		if xDelta == 7 || yDelta == 7 {
			// If it does, append the current template if nonzero score
			// and template is nonempty
			if currentTemplate.score != 0 && currentTemplate.decidedMask != [7]uint8{} {
				templateList = append(templateList, currentTemplate)
			}
			//Once the template is appended, set up new template
			currentTemplate = Template{}
			if xDelta == 7 && yDelta == 7 {
				currentTemplate.score = 0
			} else if xDelta == 7 {
				currentTemplate.score = occupancy*8 + int64(yDelta)
			} else {
				currentTemplate.score = (occupancy*8 + int64(xDelta))
			}
			continue
		}
		// Otherwise, add to template
		if currentTemplate.isDecided(xDelta, yDelta) {
			currentTemplate.score = 0
		} else {
			currentTemplate.setDecided(xDelta, yDelta)
			if blackOccupies {
				currentTemplate.setBlack(xDelta, yDelta)
			}
			if whiteOccupies {
				currentTemplate.setWhite(xDelta, yDelta)
			}
		}
	}
	return templateList
}

// Prints out a template
func (tplt *Template) printTemplate() {

	fmt.Printf("Score is %d\n", tplt.score)
	for x := uint8(0); x < 7; x++ {
		rowString := ""
		for y := uint8(0); y < 7; y++ {
			decided := tplt.isDecided(x, y)
			if !decided {
				rowString += "."
			} else {
				black := tplt.isBlack(x, y)
				white := tplt.isWhite(x, y)
				if !black && !white {
					// Empty intn
					rowString += "+"
				} else if black && !white {
					// Black stone
					rowString += "b"
				} else if !black && white {
					// White stone
					rowString += "w"
				} else {
					// Off board
					rowString += "*"
				}
			}
		}
		fmt.Println(rowString)
	}
}

// The analyzer takes a slice of bytes,
// a board, and who is to move (wrapped in a board with black to move)
// A starting intersection on the board to analyze
//
// Returns a int64
//
// The analyzer works as follows:
// Each byte is broken up into 4*8*8
// the byte ab 111 111 specifies a scoring of zero
// the byte ab 111 cde specifies a scoring of +abcde
// the byte ab cde 111 specifies a scoring of -abcde
// the byte ab cde fgh is a requirement:
// ab = 00 => empty
// ab = 01 => black
// ab = 10 => white
// ab = 11 => off-board
// cde != 111 != fgh
// cde represents x, fgh represents y coordinate of intersection
// grid is centered on starting intersection
func analyzer(data []byte, board Board, intn Intersection) int64 {
	templateList := makeTemplateList(data)
	total := int64(0)
	for _, tplt := range templateList {
		satisfied := true
		for x := uint8(0); x < 7; x++ {
			for y := uint8(0); y < 7; y++ {
				var xCoord uint8 = intn.x - 3 + x
				var yCoord uint8 = intn.y - 3 + y
				var scanIntn Intersection = Intersection{xCoord, yCoord}
				decided := tplt.isDecided(x, y)
				if decided {
					black := tplt.isBlack(x, y)
					white := tplt.isWhite(x, y)
					if xCoord >= SIZE || yCoord >= SIZE {
						satisfied = satisfied && black && white
					} else if board.isBlackStone(scanIntn) {
						satisfied = satisfied && black && !white
					} else if board.isWhiteStone(scanIntn) {
						satisfied = satisfied && !black && white
					} else {
						satisfied = satisfied && !black && !white
					}
				}
			}
		}
		if satisfied {
			total += tplt.score
		}
	}
	return total
}

// Prints how how the analysis works for a certain slice of bytes
func PrintAnalyzer(data []byte) {
	fmt.Println("Printing analyzer")
	templateList := makeTemplateList(data)
	for _, tplt := range templateList {
		tplt.printTemplate()
	}
}

// Runs a round rbin style tournament
// Each player is created from the respective data file
// Each player plays each other player, once as white, once as black
// The total scores make up the scoreboard
func RoundRobin() {
	// Initialize the array of players and the array of scores
	var players [NUM_FILES]func(Position) Intersection
	var scoreBoard [NUM_FILES]uint64
	// Fill the array of players with the players
	for i := 0; i < NUM_FILES; i++ {
		players[i] = PlayerMaker(i)
	}
	// Loop through each pair of players, playing a game and printing out
	for i := 0; i < NUM_FILES; i++ {
		for j := 0; j < NUM_FILES; j++ {
			fmt.Printf("Round Robin Challenge: %d, %d\n", i, j)
			var challengeGame Game = MakeGame(players[i], players[j])
			iScore, jScore := challengeGame.PlayGame()
			scoreBoard[i] += uint64(iScore)
			scoreBoard[j] += uint64(jScore)
		}
	}

	// Sort by score
	for i := 0; i < NUM_FILES; i++ {
		for j := 0; j < i; j++ {
			if scoreBoard[i] > scoreBoard[j] {
				scoreBoard[i], scoreBoard[j] = scoreBoard[j], scoreBoard[i]
				iData := readDatafile(i)
				jData := readDatafile(j)
				writeDatafile(i, jData)
				writeDatafile(j, iData)
			}
		}
	}
	//Print scores
	for i := 0; i < NUM_FILES; i++ {
		fmt.Printf("Scoreboard: %d has %d points \n", i, scoreBoard[i])
	}
	// Mutate the last file
	crucibleOfFire(NUM_FILES - 1)

}

func crucibleOfFire(i int) {
	fmt.Printf("Begin Crucible\n")
	for {

		newData := []byte{}
		for len(newData) < 512 {
			newData = append(newData, byte(rand.Intn(256)))
		}
		writeDatafile(i, newData)
		cruciblePlayer := PlayerMaker(i)
		gameToShow := MakeGame(cruciblePlayer, CapturePlayer)
		fmt.Printf("Play Crucible\n")
		i, j := gameToShow.PlayGame()
		if i > j+10 {
			fmt.Printf("End Crucible\n")
			return
		}
	}
}

// ****************************************************************************** Various other tournament styles

// Creates two players from files i and j
// plays them against each other
// The winner takes the i ranking (i should be better ranked than j)
func challenge(i, j int) {
	fmt.Printf("Challenge: %d, %d\n", i, j)
	if i >= j {
		panic("Better ranking challenging worse")
	}

	var challengeGame1 Game = MakeGame(PlayerMaker(i), PlayerMaker(j))
	iScore1, jScore1 := challengeGame1.PlayGame()

	var challengeGame2 Game = MakeGame(PlayerMaker(j), PlayerMaker(i))
	jScore2, iScore2 := challengeGame2.PlayGame()

	// See if j beat i
	if jScore1+jScore2 >= iScore1+iScore2 {
		fmt.Printf("%d beat %d: switched\n", j, i)
		iData := readDatafile(i)
		jData := readDatafile(j)
		writeDatafile(i, jData)
		writeDatafile(j, iData)
	}
}

// Another tourney style
func QuadEvolve() {
	// Choose four random files
	var players [4]int
	for i := 0; i < 4; i++ {
		players[i] = rand.Intn(NUM_FILES)
		fmt.Printf("Chose %d with length %d\n", players[i], len(readDatafile(i)))
	}
	var scoreBoard [4]uint64
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			var challengeGame Game = MakeGame(PlayerMaker(players[i]), PlayerMaker(players[j]))
			iScore, jScore := 0, 0
			iScore, jScore = challengeGame.PlayGame()
			if PRINT {
				challengeGame.PrintGame()
			}
			scoreBoard[i] += uint64(iScore)
			scoreBoard[j] += uint64(jScore)
		}
	}
	//Sort
	for i := 0; i < 4; i++ {
		for j := 0; j < i; j++ {
			if scoreBoard[i] > scoreBoard[j] {
				scoreBoard[i], scoreBoard[j] = scoreBoard[j], scoreBoard[i]
				players[i], players[j] = players[j], players[i]
			}
		}
	}
	for i := 0; i < 4; i++ {
		fmt.Printf("%d Scored %d\n", players[i], scoreBoard[i])
	}
	// winners child replace the losers
	// Get a random sequence from winner1 - ensure at least 1 byte
	winner1 := append(readDatafile(players[0]), byte(rand.Intn(256)))
	gene1end := rand.Intn(len(winner1)) + 1
	gene1start := rand.Intn(gene1end)
	gene1 := winner1[gene1start:gene1end]
	// Get a random sequence from winner2 - ensure at least 1 byte
	winner2 := append(readDatafile(players[1]), byte(rand.Intn(256)))
	gene2end := rand.Intn(len(winner2)) + 1
	gene2start := rand.Intn(gene2end)
	gene2 := winner2[gene2start:gene2end]
	// Make a child - length about a quarter the sum of lengths of winners
	child1 := append(gene1, gene2...)
	// Make another child with 30 random bytes
	child2 := []byte{}
	for len(child2) < 30 {
		child2 = append(child2, byte(rand.Intn(256)))
	}
	fmt.Printf("Made a child:\n")
	PrintAnalyzer(child1)
	writeDatafile(players[2], child1)
	writeDatafile(players[3], child2)
}

// Runs a tournament
func Gauntlet() {
	fmt.Println("Running tournament")
	for pres := 0; pres < FILES_PRESERVED; pres++ {
		for i := pres + 1; i < NUM_FILES; i++ {
			challenge(pres, i)
		}
	}
}

// Tests a gene by seeing
func GeneTester(gene []byte) float64 {
	fmt.Println("Testing gene")
	// Randomly seed the unpreserved files
	for i := FILES_PRESERVED; i < NUM_FILES; i++ {
		randomData := []byte{}
		for len(randomData) < len(gene) {
			randomData = append(randomData, byte(rand.Intn(256)))
		}
		writeDatafile(i, randomData)
	}
	// Now test the gene
	var geneScore uint64
	var otherScore uint64
	// Use a round robin
	for i := 0; i < NUM_FILES; i++ {
		for j := 0; j < NUM_FILES; j++ {
			black1 := DataPlayerMaker(append(readDatafile(i), gene...))
			white1 := PlayerMaker(j)
			var challengeGame1 Game = MakeGame(black1, white1)
			iScore, jScore := challengeGame1.PlayGame()
			challengeGame1.PrintGame()
			geneScore += uint64(iScore)
			otherScore += uint64(jScore)
			// Switch gene side
			black2 := PlayerMaker(i)
			white2 := DataPlayerMaker(append(readDatafile(j), gene...))
			var challengeGame2 Game = MakeGame(black2, white2)
			iScore, jScore = challengeGame2.PlayGame()
			otherScore += uint64(iScore)
			geneScore += uint64(jScore)
		}
	}
	// Get the new genes improvement coefficient
	improvement := float64(geneScore) / float64(geneScore+otherScore)
	fmt.Printf("Scored %f\n", improvement)
	return improvement
}

// Removes bytes from a gene until it starts corrupting the gene
func GeneImprover(gene []byte) []byte {
	// Try to improve n times
	currentScore := GeneTester(gene)
	for i := 0; i < 5; i++ {
		toMutate := rand.Intn(len(gene))
		mutatedGene := append(gene[:toMutate], gene[toMutate+1:]...)
		if GeneTester(mutatedGene) > currentScore {
			fmt.Println("Gene improved")
			return GeneImprover(mutatedGene)
		}
	}
	return gene
}

func BeatCapturePlayer() []byte {
	for {
		gene := []byte{}
		for len(gene) < 20 {
			gene = append(gene, byte(rand.Intn(256)))
		}

		black := CapturePlayer
		white := DataPlayerMaker(gene)
		var challengeGame1 Game = MakeGame(black, white)
		iScore1, jScore1 := challengeGame1.PlayGame()
		if iScore1 > jScore1 {
			continue
		}
		var challengeGame2 Game = MakeGame(black, white)
		iScore2, jScore2 := challengeGame2.PlayGame()
		var challengeGame3 Game = MakeGame(black, white)
		iScore3, jScore3 := challengeGame3.PlayGame()
		if iScore1+iScore2+iScore3 < jScore1+jScore2+jScore3 {
			PrintAnalyzer(gene)
			var challengeGame Game = MakeGame(black, white)
			challengeGame.PlayGame()
			challengeGame.PrintGame()
			return gene
		}
	}
}
