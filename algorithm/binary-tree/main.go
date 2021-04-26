package main

import (
	"fmt"
	"math/rand"
)

type BST struct {
	Value       int
	Left, Right *BST
}

func (tree *BST) Insert(value int) *BST {
	if value < tree.Value {
		if tree.Left == nil {
			tree.Left = &BST{Value: value}
		} else {
			tree.Left.Insert(value)
		}
	} else {
		if tree.Right == nil {
			tree.Right = &BST{Value: value}
		} else {
			tree.Right.Insert(value)
		}
	}
	return tree
}

func (tree *BST) Contains(value int) bool {
	if tree == nil {
		return false
	} else if tree.Value == value {
		return true
	}

	if value < tree.Value {
		return tree.Left.Contains(value)
	} else {
		return tree.Right.Contains(value)
	}
}

func (tree *BST) Remove(value int) *BST {
	tree.remove(value, nil)
	return tree
}

func (tree *BST) remove(value int, parent *BST) {

}

func (tree *BST) Display() {
	if tree == nil {
		return
	}
	tree.Left.Display()
	fmt.Println(tree.Value)
	tree.Right.Display()
}

func main() {
	tree := &BST{Value: rand.Intn(1000)}

	for i := 0; i < 100; i++ {
		tree.Insert(rand.Intn(1000))
	}
	tree.Display()
	fmt.Println(tree.Contains(2), tree.Contains(10))
}
