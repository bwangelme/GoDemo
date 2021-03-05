package main

import (
	"fmt"
)

type AbstractStack struct {
	elems []interface{}
}

func (s *AbstractStack) Push(elem interface{}) {
	s.elems = append(s.elems, elem)
}

func (s *AbstractStack) Pop() interface{} {
	return s.elems[0]
}

type IntStack struct {
	basicStack *AbstractStack
}

func NewIntStack() *IntStack {
	return &IntStack{
		basicStack: new(AbstractStack),
	}
}

func (s *IntStack) Push(elem int) {
	s.basicStack.Push(elem)
}

func (s *IntStack) Pop() int {
	e := s.basicStack.Pop()

	v, _ := e.(int)
	return v
}

func main() {
	var v interface{}
	v = "10"
	i := v.(int)
	fmt.Println(i)
}
