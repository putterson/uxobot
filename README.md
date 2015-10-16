# uxobot
Ultimate Tic Tac Toe Bot

This is a command line implementation of the game described at http://mathwithbaddrawings.com/ultimate-tic-tac-toe-original-post/
with the clarifying rules:

1. If you are sent to a sub board that is already won, you can play anywhere in any non-won board
2. Tied boards count towards neither player

Included are two AI players. One is based on the negamax algorithm with zobrist hashing to cache node values and uses an imperfect heuristic to evaluate unfinished games.
The second is based on a Monte Carlo tree search algorithm with the UCT extension to balance exploration vs exploitation.

# Making
```
go build
```
or
```
go install
```

uxobot does not depend on any external libraries

# Running
./uxobot [opts]
Specify the two players with the -p1 and -p2 options. Valid values are 'human', 'negamax', and 'montecarlo'

Specify the strength of the AI player with the -s1 and -s2 options.
For 'montecarlo' the strength number is the number of seconds to take per turn and accepts decimal values
For 'negamax' the strength number is the depth that it will search to and must be an integer.

To make a move at the [move]> prompt type two numbers that correspond to the board and cell you would like to move in. Boards and cells are numbered like:
```
1 2 3
4 5 6
7 8 9
```
So to move in the centre board in the top left cell you would type:
```
[move]>5 1
```
# Example
./uxobot -p1 montecarlo -s1 1 -p2 human

This sets the montecarlo player to X and human player to O. It also sets the montecarlo player to take 1 second per turn.
