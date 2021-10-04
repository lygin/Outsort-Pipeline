package main

import (
	"bufio"
	"fmt"
	"go-demo-1/mascot"
	"go-demo-1/pipeline"
	"os"
)

func main() {
	//Create file
	const filename string = "small.in"
	const n int = 1000000
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	p := pipeline.RandomSource(n)
	writer := bufio.NewWriter(file)
	pipeline.WriterSink(writer, p)
	writer.Flush()

	file, err = os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	p = pipeline.ReaderSource(file, -1)
	for v := range p {
		fmt.Println(v)
	}

	ch := pipeline.Merge(
		pipeline.InMemSort(pipeline.ArrayToChan(3, 2, 6, 9, 1)),
		pipeline.InMemSort(pipeline.ArrayToChan(4, 0, 1, 8, 7)))
	/*
		for {
			if num, ok := <-ch; ok {
				fmt.Println(num)
			} else {
				break
			}

		}
	*/
	for v := range ch {
		fmt.Println(v)
	}
	fmt.Println(mascot.Hello())
}
