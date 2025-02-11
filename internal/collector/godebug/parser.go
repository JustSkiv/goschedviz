package godebug

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/yourusername/projectname/internal/domain"
)

// Parser handles parsing of GODEBUG schedtrace output lines.
type Parser struct {
	regex *regexp.Regexp
}

// NewParser creates a new GODEBUG output parser.
func NewParser() *Parser {
	// Regex for lines like:
	// SCHED 2013ms: gomaxprocs=14 idleprocs=14 threads=22 spinningthreads=0 needspinning=0 idlethreads=17 runqueue=0 [0 0 0 0 0 0 0 0 0 0 0 0 0 0]
	return &Parser{
		regex: regexp.MustCompile(
			`^SCHED\s+(\d+)ms:\s+` +
				`gomaxprocs=(\d+)\s+` +
				`idleprocs=(\d+)\s+` +
				`threads=(\d+)\s+` +
				`spinningthreads=(\d+)\s+` +
				`needspinning=(\d+)\s+` +
				`idlethreads=(\d+)\s+` +
				`runqueue=(\d+)\s+\[([0-9\s]+)\]`,
		),
	}
}

// Parse attempts to parse a single line of schedtrace output.
// Returns the parsed snapshot and true if successful, or zero value and false if parsing failed.
func (p *Parser) Parse(line string) (domain.SchedulerSnapshot, bool) {
	matches := p.regex.FindStringSubmatch(line)
	if len(matches) != 10 {
		return domain.SchedulerSnapshot{}, false
	}

	timeMs, _ := strconv.Atoi(matches[1])
	gmp, _ := strconv.Atoi(matches[2])
	idleProcs, _ := strconv.Atoi(matches[3])
	threads, _ := strconv.Atoi(matches[4])
	spinTh, _ := strconv.Atoi(matches[5])
	needSpin, _ := strconv.Atoi(matches[6])
	idleThreads, _ := strconv.Atoi(matches[7])
	runQ, _ := strconv.Atoi(matches[8])

	// Parse LRQ values
	fields := strings.Fields(matches[9])
	lrqVals := make([]int, len(fields))
	sumLRQ := 0
	for i, s := range fields {
		if n, err := strconv.Atoi(s); err == nil {
			lrqVals[i] = n
			sumLRQ += n
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
