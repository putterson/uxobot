#!/bin/python
# coding=utf-8

# A bot for the ultimate tic-tac-toe game described here:
#   http://mathwithbaddrawings.com/2013/06/16/ultimate-tic-tac-toe/
# Utilizing alphabeta pruned minimax algorith at the moment
# Make sure that your terminal supports utf-8 line drawing characters
# 
# TODO:
#   add some way to store information on which boards have been won already
#   just realized that being sent to a won board means you play anywhere,
#   might not need to store won board info after all
#   
#
import sys
import time

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
    
    def run(self, mode):
        self.bot = XOHeuristic()
        s = GameState()
        display = Display()
        
        state = s.newBoard()
        display.printBoard(state)
        
        exits = ("exit", "quit", "q")
        
        players = ['X','O']
        
        while True:
            for curplayer in players:
                #if curplayer == 'X':
                    #move = raw_input("[move]>")
                    #if move in exits:
                        #sys.exit(0)
                    #if not self.handleMove(move, state, curplayer):
                        #print "Please enter a valid move."
                        #break
                #elif curplayer == 'O':
                    #print "Computer move."
                    #self.lastmove = self.bot.move(state, self.lastmove, curplayer)
                self.lastmove = self.bot.move(state, self.lastmove, curplayer)
                display.printBoard(state)
                #time.sleep(2)
                if self.gameWon(state, curplayer):
                    print "Player " + curplayer + " has won!"
                    print state
                    exit(0)
    
    def handleMove(self,move,state,player):
        lastmove = self.lastmove
        if len(move) <= 1:
            return False
        pos = move.split(',')
        
        if len(lastmove) > 0:
            if len(pos) != 2 or not self.validMove(pos[0]) or not self.validMove(pos[1]):
                return False
            else:
                for m in self.bot.genChildren( { 'state' : state, 'moves' : [lastmove]}, player ):
                    board = self.ymoves[pos[0][0]] * 3 + self.xmoves[pos[0][1]]
                    cell = self.ymoves[pos[1][0]] * 3 + self.xmoves[pos[1][1]]
                    if board == m[0] and cell == m[1]:
                        state[board][cell] = player
                        self.lastmove = (board,cell)
                        return True
        else:
            if len(pos) != 2 or not self.validMove(pos[0]) or not self.validMove(pos[1]):
                return False
            else:
                for m in self.bot.genChildren( { 'state' : state, 'moves' : [] }, player ):
                    board = self.ymoves[pos[0][0]] * 3 + self.xmoves[pos[0][1]]
                    cell = self.ymoves[pos[1][0]] * 3 + self.xmoves[pos[1][1]]
                    if board == m[0] and cell == m[1]:
                        state[board][cell] = player
                        self.lastmove = (board,cell)
                        return True
        return False
    
    def validMove(self,pos):
        if len(pos) == 2:
            if pos[0] in self.ymoves and pos[1] in self.xmoves:
                return True
            else:
                return False
        else:
            return False

    def gameWon(self, state, player):
        if len(self.bot.genChildren({ 'state' : state, 'moves' : [self.lastmove]}, self.bot.notPlayer(player)))==0:
            return True

class XOHeuristic:
    
    depth = 5
    colour = 'X'
    lastmove = ()
    
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
    
    
    alpha = -18000
    beta  =  18000
    
    def move(self, state, lastmove, player):
        self.colour = player
        moves = []
        moves.append(lastmove)
        # call alphabeta with depth, it will return the end state and the list of moves leading to it
        print moves
        node = { 'state' : state, 'moves' : moves}
        result = self.alphabeta( node, self.depth, self.alpha, self.beta, player)
        move = result[1]
        state[move[0]][move[1]] = player
        print move
        return move
    
    def genChildren(self, node, player):
        if node['moves']:
            lastmove = node['moves'].pop()
        else:
            lastmove = ()
            node['moves'] = []
        state = node['state']
        moves = []
        if len(lastmove) == 0 or abs(self.evalBoard(state[lastmove[1]], player)) >= 1000:
            for i,b in enumerate(state):
                if abs(self.evalBoard(b, player)) >= 1000:
                    continue
                for j,c in enumerate(state[i]):
                    if c == ' ':
                        moves.append((i,j))
        else:
            i = lastmove[1]
            for j,c in enumerate(state[i]):
                if c == ' ':
                    moves.append((i,j))
            if len(moves) == 0:
                for i,b in enumerate(state):
                    if abs(self.evalBoard(b, player)) >= 1000:
                        for j,c in enumerate(state[i]):
                            if c == ' ':
                                moves.append((i,j))
        node['moves'].append(lastmove)
        return moves
    
    # max score is 1000/-1000
    # if the score = 1000 then the board is won for player
    def evalBoard(self, board, player):
        t = 0
        opponent = self.notPlayer(self.colour)
        for i in range(8):
            players = 0
            others = 0
            for j in range(3):
                piece = board[self.wins[i][j]]
                if piece == self.colour:
                    players += 1
                elif piece == opponent:
                    others += 1
            t += self.score[players][others]
            
            #if the board is won override and return the max/min value
            if abs(self.score[players][others]) >= 1000:
                t = self.score[players][others]
                break
        #cap the max/min value to 1000
        if t > 1000:
            t = 1000
        elif t < -1000:
            t = -1000            
        return t
    
    # compute score for whole metaboard
    def evalState(self, node, maximize):
        #scores = [self.evalBoard(node['state'][i],maximize) for i in range(9)]
        #print scores
        return sum([self.evalBoard(node['state'][i],maximize) for i in range(9)])
    
    # 
    # 
    def alphabeta(self, node, depth, a, b, maximize):
        children = self.genChildren(node, maximize)
        move = []
        if depth == 0 or len(children) == 0:
            return (self.evalState(node, maximize), node['moves'][-1])
        if maximize == self.colour:
            for child in children:
                b = child[0]
                c = child[1]
                
                node['moves'].append(child)
                node['state'][b][c] = maximize
                
                comp = self.alphabeta(node, depth - 1, a, b, self.notPlayer(maximize))
                
                node['state'][b][c] = ' '
                node['moves'].pop()
                
                if comp[0] >= a:
                    a = comp[0]
                    move = child
                if b <= a:
                    break #beta cutoff
            return (a, move)
        else:
            for child in children:
                node['moves'].append(child)
                node['state'][child[0]][child[1]] = maximize
                
                comp = self.alphabeta(node, depth - 1, a, b, self.notPlayer(maximize))
                
                node['state'][child[0]][child[1]] = ' '
                node['moves'].pop()
                
                if comp[0] <= b:
                    b = comp[0]
                    move = child
                if b <= a:
                    break #alpha cutoff
            return (b, child)
    
    def notPlayer(self,p):
        if p == 'X':
            return 'O'
        else:
            return 'X'

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
    mode = ''
    game.run(mode)

