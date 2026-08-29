// Package standalone writes self-contained R2Go output projects.
package standalone

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// runtimeSources is compiled into R2Go.exe. Generated projects therefore need
// neither an installed R2Go module nor a network connection.
//
//go:embed runtime/*.go syntax/*.go cir/*.go converted/rgo/*.go converted/rprimitive/*.go LICENSE
var runtimeSources embed.FS

var moduleLine = regexp.MustCompile(`(?m)^\s*module\s+([^\s]+)\s*$`)

// WriteProgram writes generated Go plus the complete Pure-Go compatibility
// runtime into a local r2go_runtime source tree.
func WriteProgram(outputPath string, source []byte) error {
	absOutput, err := filepath.Abs(outputPath)
	if err != nil {
		return fmt.Errorf("resolve output: %w", err)
	}
	outputDir := filepath.Dir(absOutput)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	moduleRoot, modulePath, err := findOrCreateModule(outputDir, absOutput)
	if err != nil {
		return err
	}
	runtimeRoot := filepath.Join(outputDir, "r2go_runtime")
	relRuntime, err := filepath.Rel(moduleRoot, runtimeRoot)
	if err != nil || relRuntime == ".." || strings.HasPrefix(relRuntime, ".."+string(filepath.Separator)) {
		return fmt.Errorf("runtime output is outside generated module")
	}
	importPrefix := strings.TrimSuffix(modulePath, "/") + "/" + filepath.ToSlash(relRuntime)
	if err := materializeRuntime(runtimeRoot, importPrefix); err != nil {
		return err
	}
	if err := os.WriteFile(absOutput, rewriteImports(source, importPrefix), 0644); err != nil {
		return fmt.Errorf("write generated program: %w", err)
	}
	return nil
}

func findOrCreateModule(outputDir, outputPath string) (string, string, error) {
	for dir := outputDir; ; dir = filepath.Dir(dir) {
		goModPath := filepath.Join(dir, "go.mod")
		if data, err := os.ReadFile(goModPath); err == nil {
			match := moduleLine.FindSubmatch(data)
			if len(match) != 2 {
				return "", "", fmt.Errorf("cannot read module name from %s", goModPath)
			}
			return dir, string(match[1]), nil
		} else if !os.IsNotExist(err) {
			return "", "", fmt.Errorf("read %s: %w", goModPath, err)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
	}
	name := strings.TrimSuffix(filepath.Base(outputPath), filepath.Ext(outputPath))
	modulePath := "r2go.generated/" + safeName(name)
	if err := os.WriteFile(filepath.Join(outputDir, "go.mod"), []byte("module "+modulePath+"\n\ngo 1.26\n"), 0644); err != nil {
		return "", "", fmt.Errorf("write generated go.mod: %w", err)
	}
	return outputDir, modulePath, nil
}

func safeName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	var b strings.Builder
	for _, r := range name {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || strings.ContainsRune("-_.", r) {
			b.WriteRune(r)
		} else {
			b.WriteByte('-')
		}
	}
	if clean := strings.Trim(b.String(), "-_."); clean != "" {
		return clean
	}
	return "program"
}

func rewriteImports(source []byte, importPrefix string) []byte {
	return []byte(strings.ReplaceAll(string(source), `"r2go/`, `"`+importPrefix+`/`))
}

func materializeRuntime(runtimeRoot, importPrefix string) error {
	return fs.WalkDir(runtimeSources, ".", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || strings.HasSuffix(path, "_test.go") || path == "standalone.go" {
			return nil
		}
		data, err := runtimeSources.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read embedded %s: %w", path, err)
		}
		if strings.HasSuffix(path, ".go") {
			data = rewriteImports(data, importPrefix)
		}
		target := filepath.Join(runtimeRoot, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return fmt.Errorf("create runtime directory: %w", err)
		}
		if err := os.WriteFile(target, data, 0644); err != nil {
			return fmt.Errorf("write embedded %s: %w", path, err)
		}
		return nil
	})
}
