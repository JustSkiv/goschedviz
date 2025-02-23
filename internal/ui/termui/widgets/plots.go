// Package widgets provides terminal UI components using termui library.
package widgets

import (
	"math"

	tui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	"github.com/JustSkiv/goschedviz/internal/ui"
)

// BaseHistoryPlot encapsulates common plot functionality
type BaseHistoryPlot struct {
	*widgets.Plot
}

// newBasePlot creates a new base plot with common settings
func newBasePlot() *BaseHistoryPlot {
	p := &BaseHistoryPlot{
		Plot: widgets.NewPlot(),
	}

	p.DataLabels = []string{
		"GRQ",
		"LRQ",
		"Threads",
		"IdleProcs",
		"Goroutines",
	}

	p.Data = make([][]float64, 5)
	for i := 0; i < 5; i++ {
		p.Data[i] = []float64{0, 0}
	}

	p.LineColors = []tui.Color{
		tui.ColorGreen,   // GRQ
		tui.ColorMagenta, // LRQ
		tui.ColorRed,     // Threads
		tui.ColorYellow,  // IdleProcs
		tui.ColorCyan,    // Goroutines
	}

	p.AxesColor = tui.ColorWhite
	p.DrawDirection = widgets.DrawLeft

	return p
}

// LinearHistoryPlot displays metrics using linear scale
type LinearHistoryPlot struct {
	*BaseHistoryPlot
}

// NewLinearHistoryPlot creates a new linear-scale plot
func NewLinearHistoryPlot() *LinearHistoryPlot {
	p := &LinearHistoryPlot{
		BaseHistoryPlot: newBasePlot(),
	}
	p.Title = "History Plot (linear)"
	return p
}

// Update updates plot with raw values
func (p *LinearHistoryPlot) Update(history []ui.HistoricalValues) {
	if len(history) < 2 {
		for i := 0; i < 5; i++ {
			p.Data[i] = []float64{0, 0}
		}
		return
	}

	length := len(history)
	grqVals := make([]float64, length)
	lrqVals := make([]float64, length)
	threadVals := make([]float64, length)
	idleProcVals := make([]float64, length)
	goroutineVals := make([]float64, length)

	for i, h := range history {
		grqVals[i] = float64(h.GRQ)
		lrqVals[i] = float64(h.LRQSum)
		threadVals[i] = float64(h.Threads)
		idleProcVals[i] = float64(h.IdleProcs)
		goroutineVals[i] = float64(h.Goroutines)
	}

	p.Data[0] = grqVals
	p.Data[1] = lrqVals
	p.Data[2] = threadVals
	p.Data[3] = idleProcVals
	p.Data[4] = goroutineVals
}

// LogHistoryPlot displays metrics using logarithmic scale
type LogHistoryPlot struct {
	*BaseHistoryPlot
}

// NewLogHistoryPlot creates a new logarithmic-scale plot
func NewLogHistoryPlot() *LogHistoryPlot {
	p := &LogHistoryPlot{
		BaseHistoryPlot: newBasePlot(),
	}
	p.Title = "History Plot (log)"
	return p
}

// toLogScale converts a value to logarithmic scale safely
func toLogScale(value float64) float64 {
	if value <= 0 {
		return 0
	}
	return math.Log10(value)
}

// Update updates plot with logarithmically scaled values
func (p *LogHistoryPlot) Update(history []ui.HistoricalValues) {
	if len(history) < 2 {
		for i := 0; i < 5; i++ {
			p.Data[i] = []float64{0, 0}
		}
		return
	}

	length := len(history)
	grqVals := make([]float64, length)
	lrqVals := make([]float64, length)
	threadVals := make([]float64, length)
	idleProcVals := make([]float64, length)
	goroutineVals := make([]float64, length)

	for i, h := range history {
		grqVals[i] = toLogScale(float64(h.GRQ))
		lrqVals[i] = toLogScale(float64(h.LRQSum))
		threadVals[i] = toLogScale(float64(h.Threads))
		idleProcVals[i] = toLogScale(float64(h.IdleProcs))
		goroutineVals[i] = toLogScale(float64(h.Goroutines))
	}

	p.Data[0] = grqVals
	p.Data[1] = lrqVals
	p.Data[2] = threadVals
	p.Data[3] = idleProcVals
	p.Data[4] = goroutineVals
}
