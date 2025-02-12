// Package mock provides a mock collector for testing purposes.
package mock

import (
	"context"
	"time"

	"github.com/JustSkiv/goschedviz/internal/domain"
)

// Collector is a mock implementation of collector.Collector interface.
type Collector struct {
	snapshots []domain.SchedulerSnapshot
	interval  time.Duration
	done      chan struct{}
}

// New creates a mock collector that will emit the provided snapshots.
func New(snapshots []domain.SchedulerSnapshot, interval time.Duration) *Collector {
	return &Collector{
		snapshots: snapshots,
		interval:  interval,
		done:      make(chan struct{}),
	}
}

// Start implements collector.Collector interface.
func (c *Collector) Start(ctx context.Context) (<-chan domain.SchedulerSnapshot, error) {
	out := make(chan domain.SchedulerSnapshot)

	go func() {
		defer close(out)

		ticker := time.NewTicker(c.interval)
		defer ticker.Stop()

		i := 0
		for {
			select {
			case <-ctx.Done():
				return
			case <-c.done:
				return
			case <-ticker.C:
				if i >= len(c.snapshots) {
					i = 0
				}
				out <- c.snapshots[i]
				i++
			}
		}
	}()

	return out, nil
}

// Stop implements collector.Collector interface.
func (c *Collector) Stop() error {
	close(c.done)
	return nil
}
