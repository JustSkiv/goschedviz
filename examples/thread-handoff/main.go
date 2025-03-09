// Package threadhandoff demonstrates Go scheduler's thread handoff mechanism.
//
// Designed for analysis with goschedviz, this program shows how the Go scheduler
// handles blocking system calls through thread handoff.
//
// The stdin Read() syscall forces the runtime to detach P from M, allowing
// other goroutines to continue execution while the original M remains blocked.
package threadhandoff

import (
	"runtime"
	"time"

	"golang.org/x/sys/unix"
)

func main() {
	// for better visualization
	runtime.GOMAXPROCS(5)

	// Create goroutines that block on I/O
	for i := 0; i < 250; i++ {
		go func() {
			var buf [1]byte

			// Blocking syscall that triggers thread handoff:
			// M gets blocked, P is handed to another M
			unix.Read(unix.Stdin, buf[:])
		}()
	}

	// Keep the main goroutine alive for observation
	time.Sleep(1 * time.Hour)

	// With goschedviz you'll observe thread handoff:
	// OS threads growing beyond GOMAXPROCS as goroutines block
}
