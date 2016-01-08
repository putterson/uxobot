package main

import "testing"
//import "fmt"

func TestIsBlank(t *testing.T) {
	for mx := 0; mx < 9; mx++ {
		for my := 0; my < 9; my++ {
			b := new(Board)
			m := Move{x: mx, y: my}
			// drawBoard(b, m)
			bb := b.toBitBoard()
			bm := m.toBitMove()


			//Make sure a blank isBlank()
			if !bb.isBlank(&bm) {
				t.Fail()
			}

			bb.applyMove(&bm, O)

			//Make sure after that the move isn't blank
			if bb.isBlank(&bm) {
				t.Fail()
			}

		}
	}
}

func TestGetMove(t *testing.T) {

}

func TestGenBitChildren(t *testing.T) {
	lastmove := BitMove{s: 8, c: 0}
	board := new(BitBoard)

	for ms := uint8(0); ms < 9; ms++ {
		lastmove.c = ms
		for mc := uint8(0); mc < 9; mc++ {
			move := BitMove{s: ms, c: mc}

			moveslice := NewBitMoveSlice()
			slen := 0
			genBitChildren(board, &lastmove, moveslice, &slen)
			board.applyMove(&move, X)

			if slen > 9 - int(mc) {
				t.Fail()
			}
		}
	}
}

func TestToFromBoard(t *testing.T) {
	board := new(Board)

	for xy := 0; xy < 9; xy++ {
		move := Move{x: xy, y: xy}
		board.applyMove(&move, X)
	}

	drawBoard(board, NoMove())

	
	drawBoard(board.toBitBoard().toBoard(), NoMove())
}
