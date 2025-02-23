// Package widgets provides terminal UI components using termui library.
package widgets

import (
	"fmt"
	"time"

	"github.com/gizak/termui/v3/widgets"

	tui "github.com/gizak/termui/v3"

	"github.com/JustSkiv/goschedviz/internal/ui"
)

// InfoBox displays additional monitoring information.
// Shows:
// - Exit instructions
// - Last update timestamp
// - Maximum values for GRQ and goroutines
type InfoBox struct {
	*widgets.Paragraph
}

// NewInfoBox creates a new info box widget with default styling.
func NewInfoBox() *InfoBox {
	i := &InfoBox{
		Paragraph: widgets.NewParagraph(),
	}
	i.Title = "Information"
	i.BorderStyle.Fg = tui.ColorCyan
	return i
}

// Update refreshes info box with current monitoring state.
func (i *InfoBox) Update(current ui.CurrentValues, gauges ui.GaugeValues) {
	i.Text = fmt.Sprintf(
		"Last update: %s\n"+
			"Max GRQ: %d\n"+
			"Max Gs: %d\n"+
			"Exit: press 'q'",
		time.Now().Format("15:04:05"),
		gauges.GRQ.Max,
		gauges.Goroutines.Max,
	)
}
