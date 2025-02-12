package widgets

import (
	"strconv"

	tui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	"github.com/JustSkiv/goschedviz/internal/ui"
)

// TableWidget displays current scheduler values.
type TableWidget struct {
	*widgets.Table
}

// NewTableWidget creates a new table widget.
func NewTableWidget() *TableWidget {
	t := &TableWidget{
		Table: widgets.NewTable(),
	}
	t.Title = "Current Scheduler Values"
	t.TextStyle = tui.NewStyle(tui.ColorWhite)
	t.RowSeparator = false
	t.BorderStyle.Fg = tui.ColorYellow

	// Add initial empty data to prevent panic on first render
	t.Rows = [][]string{
		{"Field", "Value"},
		{"Time (ms)", "-"},
		{"gomaxprocs", "-"},
		{"idleprocs", "-"},
		{"threads", "-"},
		{"spinningthreads", "-"},
		{"needspinning", "-"},
		{"idlethreads", "-"},
		{"runqueue (GRQ)", "-"},
		{"LRQ (sum)", "-"},
		{"Number of P", "-"},
	}

	return t
}

// Update updates table with current values.
func (t *TableWidget) Update(data ui.CurrentValues) {
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
		{"LRQ (sum)", strconv.Itoa(data.LRQSum)},
		{"Number of P", strconv.Itoa(data.NumP)},
	}
}
