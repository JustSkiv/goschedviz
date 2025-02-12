// Package termui implements terminal-based UI using termui library.
package termui

import (
	"github.com/gizak/termui/v3"

	"github.com/JustSkiv/goschedviz/internal/ui"
	"github.com/JustSkiv/goschedviz/internal/ui/termui/widgets"
)

// TermUI implements ui.Presenter interface using termui library.
type TermUI struct {
	table    *widgets.TableWidget
	barChart *widgets.LRQBarChart
	grqGauge *widgets.GRQGauge
	lrqGauge *widgets.LRQGauge
	plot     *widgets.HistoryPlot
	info     *widgets.InfoBox
	grid     *termui.Grid
	done     chan struct{}
}

// New creates a new terminal UI implementation.
func New() *TermUI {
	return &TermUI{
		done: make(chan struct{}),
	}
}

// Start implements ui.Presenter interface.
func (t *TermUI) Start() error {
	if err := termui.Init(); err != nil {
		return err
	}

	// Initialize widgets
	t.table = widgets.NewTableWidget()
	t.barChart = widgets.NewLRQBarChart()
	t.grqGauge = widgets.NewGRQGauge()
	t.lrqGauge = widgets.NewLRQGauge()
	t.plot = widgets.NewHistoryPlot()
	t.info = widgets.NewInfoBox()

	// Setup grid
	t.setupGrid()

	// Start event handling
	go t.handleEvents()

	return nil
}

// Stop implements ui.Presenter interface.
func (t *TermUI) Stop() {
	termui.Close()
}

// Done implements ui.Presenter interface.
func (t *TermUI) Done() <-chan struct{} {
	return t.done
}

// Update implements ui.Presenter interface.
func (t *TermUI) Update(data ui.UIData) {
	t.table.Update(data.Current)
	t.barChart.Update(data.Current.LRQ)
	t.grqGauge.Update(data.Gauges.GRQ)
	t.lrqGauge.Update(data.Gauges.LRQ)
	t.plot.Update(data.History)
	t.info.Update(data.Current, data.Gauges)

	termui.Render(t.grid)
}

// setupGrid initializes the terminal UI layout.
func (t *TermUI) setupGrid() {
	t.grid = termui.NewGrid()
	termWidth, termHeight := termui.TerminalDimensions()
	t.grid.SetRect(0, 0, termWidth, termHeight)

	t.grid.Set(
		termui.NewRow(0.3,
			termui.NewCol(0.4, t.table),
			termui.NewCol(0.6, t.barChart),
		),
		termui.NewRow(0.3,
			termui.NewCol(0.5, t.grqGauge),
			termui.NewCol(0.5, t.lrqGauge),
		),
		termui.NewRow(0.4,
			termui.NewCol(0.8, t.plot),
			termui.NewCol(0.2, t.info),
		),
	)
}

// handleEvents processes terminal UI events.
func (t *TermUI) handleEvents() {
	uiEvents := termui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				close(t.done)
				return
			case "<Resize>":
				payload := e.Payload.(termui.Resize)
				t.grid.SetRect(0, 0, payload.Width, payload.Height)
				termui.Clear()
				termui.Render(t.grid)
			}
		case <-t.done:
			return
		}
	}
}
