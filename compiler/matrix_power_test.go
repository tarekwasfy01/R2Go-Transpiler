package compiler

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"r2go/runtime"
	"r2go/syntax"
)

type compositionFactor struct {
	stage, name, inputClass, outputClass, template string
}

type poweredComposition struct {
	name, source, outputClass string
}

func loadCompositionFactors(t *testing.T) []compositionFactor {
	t.Helper()
	path := filepath.Join("..", "reference", "composition_test_matrix.csv")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(records) < 2 {
		t.Fatal("composition matrix has no factors")
	}
	factors := make([]compositionFactor, 0, len(records)-1)
	for rowIndex, row := range records[1:] {
		if len(row) != 5 {
			t.Fatalf("composition matrix row %d has %d columns, want 5", rowIndex+2, len(row))
		}
		factor := compositionFactor{stage: row[0], name: row[1], inputClass: row[2], outputClass: row[3], template: row[4]}
		switch factor.stage {
		case "producer":
			if factor.inputClass != "" || strings.Contains(factor.template, "%s") {
				t.Fatalf("invalid producer row %q", row)
			}
		case "transform", "observer":
			if factor.inputClass == "" || strings.Count(factor.template, "%s") != 1 {
				t.Fatalf("invalid %s row %q", factor.stage, row)
			}
		default:
			t.Fatalf("unknown composition stage %q", factor.stage)
		}
		factors = append(factors, factor)
	}
	return factors
}

// powerCompositionMatrix performs boolean incidence multiplication over the
// type classes in composition_test_matrix.csv. power=2 is P*T^2*O: every
// producer is passed through every pair of type-compatible transformations
// and then every compatible observer. Invalid type paths never enter the
// product, so failures identify semantic composition defects rather than
// deliberately invalid R programs.
func powerCompositionMatrix(t *testing.T, power int) []poweredComposition {
	t.Helper()
	factors := loadCompositionFactors(t)
	var producers, transforms, observers []compositionFactor
	for _, factor := range factors {
		switch factor.stage {
		case "producer":
			producers = append(producers, factor)
		case "transform":
			transforms = append(transforms, factor)
		case "observer":
			observers = append(observers, factor)
		}
	}
	type path struct {
		name, source, class string
	}
	paths := make([]path, len(producers))
	for i, producer := range producers {
		paths[i] = path{name: producer.name, source: producer.template, class: producer.outputClass}
	}
	for exponent := 0; exponent < power; exponent++ {
		next := make([]path, 0, len(paths)*4)
		for _, current := range paths {
			for _, transform := range transforms {
				if current.class != transform.inputClass {
					continue
				}
				next = append(next, path{
					name:   current.name + "__" + transform.name,
					source: fmt.Sprintf(transform.template, current.source),
					class:  transform.outputClass,
				})
			}
		}
		paths = next
	}
	compositions := make([]poweredComposition, 0, len(paths)*4)
	for _, current := range paths {
		for _, observer := range observers {
			if current.class != observer.inputClass {
				continue
			}
			compositions = append(compositions, poweredComposition{
				name:        current.name + "__" + observer.name,
				source:      fmt.Sprintf(observer.template, current.source),
				outputClass: observer.outputClass,
			})
		}
	}
	sort.Slice(compositions, func(i, j int) bool { return compositions[i].name < compositions[j].name })
	return compositions
}

func evaluatePoweredCompositions(t *testing.T, cases []poweredComposition) ([]string, []string) {
	t.Helper()
	canonical := make([]string, len(cases))
	printed := make([]string, len(cases))
	var failures []string
	for i, test := range cases {
		program, err := syntax.Parse(test.source)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s parse: %v", test.name, err))
			continue
		}
		value, err := runtime.NewContext().EvalProgram(program)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s execute %q: %v", test.name, test.source, err))
			continue
		}
		canonical[i] = canonicalRuntimeValue(value)
		printed[i] = value.String()
	}
	return canonical, append(printed, failures...)
}

func runRCompositionBatch(t *testing.T, rscript string, cases []poweredComposition) []string {
	t.Helper()
	var source strings.Builder
	source.WriteString(`encode_value <- function(x) {
  typ <- typeof(x)
  encode_element <- function(v) {
    if (is.na(v)) return("NA")
    if (typ == "character") return(encodeString(v, quote="\""))
	if (typ == "logical") return(if (v) "TRUE" else "FALSE")
	if (typ == "double") return(sprintf("%.17g", v))
    as.character(v)
  }
  payload <- if (length(x) == 0) "" else paste(vapply(seq_along(x), function(i) encode_element(x[[i]]), character(1)), collapse="|")
  paste0(typ, ":", payload)
}
`)
	for i, test := range cases {
		fmt.Fprintf(&source, "cat(\"%06d\\t\", encode_value((%s)), \"\\n\", sep=\"\")\n", i, test.source)
	}
	path := filepath.Join(t.TempDir(), "matrix-power-reference.R")
	if err := os.WriteFile(path, []byte(source.String()), 0o600); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(rscript, "--vanilla", path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("GNU R composition batch failed: %v: %s", err, stderr.String())
	}
	values := make([]string, len(cases))
	lines := strings.Split(strings.TrimSpace(strings.ReplaceAll(string(out), "\r\n", "\n")), "\n")
	for _, line := range lines {
		fields := strings.SplitN(line, "\t", 2)
		if len(fields) != 2 {
			t.Fatalf("invalid GNU R matrix line %q", line)
		}
		var index int
		if _, err := fmt.Sscanf(fields[0], "%d", &index); err != nil || index < 0 || index >= len(values) {
			t.Fatalf("invalid GNU R matrix index %q", fields[0])
		}
		values[index] = fields[1]
	}
	return values
}

func verifyCompositionPowerAgainstGNUR(t *testing.T, power, minimum int) {
	t.Helper()
	cases := powerCompositionMatrix(t, power)
	if len(cases) < minimum {
		t.Fatalf("matrix power %d produced only %d cases, want at least %d", power, len(cases), minimum)
	}
	got, printedAndFailures := evaluatePoweredCompositions(t, cases)
	if len(printedAndFailures) > len(cases) {
		failures := printedAndFailures[len(cases):]
		limit := len(failures)
		if limit > 40 {
			limit = 40
		}
		t.Fatalf("Pure-Go execution failed in %d/%d matrix-power cases:\n%s", len(failures), len(cases), strings.Join(failures[:limit], "\n"))
	}
	rscript := `C:\Users\tarek\Desktop\GroundGIS\Release\Package\runtimes\R\bin\Rscript.exe`
	if _, err := os.Stat(rscript); err != nil {
		t.Skipf("GNU R reference is unavailable after %d Pure-Go cases: %v", len(cases), err)
	}
	want := runRCompositionBatch(t, rscript, cases)
	var mismatches []string
	for i := range cases {
		if got[i] != want[i] {
			mismatches = append(mismatches, fmt.Sprintf("%s: Go=%q R=%q source=%s", cases[i].name, got[i], want[i], cases[i].source))
		}
	}
	if len(mismatches) != 0 {
		limit := len(mismatches)
		if limit > 40 {
			limit = 40
		}
		t.Fatalf("GNU R differs in %d/%d matrix-power cases:\n%s", len(mismatches), len(cases), strings.Join(mismatches[:limit], "\n"))
	}
	t.Logf("P*T^%d*O differential matrix passed %d compositions", power, len(cases))
}

func TestMatrixPowerCompositionAgainstGNU_R(t *testing.T) {
	verifyCompositionPowerAgainstGNUR(t, 2, 4000)
}

func TestMatrixCubeCompositionAgainstGNU_R(t *testing.T) {
	verifyCompositionPowerAgainstGNUR(t, 3, 35000)
}

func TestMatrixPowerGeneratedGo(t *testing.T) {
	cases := powerCompositionMatrix(t, 2)
	_, printedAndFailures := evaluatePoweredCompositions(t, cases)
	if len(printedAndFailures) > len(cases) {
		t.Fatalf("cannot establish generated-Go oracle: %s", strings.Join(printedAndFailures[len(cases):], "\n"))
	}
	expected := printedAndFailures[:len(cases)]
	repositoryRoot, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	temporary := t.TempDir()
	const shardSize = 256
	shards := 0
	for start := 0; start < len(cases); start += shardSize {
		end := start + shardSize
		if end > len(cases) {
			end = len(cases)
		}
		shardCases, shardExpected := cases[start:end], expected[start:end]
		var rSource strings.Builder
		for _, test := range shardCases {
			fmt.Fprintf(&rSource, "print((%s))\n", test.source)
		}
		program, err := syntax.Parse(rSource.String())
		if err != nil {
			t.Fatal(err)
		}
		generated, err := GenerateMainWithOptions(program, rSource.String(), GenerateOptions{AllowIRFallback: true})
		if err != nil {
			t.Fatal(err)
		}
		text := string(generated)
		for _, forbidden := range []string{"compatibilitySource", "compatibilityBlock", "RunSourceInContext", "RunEncodedProgramInContext", "encoding/base64"} {
			if strings.Contains(text, forbidden) {
				t.Fatalf("P*T^2*O shard %d contains fallback %q", shards, forbidden)
			}
		}
		if got := strings.Count(text, "matrix-lowered block"); got != len(shardCases) {
			t.Fatalf("generated shard %d matrix blocks=%d, want %d", shards, got, len(shardCases))
		}
		if _, err := parser.ParseFile(token.NewFileSet(), "matrix-power.go", generated, parser.AllErrors); err != nil {
			t.Fatalf("generated matrix-power shard %d is invalid Go: %v", shards, err)
		}
		generatedPath := filepath.Join(temporary, fmt.Sprintf("matrix-power-%02d.go", shards))
		if err := os.WriteFile(generatedPath, generated, 0o600); err != nil {
			t.Fatal(err)
		}
		cmd := exec.Command("go", "run", generatedPath)
		cmd.Dir = repositoryRoot
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		out, err := cmd.Output()
		if err != nil {
			t.Fatalf("generated P*T^2*O shard %d failed: %v: %s", shards, err, stderr.String())
		}
		actual := strings.Split(strings.TrimSuffix(strings.ReplaceAll(string(out), "\r\n", "\n"), "\n"), "\n")
		if len(actual) != len(shardExpected) {
			t.Fatalf("generated Go shard %d output lines=%d, want %d", shards, len(actual), len(shardExpected))
		}
		var mismatches []string
		for i := range shardExpected {
			if actual[i] != shardExpected[i] {
				mismatches = append(mismatches, fmt.Sprintf("%s: generated=%q runtime=%q", shardCases[i].name, actual[i], shardExpected[i]))
			}
		}
		if len(mismatches) != 0 {
			limit := len(mismatches)
			if limit > 40 {
				limit = 40
			}
			t.Fatalf("generated Go shard %d differs in %d/%d cases:\n%s", shards, len(mismatches), len(shardCases), strings.Join(mismatches[:limit], "\n"))
		}
		shards++
	}
	t.Logf("generated Go compiled and executed all %d P*T^2*O compositions in %d shards", len(cases), shards)
}
