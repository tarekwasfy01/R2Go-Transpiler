package runtime

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func init() {
	for _, entry := range []string{"do_Cstack_info", "do_Rhome", "do_glob", "do_localeconv", "do_setlocale", "do_syssleep", "do_capabilities", "do_capabilitiesX11", "do_compilerVersion", "do_date", "do_eSoftVersion", "do_iconv", "do_interactive", "do_l10n_info", "do_pathexpand", "do_proctime", "do_setTimeLimit", "do_setSessionTimeLimit", "do_syschmod", "do_unlink", "do_syswhich"} {
		registerLoweringKernel(entry, "0", kernelSystemInfo)
	}
}
func kernelSystemInfo(c *Context, f *LoweringFrame) error {
	switch f.Plan.CEntry {
	case "do_date":
		f.Result = &CharacterVector{Data: []string{c.Host.Now().Format("Mon Jan 02 15:04:05 2006")}}
	case "do_Rhome":
		f.Result = &CharacterVector{Data: []string{"pure-go-r2go"}}
	case "do_interactive":
		f.Result = &LogicalVector{Data: []Logical{False}}
	case "do_compilerVersion":
		f.Result = &CharacterVector{Data: []string{runtime.Version()}}
	case "do_eSoftVersion":
		f.Result = &List{Names: []string{"Go"}, Data: []Value{&CharacterVector{Data: []string{runtime.Version()}}}}
	case "do_Cstack_info":
		f.Result = &List{Names: []string{"size", "current", "direction"}, Data: []Value{&DoubleVector{Data: []float64{0}}, &DoubleVector{Data: []float64{0}}, &IntegerVector{Data: []int64{1}}}}
	case "do_proctime":
		f.Result = &DoubleVector{Data: []float64{time.Since(c.Started).Seconds(), 0, time.Since(c.Started).Seconds(), 0, 0}}
	case "do_capabilities":
		f.Result = &LogicalVector{Data: []Logical{False}}
	case "do_capabilitiesX11":
		f.Result = &LogicalVector{Data: []Logical{False}}
	case "do_localeconv":
		f.Result = &List{Names: []string{"decimal_point", "thousands_sep"}, Data: []Value{&CharacterVector{Data: []string{"."}}, &CharacterVector{Data: []string{""}}}}
	case "do_setlocale":
		previous := c.Locale
		if len(f.Arguments) > 1 {
			if locale, e := frameText(c, f, 1); e == nil && locale != "" {
				c.Locale = locale
			}
		}
		f.Result = &CharacterVector{Data: []string{previous}}
	case "do_setTimeLimit", "do_setSessionTimeLimit":
		f.Result = NullValue
	case "do_l10n_info":
		f.Result = &List{Names: []string{"MBCS", "UTF-8", "Latin-1"}, Data: []Value{&LogicalVector{Data: []Logical{True}}, &LogicalVector{Data: []Logical{True}}, &LogicalVector{Data: []Logical{False}}}}
	case "do_pathexpand":
		v, e := frameText(c, f, 0)
		if e != nil {
			return e
		}
		f.Result = &CharacterVector{Data: []string{filepath.Clean(v)}}
	case "do_syswhich":
		v, e := frameValue(c, f, 0)
		if e != nil {
			return e
		}
		paths, e := characterData(v)
		if e != nil {
			return e
		}
		out := &CharacterVector{Data: make([]string, len(paths.Data)), Attr: map[string]Value{"names": &CharacterVector{Data: append([]string(nil), paths.Data...)}}}
		for i, name := range paths.Data {
			if found, err := exec.LookPath(name); err == nil {
				out.Data[i] = found
			}
		}
		f.Result = out
	case "do_glob":
		pattern, e := frameText(c, f, 0)
		if e != nil {
			return e
		}
		matches, e := filepath.Glob(pattern)
		if e != nil {
			return e
		}
		f.Result = &CharacterVector{Data: matches}
	case "do_iconv":
		v, e := frameValue(c, f, 0)
		if e != nil {
			return e
		}
		f.Result = v
	case "do_syssleep":
		v, e := frameValue(c, f, 0)
		if e != nil {
			return e
		}
		seconds, e := scalarInt(v)
		if e != nil || seconds < 0 {
			return e
		}
		time.Sleep(time.Duration(seconds) * time.Second)
		f.Result = NullValue
	case "do_syschmod":
		path, e := frameText(c, f, 0)
		if e != nil {
			return e
		}
		modeValue, e := frameValue(c, f, 1)
		if e != nil {
			return e
		}
		mode, e := scalarInt(modeValue)
		if e != nil {
			return e
		}
		if e = os.Chmod(path, os.FileMode(mode)); e != nil {
			return e
		}
		f.Result = &IntegerVector{Data: []int64{0}}
	case "do_unlink":
		path, e := frameText(c, f, 0)
		if e != nil {
			return e
		}
		if e = os.Remove(path); e != nil {
			return e
		}
		f.Result = &IntegerVector{Data: []int64{0}}
	default:
		f.Result = &CharacterVector{Data: []string{strings.TrimPrefix(f.Plan.Name, "Sys.")}}
	}
	return nil
}
