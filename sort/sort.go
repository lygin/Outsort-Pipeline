package main

import (
	"bufio"
	"fmt"
	"go-demo-1/pipeline"
	"os"
	"strconv"
)

func main() {
	p := createNetworkPipeline("small.in", 4000000, 4)
	writeToFile(p, "small.out")
	printFile("small.out")
}

func printFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	count := 0
	p := pipeline.ReaderSource(file, -1)
	for v := range p {
		fmt.Println(v)
		count++
		if count >= 100 {
			break
		}
	}
}

// Write chanel date to file
func writeToFile(p <-chan int, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	pipeline.WriterSink(file, p)
}

// Read date from file to nums of chanels
// Then sort each chanels in mem and put them into lists
// Finally, Merge these sorted Chanels
// IN: fileSize bytes
// IN: chunkCount nums of divided chunk
func createPipeline(filename string, fileSize, chunkCount int) <-chan int {
	chunkSize := fileSize / chunkCount
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	pipeline.Init()
	sortRes := []<-chan int{}
	for i := 0; i < chunkCount; i++ {
		file.Seek(int64(i*chunkSize), 0)
		source := pipeline.ReaderSource(bufio.NewReader(file), chunkSize)
		sortRes = append(sortRes, pipeline.InMemSort(source))
	}
	return pipeline.MergeN(sortRes...)
}

// Read date from file to chanel
// Then sort chanels and put them to each net addr
// sortAddr stores the sorted res
// Finally, Read date form net to chanel
// Return Merged chanel
func createNetworkPipeline(filename string, fileSize, chunkCount int) <-chan int {
	chunkSize := fileSize / chunkCount
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	pipeline.Init()
	sortAddr := []string{}
	for i := 0; i < chunkCount; i++ {
		file.Seek(int64(i*chunkSize), 0)
		source := pipeline.ReaderSource(bufio.NewReader(file), chunkSize)
		addr := ":" + strconv.Itoa(i+7000)
		pipeline.NetworkSink(addr, pipeline.InMemSort(source))
		sortAddr = append(sortAddr, addr)
	}
	sortRes := []<-chan int{}
	for _, addr := range sortAddr {
		sortRes = append(sortRes, pipeline.NetworkSource(addr))
	}
	return pipeline.MergeN(sortRes...)
}
