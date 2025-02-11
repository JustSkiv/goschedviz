// Package widgets provides UI components for scheduler metrics visualization
package widgets

import (
	"strconv"

	"schedtrace-mon/internal/domain"

	ui "github.com/gizak/termui/v3"
	twidgets "github.com/gizak/termui/v3/widgets"
)

// SchedTable displays current scheduler metrics in a table format
type SchedTable struct {
	*twidgets.Table
}

// NewSchedTable creates a new scheduler metrics table widget
func NewSchedTable() *SchedTable {
	t := &SchedTable{
		Table: twidgets.NewTable(),
	}
	t.Title = "Current Scheduler Values"
	t.TextStyle = ui.NewStyle(ui.ColorWhite)
	t.RowSeparator = false
	t.BorderStyle.Fg = ui.ColorYellow
	return t
}

// Update refreshes table with current scheduler data
func (t *SchedTable) Update(data domain.SchedData) {
	t.Rows = [][]string{
		{"Field", "Value"},
		{"Time (ms)", strconv.Itoa(data.TimeMs)},
		{"gomaxprocs", strconv.Itoa(data.GoMaxProcs)},
		{"idleprocs", strconv.Itoa(data.IdleProcs)},
		{"threads", strconv.Itoa(data.Threads)},
		{"spinningthreads", strconv.Itoa(data.SpinningThreads)},
		{"needspinning", strconv.Itoa(data.NeedSpinning)},
		{"idlethreads", strconv.Itoa(data.IdleThreads)},
		{"runqueue (GRQ)", strconv.Itoa(data.RunQueue)},
		{"LRQ (sum)", strconv.Itoa(data.LrqSum)},
		{"P count", strconv.Itoa(len(data.Lrq))},
	}
}
