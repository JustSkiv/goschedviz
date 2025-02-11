package widgets

import (
	"fmt"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// LRQBarChart displays Local Run Queues as bars.
type LRQBarChart struct {
	*widgets.BarChart
}

// NewLRQBarChart creates a new bar chart for LRQ visualization.
func NewLRQBarChart() *LRQBarChart {
	b := &LRQBarChart{
		BarChart: widgets.NewBarChart(),
	}
	b.Title = "Local Run Queues (per P)"
	b.BarWidth = 3
	b.BarGap = 1
	b.BarColors = []termui.Color{termui.ColorCyan}
	b.LabelStyles = []termui.Style{termui.NewStyle(termui.ColorYellow)}
	b.NumFormatter = func(f float64) string {
		return fmt.Sprintf("%.0f", f)
	}
	return b
}

// Update updates bar chart with new LRQ values.
func (b *LRQBarChart) Update(lrq []int) {
	b.Data = nil
	b.Labels = nil
	for i, v := range lrq {
		b.Data = append(b.Data, float64(v))
		b.Labels = append(b.Labels, fmt.Sprintf("P%d", i))
	}
}
