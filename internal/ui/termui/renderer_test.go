package termui

import (
	"testing"
	"time"

	"github.com/gizak/termui/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JustSkiv/goschedviz/internal/ui"
)

func TestTermUI_New(t *testing.T) {
	term := New()
	require.NotNil(t, term, "New should return non-nil UI")
	assert.NotNil(t, term.done, "Done channel should be initialized")
}

func TestTermUI_StartStop(t *testing.T) {
	term := newWithTerminal(newTestTerminal())

	// Test Start
	err := term.Start()
	require.NoError(t, err)

	// Ensure cleanup
	defer term.Stop()

	// Check that widgets are initialized
	require.NotNil(t, term.table, "Table widget should be initialized")
	require.NotNil(t, term.barChart, "Bar chart widget should be initialized")
	require.NotNil(t, term.grqGauge, "GRQ gauge should be initialized")
	require.NotNil(t, term.goroutinesGauge, "Goroutines gauge should be initialized")
	require.NotNil(t, term.threadsGauge, "Threads gauge should be initialized")
	require.NotNil(t, term.idleProcsGauge, "IdleProcs gauge should be initialized")
	require.NotNil(t, term.linearPlot, "Linear plot widget should be initialized")
	require.NotNil(t, term.logPlot, "Log plot widget should be initialized")
	require.NotNil(t, term.legend, "Legend widget should be initialized")
	require.NotNil(t, term.info, "Info widget should be initialized")
	require.NotNil(t, term.grid, "Grid should be initialized")
}

func TestTermUI_Update(t *testing.T) {
	term := newWithTerminal(newTestTerminal())
	err := term.Start()
	require.NoError(t, err)
	defer term.Stop()

	tests := []struct {
		name string
		data ui.UIData
	}{
		{
			name: "empty data",
			data: ui.UIData{
				Gauges: ui.GaugeValues{
					GRQ:        struct{ Current, Max int }{0, 1},
					Goroutines: struct{ Current, Max int }{0, 1},
					Threads:    struct{ Current, Max int }{0, 1},
					IdleProcs:  struct{ Current, Max int }{0, 1},
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
					Goroutines:      100,
				},
				History: struct {
					Raw    []ui.HistoricalValues
					Scaled []ui.HistoricalValues
				}{
					Raw: []ui.HistoricalValues{
						{TimeMs: 0, GRQ: 0, LRQSum: 0, Threads: 0, IdleProcs: 0, Goroutines: 0},
						{TimeMs: 500, GRQ: 2, LRQSum: 5, Threads: 4, IdleProcs: 1, Goroutines: 50},
						{TimeMs: 1000, GRQ: 5, LRQSum: 10, Threads: 8, IdleProcs: 2, Goroutines: 100},
					},
				},
				Gauges: ui.GaugeValues{
					GRQ:        struct{ Current, Max int }{5, 10},
					Goroutines: struct{ Current, Max int }{100, 200},
					Threads:    struct{ Current, Max int }{8, 16},
					IdleProcs:  struct{ Current, Max int }{2, 4},
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

func TestTermUI_Events(t *testing.T) {
	mock := newTestTerminal()
	term := newWithTerminal(mock)

	err := term.Start()
	require.NoError(t, err)

	// Test quit event
	done := make(chan struct{})
	go func() {
		mock.SendEvent(termui.Event{
			ID: "q",
		})
		close(done)
	}()

	select {
	case <-term.Done():
		// Expected - UI should close
	case <-time.After(time.Second):
		t.Fatal("Quit event was not processed")
	}

	<-done
	term.Stop()
}

func TestTermUI_ResizeEvent(t *testing.T) {
	mock := newTestTerminal()
	term := newWithTerminal(mock)

	err := term.Start()
	require.NoError(t, err)

	// Test resize event
	done := make(chan struct{})
	go func() {
		mock.SendEvent(termui.Event{
			ID: "<Resize>",
			Payload: termui.Resize{
				Width:  120,
				Height: 50,
			},
		})
		close(done)
	}()

	// Give time for event processing
	time.Sleep(100 * time.Millisecond)
	<-done
	term.Stop()
}
