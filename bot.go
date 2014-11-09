package main

//import ("fmt")

// Constants for the max and minimum board scores
const (
	SCOREMAX = 1000000
	SCOREMIN =-1000000
)

var wins = [][]int{
	{1, 1, 1, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 1, 1, 1, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 1, 1, 1},
	{1, 0, 0, 0, 1, 0, 0, 0, 1},
	{0, 0, 1, 0, 1, 0, 1, 0, 0},
	{1, 0, 0, 1, 0, 0, 1, 0, 0},
	{0, 1, 0, 0, 1, 0, 0, 1, 0},
	{0, 0, 1, 0, 0, 1, 0, 0, 1},
}

type Line struct {
	x1 int
	x2 int
	x3 int
	y1 int
	y2 int
	y3 int
}

var winlines = []Line{
	Line{x1: 0, x2: 0, x3: 0, y1: 0, y2: 1, y3: 2},
	Line{x1: 1, x2: 1, x3: 1, y1: 0, y2: 1, y3: 2},
	Line{x1: 2, x2: 2, x3: 2, y1: 0, y2: 1, y3: 2},

	Line{x1: 0, x2: 1, x3: 2, y1: 0, y2: 0, y3: 0},
	Line{x1: 0, x2: 1, x3: 2, y1: 1, y2: 1, y3: 1},
	Line{x1: 0, x2: 1, x3: 2, y1: 2, y2: 2, y3: 2},

	Line{x1: 0, x2: 1, x3: 2, y1: 0, y2: 1, y3: 2},
	Line{x1: 0, x2: 1, x3: 2, y1: 2, y2: 1, y3: 0},
}

//AI types and structs
type BHash uint64
type CacheEntry struct {
	cutoff int
	depth  int
	score  int
	move   Move
}
type AINode struct {
	board  *Board
	moves  *MoveSlice
	cache  map[BHash]CacheEntry
	scores *Scores
	hash   BHash
}

/*
*
* AI functionality
*
 */

func playerToMul(player int) int {
	if player == X {
		return 1
	} else {
		return -1
	}
}

// evaluate the score of the board specialized for the AI
func evalAIBoard(board *Board, scores *Scores) int {
	score := 0.0
	intscores := [3][3]float64{}

	for bx := 0; bx < 3; bx++ {
		for by := 0; by < 3; by++ {
			intscores[bx][by] = float64(evalSubBoard(board, bx*3, by*3))
		}
	}

	for _, l := range winlines {
		s := intscores[l.x1][l.y1] + intscores[l.x2][l.y2] + intscores[l.x3][l.y3]
		if s == 3000 {
			return SCOREMAX
		} else if s == -3000 {
			return SCOREMIN
		} else {
			//			score = (((intscores[l.x1][l.y1] + 1000.0) / 2000.0) + ((intscores[l.x2][l.y2] + 1000.0) / 2000.0) + ((intscores[l.x3][l.y3] + 1000.0) / 2000.0))/3
			score = (((intscores[l.x1][l.y1]) / 1000.0) + ((intscores[l.x2][l.y2]) / 1000.0) + ((intscores[l.x3][l.y3]) / 1000.0))/3
// 			fmt.Println(square)
// 			score += (square * square) - 0.5
		}
	}
// 	fmt.Println(score)
	return int(score * 500000)
}

func evalBoard(board *Board) int {
	score := 0.0
	intscores := [3][3]float64{}

	for bx := 0; bx < 3; bx++ {
		for by := 0; by < 3; by++ {
			intscores[bx][by] = float64(evalSubBoard(board, bx*3, by*3))
		}
	}

	for _, l := range winlines {
		s := intscores[l.x1][l.y1] + intscores[l.x2][l.y2] + intscores[l.x3][l.y3]
		if s == 3000 {
			return SCOREMAX
		} else if s == -3000 {
			return SCOREMIN
		} else {
//			score = (((intscores[l.x1][l.y1] + 1000.0) / 2000.0) + ((intscores[l.x2][l.y2] + 1000.0) / 2000.0) + ((intscores[l.x3][l.y3] + 1000.0) / 2000.0))/3
			score = (((intscores[l.x1][l.y1]) / 1000.0) + ((intscores[l.x2][l.y2]) / 1000.0) + ((intscores[l.x3][l.y3]) / 1000.0))/3
// 			fmt.Println(square)
// 			score += (square * square) - 0.5
		}
	}
// 	fmt.Println(score)
	return int(score * 500000)
}

// evaluate the score of a sub-board always with regard to X
func evalSubBoard(b *Board, bx int, by int) int {
	xS := new(Scores)
	oS := new(Scores)
	var sS *Scores
	var score int
	score = 0
	//fmt.Println("Entering evalSubBoard at location",bx,by)
	for _, l := range winlines {
		n := b[bx+l.x1][by+l.y1] | b[bx+l.x2][by+l.y2] | b[bx+l.x3][by+l.y3]
		s := b[bx+l.x1][by+l.y1] + b[bx+l.x2][by+l.y2] + b[bx+l.x3][by+l.y3]

		//fmt.Printf("OR: %b, AND: %b\n", n, s)

		// if there are mixed players or no players in a line
		if n > 2 || n == 0 {
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
		if n == 1 {
			sS = xS
		} else {
			sS = oS
		}

		if s == n {
			if b[bx+l.x1][by+l.y1] == 0 {
				sS[l.x1][l.y1] |= 2
			}
			if b[bx+l.x2][by+l.y2] == 0 {
				sS[l.x2][l.y2] |= 2
			}
			if b[bx+l.x3][by+l.y3] == 0 {
				sS[l.x3][l.y3] |= 2
			}
		} else {
			if b[bx+l.x1][by+l.y1] == 0 {
				sS[l.x1][l.y1] |= 1
			}
			if b[bx+l.x2][by+l.y2] == 0 {
				sS[l.x2][l.y2] |= 1
			}
			if b[bx+l.x3][by+l.y3] == 0 {
				sS[l.x3][l.y3] |= 1
			}
		}
	}

	// tally the score
	var ones int
	var twos int

	ones = 0
	twos = 0

	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			//		fmt.Print(xS[x][y])
			if xS[x][y]&1 == 1 {
				ones++
			} else if xS[x][y]&2 == 2 {
				twos++
			}
		}
	}
	//fmt.Printf("\nones: %d twos: %d\n", ones, twos)
	score += int(10*ones + twos)

	ones = 0
	twos = 0

	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			//		fmt.Print(oS[x][y])
			if oS[x][y]&1 == 1 {
				ones++
			} else if oS[x][y]&2 == 2 {
				twos++
			}
		}
	}
	//fmt.Printf("\nones: %d twos: %d\n", ones, twos)
	score -= int(10*ones + twos)

	return score
}

// note: we return CacheEntry because it has all the information we need to return to a higher level of negamax
func negamax(node *AINode, depth int, alpha int, beta int, player int, first bool) (CacheEntry, Move, error) {
	//fmt.Printf("depth: %d\n", depth)
	//fmt.Printf("Movecount: %d ", len(*node.moves))
	//node.moves.Print()
	lastmove := node.moves.LastMove()
	//drawBoard(node.board, lastmove)
	children := genChildren(node.board, &lastmove, node.scores)
	// FIXME: attempt at depth first
//	if len(children) > 4 {
//		children = children[0:len(children) / 2]
//	}

	if depth == 0 || len(children) == 0 {
		//fmt.Printf("FIN depth %d children %d\n",depth, len(children))
		return CacheEntry{
			//FIXME: CHANGE THIS to be correct
			cutoff: alpha,
			depth:  depth,
			score:  playerToMul(player) * evalAIBoard(node.board, node.scores),
			move:   lastmove,
		}, Move{NoMove, NoMove}, nil
	}

	maxScore := SCOREMIN
	var maxEntry CacheEntry

	bestChild := Move{NoMove, NoMove}
	for _, child := range children {
		//fmt.Printf("d%d Trying child %d (%d,%d)\n",depth, i, child.x, child.y)

		// PushMove will panic if it fails (shouldn't fail)
		node.moves.PushMove(child)
		node.board[child.x][child.y] = player

		// NOTE: alpha and beta are negated and swapped for the subcall to negamax
		entry, _, err := negamax(node, depth-1, -beta, -alpha, notPlayer(player), false)
		if err != nil {
			return *new(CacheEntry), Move{NoMove, NoMove}, err
		}
		entry.score = -entry.score
		node.board[child.x][child.y] = B
		node.moves.RemMove()

		if entry.score > maxScore {
			maxScore = entry.score
			maxEntry = entry
			if first {
				bestChild = child
			}
		}
		alpha = max(alpha, entry.score)
		if alpha >= beta {
			break
		}
	}
	return maxEntry, bestChild, nil
}
