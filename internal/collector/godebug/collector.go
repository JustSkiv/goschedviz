// Package godebug implements scheduler metrics collection from GODEBUG=schedtrace output.
package godebug

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/yourusername/projectname/internal/domain"
)

// Collector implements collector.Collector interface for GODEBUG schedtrace output.
type Collector struct {
	cmd    *exec.Cmd
	done   chan struct{}
	path   string
	period int // schedtrace period in milliseconds
}

// New creates a new GODEBUG collector that will monitor the specified program.
func New(programPath string, tracePeriod int) *Collector {
	return &Collector{
		path:   programPath,
		period: tracePeriod,
		done:   make(chan struct{}),
	}
}

// Start implements collector.Collector interface.
func (c *Collector) Start(ctx context.Context) (<-chan domain.SchedulerSnapshot, error) {
	snapshots := make(chan domain.SchedulerSnapshot)

	c.cmd = exec.Command("go", "run", c.path)
	c.cmd.Env = append(os.Environ(), fmt.Sprintf("GODEBUG=schedtrace=%d", c.period))

	stderr, err := c.cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	if err := c.cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start process: %w", err)
	}

	go func() {
		defer close(snapshots)
		defer c.cmd.Process.Kill()

		scanner := bufio.NewScanner(stderr)
		parser := NewParser()

		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case <-c.done:
				return
			default:
				if snapshot, ok := parser.Parse(scanner.Text()); ok {
					snapshots <- snapshot
				}
			}
		}
	}()

	return snapshots, nil
}

// Stop implements collector.Collector interface.
func (c *Collector) Stop() error {
	close(c.done)
	if c.cmd != nil && c.cmd.Process != nil {
		return c.cmd.Process.Kill()
	}
	return nil
}
