//go:build windows

package platform

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unicode/utf16"
)

// SaveGoFileDialog runs the Windows Forms picker in a separate STA process.
// Native dialog faults therefore cannot take down the Gio process.
func SaveGoFileDialog(defaultName string) (string, error) {
	if override := os.Getenv("R2GO_SAVE_PATH"); override != "" {
		return validateGoSavePath(override)
	}
	if defaultName == "" {
		defaultName = "output.go"
	}
	safeName := strings.ReplaceAll(defaultName, "'", "''")
	script := `$ErrorActionPreference='Stop'
[Console]::OutputEncoding = [System.Text.UTF8Encoding]::new($false)
Add-Type -AssemblyName System.Windows.Forms
$dialog = New-Object System.Windows.Forms.SaveFileDialog
$dialog.Filter = 'Go source files (*.go)|*.go|All files (*.*)|*.*'
$dialog.DefaultExt = 'go'
$dialog.AddExtension = $true
$dialog.OverwritePrompt = $true
$dialog.FileName = '` + safeName + `'
try {
  if ($dialog.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) {
    [Console]::Out.Write($dialog.FileName)
  }
} finally {
  $dialog.Dispose()
}`
	cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-STA", "-WindowStyle", "Hidden", "-EncodedCommand", encodePowerShell(script))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("open Save As dialog: %w", err)
	}
	path := strings.TrimSpace(string(output))
	if path == "" {
		return "", nil
	}
	return validateGoSavePath(path)
}

func encodePowerShell(script string) string {
	words := utf16.Encode([]rune(script))
	data := make([]byte, len(words)*2)
	for i, word := range words {
		data[i*2] = byte(word)
		data[i*2+1] = byte(word >> 8)
	}
	return base64.StdEncoding.EncodeToString(data)
}

func validateGoSavePath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", nil
	}
	if filepath.Ext(path) == "" {
		path += ".go"
	}
	if !strings.EqualFold(filepath.Ext(path), ".go") {
		return "", errors.New("please choose a .go file")
	}
	return filepath.Clean(path), nil
}
