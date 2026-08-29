package compiler

import (
	"strings"
	"testing"

	"r2go/syntax"
)

func TestNativeRegressionCasesRemainReadable(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   []string
	}{
		{
			name:   "scalar",
			source: "x <- 2\ny <- x^3 + 1\ny",
			want:   []string{"r2go: native-lowered", "math.Pow"},
		},
		{
			name:   "vector",
			source: "x <- c(1, 2, 3)\ny <- x * 2\nprint(y)",
			want:   []string{"[]float64", "rVectorBinary"},
		},
		{
			name:   "function",
			source: "twice <- function(x) x * 2\nprint(twice(21))",
			want:   []string{"func r_twice", "r_twice(21)"},
		},
		{
			name:   "named matrix primitive",
			source: "print(matrix(c(1, 2, 3, 4), nrow=2, ncol=2))",
			want:   []string{"MustCallPrimitive", "NamedPrimitiveArg(\"nrow\"", "NamedPrimitiveArg(\"ncol\""},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			program, err := syntax.Parse(tc.source)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			generated, err := GenerateNativeMain(program)
			if err != nil {
				t.Fatalf("GenerateNativeMain: %v", err)
			}
			text := string(generated)
			for _, want := range tc.want {
				if !strings.Contains(text, want) {
					t.Fatalf("generated output lacks %q:\n%s", want, text)
				}
			}
			for _, forbidden := range []string{"compiledBlock", "RunEncodedProgramInContext", "programIR", "encoding/base64"} {
				if strings.Contains(text, forbidden) {
					t.Fatalf("native output contains opaque fallback %q:\n%s", forbidden, text)
				}
			}
		})
	}
}

func TestDefaultGenerationKeepsSupportedBlocksAroundUnsupportedBlock(t *testing.T) {
	source := "x <- 1\nmystery_call(x)\ny <- 2\nprint(y)"
	program, err := syntax.Parse(source)
	if err != nil {
		t.Fatal(err)
	}
	generated, err := GenerateMainWithSource(program, source)
	if err != nil {
		t.Fatalf("GenerateMainWithSource: %v", err)
	}
	text := string(generated)
	if strings.TrimSpace(text) == "" {
		t.Fatal("mixed program generated empty output")
	}
	for _, want := range []string{`SetGlobal(ctx, "x"`, `MustCall(ctx, "mystery_call"`, `SetGlobal(ctx, "y"`, `MustCall(ctx, "print"`, "// R| mystery_call(x)"} {
		if !strings.Contains(text, want) {
			t.Fatalf("mixed output lacks %q:\n%s", want, text)
		}
	}
	if strings.Contains(text, "RunEncodedProgramInContext") {
		t.Fatalf("fully matrix-lowerable blocks unexpectedly use compatibility IR:\n%s", text)
	}
}

func TestUnsupportedCRLFBlockCommentStaysAligned(t *testing.T) {
	source := "x <- 1\r\nmystery_call(\r\n  x\r\n)\r\ny <- 2\r\n"
	program, err := syntax.Parse(source)
	if err != nil {
		t.Fatal(err)
	}
	generated, err := GenerateMainWithSource(program, source)
	if err != nil {
		t.Fatalf("GenerateMainWithSource: %v", err)
	}
	text := string(generated)
	for _, want := range []string{"// R| mystery_call(", "// R|   x", "// R| )"} {
		if !strings.Contains(text, want) {
			t.Fatalf("CRLF source comment is misaligned; missing %q:\n%s", want, text)
		}
	}
	if strings.Contains(text, "\r") {
		t.Fatalf("generated Go retained CR characters:\n%q", text)
	}
}
