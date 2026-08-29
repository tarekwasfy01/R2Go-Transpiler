package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type module struct {
	Path    string
	Version string
	Dir     string
	Main    bool
	Replace *module
}

func main() {
	cmd := exec.Command("go", "list", "-m", "-json", "all")
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "license collector: go list failed:", err)
		os.Exit(1)
	}

	dec := json.NewDecoder(bytes.NewReader(out))
	var mods []module
	for {
		var m module
		err := dec.Decode(&m)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "license collector: decode:", err)
			os.Exit(1)
		}
		if m.Main {
			continue
		}
		if m.Replace != nil && m.Replace.Dir != "" {
			m.Dir = m.Replace.Dir
		}
		if m.Dir != "" {
			mods = append(mods, m)
		}
	}

	sort.Slice(mods, func(i, j int) bool { return mods[i].Path < mods[j].Path })
	var b strings.Builder
	b.WriteString("r2go - LICENSES AND THIRD PARTY NOTICES\n")
	b.WriteString("Generated at build time and embedded into r2go.exe.\n\n")

	found := 0
	// Include r2go's own license and project-level notices (including logo
	// licensing/trademark notices) in the embedded bundle as well.
	projectRoot := filepath.Join("..", "..")
	if data, err := os.ReadFile(filepath.Join(projectRoot, "LICENSE")); err == nil {
		found++
		writeNotice(&b, "r2go", "project", "LICENSE", data)
	}
	for _, path := range discoverLicenseFiles(filepath.Join(projectRoot, "licenses")) {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		found++
		rel, _ := filepath.Rel(projectRoot, path)
		writeNotice(&b, "r2go-assets", "project", filepath.ToSlash(rel), data)
	}

	for _, m := range mods {
		for _, path := range discoverLicenseFiles(m.Dir) {
			data, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			found++
			rel, _ := filepath.Rel(m.Dir, path)
			writeNotice(&b, m.Path, m.Version, filepath.ToSlash(rel), data)
		}
	}
	if found == 0 {
		fmt.Fprintln(os.Stderr, "license collector: no license/notice files found")
		os.Exit(1)
	}

	outPath := filepath.Join("generated", "THIRD_PARTY_NOTICES.txt")
	if err := os.WriteFile(outPath, []byte(b.String()), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "license collector: write:", err)
		os.Exit(1)
	}
	fmt.Printf("embedded license bundle updated: %s (%d files)\n", outPath, found)
}

func discoverLicenseFiles(root string) []string {
	var files []string
	// Most Go modules put licenses at the module root. We also inspect a few
	// conventional metadata directories so bundled notices are not missed.
	roots := []string{root, filepath.Join(root, "licenses"), filepath.Join(root, "LICENSES")}
	seen := map[string]bool{}
	for _, dir := range roots {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() || !looksLikeLicense(e.Name()) {
				continue
			}
			p := filepath.Join(dir, e.Name())
			if !seen[p] {
				seen[p] = true
				files = append(files, p)
			}
		}
	}
	sort.Strings(files)
	return files
}

func looksLikeLicense(name string) bool {
	n := strings.ToUpper(name)
	for _, prefix := range []string{"LICENSE", "LICENCE", "COPYING", "NOTICE", "AUTHORS", "COPYRIGHT"} {
		if n == prefix || strings.HasPrefix(n, prefix+".") || strings.HasPrefix(n, prefix+"-") || strings.HasPrefix(n, prefix+"_") {
			return true
		}
	}
	return false
}

func writeNotice(b *strings.Builder, modulePath, version, file string, data []byte) {
	fmt.Fprintf(b, "================================================================================\nMODULE: %s %s\nFILE: %s\n================================================================================\n", modulePath, version, file)
	b.Write(data)
	if len(data) == 0 || data[len(data)-1] != '\n' {
		b.WriteByte('\n')
	}
	b.WriteByte('\n')
}
