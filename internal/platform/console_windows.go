//go:build windows

package platform

import (
	"os"
	"syscall"
)

const attachParentProcess = uintptr(0xFFFFFFFF)

var (
	kernel32          = syscall.NewLazyDLL("kernel32.dll")
	procAttachConsole = kernel32.NewProc("AttachConsole")
)

// EnsureCLIConsole lets a -H=windowsgui build still behave like a CLI when it
// is started from a terminal. If stdout is already redirected/usable, it is
// left untouched.
func EnsureCLIConsole() {
	if _, err := os.Stdout.Stat(); err == nil {
		return
	}
	_, _, _ = procAttachConsole.Call(attachParentProcess)
	if f, err := os.OpenFile("CONOUT$", os.O_WRONLY, 0); err == nil {
		os.Stdout = f
	}
	if f, err := os.OpenFile("CONOUT$", os.O_WRONLY, 0); err == nil {
		os.Stderr = f
	}
	if f, err := os.OpenFile("CONIN$", os.O_RDONLY, 0); err == nil {
		os.Stdin = f
	}
}
