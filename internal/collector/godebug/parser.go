package godebug

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/JustSkiv/goschedviz/internal/domain"
)

// Parser handles parsing of GODEBUG schedtrace output.
type Parser struct {
	// Example of schedtrace output:
	// SCHED 2013ms: gomaxprocs=14 idleprocs=14 threads=22 spinningthreads=0 needspinning=0 idlethreads=17 runqueue=0 [0 0 0 0 0 0 0 0 0 0 0 0 0 0]
	regex *regexp.Regexp
	// lastGoroutines holds the last seen goroutines count from metrics
	lastGoroutines int
}

// parseMetrics attempts to parse process metrics line.
// Returns number of goroutines if successful, or -1 if parsing failed.
func (p *Parser) parseMetrics(line string) int {
	if !strings.HasPrefix(line, "PROCMETR") {
		return -1
	}

	// Extract goroutines count from format "PROCMETR num_goroutines=1234"
	parts := strings.Split(line, "=")
	if len(parts) != 2 {
		return -1
	}

	count, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return -1
	}

	return count
}

// NewParser creates a new GODEBUG output parser.
func NewParser() *Parser {
	return &Parser{
		regex: regexp.MustCompile(
			// Separate into capturing groups carefully
			`^SCHED\s+` + // Prefix
				`(\d+)ms:\s+` + // TimeMs (group 1)
				`gomaxprocs=(\d+)\s+` + // GoMaxProcs (group 2)
				`idleprocs=(\d+)\s+` + // IdleProcs (group 3)
				`threads=(\d+)\s+` + // Threads (group 4)
				`spinningthreads=(\d+)\s+` + // SpinningThreads (group 5)
				`needspinning=(\d+)\s+` + // NeedSpinning (group 6)
				`idlethreads=(\d+)\s+` + // IdleThreads (group 7)
				`runqueue=(\d+)\s+` + // RunQueue (group 8)
				`\[([\d\s]+)\]`, // LRQ values (group 9)
		),
	}
}

// isValidSnapshot performs additional validation of the scheduler snapshot data.
// Returns false if any of the validation rules fail.
func (p *Parser) isValidSnapshot(s domain.SchedulerSnapshot) bool {
	// Check for valid GOMAXPROCS value
	if s.GoMaxProcs <= 0 {
		return false
	}

	// Verify that idle processors count doesn't exceed total processors
	if s.IdleProcs > s.GoMaxProcs {
		return false
	}

	// Verify that LRQ length matches GOMAXPROCS
	if len(s.LRQ) != s.GoMaxProcs {
		return false
	}

	// Validate thread counts consistency
	if s.Threads < s.SpinningThreads || s.Threads < s.IdleThreads {
		return false
	}

	// Calculate and validate queue totals
	calculatedSum := s.RunQueue
	for _, v := range s.LRQ {
		if v < 0 { // Queue lengths cannot be negative
			return false
		}
		calculatedSum += v
	}

	// Suspect data detection: all zeros might indicate invalid trace output
	if calculatedSum == 0 && s.IdleProcs == 0 && s.SpinningThreads == 0 {
		return false
	}

	return true
}

// Parse attempts to parse a single line of schedtrace output.
// Returns the parsed snapshot and true if successful, or zero value and false if parsing failed.
func (p *Parser) Parse(line string) (domain.SchedulerSnapshot, bool) {
	// Try to parse metrics line first
	if goroutines := p.parseMetrics(line); goroutines >= 0 {
		p.lastGoroutines = goroutines
		return domain.SchedulerSnapshot{}, false
	}

	matches := p.regex.FindStringSubmatch(line)
	if len(matches) != 10 { // 1 full match + 9 groups
		return domain.SchedulerSnapshot{}, false
	}

	// Parse all integer values
	timeMs, err := strconv.Atoi(matches[1])
	if err != nil {
		return domain.SchedulerSnapshot{}, false
	}

	gmp, err := strconv.Atoi(matches[2])
	if err != nil {
		return domain.SchedulerSnapshot{}, false
	}

	idleProcs, err := strconv.Atoi(matches[3])
	if err != nil {
		return domain.SchedulerSnapshot{}, false
	}

	threads, err := strconv.Atoi(matches[4])
	if err != nil {
		return domain.SchedulerSnapshot{}, false
	}

	spinTh, err := strconv.Atoi(matches[5])
	if err != nil {
		return domain.SchedulerSnapshot{}, false
	}

	needSpin, err := strconv.Atoi(matches[6])
	if err != nil {
		return domain.SchedulerSnapshot{}, false
	}

	idleThreads, err := strconv.Atoi(matches[7])
	if err != nil {
		return domain.SchedulerSnapshot{}, false
	}

	runQ, err := strconv.Atoi(matches[8])
	if err != nil {
		return domain.SchedulerSnapshot{}, false
	}

	// Parse LRQ values
	fields := strings.Fields(matches[9])
	lrqVals := make([]int, len(fields))
	sumLRQ := 0
	for i, s := range fields {
		if n, err := strconv.Atoi(s); err == nil {
			lrqVals[i] = n
			sumLRQ += n
		} else {
			return domain.SchedulerSnapshot{}, false
		}
	}

	snapshot := domain.SchedulerSnapshot{
		TimeMs:          timeMs,
		GoMaxProcs:      gmp,
		IdleProcs:       idleProcs,
		Threads:         threads,
		SpinningThreads: spinTh,
		NeedSpinning:    needSpin,
		IdleThreads:     idleThreads,
		RunQueue:        runQ,
		LRQSum:          sumLRQ,
		LRQ:             lrqVals,
		Goroutines:      p.lastGoroutines,
	}

	// Perform additional validation of the snapshot
	if !p.isValidSnapshot(snapshot) {
		return domain.SchedulerSnapshot{}, false
	}

	return snapshot, true
}
