package compiler

import (
	"r2go/syntax"
	"strings"
	"testing"
)

func TestGenerateMainExecutesEveryMixedProgramBlock(t *testing.T) {
	source := "x <- 1\nmystery_call(x)\ny <- 2"
	p, err := syntax.Parse(source)
	if err != nil {
		t.Fatal(err)
	}
	generated, err := GenerateMainWithSource(p, source)
	if err != nil {
		t.Fatal(err)
	}
	text := string(generated)
	for _, forbidden := range []string{"syntax.Parse", "programIR", "r2go/syntax", "DecodeProgramIR", "syntax.Program", "encoding/base64"} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("generated source contains opaque compatibility artifact %q:\n%s", forbidden, text)
		}
	}
	if !strings.Contains(text, "// R| mystery_call(x)") {
		t.Fatalf("unsupported R block is not preserved as a readable comment:\n%s", text)
	}
	for _, want := range []string{`SetGlobal(ctx, "x"`, `MustCall(ctx, "mystery_call"`, `SetGlobal(ctx, "y"`} {
		if !strings.Contains(text, want) {
			t.Fatalf("matrix output lacks %q:\n%s", want, text)
		}
	}
}

func TestGenerateMainNativeScalarLowering(t *testing.T) {
	p, e := syntax.Parse("x <- 2\ny <- x^3 + 1\ny")
	if e != nil {
		t.Fatal(e)
	}
	generated, e := GenerateMain(p)
	if e != nil {
		t.Fatal(e)
	}
	text := string(generated)
	for _, forbidden := range []string{"r2go/runtime", "r2go/syntax", "EvalProgram", "syntax.Program"} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("native output contains %s", forbidden)
		}
	}
	if !strings.Contains(text, "r2go: native-lowered") || !strings.Contains(text, "math.Pow") {
		t.Fatal("native lowering missing")
	}
}

func TestGenerateMainNativeVectorLowering(t *testing.T) {
	p, err := syntax.Parse("x <- c(1, 2, 3)\ny <- x * 2\nprint(y)")
	if err != nil {
		t.Fatal(err)
	}
	generated, err := GenerateMain(p)
	if err != nil {
		t.Fatal(err)
	}
	text := string(generated)
	for _, want := range []string{"[]float64", "rVectorBinary", `"*"`} {
		if !strings.Contains(text, want) {
			t.Fatalf("native vector output lacks %s:\n%s", want, text)
		}
	}
	for _, forbidden := range []string{"programIR", "DecodeProgramIR", "EvalProgram", "r2go/runtime"} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("native vector output contains interpreter fallback %s", forbidden)
		}
	}
}

func TestUnsupportedBlocksStayAsOriginalRComments(t *testing.T) {
	source := "f <- function(x = 2) x + 1\nf()"
	p, err := syntax.Parse(source)
	if err != nil {
		t.Fatal(err)
	}
	generated, err := GenerateMainWithSource(p, source)
	if err != nil {
		t.Fatal(err)
	}
	text := string(generated)
	for _, want := range []string{"// R| f <- function(x = 2) x + 1", "// R| f()"} {
		if !strings.Contains(text, want) {
			t.Fatalf("fallback output lacks %q:\n%s", want, text)
		}
	}
	for _, required := range []string{"compatibilitySource0001", "RunSourceInContext", `MustCall(ctx, "f"`} {
		if !strings.Contains(text, required) {
			t.Fatalf("unsupported block is not executable via %q:\n%s", required, text)
		}
	}
}

func TestFallbackCommentsRemainAlignedWithWindowsCRLF(t *testing.T) {
	source := "f <- function(x = 2) x + 1\r\nf()\r\ny <- 3\r\n"
	p, err := syntax.Parse(source)
	if err != nil {
		t.Fatal(err)
	}
	generated, err := GenerateMainWithSource(p, source)
	if err != nil {
		t.Fatal(err)
	}
	text := string(generated)
	for _, want := range []string{"// R| f <- function(x = 2) x + 1", "// R| f()", "// R| y <- 3"} {
		if !strings.Contains(text, want) {
			t.Fatalf("CRLF fallback comment/native block is misaligned; missing %q:\n%s", want, text)
		}
	}
}

func TestStrictNativeModeRejectsCompatibilityFallback(t *testing.T) {
	p, err := syntax.Parse("f <- function(x = 2) x + 1\nf()")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = GenerateNativeMain(p); err != ErrNativeLoweringIncomplete {
		t.Fatalf("strict native error=%v", err)
	}

	p, err = syntax.Parse("x <- c(1,2,3)\nprint(x*2)")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = GenerateNativeMain(p); err != nil {
		t.Fatalf("strict native rejected supported vector program: %v", err)
	}
}

func TestFallbackOptions(t *testing.T) {
	source := "f <- function(x = 2) x + 1\nf()"
	p, err := syntax.Parse(source)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = GenerateMainWithOptions(p, source, GenerateOptions{}); err != ErrNativeLoweringIncomplete {
		t.Fatalf("disabled fallback error=%v", err)
	}
	generated, err := GenerateMainWithOptions(p, source, GenerateOptions{AllowIRFallback: true})
	if err != nil {
		t.Fatal(err)
	}
	text := string(generated)
	if strings.Contains(text, "// R|") {
		t.Fatal("original R comments emitted while PreserveOriginal is false")
	}
	if !strings.Contains(text, "matrix lowering unavailable for block") {
		t.Fatalf("enabled partial fallback did not mark unsupported blocks:\n%s", text)
	}
	for _, required := range []string{"compatibilitySource", "RunSourceInContext"} {
		if !strings.Contains(text, required) {
			t.Fatalf("enabled fallback lacks executable artifact %q:\n%s", required, text)
		}
	}
	for _, forbidden := range []string{"compatibilityBlock", "RunEncodedProgramInContext"} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("source-backed fallback unexpectedly contains IR artifact %q:\n%s", forbidden, text)
		}
	}
}

func TestSimpleFunctionIsDirectGoWithoutRuntime(t *testing.T) {
	p, err := syntax.Parse("twice <- function(x) x * 2\nprint(twice(21))")
	if err != nil {
		t.Fatal(err)
	}
	generated, err := GenerateNativeMain(p)
	if err != nil {
		t.Fatal(err)
	}
	text := string(generated)
	for _, want := range []string{"func r_twice(r_x float64) float64", "r_twice(21)"} {
		if !strings.Contains(text, want) {
			t.Fatalf("direct function output lacks %q:\n%s", want, text)
		}
	}
	for _, forbidden := range []string{"r2go/runtime", "programIR", "RunEncodedProgram"} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("simple function unnecessarily contains %q", forbidden)
		}
	}
}

func TestKnownMatrixFunctionUsesTargetedDispatchNotIR(t *testing.T) {
	p, err := syntax.Parse("print(sum(c(1, 2, 3)))")
	if err != nil {
		t.Fatal(err)
	}
	generated, err := GenerateNativeMain(p)
	if err != nil {
		t.Fatal(err)
	}
	text := string(generated)
	for _, want := range []string{`"r2go/runtime"`, `MustCallPrimitive(rContext, "sum"`} {
		if !strings.Contains(text, want) {
			t.Fatalf("matrix forwarding output lacks %q:\n%s", want, text)
		}
	}
	for _, forbidden := range []string{"programIR", "RunEncodedProgram"} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("matrix call fell back to whole-program IR: %s", forbidden)
		}
	}
}
