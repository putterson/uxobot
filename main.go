package main

import (
	"fmt"
	"bufio"
	"os"
// 	"os/signal"
	"strings"
	//"errors"
)

var Stdin = bufio.NewReader(os.Stdin)

// Constants for the piece values (blank, X, O)
const (
	B int = 0
	X int = 1
	O int = 2
)
var maptochar = []string{" ", "X", "O"}
func notPlayer(p int) int {
	if p == X {
		return O
	} else {
		return X
	}
}

// a Move is the x,y coordinates of the move, 0 based and starting at the top left of the board
type Move struct {
	x int
	y int
}
const NoMove int = 254

func NewMove() *Move {
	return &Move{
		x: NoMove,
		y: NoMove,
	}
}

// The board is stored as a 2D array
type Board [9][9]int
type Scores [3][3]int


type GameSettings struct {
	players [2]string
	depth [2]int
	curplayer int
	running bool
	paused bool
	lastmove *Move
	board *Board
}

/*
*
* Main loops
*
*/
func main() {
// 	sigs := make(chan os.Signal, 1)
// 	done := make(chan bool, 1)


	gamerunning := false
	var b *Board

	settings := &GameSettings{
		[2]string{"cpu","cpu",},
		[2]int{5,8,},
		X,
		false,
		false,
		NewMove(),
		new(Board),
	}

	fmt.Println("Initializing Zobrist keys...");
	init_zobrist_keys()
	for {

		fmt.Printf(colourString("[uxobot]>"))
		var input, _ = Stdin.ReadString('\n')

		if strings.Contains(input, "help"){
			fmt.Println("There should be some help here...")
			fmt.Println("Available commands are: start, quit, evala, evalb")
		} else if strings.Contains(input, "quit") || strings.Contains(input, "exit") {
			os.Exit(0)
		} else if strings.Contains(input, "evalb") && gamerunning {
			var i int
			for i = 1; i < 10; i++ {
				move := move_to_subboard(move_notation(i, i))
				fmt.Printf("Score of board %d is: %d\n", i, evalSubBoard(b, move.x, move.y))
			}
		} else if strings.Contains(input, "evala") && settings.running {
			fmt.Printf("Score of whole board is: %d\n", evalBoard(b))
		} else if strings.Contains(input, "clearcache") {
			ai_cache = make(AICache)
		} else if len(input) < 2 || strings.Contains(input, "start") && !settings.running {
			b = new(Board)
			for x := range b {
				for y := range b[x] {
					b[x][y] = B
				}
			}

			settings.board = b
			settings.curplayer = X
			gameloop(settings)
		} else if strings.Contains(input, "up") {
			settings.depth[0]++
			settings.depth[1]++
			fmt.Printf("Difficulties: X: %d -- O: %d\n", settings.depth[0], settings.depth[1]);
		} else if strings.Contains(input, "down") {
			settings.depth[0]--
			settings.depth[1]--
			fmt.Printf("Difficulties: X: %d -- O: %d\n", settings.depth[0], settings.depth[1]);
		} else if strings.Contains(input, "stats") {
			fmt.Printf("Difficulties: X: %d -- O: %d\n", settings.depth[0], settings.depth[1]);
			fmt.Printf("Cache size: %d", len(ai_cache))
		} else {
			fmt.Println("Enter a valid command or type help.")
		}
	}
}


func gameloop(s *GameSettings){
	lastmove := *NewMove()
	var move Move
	for {

		if s.players[s.curplayer-1] == "human" {
			if len(genHumanChildren(s.board, lastmove)) > 0 {
				move = getMove(s.board, lastmove)
			} else {
				move = *NewMove();
			}
		} else {
			move = getCpuMove(s.board, &lastmove, s.curplayer, s.depth[s.curplayer-1])
		}

		if (move.x == NoMove) && (move.y == NoMove) {
			drawBoard(s.board, move)
			fmt.Println("Game Over!")
			fmt.Printf("Hash: %x\n", hash_board(s.board))
			return
		}

		(*s.board)[move.x][move.y] = s.curplayer
		drawBoard(s.board, move)
		fmt.Printf("Hash: %x\n", hash_board(s.board))



		if evalBoard(s.board) == SCOREMAX {
			fmt.Println("X Wins the game!")
			return
		} else if evalBoard(s.board) == SCOREMIN {
			fmt.Println("O Wins the game!")
			return
		}

		lastmove = move
		s.curplayer = notPlayer(s.curplayer)
	}
}

/*
*
* Utility functions
*
*/

// func move_to_pos(move Move) int {
// 	return move.y*9 + move.x
// }
// 
// func xy_to_pos(x, y int) int {
// 	//fmt.Printf("%d", x*9+y)
// 	return x*9 + y
// }
// 
// func pos_to_move(pos int) Move {
// 	return Move{
// 		x: pos % 9,
// 		y: pos / 9,
// 	}
// }

// Takes a move and returns the top left corner of the board that the move was in
func move_to_subboard(m Move) Move {
	return Move{
		x: (m.x%3)*3,
		y: (m.y%3)*3,
	}
}

// Takes human move notation ie. board#,cell# and returns a Move
func move_notation(b, c int) Move {
	b = b-1
	c = c-1
	return Move{
		x: (b%3)*3 + (c%3),
		y: (b/3)*3 + (c/3),
	}
}

func max(a int, b int) int{
	if a > b{
		return a
	}
	return b
}

func min(a int, b int) int{
	if a < b{
		return a
	}
	return b
}

func getMove(board *Board, lastmove Move) (move Move) {
	move = *(new(Move))
	var b, c int
	for {
		fmt.Printf(colourString("[move]>"))
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
		} else if strings.Contains(input, "evalb") {
			var i int
			for i = 1; i < 10; i++ {
				move = move_to_subboard(move_notation(i, i))
				fmt.Printf("Score of board %d is: %d\n", i, evalSubBoard(board, move.x, move.y))
			}
		} else if strings.Contains(input, "evala") {
			fmt.Printf("Score of whole board is: %d\n", evalBoard(board))
		} else {
			fmt.Println("Enter a valid command or type help.")
		}
	}
}

func getCpuMove(b *Board, lastmove *Move, player int, depth int) Move {
	moves := make(MoveSlice, 0, depth + 1)
	ai_cache = make(AICache)
	node := new(AINode)
	node.board = b
	node.moves = &moves
	node.hashes = hash_board(b)

	(*node.moves).PushMove(*lastmove)
	//node.moves.Print()
	entry, move, err := negamax(node, depth, SCOREMIN - 1, SCOREMAX + 1, player, true)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("Move: [%d,%d] Score: %d Cache_s: %d\n", move.x, move.y, entry.score, len(ai_cache))
	return move
}

func genHumanChildren(b *Board, lastmove Move) MoveSlice {
	return genChildren(b, &lastmove, new(Scores))
}

func genChildren(b *Board, lastmove *Move, scores *Scores) MoveSlice {
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

/*
*
* Board drawing functionality
*
*/
func drawBoard(b *Board, move Move) {
	thins  := "─┼─┼─┃─┼─┼─┃─┼─┼─"
	thicks := "━━━━━╋━━━━━╋━━━━━"
	piece := ""

	var x, y int
	for y = 0; y < 9; y++ {
		if y > 0 {
			if y % 3 == 0 {
				fmt.Println(thicks)
			} else {
				fmt.Println(thins)
			}
		}

		for x = 0; x < 9; x++ {
			if move.x == x && move.y == y {
				piece = colourString(maptochar[(*b)[x][y]])
			} else {
				piece = maptochar[(*b)[x][y]]
			}


			if x > 0{
				if x % 3 == 0{
					fmt.Printf("┃%s", piece)
				} else {
					fmt.Printf("│%s", piece)
				}
			} else {
				fmt.Printf("%s", piece)
			}
		}
		fmt.Printf("\n")
	}
	return
}

func colourString(c string) string{
//	return c
	return "\x1b[31;1m" + c + "\x1b[0m"
}


