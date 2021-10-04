package pipeline

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"time"
)

func ArrayToChan(a ...int) <-chan int { //返回的chan只拿不放
	out := make(chan int, 1024)
	go func() { //新开一个协程放数据 避免等待
		for _, v := range a {
			out <- v
		}
		close(out)
	}()
	return out
}

var startTime time.Time

func Init() {
	startTime = time.Now()
}

// Sort chanel in mem
func InMemSort(in <-chan int) <-chan int {
	out := make(chan int, 1024)
	go func() { //新开协程进行排序，不消耗时间
		a := []int{}
		for v := range in {
			a = append(a, v)
		}
		fmt.Println("Read Done: ", time.Since(startTime))
		//Sort
		sort.Ints(a)
		fmt.Println("Sort Done: ", time.Since(startTime))
		//Out
		for _, v := range a {
			out <- v
		}
		close(out)
	}()
	return out
}

// Merge two ordered channel into one which is also ordered
func Merge(in1, in2 <-chan int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		v1, ok1 := <-in1
		v2, ok2 := <-in2
		for ok1 || ok2 {
			if !ok2 || (ok1 && v1 <= v2) {
				out <- v1
				v1, ok1 = <-in1
			} else {
				out <- v2
				v2, ok2 = <-in2
			}
		}
		close(out)
		fmt.Println("Merge Done: ", time.Since(startTime))
	}()
	return out
}

// Read less than chunkSize bytes file into chanel
func ReaderSource(reader io.Reader, chunkSize int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		buffer := make([]byte, 8)
		bytesRead := 0
		for {
			n, err := reader.Read(buffer)
			bytesRead += n
			if n > 0 {
				v := int(binary.BigEndian.Uint64(buffer))
				out <- v
			}
			if err != nil || (chunkSize != -1 && bytesRead >= chunkSize) {
				break
			}
		}
		close(out)
	}()
	return out
}

// Write chanel date into file
func WriterSink(writer io.Writer, in <-chan int) {
	for v := range in {
		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, uint64(v))
		writer.Write(buffer)
	}
}

// Generate random chanel which size is count
func RandomSource(count int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		for i := 0; i < count; i++ {
			out <- rand.Int()
		}
		close(out)
	}()
	return out
}

// Merge list of ordered chanel into one ordered channel
// in : list of chanel
// out : one chanel
func MergeN(inputs ...<-chan int) <-chan int { //多路两两归并
	if len(inputs) == 1 {
		return inputs[0]
	}
	m := len(inputs) / 2
	//input[0:m)and input[m..end)
	return Merge(
		MergeN(inputs[:m]...), //左边排好了
		MergeN(inputs[m:]...)) //右边排好了
}
