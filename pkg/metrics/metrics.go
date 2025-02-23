// Package metrics provides functionality for exporting runtime metrics
// that can be consumed by goschedviz monitoring tool.
//
// Example usage:
//
//	reporter := metrics.NewReporter(time.Second)
//	reporter.Start()
//	defer reporter.Stop()
package metrics

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	// prefix is used to distinguish process metrics output from other stderr content
	prefix = "PROCMETR"
)

// Reporter handles periodic metrics reporting
type Reporter struct {
	interval time.Duration
	done     chan struct{}
	stopOnce sync.Once
}

// NewReporter creates a new metrics reporter that will output metrics
// at the specified interval
func NewReporter(interval time.Duration) *Reporter {
	return &Reporter{
		interval: interval,
		done:     make(chan struct{}),
	}
}

// Start begins periodic metrics reporting.
// It runs in a separate goroutine until Stop is called.
// Multiple calls to Start are safe but only the first one takes effect.
func (r *Reporter) Start() {
	go func() {
		ticker := time.NewTicker(r.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				r.report()
			case <-r.done:
				return
			}
		}
	}()
}

// Stop terminates metrics reporting.
// Multiple calls to Stop are safe but only the first one takes effect.
func (r *Reporter) Stop() {
	r.stopOnce.Do(func() {
		close(r.done)
	})
}

// report outputs current metrics to stderr
func (r *Reporter) report() {
	// Get number of goroutines
	numG := runtime.NumGoroutine()

	// Output in a format that can be parsed by the monitor
	fmt.Fprintf(os.Stderr, "%s num_goroutines=%d\n", prefix, numG)
}
