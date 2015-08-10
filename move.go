package main

import "fmt"
// a Move is the x,y coordinates of the move, 0 based and starting at the top left of the board
type Move struct {
	x int
	y int
}

const noMove int = 254

func (m Move) isNoMove() bool {
	if m.x == noMove || m.y == noMove {
		return true
	} else {
		return false
	}
}

func NoMove() Move {
	return Move{noMove, noMove}
}

func NewMove() *Move {
	return &Move{
		x: noMove,
		y: noMove,
	}
}

// Takes a move and returns the top left corner of the board that the move was in
func move_to_subboard(m Move) Move {
	return Move{
		x: (m.x/3)*3,
		y: (m.y/3)*3,
	}
}

// Takes human move notation ie. board#,cell# and returns a Move
func move_notation(b, c int) Move {
	b = b-1
	c = c-1
	return Move{
		x: (b%3)*3 + (c%3),
		y: (b/3)*3 + (c/3),
	}
}

func (m *Move) Print() {
	fmt.Printf("[%d, %d]\n", m.x, m.y)
}
