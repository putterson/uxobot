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

//The tuning value for UCB algorithm
const C = math.Sqrt2

func (m *MonteCarlo) makeMove(move Move) error {
	if m.root == nil {
		return errors.New("root was nil")
	}
	node, ok := m.root.childNodes[move]
	if ok {
		fmt.Println("Found childNode")
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

	for time.Since(start_t).Seconds() < m.timeout {
		count++
		board := origBoard
		player := origPlayer

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
			var optimalUCB float64
			var optimalNode *TreeNode
			var optimalOk bool
			var optimalMove Move

			//Find the largest UCB value for all the moves
			for _, m := range lastNode.childMoves {
				ratio := float64(0.0)
				nextNode, ok := lastNode.childNodes[m]
				nextOutcomes := 0.0
				if ok {
					ratio = float64(nextNode.wincomes) / float64(nextNode.outcomes)
					nextOutcomes = float64(nextNode.outcomes)
				}
				
				ucbval := ratio + C*(math.Sqrt(math.Log(float64(lastNode.outcomes))/nextOutcomes))
				if(ucbval > optimalUCB){
						optimalUCB = ucbval
						optimalNode = nextNode
						optimalMove = m
						optimalOk = ok
				}
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

		//drawBoard(&board, move)
	
		//Simulation phase
		//fmt.Println("Simulation phase")
		//make random moves until the game is over

		simBoard := board
		score := boardScore(simBoard, move, player)
		simdepth := 0
		for !finished(score) {
			simdepth++
			moves := genChildren(&simBoard, &move)
			size_moves := len(moves)

			if size_moves == 0 || simdepth > 81 {
				drawBoard(&simBoard, move)
				fmt.Printf("Error in simulation| simdepth: %d | size_moves: %d\n", simdepth, size_moves)
				return NoMove(), errors.New("")
			}
			
			rnd_move_index := rand.Intn(size_moves)
			move = moves[rnd_move_index]

			simBoard.applyMove(&move, player)

			//make random move on board
			
			score = boardScore(simBoard, move, player)
			player = notPlayer(player)
		}



		scorePlayer := notPlayer(player)

		//Backpropogation
		//fmt.Println("Backpropogation phase")
		if score == 0 {
			for _, node := range nodePath {
				node.outcomes += 1
			}
		} else {
			if scorePlayer != origPlayer {
				score = -score
			}

			if score < 0 {
				score = 0
			}
			
			for _, node := range nodePath {
				node.outcomes += 1
				node.wincomes += score
				score = score ^ 1
			}
		}
		//fmt.Println("Finished Backprop")
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
		ratio := float64(node.wincomes) / float64(node.outcomes)
		fmt.Printf("%d move had ratio of %d / %d = %.2f\nout of %d rounds\n", move, node.wincomes, node.outcomes, ratio, count)
	}
	ratio := float64(optimalNode.wincomes) / float64(optimalNode.outcomes)
	fmt.Printf("optimal move had ratio of %d / %d = %.2f\nout of %d rounds\n", optimalNode.wincomes, optimalNode.outcomes, ratio, count)
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

func montePlayerToMul(origPlayer Player, player Player) int {
	if player == origPlayer {
		return 1
	} else {
		return 0
	}
}
