package main

import "fmt"

type IDynamicArray interface {
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

type DynamicArray struct {
	arr []interface{}
}

func NewDynamicArray() IDynamicArray {
	return &DynamicArray{
		arr: make([]interface{}, 0),
	}
}

func (d DynamicArray) Size() int {
	return len(d.arr)
}

func (d DynamicArray) IsEmpty() bool {
	if d.Size() == 0 {
		return true
	}
	return false
}

func (d DynamicArray) Contains(e interface{}) bool {
	for _, item := range d.arr {
		if item == e {
			return true
		}
	}
	return false
}

func (d *DynamicArray) Add(e interface{}) {
	d.arr = append(d.arr, e)
}

func (d DynamicArray) Get(i int) interface{} {
	return d.arr[i]
}

func (d *DynamicArray) Set(i int, e interface{}) {
	d.arr[i] = e
}

func (d *DynamicArray) AddToIndex(i int, e interface{}) {
	former := d.arr[0:i]
	copyLen := d.Size() - i
	latter := make([]interface{}, copyLen)
	copy(latter, d.arr[i:d.Size()])
	newArr := append(former, e)
	newArr = append(newArr, latter...)
	d.arr = newArr
}

func (d *DynamicArray) Remove(i int) interface{} {
	eToDelete := d.arr[i]
	former := d.arr[0:i]
	copyLen := d.Size() - i - 1
	latter := make([]interface{}, copyLen)
	copy(d.arr[i+1:d.Size()], latter)
	newArr := append(former, latter...)

	d.arr = newArr
	return eToDelete
}

func (d DynamicArray) IndexOf(e interface{}) int {
	for i, item := range d.arr {
		if item == e {
			return i
		}
	}
	panic("out of indices")
}

func (d *DynamicArray) Clear() {
	d.arr = make([]interface{}, 0)
}

func (d DynamicArray) String() {
	for _, item := range d.arr {
		fmt.Printf("%v ", item)
	}
	fmt.Println()
}

func main() {
	arr := NewDynamicArray()
	arr.Add(1)
	arr.Add(2)
	arr.Add(3)
	arr.Add(4)

	arr.String()

	fmt.Println(arr.Get(0))
	fmt.Println(arr.Get(3))

	arr.AddToIndex(2, 5)
	arr.String()

	arr.Set(2, 9)
	arr.String()

	arr.AddToIndex(2, 6)
	arr.String()

	fmt.Println(arr.Contains(9))
	fmt.Println(arr.Contains(11))
	fmt.Println(arr.IndexOf(9))
	fmt.Println(arr.Size())
	arr.Clear()

	arr.String()
}
