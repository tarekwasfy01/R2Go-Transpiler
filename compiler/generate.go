package compiler

import (
	"encoding/base64"
	"fmt"
	"go/format"
	"strconv"
	"strings"

	"r2go/syntax"
)

var ErrNativeLoweringIncomplete = fmt.Errorf("program contains R structures not yet supported by the native Go lowerer")

// GenerateOptions controls the honest boundary between native Go lowering and
// the Pure-Go compatibility runtime.
type GenerateOptions struct {
	AllowIRFallback  bool
	PreserveOriginal bool
}

var DefaultGenerateOptions = GenerateOptions{AllowIRFallback: true, PreserveOriginal: true}

// GenerateNativeMain requires genuine AST-to-Go lowering and never emits the
// compatibility IR path.
func GenerateNativeMain(program *syntax.Program) ([]byte, error) {
	if native, ok := generateNativeMain(program); ok {
		return native, nil
	}
	return nil, ErrNativeLoweringIncomplete
}

// GenerateMain lowers an already parsed R program into a standalone Go main
// package.
func GenerateMain(program *syntax.Program) ([]byte, error) {
	return GenerateMainWithSource(program, "")
}

// GenerateMainWithSource preserves unsupported original R text as comments
// while keeping every top-level block that can be lowered as readable Go.
func GenerateMainWithSource(program *syntax.Program, originalR string) ([]byte, error) {
	return GenerateMainWithOptions(program, originalR, DefaultGenerateOptions)
}

// GenerateMainWithOptions exposes the fallback policy to CLI and GUI callers.
// The normal fallback is partial native lowering: supported top-level blocks
// stay readable Go and unsupported blocks are retained only as R comments.
func GenerateMainWithOptions(program *syntax.Program, originalR string, options GenerateOptions) ([]byte, error) {
	if native, err := GenerateNativeMain(program); err == nil {
		return native, nil
	}
	if !options.AllowIRFallback {
		return nil, ErrNativeLoweringIncomplete
	}
	partial, ok := generatePartialNativeMain(program, originalR, options.PreserveOriginal)
	if !ok {
		return nil, fmt.Errorf("format partially lowered Go: %w", ErrNativeLoweringIncomplete)
	}
	return partial, nil
}

func writeOriginalBlockComment(b *strings.Builder, originalR string, expression syntax.Expr) {
	source := originalBlockSource(originalR, expression)
	if source == "" {
		return
	}
	for _, line := range strings.Split(source, "\n") {
		b.WriteString("\t// R| ")
		b.WriteString(line)
		b.WriteByte('\n')
	}
}

func originalBlockSource(originalR string, expression syntax.Expr) string {
	if originalR == "" || expression == nil {
		return ""
	}
	span := expression.SourceSpan()
	start, end := span.Start.Offset, span.End.Offset
	if start < 0 || end < start || start > len(originalR) {
		return ""
	}
	if end > len(originalR) {
		end = len(originalR)
	}
	return strings.TrimSpace(strings.ReplaceAll(originalR[start:end], "\r\n", "\n"))
}

// GeneratePackage emits a reusable Go package loader for precompiled R code.
// It is used to turn GNU-R's pure-R base library into a build-time artifact.
func GeneratePackage(program *syntax.Program, packageName, loaderName string) ([]byte, error) {
	if packageName == "" || loaderName == "" {
		return nil, fmt.Errorf("package and loader names are required")
	}
	ir, e := syntax.EncodeProgramIR(program)
	if e != nil {
		return nil, fmt.Errorf("encode IR: %w", e)
	}
	encoded := base64.StdEncoding.EncodeToString(ir)
	var b strings.Builder
	fmt.Fprintf(&b, "package %s\n\nimport (\n\t\"encoding/base64\"\n\t\"r2go/runtime\"\n\t\"r2go/syntax\"\n)\n\nconst programIR = %s\n", packageName, strconv.Quote(encoded))
	fmt.Fprintf(&b, "\nfunc %s(ctx *runtime.Context) error {\n\tdata, err := base64.StdEncoding.DecodeString(programIR)\n\tif err != nil { return err }\n\tprogram, err := syntax.DecodeProgramIR(data)\n\tif err != nil { return err }\n\t_, err = ctx.EvalProgram(program)\n\treturn err\n}\n", loaderName)
	out, e := format.Source([]byte(b.String()))
	if e != nil {
		return nil, fmt.Errorf("format generated package: %w", e)
	}
	return out, nil
}
func emitProgram(b *strings.Builder, p *syntax.Program) {
	b.WriteString("&syntax.Program{Expressions: []syntax.Expr{")
	for _, e := range p.Expressions {
		emitExpr(b, e)
		b.WriteByte(',')
	}
	b.WriteString("}}")
}
func emitExpr(b *strings.Builder, e syntax.Expr) {
	switch x := e.(type) {
	case *syntax.Literal:
		b.WriteString("&syntax.Literal{Kind:syntax.")
		switch x.Kind {
		case syntax.NumberLiteral:
			b.WriteString("NumberLiteral")
		case syntax.StringLiteral:
			b.WriteString("StringLiteral")
		case syntax.LogicalLiteral:
			b.WriteString("LogicalLiteral")
		case syntax.NullLiteral:
			b.WriteString("NullLiteral")
		case syntax.NALiteral:
			b.WriteString("NALiteral")
		}
		b.WriteString(",Text:")
		b.WriteString(strconv.Quote(x.Text))
		b.WriteString(",At:")
		emitSpan(b, x.At)
		b.WriteByte('}')
	case *syntax.Symbol:
		b.WriteString("&syntax.Symbol{Name:")
		b.WriteString(strconv.Quote(x.Name))
		b.WriteString(",At:")
		emitSpan(b, x.At)
		b.WriteByte('}')
	case *syntax.Call:
		b.WriteString("&syntax.Call{Function:")
		emitExpr(b, x.Function)
		b.WriteString(",Arguments:[]syntax.Argument{")
		for _, a := range x.Arguments {
			emitArgument(b, a)
			b.WriteByte(',')
		}
		b.WriteString("},At:")
		emitSpan(b, x.At)
		b.WriteByte('}')
	case *syntax.Block:
		b.WriteString("&syntax.Block{Expressions:[]syntax.Expr{")
		for _, item := range x.Expressions {
			emitExpr(b, item)
			b.WriteByte(',')
		}
		b.WriteString("},At:")
		emitSpan(b, x.At)
		b.WriteByte('}')
	case *syntax.Function:
		b.WriteString("&syntax.Function{Parameters:[]syntax.Parameter{")
		for _, p := range x.Parameters {
			emitParameter(b, p)
			b.WriteByte(',')
		}
		b.WriteString("},Body:")
		emitExpr(b, x.Body)
		b.WriteString(",At:")
		emitSpan(b, x.At)
		b.WriteByte('}')
	case *syntax.If:
		b.WriteString("&syntax.If{Condition:")
		emitExpr(b, x.Condition)
		b.WriteString(",Then:")
		emitExpr(b, x.Then)
		if x.Else != nil {
			b.WriteString(",Else:")
			emitExpr(b, x.Else)
		}
		b.WriteString(",At:")
		emitSpan(b, x.At)
		b.WriteByte('}')
	case *syntax.While:
		b.WriteString("&syntax.While{Condition:")
		emitExpr(b, x.Condition)
		b.WriteString(",Body:")
		emitExpr(b, x.Body)
		b.WriteString(",At:")
		emitSpan(b, x.At)
		b.WriteByte('}')
	case *syntax.For:
		b.WriteString("&syntax.For{Variable:")
		b.WriteString(strconv.Quote(x.Variable))
		b.WriteString(",Sequence:")
		emitExpr(b, x.Sequence)
		b.WriteString(",Body:")
		emitExpr(b, x.Body)
		b.WriteString(",At:")
		emitSpan(b, x.At)
		b.WriteByte('}')
	case *syntax.Repeat:
		b.WriteString("&syntax.Repeat{Body:")
		emitExpr(b, x.Body)
		b.WriteString(",At:")
		emitSpan(b, x.At)
		b.WriteByte('}')
	default:
		panic(fmt.Sprintf("unsupported AST node %T", e))
	}
}
func emitArgument(b *strings.Builder, a syntax.Argument) {
	b.WriteString("{Name:")
	b.WriteString(strconv.Quote(a.Name))
	if a.Value != nil {
		b.WriteString(",Value:")
		emitExpr(b, a.Value)
	}
	b.WriteString(",At:")
	emitSpan(b, a.At)
	b.WriteByte('}')
}
func emitParameter(b *strings.Builder, p syntax.Parameter) {
	b.WriteString("{Name:")
	b.WriteString(strconv.Quote(p.Name))
	if p.Default != nil {
		b.WriteString(",Default:")
		emitExpr(b, p.Default)
	}
	b.WriteString(",At:")
	emitSpan(b, p.At)
	b.WriteByte('}')
}
func emitSpan(b *strings.Builder, s syntax.Span) {
	b.WriteString("syntax.Span{Start:")
	emitPosition(b, s.Start)
	b.WriteString(",End:")
	emitPosition(b, s.End)
	b.WriteByte('}')
}
func emitPosition(b *strings.Builder, p syntax.Position) {
	fmt.Fprintf(b, "syntax.Position{Offset:%d,Line:%d,Column:%d}", p.Offset, p.Line, p.Column)
}
