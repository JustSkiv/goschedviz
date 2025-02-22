package main

import (
	"net/http"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(2)

	// Используем IP напрямую вместо доменного имени
	for i := 0; i < 10000; i++ {
		go func() {
			for {
				http.Get("http://216.58.210.142") // IP Google
			}
		}()
	}

	time.Sleep(time.Hour)
}
