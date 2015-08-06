package main

// Constants for the piece values (blank, X, O)
type Player int

const (
	B Player = 0
	X Player = 1
	O Player = 2
)

func notPlayer(p Player) Player {
	if p == X {
		return O
	} else {
		return X
	}
}
