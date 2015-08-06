package main

import (
	"fmt"
	"math/rand"
	"math"
	"sort"
)

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

var zobrist_keys [4][9][9][3]BHash
var ai_cache AICache = make(AICache)

//AI types and structs
type BHash uint64
type BHashes [4]BHash

type AICache map[BHash]CacheEntry

type CacheEntry struct {
	depth  int
	score  int
	flag   int
	move   Move
	player int
}

type AINode struct {
	board  *Board
	moves  *MoveSlice
	cache  AICache
	scores *Scores
	hashes   BHash
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

//Return the score as a float from 0 to 1, or NaN in the case of a drawn board
func normSubScore(score int) float64 {
	if score == SCOREDRAW {
		return math.NaN()
	} else {
		return float64(score + 1000) / 2000.0
	}
}

func evalAIBoard(node *AINode) int {
	//Really only need to check one of the axes
	// If there was a last move we only need to update one of the subboard scores (which the move was made in)
	// return evalSuperBoard(getSuperScores(node.board))
	
	if !(*node.moves).LastMove().isNoMove() {
		subboard := move_to_subboard((*node.moves).LastMove())
		node.scores[subboard.x/3][subboard.y/3] = normSubScore(evalSubBoard(node.board, subboard.x, subboard.y))
	}

	return evalSuperBoard(node.scores)
}

func getSuperScores(board *Board) *Scores {
	floatscores := new(Scores)
	for bx := 0; bx < 3; bx++ {
		for by := 0; by < 3; by++ {
			floatscores[bx][by] = normSubScore(evalSubBoard(board, bx*3, by*3))
		}
	}
	return floatscores
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

func spreadOfSlice(slice *[]float64) float64 {
	maxEntry := 0.0
	minEntry := 0.0

	for _, entry := range *slice {
		maxEntry = math.Max(maxEntry, entry)
		minEntry = math.Min(minEntry, entry)
	}

	//minEntry must be <= 0 and maxEntry >= 0
	return minEntry + maxEntry
}

// evaluate the score of a sub-board always with regard to X
// bx and by are the top left corner of the subboard to score	
func evalSubBoard(board *Board, bx int, by int) int {
	// fmt.Println(len(board[bx:bx+3]))
	bcols := board[bx:bx+3]
	b := make([][]int, 3)
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

func orderChildren(node *AINode, moves *MoveSlice, player int){
	moveScoreSlice := getMoveScores(node, moves, player)
	sort.Sort(moveScoreSlice)
	// fmt.Println("Moves:")
	// moveScoreSlice.Print()
	for i, movescore := range moveScoreSlice {
		(*moves)[i] = movescore.move
	}
}

func getMoveScores(node *AINode, moves *MoveSlice, player int) MoveByScore {
	moveScoreSlice := make(MoveByScore, len(*moves))

	for i, move := range *moves {
		moveScoreSlice[i].move = move
		hash := node.hashes ^ zobrist_keys[0][move.x][move.y][player]
		cache_entry, exists := ai_cache[hash]
		if exists {
			moveScoreSlice[i].score = cache_entry.score
		} else {
			node.board[move.x][move.y] = player
			moveScoreSlice[i].score = evalAIBoard(node)
			node.board[move.x][move.y] = B
			// moveScoreSlice[i].score = 0
		}
	}

	return moveScoreSlice
}

// note: we return CacheEntry because it has all the information we need to return to a higher level of negamax
func negamax(node *AINode, depth int, alpha int, beta int, player int, first bool) (CacheEntry, Move, error) {
	//fmt.Printf("depth: %d\n", depth)
	//fmt.Printf("Movecount: %d ", len(*node.moves))
	//node.moves.Print()

	originalAlpha := alpha

	hash := node.hashes

	lastmove := node.moves.LastMove()

	cache_entry, exists := ai_cache[hash]

	if exists == true && !first {
		if cache_entry.depth >= depth {
			if cache_entry.flag == SCORE_EXACT {
				// fmt.Println("Cache entry exact. depth ", depth)
				return cache_entry, NoMove(), nil
			} else if cache_entry.flag == SCORE_LOWER_BOUND {
				alpha = max(alpha, cache_entry.score)
			} else if cache_entry.flag == SCORE_UPPER_BOUND {
				beta = min(beta, cache_entry.score)
			}

			if alpha >= beta {
				// fmt.Println("Cache entry a > b. depth ", depth)
				return cache_entry, NoMove(), nil
			}
		} else {
			delete(ai_cache, hash)
		}
	}


	children := genChildren(node.board, &lastmove, node.scores)
	//TODO: Order moves to increase performance of pruning
	//orderChildren(node, &children, player)

	//TODO: check for won game or finished search
	curScore := evalAIBoard(node)
	
	if depth == 0 || len(children) == 0 || won(curScore) {
//		fmt.Printf("FIN depth %d children %d\n",depth, len(children))
		return CacheEntry{
			//FIXME: CHANGE THIS to be correct
			depth:  depth,
			score:  playerToMul(player) * curScore,
			move:   lastmove,
			player: player,
		}, NoMove(), nil
	}

	maxScore := SCOREMIN - 1
	var maxEntry CacheEntry

	bestChild := NoMove()
	for _, child := range children {
		//fmt.Printf("d%d Trying child %d (%d,%d)\n",depth, i, child.x, child.y)

		node.moves.PushMove(child)
		node.board[child.x][child.y] = player
		node.hashes ^= zobrist_keys[0][child.x][child.y][player]

		// Store subboard score to be updated
		childSubBoard := move_to_subboard(child)
		modscore := node.scores[childSubBoard.x/3][childSubBoard.y/3]
		

		// NOTE: alpha and beta are negated and swapped for the subcall to negamax
		entry, _, err := negamax(node, depth-1, -beta, -alpha, notPlayer(player), false)
		if err != nil {
			return *new(CacheEntry), NoMove(), err
		}

		entry.score = -entry.score
		ai_cache[node.hashes] = entry

		//Undo all the updating that happened
		node.hashes ^= zobrist_keys[0][child.x][child.y][player]
		node.board[child.x][child.y] = B
		node.scores[childSubBoard.x/3][childSubBoard.y/3] = modscore
		node.moves.RemMove()

		if entry.score > maxScore || (entry.score == maxScore && entry.depth > maxEntry.depth){
			maxScore = entry.score
			maxEntry = entry
			//We only have to return the best move if at the top level negamax call
			if first {
				bestChild = child
			}
		}
		alpha = max(alpha, entry.score)
		if alpha >= beta {
			break
		}
	}

	if maxEntry.score <= originalAlpha {
		maxEntry.flag = SCORE_UPPER_BOUND
	}else if maxEntry.score >= beta {
		maxEntry.flag = SCORE_LOWER_BOUND
	}else {
		maxEntry.flag = SCORE_EXACT
	}
	//TODO: suspect line
	maxEntry.depth = depth

	//TODO: remember to use canonical hash here
	// ai_cache[node.hashes] = maxEntry

	return maxEntry, bestChild, nil
}

func hash_cell(board *Board, x int, y int, orientation int) BHash {
	var player int = board[x][y]
	var hash BHash
	hash = zobrist_keys[orientation][x][y][player]
	return hash
}

func hash_board(board *Board) BHash{
        var hash BHash = 0
        for x := 0; x < 9; x++ {
                for y := 0; y < 9; y++ {
                        hash ^= hash_cell(board, x, y, 0)
                }
        }
        


	return hash
}

func canon_hash(hashes BHashes) BHash{
	var canonicalHash BHash = 0
	for _, hash := range hashes {
		canonicalHash = minBHash(canonicalHash, hash)
	}
	return canonicalHash
}

func minBHash(h1 BHash, h2 BHash) BHash {
	if h1 < h2 {
		return h1
	}else {
		return h2
	}
}

func init_zobrist_keys(){
	for ori, a := range zobrist_keys{
		for x, b := range a {
			for y, c := range b {
				for player, _ := range c {
					zobrist_keys[ori][x][y][player] = randBHash()
				}
			}
		}
	}
}

func print_zobrist_keys(){
	for ori, a := range zobrist_keys{
		for x, b := range a {
			for y, c := range b {
				for player, _ := range c {
					fmt.Println(zobrist_keys[ori][x][y][player])
				}
			}
		}
	}
}

func randBHash() BHash {
	return BHash(rand.Uint32())<<32 + BHash(rand.Uint32())
}

func won(score int) bool {
	if score == SCOREMAX || score == SCOREMIN {
		return true
	}
	return false
}
