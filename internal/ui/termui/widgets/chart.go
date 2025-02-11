package widgets

import (
	ui "github.com/gizak/termui/v3"
	twidgets "github.com/gizak/termui/v3/widgets"
)

// QueueChart displays historical GRQ and LRQ values
type QueueChart struct {
	*twidgets.Plot
}

// NewQueueChart creates a new queue history chart
func NewQueueChart() *QueueChart {
	p := &QueueChart{
		Plot: twidgets.NewPlot(),
	}
	p.Title = "GRQ / LRQ History"
	p.Data = make([][]float64, 2)
	p.AxesColor = ui.ColorWhite
	p.LineColors[0] = ui.ColorGreen   // GRQ
	p.LineColors[1] = ui.ColorMagenta // LRQ
	p.DrawDirection = twidgets.DrawLeft
	return p
}

// Update refreshes chart with historical data
func (p *QueueChart) Update(grqHistory, lrqHistory []float64) {
	p.Data[0] = grqHistory
	p.Data[1] = lrqHistory
}
