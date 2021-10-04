package main

import (
	"fmt"
)

func main() {
	ch := make(chan string)
	for i := 0; i < 5000; i++ {
		go printHello(i, ch)
	}
	for {
		msg := <-ch
		fmt.Println(msg)
	}
}

func printHello(x int, ch chan string) {
	for {
		ch <- fmt.Sprintf("Hello %d\n", x)
	}
}
