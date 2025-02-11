package godebug

import (
	"regexp"
	"strconv"
	"strings"

	"schedtrace-mon/internal/domain"
)

// Regular expression for SCHED trace lines with format:
// SCHED 2013ms: gomaxprocs=14 idleprocs=14 threads=22 spinningthreads=0 needspinning=0 idlethreads=17 runqueue=0 [0 0 0 0 0 0 0 0 0 0 0 0 0 0]
var schedRegex = regexp.MustCompile(
	`^SCHED\s+(\d+)ms:\s+` +
		`gomaxprocs=(\d+)\s+` +
		`idleprocs=(\d+)\s+` +
		`threads=(\d+)\s+` +
		`spinningthreads=(\d+)\s+` +
		`needspinning=(\d+)\s+` +
		`idlethreads=(\d+)\s+` +
		`runqueue=(\d+)\s+\[([0-9\s]+)\]`,
)

// ParseSchedLine parses a single GODEBUG=schedtrace output line into SchedData.
// Returns nil for both error and SchedData if the line doesn't match expected format.
func ParseSchedLine(line string) (*domain.SchedData, error) {
	matches := schedRegex.FindStringSubmatch(line)
	if len(matches) != 10 {
		return nil, nil // not a sched trace line - skip
	}

	timeMs, _ := strconv.Atoi(matches[1])
	gmp, _ := strconv.Atoi(matches[2])
	idleProcs, _ := strconv.Atoi(matches[3])
	threads, _ := strconv.Atoi(matches[4])
	spinTh, _ := strconv.Atoi(matches[5])
	needSpin, _ := strconv.Atoi(matches[6])
	idleThreads, _ := strconv.Atoi(matches[7])
	runQ, _ := strconv.Atoi(matches[8])

	// Parse Local Run Queues
	brackets := matches[9]
	fields := strings.Fields(brackets)
	lrqVals := make([]int, len(fields))
	sumLRQ := 0
	for i, s := range fields {
		if n, err := strconv.Atoi(s); err == nil {
			lrqVals[i] = n
			sumLRQ += n
		}
	}

	return &domain.SchedData{
		TimeMs:          timeMs,
		GoMaxProcs:      gmp,
		IdleProcs:       idleProcs,
		Threads:         threads,
		SpinningThreads: spinTh,
		NeedSpinning:    needSpin,
		IdleThreads:     idleThreads,
		RunQueue:        runQ,
		LrqSum:          sumLRQ,
		Lrq:             lrqVals,
	}, nil
}
