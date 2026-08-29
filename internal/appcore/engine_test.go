package appcore

import (
	"context"
	"strings"
	"testing"
)

func TestR2GoEngineDefaultMixedOutputKeepsEveryBlockExecutable(t *testing.T) {
	source := "x <- 1\nmystery_call(x)\ny <- 2\nprint(y)"
	output, err := (R2GoEngine{}).Transpile(context.Background(), source)
	if err != nil {
		t.Fatalf("Transpile returned error for partially supported source: %v\n%s", err, output)
	}
	if strings.TrimSpace(output) == "" {
		t.Fatal("Transpile returned empty output")
	}
	for _, want := range []string{`SetGlobal(ctx, "x"`, `MustCall(ctx, "mystery_call"`, `SetGlobal(ctx, "y"`, `MustCall(ctx, "print"`, "// R| x <- 1", "// R| mystery_call(x)", "// R| y <- 2"} {
		if !strings.Contains(output, want) {
			t.Fatalf("GUI/default output lacks %q:\n%s", want, output)
		}
	}
	if strings.Contains(output, "RunEncodedProgramInContext") {
		t.Fatalf("GUI/default output used compatibility IR for matrix-lowerable blocks:\n%s", output)
	}
}

func TestR2GoEngineParseFailureReturnsVisibleOutputAndError(t *testing.T) {
	output, err := (R2GoEngine{}).Transpile(context.Background(), "x <- (")
	if err == nil {
		t.Fatal("invalid R unexpectedly transpiled without error")
	}
	if strings.TrimSpace(output) == "" {
		t.Fatalf("parse failure returned empty output; err=%v", err)
	}
	lower := strings.ToLower(output)
	if !strings.Contains(lower, "error") && !strings.Contains(lower, "parse") {
		t.Fatalf("failure output does not explain the conversion error: %q", output)
	}
}

func TestR2GoEngineStrictGenerationFailureReturnsVisibleOutputAndError(t *testing.T) {
	output, err := (R2GoEngine{}).TranspileWithOptions(context.Background(), "mystery_call(1)", TranspileOptions{
		AllowIRFallback:  false,
		PreserveOriginal: true,
	})
	if err == nil {
		t.Fatal("strict unsupported source unexpectedly transpiled without error")
	}
	if strings.TrimSpace(output) == "" {
		t.Fatalf("generation failure returned empty output; err=%v", err)
	}
}
