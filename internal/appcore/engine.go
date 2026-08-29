package appcore

import (
	"context"
	"fmt"
	"strings"

	"r2go/compiler"
	"r2go/syntax"
)

// Engine is the single conversion boundary shared by the GUI and command line.
type Engine interface {
	Transpile(ctx context.Context, rCode string) (string, error)
}

type TranspileOptions struct {
	AllowIRFallback  bool
	PreserveOriginal bool
}

var DefaultTranspileOptions = TranspileOptions{AllowIRFallback: true, PreserveOriginal: true}

type ConfigurableEngine interface {
	TranspileWithOptions(ctx context.Context, rCode string, options TranspileOptions) (string, error)
}

// R2GoEngine lowers R source through the production parser and Go generator.
// Context is checked on either side of parsing, so a superseded GUI request
// never publishes an obsolete result.
type R2GoEngine struct{}

func (R2GoEngine) Transpile(ctx context.Context, rCode string) (string, error) {
	return R2GoEngine{}.TranspileWithOptions(ctx, rCode, DefaultTranspileOptions)
}

func (R2GoEngine) TranspileWithOptions(ctx context.Context, rCode string, options TranspileOptions) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}
	program, err := syntax.Parse(rCode)
	if err != nil {
		wrapped := fmt.Errorf("parse R source: %w", err)
		return conversionErrorOutput("parse", wrapped, rCode, options.PreserveOriginal), wrapped
	}
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}
	result, err := compiler.GenerateMainWithOptions(program, rCode, compiler.GenerateOptions{
		AllowIRFallback:  options.AllowIRFallback,
		PreserveOriginal: options.PreserveOriginal,
	})
	if err != nil {
		wrapped := fmt.Errorf("generate Go: %w", err)
		return conversionErrorOutput("generation", wrapped, rCode, options.PreserveOriginal), wrapped
	}
	return string(result), nil
}

func conversionErrorOutput(stage string, err error, rCode string, preserveOriginal bool) string {
	var b strings.Builder
	b.WriteString("// r2go conversion error during ")
	b.WriteString(stage)
	b.WriteString(":\n")
	for _, line := range strings.Split(strings.ReplaceAll(err.Error(), "\r\n", "\n"), "\n") {
		b.WriteString("// ERROR| ")
		b.WriteString(line)
		b.WriteByte('\n')
	}
	if preserveOriginal && strings.TrimSpace(rCode) != "" {
		b.WriteString("// Original R source:\n")
		for _, line := range strings.Split(strings.ReplaceAll(rCode, "\r\n", "\n"), "\n") {
			b.WriteString("// R| ")
			b.WriteString(line)
			b.WriteByte('\n')
		}
	}
	return b.String()
}
