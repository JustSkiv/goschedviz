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

type MockPresenter struct {
	done       chan struct{}
	updateFunc func(ui.UIData)
}

func (m *MockPresenter) Start() error { return nil }
func (m *MockPresenter) Stop()        {}
func (m *MockPresenter) Done() <-chan struct{} {
	return m.done
}

func (m *MockPresenter) Update(data ui.UIData) {
	if m.updateFunc != nil {
		m.updateFunc(data)
	}
}

func TestMonitorScheduler(t *testing.T) {
	mockCollector := &MockCollector{
		snapshots: make(chan domain.SchedulerSnapshot),
	}
	mockPresenter := &MockPresenter{
		done: make(chan struct{}),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var updates []ui.UIData
	mockPresenter.updateFunc = func(data ui.UIData) {
		updates = append(updates, data)
	}

	errChan := make(chan error)
	go func() {
		errChan <- monitorScheduler(ctx, mockCollector, mockPresenter)
	}()

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

	time.Sleep(600 * time.Millisecond)
	close(mockCollector.snapshots)

	err := <-errChan
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(updates), 1)
}

func TestMonitorScheduler_Scenarios(t *testing.T) {
	tests := []struct {
		name        string
		setupMocks  func(collector *MockCollector, presenter *MockPresenter)
		expectError bool
	}{
		{
			name: "collector_start_error",
			setupMocks: func(collector *MockCollector, presenter *MockPresenter) {
				collector.startError = fmt.Errorf("start failed")
			},
			expectError: true,
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
			name:       "context_cancelled",
			setupMocks: func(collector *MockCollector, presenter *MockPresenter) {},
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
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
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
				go func() {
					collector.snapshots <- domain.SchedulerSnapshot{TimeMs: 100}
					time.Sleep(100 * time.Millisecond)
					close(collector.snapshots)
				}()

				ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
				defer cancel()

				err := monitorScheduler(ctx, collector, presenter)
				require.NoError(t, err)
				assert.True(t, collector.stopCalled)
			},
		},
		{
			name: "presenter_signals_done",
			runTest: func(t *testing.T, collector *MockCollector, presenter *MockPresenter) {
				go func() {
					collector.snapshots <- domain.SchedulerSnapshot{TimeMs: 100}
					time.Sleep(100 * time.Millisecond)
					close(presenter.done)
				}()

				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()

				err := monitorScheduler(ctx, collector, presenter)
				assert.NoError(t, err)
			},
		},
		{
			name: "state_updates_continue",
			runTest: func(t *testing.T, collector *MockCollector, presenter *MockPresenter) {
				var updates []ui.UIData
				presenter.updateFunc = func(data ui.UIData) {
					updates = append(updates, data)
				}

				go func() {
					collector.snapshots <- domain.SchedulerSnapshot{TimeMs: 100}
					time.Sleep(600 * time.Millisecond)
					collector.snapshots <- domain.SchedulerSnapshot{TimeMs: 200}
				}()

				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()

				err := monitorScheduler(ctx, collector, presenter)
				require.NoError(t, err)
				assert.GreaterOrEqual(t, len(updates), 1)
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

func TestMonitorScheduler_ErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func(collector *MockCollector, presenter *MockPresenter)
		expectedError string
	}{
		{
			name: "collector_returns_nil_channel",
			setupMocks: func(c *MockCollector, p *MockPresenter) {
				c.snapshots = nil
				c.startError = nil
			},
			expectedError: "",
		},
		{
			name: "collector_returns_closed_channel",
			setupMocks: func(c *MockCollector, p *MockPresenter) {
				ch := make(chan domain.SchedulerSnapshot)
				close(ch)
				c.snapshots = ch
			},
			expectedError: "",
		},
		{
			name: "presenter_signals_done",
			setupMocks: func(c *MockCollector, p *MockPresenter) {
				c.snapshots = make(chan domain.SchedulerSnapshot)
				go func() {
					c.snapshots <- domain.SchedulerSnapshot{TimeMs: 100}
					time.Sleep(50 * time.Millisecond)
					close(p.done)
				}()
			},
			expectedError: "",
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

			ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
			defer cancel()

			err := monitorScheduler(ctx, mockCollector, mockPresenter)
			if tt.expectedError != "" {
				assert.ErrorContains(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMonitorScheduler_Metrics(t *testing.T) {
	tests := []struct {
		name     string
		metrics  []domain.SchedulerSnapshot
		validate func(t *testing.T, updates []ui.UIData)
	}{
		{
			name: "increasing_metrics",
			metrics: []domain.SchedulerSnapshot{
				{TimeMs: 100, RunQueue: 5, Goroutines: 100},
				{TimeMs: 200, RunQueue: 10, Goroutines: 200},
				{TimeMs: 300, RunQueue: 15, Goroutines: 300},
			},
			validate: func(t *testing.T, updates []ui.UIData) {
				require.GreaterOrEqual(t, len(updates), 1)
				last := updates[len(updates)-1]
				assert.Equal(t, 15, last.Current.RunQueue)
				assert.Equal(t, 300, last.Current.Goroutines)
			},
		},
		{
			name: "fluctuating_metrics",
			metrics: []domain.SchedulerSnapshot{
				{TimeMs: 100, RunQueue: 10, Threads: 5},
				{TimeMs: 200, RunQueue: 5, Threads: 10},
				{TimeMs: 300, RunQueue: 15, Threads: 7},
			},
			validate: func(t *testing.T, updates []ui.UIData) {
				require.GreaterOrEqual(t, len(updates), 1)
				last := updates[len(updates)-1]
				assert.Equal(t, 15, last.Current.RunQueue)
				assert.GreaterOrEqual(t, last.Gauges.GRQ.Max, 15)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var updates []ui.UIData
			mockCollector := &MockCollector{
				snapshots: make(chan domain.SchedulerSnapshot),
			}
			mockPresenter := &MockPresenter{
				done: make(chan struct{}),
				updateFunc: func(data ui.UIData) {
					updates = append(updates, data)
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			errCh := make(chan error)
			go func() {
				errCh <- monitorScheduler(ctx, mockCollector, mockPresenter)
			}()

			for _, m := range tt.metrics {
				mockCollector.snapshots <- m
				time.Sleep(200 * time.Millisecond)
			}

			time.Sleep(200 * time.Millisecond)
			cancel()
			err := <-errCh
			assert.NoError(t, err)

			tt.validate(t, updates)
		})
	}
}
