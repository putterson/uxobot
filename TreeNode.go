package main

type NodeChildren map[BitMove]*TreeNode

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

func NewTreeNode(board *BitBoard, lastmove *BitMove) *TreeNode {
	moveslice := *NewBitMoveSlice()
	slen := 0
	genBitChildren(board, lastmove, &moveslice, &slen)
	return &TreeNode{
		outcomes:   1,
		wincomes:   [2]int{0,0},
		childMoves: moveslice,
		nChildMoves: slen,
		childNodes: make(NodeChildren),
	}
}