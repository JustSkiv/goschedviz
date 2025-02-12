package widgets

import (
	"fmt"
	"testing"

	"github.com/gizak/termui/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLRQBarChart_New(t *testing.T) {
	chart := NewLRQBarChart()
	require.NotNil(t, chart, "NewLRQBarChart should return non-nil chart")

	// Check initial configuration
	assert.Equal(t, "Local Run Queues (per P)", chart.Title,
		"Bar chart should have correct title")
	assert.Equal(t, 3, chart.BarWidth,
		"Bar chart should have correct bar width")
	assert.Equal(t, 1, chart.BarGap,
		"Bar chart should have correct gap between bars")

	// Check color configuration
	assert.Equal(t, []termui.Color{termui.ColorCyan}, chart.BarColors,
		"Bar chart should have correct bar color")
	assert.Equal(t, []termui.Style{termui.NewStyle(termui.ColorYellow)}, chart.LabelStyles,
		"Bar chart should have correct label style")

	// Test number formatter
	formatter := chart.NumFormatter
	require.NotNil(t, formatter, "Number formatter should not be nil")
	assert.Equal(t, "42", formatter(42.0),
		"Number formatter should format integers correctly")
	assert.Equal(t, "0", formatter(0.0),
		"Number formatter should handle zero correctly")
}

func TestLRQBarChart_Update(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected struct {
			data   []float64
			labels []string
		}
	}{
		{
			name:  "empty input",
			input: []int{},
			expected: struct {
				data   []float64
				labels []string
			}{
				data:   nil, // Changed from empty slice to nil
				labels: nil, // Changed from empty slice to nil
			},
		},
		{
			name:  "single processor",
			input: []int{5},
			expected: struct {
				data   []float64
				labels []string
			}{
				data:   []float64{5},
				labels: []string{"P0"},
			},
		},
		{
			name:  "multiple processors",
			input: []int{1, 2, 3, 4},
			expected: struct {
				data   []float64
				labels []string
			}{
				data:   []float64{1, 2, 3, 4},
				labels: []string{"P0", "P1", "P2", "P3"},
			},
		},
		{
			name:  "zero values",
			input: []int{0, 0, 0},
			expected: struct {
				data   []float64
				labels []string
			}{
				data:   []float64{0, 0, 0},
				labels: []string{"P0", "P1", "P2"},
			},
		},
		{
			name:  "mixed values",
			input: []int{0, 5, 0, 10, 0},
			expected: struct {
				data   []float64
				labels []string
			}{
				data:   []float64{0, 5, 0, 10, 0},
				labels: []string{"P0", "P1", "P2", "P3", "P4"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chart := NewLRQBarChart()
			chart.Update(tt.input)

			if len(tt.input) == 0 {
				// For empty input, check that both slices are nil
				assert.Nil(t, chart.Data, "Bar chart data should be nil for empty input")
				assert.Nil(t, chart.Labels, "Bar chart labels should be nil for empty input")
			} else {
				// Check data values
				assert.Equal(t, tt.expected.data, chart.Data,
					"Bar chart data should match expected values")

				// Check labels
				assert.Equal(t, tt.expected.labels, chart.Labels,
					"Bar chart labels should match expected values")

				// Check data and labels length match
				assert.Equal(t, len(chart.Data), len(chart.Labels),
					"Data and labels should have same length")
			}
		})
	}
}

func TestLRQBarChart_LargeDataSet(t *testing.T) {
	// Test with a larger number of processors
	numProcs := 32
	input := make([]int, numProcs)
	expectedData := make([]float64, numProcs)
	expectedLabels := make([]string, numProcs)

	// Fill test data
	for i := 0; i < numProcs; i++ {
		input[i] = i * 2 // Some arbitrary pattern
		expectedData[i] = float64(i * 2)
		expectedLabels[i] = fmt.Sprintf("P%d", i)
	}

	chart := NewLRQBarChart()
	chart.Update(input)

	// Verify all data points
	assert.Equal(t, expectedData, chart.Data,
		"Bar chart should handle large data sets correctly")
	assert.Equal(t, expectedLabels, chart.Labels,
		"Bar chart should generate correct labels for large data sets")
}

func TestLRQBarChart_UpdateConsistency(t *testing.T) {
	chart := NewLRQBarChart()

	// Test multiple updates
	updates := [][]int{
		{1, 2, 3},        // First update
		{4, 5, 6, 7},     // More processors
		{1},              // Fewer processors
		{0, 0, 0},        // All zeros
		{10, 20, 30, 40}, // Different values
	}

	for i, update := range updates {
		t.Run(fmt.Sprintf("update_%d", i), func(t *testing.T) {
			chart.Update(update)

			// Check consistency after each update
			assert.Equal(t, len(update), len(chart.Data),
				"Data length should match input length")
			assert.Equal(t, len(update), len(chart.Labels),
				"Labels length should match input length")

			// Verify data conversion
			for j, val := range update {
				assert.Equal(t, float64(val), chart.Data[j],
					"Data conversion should be correct at index %d", j)
				assert.Equal(t, fmt.Sprintf("P%d", j), chart.Labels[j],
					"Label should be correct at index %d", j)
			}
		})
	}
}

func TestLRQBarChart_NumberFormatter(t *testing.T) {
	chart := NewLRQBarChart()
	formatter := chart.NumFormatter

	tests := []struct {
		input    float64
		expected string
	}{
		{0.0, "0"},
		{42.0, "42"},
		{99.9, "100"}, // Should round to nearest integer
		{-1.0, "-1"},  // Handle negative numbers (though unlikely in practice)
		{1000.0, "1000"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("format_%.1f", tt.input), func(t *testing.T) {
			result := formatter(tt.input)
			assert.Equal(t, tt.expected, result,
				"Number formatter should correctly format %.1f", tt.input)
		})
	}
}
