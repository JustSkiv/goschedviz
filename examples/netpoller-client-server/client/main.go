package main

import (
	"fmt"
	"net"
	"time"

	"github.com/JustSkiv/goschedviz/pkg/metrics"
)

const numClients = 150 // Number of connections

func startClient(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Connection error for client: %v\n", err)
		return
	}

	localAddr := conn.LocalAddr()

	defer func() {
		fmt.Printf("Client disconnecting: %v\n", localAddr)
		conn.Close()
	}()

	for {
		n, err := conn.Write([]byte("ping"))
		if err != nil {
			fmt.Printf("Write error for %v: %v (bytes written: %d)\n",
				localAddr, err, n)
			return
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func main() {
	// Start metrics export
	reporter := metrics.NewReporter(time.Second)
	reporter.Start()
	defer reporter.Stop()

	// Start all clients
	for i := 0; i < numClients; i++ {
		go startClient("127.0.0.1:8080")
	}

	time.Sleep(1 * time.Hour)
}
