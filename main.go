package main

import (
	"fmt"
	"bufio"
	"os"
// 	"os/signal"
	"strings"
	"flag"
	//"errors"
	"time"
)

var Stdin = bufio.NewReader(os.Stdin)


var maptochar = []string{" ", "X", "O"}


type GameSettings struct {
	players [2]string
	depth [2]float64
	curplayer Player
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


	
	

	settings := &GameSettings{
		[2]string{"montecarlo","negamax",},
		[2]float64{1,9,},
		X,
		false,
		false,
		NewMove(),
		new(Board),
	}

	prompt := flag.Bool("prompt", false, "Don't start game immediately.")
	ai_one := flag.Float64("s1", 2, "Set AI 1 strength.")
	ai_two := flag.Float64("s2", 2, "Set AI 2 strength.")

	player_one := flag.String("p1", "montecarlo", "Set player 1 to cpu or human.")
	player_two := flag.String("p2", "negamax", "Set player 2 to cpu or human.")

	flag.Parse()

	if !(*prompt) {
		settings.players[0] = *player_one
		settings.players[1] = *player_two
		settings.depth[0] = *ai_one
		settings.depth[1] = *ai_two
		settings.board = newBoard()
		settings.curplayer = X
		gameloop(settings)
		os.Exit(0)
	}

	for {

		fmt.Printf(colourString("[uxobot]>"))
		var input, _ = Stdin.ReadString('\n')

		if strings.Contains(input, "help"){
			fmt.Println("Available commands are: start, quit")
		} else if strings.Contains(input, "quit") || strings.Contains(input, "exit") {
			os.Exit(0)
		} else if len(input) < 2 || strings.Contains(input, "start") && !settings.running {
			settings.board = newBoard()
			settings.curplayer = X
			gameloop(settings)
		} else {
			fmt.Println("Enter a valid command or type help.")
		}
	}
}

func newBoard() *Board {
	b := new(Board)
	for x := range b {
		for y := range b[x] {
			b[x][y] = B
		}
	}
	return b
}

func makeBot(bottype string, depth float64) UXOBot {
	if bottype == "negamax" {
		negabot := new(NegaMax)
		negabot.setDepth(int(depth))
		return negabot
	} else if bottype == "montecarlo" {
		return NewMonteCarlo(depth)
	} else {
		return nil
	}
}

func gameloop(s *GameSettings){
	lastmove := *NewMove()
	var move Move
	bots := [2]UXOBot{
		makeBot(s.players[0], s.depth[0]),
		makeBot(s.players[1], s.depth[1]),
	}
	
	for {

		if bots[s.curplayer-1] != nil {
			move = getCpuMove(bots[s.curplayer-1], s.board, &lastmove, s.curplayer)
		} else {
			if len(genHumanChildren(s.board, lastmove)) > 0 {
				move = getHumanMove(s.board, lastmove)
			} else {
				move = *NewMove();
			}
		}



		if move.isNoMove() {
			drawBoard(s.board, move)
			fmt.Println("Draw! Game Over!")
			return
		}

		for _, bot := range(bots){
			if bot != nil {
				err := bot.makeMove(move)
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		(*s.board)[move.x][move.y] = s.curplayer
		getSuperScores(s.board).Print()
		drawBoard(s.board, move)

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


func getHumanMove(board *Board, lastmove Move) (move Move) {
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
				fmt.Printf("Score of board %d is: %d\n",
					i,
					evalSubBoard(board, move.x, move.y))
			}
		} else if strings.Contains(input, "evala") {
			fmt.Printf("Score of whole board is: %d\n", evalBoard(board))
		} else {
			fmt.Println("Enter a valid command or type help.")
		}
	}
}

func getCpuMove(bot UXOBot, b *Board, lastmove *Move, player Player) Move {


	
	//node.moves.Print()
	start_t := time.Now()
	
	move, err := bot.getMove(*b, *lastmove, player)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	
	duration := time.Since(start_t).Seconds()


	fmt.Printf("Move: [%d,%d] Time (s): %.2f\n",
		move.x, move.y,	duration)
	return move
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
