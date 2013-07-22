#!/bin/python
# coding=utf-8

import sys

class Game:
    lastmove = ()
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
    
    def run(self):
        bot = XOHeuristic()
        s = GameState()
        display = Display()
        
        state = s.newBoard()
        display.printBoard(state)
        
        exits = ("exit", "quit", "q")
        
        players = ['X','O']
        
        while True:
            for curplayer in players:
                if curplayer == 'X':
                    move = raw_input("move=>")
                    if move in exits:
                        sys.exit(0)
                    if not self.handleMove(move, state, curplayer):
                        print "Please enter a valid move."
                        continue
                elif curplayer == 'O':
                    self.lastmove = bot.move(state, self.lastmove, curplayer)
                display.printBoard(state)
    
    def handleMove(self,move,state,player):
        if len(move) <= 1:
            return False
        pos = move.split(',')
        if len(self.lastmove) > 0:
            if len(pos) != 1 or not self.validMove(pos[0]):
                return False
            else:
                board = self.lastmove[1]
                cell = self.ymoves[pos[0][0]] * 3 + self.xmoves[pos[0][1]]
        else:
            if len(pos) != 2 or not self.validMove(pos[0]) or not self.validMove(pos[1]):
                return False
            board = self.ymoves[pos[0][0]] * 3 + self.xmoves[pos[0][1]]
            cell = self.ymoves[pos[1][0]] * 3 + self.xmoves[pos[1][1]]
        state[board][cell] = player
        self.lastmove = (board,cell)
        return True
    
    def validMove(self,pos):
        if len(pos) == 2:
            if pos[0] in self.ymoves and pos[1] in self.xmoves:
                return True
            else:
                return False
        else:
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
    
    def move(self, state, lastmove, player):
        
    
    def genChildren()
    def evalPos( self, board, player ):
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
    
    def evalState(self, node ):
        return 1000
    
    def alphabeta(self, node, depth, a, b, maximize):
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
    e = [ ' ' for i in range (9) ]
    
    def newBoard(self):
        e = self.e
        return [ list(e) for i in range(9) ]
    
    def move( state, player, board, pos ):
        state[board][pos] = player
    
    def sampleBoard(self):
        e = self.e
        
        x = [
            'X','O',' ',
            ' ','X',' ',
            'O',' ',' ']
        o = [
            'O','X',' ',
            ' ','O',' ',
            'X',' ',' ']

        aboard = [
            x,o,e,
            e,x,e,
            o,e,e
            ]

        return aboard

class Display:
    hthick = "━"
    vthick = "┃"
    cthick = "╋"

    hthin  = "─"
    vthin  = "│"
    cthin  = "┼"
    
    bline = [ vthick if i % 6 == 0 else vthin if i % 2 == 0 else ' ' for i in range(1,18) ]
    clthin = [ vthick if i % 6 == 0 else hthin if i % 2 == 1 else cthin for i in range(1,18) ]
    clthick = [ cthick if i % 6 == 0 else hthick for i in range(1,18) ]
    
    boardrow = bline + clthin + bline + clthin + bline
    
    board =  boardrow + clthick + boardrow + clthick + boardrow

    def printSubBoard(self, board, disp, start):
        return    
    
    def printBoard(self,state):
        board = list(self.board)
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
    
    def printTest(self):
        for i in range(len(self.board)):
            if i % 17 == 0:
                sys.stdout.write('\n')
            sys.stdout.write(self.board[i])
        sys.stdout.write('\n')

if __name__ == "__main__":
    game = Game()
    game.run()

