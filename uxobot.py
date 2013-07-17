#!/bin/python

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

	evalPos( board, player ):
		if player == 'X':
			opponent = 'O'
		else :
			opponent = 'X'
		opponent = (player == 'X') ? 'O' : 'X'
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
	sampleBoard():
		return 
