package godebug

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		want     bool // whether parsing should succeed
		expected struct {
			timeMs          int
			gomaxprocs      int
			idleprocs       int
			threads         int
			spinningThreads int
			needSpinning    int
			idleThreads     int
			runQueue        int
			lrqSum          int
			lrq             []int
		}
	}{
		{
			name:  "empty input",
			input: "",
			want:  false,
		},
		{
			name:  "invalid prefix",
			input: "SCHEDx 2013ms: gomaxprocs=14 idleprocs=14 threads=22 spinningthreads=0 needspinning=0 idlethreads=17 runqueue=0 [0 0 0 0]",
			want:  false,
		},
		{
			name:  "malformed input - missing brackets",
			input: "SCHED 2013ms: gomaxprocs=14 idleprocs=14 threads=22 spinningthreads=0 needspinning=0 idlethreads=17 runqueue=0 0 0 0 0",
			want:  false,
		},
		{
			name:  "malformed input - invalid number",
			input: "SCHED 2013ms: gomaxprocs=14 idleprocs=bad threads=22 spinningthreads=0 needspinning=0 idlethreads=17 runqueue=0 [0 0 0 0]",
			want:  false,
		},
		{
			name:  "normal case with zero values",
			input: "SCHED 2013ms: gomaxprocs=14 idleprocs=14 threads=22 spinningthreads=0 needspinning=0 idlethreads=17 runqueue=0 [0 0 0 0 0 0 0 0 0 0 0 0 0 0]",
			want:  true,
			expected: struct {
				timeMs          int
				gomaxprocs      int
				idleprocs       int
				threads         int
				spinningThreads int
				needSpinning    int
				idleThreads     int
				runQueue        int
				lrqSum          int
				lrq             []int
			}{
				timeMs:          2013,
				gomaxprocs:      14,
				idleprocs:       14,
				threads:         22,
				spinningThreads: 0,
				needSpinning:    0,
				idleThreads:     17,
				runQueue:        0,
				lrqSum:          0,
				lrq:             []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
		{
			name:  "normal case with non-zero values",
			input: "SCHED 5000ms: gomaxprocs=8 idleprocs=6 threads=12 spinningthreads=1 needspinning=1 idlethreads=4 runqueue=5 [2 1 0 3 0 1 2 0]",
			want:  true,
			expected: struct {
				timeMs          int
				gomaxprocs      int
				idleprocs       int
				threads         int
				spinningThreads int
				needSpinning    int
				idleThreads     int
				runQueue        int
				lrqSum          int
				lrq             []int
			}{
				timeMs:          5000,
				gomaxprocs:      8,
				idleprocs:       6,
				threads:         12,
				spinningThreads: 1,
				needSpinning:    1,
				idleThreads:     4,
				runQueue:        5,
				lrqSum:          9,
				lrq:             []int{2, 1, 0, 3, 0, 1, 2, 0},
			},
		},
	}

	parser := NewParser()
	require.NotNil(t, parser, "NewParser() returned nil")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing input: %q", tt.input)
			snapshot, ok := parser.Parse(tt.input)

			assert.Equal(t, tt.want, ok, "Parse() success status")

			if !tt.want {
				return // don't check values for tests that should fail
			}

			assert.Equal(t, tt.expected.timeMs, snapshot.TimeMs, "TimeMs mismatch")
			assert.Equal(t, tt.expected.gomaxprocs, snapshot.GoMaxProcs, "GoMaxProcs mismatch")
			assert.Equal(t, tt.expected.idleprocs, snapshot.IdleProcs, "IdleProcs mismatch")
			assert.Equal(t, tt.expected.threads, snapshot.Threads, "Threads mismatch")
			assert.Equal(t, tt.expected.spinningThreads, snapshot.SpinningThreads, "SpinningThreads mismatch")
			assert.Equal(t, tt.expected.needSpinning, snapshot.NeedSpinning, "NeedSpinning mismatch")
			assert.Equal(t, tt.expected.idleThreads, snapshot.IdleThreads, "IdleThreads mismatch")
			assert.Equal(t, tt.expected.runQueue, snapshot.RunQueue, "RunQueue mismatch")
			assert.Equal(t, tt.expected.lrqSum, snapshot.LRQSum, "LRQSum mismatch")
			assert.Equal(t, tt.expected.lrq, snapshot.LRQ, "LRQ mismatch")
		})
	}
}

func TestParser_ParseMetrics(t *testing.T) {
	parser := NewParser()
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "valid metrics line",
			input:    "PROCMETR num_goroutines=1234",
			expected: 1234,
		},
		{
			name:     "invalid prefix",
			input:    "WRONGPREFIX num_goroutines=1234",
			expected: -1,
		},
		{
			name:     "invalid format",
			input:    "PROCMETR bad_format",
			expected: -1,
		},
		{
			name:     "invalid number",
			input:    "PROCMETR num_goroutines=abc",
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.parseMetrics(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParser_GoroutinesUpdate(t *testing.T) {
	parser := NewParser()

	// First parse metrics line
	snapshot, ok := parser.Parse("PROCMETR num_goroutines=1234")
	assert.False(t, ok, "Metrics line should not produce snapshot")

	// Then parse sched line to get snapshot with goroutines count
	schedLine := "SCHED 1000ms: gomaxprocs=4 idleprocs=2 threads=8 spinningthreads=1 needspinning=0 idlethreads=3 runqueue=5 [1 2 1 0]"
	snapshot, ok = parser.Parse(schedLine)
	assert.True(t, ok, "Should parse sched line")
	assert.Equal(t, 1234, snapshot.Goroutines, "Should include goroutines count")

	// Test that goroutines count persists between sched lines
	snapshot, ok = parser.Parse(schedLine)
	assert.True(t, ok, "Should parse sched line")
	assert.Equal(t, 1234, snapshot.Goroutines, "Should maintain goroutines count")

	// Test updating goroutines count
	snapshot, ok = parser.Parse("PROCMETR num_goroutines=5678")
	assert.False(t, ok, "Metrics line should not produce snapshot")

	snapshot, ok = parser.Parse(schedLine)
	assert.True(t, ok, "Should parse sched line")
	assert.Equal(t, 5678, snapshot.Goroutines, "Should update goroutines count")
}
