package main

type BitBoard [9]uint32

func (b *BitBoard) applyMove(move *Move, player Player){
	subboard := bitSubBoard(move)
	
	b[subboard] = b[subboard] | bitMaskFromMove(move, player)
}

func bitSubBoard(move *Move) int {
	x := move.x / 3
	y := move.y / 3

	return x + 3*y
}

func bitMaskFromMove(move *Move, player Player) uint32 {

	shift := int(player * 9)
	
	x := move.x % 3
	y := move.y % 3

	return uint32(1) << uint32(shift + x + 3*y)
}


