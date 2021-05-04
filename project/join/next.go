package main

import "fmt"

type Node struct {
	Key      string  `json:"name"`
	Count    int     `json:"count"`
	Parent   *Node   `json:"-"`
	Children []*Node `json:"children,omitempty"`
}

func (node *Node) Print() {
	if node == nil {
		return
	}

	fmt.Println(node.Key)

	for _, child := range node.Children {
		child.Print()
	}
}

func (node *Node) GetParentKey() string {
	return node.Key[:len(node.Key)-1]
}

func (node *Node) Split() {
	if len(node.Children) == 0 {
		for i := 65; i < 81; i++ {
			child := &Node{Key: node.Key + string(byte(i)), Parent: node}
			node.Children = append(node.Children, child)
		}
	}
}
