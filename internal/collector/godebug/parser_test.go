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
