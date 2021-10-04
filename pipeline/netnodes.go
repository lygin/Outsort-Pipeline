package pipeline

import (
	"bufio"
	"net"
)

// Write chanel date to net
func NetworkSink(addr string, in <-chan int) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	go func() {
		defer listener.Close()
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		writer := bufio.NewWriter(conn)
		// Remember to flush buffer
		defer writer.Flush()
		WriterSink(writer, in)
	}()
}

// Get date to chanel from net
func NetworkSource(addr string) <-chan int {
	out := make(chan int)
	go func() {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			panic(err)
		}
		// Read net file to chanel
		r := ReaderSource(bufio.NewReader(conn), -1)
		for v := range r {
			out <- v
		}
		close(out)
	}()
	return out
}
