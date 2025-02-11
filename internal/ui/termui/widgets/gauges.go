package widgets

import (
	"fmt"

	ui "github.com/gizak/termui/v3"
	twidgets "github.com/gizak/termui/v3/widgets"
)

// QueueGauges displays GRQ and LRQ metrics as gauge bars
type QueueGauges struct {
	GRQ *twidgets.Gauge // Global Run Queue gauge
	LRQ *twidgets.Gauge // Local Run Queues sum gauge
}

// NewQueueGauges creates new queue gauge widgets
func NewQueueGauges() *QueueGauges {
	grq := twidgets.NewGauge()
	grq.Title = "GRQ"
	grq.BarColor = ui.ColorGreen
	grq.TitleStyle.Fg = ui.ColorCyan

	lrq := twidgets.NewGauge()
	lrq.Title = "LRQ (sum)"
	lrq.BarColor = ui.ColorMagenta
	lrq.TitleStyle.Fg = ui.ColorCyan

	return &QueueGauges{
		GRQ: grq,
		LRQ: lrq,
	}
}

// Update refreshes gauges with current values
func (g *QueueGauges) Update(grqVal, grqMax, lrqVal, lrqMax int) {
	if grqMax == 0 {
		grqMax = 1
	}
	if lrqMax == 0 {
		lrqMax = 1
	}

	g.GRQ.Percent = grqVal * 100 / grqMax
	g.GRQ.Label = fmt.Sprintf("%d / %d", grqVal, grqMax)

	g.LRQ.Percent = lrqVal * 100 / lrqMax
	g.LRQ.Label = fmt.Sprintf("%d / %d", lrqVal, lrqMax)
}
