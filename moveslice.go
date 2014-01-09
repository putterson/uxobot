package main

import "fmt"

type MoveSlice []Move

func (m *MoveSlice) PushMove(move Move) {
    moves := *m
    l := len(moves)
//     if l == cap(moves) {
//         panic(fmt.Sprintln("PushMove: The MoveSlice was full when trying to push a move"))
//     }
    moves = moves[:l+1]
    moves[l] = move
    *m = moves
}

func (m *MoveSlice) PopMove() Move {
    moves := *m
    l := len(moves)
//     if 0 == l {
//         panic(fmt.Sprintln("PopMove: The MoveSlice was empty when trying to pop a move"))
//     }
    move := moves[l-1]
    moves = moves[:l-2]
    *m = moves
    return move
}

func (m *MoveSlice) LastMove() Move {
    moves := *m
    return moves[len(moves)-1]
}

func (m *MoveSlice) RemMove() {
    moves := *m
    l := len(moves)
//     if 0 == l {
//         panic(fmt.Sprintln("RemMove: The MoveSlice was empty when trying to remove a move"))
//     }
    moves = moves[:l-1]
    *m = moves
}

func (m *MoveSlice) Print() {
    moves := *m
    for _, move := range moves {
        fmt.Printf("(%d,%d) ", move.x, move.y)
    }
    fmt.Printf("\n")
}