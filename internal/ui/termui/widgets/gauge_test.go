package widgets

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGRQGauge_New(t *testing.T) {
	// Test gauge initialization
	gauge := NewGRQGauge()
	require.NotNil(t, gauge, "NewGRQGauge should return non-nil gauge")
	assert.Equal(t, "GRQ", gauge.Title, "Gauge should have correct title")
}

func TestGRQGauge_Update(t *testing.T) {
	tests := []struct {
		name     string
		input    struct{ Current, Max int }
		expected struct {
			percent int
			label   string
		}
	}{
		{
			name:  "zero values",
			input: struct{ Current, Max int }{0, 1},
			expected: struct {
				percent int
				label   string
			}{0, "0 / 1"},
		},
		{
			name:  "half full",
			input: struct{ Current, Max int }{50, 100},
			expected: struct {
				percent int
				label   string
			}{50, "50 / 100"},
		},
		{
			name:  "full gauge",
			input: struct{ Current, Max int }{100, 100},
			expected: struct {
				percent int
				label   string
			}{100, "100 / 100"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gauge := NewGRQGauge()
			gauge.Update(tt.input)

			assert.Equal(t, tt.expected.percent, gauge.Percent,
				"Gauge percent should match expected value")
			assert.Equal(t, tt.expected.label, gauge.Label,
				"Gauge label should match expected value")
		})
	}
}

func TestLRQGauge_New(t *testing.T) {
	gauge := NewLRQGauge()
	require.NotNil(t, gauge, "NewLRQGauge should return non-nil gauge")
	assert.Equal(t, "LRQ (sum)", gauge.Title, "Gauge should have correct title")
}

func TestLRQGauge_Update(t *testing.T) {
	tests := []struct {
		name     string
		input    struct{ Current, Max int }
		expected struct {
			percent int
			label   string
		}
	}{
		{
			name:  "empty queue",
			input: struct{ Current, Max int }{0, 100},
			expected: struct {
				percent int
				label   string
			}{0, "0 / 100"},
		},
		{
			name:  "partial load",
			input: struct{ Current, Max int }{75, 150},
			expected: struct {
				percent int
				label   string
			}{50, "75 / 150"},
		},
		{
			name:  "overload simulation",
			input: struct{ Current, Max int }{200, 100},
			expected: struct {
				percent int
				label   string
			}{200, "200 / 100"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gauge := NewLRQGauge()
			gauge.Update(tt.input)

			assert.Equal(t, tt.expected.percent, gauge.Percent,
				"Gauge percent should match expected value")
			assert.Equal(t, tt.expected.label, gauge.Label,
				"Gauge label should match expected value")
		})
	}
}
