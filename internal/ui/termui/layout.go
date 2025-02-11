package termui

import (
	ui "github.com/gizak/termui/v3"
	"schedtrace-mon/internal/ui/termui/widgets"
)

//-----------------------------------------------------------------------------
// Layout manages UI components arrangement in a 3-row grid:
//
//  ┌─────────────────────────────────┬─────────────────────────────────────┐
//  │  Table (0.3 height, col 0.4)    │  BarChart (0.3 height, col 0.6)    │
//  ├─────────────────────────────────┴─────────────────────────────────────┤
//  │  GRQ Gauge (col 0.5) | LRQ Gauge (col 0.5)  (0.3 height)            │
//  ├─────────────────────────────────┬─────────────────────────────────────┤
//  │  Plot (col 0.8)                 │  Info (col 0.2) (0.4 height)       │
//  └─────────────────────────────────┴─────────────────────────────────────┘
//-----------------------------------------------------------------------------

// Layout manages UI components arrangement
type Layout struct {
	Grid      *ui.Grid
	Table     *widgets.SchedTable
	Gauges    *widgets.QueueGauges
	Chart     *widgets.QueueChart
	InfoPanel *widgets.InfoPanel
}

// NewLayout creates a new UI layout with all widgets
func NewLayout() *Layout {
	return &Layout{
		Grid:      ui.NewGrid(),
		Table:     widgets.NewSchedTable(),
		Gauges:    widgets.NewQueueGauges(),
		Chart:     widgets.NewQueueChart(),
		InfoPanel: widgets.NewInfoPanel(),
	}
}

// UpdateSize recalculates layout for new terminal dimensions
func (l *Layout) UpdateSize(width, height int) {
	l.Grid.SetRect(0, 0, width, height)
	l.Grid.Set(
		ui.NewRow(0.3,
			ui.NewCol(0.4, l.Table),
			ui.NewCol(0.6, l.Chart),
		),
		ui.NewRow(0.3,
			ui.NewCol(0.5, l.Gauges.GRQ),
			ui.NewCol(0.5, l.Gauges.LRQ),
		),
		ui.NewRow(0.4,
			ui.NewCol(0.8, l.Chart),
			ui.NewCol(0.2, l.InfoPanel),
		),
	)
}
