package main

//import "fmt"

type Board [9][9]Player

func (b *BitBoard) toBoard() *Board {
	newboard := new(Board)

	for s := uint8(0); s < 9; s++ {
		for c := uint8(0); c < 9; c++ {
			move := move_notation(int(s+1), int(c+1))
			newboard[move.x][move.y] = b.getMove(&BitMove{s: s, c: c})
		}
	}

	return newboard
}

func (b *Board) applyMove(move *Move, player Player) {
	if (*b)[move.x][move.y] != B {
		panic("Bad move")
	}
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
	won := scoreSubBoard(b, ox , oy)
	if  won == 1 || won == -1 {
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



func genPartialChildren(subscores *SubScores, b *Board, lastmove *Move) *MoveSlice {
	var moves *MoveSlice
	if lastmove.isNoMove() {
		return genPartialAllChildren(subscores, b, lastmove)
	}

	moves = genPartialBoardChildren(subscores, b, lastmove)
	if len(*moves) == 0 {
		return genPartialAllChildren(subscores, b, lastmove)
	}
	return moves
}

func genPartialBoardChildren(subscores *SubScores, b *Board, lastmove *Move) *MoveSlice {
	moves := make(MoveSlice, 0, 9)
	sx := lastmove.x%3
	sy := lastmove.y%3
	won := subscores[sx][sy]
	

	if  won == 1 || won == -1 {
		//fmt.Println("Board is won:", ox, oy)
		return &moves
	}

	ox := sx * 3
	oy := sy * 3

	for x := ox; x < ox+3; x++ {
		for y := oy; y < oy+3; y++ {
			if b[x][y] == B {
				move := Move{x: x, y: y}
				moves = append(moves, move)
			}
		}
	}
	return &moves
}

func genPartialAllChildren(subscores *SubScores, b *Board, lastmove *Move) *MoveSlice {
	moves := make(MoveSlice,0,81)
	var x,y int

	// FIXME: Horrible horrible hack, need to refactor this to make sense
	fakemove := new(Move)
	for x = 0; x < 3; x++ {
		for y = 0; y < 3; y++ {
			fakemove.x = x
			fakemove.y = y
			moves = append(moves, *genPartialBoardChildren(subscores, b, fakemove)...)
		}
	}
	return &moves
}


func scoreSubBoard(b *Board, bx int, by int) int {
	// fmt.Println(len(board[bx:bx+3]))
	// bcols := board[bx:bx+3]
	// b := make([][]Player, 3)
	// for x := 0; x < 3; x++ {
	// 	b[x] = bcols[x][by:by+3]
	// }

	//fmt.Println("Entering evalSubBoard at location",bx,by)
	for _, l := range winlines {
		n := b[l.x1+bx][l.y1+by] | b[l.x2+bx][l.y2+by] | b[l.x3+bx][l.y3+by]
		s := b[l.x1+bx][l.y1+by] + b[l.x2+bx][l.y2+by] + b[l.x3+bx][l.y3+by]

		if s == 0 || n == 3 {
			continue
		}
		
		// if the whole line is one player
		if s == 3 {
			return 1
		} else if s == 6 {
			return -1
		}

	}

	return 0
}
