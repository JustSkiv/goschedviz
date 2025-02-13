package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/JustSkiv/goschedviz/internal/domain"
)

func TestConvertToUIData(t *testing.T) {
	tests := []struct {
		name     string
		latest   domain.SchedulerSnapshot
		history  []domain.SchedulerSnapshot
		expected struct {
			maxGRQ  int
			maxLRQ  int
			current int
			lrqSum  int
		}
	}{
		{
			name: "empty_history",
			latest: domain.SchedulerSnapshot{
				TimeMs:          1000,
				GoMaxProcs:      4,
				IdleProcs:       2,
				Threads:         8,
				SpinningThreads: 1,
				NeedSpinning:    0,
				IdleThreads:     3,
				RunQueue:        5,
				LRQSum:          10,
				LRQ:             []int{2, 3, 1, 4},
			},
			history: nil,
			expected: struct {
				maxGRQ  int
				maxLRQ  int
				current int
				lrqSum  int
			}{
				maxGRQ:  1,
				maxLRQ:  1,
				current: 5,
				lrqSum:  10,
			},
		},
		{
			name: "growing_load",
			latest: domain.SchedulerSnapshot{
				TimeMs:   3000,
				RunQueue: 15,
				LRQSum:   30,
				LRQ:      []int{10, 10, 10},
			},
			history: []domain.SchedulerSnapshot{
				{
					TimeMs:   1000,
					RunQueue: 5,
					LRQSum:   10,
				},
				{
					TimeMs:   2000,
					RunQueue: 10,
					LRQSum:   20,
				},
				{
					TimeMs:   3000,
					RunQueue: 15,
					LRQSum:   30,
				},
			},
			expected: struct {
				maxGRQ  int
				maxLRQ  int
				current int
				lrqSum  int
			}{
				maxGRQ:  15,
				maxLRQ:  30,
				current: 15,
				lrqSum:  30,
			},
		},
		{
			name: "zero_values",
			latest: domain.SchedulerSnapshot{
				TimeMs:   1000,
				RunQueue: 0,
				LRQSum:   0,
				LRQ:      []int{0, 0},
			},
			history: []domain.SchedulerSnapshot{
				{
					TimeMs:   1000,
					RunQueue: 0,
					LRQSum:   0,
				},
			},
			expected: struct {
				maxGRQ  int
				maxLRQ  int
				current int
				lrqSum  int
			}{
				maxGRQ:  1,
				maxLRQ:  1,
				current: 0,
				lrqSum:  0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertToUIData(tt.latest, tt.history)

			assert.Equal(t, tt.latest.TimeMs, result.Current.TimeMs)
			assert.Equal(t, tt.latest.GoMaxProcs, result.Current.GoMaxProcs)
			assert.Equal(t, tt.latest.IdleProcs, result.Current.IdleProcs)
			assert.Equal(t, tt.latest.Threads, result.Current.Threads)
			assert.Equal(t, tt.latest.SpinningThreads, result.Current.SpinningThreads)
			assert.Equal(t, tt.latest.NeedSpinning, result.Current.NeedSpinning)
			assert.Equal(t, tt.latest.IdleThreads, result.Current.IdleThreads)
			assert.Equal(t, tt.latest.RunQueue, result.Current.RunQueue)
			assert.Equal(t, tt.latest.LRQSum, result.Current.LRQSum)
			assert.Equal(t, len(tt.latest.LRQ), result.Current.NumP)
			assert.Equal(t, tt.latest.LRQ, result.Current.LRQ)

			assert.Equal(t, tt.expected.current, result.Gauges.GRQ.Current)
			assert.Equal(t, tt.expected.maxGRQ, result.Gauges.GRQ.Max)
			assert.Equal(t, tt.expected.lrqSum, result.Gauges.LRQ.Current)
			assert.Equal(t, tt.expected.maxLRQ, result.Gauges.LRQ.Max)

			assert.Equal(t, len(tt.history), len(result.History))
			for i := range tt.history {
				assert.Equal(t, tt.history[i].TimeMs, result.History[i].TimeMs)
				assert.Equal(t, tt.history[i].RunQueue, result.History[i].GRQ)
				assert.Equal(t, tt.history[i].LRQSum, result.History[i].LRQSum)
			}
		})
	}
}
