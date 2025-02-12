// Package widgets provides terminal UI components using termui library.
package widgets

import (
	"github.com/gizak/termui/v3/widgets"

	tui "github.com/gizak/termui/v3"

	"github.com/JustSkiv/goschedviz/internal/ui"
)

// HistoryPlot displays GRQ/LRQ history as a line graph.
// It shows two lines:
// - Green line for Global Run Queue history
// - Magenta line for total Local Run Queues history
type HistoryPlot struct {
	*widgets.Plot
}

// NewHistoryPlot creates a new history plot widget with default styling.
func NewHistoryPlot() *HistoryPlot {
	p := &HistoryPlot{
		Plot: widgets.NewPlot(),
	}
	p.Title = "GRQ / LRQ History"

	// Initialize with default values
	p.DataLabels = make([]string, 2)
	p.Data = make([][]float64, 2)
	p.Data[0] = []float64{0, 0} // At least two values are required for plotting
	p.Data[1] = []float64{0, 0}
	p.LineColors = make([]tui.Color, 2)
	p.LineColors[0] = tui.ColorGreen
	p.LineColors[1] = tui.ColorMagenta

	p.AxesColor = tui.ColorWhite
	p.DrawDirection = widgets.DrawLeft

	return p
}

// Update updates plot with new historical values.
// Converts integer metrics to float64 for plotting.
func (p *HistoryPlot) Update(history []ui.HistoricalValues) {
	if len(history) == 0 {
		// If there is no data, keep initial values
		return
	}

	grqVals := make([]float64, len(history))
	lrqVals := make([]float64, len(history))

	for i, h := range history {
		grqVals[i] = float64(h.GRQ)
		lrqVals[i] = float64(h.LRQSum)
	}

	p.Data[0] = grqVals
	p.Data[1] = lrqVals
}
