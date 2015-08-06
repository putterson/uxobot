package main

type MonteCarlo struct {
	root TreeNode
}

type TreeNode struct {
	children []*TreeNode
}

func (m MonteCarlo)  getMove(node *AINode) Move {
	return NoMove()
}
