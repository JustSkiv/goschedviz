package main

import (
	"time"
)

func main() {
	for i := 0; i < 1000; i++ {
		go func(id int) {
			for {
				// Активная работа вместо сна
				for i := 0; i < 10000000; i++ {
					_ = i * i
				}
				time.Sleep(time.Millisecond) // небольшая пауза
			}
		}(i)
	}

	time.Sleep(1 * time.Minute)
}
