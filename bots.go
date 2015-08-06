package main

type UXOBot interface {
	getMove(board *Board, lastmove *Move, player Player) (Move, error)
}
