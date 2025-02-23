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
	require.NotNil(t, info)
	assert.Equal(t, "Information", info.Title)
	assert.Equal(t, termui.ColorCyan, info.BorderStyle.Fg)
}

func TestInfoBox_Update(t *testing.T) {
	tests := []struct {
		name          string
		currentValues ui.CurrentValues
		gaugeValues   ui.GaugeValues
		expectations  []string
	}{
		{
			name:          "zero values",
			currentValues: ui.CurrentValues{},
			gaugeValues: ui.GaugeValues{
				GRQ:        struct{ Current, Max int }{0, 0},
				Goroutines: struct{ Current, Max int }{0, 0},
			},
			expectations: []string{
				"Last update:",
				"Max GRQ: 0",
				"Max Gs: 0",
				"Exit: press 'q'",
			},
		},
		{
			name: "normal values",
			currentValues: ui.CurrentValues{
				TimeMs:     1000,
				RunQueue:   5,
				Goroutines: 100,
			},
			gaugeValues: ui.GaugeValues{
				GRQ:        struct{ Current, Max int }{5, 20},
				Goroutines: struct{ Current, Max int }{100, 200},
			},
			expectations: []string{
				"Last update:",
				"Max GRQ: 20",
				"Max Gs: 200",
				"Exit: press 'q'",
			},
		},
		{
			name: "high load values",
			currentValues: ui.CurrentValues{
				TimeMs:     5000,
				RunQueue:   100,
				Goroutines: 1000,
			},
			gaugeValues: ui.GaugeValues{
				GRQ:        struct{ Current, Max int }{100, 200},
				Goroutines: struct{ Current, Max int }{1000, 2000},
			},
			expectations: []string{
				"Last update:",
				"Max GRQ: 200",
				"Max Gs: 2000",
				"Exit: press 'q'",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := NewInfoBox()
			info.Update(tt.currentValues, tt.gaugeValues)

			lines := strings.Split(info.Text, "\n")
			assert.Equal(t, 4, len(lines), "Should have exactly 4 lines")

			for i, expected := range tt.expectations {
				assert.Contains(t, lines[i], expected,
					"Line %d should contain %q", i, expected)
			}
		})
	}
}

func TestInfoBox_UpdateConsistency(t *testing.T) {
	info := NewInfoBox()
	gaugeValues := ui.GaugeValues{
		GRQ:        struct{ Current, Max int }{5, 10},
		Goroutines: struct{ Current, Max int }{100, 200},
	}

	// Do multiple updates
	for i := 0; i < 3; i++ {
		info.Update(ui.CurrentValues{}, gaugeValues)

		lines := strings.Split(info.Text, "\n")
		require.Equal(t, 4, len(lines), "Should always have 4 lines")

		assert.Regexp(t, `^Last update: \d{2}:\d{2}:\d{2}$`, lines[0])
		assert.Contains(t, lines[1], "Max GRQ: 10")
		assert.Contains(t, lines[2], "Max Gs: 200")
		assert.Equal(t, "Exit: press 'q'", lines[3])
	}
}
