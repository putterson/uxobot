package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
	"errors"
)

type MonteCarlo struct {
	timeout float64
	root    *TreeNode
}

func NewMonteCarlo(timeout float64) *MonteCarlo {
	return &MonteCarlo{
		timeout: timeout,
		root:    nil,
	}
}

type NodeChildren map[Move]*TreeNode

type TreeNode struct {
	outcomes   int
	wincomes   [2]int
	move       Move
	childMoves MoveSlice
	childNodes NodeChildren
}

func (t *TreeNode) hasChildNodes() bool {
	return len(t.childNodes) > 0
}

func (t *TreeNode) nNextMoves() int {
	return len(t.childMoves)
}

func (t *TreeNode) getMove(n int) Move {
	return t.childMoves[n]
}

func NewTreeNode(board *Board, lastmove *Move) *TreeNode {
	return &TreeNode{
		outcomes:   1,
		wincomes:   [2]int{0,0},
		childMoves: genChildren(board, lastmove),
		childNodes: make(NodeChildren),
	}
}

func getLastNode(nodePath []*TreeNode) *TreeNode {
	return nodePath[len(nodePath)-1]
}

//The tuning value for UCB algorithm
const C = 1.1 //math.Sqrt2

func (m *MonteCarlo) makeMove(move Move) error {
	if m.root == nil {
		return errors.New("root was nil")
	}
	node, ok := m.root.childNodes[move]
	if ok {
		//fmt.Println("Found childNode")
		m.root = node
		return nil
	} else {
		return errors.New("Move was not available")
	}
}

func (m *MonteCarlo) getMove(board Board, lastmove Move, player Player) (Move, error) {
	if m.root == nil {
		fmt.Println("getMove with nil root")
		m.root = NewTreeNode(&board, &lastmove)
	}
	
	start_t := time.Now()

	//Save the original board and player
	origBoard := board
	origPlayer := player
	
	count := 0

	for (time.Since(start_t).Seconds() < m.timeout) {
		count++
		player := origPlayer

		board, player, move, nodePath := m.selection(origBoard, origPlayer)

		score := m.simulate(board, player, move)

		//drawBoard(&board, move)
	


		m.backpropagate(nodePath, score, origPlayer)
	}
	//fmt.Println("Ran out of time")

	//Find optimal toplevel move
	var visits, newVisits int
	var optimalMove Move
	var optimalNode *TreeNode
	optimalMove = NoMove()

	//fmt.Println("Eval final move")
	//Final move should be the move with the most visits, not the best ratio
	for move, node := range m.root.childNodes {

		//move.Print()
		newVisits = node.outcomes
		if newVisits > visits {
			optimalMove = move
			optimalNode = node
			visits = newVisits
		}
		ratio := float64(node.wincomes[player-1]) / float64(node.outcomes)
		fmt.Printf("%d move had ratio of %d / %d = %.2f\n", move, node.wincomes, node.outcomes, ratio)
	}
	ratio := float64(optimalNode.wincomes[player-1]) / float64(optimalNode.outcomes)
	fmt.Printf("optimal move %d had ratio of %d / %d = %.2f\nout of %d rounds\n", optimalMove, optimalNode.wincomes, optimalNode.outcomes, ratio, count)
	return optimalMove, nil
}

func finished(score int) bool {
	return score != 0
}

func boardScore(board Board, lastmove Move) int {
	score := scoreBoard(&board)

	if score == 0 {
		if len(genChildren(&board, &lastmove)) == 0 {
			return -2
		} else {
			return 0
		}
	} else {
		return score
	}
}

func montePlayerToMul(origPlayer Player, player Player) int {
	if player == origPlayer {
		return 1
	} else {
		return 0
	}
}

func (m *MonteCarlo) selection(board Board, player Player) (Board, Player, Move, []*TreeNode) {
	var nodePath []*TreeNode
	nodePath = make([]*TreeNode, 0, 81)
	nodePath = append(nodePath, m.root)

	var move Move
	nomove := NoMove()
	move = NoMove()

	//Selection phase
	//fmt.Println("Selection phase")
	//While the node has visited children move to a selected child
	var lastNode *TreeNode
	for !finished(boardScore(board,move)) {
		lastNode = getLastNode(nodePath);
		optimalUCB := math.Inf(-1)
		var optimalNode *TreeNode
		var optimalOk bool
		optimalMove := NoMove()

		//Find the largest UCB value for all the moves
		for _, m := range lastNode.childMoves {
			ratio := float64(0.0)
			nextNode, ok := lastNode.childNodes[m]
			nextOutcomes := 1.0
			if ok {
				ratio = float64(nextNode.wincomes[player-1]) / float64(nextNode.outcomes)
				nextOutcomes = float64(nextNode.outcomes)
			}
			
			ucbval := ratio + C*(math.Sqrt(math.Log(float64(lastNode.outcomes))/nextOutcomes))

			// fmt.Println(ucbval)

			//If the move is better for us (or there hasn't been an optimal move yet)
			if ucbval >= optimalUCB || optimalMove == nomove {
				optimalUCB = ucbval
				optimalNode = nextNode
				optimalMove = m
				optimalOk = ok
			}

		}

		if optimalMove == NoMove() {
			fmt.Printf("Bad length of childMoves? %d\n", len(lastNode.childMoves))
			drawBoard(&board, move)
			fmt.Println("Failed to find optimalMove during selection")
		}
		
		
		move = optimalMove
		board.applyMove(&move, player)
		player = notPlayer(player)

		if !optimalOk {
			//Go to Extension phase
			//fmt.Println("Extension phase")
			optimalNode := NewTreeNode(&board, &move)
			lastNode.childNodes[move] = optimalNode
			nodePath = append(nodePath, optimalNode)
			break
		} else {
			nodePath = append(nodePath, optimalNode)
		}
	}

	return board, player, move, nodePath
}

func (m *MonteCarlo) simulate(board Board, player Player, move Move) int {
	//Simulation phase
	//fmt.Println("Simulation phase")
	//make random moves until the game is over

	simBoard := board
	score := boardScore(simBoard, move)

	for !finished(score) {
		moves := genChildren(&simBoard, &move)
		size_moves := len(moves)

		if size_moves == 0 {
			drawBoard(&simBoard, move)
			fmt.Printf("Error in simulation | score: %d | size_moves: %d\n",score, size_moves)
			return 0
		}
		
		//Check for one move away wins
		for _, p_move := range moves {
			simBoard.applyMove(&p_move, player)
			score = boardScore(simBoard, p_move)
			if (score & 1) == 1 {
				//fmt.Printf("Found one move away win score: %d\n", score )
				//drawBoard(&board, p_move)
				return score
			}
			simBoard.clearMove(&p_move)
		}

		//otherwise make random move on board
		rnd_move_index := rand.Intn(size_moves)
		move = moves[rnd_move_index]
		simBoard.applyMove(&move, player)
		
		score = boardScore(simBoard, move)

		player = notPlayer(player)
	}

	return score
}

func (m *MonteCarlo) backpropagate(nodePath []*TreeNode, score int, origPlayer Player) {
	//Backpropogation
	//fmt.Println("Backpropogation phase")
	player := B
	if score == -2 {
		for _, node := range nodePath {
			node.outcomes += 1
		}
		return
	} else if score == 1 {
		player = X
	} else if score == -1 {
		player = O
	} else {
		fmt.Println("bad score during backpropagation")
	}

	for _, node := range nodePath {
		node.outcomes += 1
		node.wincomes[player-1] += 1
	}
	
	//fmt.Println("Finished Backprop")
}

func scoreBoard(board *Board) int {
	subscores := new(SubScores)

	for bx := 0; bx < 3; bx++ {
		for by := 0; by < 3; by++ {
			subscores[bx][by] = scoreSubBoard(board, bx*3, by*3)
		}
	}
	return scoreSuperBoard(subscores)
}

func scoreSuperBoard(b *SubScores) int {
	for _, l := range winlines {
		s := b[l.x1][l.y1] + b[l.x2][l.y2] + b[l.x3][l.y3]

		if s == 3 {
			return 1
		} else if s == -3 {
			return -1
		}
	}

	return 0
}

func scoreSubBoard(board *Board, bx int, by int) int {
	// fmt.Println(len(board[bx:bx+3]))
	bcols := board[bx:bx+3]
	b := make([][]Player, 3)
	for x := 0; x < 3; x++ {
		b[x] = bcols[x][by:by+3]
	}

	//fmt.Println("Entering evalSubBoard at location",bx,by)
	for _, l := range winlines {
		n := b[l.x1][l.y1] | b[l.x2][l.y2] | b[l.x3][l.y3]
		s := b[l.x1][l.y1] + b[l.x2][l.y2] + b[l.x3][l.y3]

		// if the whole line is one player
		if s == 3*n {
			if s == 3 {
				return 1
			} else if s == 6 {
				return -1
			}
		}
	}

	return 0
}
