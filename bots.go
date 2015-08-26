package main

import "math"

type UXOBot interface {
	getMove(board Board, lastmove Move, player Player) (Move, error)
	makeMove(move Move) error
}

// Constants for the max and minimum board scores
const (
	SCOREMAX = 1000000
	SCOREMIN =-1000000
	SCORE_EXACT = 0
	SCORE_UPPER_BOUND = 1
	SCORE_LOWER_BOUND = 2

	//Draw on a sub-board is defined as 1024 
	SCOREDRAW = 1 << 10
)

func playerToMul(player Player) int {
	if player == X {
		return 1
	} else {
		return -1
	}
}

func evalBoard(board *Board) int {
	floatscores := new(Scores)

	for bx := 0; bx < 3; bx++ {
		for by := 0; by < 3; by++ {
			floatscores[bx][by] = normSubScore(evalSubBoard(board, bx*3, by*3))
		}
	}
	return evalSuperBoard(floatscores)
}

func evalSuperBoard(floatscores *Scores) int {
	//var xsum, osum float64
	var xmax, omax float64

	for _, l := range winlines {
		xs := floatscores[l.x1][l.y1] * floatscores[l.x2][l.y2] * floatscores[l.x3][l.y3]
		os := (1.0 - floatscores[l.x1][l.y1]) * (1.0 - floatscores[l.x2][l.y2]) * (1.0 - floatscores[l.x3][l.y3])
		if xs == 1.0 {
			return SCOREMAX
		} else if os == 1.0 {
			return SCOREMIN
		}


		if !math.IsNaN(xs) { //If the line isn't a draw
			xmax = math.Max(xmax, xs)
			omax = math.Max(omax, os)
		}
	}
	
	return int((xmax - omax) * 100000)
}

// func spreadOfSlice(slice *[]float64) float64 {
// 	maxEntry := 0.0
// 	minEntry := 0.0

// 	for _, entry := range *slice {
// 		maxEntry = math.Max(maxEntry, entry)
// 		minEntry = math.Min(minEntry, entry)
// 	}

// 	//minEntry must be <= 0 and maxEntry >= 0
// 	return minEntry + maxEntry
// }

// evaluate the score of a sub-board always with regard to X
// bx and by are the top left corner of the subboard to score	
func evalSubBoard(board *Board, bx int, by int) int {
	// fmt.Println(len(board[bx:bx+3]))
	bcols := board[bx:bx+3]
	b := make([][]Player, 3)
	for x := 0; x < 3; x++ {
		b[x] = bcols[x][by:by+3]
	}

	pieces := false

	xS := new(SubScores)
	oS := new(SubScores)
	var sS *SubScores
	var score int
	score = 0
	//fmt.Println("Entering evalSubBoard at location",bx,by)
	for _, l := range winlines {
		n := b[l.x1][l.y1] | b[l.x2][l.y2] | b[l.x3][l.y3]
		s := b[l.x1][l.y1] + b[l.x2][l.y2] + b[l.x3][l.y3]

		//fmt.Printf("OR: %b, AND: %b\n", n, s)

		// if there are mixed players or no players in a line
		if n > 2 {
			pieces = true
			continue
		} else if n == 0 {
			continue
		}

		// if the whole line is one player
		if s == 3*n {
			if s == 3 {
				return 1000
			} else {
				return -1000
			}
		}

		// which player can win this line
		if n == X {
			sS = xS
		} else {
			sS = oS
		}

		// if this line has a single piece
		if s == n {
			
			if b[l.x1][l.y1] == B {
				sS[l.x1][l.y1] |= 2
			}
			if b[l.x2][l.y2] == B {
				sS[l.x2][l.y2] |= 2
			}
			if b[l.x3][l.y3] == B {
				sS[l.x3][l.y3] |= 2
			}
		} else {
			if b[l.x1][l.y1] == B {
				sS[l.x1][l.y1] |= 1
			}
			if b[l.x2][l.y2] == B {
				sS[l.x2][l.y2] |= 1
			}
			if b[l.x3][l.y3] == B {
				sS[l.x3][l.y3] |= 1
			}
		}
	}

	// tally the score
	var Xones int
	var Xtwos int

	Xones = 0
	Xtwos = 0

	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			//		fmt.Print(xS[x][y])
			if xS[x][y]&1 == 1 {
				Xones++
			} else if xS[x][y]&2 == 2 {
				Xtwos++
			}
		}
	}
	//fmt.Printf("\nones: %d twos: %d\n", ones, twos)
	score += int(10*Xones + Xtwos)

	// tally the score
	var Oones int
	var Otwos int
	
	Oones = 0
	Otwos = 0

	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			//		fmt.Print(oS[x][y])
			if oS[x][y]&1 == 1 {
				Oones++
			} else if oS[x][y]&2 == 2 {
				Otwos++
			}
		}
	}
	//fmt.Printf("\nones: %d twos: %d\n", ones, twos)
	score -= int(10*Oones + Otwos)

	if Xones + Xtwos + Oones + Otwos == 0 && pieces {
		return SCOREDRAW
	}

	return score
}
