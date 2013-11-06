package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"errors"
)

var Stdin = bufio.NewReader(os.Stdin)

// Constants for the piece values (blank, X, O)
const (
	B byte = 0
	X byte = 1
	O byte = 2
)

var maptochar = []rune{' ', 'X', 'O'}
var playerconf = []string{"human", "cpu"}
var wins = [][]byte{
	{1, 1, 1, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 1, 1, 1, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 1, 1, 1},
	{1, 0, 0, 0, 1, 0, 0, 0, 1},
	{0, 0, 1, 0, 1, 0, 1, 0, 0},
	{1, 0, 0, 1, 0, 0, 1, 0, 0},
	{0, 1, 0, 0, 1, 0, 0, 1, 0},
	{0, 0, 1, 0, 0, 1, 0, 0, 1},
}

func notPlayer(p byte) byte {
	if p == X {
		return O
	} else {
		return X
	}
}

// The board is stored as a flat array, interpreted as a 9*9 array using modular arithmancy
type Board [81]byte
type Won [9]bool

// a Move is the x,y coordinates of the move, 0 based and starting at the top left of the board
type Move struct {
	x byte
	y byte
}
type MoveSlice []Move

//AI types and structs
type BHash uint64
type CacheEntry struct {
	cutoff	int
	depth	int
	score	int
	move	Move
}
type AINode struct {
	board	*Board
	moves	*MoveSlice
	cache	map[BHash]CacheEntry
	won	*Won
	hash	BHash
}


const NoMove byte = 10
func NewMove() *Move {
	return &Move{
		x: NoMove,
		y: NoMove,
	}
}

/*
*
* Main loops are here
*
*/
func main() {
	run()
}

func run() {
	b := new(Board)
	for i := range b {
		b[i] = B
	}
	var curplayer = X
	lastmove := *NewMove()
	var move Move
	for {
		move = makeMove(b, curplayer, lastmove)
		drawBoard(b, move)
		lastmove = move
		curplayer = notPlayer(curplayer)
	}
}

/*
*
* Utility functions for the different move notations
*
*/

func move_to_pos(move Move) byte {
	return move.y*9 + move.x
}

func xy_to_pos(x, y byte) byte {
	//fmt.Printf("%d", x*9+y)
	return x*9 + y
}

func pos_to_move(pos byte) Move {
	return Move{
		x: pos % 9,
		y: pos / 9,
	}
}

func move_to_board(m Move) Move {
	return Move{
		x: (m.x%3)*3,
		y: (m.y%3)*3,
	}
}

// Takes human move notation ie. board#,cell# and returns a Move
func move_notation(b, c byte) Move {
	b = b-1
	c = c-1
	return Move{
		x: (b%3)*3 + (c%3),
		y: (b/3)*3 + (c/3),
	}
}

/*
*
* Move making and validation functions
*
*/
func makeMove(b *Board, curplayer byte, lastmove Move) (move Move) {
	if playerconf[curplayer-1] == "human" {
		move = getMove(b, lastmove)
	} else {
		move = getCpuMove(b, &lastmove, curplayer)
	}
	(*b)[move_to_pos(move)] = curplayer
	return
}

func getMove(board *Board, lastmove Move) (move Move) {
	move = *(new(Move))
	var b, c byte
	for {
		fmt.Printf("[uxobot]>")
		var input, _ = Stdin.ReadString('\n')
		var _, err = fmt.Sscanln(input, &b, &c)
		if err == nil {
			if b > 0 && b < 10 && c > 0 && c < 10 {
				move = move_notation(b, c)
				//fmt.Printf("x:%d y:%d\n",move.x,move.y)
				for _, vm := range genHumanChildren(board, lastmove) {
					if move == vm {
						return
					}
				}
			}
			// if there were no matching moves
			fmt.Println("Please make a valid move, or type help.")
		} else if strings.Contains(input, "help"){
			fmt.Println("There should be some help here...")
			fmt.Println("Valid moves are of the form \"1 9\" ie. a board number followed by a cell number (1-9)")
		} else if strings.Contains(input, "quit") || strings.Contains(input, "exit") {
			os.Exit(0)
		} else {
			fmt.Println("Enter a valid command or type help.")
		}
	}
}

func genHumanChildren(b *Board, lastmove Move) []Move {
	return genChildren(b, &lastmove, new(Won))
}

func genChildren(b *Board, lastmove *Move, mboard *Won) []Move {
	var moves []Move
	if (lastmove.x == NoMove) && (lastmove.y == NoMove) {
		return genAllChildren(b, lastmove)
	}

	moves = genBoardChildren(b, lastmove)
	if len(moves) == 0 {
		return genAllChildren(b, lastmove)
	}
	return moves
}

// Generate children for all the boards
func genAllChildren(b *Board, lastmove *Move) []Move {
	moves := []Move{}
	var pos byte
	for pos = 0; pos < 81; pos++ {
		if b[pos] == B {
			moves = append(moves, pos_to_move(pos))
		}
	}
	return moves
}

// Generate children for a specific board
func genBoardChildren(b *Board, lastmove *Move) []Move {
	moves := []Move{}
	ox := (lastmove.x%3) * 3
	oy := (lastmove.y%3) * 3
	//fmt.Printf("Origin x:%d y:%d\n", ox, oy)
	for x := ox; x < ox+3; x++ {
		for y := oy; y < oy+3; y++ {
			if b[xy_to_pos(x,y)] == B {
				move := Move{x: x, y: y}
				moves = append(moves, move)
			}
		}
	}

	return moves
}

/*
*
* Board drawing functionality
*
*/
func drawBoard(b *Board, move Move) {
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			fmt.Printf("%c ", maptochar[(*b)[y*9+x]])
		}
		fmt.Printf("\n\n")
	}
	return
}


/*
*
* AI functionality
*
*/

func getCpuMove(b *Board, lastmove *Move, player byte) Move {
	//TODO: replace depth with config
	depth := 6
	moves := make(MoveSlice, depth + 1)

	node := new(AINode)
	node.board = b
	node.moves = &moves

	(*node.moves).PushMove(*lastmove)

	entry, err := negamax(node, depth, -10000, 10000, player)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return entry.move
}

//TODO: evaluate the score of the board
func evalBoard(board *Board) int {
	return 0
}

// note: we return CacheEntry because it has all the information we need to return to a higher level of negamax
func negamax(node *AINode, depth int, alpha int, beta int, player byte) (CacheEntry, error) {
	lastmove, err := node.moves.PopMove()
	if err != nil {
		return *new(CacheEntry), err
	}

	children := genChildren(node.board, &lastmove, node.won)
	
	if depth == 0 || len(children) == 0 {
		return CacheEntry{
			//FIXME: CHANGE THIS to be correct
			cutoff : alpha,
			depth  : 0,
			score  : evalBoard(node.board),
			move   : lastmove,
		}, nil
	}
	return *new(CacheEntry), nil
}

func (m *MoveSlice) PushMove(move Move) error {
	moves := *m
	l := len(moves)
	if l == cap(moves) {
		return errors.New("PushMove: the MoveSlice has reached it's maximum size")
	}
	moves[l+1] = move
	*m = moves
	return nil
}

func (m *MoveSlice) PopMove() (Move, error) {
	moves := *m
	l := len(moves)
	if 0 == l {
		return *NewMove(), errors.New("PopMove: the MoveSlice is empty")
	}
	move := moves[l-1]
	moves = moves[:l-1]
	return move, nil
	*m = moves
	return move, nil
}
