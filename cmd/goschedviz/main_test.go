package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/JustSkiv/goschedviz/internal/domain"
)

func TestConvertToUIData_Simple(t *testing.T) {
	// Prepare test data
	latest := domain.SchedulerSnapshot{
		TimeMs:     1000,
		GoMaxProcs: 2,
		LRQ:        []int{1, 2},
		LRQSum:     3,
	}

	history := []domain.SchedulerSnapshot{latest}

	// Call the function under test
	result := convertToUIData(latest, history)

	// Assert results
	assert.Equal(t, 1000, result.Current.TimeMs)
	assert.Equal(t, 2, result.Current.GoMaxProcs)
	assert.Equal(t, []int{1, 2}, result.Current.LRQ)
	assert.Equal(t, 3, result.Current.LRQSum)
}

func TestConvertToUIData2(t *testing.T) {
	tests := []struct {
		name     string
		latest   domain.SchedulerSnapshot
		history  []domain.SchedulerSnapshot
		expected struct {
			maxGRQ       int
			maxLRQ       int
			maxThreads   int
			maxIdleProcs int
		}
	}{
		{
			name: "empty_history",
			latest: domain.SchedulerSnapshot{
				TimeMs:     1000,
				GoMaxProcs: 4,
				IdleProcs:  2,
				Threads:    8,
				RunQueue:   5,
				LRQ:        []int{1, 2, 1, 0},
				LRQSum:     4,
			},
			history: nil,
			expected: struct {
				maxGRQ       int
				maxLRQ       int
				maxThreads   int
				maxIdleProcs int
			}{
				maxGRQ:       1,
				maxLRQ:       1,
				maxThreads:   1,
				maxIdleProcs: 1,
			},
		},
		{
			name: "single_processor",
			latest: domain.SchedulerSnapshot{
				TimeMs:     1000,
				GoMaxProcs: 1,
				IdleProcs:  0,
				Threads:    2,
				RunQueue:   5,
				LRQ:        []int{3},
				LRQSum:     3,
			},
			history: []domain.SchedulerSnapshot{
				{TimeMs: 0, RunQueue: 2, Threads: 1, LRQ: []int{1}, LRQSum: 1},
				{TimeMs: 500, RunQueue: 3, Threads: 2, LRQ: []int{2}, LRQSum: 2},
				{TimeMs: 1000, RunQueue: 5, Threads: 2, LRQ: []int{3}, LRQSum: 3},
			},
			expected: struct {
				maxGRQ       int
				maxLRQ       int
				maxThreads   int
				maxIdleProcs int
			}{
				maxGRQ:       5,
				maxLRQ:       3,
				maxThreads:   2,
				maxIdleProcs: 1,
			},
		},
		{
			name: "growing_load",
			latest: domain.SchedulerSnapshot{
				TimeMs:     3000,
				GoMaxProcs: 4,
				IdleProcs:  0,
				Threads:    16,
				RunQueue:   15,
				LRQ:        []int{5, 5, 5, 5},
				LRQSum:     20,
			},
			history: []domain.SchedulerSnapshot{
				{TimeMs: 1000, RunQueue: 5, Threads: 8, IdleProcs: 2, LRQSum: 10},
				{TimeMs: 2000, RunQueue: 10, Threads: 12, IdleProcs: 1, LRQSum: 15},
				{TimeMs: 3000, RunQueue: 15, Threads: 16, IdleProcs: 0, LRQSum: 20},
			},
			expected: struct {
				maxGRQ       int
				maxLRQ       int
				maxThreads   int
				maxIdleProcs int
			}{
				maxGRQ:       15,
				maxLRQ:       20,
				maxThreads:   16,
				maxIdleProcs: 2,
			},
		},
		{
			name: "zero_values",
			latest: domain.SchedulerSnapshot{
				TimeMs:     1000,
				GoMaxProcs: 2,
				IdleProcs:  0,
				Threads:    0,
				RunQueue:   0,
				LRQ:        []int{0, 0},
				LRQSum:     0,
			},
			history: []domain.SchedulerSnapshot{
				{TimeMs: 1000, RunQueue: 0, Threads: 0, LRQSum: 0},
			},
			expected: struct {
				maxGRQ       int
				maxLRQ       int
				maxThreads   int
				maxIdleProcs int
			}{
				maxGRQ:       1,
				maxLRQ:       1,
				maxThreads:   1,
				maxIdleProcs: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertToUIData(tt.latest, tt.history)

			// Verify current values
			assert.Equal(t, tt.latest.TimeMs, result.Current.TimeMs)
			assert.Equal(t, tt.latest.GoMaxProcs, result.Current.GoMaxProcs)
			assert.Equal(t, tt.latest.IdleProcs, result.Current.IdleProcs)
			assert.Equal(t, tt.latest.Threads, result.Current.Threads)
			assert.Equal(t, tt.latest.RunQueue, result.Current.RunQueue)
			assert.Equal(t, tt.latest.LRQ, result.Current.LRQ)
			assert.Equal(t, tt.latest.LRQSum, result.Current.LRQSum)
			assert.Equal(t, len(tt.latest.LRQ), result.Current.NumP)

			// Verify max values
			assert.Equal(t, tt.expected.maxGRQ, result.Gauges.GRQ.Max)
			assert.Equal(t, tt.expected.maxLRQ, result.Gauges.LRQ.Max)
			assert.Equal(t, tt.expected.maxThreads, result.Gauges.Threads.Max)
			assert.Equal(t, tt.expected.maxIdleProcs, result.Gauges.IdleProcs.Max)

			// Verify history conversion
			if tt.history != nil {
				assert.Equal(t, len(tt.history), len(result.History))
				for i, h := range tt.history {
					assert.Equal(t, h.TimeMs, result.History[i].TimeMs)
					assert.Equal(t, h.RunQueue, result.History[i].GRQ)
					assert.Equal(t, h.LRQSum, result.History[i].LRQSum)
					assert.Equal(t, h.IdleProcs, result.History[i].IdleProcs)
					assert.Equal(t, h.Threads, result.History[i].Threads)
				}
			} else {
				assert.Empty(t, result.History)
			}
		})
	}
}
