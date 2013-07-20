#!/bin/python
# coding=utf-8

import sys

class Game:

	lastmove = ""
	xmoves = {
		"l" : 0,
		"m" : 1,
		"r" : 2
		}

	ymoves = {
		"t" : 0,
		"m" : 1,
		"b" : 2
		}

	def run(this):
		bot = XOHeuristic()
		s = GameState()
		display = Display()

		state = s.sampleBoard()
		display.printBoard(state)

		while True:
			move = raw_input("move=>")
			if move == "exit":
				sys.exit(0)
			if not this.handleMove(move, state):
				print "Please enter a valid move."
				continue
			display.printBoard(state)
	
	def handleMove(this,move,state):
		if len(move) <= 1:
			return False
		pos = move.split(',')
		
		print pos
		return False
		
		
		


class XOHeuristic:
	wins = (
		(0,1,2),
		(3,4,5),
		(6,7,8),
		(0,3,6),
		(1,4,7),
		(2,5,8),
		(0,4,8),
		(2,4,6)
		)
	score = (
		(0,   -10,  -100, -1000),
		(10,    0,     0,     0),
		(100,   0,     0,     0),
		(1000,  0,     0,     0),
		)
	
	def evalPos( board, player ):
		if player == 'X':
			opponent = 'O'
		else :
			opponent = 'X'
		for i in range(8):
			players = 0
			others = 0
			for j in range(3):
				piece = board[wins[i][j]]
				if piece == player:
					players += 1
				elif piece == opponent:
					others += 1
			t += score[players][others]
		return t
	
	def evalState( node ):
		return 1000
	
	def alphabeta(node, depth, a, b, maximize):
		if depth == 0 or len(node.children) == 0:
			return evalState(node)
		if maximize:
			for child in node.children:
				a = max(a, alphabeta(child, depth - 1, a, b, not maximize))
				if b <= a:
					break #beta cutoff
			return a
		else:
			for child in node.children:
				b = min(b, alphabeta(child, depth - 1, a, b, not maximize))
				if b <= a:
					break #alpha cutoff
			return b

class GameState:
	e = (' ',' ',' ',
		' ',' ',' ',
		' ',' ',' ')
	
	def newBoard(this):
		e = this.e
		return (
			e,e,e,
			e,e,e,
			e,e,e
		)
	
	def move( state, player, board, pos ):
		state[board][pos] = player
	
	def sampleBoard(this):
		e = this.e
		
		x = (
			'X','O',' ',
			' ','X',' ',
			'O',' ',' ')
		o = (
			'O','X',' ',
			' ','O',' ',
			'X',' ',' ')

		aboard = (
			x,o,e,
			e,x,e,
			o,e,e
			)

		return aboard

class Display:
	hthick = "━"
	vthick = "┃"
	
	hthin  = "─"
	vthin  = "│"
	
	cthick = "╋"
	cthin  = "┼"
	
	bline = (
		' ',vthin,' ',vthin,' ',vthick,
		' ',vthin,' ',vthin,' ',vthick,
		' ',vthin,' ',vthin,' ')
	cthin = (
		hthin,cthin,hthin,cthin,hthin,vthick,
		hthin,cthin,hthin,cthin,hthin,vthick,
		hthin,cthin,hthin,cthin,hthin)
	cthick = (
		hthick,hthick,hthick,hthick,hthick,cthick,
		hthick,hthick,hthick,hthick,hthick,cthick,
		hthick,hthick,hthick,hthick,hthick)
	
	boardrow = bline + cthin + bline + cthin + bline
	
	board =  boardrow + cthick + boardrow + cthick + boardrow

	def printSubBoard(this, board, disp, start):
		return	
	
	def printBoard(this,state):
		board = list(this.board)
		game = state
		for b in range(len(game)):
			for c in range(len(game[b])):
				b_location = divmod(b, 3)
				by = b_location[0] * 17 * 6
				bx = b_location[1] * 6
				
				c_location = divmod(c, 3)
				cy = c_location[0] * 17 * 2
				cx = c_location[1] * 2
				
				board[by+bx+cy+cx] = game[b][c]
		for i in range(len(board)):
			if i % 17 == 0:
				sys.stdout.write('\n')
			sys.stdout.write(board[i])
		sys.stdout.write('\n')
	
	def printTest(this):
		for i in range(len(this.board)):
			if i % 17 == 0:
				sys.stdout.write('\n')
			sys.stdout.write(this.board[i])
		sys.stdout.write('\n')

if __name__ == "__main__":
	game = Game()
	game.run()

