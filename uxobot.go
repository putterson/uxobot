package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	//"errors"
)

var Stdin = bufio.NewReader(os.Stdin)

// Constants for the piece values (blank, X, O)
const (
	B byte = 0
	X byte = 1
	O byte = 2
)

// Constants for the max and minimum board scores
const (
	SCOREMAX = 10000
	SCOREMIN = -10000
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
type Board [9][9]byte
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
	for x := range b {
		for y := range b[x] {
			b[x][y] = B
		}
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
	(*b)[move.x][move.y] = curplayer
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

func genHumanChildren(b *Board, lastmove Move) MoveSlice {
	return genChildren(b, &lastmove, new(Won))
}

func genChildren(b *Board, lastmove *Move, mboard *Won) MoveSlice {
	var moves MoveSlice
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
func genAllChildren(b *Board, lastmove *Move) MoveSlice {
	moves := MoveSlice{}
	var x,y byte
	for x = 0; x < 9; x++ {
		for y = 0; y < 9; y++ {
			if b[x][y] == B {
				moves = append(moves, Move{x: x, y: y})
			}
		}
	}
	return moves
}

// Generate children for a specific board
func genBoardChildren(b *Board, lastmove *Move) MoveSlice {
	moves := MoveSlice{}
	ox := (lastmove.x%3) * 3
	oy := (lastmove.y%3) * 3
	//fmt.Printf("Origin x:%d y:%d\n", ox, oy)
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

/*
*
* Board drawing functionality
*
*/
func drawBoard(b *Board, move Move) {
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			fmt.Printf("%c ", maptochar[(*b)[x][y]])
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

func playerToMul(player byte) int {
	if player == X {
		return 1
	} else {
		return -1
	}
}

func getCpuMove(b *Board, lastmove *Move, player byte) Move {
	//TODO: replace depth with config
	depth := 4
	moves := make(MoveSlice, 0, depth + 1)

	node := new(AINode)
	node.board = b
	node.moves = &moves

	(*node.moves).PushMove(*lastmove)
	//node.moves.Print()
	_, move, err := negamax(node, depth, -10000, 10000, player, true)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return move
}

//TODO: evaluate the score of the board
func evalBoard(board *Board, won *Won) int {
	return 0
}

// evaluate the score of a sub-board
func evalSubBoard(board *Board, won *Won) int {
	return 0
}

// note: we return CacheEntry because it has all the information we need to return to a higher level of negamax
func negamax(node *AINode, depth int, alpha int, beta int, player byte, first bool) (CacheEntry, Move, error) {
	//fmt.Printf("depth: %d\n", depth)
	//fmt.Printf("Movecount: %d ", len(*node.moves))
	//node.moves.Print()
	lastmove := node.moves.LastMove()
	//drawBoard(node.board, lastmove)
	children := genChildren(node.board, &lastmove, node.won)
	if depth == 0 || len(children) == 0 {
		//fmt.Printf("FIN depth %d children %d\n",depth, len(children))
		return CacheEntry{
			//FIXME: CHANGE THIS to be correct
			cutoff : alpha,
			depth  : depth,
			score  : playerToMul(player) * evalBoard(node.board, node.won),
			move   : lastmove,
		}, Move{}, nil
	}

	maxScore := SCOREMIN
	var maxEntry CacheEntry

	//defer func(){
	//	if r := recover(); r != nil {
	//		fmt.Printf("Depth: %d\n", depth)
	//	}
	//}()

	var bestChild Move
	for _, child := range children {
		//fmt.Printf("d%d Trying child %d (%d,%d)\n",depth, i, child.x, child.y)

		// PushMove will panic if it fails (shouldn't fail)
		node.moves.PushMove(child)
		node.board[child.x][child.y] = player

		// NOTE: alpha and beta are negated and swapped for the subcall to negamax
		entry, _, err := negamax(node, depth-1, -beta, -alpha, notPlayer(player), false)
		if err != nil {
			return *new(CacheEntry), Move{}, err
		}

		node.board[child.x][child.y] = B
		node.moves.RemMove()

		if -entry.score > maxScore {
			maxScore = -entry.score
			maxEntry = entry
			if first {
				bestChild = child
			}
		}
		alpha = max(alpha, entry.score)
		if alpha >= beta {
			break
		}
	}
	return maxEntry, bestChild, nil
}

func (m *MoveSlice) PushMove(move Move) {
	moves := *m
	l := len(moves)
	if l == cap(moves) {
		panic(fmt.Sprintln("PushMove: The MoveSlice was full when trying to push a move"))
	}
	moves = moves[:l+1]
	moves[l] = move
	*m = moves
}

func (m *MoveSlice) PopMove() Move {
	moves := *m
	l := len(moves)
	if 0 == l {
		panic(fmt.Sprintln("PopMove: The MoveSlice was empty when trying to pop a move"))
	}
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
	if 0 == l {
		panic(fmt.Sprintln("RemMove: The MoveSlice was empty when trying to remove a move"))
	}
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

func max(a int, b int) int{
	if a > b{
		return a
	}
	return b
}
