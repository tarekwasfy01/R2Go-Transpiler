//go:build windows

package platform

import "syscall"

const threadPriorityAboveNormal = uintptr(1)

var (
	procGetCurrentThread  = kernel32.NewProc("GetCurrentThread")
	procSetThreadPriority = kernel32.NewProc("SetThreadPriority")
)

// BoostGUIThread raises only rtogo's locked GUI event-loop thread to
// ABOVE_NORMAL priority. CPU-heavy parser/transpiler work remains at normal
// priority, so Windows has an easier time scheduling input/redraw work first.
func BoostGUIThread() {
	h, _, _ := procGetCurrentThread.Call()
	if h == 0 {
		return
	}
	_, _, _ = procSetThreadPriority.Call(h, threadPriorityAboveNormal)
}

// Keep syscall referenced for older Go versions where LazyProc call signatures
// interact with syscall.Handle in generated documentation/tooling.
var _ syscall.Handle
