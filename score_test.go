package main

import "testing"

func sampleBoard() *BitBoard {
	board := new(BitBoard)

	//Won by X
	board.applyMove(&BitMove{0,0},X)
	board.applyMove(&BitMove{0,4},X)
	board.applyMove(&BitMove{0,8},X)

	//Won by O
	board.applyMove(&BitMove{2,2},O)
	board.applyMove(&BitMove{2,4},O)
	board.applyMove(&BitMove{2,6},O)
	
	//Full tied board
	board.applyMove(&BitMove{1,0},X)
	board.applyMove(&BitMove{1,1},X)
	board.applyMove(&BitMove{1,2},O)
	board.applyMove(&BitMove{1,3},O)
	board.applyMove(&BitMove{1,4},O)
	board.applyMove(&BitMove{1,5},X)
	board.applyMove(&BitMove{1,6},X)
	board.applyMove(&BitMove{1,7},X)
	board.applyMove(&BitMove{1,8},O)

	//One move to win X
	board.applyMove(&BitMove{3,2},X)
	board.applyMove(&BitMove{3,4},X)
	
	//One move to win O
	board.applyMove(&BitMove{4,2},O)
	board.applyMove(&BitMove{4,4},O)

	//Won by X
	board.applyMove(&BitMove{6,0},X)
	board.applyMove(&BitMove{6,4},X)
	board.applyMove(&BitMove{6,8},X)
	
	return board
}

func sampleBoardSubScores() *BitSubScores {
	return &BitSubScores{1, 0, -1, 0, 0, 0, 1, 0, 0}
}

func TestSubScoresBoard(t *testing.T) {
	board := sampleBoard()
	expected := sampleBoardSubScores()	

	subscores := subScoresBoard(board)

	for i:=0; i < 9; i++ {
		if expected[i] != subscores[i] {
			t.Fail()
		}
	}
}
                 
func TestBoardPartialScore(t *testing.T) {
	board := sampleBoard()
	expected := sampleBoardSubScores()
	realSubscores := subScoresBoard(board)

	//This move should cause a win in board 3
	move := &BitMove{3,6}
	board.applyMove(move, X)

	score := boardPartialScore(realSubscores, board, move)

	if expected[3] == realSubscores[3] {
		t.Fail()
		t.Logf("Recalculated Scores are %d\n", subScoresBoard(board))
		
	}

	if score == 0 {
		t.Fail()
		t.Log("Game winning move did not win")
	}


	t.Logf("old Scores are %d\n", expected)
	t.Logf("new Scores are %d\n", realSubscores)
}
