package termui

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JustSkiv/goschedviz/internal/ui"
)

func TestTermUI_New(t *testing.T) {
	term := New()
	require.NotNil(t, term, "New should return non-nil UI")
	assert.NotNil(t, term.done, "Done channel should be initialized")
}

func TestTermUI_DoneChannel(t *testing.T) {
	term := New()
	done := term.Done()
	require.NotNil(t, done, "Done() should return non-nil channel")

	select {
	case <-done:
		t.Fatal("Done channel should not be closed initially")
	default:
		// This is expected - channel should be open
	}

	close(term.done)

	select {
	case <-done:
		// This is expected - channel should be closed
	default:
		t.Fatal("Done channel should be closed after closing internal channel")
	}
}

func TestTermUI_StartStop(t *testing.T) {
	term := New()

	// Test Start
	err := term.Start()
	if err != nil {
		t.Skipf("Skipping test in non-terminal environment: %v", err)
		return
	}

	// Ensure cleanup
	defer term.Stop()

	// Check that widgets are initialized
	require.NotNil(t, term.table, "Table widget should be initialized")
	require.NotNil(t, term.barChart, "Bar chart widget should be initialized")
	require.NotNil(t, term.grqGauge, "GRQ gauge should be initialized")
	require.NotNil(t, term.lrqGauge, "LRQ gauge should be initialized")
	require.NotNil(t, term.plot, "Plot widget should be initialized")
	require.NotNil(t, term.info, "Info widget should be initialized")
	require.NotNil(t, term.grid, "Grid should be initialized")
}

func TestTermUI_Update(t *testing.T) {
	term := New()
	err := term.Start()
	if err != nil {
		t.Skipf("Skipping test in non-terminal environment: %v", err)
		return
	}
	defer term.Stop()

	// Test various data scenarios
	tests := []struct {
		name string
		data ui.UIData
	}{
		{
			name: "empty data",
			data: ui.UIData{
				Gauges: ui.GaugeValues{
					GRQ: struct{ Current, Max int }{0, 1}, // Установили Max = 1 вместо 0
					LRQ: struct{ Current, Max int }{0, 1}, // Установили Max = 1 вместо 0
				},
			},
		},
		{
			name: "normal load",
			data: ui.UIData{
				Current: ui.CurrentValues{
					TimeMs:          1000,
					GoMaxProcs:      4,
					IdleProcs:       2,
					Threads:         8,
					SpinningThreads: 1,
					NeedSpinning:    0,
					IdleThreads:     3,
					RunQueue:        5,
					LRQSum:          10,
					NumP:            4,
					LRQ:             []int{2, 3, 1, 4},
				},
				History: []ui.HistoricalValues{
					{TimeMs: 0, GRQ: 0, LRQSum: 0},
					{TimeMs: 500, GRQ: 2, LRQSum: 5},
					{TimeMs: 1000, GRQ: 5, LRQSum: 10},
				},
				Gauges: ui.GaugeValues{
					GRQ: struct{ Current, Max int }{5, 10},
					LRQ: struct{ Current, Max int }{10, 20},
				},
			},
		},
		{
			name: "high load",
			data: ui.UIData{
				Current: ui.CurrentValues{
					TimeMs:          5000,
					GoMaxProcs:      32,
					IdleProcs:       0,
					Threads:         64,
					SpinningThreads: 8,
					NeedSpinning:    4,
					IdleThreads:     0,
					RunQueue:        100,
					LRQSum:          500,
					NumP:            32,
					LRQ:             make([]int, 32),
				},
				History: []ui.HistoricalValues{
					{TimeMs: 4800, GRQ: 90, LRQSum: 450},
					{TimeMs: 4900, GRQ: 95, LRQSum: 480},
					{TimeMs: 5000, GRQ: 100, LRQSum: 500},
				},
				Gauges: ui.GaugeValues{
					GRQ: struct{ Current, Max int }{100, 200},
					LRQ: struct{ Current, Max int }{500, 1000},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				term.Update(tt.data)
			}, "Update should not panic with %s", tt.name)
		})
	}
}

func TestTermUI_GridLayout(t *testing.T) {
	term := New()
	err := term.Start()
	if err != nil {
		t.Skipf("Skipping test in non-terminal environment: %v", err)
		return
	}
	defer term.Stop()

	// Check that grid is initialized
	require.NotNil(t, term.grid, "Grid should be initialized")

	// Since we can't directly access grid internals, we'll just verify
	// that the grid exists and basic UI elements are present
	assert.NotNil(t, term.table, "Table widget should be present in grid")
	assert.NotNil(t, term.barChart, "Bar chart widget should be present in grid")
	assert.NotNil(t, term.grqGauge, "GRQ gauge should be present in grid")
	assert.NotNil(t, term.lrqGauge, "LRQ gauge should be present in grid")
	assert.NotNil(t, term.plot, "Plot widget should be present in grid")
	assert.NotNil(t, term.info, "Info widget should be present in grid")
}

func TestTermUI_LargeDataSet(t *testing.T) {
	term := New()
	err := term.Start()
	if err != nil {
		t.Skipf("Skipping test in non-terminal environment: %v", err)
		return
	}
	defer term.Stop()

	// Create large test data
	largeData := ui.UIData{
		Current: ui.CurrentValues{
			GoMaxProcs: 32,
			NumP:       32,
			LRQ:        make([]int, 32),
		},
		Gauges: ui.GaugeValues{ // Добавили значения для gauge
			GRQ: struct{ Current, Max int }{0, 100},
			LRQ: struct{ Current, Max int }{0, 100},
		},
	}

	// Fill history with test data
	largeData.History = make([]ui.HistoricalValues, 1000)
	for i := range largeData.History {
		largeData.History[i] = ui.HistoricalValues{
			TimeMs: i * 100,
			GRQ:    i % 50,
			LRQSum: i % 100,
		}
	}

	assert.NotPanics(t, func() {
		term.Update(largeData)
	}, "Update with large data set should not panic")
}
