package domain

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMonitorState_Update(t *testing.T) {
	tests := []struct {
		name           string
		initialState   []SchedulerSnapshot // pre-populate history
		update         SchedulerSnapshot   // new snapshot to add
		wantLatest     SchedulerSnapshot   // expected latest snapshot
		wantHistoryLen int                 // expected history length
	}{
		{
			name:         "first update",
			initialState: nil,
			update: SchedulerSnapshot{
				TimeMs:     1000,
				GoMaxProcs: 4,
				RunQueue:   2,
				LRQ:        []int{1, 0, 1, 0},
			},
			wantLatest: SchedulerSnapshot{
				TimeMs:     1000,
				GoMaxProcs: 4,
				RunQueue:   2,
				LRQ:        []int{1, 0, 1, 0},
			},
			wantHistoryLen: 1,
		},
		{
			name: "history within limit",
			initialState: []SchedulerSnapshot{
				{TimeMs: 1000, RunQueue: 1},
				{TimeMs: 2000, RunQueue: 2},
			},
			update: SchedulerSnapshot{
				TimeMs:   3000,
				RunQueue: 3,
			},
			wantLatest: SchedulerSnapshot{
				TimeMs:   3000,
				RunQueue: 3,
			},
			wantHistoryLen: 3,
		},
		{
			name:         "history exceeds limit",
			initialState: make([]SchedulerSnapshot, MaxHistoryPoints), // fill to limit
			update: SchedulerSnapshot{
				TimeMs:   9999,
				RunQueue: 42,
			},
			wantLatest: SchedulerSnapshot{
				TimeMs:   9999,
				RunQueue: 42,
			},
			wantHistoryLen: MaxHistoryPoints,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MonitorState{}

			// Pre-populate history if needed
			for _, s := range tt.initialState {
				ms.Update(s)
			}

			// Perform update
			ms.Update(tt.update)

			// Check latest snapshot
			latest, history := ms.GetSnapshot()
			assert.Equal(t, tt.wantLatest, latest, "latest snapshot mismatch")
			assert.Equal(t, tt.wantHistoryLen, len(history), "history length mismatch")

			// Latest snapshot should be both in .latest and at the end of history
			if len(history) > 0 {
				assert.Equal(t, latest, history[len(history)-1], "latest snapshot should match last history entry")
			}
		})
	}
}

func TestMonitorState_ConcurrentAccess(t *testing.T) {
	ms := &MonitorState{}
	const numGoroutines = 10
	const updatesPerGoroutine = 100

	// Create wait group to sync goroutines
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Launch multiple goroutines to update state concurrently
	for i := 0; i < numGoroutines; i++ {
		go func(routineID int) {
			defer wg.Done()
			for j := 0; j < updatesPerGoroutine; j++ {
				snapshot := SchedulerSnapshot{
					TimeMs:   routineID*1000 + j,
					RunQueue: routineID*100 + j,
				}
				ms.Update(snapshot)
			}
		}(i)
	}

	// Launch goroutine that reads state while updates are happening
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				latest, history := ms.GetSnapshot()
				// Verify that we always get consistent state
				if len(history) > 0 {
					assert.Equal(t, latest, history[len(history)-1],
						"latest snapshot should match last history entry")
				}
				// History should never exceed max size
				assert.LessOrEqual(t, len(history), MaxHistoryPoints,
					"history should not exceed max size")
			}
		}
	}()

	// Wait for all writers to finish
	wg.Wait()
	close(done)

	// Final verification
	latest, history := ms.GetSnapshot()
	require.NotEmpty(t, history, "history should not be empty after updates")
	assert.Equal(t, latest, history[len(history)-1],
		"final latest snapshot should match last history entry")
	assert.LessOrEqual(t, len(history), MaxHistoryPoints,
		"final history should not exceed max size")
}

func TestMonitorState_GetSnapshot(t *testing.T) {
	// Test that GetSnapshot returns a copy, not references to internal state
	ms := &MonitorState{}

	original := SchedulerSnapshot{
		TimeMs:     1000,
		GoMaxProcs: 4,
		LRQ:        []int{1, 2, 3, 4},
	}
	ms.Update(original)

	// Get snapshot and modify it
	latest, history := ms.GetSnapshot()
	latest.TimeMs = 9999
	latest.LRQ[0] = 9999
	history[0].TimeMs = 8888
	history[0].LRQ[0] = 8888

	// Get another snapshot and verify it wasn't affected
	newLatest, newHistory := ms.GetSnapshot()
	assert.Equal(t, original.TimeMs, newLatest.TimeMs, "internal latest TimeMs was modified")
	assert.Equal(t, original.LRQ[0], newLatest.LRQ[0], "internal latest LRQ was modified")
	assert.Equal(t, original.TimeMs, newHistory[0].TimeMs, "internal history TimeMs was modified")
	assert.Equal(t, original.LRQ[0], newHistory[0].LRQ[0], "internal history LRQ was modified")
}
