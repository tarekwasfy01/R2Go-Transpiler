package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gioui.org/app"
	"r2go/compiler"
	"r2go/internal/appcore"
	"r2go/internal/licenses"
	"r2go/internal/platform"
	"r2go/internal/ui"
	"r2go/runtime"
	"r2go/syntax"
)

func main() {
	if len(os.Args) < 2 {
		launchGUI()
		return
	}
	var err error
	switch os.Args[1] {
	case "gui":
		launchGUI()
		return
	case "--licenses", "licenses":
		platform.EnsureCLIConsole()
		fmt.Print(licenses.ThirdPartyNotices())
		return
	case "--help", "-h", "help":
		platform.EnsureCLIConsole()
		usage()
		return
	case "run":
		platform.EnsureCLIConsole()
		err = run(os.Args[2:])
	case "ast":
		platform.EnsureCLIConsole()
		err = ast(os.Args[2:])
	case "transpile":
		platform.EnsureCLIConsole()
		err = transpile(os.Args[2:])
	case "version":
		platform.EnsureCLIConsole()
		fmt.Println("r2go development")
	case "coverage":
		platform.EnsureCLIConsole()
		err = coverage(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "r2go:", err)
		os.Exit(1)
	}
}

func launchGUI() {
	go func() {
		gui := ui.New(appcore.R2GoEngine{})
		if err := gui.Run(); err != nil {
			fmt.Fprintln(os.Stderr, "r2go GUI:", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

func usage() {
	fmt.Fprintln(os.Stderr, `r2go - Pure-Go R to Go transpiler

usage: r2go <gui|run|ast|transpile|coverage|version> [options]

  r2go                         start the GUI
  r2go gui                     start the GUI
  r2go run input.R             execute R source with the Pure-Go runtime
  r2go ast input.R             print the parsed syntax tree as JSON
  r2go transpile input.R -o out.go
				generate a Go main package
  r2go coverage                print primitive coverage metadata
  r2go version                 print the version
  r2go --licenses              print embedded third-party notices`)
}
func readOne(args []string) (string, []byte, error) {
	if len(args) != 1 {
		return "", nil, fmt.Errorf("expected one R source file")
	}
	b, e := os.ReadFile(args[0])
	return args[0], b, e
}

func coverage(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("coverage takes no arguments; its table is generated from pinned GNU R sources")
	}
	seen := map[string]bool{}
	missing := []string{}
	implemented := 0
	implementedFamilies := map[string]bool{}
	for _, status := range runtime.PrimitiveCoverageMatrix() {
		name := status.Name
		if seen[name] {
			continue
		}
		seen[name] = true
		if status.Implemented {
			implemented++
			implementedFamilies[status.CEntry] = true
		} else {
			missing = append(missing, name)
		}
	}
	report := struct {
		Inventory       int      `json:"inventory"`
		Recognized      int      `json:"recognized"`
		Missing         int      `json:"missing_count"`
		Coverage        float64  `json:"coverage_percent"`
		Families        int      `json:"dispatcher_families"`
		TouchedFamilies int      `json:"touched_dispatcher_families"`
		MissingNames    []string `json:"missing_names"`
	}{Inventory: len(seen), Recognized: implemented, Missing: len(missing), Families: len(runtime.PrimitiveByEntry), TouchedFamilies: len(implementedFamilies), MissingNames: missing}
	if report.Inventory > 0 {
		report.Coverage = 100 * float64(implemented) / float64(report.Inventory)
	}
	out, _ := json.MarshalIndent(report, "", "  ")
	fmt.Println(string(out))
	return nil
}
func run(args []string) error {
	_, b, e := readOne(args)
	if e != nil {
		return e
	}
	p, e := syntax.Parse(string(b))
	if e != nil {
		return e
	}
	v, e := runtime.NewContext().EvalProgram(p)
	if e != nil {
		return e
	}
	if v != nil && v.Kind() != runtime.NullKind {
		fmt.Println(v.String())
	}
	return nil
}
func ast(args []string) error {
	_, b, e := readOne(args)
	if e != nil {
		return e
	}
	p, e := syntax.Parse(string(b))
	if e != nil {
		return e
	}
	out, e := json.MarshalIndent(p, "", "  ")
	if e == nil {
		fmt.Println(string(out))
	}
	return e
}
func transpile(args []string) error {
	fs := flag.NewFlagSet("transpile", flag.ContinueOnError)
	out := fs.String("o", "", "output Go file")
	strictNative := fs.Bool("strict-native", false, "fail instead of emitting the compatibility IR fallback")
	sourceComments := fs.Bool("source-comments", true, "retain original R as comments when compatibility fallback is used")
	if e := fs.Parse(args); e != nil {
		return e
	}
	name, b, e := readOne(fs.Args())
	if e != nil {
		return e
	}
	p, e := syntax.Parse(string(b))
	if e != nil {
		return e
	}
	if *out == "" {
		*out = strings.TrimSuffix(name, filepath.Ext(name)) + ".go"
	}
	var source []byte
	if *strictNative {
		source, e = compiler.GenerateNativeMain(p)
	} else {
		source, e = compiler.GenerateMainWithOptions(p, string(b), compiler.GenerateOptions{
			AllowIRFallback:  true,
			PreserveOriginal: *sourceComments,
		})
	}
	if e != nil {
		return e
	}
	return os.WriteFile(*out, source, 0644)
}
