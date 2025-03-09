package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/JustSkiv/goschedviz/pkg/metrics"
)

func handleConnection(conn net.Conn) {
	defer func() {
		fmt.Printf("Connection closed: %v\n", conn.RemoteAddr())
		conn.Close()
	}()

	written, err := io.Copy(io.Discard, conn)
	if err != nil {
		fmt.Printf("Error reading from %v: %v (bytes read: %d)\n",
			conn.RemoteAddr(), err, written)
	}
}

func main() {
	// Start metrics export
	reporter := metrics.NewReporter(time.Second)
	reporter.Start()
	defer reporter.Stop()

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Server start error:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server started on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				fmt.Println("Temporary Accept error:", err)
				continue
			}
			fmt.Println("Critical Accept error:", err)
			return
		}
		go handleConnection(conn)
	}
}
