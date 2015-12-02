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

func (m *BitMove) toMove() Move {
	return move_notation(int(m.s)+1, int(m.c)+1)
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

func (m *BitMove) Print() {
	fmt.Printf("[%d, %d]\n", m.s, m.c)
}
