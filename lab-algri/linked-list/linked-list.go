package main

type ILinkedList interface {
	Size() int
	IsEmpty() bool
	Contains(e interface{}) bool
	Add(e interface{})
	Get(i int) interface{}
	Set(i int, e interface{})
	AddToIndex(i int, e interface{})
	Remove(i int) interface{}
	IndexOf(e interface{}) int
	Clear()
	String()
}

type LinkedList struct {
	size int
	node *node
}

type node struct {
	val  interface{}
	next *node
}

func (l LinkedList) Size() int {
	panic("implement me")
}

func (l LinkedList) IsEmpty() bool {
	panic("implement me")
}

func (l LinkedList) Contains(e interface{}) bool {
	panic("implement me")
}

func (l LinkedList) Add(e interface{}) {
	panic("implement me")
}

func (l LinkedList) Get(i int) interface{} {
	panic("implement me")
}

func (l LinkedList) Set(i int, e interface{}) {
	panic("implement me")
}

func (l LinkedList) AddToIndex(i int, e interface{}) {
	panic("implement me")
}

func (l LinkedList) Remove(i int) interface{} {
	panic("implement me")
}

func (l LinkedList) IndexOf(e interface{}) int {
	panic("implement me")
}

func (l LinkedList) Clear() {
	panic("implement me")
}

func (l LinkedList) String() {
	panic("implement me")
}
