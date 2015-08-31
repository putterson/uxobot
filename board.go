package main

type Board [9][9]Player

func (b *Board) applyMove(move *Move, player Player) {
	(*b)[move.x][move.y] = player
}

func (b *Board) clearMove(move *Move){
	(*b)[move.x][move.y] = B
}

func genHumanChildren(b *Board, lastmove Move) MoveSlice {
	return genChildren(b, &lastmove)
}

func genChildren(b *Board, lastmove *Move) MoveSlice {
	var moves MoveSlice
	if lastmove.isNoMove() {
		return genAllChildren(b, lastmove)
	}

	moves = genBoardChildren(b, lastmove)
	if len(moves) == 0 {
		return genAllChildren(b, lastmove)
	}
	return moves
}

// Generate children for all the boards
func genAllChildren(b *Board, lastmove *Move) MoveSlice {
	moves := make(MoveSlice,0,81)
	var x,y int

	// FIXME: Horrible horrible hack, need to refactor this to make sense
	fakemove := new(Move)
	for x = 0; x < 3; x++ {
		for y = 0; y < 3; y++ {
			fakemove.x = x
			fakemove.y = y
			moves = append(moves, genBoardChildren(b, fakemove)...)
		}
	}
	return moves
}

// Generate children for a specific board
func genBoardChildren(b *Board, lastmove *Move) MoveSlice {
	moves := make(MoveSlice, 0, 9)
	ox := (lastmove.x%3) * 3
	oy := (lastmove.y%3) * 3
	won := evalSubBoard(b, ox , oy)
	if  won == 1000 || won == -1000 {
	//	fmt.Println("Board is won:", ox, oy)
		return moves
	}
	for x := ox; x < ox+3; x++ {
		for y := oy; y < oy+3; y++ {
			if b[x][y] == B {
				move := Move{x: x, y: y}
				moves = append(moves, move)
			}
		}
	}
	return moves
}
