#!/bin/python
# coding=utf-8

# A bot for the ultimate tic-tac-toe game described here:
#   http://mathwithbaddrawings.com/2013/06/16/ultimate-tic-tac-toe/
# Utilizing alphabeta pruned negamax algorithm with transposition tables
# Make sure that your terminal supports utf-8 line drawing characters
# 
#   
#
import sys
import time
import random

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
        self.botX = XOHeuristic()
        self.botO = XOHeuristic()
        s = GameState()
        display = Display()
        
        state = s.newBoard()
        display.printBoard(state)
        
        exits = ("exit", "quit", "q")
        
        players = ['X','O']
        
        count = 0
        
        while True:
            for curplayer in players:
                count += 1
                if curplayer == 'O':
                    while True:
                        move = raw_input("[move]>")
                        if move in exits:
                            sys.exit(0)
                        if self.handleNMove(move, state, curplayer):
                            break
                        else:
                            print "Please enter a valid move."
                elif curplayer == 'X':
                    print "Computer move."
                    self.lastmove = self.botX.move(state, self.lastmove, curplayer)
                    
                #if curplayer == 'X':
                    #self.lastmove = self.botX.move(state, self.lastmove, curplayer)
                #else:
                    #self.lastmove = self.botO.move(state, self.lastmove, curplayer)
                
                
                display.printBoard(state)
                #time.sleep(2)
                win = self.gameWon(state, curplayer)
                if win == 0:
                    continue
                elif win == 1:
                    print "Player " + curplayer + " has won in " + str(count) + " moves!"
                    exit(0)
                elif win == 2:
                    print "Player " + curplayer + " has drawn the game in " + str(count) + " moves!"
                    exit(0)
    
    def handleNMove(self, move, state, player):
        lastmove = self.lastmove
        if len(move) != 2 or not self.validNMove(move[0]) or not self.validNMove(move[1]):
            return False
        board = int(move[0]) - 1
        cell = int(move[1]) - 1
        children = self.botX.genChildren( { 'state' : state, 'moves' : [lastmove]}, player )
        #print children
        print lastmove
        for m in children:
            if board == m[0] and cell == m[1]:
                state[board][cell] = player
                self.lastmove = (board,cell)
                return True     
        return False
    
    def validNMove(self, num):
        try:
            n = int(num)
        except ValueError:
                return False
        if n >= 1 and n <= 9:
            return True
        else:
            return False
    
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
                    cell  = self.ymoves[pos[1][0]] * 3 + self.xmoves[pos[1][1]]
                    if board == m[0] and cell == m[1]:
                        state[board][cell] = player
                        self.lastmove = (board,cell)
                        return True
                else:
                    return False
        else:
            if len(pos) != 2 or not self.validMove(pos[0]) or not self.validMove(pos[1]):
                return False
            else:
                for m in self.bot.genChildren( { 'state' : state, 'moves' : [] }, player ):
                    board = self.ymoves[pos[0][0]] * 3 + self.xmoves[pos[0][1]]
                    cell  = self.ymoves[pos[1][0]] * 3 + self.xmoves[pos[1][1]]
                    if board == m[0] and cell == m[1]:
                        state[board][cell] = player
                        self.lastmove = (board,cell)
                        return True
                else:
                    return False
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
        if abs(self.botX.evalState({ 'state' : state, 'moves' : [self.lastmove]}, self.botX.notPlayer(player))) == 1000000:
            return 1
        elif len(self.botX.genChildren({ 'state' : state, 'moves' : [self.lastmove]}, self.botX.notPlayer(player))) == 0:
            return 2
        else:
            return 0

SCORE_EXACT = 0
SCORE_UPPER_BOUND = 1
SCORE_LOWER_BOUND = 2


def state_to_xy(state):
    return [
        [state[b][c] for b in range(x,9,3) for c in range(y,9,3)]
        for x in range(3) for y in range(3)
        ]
    
def xy_to_state(xy):
    return [
        [xy[x][y] for y in range(b,b+3) for x in range(r,r+3)]
        for b in range(0,9,3) for r in range(0,9,3)
        ]

def rot_state(orient):
    def rot(state):
        if orient == 0:
            return state
        elif orient == 1:
            return xy_to_state(zip(*state_to_xy(state)[::-1]))
        elif orient == 2:
            return xy_to_state(map(list,map(reversed,state_to_xy(state)))[::-1])
        elif orient == 3:
            return xy_to_state(zip(*state_to_xy(state))[::-1])
    return rot

def rot_move(move, orient):
    newmv = []
    for i in range(2):
        for d in range(3):
            if move[i] in rot_dict[d]:
                newmv.append(rot_dict[d][(rot_dict[d].index(move[i]) + orient) % 4])
    return newmv

def move_to_xy(b, c):
    bxy = divmod(b,3)
    cxy = divmod(c,3)
    x = cxy[0]*3 + bxy[1]
    y = bxy[0]*3 + cxy[1]
    return (x,y)

# Generate 3 9x9 boxes of zobrist keys intended for the different cases:
# - Square is empty
# - Square has X
# - Square has O
def gen_zobrist_keys():
    grids = []
    for i in range(3):
        grid = [[random.getrandbits(64) for b in range(9)] for c in range(9)]
        grids.append(grid)
    return grids

zobrist_keys = [map(rot_state(i), gen_zobrist_keys()) for i in range(4)]

# rotation lists, each next number is a rotation counter clockwise
rot_dict = [
    [0,6,8,2],
    [1,3,7,5],
    [4,4,4,4]
    ]

#def orient_of_hash(h):
    #for i in range(4):
        #if h in hash_tables[i]:
            #return i

def canon_hash(h):
    for i in range(4):
        if i == 0:
            minh = h[i]
        else:
            minh = min(minh,hash)
    return minh
    

def hash_square(state, b, c, orient):
    contents = state[b][c]
    ob, oc = rot_move((b,c), orient)
    if contents == 'X':
        hash = zobrist_keys[orient][0][ob][oc]
    elif contents == 'O':
        hash = zobrist_keys[orient][1][ob][oc]
    else:
        hash = zobrist_keys[orient][2][ob][oc]
    return hash

def hash_board(state):
    hashes = []
    for ori in range(4):
        hash = 0
        for b in range(9):
            for c in range(9):
                hash ^= hash_square(state, b, c, ori)
        hashes.append(hash)
    return hashes

nodecount = 0

class XOHeuristic:
    
    depth = 1
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
    
    
    alpha = -1000000
    beta  =  1000000

    cache = {}
    
    def move(self, state, lastmove, player):
        global nodecount
        nodecount = 0
        self.colour = player
        start_clock = time.clock()

        if player == 'X':
            self.depth = 8
        else:
            self.depth = 8
        self.colour = player

        moves = []
        moves.append(lastmove)
        # call alphabeta with depth, it will return the end state and the list of moves leading to it
        node = { 'state' : state, 'moves' : moves, 'cache' : self.cache, 'hash' : hash_board(state) }

        result = self.negamax(node, self.depth, self.alpha, self.beta, player)
        
        print "Nodes visited: ", nodecount
        print "Cache size: ", len(node['cache'])
        move = result[1]
        print "Score for " + str(move) + " : " +  str(result[0])
        
        state[move[0]][move[1]] = player
        print 'move took', start_clock - time.clock()
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
        random.shuffle(moves)
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
            #if the board is won override and return the max/min value
            if abs(self.score[players][others]) == 1000:
                return self.score[players][others]
            t += self.score[players][others]
        return t
    
    # compute score for whole metaboard
    def evalState(self, node, player):
        t = 0
        state = node['state']
        scores = [self.evalBoard(node['state'][i],player) for i in range(9)]
        for i in range(8):
            players = 0
            others = 0
            for j in range(3):
                score = scores[self.wins[i][j]]
                if score == 1000:
                    players += 1
                elif score == -1000:
                    others += 1
            t += 1000 * self.score[players][others]
            
            #if the board is won override and return the max/min value
            if 1000 * abs(self.score[players][others]) >= 1000000:
                t = 1000 * self.score[players][others]
                break
        t += sum(scores)
        #cap the max/min value to 1000
        if t > 1000000:
            t = 1000000
        elif t < -1000000:
            t = -1000000
        #print player + ": " + str(t)#lastmove[1]
        return t

    def update_cache(self, node, entry):
        cache = node["cache"]
        h = canon_hash(node["hash"])
        if h in cache:
            old = cache[h]
            if old[2] > entry[2]:
                return # Value was already better
            if old[2] == entry[2]:
                if old[3] > entry[3]:
                    return # Value has better bound type
        cache[h] = entry

    def negamax(self, node, depth, a, b, maximize):
        #print depth
        cache = node["cache"]
        h = canon_hash(node["hash"])
        if h in cache:
            entry = cache[h]
            if entry[2] >= depth:
                if entry[3] == SCORE_EXACT:
                    return entry
                b = min(entry[1], b)
                if a >= b:
                    a = b
        children = self.genChildren(node, maximize)
        move = []

        if depth == 0 or len(children) == 0:
            score = self.evalState(node, maximize)
            ret = [score, []]
            return ret
        for child in children:
            if move == []: move = child
            board = child[0]
            cell = child[1]
            
            node['moves'].append(child)
            node['state'][board][cell] = maximize
            
            for i in range(4):
                hsquare = hash_square(node['state'],board,cell,i)
                node['hash'][i] ^= hsquare
            
            global nodecount
            nodecount += 1
            comp = self.negamax(node, depth - 1, -b, -a, self.notPlayer(maximize))
            
           # print "Children " + str(children)
           # print "Depth: " + str(depth) + " ret: " + str(comp)
           # print node['moves']
           # if abs(comp[0]) < 10 and comp[0] != 0:
           #    exit(0)
           
           
            for i in range(4):
                hsquare = hash_square(node['state'],board,cell,i)
                node['hash'][i] ^= hsquare
            node['state'][board][cell] = ' '
            node['moves'].pop()
            
            move = child #rot_move(child, ori)
            
            if comp[0] >= b:
                movepair = [comp[0], move, depth, SCORE_UPPER_BOUND]
                self.update_cache(node, movepair)
                return movepair
            if comp[0] > a:
                a = comp[0]

        #print a
        movepair = [a, move, depth, SCORE_EXACT]
        self.update_cache(node, movepair)
        return movepair
        
    # 
    # 
    def alphabeta(self, node, depth, a, b, maximize):
        global nodecount
        children = self.genChildren(node, maximize)
        move = []

        if depth == 0 or len(children) == 0:
            score = self.evalState(node, maximize)
            ret = [score, []]
            return ret
        if maximize == self.colour:
            for child in children:
                board = child[0]
                cell = child[1]
                
                node['moves'].append(child)
                node['state'][board][cell] = maximize
                
                nodecount += 1
                comp = self.alphabeta(node, depth - 1, a, b, self.notPlayer(maximize))
                
                #print "Children " + str(children)
                #print "Depth: " + str(depth) + " ret: " + str(comp)
                #print node['moves']
                #if abs(comp[0]) < 10 and comp[0] != 0:
                    #exit(0)
                
                node['state'][board][cell] = ' '
                node['moves'].pop()
                
                if comp[0] >= a:
                    a = comp[0]
                    move = child
                if b <= a:
                    break #beta cutoff
            #print a
            return [a, move]
        else:
            for child in children:
                board = child[0]
                cell = child[1]
                
                node['moves'].append(child)
                node['state'][board][cell] = maximize
                
                nodecount += 1
                comp = self.alphabeta(node, depth - 1, a, b, self.notPlayer(maximize))
                
                #print "Children " + str(children)
                #print "Depth: " + str(depth) + " ret: " + str(comp)
                #print node['moves']
                #if abs(comp[0]) < 10 and comp[0] != 0:
                    #exit(0)
                    
                node['state'][board][cell] = ' '
                node['moves'].pop()
                
                if comp[0] <= b:
                    b = comp[0]
                    move = child
                if b <= a:
                    break #alpha cutoff
            return [b, move]
    
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

