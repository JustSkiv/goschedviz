package widgets

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JustSkiv/goschedviz/internal/ui"
)

func TestTableWidget_New(t *testing.T) {
	// Test table initialization and default values
	table := NewTableWidget()
	require.NotNil(t, table, "NewTableWidget should return non-nil table")

	// Check initial configuration
	assert.Equal(t, "Current Values", table.Title,
		"Table should have correct title")
	assert.False(t, table.RowSeparator,
		"Table should not have row separators by default")

	// Verify initial rows structure
	assert.Greater(t, len(table.Rows), 0,
		"Table should have initial rows")
}

func TestTableWidget_Update(t *testing.T) {
	tests := []struct {
		name     string
		input    ui.CurrentValues
		expected [][]string
	}{
		{
			name: "zero values",
			input: ui.CurrentValues{
				TimeMs:          0,
				GoMaxProcs:      1,
				IdleProcs:       0,
				Threads:         0,
				SpinningThreads: 0,
				NeedSpinning:    0,
				IdleThreads:     0,
				RunQueue:        0,
				LRQSum:          0,
				NumP:            1,
			},
			expected: [][]string{
				{"Time (ms)", "0"},
				{"gomaxprocs", "1"},
				{"idleprocs", "0"},
				{"threads", "0"},
				{"spinningthreads", "0"},
				{"needspinning", "0"},
				{"idlethreads", "0"},
				{"runqueue (GRQ)", "0"},
				{"LRQ (sum)", "0"},
				{"Number of P", "1"},
			},
		},
		{
			name: "typical load values",
			input: ui.CurrentValues{
				TimeMs:          1500,
				GoMaxProcs:      4,
				IdleProcs:       2,
				Threads:         8,
				SpinningThreads: 1,
				NeedSpinning:    1,
				IdleThreads:     3,
				RunQueue:        5,
				LRQSum:          10,
				NumP:            4,
			},
			expected: [][]string{
				{"Time (ms)", "1500"},
				{"gomaxprocs", "4"},
				{"idleprocs", "2"},
				{"threads", "8"},
				{"spinningthreads", "1"},
				{"needspinning", "1"},
				{"idlethreads", "3"},
				{"runqueue (GRQ)", "5"},
				{"LRQ (sum)", "10"},
				{"Number of P", "4"},
			},
		},
		{
			name: "high load values",
			input: ui.CurrentValues{
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
			},
			expected: [][]string{
				{"Time (ms)", "5000"},
				{"gomaxprocs", "32"},
				{"idleprocs", "0"},
				{"threads", "64"},
				{"spinningthreads", "8"},
				{"needspinning", "4"},
				{"idlethreads", "0"},
				{"runqueue (GRQ)", "100"},
				{"LRQ (sum)", "500"},
				{"Number of P", "32"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := NewTableWidget()
			table.Update(tt.input)

			// Check total number of rows
			assert.Equal(t, len(tt.expected), len(table.Rows),
				"Table should have correct number of rows")

			// Check each row's content
			for i, expectedRow := range tt.expected {
				assert.Equal(t, expectedRow, table.Rows[i],
					"Row %d content should match expected values", i)
			}
		})
	}
}

func TestTableWidget_UpdateDataConsistency(t *testing.T) {
	// Test that multiple updates don't corrupt the table structure
	table := NewTableWidget()
	initialRows := len(table.Rows)

	// Perform multiple updates
	for i := 0; i < 5; i++ {
		data := ui.CurrentValues{
			TimeMs:     i * 1000,
			GoMaxProcs: 4,
			NumP:       4,
		}
		table.Update(data)

		assert.Equal(t, initialRows, len(table.Rows),
			"Number of rows should remain constant after update")

		// Check first row format
		assert.Equal(t, 2, len(table.Rows[0]),
			"Each row should have exactly 2 columns")
	}
}
