package godebug

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JustSkiv/goschedviz/internal/domain"
)

// TestProgram is a simple Go program that will generate GODEBUG schedtrace output.
// It creates lots of goroutines and CPU activity to ensure scheduler events.
const TestProgram = `package main

import (
	"runtime"
	"time"
)

func main() {
	// Force multiple processors for more scheduler activity
	runtime.GOMAXPROCS(4)
	
	// Create worker pool
	jobs := make(chan int, 100)
	results := make(chan int, 100)
	
	// Start workers
	for i := 0; i < 20; i++ {
		go func() {
			for j := range jobs {
				// CPU intensive work
				result := 0
				for k := 0; k < j*1000; k++ {
					result += k
				}
				results <- result
				runtime.Gosched() // Explicitly yield
			}
		}()
	}
	
	// Send jobs
	go func() {
		for i := 0; i < 1000; i++ {
			jobs <- i
			time.Sleep(time.Millisecond)
		}
		close(jobs)
	}()
	
	// Main goroutine collects results
	timeout := time.After(3 * time.Second)
	count := 0
	for {
		select {
		case <-results:
			count++
		case <-timeout:
			return
		}
	}
}`

func setupTestProgram(t *testing.T) string {
	t.Helper()

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "collector-test-*")
	require.NoError(t, err, "failed to create temp dir")

	// Create test program file
	programPath := filepath.Join(tmpDir, "test.go")
	err = os.WriteFile(programPath, []byte(TestProgram), 0666)
	require.NoError(t, err, "failed to write test program")

	// Return path and cleanup function
	t.Cleanup(func() { os.RemoveAll(tmpDir) })
	return programPath
}

func TestCollector_StartStop(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skip("Skipping test on non-Unix platform")
	}

	programPath := setupTestProgram(t)
	collector := New(programPath, 100) // 100ms period
	require.NotNil(t, collector)

	// Start collector
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	snapshots, err := collector.Start(ctx)
	require.NoError(t, err, "Start should not return error")
	require.NotNil(t, snapshots, "Snapshots channel should not be nil")

	// Read some snapshots
	var count int
	var lastSnapshot domain.SchedulerSnapshot
	timeout := time.After(2 * time.Second)

readLoop:
	for {
		select {
		case snapshot, ok := <-snapshots:
			if !ok {
				break readLoop
			}
			lastSnapshot = snapshot

			// Basic validation of snapshot data
			assert.GreaterOrEqual(t, snapshot.GoMaxProcs, 1, "GoMaxProcs should be >= 1")
			assert.Len(t, snapshot.LRQ, snapshot.GoMaxProcs, "LRQ length should match GoMaxProcs")

			// Count only non-zero time samples
			if snapshot.TimeMs > 0 {
				count++
				// Break after we get enough non-zero samples
				if count >= 3 {
					break readLoop
				}
			}

		case <-timeout:
			t.Log("Test timed out waiting for snapshots")
			break readLoop
		}
	}

	// We should receive at least a few snapshots
	assert.Greater(t, count, 0, "Should receive some snapshots")
	t.Logf("Received %d snapshots, last TimeMs: %d", count, lastSnapshot.TimeMs)

	// Stop collector
	err = collector.Stop()
	assert.NoError(t, err, "Stop should not return error")

	// Channel should be closed
	_, ok := <-snapshots
	assert.False(t, ok, "Channel should be closed after Stop")
}

func TestCollector_InvalidProgram(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skip("Skipping test on non-Unix platform")
	}

	tests := []struct {
		name        string
		programPath string
		wantErr     bool
	}{
		{
			name:        "non-existent program",
			programPath: "/path/to/nonexistent/program.go",
			wantErr:     true,
		},
		{
			name:        "empty path",
			programPath: "",
			wantErr:     true,
		},
		{
			name:        "directory instead of file",
			programPath: os.TempDir(),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := New(tt.programPath, 100)
			require.NotNil(t, collector)

			ctx := context.Background()
			snapshots, err := collector.Start(ctx)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, snapshots)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, snapshots)
				collector.Stop()
			}
		})
	}
}

func TestCollector_ContextCancellation(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skip("Skipping test on non-Unix platform")
	}

	programPath := setupTestProgram(t)
	collector := New(programPath, 100)
	require.NotNil(t, collector)

	// Create cancelable context
	ctx, cancel := context.WithCancel(context.Background())

	// Start collector
	snapshots, err := collector.Start(ctx)
	require.NoError(t, err)
	require.NotNil(t, snapshots)

	// Cancel context after short delay
	time.Sleep(200 * time.Millisecond)
	cancel()

	// Channel should be closed soon
	timeout := time.After(1 * time.Second)
	select {
	case _, ok := <-snapshots:
		if !ok {
			// Channel closed as expected
			return
		}
	case <-timeout:
		t.Fatal("Channel was not closed after context cancellation")
	}
}

func TestCollector_PeriodValidation(t *testing.T) {
	tests := []struct {
		name    string
		period  int
		wantErr bool
	}{
		{"zero period", 0, true},
		{"negative period", -100, true},
		{"valid period", 100, false},
		{"large period", 5000, false},
	}

	programPath := setupTestProgram(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := New(programPath, tt.period)
			require.NotNil(t, collector)

			ctx := context.Background()
			snapshots, err := collector.Start(ctx)
			defer collector.Stop()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, snapshots)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, snapshots)
			}
		})
	}
}
