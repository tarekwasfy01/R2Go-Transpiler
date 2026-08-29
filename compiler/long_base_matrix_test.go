package compiler

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"r2go/syntax"
)

func TestLongBaseScriptIsFullyGoLowered(t *testing.T) {
	path := filepath.Join("..", "examples", "long_base_matrix.R")
	sourceBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	source := string(sourceBytes)
	program, err := syntax.Parse(source)
	if err != nil {
		t.Fatalf("parse long Base-R fixture: %v", err)
	}
	if len(program.Expressions) != 100 {
		t.Fatalf("fixture block count=%d, want 100", len(program.Expressions))
	}
	generated, err := GenerateMainWithOptions(program, source, GenerateOptions{AllowIRFallback: true})
	if err != nil {
		t.Fatalf("generate long Base-R fixture: %v", err)
	}
	text := string(generated)
	for _, forbidden := range []string{
		"compatibilitySource", "compatibilityBlock", "RunSourceInContext",
		"RunEncodedProgramInContext", "matrix lowering unavailable", "encoding/base64",
	} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("long fixture contains non-Go block %q", forbidden)
		}
	}
	if got := strings.Count(text, "matrix-lowered block"); got != 100 {
		t.Fatalf("matrix-lowered blocks=%d, want 100", got)
	}
	if _, err := parser.ParseFile(token.NewFileSet(), "generated.go", generated, parser.AllErrors); err != nil {
		t.Fatalf("generated Go is invalid: %v", err)
	}
}
