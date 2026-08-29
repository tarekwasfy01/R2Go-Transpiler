//go:build windows

package platform

import (
	"path/filepath"
	"testing"
)

func TestSaveDialogAutomationPath(t *testing.T) {
	want := filepath.Join(t.TempDir(), "generated.go")
	t.Setenv("R2GO_SAVE_PATH", want)
	got, err := SaveGoFileDialog("output.go")
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("path=%q want %q", got, want)
	}
}

func TestSaveDialogRejectsWrongExtension(t *testing.T) {
	t.Setenv("R2GO_SAVE_PATH", filepath.Join(t.TempDir(), "generated.txt"))
	if _, err := SaveGoFileDialog("output.go"); err == nil {
		t.Fatal("expected a .go extension error")
	}
}
