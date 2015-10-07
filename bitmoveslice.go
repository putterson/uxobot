package main

import "fmt"

type BitMoveSlice [81]BitMove

func NewBitMoveSlice() *BitMoveSlice {
	return new(BitMoveSlice)
}

// func (m *BitMoveSlice) PushMove(move BitMove) {
//     moves := *m
//     l := len(moves)
// //     if l == cap(moves) {
// //         panic(fmt.Sprintln("PushMove: The BitMoveSlice was full when trying to push a move"))
// //     }
//     moves = moves[:l+1]
//     moves[l] = move
//     *m = moves
// }

// func (m *BitMoveSlice) PopMove() BitMove {
//     moves := *m
//     l := len(moves)
// //     if 0 == l {
// //         panic(fmt.Sprintln("PopMove: The BitMoveSlice was empty when trying to pop a move"))
// //     }
//     move := moves[l-1]
//     moves = moves[:l-2]
//     *m = moves
//     return move
// }

// func (m *BitMoveSlice) LastMove() BitMove {
//     moves := *m
//     return moves[len(moves)-1]
// }

// func (m *BitMoveSlice) RemMove() {
//     moves := *m
//     l := len(moves)
// //     if 0 == l {
// //         panic(fmt.Sprintln("RemMove: The BitMoveSlice was empty when trying to remove a move"))
// //     }
//     moves = moves[:l-1]
//     *m = moves
// }

func (m *BitMoveSlice) Print() {
    moves := *m
    for _, move := range moves {
        fmt.Printf("(%d,%d) ", move.s, move.c)
    }
    fmt.Printf("\n")
}
