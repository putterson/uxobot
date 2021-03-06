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

func getLastNode(nodePath []*TreeNode) *TreeNode {
	return nodePath[len(nodePath)-1]
}

//The tuning value for UCB algorithm
const C = math.Sqrt2

func (m *MonteCarlo) makeMove(move Move) error {
	if m.root == nil {
		return errors.New("root was nil")
	}
	node, ok := m.root.getChild(move.toBitMove())
	if ok {
		//fmt.Println("Found childNode")
		m.root = node
		return nil
	} else {
		//Technically this is not an error because montecarlo may just not have visited that node
		//it will continue just fine by setting the root to nil
		m.root = nil
		return nil
	}
}

func (m *MonteCarlo) getMove(board Board, lastmove Move, player Player) (Move, error) {
	origBitBoard := *(board.toBitBoard())
	origBitMove  := (lastmove.toBitMove())
	
	if m.root == nil {
		fmt.Println("getMove with nil root")
		m.root = NewTreeNode(&origBitBoard, &origBitMove)
	}
	
	start_t := time.Now()

	//Save the original board and player
	origPlayer := player
	
	count := 0

	var nodePath NodePath
	nodePath = make(NodePath, 0, 81)
	nodePath = append(nodePath, m.root)

	for (time.Since(start_t).Seconds() < m.timeout) {
		count++
		player := origPlayer

		nodePath = nodePath[:1]
		bitboard, player, move, nodePath := m.selection(nodePath, origBitBoard, origPlayer)

		score := m.simulate(bitboard, player, &move)

		//drawBoard(&board, move)
	


		m.backpropagate(nodePath, score, origPlayer)
	}
	//fmt.Println("Ran out of time")

	//Find optimal toplevel move
	var visits, newVisits int
	var optimalMove BitMove
	var optimalNode *TreeNode
	optimalMove = NoBitMove()

	//fmt.Println("Eval final move")
	//Final move should be the move with the most visits, not the best ratio
	for _, node := range m.root.childNodes {

		//move.Print()
		newVisits = node.outcomes
		if newVisits > visits {
			optimalMove = node.move
			optimalNode = node
			visits = newVisits
		}
		//ratio := float64(node.wincomes[player-1]) / float64(node.outcomes)
//		fmt.Printf("%d move had ratio of %d / %d = %.2f\n", move, node.wincomes, node.outcomes, ratio)
	}

	ratio := float64(optimalNode.wincomes[player-1]) / float64(optimalNode.outcomes)
	//fmt.Printf("optimal move %d had ratio of %d / %d = %.2f\nout of %d rounds\n", optimalMove, optimalNode.wincomes, optimalNode.outcomes, ratio, count)
	fmt.Printf("Examined %d playouts with win probability of %.2f\n", count, ratio)
	return optimalMove.toMove(), nil
}

func finished(score int) bool {
	return score != 0
}

func boardPartialScore(subscores *BitSubScores, board *BitBoard, lastmove *BitMove) int {
	score := scorePartialBoard(subscores, board, lastmove)
	if score == 0 {
		if !areBitPartialChildren(subscores, board, lastmove) {
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

type NodePath []*TreeNode

func (m *MonteCarlo) selection(nodePath NodePath, board BitBoard, player Player) (BitBoard, Player, BitMove, NodePath) {
	var move BitMove
	nomove := NoBitMove()
	move = NoBitMove()


	subscores := subScoresBoard(&board)
	//Selection phase
	//fmt.Println("Selection phase")
	//While the node has visited children move to a selected child
	var lastNode *TreeNode
	for !finished(boardPartialScore(subscores, &board,&move)) {
		lastNode = getLastNode(nodePath);
		optimalUCB := math.Inf(-1)
		var optimalNode *TreeNode
		var optimalOk bool
		optimalMove := NoBitMove()

		//Find the largest UCB value for all the moves
		//fmt.Println("Selection:")
		for i := 0; i < lastNode.nChildMoves; i++ {
			cmove := lastNode.childMoves[i]
			ratio := float64(0.0)
			nextNode, ok := lastNode.getChild(cmove)
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
				optimalMove = cmove
				optimalOk = ok
			}


			//fmt.Printf("%d move had ucb of %.2f\n", cmove, ucbval)
		}
		
		move = optimalMove
		board.applyMove(&move, player)
		player = notPlayer(player)

		//fmt.Printf("Nodepath %d\n",nodePath)

		if !optimalOk {
			//Go to Extension phase
			//fmt.Println("Extension phase")
			optimalNode := lastNode.addChild(&board, &move)
			nodePath = append(nodePath, optimalNode)
			break
		} else {
			nodePath = append(nodePath, optimalNode)
		}
	}

	return board, player, move, nodePath
}

func (m *MonteCarlo) simulate(board BitBoard, player Player, move *BitMove) int {
	//Simulation phase
	//fmt.Println("Simulation phase")
	//make random moves until the game is over

	moveslice := NewBitMoveSlice()
	slen := 0
	simBoard := board
	nomove := NoBitMove()

	subscores := BitSubScores{}
	score := boardPartialScore(&subscores, &simBoard, &nomove)

	for !finished(score) {
		slen = 0
		genBitPartialChildren(&subscores, &simBoard, move, moveslice, &slen)


//		if slen == 0 {
//			drawBoard(simBoard.toBoard(), move.toMove())
//		}

		//Check for one move away wins
		oldscores := subscores
		for i := 0; i < slen; i++ {
			p_move := &(*moveslice)[i]
			simBoard.applyMove(p_move, player)
			score = boardPartialScore(&subscores, &simBoard, p_move)
			if (score & 1) == 1 {
				//fmt.Printf("Found one move away win score: %d\n", score )
				//drawBoard(&board, p_move)
				return score
			}
			simBoard.applyMove(p_move, B)
			subscores = oldscores
		}

		//otherwise make random move on board
		rnd_move_index := rand.Intn(slen)
		move = &(*moveslice)[rnd_move_index]
		simBoard.applyMove(move, player)
		
		score = boardPartialScore(&subscores, &simBoard, move)

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

func scoreBoard(board *BitBoard) int {
	return scoreSuperBoard(subScoresBoard(board))
}

func subScoresBoard(board *BitBoard) *BitSubScores {
	subscores := new(BitSubScores)

	for s:= uint8(0); s < 9; s++{
			subscores[s] = scoreBitSubBoard(board, s)
	}
	return subscores
}

func scorePartialBoard(b *BitSubScores, board *BitBoard, lastmove *BitMove) int {
	if lastmove.isNoMove() {
		for s := uint8(0); s < 9; s++ {
			b[s] = scoreBitSubBoard(board, s)
		}
	}else {

		s := lastmove.s
		b[s] = scoreBitSubBoard(board, s)
	}
	return scoreSuperBoard(b)	
}

func scoreSuperBoard(b *BitSubScores) int {
	for _, l := range superlines {
		s := b[l.a] + b[l.b] + b[l.c]

		if s == 3 {
			return 1
		} else if s == -3 {
			return -1
		}
	}

	return 0
}
