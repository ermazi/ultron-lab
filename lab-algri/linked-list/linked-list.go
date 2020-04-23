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

func (l *LinkedList) Size() int {
	// 遍历链表获取
	return l.size
}

func (l *LinkedList) IsEmpty() bool {
	return l.size == 0
}

func (l *LinkedList) Contains(e interface{}) bool {
	n := l.node

	for n != nil {
		if n.val == e {
			return true
		}
	}

	return false
}

func (l *LinkedList) Add(e interface{}) {
	newNode := &node{
		val:  e,
		next: nil,
	}
	n := l.node

	if n == nil {
		n = newNode
	}
	for {
		if n.next == nil {
			n.next = newNode
			break
		}
	}
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
