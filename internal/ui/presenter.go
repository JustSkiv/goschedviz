// Package ui provides interfaces and implementations for visualizing scheduler metrics.
package ui

// Presenter defines interface for any UI implementation that can visualize scheduler metrics.
//
// UI Layout Reference:
//
//	┌─────────────────────────────────┬─────────────────────────────────┐
//	│     Current Values Table        │      Local Run Queue Bars       │
//	│    (30% height, 40% width)      │     (30% height, 60% width)     │
//	├─────────────────────────────────┴─────────────────────────────────┤
//	│                    Run Queue Gauges                               │
//	│           GRQ (50% width) | LRQ Sum (50% width)                  │
//	│                       (30% height)                               │
//	├─────────────────────────────────┬─────────────────────────────────┤
//	│       History Plot              │           Info Box              │
//	│    (40% height, 80% width)      │    (40% height, 20% width)      │
//	└─────────────────────────────────┴─────────────────────────────────┘
type Presenter interface {
	// Start initializes and starts the UI.
	// Returns error if initialization fails.
	Start() error

	// Stop gracefully shuts down the UI.
	Stop()

	// Update updates UI with new metrics data.
	Update(data UIData)

	// Done returns a channel that's closed when UI should exit
	// (e.g., user pressed 'q' or Ctrl+C).
	Done() <-chan struct{}
}
