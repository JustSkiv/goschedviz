package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JustSkiv/goschedviz/internal/domain"
	"github.com/JustSkiv/goschedviz/internal/ui"
)

type MockCollector struct {
	snapshots  chan domain.SchedulerSnapshot
	startError error
	stopCalled bool
}

func (m *MockCollector) Start(ctx context.Context) (<-chan domain.SchedulerSnapshot, error) {
	if m.startError != nil {
		return nil, m.startError
	}
	return m.snapshots, nil
}

func (m *MockCollector) Stop() error {
	m.stopCalled = true
	return nil
}

// Add to MockPresenter
type MockPresenter struct {
	done        chan struct{}
	updateFunc  func(ui.UIData)
	updateError error
}

func (m *MockPresenter) Start() error {
	return nil
}

func (m *MockPresenter) Stop() {}

func (m *MockPresenter) Update(data ui.UIData) {
	if m.updateError != nil {
		return
	}
	if m.updateFunc != nil {
		m.updateFunc(data)
	}
}

func (m *MockPresenter) Done() <-chan struct{} {
	return m.done
}

func TestMonitorScheduler(t *testing.T) {
	// Setup mocks
	mockCollector := &MockCollector{
		snapshots: make(chan domain.SchedulerSnapshot),
	}
	mockPresenter := &MockPresenter{
		done: make(chan struct{}),
	}

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Track UI updates
	var updates []ui.UIData
	mockPresenter.updateFunc = func(data ui.UIData) {
		updates = append(updates, data)
	}

	// Start monitoring in background
	errChan := make(chan error)
	go func() {
		errChan <- monitorScheduler(ctx, mockCollector, mockPresenter)
	}()

	// Send test data
	testData := []domain.SchedulerSnapshot{
		{
			TimeMs:     100,
			GoMaxProcs: 2,
			RunQueue:   1,
			LRQ:        []int{1, 1},
			LRQSum:     2,
		},
		{
			TimeMs:     200,
			GoMaxProcs: 2,
			RunQueue:   2,
			LRQ:        []int{2, 2},
			LRQSum:     4,
		},
	}

	for _, data := range testData {
		mockCollector.snapshots <- data
	}

	// Wait for at least one UI update (ticker period is 500ms)
	time.Sleep(600 * time.Millisecond)

	// Clean shutdown
	close(mockCollector.snapshots)

	// Wait for completion
	err := <-errChan
	require.NoError(t, err)

	// Verify updates
	assert.GreaterOrEqual(t, len(updates), 1, "should receive UI updates")
	if len(updates) > 0 {
		lastUpdate := updates[len(updates)-1]
		assert.Equal(t, testData[1].TimeMs, lastUpdate.Current.TimeMs)
		assert.Equal(t, testData[1].GoMaxProcs, lastUpdate.Current.GoMaxProcs)
		assert.Equal(t, testData[1].RunQueue, lastUpdate.Current.RunQueue)
		assert.Equal(t, testData[1].LRQ, lastUpdate.Current.LRQ)
		assert.Equal(t, testData[1].LRQSum, lastUpdate.Current.LRQSum)
	}
}

func TestMonitorScheduler_Scenarios(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func(collector *MockCollector, presenter *MockPresenter)
		expectError   bool
		expectedError string
	}{
		{
			name: "collector_start_error",
			setupMocks: func(collector *MockCollector, presenter *MockPresenter) {
				collector.startError = fmt.Errorf("failed to start")
			},
			expectError:   true,
			expectedError: "failed to start collector: failed to start",
		},
		{
			name: "presenter_done",
			setupMocks: func(collector *MockCollector, presenter *MockPresenter) {
				go func() {
					time.Sleep(100 * time.Millisecond)
					close(presenter.done)
				}()
			},
		},
		{
			name: "context_cancelled",
			setupMocks: func(collector *MockCollector, presenter *MockPresenter) {
				// Context will be cancelled by the test timeout
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCollector := &MockCollector{
				snapshots: make(chan domain.SchedulerSnapshot),
			}
			mockPresenter := &MockPresenter{
				done: make(chan struct{}),
			}

			if tt.setupMocks != nil {
				tt.setupMocks(mockCollector, mockPresenter)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
			defer cancel()

			err := monitorScheduler(ctx, mockCollector, mockPresenter)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMonitorScheduler_ErrorsAndCleanup(t *testing.T) {
	tests := []struct {
		name    string
		runTest func(t *testing.T, collector *MockCollector, presenter *MockPresenter)
	}{
		{
			name: "collector_stop_called",
			runTest: func(t *testing.T, collector *MockCollector, presenter *MockPresenter) {
				// Send some data then cancel context
				go func() {
					collector.snapshots <- domain.SchedulerSnapshot{TimeMs: 100}
					time.Sleep(100 * time.Millisecond)
					close(collector.snapshots)
				}()

				ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
				defer cancel()

				err := monitorScheduler(ctx, collector, presenter)
				require.NoError(t, err)
				assert.True(t, collector.stopCalled, "Stop should be called")
			},
		},
		{
			name: "presenter_error_handled",
			runTest: func(t *testing.T, collector *MockCollector, presenter *MockPresenter) {
				presenter.updateError = fmt.Errorf("update failed")

				go func() {
					collector.snapshots <- domain.SchedulerSnapshot{TimeMs: 100}
					time.Sleep(600 * time.Millisecond) // Wait for ticker
				}()

				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()

				err := monitorScheduler(ctx, collector, presenter)
				require.NoError(t, err) // Should continue despite presenter errors
			},
		},
		{
			name: "state_updates_after_error",
			runTest: func(t *testing.T, collector *MockCollector, presenter *MockPresenter) {
				var updates []ui.UIData
				presenter.updateFunc = func(data ui.UIData) {
					updates = append(updates, data)
				}

				// Send data, trigger error, send more data
				go func() {
					collector.snapshots <- domain.SchedulerSnapshot{TimeMs: 100}
					presenter.updateError = fmt.Errorf("temporary error")
					time.Sleep(100 * time.Millisecond)
					presenter.updateError = nil
					collector.snapshots <- domain.SchedulerSnapshot{TimeMs: 200}
					time.Sleep(600 * time.Millisecond) // Wait for ticker
				}()

				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()

				err := monitorScheduler(ctx, collector, presenter)
				require.NoError(t, err)
				assert.GreaterOrEqual(t, len(updates), 1, "should receive updates after error recovery")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCollector := &MockCollector{
				snapshots: make(chan domain.SchedulerSnapshot),
			}
			mockPresenter := &MockPresenter{
				done: make(chan struct{}),
			}

			tt.runTest(t, mockCollector, mockPresenter)
		})
	}
}
