package godebug

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/yourusername/projectname/internal/domain"
)

// Parser handles parsing of GODEBUG schedtrace output lines.
type Parser struct {
	// Example of schedtrace output:
	// SCHED 2013ms: gomaxprocs=14 idleprocs=14 threads=22 spinningthreads=0 needspinning=0 idlethreads=17 runqueue=0 [0 0 0 0 0 0 0 0 0 0 0 0 0 0]
	regex *regexp.Regexp
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

// Parse attempts to parse a single line of schedtrace output.
// Returns the parsed snapshot and true if successful, or zero value and false if parsing failed.
func (p *Parser) Parse(line string) (domain.SchedulerSnapshot, bool) {
	matches := p.regex.FindStringSubmatch(line)
	if len(matches) != 10 { // 1 full match + 9 groups
		return domain.SchedulerSnapshot{}, false
	}

	// Parse all integer values
	timeMs, err := strconv.Atoi(matches[1])
	if err != nil {
		fmt.Printf("DEBUG: Failed to parse TimeMs: %v\n", err)
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

	return domain.SchedulerSnapshot{
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
	}, true
}
