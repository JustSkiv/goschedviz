package ui

// UIData represents the data structure passed to UI for visualization.
type UIData struct {
	// Current values
	Current CurrentValues

	// Historical values for plotting
	History []HistoricalValues

	// Gauge values
	Gauges GaugeValues
}

// CurrentValues contains the latest scheduler metrics.
type CurrentValues struct {
	TimeMs          int
	GoMaxProcs      int
	IdleProcs       int
	Threads         int
	SpinningThreads int
	NeedSpinning    int
	IdleThreads     int
	RunQueue        int
	LRQSum          int
	NumP            int   // Number of P (processors)
	LRQ             []int // Local run queues by P
}

// HistoricalValues contains metrics used for plotting history.
type HistoricalValues struct {
	TimeMs    int
	GRQ       int
	LRQSum    int
	IdleProcs int
	Threads   int
}

// GaugeValues contains data for all gauges
type GaugeValues struct {
	GRQ struct {
		Current int
		Max     int
	}
	LRQ struct {
		Current int
		Max     int
	}
	IdleProcs struct {
		Current int
		Max     int
	}
	Threads struct {
		Current int
		Max     int
	}
}
