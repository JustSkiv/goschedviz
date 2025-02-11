package termui

import (
	"fmt"
	"time"

	ui "github.com/gizak/termui/v3"

	"schedtrace-mon/internal/collector"
	"schedtrace-mon/internal/domain"
)

// Renderer implements the ui.Presenter interface using termui library
type Renderer struct {
	layout *Layout
	events chan domain.Event
	ticker *time.Ticker
}

// NewRenderer creates a new terminal UI renderer
func NewRenderer() *Renderer {
	return &Renderer{
		layout: NewLayout(),
		events: make(chan domain.Event),
		ticker: time.NewTicker(500 * time.Millisecond),
	}
}

// Init initializes the terminal UI system
func (r *Renderer) Init() error {
	if err := ui.Init(); err != nil {
		return fmt.Errorf("initializing UI: %w", err)
	}

	width, height := ui.TerminalDimensions()
	r.layout.UpdateSize(width, height)

	go r.handleEvents()
	return nil
}

// Close performs cleanup of UI resources
func (r *Renderer) Close() error {
	r.ticker.Stop()
	ui.Close()
	return nil
}

// Display shows current metrics from collector
func (r *Renderer) Display(c collector.Collector) {
	current := c.GetCurrent()
	history := c.GetHistory()

	r.updateWidgets(current, history)
	ui.Render(r.layout.Grid)
}

// HandleEvents sets up UI event handling
func (r *Renderer) HandleEvents(handler func(domain.Event)) {
	go func() {
		for e := range r.events {
			handler(e)
		}
	}()
}

func (r *Renderer) handleEvents() {
	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				r.events <- domain.EventQuit
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				r.layout.UpdateSize(payload.Width, payload.Height)
				r.events <- domain.EventResize
			}
		case <-r.ticker.C:
			// UI update tick
		}
	}
}

func (r *Renderer) updateWidgets(current domain.SchedData, history []domain.SchedData) {
	// Update table with current values
	r.layout.Table.Update(current)

	// Calculate maximums for gauges
	maxGRQ, maxLRQ := 0, 0
	var grqVals, lrqVals []float64
	for _, d := range history {
		if d.RunQueue > maxGRQ {
			maxGRQ = d.RunQueue
		}
		if d.LrqSum > maxLRQ {
			maxLRQ = d.LrqSum
		}
		grqVals = append(grqVals, float64(d.RunQueue))
		lrqVals = append(lrqVals, float64(d.LrqSum))
	}

	// Update all widgets
	r.layout.Gauges.Update(current.RunQueue, maxGRQ, current.LrqSum, maxLRQ)
	r.layout.Chart.Update(grqVals, lrqVals)
	r.layout.InfoPanel.Update(len(history), maxGRQ, maxLRQ)
}
