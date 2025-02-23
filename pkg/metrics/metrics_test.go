package metrics

import (
	"bufio"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReporter_Output(t *testing.T) {
	// Create a pipe to capture stderr
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stderr = w

	// Restore stderr when test completes
	defer func() {
		os.Stderr = oldStderr
	}()

	// Create and start reporter with short interval
	reporter := NewReporter(100 * time.Millisecond)
	reporter.Start()

	// Let it generate some output
	time.Sleep(250 * time.Millisecond)
	reporter.Stop()

	// Close write end of pipe to unblock reads
	w.Close()

	// Read captured output
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	require.NoError(t, scanner.Err())

	// Verify output format
	require.NotEmpty(t, lines, "Should have captured some output")
	for _, line := range lines {
		assert.True(t, strings.HasPrefix(line, prefix),
			"Each line should start with metrics prefix")
		assert.Contains(t, line, "num_goroutines=",
			"Each line should contain goroutines metric")
	}
}

func TestReporter_StartStop(t *testing.T) {
	reporter := NewReporter(time.Second)

	// Should not panic on multiple starts
	reporter.Start()
	reporter.Start()

	// Should not panic on multiple stops
	reporter.Stop()
	reporter.Stop()
}
