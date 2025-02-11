// Package main provides the entry point for the Go scheduler monitoring tool
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"schedtrace-mon/internal/collector/godebug"
	"schedtrace-mon/internal/domain"
	"schedtrace-mon/internal/ui/termui"
)

//-----------------------------------------------------------------------------
// schedtrace-mon is a terminal-based monitoring tool for Go scheduler.
// It visualizes GODEBUG=schedtrace metrics in real-time using terminal UI.
//
// Usage:
//   schedtrace-mon -cmd "your_program_to_monitor"
//
// Press 'q' to exit.
//-----------------------------------------------------------------------------

func main() {
	targetCmd := flag.String("cmd", "", "Command to run and monitor (required)")
	flag.Parse()

	if *targetCmd == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Create cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize collector
	collector := godebug.NewCollector()

	// Initialize UI
	renderer := termui.NewRenderer()
	if err := renderer.Init(); err != nil {
		log.Fatalf("Failed to initialize UI: %v", err)
	}
	defer renderer.Close()

	// Handle UI events
	renderer.HandleEvents(func(event domain.Event) {
		switch event {
		case domain.EventQuit:
			cancel()
		case domain.EventResize:
			renderer.Display(collector)
		}
	})

	// Handle OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	// Start metrics collection
	if err := collector.Start(ctx, *targetCmd); err != nil {
		log.Fatalf("Failed to start collector: %v", err)
	}
	defer collector.Stop()

	// Main loop
	for {
		select {
		case <-ctx.Done():
			return
		default:
			renderer.Display(collector)
		}
	}
}
