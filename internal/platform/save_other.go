//go:build !windows

package platform

import (
	"fmt"
	"path/filepath"
)

func SaveGoFileDialog(defaultName string) (string, error) {
	if defaultName == "" {
		defaultName = "output.go"
	}
	path, err := filepath.Abs(defaultName)
	if err != nil {
		return "", fmt.Errorf("save path: %w", err)
	}
	return path, nil
}
