// Package ui provides interfaces and implementations for metrics visualization
package ui

import (
	"schedtrace-mon/internal/collector"
	"schedtrace-mon/internal/domain"
)

// Presenter defines interface for displaying scheduler metrics
type Presenter interface {
	// Init initializes the UI system
	Init() error

	// Close performs cleanup of UI resources
	Close() error

	// Display shows current metrics from the collector
	Display(collector.Collector)

	// HandleEvents sets up UI event handling
	HandleEvents(handler func(domain.Event))
}
