// Command monitor provides real-time visualization of Go scheduler metrics.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/JustSkiv/goschedviz/internal/collector/godebug"
	"github.com/JustSkiv/goschedviz/internal/domain"
	"github.com/JustSkiv/goschedviz/internal/ui"
	"github.com/JustSkiv/goschedviz/internal/ui/termui"
)

func main() {
	var (
		targetPath = flag.String("target", "", "Path to Go program to monitor")
		period     = flag.Int("period", 1000, "GODEBUG schedtrace period in milliseconds")
	)

	flag.Parse()

	if *targetPath == "" {
		fmt.Println("Please specify target program path with -target flag")
		os.Exit(1)
	}

	// Create collector
	collector := godebug.New(*targetPath, *period)

	// Create UI
	presenter := termui.New()
	if err := presenter.Start(); err != nil {
		log.Fatalf("Failed to initialize UI: %v", err)
	}
	defer presenter.Stop()

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		cancel()
	}()

	// Start collecting metrics
	snapshots, err := collector.Start(ctx)
	if err != nil {
		log.Fatalf("Failed to start collector: %v", err)
	}
	defer collector.Stop()

	// Create monitor state
	state := &domain.MonitorState{}

	// Main update loop
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case snapshot, ok := <-snapshots:
			if !ok {
				return
			}
			state.Update(snapshot)

		case <-ticker.C:
			latest, history := state.GetSnapshot()

			// Convert domain data to UI data
			uiData := convertToUIData(latest, history)
			presenter.Update(uiData)

		case <-presenter.Done():
			return

		case <-ctx.Done():
			return
		}
	}
}

// convertToUIData converts domain data to UI-specific format
func convertToUIData(latest domain.SchedulerSnapshot, history []domain.SchedulerSnapshot) ui.UIData {
	// Calculate max values for gauges
	maxGRQ, maxLRQ, maxThreads, maxIdleProcs := 0, 0, 0, 0
	histValues := make([]ui.HistoricalValues, len(history))

	for i, h := range history {
		if h.RunQueue > maxGRQ {
			maxGRQ = h.RunQueue
		}
		if h.LRQSum > maxLRQ {
			maxLRQ = h.LRQSum
		}
		if h.Threads > maxThreads {
			maxThreads = h.Threads
		}
		if h.IdleProcs > maxIdleProcs {
			maxIdleProcs = h.IdleProcs
		}

		histValues[i] = ui.HistoricalValues{
			TimeMs:    h.TimeMs,
			GRQ:       h.RunQueue,
			LRQSum:    h.LRQSum,
			IdleProcs: h.IdleProcs,
			Threads:   h.Threads,
		}
	}

	// Ensure non-zero max values for gauges
	if maxGRQ == 0 {
		maxGRQ = 1
	}
	if maxLRQ == 0 {
		maxLRQ = 1
	}
	if maxThreads == 0 {
		maxThreads = 1
	}
	if maxIdleProcs == 0 {
		maxIdleProcs = 1
	}

	return ui.UIData{
		Current: ui.CurrentValues{
			TimeMs:          latest.TimeMs,
			GoMaxProcs:      latest.GoMaxProcs,
			IdleProcs:       latest.IdleProcs,
			Threads:         latest.Threads,
			SpinningThreads: latest.SpinningThreads,
			NeedSpinning:    latest.NeedSpinning,
			IdleThreads:     latest.IdleThreads,
			RunQueue:        latest.RunQueue,
			LRQSum:          latest.LRQSum,
			NumP:            len(latest.LRQ),
			LRQ:             latest.LRQ,
		},
		History: histValues,
		Gauges: ui.GaugeValues{
			GRQ: struct {
				Current int
				Max     int
			}{
				Current: latest.RunQueue,
				Max:     maxGRQ,
			},
			LRQ: struct {
				Current int
				Max     int
			}{
				Current: latest.LRQSum,
				Max:     maxLRQ,
			},
			Threads: struct {
				Current int
				Max     int
			}{
				Current: latest.Threads,
				Max:     maxThreads,
			},
			IdleProcs: struct {
				Current int
				Max     int
			}{
				Current: latest.IdleProcs,
				Max:     maxIdleProcs,
			},
		},
	}
}
