package widgets

import (
	"testing"

	"github.com/gizak/termui/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JustSkiv/goschedviz/internal/ui"
)

func TestHistoryPlot_New(t *testing.T) {
	plot := NewHistoryPlot()
	require.NotNil(t, plot, "NewHistoryPlot should return non-nil plot")

	assert.Equal(t, "GRQ / LRQ History", plot.Title,
		"Plot should have correct title")

	assert.Equal(t, 4, len(plot.Data),
		"Plot should have two data series (GRQ and LRQ)")
	assert.Equal(t, 4, len(plot.LineColors),
		"Plot should have two line colors")

	assert.Equal(t, termui.ColorGreen, plot.LineColors[0],
		"GRQ line should be green")
	assert.Equal(t, termui.ColorMagenta, plot.LineColors[1],
		"LRQ line should be magenta")
	assert.Equal(t, termui.ColorYellow, plot.LineColors[2],
		"IdleProcs line should be yellow")
	assert.Equal(t, termui.ColorRed, plot.LineColors[3],
		"Threads line should be red")

	for i := 0; i < 4; i++ {
		assert.Equal(t, []float64{0, 0}, plot.Data[i],
			"Initial data should be zero")
	}

}

func TestHistoryPlot_Update(t *testing.T) {
	tests := []struct {
		name     string
		history  []ui.HistoricalValues
		expected struct {
			grqData []float64
			lrqData []float64
		}
	}{
		{
			name:    "empty history",
			history: []ui.HistoricalValues{},
			expected: struct {
				grqData []float64
				lrqData []float64
			}{
				grqData: []float64{0, 0},
				lrqData: []float64{0, 0},
			},
		},
		{
			name: "single data point",
			history: []ui.HistoricalValues{
				{TimeMs: 1000, GRQ: 5, LRQSum: 10},
			},
			expected: struct {
				grqData []float64
				lrqData []float64
			}{
				grqData: []float64{0, 0},
				lrqData: []float64{0, 0},
			},
		},
		{
			name: "multiple data points",
			history: []ui.HistoricalValues{
				{TimeMs: 1000, GRQ: 5, LRQSum: 10},
				{TimeMs: 2000, GRQ: 8, LRQSum: 15},
				{TimeMs: 3000, GRQ: 3, LRQSum: 7},
			},
			expected: struct {
				grqData []float64
				lrqData []float64
			}{
				grqData: []float64{5, 8, 3},
				lrqData: []float64{10, 15, 7},
			},
		},
		{
			name: "zero values in sequence",
			history: []ui.HistoricalValues{
				{TimeMs: 1000, GRQ: 0, LRQSum: 0},
				{TimeMs: 2000, GRQ: 5, LRQSum: 10},
				{TimeMs: 3000, GRQ: 0, LRQSum: 0},
			},
			expected: struct {
				grqData []float64
				lrqData []float64
			}{
				grqData: []float64{0, 5, 0},
				lrqData: []float64{0, 10, 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plot := NewHistoryPlot()
			plot.Update(tt.history)

			assert.Equal(t, tt.expected.grqData, plot.Data[0],
				"GRQ data series should match expected values")
			assert.Equal(t, tt.expected.lrqData, plot.Data[1],
				"LRQ data series should match expected values")
		})
	}
}

func TestHistoryPlot_DataIntegrity(t *testing.T) {
	plot := NewHistoryPlot()

	updates := []struct {
		history []ui.HistoricalValues
		length  int
	}{
		{
			history: []ui.HistoricalValues{
				{TimeMs: 1000, GRQ: 5, LRQSum: 10},
				{TimeMs: 2000, GRQ: 8, LRQSum: 15},
			},
			length: 2,
		},
		{
			history: []ui.HistoricalValues{
				{TimeMs: 1000, GRQ: 5, LRQSum: 10},
				{TimeMs: 2000, GRQ: 8, LRQSum: 15},
				{TimeMs: 3000, GRQ: 3, LRQSum: 7},
			},
			length: 3,
		},
		{
			history: []ui.HistoricalValues{
				{TimeMs: 4000, GRQ: 2, LRQSum: 5},
			},
			length: 2,
		},
	}

	for i, update := range updates {
		t.Run(string(rune('A'+i)), func(t *testing.T) {
			plot.Update(update.history)

			assert.Equal(t, update.length, len(plot.Data[0]),
				"GRQ data series should have correct length")
			assert.Equal(t, update.length, len(plot.Data[1]),
				"LRQ data series should have correct length")

			if len(update.history) >= 2 {
				for j, val := range update.history {
					assert.Equal(t, float64(val.GRQ), plot.Data[0][j],
						"GRQ value at index %d should match input", j)
					assert.Equal(t, float64(val.LRQSum), plot.Data[1][j],
						"LRQ value at index %d should match input", j)
				}
			} else {
				assert.Equal(t, []float64{0, 0}, plot.Data[0],
					"GRQ should have default values when insufficient data")
				assert.Equal(t, []float64{0, 0}, plot.Data[1],
					"LRQ should have default values when insufficient data")
			}
		})
	}
}
