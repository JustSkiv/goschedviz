//gocoverage:ignore
package main

import (
	"time"

	"github.com/JustSkiv/goschedviz/pkg/metrics"
)

func main() {
	reporter := metrics.NewReporter(time.Second)
	reporter.Start()
	defer reporter.Stop()

	for i := 0; i < 1000; i++ {
		go func(id int) {
			for {
				// Активная работа вместо сна
				for i := 0; i < 10000000; i++ {
					_ = i * i
				}
			}
		}(i)
	}

	time.Sleep(1 * time.Minute)
}
