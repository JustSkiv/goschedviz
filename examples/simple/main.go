package main

import (
	"time"
)

func main() {
	// Simple program that creates some scheduler load
	for i := 0; i < 1000; i++ {
		go func() {
			time.Sleep(time.Second)
		}()
	}
	time.Sleep(10 * time.Second)
}
