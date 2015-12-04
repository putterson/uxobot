package main

type NodeChildren []TreeNode

type TreeNode struct {
	outcomes   int
	wincomes   [2]int
	move       BitMove
	childMoves BitMoveSlice
	nChildMoves int
	childNodes NodeChildren
}

func (t *TreeNode) hasChildNodes() bool {
	return len(t.childNodes) > 0
}

func (t *TreeNode) nNextMoves() int {
	return len(t.childMoves)
}

func (t *TreeNode) getMove(n int) BitMove {
	return t.childMoves[n]
}

func (t *TreeNode) getChild(move BitMove) (*TreeNode, bool) {
	for _, node := range (t.childNodes) {
		if node.move == move {
			return &node, true
		}
	}
	return nil, false
}

func (t *TreeNode) addChild(board *BitBoard, move *BitMove) *TreeNode {
	newNode := NewTreeNode(board, move)
	t.childNodes = append(t.childNodes, *newNode)
	return newNode
}

func NewTreeNode(board *BitBoard, lastmove *BitMove) *TreeNode {
	moveslice := *NewBitMoveSlice()
	slen := 0
	genBitChildren(board, lastmove, &moveslice, &slen)
	return &TreeNode{
		outcomes:   1,
		wincomes:   [2]int{0,0},
		move: BitMove{lastmove.s, lastmove.c,},
		childMoves: moveslice,
		nChildMoves: slen,
		childNodes: make(NodeChildren, 0, 9),
	}
}
