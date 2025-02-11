package godebug

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"schedtrace-mon/internal/collector"
	"schedtrace-mon/internal/domain"
)

// Collector implements the collector.Collector interface for GODEBUG=schedtrace metrics
type Collector struct {
	mu      sync.RWMutex
	cmd     *exec.Cmd
	latest  domain.SchedData
	history []domain.SchedData
}

// NewCollector creates a new GODEBUG metrics collector
func NewCollector() *Collector {
	return &Collector{
		history: make([]domain.SchedData, 0, collector.MaxHistoryPoints),
	}
}

// Start begins collecting metrics by running the specified command with GODEBUG=schedtrace=1000
func (c *Collector) Start(ctx context.Context, targetCmd string) error {
	cmdParts := strings.Fields(targetCmd)
	if len(cmdParts) == 0 {
		return fmt.Errorf("empty command")
	}

	c.cmd = exec.Command(cmdParts[0], cmdParts[1:]...)
	c.cmd.Env = append(os.Environ(), "GODEBUG=schedtrace=1000")

	stderr, err := c.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("getting stderr pipe: %w", err)
	}

	if err := c.cmd.Start(); err != nil {
		return fmt.Errorf("starting command: %w", err)
	}

	// Parse stderr lines in a goroutine
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			if data, err := ParseSchedLine(scanner.Text()); err == nil && data != nil {
				c.update(*data)
			}
		}
	}()

	return nil
}

// Stop terminates the monitored command and stops collection
func (c *Collector) Stop() error {
	if c.cmd != nil && c.cmd.Process != nil {
		return c.cmd.Process.Kill()
	}
	return nil
}

// update adds new metric data to the collector state
func (c *Collector) update(data domain.SchedData) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.latest = data
	c.history = append(c.history, data)
	if len(c.history) > collector.MaxHistoryPoints {
		c.history = c.history[1:]
	}
}

// GetCurrent returns the most recent scheduler state
func (c *Collector) GetCurrent() domain.SchedData {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.latest
}

// GetHistory returns a copy of historical scheduler states
func (c *Collector) GetHistory() []domain.SchedData {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make([]domain.SchedData, len(c.history))
	copy(result, c.history)
	return result
}
