package main

import (
	"fmt"
)


type BitBoard [9]uint32

var pieces = [3]Player{B, X, O}

func (b *Board) toBitBoard() *BitBoard {
	bb := new(BitBoard)

	move := NoMove()

	for x := 0; x < 9; x++ {
		for y := 0; y < 9; y++ {
			move.x = x
			move.y = y

			bmove := move.toBitMove()
			
			bb.applyMove(&bmove, b[x][y])
		}
	}
	
	return bb
}

func (b *BitBoard) applyMove(move *BitMove, player Player){
	subboard := move.s

	if !b.isBlank(move) && player != B{
		fmt.Printf("MOVE %d, %d WAS NOT BLANK\n", move, player)
		drawBoard(b.toBoard(), move.toMove())
		panic("Bad move")
	}
	
	b[subboard] = b[subboard] & ^bitMaskFromMove(move) | bitPlayerMaskFromMove(move, player)
}

func (b *BitBoard) isMove(move *BitMove, player Player) bool {
	subboard := move.s

	mask := bitMaskFromMove(move)
	pmask := bitPlayerMaskFromMove(move, player)

	if player != B {
		fmt.Printf("playr: %d\n", uint32(player))
		fmt.Printf(" mask: %18b\n", mask)
		fmt.Printf("pmask: %18b\n", pmask)
	}

	return b[subboard] & mask == pmask
}

func (b *BitBoard) isBlank(move *BitMove) bool {

	mask := bitMaskFromMove(move)
	sub  := b[move.s]
	mns  := mask & sub
	
	
	// fmt.Printf("\n mask: %18b\n", mask)
	// fmt.Printf("board: %18b\n", sub)
	// fmt.Printf("  mns: %18b\n", mns)

	return mns == uint32(0)
}


func (b *BitBoard) getMove(move *BitMove) Player {
	subboard := move.s

	return Player(b[subboard] & bitMaskFromMove(move) >> uint32(move.c*2))
}

func bitPlayerMaskFromMove(move *BitMove, player Player) uint32 {
	bits := uint32(player)

	return bits << uint32(move.c*2)
}

func bitMaskFromMove(move *BitMove) uint32 {
	bits := uint32(3)

	return bits << uint32(move.c*2)
}


/* Move generation */

func genBitChildren(b *BitBoard, lastmove *BitMove) *BitMoveSlice {
	subscores := subScoresBoard(b)

	return genBitPartialChildren(subscores, b, lastmove)
}
	

func genBitPartialChildren(subscores *BitSubScores, b *BitBoard, lastmove *BitMove) *BitMoveSlice {
	var moves *BitMoveSlice
	if lastmove.isNoMove() {
		return genBitPartialAllChildren(subscores, b, lastmove)
	}

	moves = genBitPartialBoardChildren(subscores, b, lastmove)
	if len(*moves) == 0 {
		return genBitPartialAllChildren(subscores, b, lastmove)
	}
	return moves
}

func genBitPartialBoardChildren(subscores *BitSubScores, b *BitBoard, lastmove *BitMove) *BitMoveSlice {
	moves := make(BitMoveSlice, 0, 9)
	s := lastmove.c

	if subscores[s] != 0 {
		return &moves
	}

	// if(lastmove.isNoMove()){
	// 	won := subscores[s]
	// 	if  won == 1 || won == -1 {
	// 		//fmt.Println("Board is won:", ox, oy)
	// 		return &moves
	// 	}
	// }b

	for cell := uint8(0); cell < 9; cell++ {
		move := BitMove{s: s, c: cell}
		if b.isBlank(&move) {
			moves = append(moves, move)
		}
	}

	return &moves
}

func genBitPartialAllChildren(subscores *BitSubScores, b *BitBoard, lastmove *BitMove) *BitMoveSlice {
	moves := make(BitMoveSlice,0,81)
	var cell uint8

	// FIXME: Horrible horrible hack, need to refactor this to make sense
	fakemove := new(BitMove)
	for cell = 0; cell < 9; cell++ {
		fakemove.c = cell
		moves = append(moves, *genBitPartialBoardChildren(subscores, b, fakemove)...)
		
	}
	return &moves
}

/* Scoring */

var bitX = uint32(X)
var bitO = uint32(O)

func shiftX(dist uint32) uint32 {
	return bitX << (dist*2)
}
func shiftO(dist uint32) uint32 {
	return bitO << (dist*2)
}

func shiftMask(dist uint32) uint32 {
	return 3 << (dist*2)
}

var bitMasklines = []uint32{
	shiftMask(0) + shiftMask(3) + shiftMask(6),
	shiftMask(1) + shiftMask(4) + shiftMask(7),
	shiftMask(2) + shiftMask(5) + shiftMask(8),

	shiftMask(0) + shiftMask(1) + shiftMask(2),
	shiftMask(3) + shiftMask(4) + shiftMask(5),
	shiftMask(6) + shiftMask(7) + shiftMask(8),

	shiftMask(0) + shiftMask(4) + shiftMask(8),
	shiftMask(2) + shiftMask(4) + shiftMask(6),
}

var bitXlines = []uint32{
	shiftX(0) + shiftX(3) + shiftX(6),
	shiftX(1) + shiftX(4) + shiftX(7),
	shiftX(2) + shiftX(5) + shiftX(8),

	shiftX(0) + shiftX(1) + shiftX(2),
	shiftX(3) + shiftX(4) + shiftX(5),
	shiftX(6) + shiftX(7) + shiftX(8),

	shiftX(0) + shiftX(4) + shiftX(8),
	shiftX(2) + shiftX(4) + shiftX(6),
}

var bitOlines = []uint32{
	shiftO(0) + shiftO(3) + shiftO(6),
	shiftO(1) + shiftO(4) + shiftO(7),
	shiftO(2) + shiftO(5) + shiftO(8),

	shiftO(0) + shiftO(1) + shiftO(2),
	shiftO(3) + shiftO(4) + shiftO(5),
	shiftO(6) + shiftO(7) + shiftO(8),

	shiftO(0) + shiftO(4) + shiftO(8),
	shiftO(2) + shiftO(4) + shiftO(6),
}

type SuperLine struct {
	a int
	b int
	c int
}

var superlines = []SuperLine{
	SuperLine{0 , 3 , 6},
	SuperLine{1 , 4 , 7},
	SuperLine{2 , 5 , 8},

	SuperLine{0 , 1 , 2},
	SuperLine{3 , 4 , 5},
	SuperLine{6 , 7 , 8},

	SuperLine{0 , 4 , 8},
	SuperLine{2 , 4 , 6},
}


/**
 * Return the score for a subboard s
 * 1 if X is the winner, -1 if O is the winner, 0 otherwise (doesn't take into account draws)
 * @param s
 */
func scoreBitSubBoard(b *BitBoard, s uint8) int {
		// fmt.Println(len(board[bx:bx+3]))
	// bcols := board[bx:bx+3]
	// b := make([][]Player, 3)
	// for x := 0; x < 3; x++ {
	// 	b[x] = bcols[x][by:by+3]
	// }

	//fmt.Println("Entering evalSubBoard at location", s)

	for i, l := range bitMasklines {
		masked := b[s] & l
		if masked == bitXlines[i]{
			return 1
		} else if masked == bitOlines[i]{
			return -1
		}
	}

	return 0
}
