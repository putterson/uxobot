package main

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
