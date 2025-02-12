package widgets

import (
	"strings"
	"testing"

	"github.com/gizak/termui/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JustSkiv/goschedviz/internal/ui"
)

func TestInfoBox_New(t *testing.T) {
	info := NewInfoBox()
	require.NotNil(t, info, "NewInfoBox should return non-nil info box")

	// Check initial configuration
	assert.Equal(t, "Information", info.Title,
		"Info box should have correct title")
	assert.Equal(t, termui.ColorCyan, info.BorderStyle.Fg,
		"Info box should have cyan border")
}

func TestInfoBox_Update(t *testing.T) {
	tests := []struct {
		name          string
		currentValues ui.CurrentValues
		gaugeValues   ui.GaugeValues
		expectations  []string // Strings that should be present in the text
	}{
		{
			name:          "zero values",
			currentValues: ui.CurrentValues{},
			gaugeValues: ui.GaugeValues{
				GRQ: struct {
					Current int
					Max     int
				}{0, 0},
				LRQ: struct {
					Current int
					Max     int
				}{0, 0},
			},
			expectations: []string{
				"Exit: press 'q'",
				"Last update:", // Just check presence of the prefix
				"Max GRQ: 0",
				"Max LRQ (sum): 0",
			},
		},
		{
			name: "normal load values",
			currentValues: ui.CurrentValues{
				TimeMs:   1000,
				RunQueue: 5,
				LRQSum:   10,
			},
			gaugeValues: ui.GaugeValues{
				GRQ: struct {
					Current int
					Max     int
				}{5, 20},
				LRQ: struct {
					Current int
					Max     int
				}{10, 50},
			},
			expectations: []string{
				"Exit: press 'q'",
				"Last update:",
				"Max GRQ: 20",
				"Max LRQ (sum): 50",
			},
		},
		{
			name: "high load values",
			currentValues: ui.CurrentValues{
				TimeMs:   5000,
				RunQueue: 100,
				LRQSum:   500,
			},
			gaugeValues: ui.GaugeValues{
				GRQ: struct {
					Current int
					Max     int
				}{100, 200},
				LRQ: struct {
					Current int
					Max     int
				}{500, 1000},
			},
			expectations: []string{
				"Exit: press 'q'",
				"Last update:",
				"Max GRQ: 200",
				"Max LRQ (sum): 1000",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := NewInfoBox()
			info.Update(tt.currentValues, tt.gaugeValues)

			// Check that all expected strings are present
			for _, expectedStr := range tt.expectations {
				assert.Contains(t, info.Text, expectedStr,
					"Info box text should contain '%s'", expectedStr)
			}

			// Verify timestamp format
			lines := strings.Split(info.Text, "\n")
			var hasTimestamp bool
			for _, line := range lines {
				if strings.HasPrefix(line, "Last update: ") {
					timeStr := strings.TrimPrefix(line, "Last update: ")
					timeStr = strings.TrimSpace(timeStr)
					// Check time format (HH:MM:SS)
					assert.Regexp(t, `^\d{2}:\d{2}:\d{2}$`, timeStr,
						"Timestamp should be in format HH:MM:SS")
					hasTimestamp = true
					break
				}
			}
			assert.True(t, hasTimestamp, "Info box should contain timestamp")
		})
	}
}

func TestInfoBox_UpdateConsistency(t *testing.T) {
	info := NewInfoBox()

	// Test multiple updates
	updates := 5
	var prevLines []string

	for i := 0; i < updates; i++ {
		currentValues := ui.CurrentValues{
			TimeMs:   i * 1000,
			RunQueue: i * 2,
			LRQSum:   i * 3,
		}

		gaugeValues := ui.GaugeValues{
			GRQ: struct {
				Current int
				Max     int
			}{i * 2, 10},
			LRQ: struct {
				Current int
				Max     int
			}{i * 3, 20},
		}

		info.Update(currentValues, gaugeValues)

		// Check text structure consistency
		lines := strings.Split(info.Text, "\n")
		assert.GreaterOrEqual(t, len(lines), 4,
			"Info box should maintain at least 4 lines of text")

		if prevLines != nil {
			// Check that static content remains unchanged
			assert.Equal(t, prevLines[0], lines[0],
				"Exit instruction line should remain constant")
			assert.Equal(t, prevLines[len(prevLines)-2], lines[len(lines)-2],
				"Max GRQ line should remain constant")
			assert.Equal(t, prevLines[len(prevLines)-1], lines[len(lines)-1],
				"Max LRQ line should remain constant")

			// Check timestamp line format
			assert.Regexp(t, `^Last update: \d{2}:\d{2}:\d{2}$`, strings.TrimSpace(lines[1]),
				"Timestamp line should maintain correct format")
		}

		prevLines = lines
	}
}
