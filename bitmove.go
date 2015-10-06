package main

import "fmt"
// a BitMove is of the form [subboard, cell] from 0,0 (top left) to 9,9 (bottom right)
type BitMove struct {
	//subboard
	s uint8
	//cell
	c uint8
}

const noBitMove uint8 = uint8(254)

func (m *Move) toBitMove() BitMove {
	if m.isNoMove() {
		return NoBitMove()
	}
	
	s := uint8((m.x / 3) + 3*(m.y / 3))
	c := uint8((m.x % 3) + 3*(m.y % 3))

	return BitMove{s: s, c: c}
}

func (m *BitMove) isNoMove() bool {
	if m.s == noBitMove || m.c == noBitMove {
		return true
	} else {
		return false
	}
}

func NoBitMove() BitMove {
	return BitMove{noBitMove, noBitMove}
}

func NewBitMove() *BitMove {
	return &BitMove{
		s: noBitMove,
		c: noBitMove,
	}
}

func (m *BitMove) Print() {
	fmt.Printf("[%d, %d]\n", m.s, m.c)
}
