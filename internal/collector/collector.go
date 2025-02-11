// Package collector provides interfaces and implementations for gathering scheduler metrics
package collector

import (
	"context"
	"schedtrace-mon/internal/domain"
)

// Collector defines interface for gathering and storing scheduler metrics
type Collector interface {
	// Start begins collecting metrics for the specified command.
	// The collection continues until context is cancelled or error occurs.
	Start(ctx context.Context, targetCmd string) error

	// Stop terminates metric collection and cleans up resources
	Stop() error

	// GetCurrent returns the most recent scheduler state
	GetCurrent() domain.SchedData

	// GetHistory returns historical scheduler states
	GetHistory() []domain.SchedData
}

// MaxHistoryPoints defines how many data points to keep in history
const MaxHistoryPoints = 60
