package runtime

import (
	"encoding/base64"
	"fmt"

	"r2go/syntax"
)

// RunEncodedProgram is the compatibility boundary for R constructs that the
// native Go lowerer has not covered yet. Generated programs depend only on the
// public Pure-Go runtime; syntax and evaluator internals stay encapsulated.
func RunEncodedProgram(encoded string) (Value, error) {
	return RunEncodedProgramInContext(NewContext(), encoded)
}

// RunEncodedProgramInContext executes one compiled block while retaining all
// bindings, closures, options and host state created by preceding blocks.
func RunEncodedProgramInContext(ctx *Context, encoded string) (Value, error) {
	if ctx == nil {
		ctx = NewContext()
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("decode compiled R program: %w", err)
	}
	program, err := syntax.DecodeProgramIR(data)
	if err != nil {
		return nil, fmt.Errorf("load compiled R program: %w", err)
	}
	return ctx.EvalProgram(program)
}

// RunSourceInContext is the readable compatibility boundary used by generated
// source. The original unsupported R block remains visible as text instead of
// being hidden in a Base64 IR blob; parsing and execution are still Pure Go.
func RunSourceInContext(ctx *Context, source string) (Value, error) {
	if ctx == nil {
		ctx = NewContext()
	}
	program, err := syntax.Parse(source)
	if err != nil {
		return nil, fmt.Errorf("parse compatibility block: %w", err)
	}
	return ctx.EvalProgram(program)
}
