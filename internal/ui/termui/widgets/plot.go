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
// - Yellow line for IdleProcs history
// - Blue line for Threads history
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
	p.DataLabels = make([]string, 4)
	p.Data = make([][]float64, 4)

	for i := 0; i < 4; i++ {
		p.Data[i] = []float64{0, 0}
	}

	p.LineColors = []tui.Color{
		tui.ColorGreen,   // GRQ
		tui.ColorMagenta, // LRQ
		tui.ColorYellow,  // IdleProcs
		tui.ColorRed,     // Threads
	}

	p.AxesColor = tui.ColorWhite
	p.DrawDirection = widgets.DrawLeft

	return p
}

// Update updates plot with new historical values.
// Converts integer metrics to float64 for plotting.
func (p *HistoryPlot) Update(history []ui.HistoricalValues) {
	if len(history) < 2 {
		for i := 0; i < 4; i++ {
			p.Data[i] = []float64{0, 0}
		}
		return
	}

	length := len(history)
	grqVals := make([]float64, length)
	lrqVals := make([]float64, length)
	idleProcVals := make([]float64, length)
	threadVals := make([]float64, length)

	for i, h := range history {
		grqVals[i] = float64(h.GRQ)
		lrqVals[i] = float64(h.LRQSum)
		idleProcVals[i] = float64(h.IdleProcs)
		threadVals[i] = float64(h.Threads)
	}

	p.Data[0] = grqVals
	p.Data[1] = lrqVals
	p.Data[2] = idleProcVals
	p.Data[3] = threadVals
}
