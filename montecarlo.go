package main

import (
	"fmt"
	"math/rand"
	"time"
)

type MonteCarlo struct {
	timeout float64
	root    *TreeNode
}

func NewMonteCarlo(board *Board, lastmove *Move, timeout float64) *MonteCarlo {
	return &MonteCarlo{
		timeout: timeout,
		root:    NewTreeNode(board, lastmove),
	}
}

type NodeChildren map[Move]*TreeNode

type TreeNode struct {
	outcomes   int
	wincomes   int
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

//Populates the child nodes
func (t *TreeNode) genChildren(board *Board, lastmove *Move) int {
	//	getChildren(
	return 0
}

func NewTreeNode(board *Board, lastmove *Move) *TreeNode {
	return &TreeNode{
		outcomes:   0,
		wincomes:   0,
		childMoves: genChildren(board, lastmove),
		childNodes: make(NodeChildren),
	}
}

func getLastNode(nodePath []*TreeNode) *TreeNode {
	return nodePath[len(nodePath)-1]
}

func (m MonteCarlo) getMove(board Board, lastmove Move, player Player) (Move, error) {
	start_t := time.Now()

	origBoard := board
	count := 0

	for time.Since(start_t).Seconds() < m.timeout {
		//player
		count += 1
		board := origBoard
		origPlayer := player

		var nodePath []*TreeNode
		nodePath = make([]*TreeNode, 0, 81)
		nodePath = append(nodePath, m.root)

		var move Move
		move = NoMove()

		//Selection phase
		//fmt.Println("Selection phase")
		//While the node has visited children move to a selected child
		var lastNode *TreeNode
		for true {
			lastNode = getLastNode(nodePath);
			size_moves := lastNode.nNextMoves()
			idx_rand_move := rand.Intn(size_moves)
			//Will be using UCB method here instead of random selection
			move = lastNode.getMove(idx_rand_move)


			board.applyMove(&move, player)
			player = notPlayer(player)
			nextNode, ok := lastNode.childNodes[move]
			if !ok {
				//Go to Extension phase
				//fmt.Println("Extension phase")
				nextNode := NewTreeNode(&board, &move)
				lastNode.childNodes[move] = nextNode
				nodePath = append(nodePath, nextNode)
				break
			} else {
				nodePath = append(nodePath, nextNode)
			}
		}

		//drawBoard(&board, move)
	
		
		//fmt.Printf("NodePath length: %d\n", len(nodePath))
		// //Extension phase
		// fmt.Println("Extension phase")
		// size_moves := lastNode.nNextMoves()
		// idx_rand_move := rand.Intn(size_moves)
		// move = lastNode.getMove(idx_rand_move)
		// board.applyMove(&move, player)
		// // move.Print()
		// nextNode := NewTreeNode(&board, &move)
		// lastNode.childNodes[move] = nextNode

		// fmt.Printf("Length of childNodes: %d\n",len(getLastNode(nodePath).childNodes))

		//player = notPlayer(player)
		//nodePath = append(nodePath, nextNode)

		//getLastNode(nodePath).childMoves.Print()

		//Simulation phase
		//fmt.Println("Simulation phase")
		//make random moves until the game is over

		simBoard := board
		score := boardScore(simBoard, move, player)

		for !finished(score) {
			moves := genChildren(&board, &move)
			size_moves := len(moves)
			rnd_move_index := rand.Intn(size_moves)
			move = moves[rnd_move_index]

			simBoard.applyMove(&move, player)

			//make random move on board
			score = boardScore(simBoard, move, player)

			player = notPlayer(player)
		}

		//drawBoard(&simBoard, move)


		player = origPlayer

		//Backpropogation
		//fmt.Println("Backpropogation phase")
		for _, node := range nodePath {
			node.outcomes += 1
			node.wincomes += (1 & score) & montePlayerToMul(player)
			player = notPlayer(player)
		}
		//fmt.Println("Finished Backprop")
	}
	fmt.Println("Ran out of time")

	//Find optimal toplevel move
	var ratio, newRatio float64
	ratio = -1.0
	var optimalMove Move
	var optimalNode *TreeNode
	optimalMove = NoMove()

	fmt.Println("Eval final move")
	for move, node := range m.root.childNodes {

		//move.Print()
		newRatio = (float64(node.wincomes) / float64(node.outcomes))
		if newRatio > ratio {
			optimalMove = move
			optimalNode = node
			ratio = newRatio
		}
	}
	fmt.Printf("optimal move had ratio of %d / %d = %.2f\nout of %d rounds", optimalNode.wincomes, optimalNode.outcomes, ratio, count)
	return optimalMove, nil
}

func finished(score int) bool {
	return score >= -1
}

func boardScore(board Board, lastmove Move, player Player) int {
	score := evalBoard(&board)

	if score == SCOREMAX {
		return playerToMul(player)
	} else if score == SCOREMIN {
		return playerToMul(notPlayer(player))
	} else {
		if len(genChildren(&board, &lastmove)) == 0 {
			return 0
		} else {
			return -2
		}
	}
}

func montePlayerToMul(player Player) int {
	if player == X {
		return 1
	} else {
		return 0
	}
}
