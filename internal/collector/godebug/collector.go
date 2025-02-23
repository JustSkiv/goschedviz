package godebug

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/JustSkiv/goschedviz/internal/domain"
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

// validateConfig checks if the collector configuration is valid
func (c *Collector) validateConfig() error {
	// Check period
	if c.period <= 0 {
		return fmt.Errorf("period must be positive, got %d", c.period)
	}

	// Check if path is empty
	if c.path == "" {
		return fmt.Errorf("program path cannot be empty")
	}

	// Check if path exists and is a file
	info, err := os.Stat(c.path)
	if err != nil {
		return fmt.Errorf("invalid program path: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("program path must be a file, got directory: %s", c.path)
	}

	// Verify it's a .go file
	ext := filepath.Ext(c.path)
	if ext != ".go" {
		return fmt.Errorf("program must be a .go file, got: %s", c.path)
	}

	return nil
}

// Start implements collector.Collector interface.
func (c *Collector) Start(ctx context.Context) (_ <-chan domain.SchedulerSnapshot, err error) {
	// Validate configuration
	if err := c.validateConfig(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	snapshots := make(chan domain.SchedulerSnapshot)

	// Get temporary binary path
	tmpBinary := "tmp_program"
	if runtime.GOOS == "windows" {
		tmpBinary += ".exe"
	}

	// Ensure cleanup in case of errors during setup
	defer func() {
		if err != nil {
			os.Remove(tmpBinary)
		}
	}()

	// First compile the program
	buildCmd := exec.Command("go", "build", "-o", tmpBinary, c.path)
	if err := buildCmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to build program: %w", err)
	}

	// Then run the compiled binary
	c.cmd = exec.Command("./" + tmpBinary)
	c.cmd.Env = append(os.Environ(), fmt.Sprintf("GODEBUG=schedtrace=%d", c.period))
	c.cmd.Stdin = os.Stdin

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
		defer os.Remove(tmpBinary)

		scanner := bufio.NewScanner(stderr)
		parser := NewParser()

		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case <-c.done:
				return
			default:
				line := scanner.Text()
				if snapshot, ok := parser.Parse(line); ok {
					select {
					case snapshots <- snapshot:
					case <-ctx.Done():
						return
					case <-c.done:
						return
					}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stderr: %v\n", err)
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
