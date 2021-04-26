package main

import "fmt"

type LinkedList struct {
	Value int
	Next  *LinkedList
}

func main() {

	linkedList := &LinkedList{Value: 0}

	for i, tmp := 1, linkedList; i < 1024; i++ {
		tmp.Next = &LinkedList{Value: i}
		tmp = tmp.Next
	}

	traverseLinkedList(linkedList)
}

func traverseLinkedList(node *LinkedList) {
	if node == nil {
		return
	}
	fmt.Println(node.Value)
	traverseLinkedList(node.Next)
}
