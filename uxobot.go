package main

import (
	"fmt"
)

// Constants for the piece values (blank, X, O)
const (
	B = 0
	X = 1
	O = 2
)

var maptochar = []rune{' ', 'X', 'O'}
var playerconf = []string{"human", "human"}
var wins = [][]uint8{
	{1, 1, 1, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 1, 1, 1, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 1, 1, 1},
	{1, 0, 0, 0, 1, 0, 0, 0, 1},
	{0, 0, 1, 0, 1, 0, 1, 0, 0},
	{1, 0, 0, 1, 0, 0, 1, 0, 0},
	{0, 1, 0, 0, 1, 0, 0, 1, 0},
	{0, 0, 1, 0, 0, 1, 0, 0, 1},
}

type Board [81]uint8

type Won [9]uint8

type Node struct {
	board Board
	won   Won
}

type Move struct {
	x int
	y int
}

func NewMove() *Move {
	return &Move{
		x: -1,
		y: -1,
	}
}

func main() {
	run()
}

func run() {
	b := new(Board)
	for i := range b {
		b[i] = B
	}
	var curplayer = X
	var lastmove Move
	var move Move
	for {
		move = makeMove(b, curplayer, lastmove)
		drawBoard(b, move)
		curplayer = notPlayer(curplayer)
	}
}

/*
*
* Utility functions for the different move notations
*
*/

func move_to_pos(move Move) int {
	return move.y*9 + move.x
}

func xy_to_pos(x, y int) int {
	//fmt.Printf("%d", x*9+y)
	return x*9 + y
}

func pos_to_move(pos int) Move {
	return Move{
		x: pos / 9,
		y: pos % 9,
	}
}

// Takes human move notation ie. board#,cell# and returns a Move
func move_notation(b, c int) Move {
	b = b-1
	c = c-1

	return Move{}
}

func makeMove(b *Board, curplayer int, move Move) (newmove Move) {
	if playerconf[curplayer-1] == "human" {
		newmove = getMove(b, move)
	}
	return
}

func getMove(b *Board, lastmove Move) (move Move) {
	move = *(new(Move))
	for {
		fmt.Printf("[move]>")
		var _, succ = fmt.Scanf("%d%d", &move.x, &move.y)
		if succ == nil {
			if move.x > 0 && move.x < 10 && move.y > 0 && move.y < 10 {
				for _, vm := range genHumanChildren(b, move) {
					if move == vm {
						return
					}
				}
			}
		}
		fmt.Printf("Please make a valid move of the form \"0 9\"\n")
	}
}

func genHumanChildren(b *Board, lastmove Move) []Move {
	return genChildren(b, lastmove, new(Won))
}

func genChildren(b *Board, lastmove Move, mboard *Won) []Move {
	var moves []Move
	if (lastmove.x == -1) && (lastmove.y == -1) {
		return genAllChildren(b, lastmove)
	}

	moves = genBoardChildren(b, lastmove)
	if len(moves) == 0 {
		return genAllChildren(b, lastmove)
	}
	return moves
}

// Generate children for all the boards
func genAllChildren(b *Board, lastmove Move) []Move {
	moves := []Move{}
	for pos := 0; pos < 81; pos++ {
		if b[pos] == B {
			moves = append(moves, pos_to_move(pos))
		}
	}
	return moves
}

// Generate children for a specific board
func genBoardChildren(b *Board, lastmove Move) []Move {
	moves := []Move{}
	ox := (lastmove.x % 3) * 3 * 9
	oy := (lastmove.y % 3) * 3
	for x := ox; x <= ox+3; x++ {
		for y := oy; y <= oy+3; y++ {
			if b[xy_to_pos(x, y)] == B {
				move := Move{x: x, y: y}
				moves = append(moves, move)
			}
		}
	}

	return moves
}

func notPlayer(p int) int {
	if p == X {
		return O
	} else {
		return X
	}
}

func drawBoard(b *Board, move Move) {
	for x := 0; x < 9; x++ {
		for y := 0; y < 9; y++ {
			fmt.Printf("%c ", maptochar[(*b)[y*9+x]])
		}
		fmt.Printf("\n\n")
	}
	return
}
