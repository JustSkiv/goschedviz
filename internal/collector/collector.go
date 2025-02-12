// Package collector provides interfaces and implementations for gathering Go scheduler metrics.
package collector

import (
	"context"

	"github.com/JustSkiv/goschedviz/internal/domain"
)

// Collector defines interface for gathering scheduler metrics from various sources.
//
// System layout:
//
//	┌──────────────────┐
//	│  Target Process  │
//	│  with GODEBUG    │──┐
//	└──────────────────┘  │
//	                      │ stderr
//	┌──────────────────┐  │
//	│    Collector     │◄─┘
//	│  ┌────────────┐  │
//	│  │  Parser    │  │
//	│  └────────────┘  │
//	└──────────────────┘
type Collector interface {
	// Start begins collecting metrics from the target process.
	// It returns a channel that will receive scheduler snapshots.
	Start(ctx context.Context) (<-chan domain.SchedulerSnapshot, error)

	// Stop gracefully stops the collection process.
	Stop() error
}
