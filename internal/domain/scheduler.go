// Package domain defines core types and interfaces for the scheduler monitoring
package domain

// SchedData contains parsed values from a single "SCHED ..." trace line.
// This data represents a snapshot of Go scheduler state at a specific moment.
type SchedData struct {
	TimeMs          int   // Time since start in milliseconds
	GoMaxProcs      int   // Number of processors (P)
	IdleProcs       int   // Number of idle processors
	Threads         int   // Total number of threads
	SpinningThreads int   // Number of spinning threads
	NeedSpinning    int   // Required number of spinning threads
	IdleThreads     int   // Number of idle threads
	RunQueue        int   // Global Run Queue (GRQ)
	LrqSum          int   // Sum of all Local Run Queues
	Lrq             []int // Individual Local Run Queues per P
}

// Event represents UI events that can occur during monitoring
type Event string

const (
	// EventQuit signals application shutdown
	EventQuit Event = "quit"
	// EventResize signals terminal resize
	EventResize Event = "resize"
)
