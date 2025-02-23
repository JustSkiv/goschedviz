package widgets

import (
	"math"
	"testing"

	tui "github.com/gizak/termui/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JustSkiv/goschedviz/internal/ui"
)

func TestLinearHistoryPlot_New(t *testing.T) {
	plot := NewLinearHistoryPlot()
	require.NotNil(t, plot, "NewLinearHistoryPlot should return non-nil plot")

	assert.Equal(t, "History Plot (linear)", plot.Title,
		"Plot should have correct title")

	assert.Equal(t, 5, len(plot.Data),
		"Plot should have five data series")
	assert.Equal(t, 5, len(plot.LineColors),
		"Plot should have five line colors")

	assert.Equal(t, tui.ColorGreen, plot.LineColors[0], "GRQ line should be green")
	assert.Equal(t, tui.ColorMagenta, plot.LineColors[1], "LRQ line should be magenta")
	assert.Equal(t, tui.ColorRed, plot.LineColors[2], "Threads line should be red")
	assert.Equal(t, tui.ColorYellow, plot.LineColors[3], "IdleProcs line should be yellow")
	assert.Equal(t, tui.ColorCyan, plot.LineColors[4], "Goroutines line should be cyan")

	// Check initial data
	for i := 0; i < 5; i++ {
		assert.Equal(t, []float64{0, 0}, plot.Data[i], "Initial data should be zero")
	}
}

func TestLogHistoryPlot_New(t *testing.T) {
	plot := NewLogHistoryPlot()
	require.NotNil(t, plot, "NewLogHistoryPlot should return non-nil plot")

	assert.Equal(t, "History Plot (log)", plot.Title,
		"Plot should have correct title")

	assert.Equal(t, 5, len(plot.Data),
		"Plot should have five data series")
	assert.Equal(t, 5, len(plot.LineColors),
		"Plot should have five line colors")

	// Check colors same as linear plot
	assert.Equal(t, tui.ColorGreen, plot.LineColors[0], "GRQ line should be green")
	assert.Equal(t, tui.ColorMagenta, plot.LineColors[1], "LRQ line should be magenta")
	assert.Equal(t, tui.ColorRed, plot.LineColors[2], "Threads line should be red")
	assert.Equal(t, tui.ColorYellow, plot.LineColors[3], "IdleProcs line should be yellow")
	assert.Equal(t, tui.ColorCyan, plot.LineColors[4], "Goroutines line should be cyan")
}

func TestLogHistoryPlot_ScaleConversion(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  float64
	}{
		{"zero value", 0, 0},
		{"negative value", -1, 0},
		{"small value", 1, 0},    // log10(1) = 0
		{"large value", 1000, 3}, // log10(1000) = 3
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toLogScale(tt.input)
			assert.InDelta(t, tt.want, got, 0.0001,
				"toLogScale(%v) should be close to %v", tt.input, tt.want)
		})
	}
}

func TestPlot_Update(t *testing.T) {
	tests := []struct {
		name    string
		history []ui.HistoricalValues
		check   func(t *testing.T, linear *LinearHistoryPlot, log *LogHistoryPlot)
	}{
		{
			name:    "empty history",
			history: []ui.HistoricalValues{},
			check: func(t *testing.T, linear *LinearHistoryPlot, log *LogHistoryPlot) {
				// Both plots should have default values
				for i := 0; i < 5; i++ {
					assert.Equal(t, []float64{0, 0}, linear.Data[i])
					assert.Equal(t, []float64{0, 0}, log.Data[i])
				}
			},
		},
		{
			name: "single point",
			history: []ui.HistoricalValues{
				{TimeMs: 100, GRQ: 5, LRQSum: 10, Threads: 20, IdleProcs: 2, Goroutines: 100},
			},
			check: func(t *testing.T, linear *LinearHistoryPlot, log *LogHistoryPlot) {
				// Both plots should have default values with single point
				for i := 0; i < 5; i++ {
					assert.Equal(t, []float64{0, 0}, linear.Data[i])
					assert.Equal(t, []float64{0, 0}, log.Data[i])
				}
			},
		},
		{
			name: "multiple points",
			history: []ui.HistoricalValues{
				{TimeMs: 100, GRQ: 5, LRQSum: 10, Threads: 20, IdleProcs: 2, Goroutines: 100},
				{TimeMs: 200, GRQ: 8, LRQSum: 15, Threads: 25, IdleProcs: 1, Goroutines: 200},
				{TimeMs: 300, GRQ: 3, LRQSum: 7, Threads: 15, IdleProcs: 3, Goroutines: 150},
			},
			check: func(t *testing.T, linear *LinearHistoryPlot, log *LogHistoryPlot) {
				// Check linear values
				assert.Equal(t, []float64{5, 8, 3}, linear.Data[0], "Linear GRQ values")
				assert.Equal(t, []float64{10, 15, 7}, linear.Data[1], "Linear LRQ values")
				assert.Equal(t, []float64{20, 25, 15}, linear.Data[2], "Linear Thread values")
				assert.Equal(t, []float64{2, 1, 3}, linear.Data[3], "Linear IdleProcs values")
				assert.Equal(t, []float64{100, 200, 150}, linear.Data[4], "Linear Goroutine values")

				// Check log values (approximate due to floating point)
				for i, vals := range log.Data[0] {
					assert.InDelta(t, math.Log10(float64([]int{5, 8, 3}[i])), vals, 0.0001)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			linear := NewLinearHistoryPlot()
			log := NewLogHistoryPlot()

			linear.Update(tt.history)
			log.Update(tt.history)

			tt.check(t, linear, log)
		})
	}
}
