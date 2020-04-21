package main

import (
	"container/heap"
	"fmt"
)

type Item struct {
	Name   string
	Expiry int
	Index  int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Expiry < pq[j].Expiry
}

func (pq *PriorityQueue) Pop() interface{} {
	item := (*pq)[len(*pq)-1]
	item.Index = -1
	*pq = (*pq)[0 : len(*pq)-1]
	return item
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	item.Index = len(*pq)

	*pq = append(*pq, item)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = pq[j].Index
	pq[j].Index = pq[i].Index
}

func main() {
	listItems := []*Item{
		{Name: "Carrot", Expiry: 30},
		{Name: "Potato", Expiry: 45},
		{Name: "Rice", Expiry: 100},
		{Name: "Spinach", Expiry: 5},
	}

	priorityQueue := make(PriorityQueue, 0)

	for _, item := range listItems {
		priorityQueue.Push(item)
	}

	heap.Init(&priorityQueue)

	for priorityQueue.Len() > 0 {
		item := heap.Pop(&priorityQueue).(*Item)
		fmt.Printf("Name: %s Expiry: %d\n", item.Name, item.Expiry)
	}
}
