package compiler

import (
	"encoding/base64"
	"fmt"
	"go/format"
	"strconv"
	"strings"

	"r2go/syntax"
)

func cloneNativeEmitter(e *nativeEmitter) *nativeEmitter {
	copyEmitter := &nativeEmitter{
		variables:    make(map[string]nativeKind, len(e.variables)),
		functionSet:  make(map[string]nativeFunction, len(e.functionSet)),
		declared:     make(map[string]bool, len(e.declared)),
		indent:       e.indent,
		controlDepth: e.controlDepth,
		allowDeclare: e.allowDeclare,
		temporary:    e.temporary,
		usesVector:   e.usesVector,
		usesRuntime:  e.usesRuntime,
	}
	copyEmitter.body.WriteString(e.body.String())
	copyEmitter.functions.WriteString(e.functions.String())
	for name, kind := range e.variables {
		copyEmitter.variables[name] = kind
	}
	for name, function := range e.functionSet {
		copyEmitter.functionSet[name] = nativeFunction{
			parameters: append([]nativeKind(nil), function.parameters...),
			result:     function.result,
		}
	}
	for name, declared := range e.declared {
		copyEmitter.declared[name] = declared
	}
	return copyEmitter
}

func generatePartialNativeMain(program *syntax.Program, originalR string, preserveComments bool) ([]byte, bool) {
	type loweredBlock struct {
		matrixCode string
		encoded    string
		source     string
	}
	blocks := make([]loweredBlock, len(program.Expressions))
	fallbacks := 0
	for i, expression := range program.Expressions {
		if code, ok := (matrixNativeEmitter{}).statement(expression); ok {
			blocks[i].matrixCode = code
			continue
		}
		if readable := originalBlockSource(originalR, expression); readable != "" {
			blocks[i].source = readable
			fallbacks++
			continue
		}
		ir, err := syntax.EncodeProgramIR(&syntax.Program{Expressions: []syntax.Expr{expression}})
		if err != nil {
			return nil, false
		}
		blocks[i].encoded = base64.StdEncoding.EncodeToString(ir)
		fallbacks++
	}

	var source strings.Builder
	source.WriteString("package main\n\n// r2go: readable matrix-lowered Go with executable compatibility only where required.\nimport (\n")
	if fallbacks > 0 {
		source.WriteString("\t\"fmt\"\n")
	}
	source.WriteString("\t\"r2go/runtime\"\n)\n\n")
	for i, block := range blocks {
		if block.source != "" {
			fmt.Fprintf(&source, "const compatibilitySource%04d = %s\n", i+1, readableGoString(block.source))
		} else if block.encoded != "" {
			fmt.Fprintf(&source, "const compatibilityBlock%04d = %s\n", i+1, strconv.Quote(block.encoded))
		}
	}
	source.WriteString("\nfunc main() {\n\tctx := runtime.NewContext()\n")
	for i, expression := range program.Expressions {
		block := blocks[i]
		commentSuffix := ""
		if preserveComments {
			commentSuffix = "; original R follows"
		}
		if block.matrixCode == "" {
			fmt.Fprintf(&source, "\n\t// r2go: matrix lowering unavailable for block %d/%d%s.\n", i+1, len(program.Expressions), commentSuffix)
		} else {
			fmt.Fprintf(&source, "\n\t// r2go: matrix-lowered block %d/%d%s.\n", i+1, len(program.Expressions), commentSuffix)
		}
		if preserveComments {
			writeOriginalBlockComment(&source, originalR, expression)
		}
		if block.matrixCode != "" {
			fmt.Fprintln(&source, "\t"+block.matrixCode)
		} else if block.source != "" {
			fmt.Fprintf(&source, "\tif _, err := runtime.RunSourceInContext(ctx, compatibilitySource%04d); err != nil { panic(fmt.Errorf(\"R block %d: %%w\", err)) }\n", i+1, i+1)
		} else {
			fmt.Fprintf(&source, "\tif _, err := runtime.RunEncodedProgramInContext(ctx, compatibilityBlock%04d); err != nil { panic(fmt.Errorf(\"R block %d: %%w\", err)) }\n", i+1, i+1)
		}
	}
	source.WriteString("}\n")
	out, err := format.Source([]byte(source.String()))
	return out, err == nil
}

func readableGoString(source string) string {
	if !strings.Contains(source, "`") {
		return "`" + source + "`"
	}
	return strconv.Quote(source)
}

func renderPartialNativeMain(e *nativeEmitter) ([]byte, bool) {
	imports := "\"fmt\"\n\"math\"\n\"strconv\"\n"
	helpers := ""
	if e.usesVector {
		helpers = `
func rVectorBinary(op string, a, b []float64) []float64 {
	n := len(a); if len(b) > n { n = len(b) }
	if n == 0 || len(a) == 0 || len(b) == 0 { return []float64{} }
	out := make([]float64, n)
	for i := range out { x,y := a[i%len(a)],b[i%len(b)]; switch op { case "+": out[i]=x+y; case "-": out[i]=x-y; case "*": out[i]=x*y; case "/": out[i]=x/y; case "^": out[i]=math.Pow(x,y) } }
	return out
}
func rScalarVector(x float64) []float64 { return []float64{x} }
func rPrintVector(x []float64) { fmt.Print("[1]"); for _,v := range x { fmt.Print(" ",rNumber(v)) }; fmt.Println() }
`
	}
	if e.usesRuntime {
		imports += "\"r2go/runtime\"\n"
		helpers += "\nvar rContext = runtime.NewContext()\n"
	}
	source := "package main\n\n// r2go: partially native-lowered; unsupported R is preserved as comments\nimport (\n" + imports + ")\n\nvar _ = fmt.Println\nvar _ = math.Pow\nfunc rNumber(x float64) string{return strconv.FormatFloat(x,'g',-1,64)}\n" + helpers + e.functions.String() + "func main(){\n" + e.body.String() + "}\n"
	out, err := format.Source([]byte(source))
	return out, err == nil
}
