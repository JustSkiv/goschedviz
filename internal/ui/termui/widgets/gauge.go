// Package widgets provides terminal UI components using termui library.
// Each widget encapsulates specific visualization logic and styling,
// implementing a consistent Update interface for data refresh.
package widgets

import (
	"fmt"

	tui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// GRQGauge displays Global Run Queue gauge.
// It shows current GRQ value as both a percentage bar and absolute value.
type GRQGauge struct {
	*widgets.Gauge
}

// NewGRQGauge creates a new GRQ gauge widget with default styling.
func NewGRQGauge() *GRQGauge {
	g := &GRQGauge{
		Gauge: widgets.NewGauge(),
	}
	g.Title = "GRQ"
	g.BarColor = tui.ColorGreen
	g.TitleStyle.Fg = tui.ColorCyan
	return g
}

// Update updates GRQ gauge values showing current/max ratio.
func (g *GRQGauge) Update(data struct{ Current, Max int }) {
	g.Percent = data.Current * 100 / data.Max
	g.Label = fmt.Sprintf("%d / %d", data.Current, data.Max)
}

// LRQGauge displays Local Run Queues sum gauge.
// It visualizes the total load across all P's local queues.
type LRQGauge struct {
	*widgets.Gauge
}

// NewLRQGauge creates a new LRQ gauge widget with default styling.
func NewLRQGauge() *LRQGauge {
	g := &LRQGauge{
		Gauge: widgets.NewGauge(),
	}
	g.Title = "LRQ (sum)"
	g.BarColor = tui.ColorMagenta
	g.TitleStyle.Fg = tui.ColorCyan
	return g
}

// Update updates LRQ gauge values showing current/max ratio.
func (g *LRQGauge) Update(data struct{ Current, Max int }) {
	g.Percent = data.Current * 100 / data.Max
	g.Label = fmt.Sprintf("%d / %d", data.Current, data.Max)
}
