package main

import "fmt"

func main() {
	size := 4
	a := NewBuffer(size, 0)
	s := []string{"a", "b", "c", "d", "e", "f"}

	for i, str := range s {
		a.Insert(str)
		fmt.Printf("%d %2v IsFull: %v Len: %d Head: %d\n", i, a.list, a.IsFull(), a.len, a.head)
	}
	fmt.Println("A Flush:", a.Flush())

	b := NewBuffer(size, 0)
	b.InsertMultiple(s)
	fmt.Println(b.list)
}
