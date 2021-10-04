package main

import (
	"fmt"
	"sort"
)

func main() {
	a := []int{1, 5, 3, 2, 4}
	sort.Ints(a)
	for i, v := range a {
		fmt.Println(i, v)
	}
}
