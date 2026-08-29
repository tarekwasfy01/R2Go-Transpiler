package compiler

import (
	"encoding/csv"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"r2go/runtime"
	"r2go/syntax"
)

// TestTranspilationMatrix validates the actual public pipeline for every
// matrix row: R parsing, mixed native lowering, and syntactically valid Go.
// This deliberately tests more than primitive registration.
func TestTranspilationMatrix(t *testing.T) {
	path := filepath.Join("..", "reference", "test_matrix.csv")
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
		t.Fatal("transpilation matrix has no cases")
	}
	for _, row := range records[1:] {
		if len(row) < 3 {
			t.Fatalf("invalid matrix row: %q", row)
		}
		id, source := row[0], row[2]
		t.Run(id, func(t *testing.T) {
			program, err := syntax.Parse(source)
			if err != nil {
				t.Fatalf("parse R: %v", err)
			}
			generated, err := GenerateMainWithSource(program, source)
			if err != nil {
				t.Fatalf("generate Go: %v", err)
			}
			if _, err := parser.ParseFile(token.NewFileSet(), "generated.go", generated, parser.AllErrors); err != nil {
				t.Fatalf("parse generated Go: %v\n%s", err, generated)
			}
			if strings.Contains(string(generated), "RunEncodedProgramInContext") {
				t.Fatal("opaque IR fallback emitted")
			}
			ctx := runtime.NewContext()
			if _, err := ctx.EvalProgram(program); err != nil {
				t.Fatalf("execute Pure-Go runtime: %v", err)
			}
		})
	}
}
