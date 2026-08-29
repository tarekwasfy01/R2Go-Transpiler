package runtime

// Host is the sole boundary between generated Pure-Go R programs and their
// environment.  Embedders can provide a sandbox, an in-memory filesystem, or
// a policy-enforcing host without changing primitive implementations.
import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Host interface {
	Getenv(string) string
	Setenv(string, string) error
	Unsetenv(string) error
	Now() time.Time
	Stdout() io.Writer
	Getwd() (string, error)
	Chdir(string) error
	ReadFile(string) ([]byte, error)
	WriteFile(string, []byte) error
	Exists(string) (exists, isDir bool)
	Mkdir(string, bool) error
}

type LocalHost struct{}

func (LocalHost) Getenv(name string) string          { return os.Getenv(name) }
func (LocalHost) Setenv(name, value string) error    { return os.Setenv(name, value) }
func (LocalHost) Unsetenv(name string) error         { return os.Unsetenv(name) }
func (LocalHost) Now() time.Time                     { return time.Now() }
func (LocalHost) Stdout() io.Writer                  { return os.Stdout }
func (LocalHost) Getwd() (string, error)             { return os.Getwd() }
func (LocalHost) Chdir(p string) error               { return os.Chdir(p) }
func (LocalHost) ReadFile(p string) ([]byte, error)  { return os.ReadFile(p) }
func (LocalHost) WriteFile(p string, b []byte) error { return os.WriteFile(p, b, 0666) }
func (LocalHost) Exists(p string) (bool, bool) {
	i, e := os.Stat(p)
	return e == nil, e == nil && i.IsDir()
}
func (LocalHost) Mkdir(p string, recursive bool) error {
	if recursive {
		return os.MkdirAll(p, 0777)
	}
	return os.Mkdir(p, 0777)
}

// MemoryHost is deterministic and side-effect free. It is used by the effect
// matrix and is available to transpiled applications that must not inherit the
// parent process environment.
type MemoryHost struct {
	Environment map[string]string
	Time        time.Time
	Output      io.Writer
	WorkingDir  string
	Files       map[string][]byte
	Directories map[string]bool
}

func NewMemoryHost() *MemoryHost {
	return &MemoryHost{Environment: map[string]string{}, Time: time.Unix(0, 0), WorkingDir: "/", Files: map[string][]byte{}, Directories: map[string]bool{"/": true}}
}
func (h *MemoryHost) Getenv(name string) string       { return h.Environment[name] }
func (h *MemoryHost) Setenv(name, value string) error { h.Environment[name] = value; return nil }
func (h *MemoryHost) Unsetenv(name string) error      { delete(h.Environment, name); return nil }
func (h *MemoryHost) Now() time.Time                  { return h.Time }
func (h *MemoryHost) Stdout() io.Writer               { return h.Output }
func (h *MemoryHost) path(p string) string {
	if !strings.HasPrefix(p, "/") {
		p = filepath.Join(h.WorkingDir, p)
	}
	return filepath.ToSlash(filepath.Clean(p))
}
func (h *MemoryHost) Getwd() (string, error) { return h.WorkingDir, nil }
func (h *MemoryHost) Chdir(p string) error {
	p = h.path(p)
	if !h.Directories[p] {
		return errors.New("directory does not exist")
	}
	h.WorkingDir = p
	return nil
}
func (h *MemoryHost) ReadFile(p string) ([]byte, error) {
	b, ok := h.Files[h.path(p)]
	if !ok {
		return nil, errors.New("file does not exist")
	}
	return append([]byte(nil), b...), nil
}
func (h *MemoryHost) WriteFile(p string, b []byte) error {
	p = h.path(p)
	parent := filepath.ToSlash(filepath.Dir(p))
	if !h.Directories[parent] {
		return errors.New("parent directory does not exist")
	}
	h.Files[p] = append([]byte(nil), b...)
	return nil
}
func (h *MemoryHost) Exists(p string) (bool, bool) {
	p = h.path(p)
	if h.Directories[p] {
		return true, true
	}
	_, ok := h.Files[p]
	return ok, false
}
func (h *MemoryHost) Mkdir(p string, recursive bool) error {
	p = h.path(p)
	parent := filepath.ToSlash(filepath.Dir(p))
	if !recursive && !h.Directories[parent] {
		return errors.New("parent directory does not exist")
	}
	if recursive {
		for parent != "." && !h.Directories[parent] {
			h.Directories[parent] = true
			parent = filepath.ToSlash(filepath.Dir(parent))
		}
	}
	h.Directories[p] = true
	return nil
}
