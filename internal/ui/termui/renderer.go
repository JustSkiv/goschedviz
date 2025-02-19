// Package termui implements terminal-based UI using termui library.
package termui

import (
	"github.com/gizak/termui/v3"

	"github.com/JustSkiv/goschedviz/internal/ui"
	"github.com/JustSkiv/goschedviz/internal/ui/termui/widgets"
)

// terminalAPI defines the interface for terminal operations
type terminalAPI interface {
	Init() error
	Close()
	PollEvents() <-chan termui.Event
	Render(ps ...termui.Drawable)
	TerminalDimensions() (int, int)
	Clear()
}

// realTerminal implements terminalAPI using actual termui library
type realTerminal struct{}

func (r *realTerminal) Init() error                     { return termui.Init() }
func (r *realTerminal) Close()                          { termui.Close() }
func (r *realTerminal) PollEvents() <-chan termui.Event { return termui.PollEvents() }
func (r *realTerminal) Render(ps ...termui.Drawable)    { termui.Render(ps...) }
func (r *realTerminal) TerminalDimensions() (int, int)  { return termui.TerminalDimensions() }
func (r *realTerminal) Clear()                          { termui.Clear() }

// testTerminal implements terminalAPI for testing purposes
type testTerminal struct {
	events chan termui.Event
}

func newTestTerminal() *testTerminal {
	return &testTerminal{
		events: make(chan termui.Event),
	}
}

func (t *testTerminal) Init() error                     { return nil }
func (t *testTerminal) Close()                          { close(t.events) }
func (t *testTerminal) PollEvents() <-chan termui.Event { return t.events }
func (t *testTerminal) Render(ps ...termui.Drawable)    { /* no-op */ }
func (t *testTerminal) TerminalDimensions() (int, int)  { return 100, 40 }
func (t *testTerminal) Clear()                          { /* no-op */ }

// SendEvent sends an event to the test terminal event channel
func (t *testTerminal) SendEvent(e termui.Event) {
	t.events <- e
}

// TermUI implements ui.Presenter interface using termui library.
type TermUI struct {
	table          *widgets.TableWidget
	barChart       *widgets.LRQBarChart
	grqGauge       *widgets.GRQGauge
	lrqGauge       *widgets.LRQGauge
	threadsGauge   *widgets.ThreadsGauge
	idleProcsGauge *widgets.IdleProcsGauge
	plot           *widgets.HistoryPlot
	info           *widgets.InfoBox
	grid           *termui.Grid
	done           chan struct{}
	term           terminalAPI
}

// New creates a new terminal UI implementation.
func New() *TermUI {
	return &TermUI{
		done: make(chan struct{}),
		term: &realTerminal{},
	}
}

// newWithTerminal creates a new terminal UI with custom terminal implementation for testing.
func newWithTerminal(term terminalAPI) *TermUI {
	return &TermUI{
		done: make(chan struct{}),
		term: term,
	}
}

// Start implements ui.Presenter interface.
func (t *TermUI) Start() error {
	if err := t.term.Init(); err != nil {
		return err
	}

	// Initialize widgets
	t.table = widgets.NewTableWidget()
	t.barChart = widgets.NewLRQBarChart()
	t.grqGauge = widgets.NewGRQGauge()
	t.lrqGauge = widgets.NewLRQGauge()
	t.threadsGauge = widgets.NewThreadsGauge()
	t.idleProcsGauge = widgets.NewIdleProcsGauge()
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
	t.term.Close()
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
	t.threadsGauge.Update(data.Gauges.Threads)
	t.idleProcsGauge.Update(data.Gauges.IdleProcs)
	t.plot.Update(data.History)
	t.info.Update(data.Current, data.Gauges)

	t.term.Render(t.grid)
}

// setupGrid initializes the terminal UI layout.
func (t *TermUI) setupGrid() {
	t.grid = termui.NewGrid()
	width, height := t.term.TerminalDimensions()
	t.grid.SetRect(0, 0, width, height)

	t.grid.Set(
		termui.NewRow(0.3,
			termui.NewCol(0.4, t.table),
			termui.NewCol(0.6, t.barChart),
		),
		termui.NewRow(0.3,
			termui.NewCol(0.5,
				termui.NewRow(0.5, t.threadsGauge),
				termui.NewRow(0.5, t.idleProcsGauge),
			),
			termui.NewCol(0.5,
				termui.NewRow(0.5, t.lrqGauge),
				termui.NewRow(0.5, t.grqGauge),
			),
		),
		termui.NewRow(0.4,
			termui.NewCol(0.8, t.plot),
			termui.NewCol(0.2, t.info),
		),
	)
}

// handleEvents processes terminal UI events.
func (t *TermUI) handleEvents() {
	uiEvents := t.term.PollEvents()
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
				t.term.Clear()
				t.term.Render(t.grid)
			}
		case <-t.done:
			return
		}
	}
}
