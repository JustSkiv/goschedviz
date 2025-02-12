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

	// Check initial configuration
	assert.Equal(t, "GRQ / LRQ History", plot.Title,
		"Plot should have correct title")

	// Verify initial data structure
	assert.Equal(t, 2, len(plot.Data),
		"Plot should have two data series (GRQ and LRQ)")
	assert.Equal(t, 2, len(plot.LineColors),
		"Plot should have two line colors")

	// Check color configuration
	assert.Equal(t, termui.ColorGreen, plot.LineColors[0],
		"GRQ line should be green")
	assert.Equal(t, termui.ColorMagenta, plot.LineColors[1],
		"LRQ line should be magenta")

	// Verify initial data points
	assert.Equal(t, []float64{0, 0}, plot.Data[0],
		"Initial GRQ data should be zero")
	assert.Equal(t, []float64{0, 0}, plot.Data[1],
		"Initial LRQ data should be zero")
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
				grqData: []float64{0, 0}, // Initial values should remain
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
				grqData: []float64{5},
				lrqData: []float64{10},
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

			if len(tt.history) == 0 {
				// For empty history, check that initial values are preserved
				assert.Equal(t, tt.expected.grqData, plot.Data[0],
					"Empty history should preserve initial GRQ values")
				assert.Equal(t, tt.expected.lrqData, plot.Data[1],
					"Empty history should preserve initial LRQ values")
			} else {
				// For non-empty history, check the actual data
				assert.Equal(t, tt.expected.grqData, plot.Data[0],
					"GRQ data series should match expected values")
				assert.Equal(t, tt.expected.lrqData, plot.Data[1],
					"LRQ data series should match expected values")
			}
		})
	}
}

func TestHistoryPlot_DataIntegrity(t *testing.T) {
	plot := NewHistoryPlot()

	// Test data integrity with sequential updates
	updates := []struct {
		history []ui.HistoricalValues
		length  int
	}{
		{
			// First update with some data
			history: []ui.HistoricalValues{
				{TimeMs: 1000, GRQ: 5, LRQSum: 10},
				{TimeMs: 2000, GRQ: 8, LRQSum: 15},
			},
			length: 2,
		},
		{
			// Second update with more data
			history: []ui.HistoricalValues{
				{TimeMs: 1000, GRQ: 5, LRQSum: 10},
				{TimeMs: 2000, GRQ: 8, LRQSum: 15},
				{TimeMs: 3000, GRQ: 3, LRQSum: 7},
			},
			length: 3,
		},
		{
			// Third update with less data
			history: []ui.HistoricalValues{
				{TimeMs: 4000, GRQ: 2, LRQSum: 5},
			},
			length: 1,
		},
	}

	for i, update := range updates {
		t.Run(string(rune('A'+i)), func(t *testing.T) {
			plot.Update(update.history)

			// Check data series lengths
			assert.Equal(t, update.length, len(plot.Data[0]),
				"GRQ data series should have correct length")
			assert.Equal(t, update.length, len(plot.Data[1]),
				"LRQ data series should have correct length")

			// Verify data matches input
			for j, val := range update.history {
				assert.Equal(t, float64(val.GRQ), plot.Data[0][j],
					"GRQ value at index %d should match input", j)
				assert.Equal(t, float64(val.LRQSum), plot.Data[1][j],
					"LRQ value at index %d should match input", j)
			}
		})
	}
}
