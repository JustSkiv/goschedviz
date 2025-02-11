// Package domain defines the core types and interfaces for the Go scheduler monitoring system.
package domain

import "sync"

// SchedulerSnapshot represents parsed values from a single "SCHED" trace line.
// It contains various metrics about Go runtime scheduler state at a specific moment.
type SchedulerSnapshot struct {
	TimeMs          int   // Time since start in milliseconds
	GoMaxProcs      int   // Current GOMAXPROCS value
	IdleProcs       int   // Number of idle processors
	Threads         int   // Total number of threads
	SpinningThreads int   // Number of spinning threads
	NeedSpinning    int   // Number of threads that need spinning
	IdleThreads     int   // Number of idle threads
	RunQueue        int   // Global Run Queue (GRQ) length
	LRQSum          int   // Sum of all Local Run Queues
	LRQ             []int // Local Run Queue length for each P
}

// MonitorState maintains the current state and history of scheduler metrics.
//
// Layout visualization:
//
//	┌─────────────────────────────┐
//	│ Current State (latest)      │
//	├─────────────────────────────┤
//	│ History (last N snapshots)  │
//	└─────────────────────────────┘
type MonitorState struct {
	mu      sync.Mutex
	latest  SchedulerSnapshot
	history []SchedulerSnapshot
}

// MaxHistoryPoints defines how many data points we keep for plotting
const MaxHistoryPoints = 60

// Update saves new snapshot and adds it to history, maintaining max history size
func (ms *MonitorState) Update(data SchedulerSnapshot) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.latest = data
	ms.history = append(ms.history, data)
	if len(ms.history) > MaxHistoryPoints {
		ms.history = ms.history[len(ms.history)-MaxHistoryPoints:]
	}
}

// GetSnapshot returns a copy of the latest state and history
func (ms *MonitorState) GetSnapshot() (SchedulerSnapshot, []SchedulerSnapshot) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	h := make([]SchedulerSnapshot, len(ms.history))
	copy(h, ms.history)
	return ms.latest, h
}
