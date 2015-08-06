package main

type UXOBot interface {
	getMove(node *AINode) Move
}
