package main

import (
	"math/rand"
	"math"
)

type NegaMax struct {
	zobrist_keys [4][9][9][3]BHash
	ai_cache AICache
	depth int
}


//AI types and structs
type BHash uint64
type BHashes [4]BHash

type AICache map[BHash]CacheEntry

type CacheEntry struct {
	depth  int
	score  int
	flag   int
	move   Move
	player Player
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

func (bot *NegaMax) makeMove(move Move) error {
	return nil
}

func (bot *NegaMax) getMove(board Board, lastmove Move, player Player) (Move, error) {
	moves := make(MoveSlice, 0, bot.depth + 1)

	bot.init_zobrist_keys()
	//Clear the cache between moves
	bot.ai_cache = make(AICache)
	node := new(AINode)
	node.board = &board
	node.scores = getSuperScores(&board)
	node.moves = &moves
	node.hashes = bot.hash_board(&board)

	(*node.moves).PushMove(lastmove)
	
	_, move, err := bot.negamax(node, bot.depth, SCOREMIN - 1, SCOREMAX + 1, player, true)


	return move, err
}

func (bot *NegaMax) setDepth(depth int) {
	bot.depth = depth
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


// func orderChildren(node *AINode, moves *MoveSlice, player int){
// 	moveScoreSlice := getMoveScores(node, moves, player)
// 	sort.Sort(moveScoreSlice)
// 	// fmt.Println("Moves:")
// 	// moveScoreSlice.Print()
// 	for i, movescore := range moveScoreSlice {
// 		(*moves)[i] = movescore.move
// 	}
// }

// func getMoveScores(node *AINode, moves *MoveSlice, player int) MoveByScore {
// 	moveScoreSlice := make(MoveByScore, len(*moves))

// 	for i, move := range *moves {
// 		moveScoreSlice[i].move = move
// 		hash := node.hashes ^ zobrist_keys[0][move.x][move.y][player]
// 		cache_entry, exists := ai_cache[hash]
// 		if exists {
// 			moveScoreSlice[i].score = cache_entry.score
// 		} else {
// 			node.board[move.x][move.y] = player
// 			moveScoreSlice[i].score = evalAIBoard(node)
// 			node.board[move.x][move.y] = B
// 			// moveScoreSlice[i].score = 0
// 		}
// 	}

// 	return moveScoreSlice
// }

// note: we return CacheEntry because it has all the information we need to return to a higher level of negamax
func (bot *NegaMax) negamax(node *AINode, depth int, alpha int, beta int, player Player, first bool) (CacheEntry, Move, error) {
	//fmt.Printf("depth: %d\n", depth)
	//fmt.Printf("Movecount: %d ", len(*node.moves))
	//node.moves.Print()

	originalAlpha := alpha

	hash := node.hashes

	lastmove := node.moves.LastMove()

	cache_entry, exists := bot.ai_cache[hash]

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
			delete(bot.ai_cache, hash)
		}
	}


	children := genChildren(node.board, &lastmove)
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
		node.hashes ^= bot.zobrist_keys[0][child.x][child.y][player]

		// Store subboard score to be updated
		childSubBoard := move_to_subboard(child)
		modscore := node.scores[childSubBoard.x/3][childSubBoard.y/3]
		

		// NOTE: alpha and beta are negated and swapped for the subcall to negamax
		entry, _, err := bot.negamax(node, depth-1, -beta, -alpha, notPlayer(player), false)
		if err != nil {
			return *new(CacheEntry), NoMove(), err
		}

		entry.score = -entry.score
		bot.ai_cache[node.hashes] = entry

		//Undo all the updating that happened
		node.hashes ^= bot.zobrist_keys[0][child.x][child.y][player]
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

func (bot *NegaMax) hash_cell(board *Board, x int, y int, orientation int) BHash {
	var player Player = board[x][y]
	var hash BHash
	hash = bot.zobrist_keys[orientation][x][y][player]
	return hash
}

func (bot *NegaMax) hash_board(board *Board) BHash{
        var hash BHash = 0
        for x := 0; x < 9; x++ {
                for y := 0; y < 9; y++ {
                        hash ^= bot.hash_cell(board, x, y, 0)
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

func (bot *NegaMax) init_zobrist_keys(){
	for ori, a := range bot.zobrist_keys {
		for x, b := range a {
			for y, c := range b {
				for player, _ := range c {
					bot.zobrist_keys[ori][x][y][player] = randBHash()
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

func max(a int, b int) int{
	if a > b{
		return a
	}
	return b
}

func min(a int, b int) int{
	if a < b{
		return a
	}
	return b
}
