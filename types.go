package main



type MoveScore struct {
	move Move
	score int
}


type SubScores [3][3]int
type BitSubScores [9]int
type Scores [3][3]float64


//Definitions usefull for bots and evaluating the board
type Line struct {
	x1 int
	x2 int
	x3 int
	y1 int
	y2 int
	y3 int
}


var winlines = []Line{
	Line{x1: 0, x2: 0, x3: 0, y1: 0, y2: 1, y3: 2},
	Line{x1: 1, x2: 1, x3: 1, y1: 0, y2: 1, y3: 2},
	Line{x1: 2, x2: 2, x3: 2, y1: 0, y2: 1, y3: 2},

	Line{x1: 0, x2: 1, x3: 2, y1: 0, y2: 0, y3: 0},
	Line{x1: 0, x2: 1, x3: 2, y1: 1, y2: 1, y3: 1},
	Line{x1: 0, x2: 1, x3: 2, y1: 2, y2: 2, y3: 2},

	Line{x1: 0, x2: 1, x3: 2, y1: 0, y2: 1, y3: 2},
	Line{x1: 0, x2: 1, x3: 2, y1: 2, y2: 1, y3: 0},
}

