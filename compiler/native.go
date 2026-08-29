package compiler

import (
	"fmt"
	"go/format"
	rruntime "r2go/runtime"
	"r2go/syntax"
	"strconv"
	"strings"
)

type nativeKind int

const (
	nativeNumber nativeKind = iota
	nativeLogical
	nativeString
	nativeVector
	nativeValue
)

type nativeExpression struct {
	code string
	kind nativeKind
}
type nativeEmitter struct {
	body         strings.Builder
	functions    strings.Builder
	variables    map[string]nativeKind
	functionSet  map[string]nativeFunction
	declared     map[string]bool
	indent       int
	controlDepth int
	allowDeclare bool
	temporary    int
	usesVector   bool
	usesRuntime  bool
}

type nativeFunction struct {
	parameters []nativeKind
	result     nativeKind
}

func generateNativeMain(p *syntax.Program) ([]byte, bool) {
	e := &nativeEmitter{variables: map[string]nativeKind{}, declared: map[string]bool{}, functionSet: map[string]nativeFunction{}, indent: 1, allowDeclare: true}
	for i, x := range p.Expressions {
		if !e.statement(x, i == len(p.Expressions)-1) {
			return nil, false
		}
	}
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
	source := "package main\n\n// r2go: native-lowered\nimport (\n" + imports + ")\n\nvar _ = math.Pow\nfunc rNumber(x float64) string{return strconv.FormatFloat(x,'g',-1,64)}\n" + helpers + e.functions.String() + "func main(){\n" + e.body.String() + "}\n"
	out, err := format.Source([]byte(source))
	return out, err == nil
}
func (e *nativeEmitter) line(s string) {
	e.body.WriteString(strings.Repeat("\t", e.indent))
	e.body.WriteString(s)
	e.body.WriteByte('\n')
}
func (e *nativeEmitter) statement(x syntax.Expr, visible bool) bool {
	if call, ok := x.(*syntax.Call); ok {
		if fn, ok := call.Function.(*syntax.Symbol); ok {
			switch fn.Name {
			case "<-", "=":
				if len(call.Arguments) != 2 {
					return false
				}
				sym, ok := call.Arguments[0].Value.(*syntax.Symbol)
				if !ok {
					return false
				}
				if function, isFunction := call.Arguments[1].Value.(*syntax.Function); isFunction {
					return e.emitFunction(sym.Name, function)
				}
				v, ok := e.expression(call.Arguments[1].Value)
				if !ok {
					return false
				}
				if old, exists := e.variables[sym.Name]; exists && old != v.kind {
					return false
				}
				if !e.declared[sym.Name] && !e.allowDeclare {
					return false
				}
				e.variables[sym.Name] = v.kind
				op := ":="
				if e.declared[sym.Name] {
					op = "="
				}
				e.declared[sym.Name] = true
				e.line(goName(sym.Name) + " " + op + " " + v.code)
				return true
			case "print":
				if len(call.Arguments) != 1 {
					return false
				}
				v, ok := e.expression(call.Arguments[0].Value)
				if !ok {
					return false
				}
				e.line(nativePrint(v))
				return true
			case "break", "next":
				if e.controlDepth == 0 {
					return false
				}
				e.line(fn.Name)
				return true
			}
		}
	}
	if block, ok := x.(*syntax.Block); ok {
		for i, item := range block.Expressions {
			if !e.statement(item, visible && i == len(block.Expressions)-1) {
				return false
			}
		}
		return true
	}
	if branch, ok := x.(*syntax.If); ok {
		return e.ifStatement(branch, visible)
	}
	if loop, ok := x.(*syntax.While); ok {
		return e.whileStatement(loop)
	}
	if loop, ok := x.(*syntax.Repeat); ok {
		return e.repeatStatement(loop)
	}
	if loop, ok := x.(*syntax.For); ok {
		return e.forStatement(loop)
	}
	v, ok := e.expression(x)
	if !ok {
		return false
	}
	if visible {
		e.line(nativePrint(v))
	} else {
		e.line("_ = " + v.code)
	}
	return true
}

func (e *nativeEmitter) emitFunction(name string, function *syntax.Function) bool {
	if e.controlDepth != 0 || e.declared[name] || len(function.Parameters) == 0 && name == "" {
		return false
	}
	parameterKinds := make([]nativeKind, len(function.Parameters))
	child := &nativeEmitter{
		variables:    map[string]nativeKind{},
		declared:     map[string]bool{},
		functionSet:  make(map[string]nativeFunction, len(e.functionSet)+1),
		indent:       1,
		allowDeclare: true,
	}
	for key, value := range e.functionSet {
		child.functionSet[key] = value
	}
	for i, parameter := range function.Parameters {
		if parameter.Name == "..." || parameter.Default != nil {
			return false
		}
		parameterKinds[i] = nativeNumber
		child.variables[parameter.Name] = nativeNumber
		child.declared[parameter.Name] = true
	}
	// Register a numeric provisional signature so direct numeric recursion can
	// be lowered without routing the entire program through the IR evaluator.
	child.functionSet[name] = nativeFunction{parameters: parameterKinds, result: nativeNumber}
	result, ok := child.returnExpression(function.Body)
	if !ok || result.kind == nativeVector || result.kind == nativeValue {
		return false
	}
	signature := nativeFunction{parameters: parameterKinds, result: result.kind}
	e.functionSet[name] = signature
	params := make([]string, len(function.Parameters))
	for i, parameter := range function.Parameters {
		params[i] = goName(parameter.Name) + " float64"
	}
	e.functions.WriteString("func " + goName(name) + "(" + strings.Join(params, ",") + ") " + nativeGoType(result.kind) + " {\n")
	e.functions.WriteString(child.body.String())
	e.functions.WriteString("\treturn " + result.code + "\n}\n")
	e.usesVector = e.usesVector || child.usesVector
	e.usesRuntime = e.usesRuntime || child.usesRuntime
	return true
}

func (e *nativeEmitter) returnExpression(body syntax.Expr) (nativeExpression, bool) {
	if block, ok := body.(*syntax.Block); ok {
		if len(block.Expressions) == 0 {
			return nativeExpression{}, false
		}
		for _, item := range block.Expressions[:len(block.Expressions)-1] {
			if !e.statement(item, false) {
				return nativeExpression{}, false
			}
		}
		return e.expression(block.Expressions[len(block.Expressions)-1])
	}
	return e.expression(body)
}

func nativeGoType(kind nativeKind) string {
	switch kind {
	case nativeNumber:
		return "float64"
	case nativeLogical:
		return "bool"
	case nativeString:
		return "string"
	default:
		return "any"
	}
}
func (e *nativeEmitter) controlledBody(x syntax.Expr, visible bool) bool {
	old := e.allowDeclare
	e.allowDeclare = false
	e.controlDepth++
	ok := e.statement(x, visible)
	e.controlDepth--
	e.allowDeclare = old
	return ok
}
func (e *nativeEmitter) ifStatement(x *syntax.If, visible bool) bool {
	cond, ok := e.expression(x.Condition)
	if !ok || cond.kind != nativeLogical {
		return false
	}
	old := e.allowDeclare
	e.allowDeclare = false
	e.line("if " + cond.code + " {")
	e.indent++
	if !e.statement(x.Then, visible) {
		return false
	}
	e.indent--
	if x.Else != nil {
		e.line("} else {")
		e.indent++
		if !e.statement(x.Else, visible) {
			return false
		}
		e.indent--
	}
	e.line("}")
	e.allowDeclare = old
	return true
}
func (e *nativeEmitter) whileStatement(x *syntax.While) bool {
	cond, ok := e.expression(x.Condition)
	if !ok || cond.kind != nativeLogical {
		return false
	}
	e.line("for " + cond.code + " {")
	e.indent++
	if !e.controlledBody(x.Body, false) {
		return false
	}
	e.indent--
	e.line("}")
	return true
}
func (e *nativeEmitter) repeatStatement(x *syntax.Repeat) bool {
	e.line("for {")
	e.indent++
	if !e.controlledBody(x.Body, false) {
		return false
	}
	e.indent--
	e.line("}")
	return true
}
func (e *nativeEmitter) forStatement(x *syntax.For) bool {
	call, ok := x.Sequence.(*syntax.Call)
	if !ok {
		return false
	}
	fn, ok := call.Function.(*syntax.Symbol)
	if !ok || fn.Name != ":" || len(call.Arguments) != 2 {
		return false
	}
	from, ok := e.expression(call.Arguments[0].Value)
	if !ok || from.kind != nativeNumber {
		return false
	}
	to, ok := e.expression(call.Arguments[1].Value)
	if !ok || to.kind != nativeNumber {
		return false
	}
	e.temporary++
	id := strconv.Itoa(e.temporary)
	a, b, step := "r_from"+id, "r_to"+id, "r_step"+id
	e.line(a + " := " + from.code)
	e.line(b + " := " + to.code)
	e.line(step + " := 1.0")
	e.line("if " + a + " > " + b + " { " + step + " = -1 }")
	if !e.declared[x.Variable] {
		e.line("var " + goName(x.Variable) + " float64")
	}
	e.variables[x.Variable] = nativeNumber
	e.declared[x.Variable] = true
	e.line("for " + goName(x.Variable) + " = " + a + "; (" + step + ">0 && " + goName(x.Variable) + "<=" + b + ") || (" + step + "<0 && " + goName(x.Variable) + ">=" + b + "); " + goName(x.Variable) + " += " + step + " {")
	e.indent++
	if !e.controlledBody(x.Body, false) {
		return false
	}
	e.indent--
	e.line("}")
	return true
}
func (e *nativeEmitter) expression(x syntax.Expr) (nativeExpression, bool) {
	switch n := x.(type) {
	case *syntax.Literal:
		switch n.Kind {
		case syntax.NumberLiteral:
			if strings.HasSuffix(n.Text, "L") || strings.HasSuffix(n.Text, "i") || n.Text == "Inf" || n.Text == "NaN" {
				return nativeExpression{}, false
			}
			return nativeExpression{n.Text, nativeNumber}, true
		case syntax.LogicalLiteral:
			return nativeExpression{strconv.FormatBool(n.Text == "TRUE" || n.Text == "T"), nativeLogical}, true
		case syntax.StringLiteral:
			return nativeExpression{strconv.Quote(n.Text), nativeString}, true
		}
	case *syntax.Symbol:
		k, ok := e.variables[n.Name]
		return nativeExpression{goName(n.Name), k}, ok
	case *syntax.Call:
		fn, ok := n.Function.(*syntax.Symbol)
		if !ok {
			return nativeExpression{}, false
		}
		if fn.Name == "(" && len(n.Arguments) == 1 {
			return e.expression(n.Arguments[0].Value)
		}
		if function, exists := e.functionSet[fn.Name]; exists {
			if len(n.Arguments) != len(function.parameters) {
				return nativeExpression{}, false
			}
			arguments := make([]string, len(n.Arguments))
			for i, argument := range n.Arguments {
				value, ok := e.expression(argument.Value)
				if !ok || value.kind != function.parameters[i] || argument.Name != "" {
					return nativeExpression{}, false
				}
				arguments[i] = value.code
			}
			return nativeExpression{goName(fn.Name) + "(" + strings.Join(arguments, ",") + ")", function.result}, true
		}
		if len(n.Arguments) == 1 && (fn.Name == "+" || fn.Name == "-") {
			a, ok := e.expression(n.Arguments[0].Value)
			if !ok || a.kind != nativeNumber {
				return nativeExpression{}, false
			}
			return nativeExpression{"(" + fn.Name + a.code + ")", nativeNumber}, true
		}
		if fn.Name == "!" && len(n.Arguments) == 1 {
			a, ok := e.expression(n.Arguments[0].Value)
			if !ok || a.kind != nativeLogical {
				return nativeExpression{}, false
			}
			return nativeExpression{"(!" + a.code + ")", nativeLogical}, true
		}
		if fn.Name == "c" {
			parts := make([]string, len(n.Arguments))
			for i, argument := range n.Arguments {
				value, ok := e.expression(argument.Value)
				if !ok || value.kind != nativeNumber {
					return nativeExpression{}, false
				}
				parts[i] = value.code
			}
			e.usesVector = true
			return nativeExpression{"[]float64{" + strings.Join(parts, ",") + "}", nativeVector}, true
		}
		if rruntime.PrimitiveKnown(fn.Name) && !nativeOperator(fn.Name) {
			arguments := make([]string, len(n.Arguments))
			for i, argument := range n.Arguments {
				value, ok := e.expression(argument.Value)
				if !ok {
					return nativeExpression{}, false
				}
				converted := nativeRuntimeValueCode(value)
				if argument.Name == "" {
					arguments[i] = "runtime.PrimitiveArg(" + converted + ")"
				} else {
					arguments[i] = "runtime.NamedPrimitiveArg(" + strconv.Quote(argument.Name) + "," + converted + ")"
				}
			}
			e.usesRuntime = true
			return nativeExpression{"runtime.MustCallPrimitive(rContext," + strconv.Quote(fn.Name) + optionalComma(arguments) + strings.Join(arguments, ",") + ")", nativeValue}, true
		}
		if len(n.Arguments) != 2 {
			return nativeExpression{}, false
		}
		a, ok := e.expression(n.Arguments[0].Value)
		if !ok {
			return nativeExpression{}, false
		}
		b, ok := e.expression(n.Arguments[1].Value)
		if !ok {
			return nativeExpression{}, false
		}
		switch fn.Name {
		case "+", "-", "*", "/":
			if a.kind == nativeValue || b.kind == nativeValue {
				e.usesRuntime = true
				return nativeExpression{"runtime.MustBinary(" + strconv.Quote(fn.Name) + "," + nativeRuntimeValueCode(a) + "," + nativeRuntimeValueCode(b) + ")", nativeValue}, true
			}
			if a.kind == nativeVector || b.kind == nativeVector {
				e.usesVector = true
				return nativeExpression{"rVectorBinary(" + strconv.Quote(fn.Name) + "," + nativeVectorCode(a) + "," + nativeVectorCode(b) + ")", nativeVector}, true
			}
			if a.kind != nativeNumber || b.kind != nativeNumber {
				return nativeExpression{}, false
			}
			return nativeExpression{"(" + a.code + fn.Name + b.code + ")", nativeNumber}, true
		case "^":
			if a.kind == nativeValue || b.kind == nativeValue {
				e.usesRuntime = true
				return nativeExpression{"runtime.MustBinary(\"^\"," + nativeRuntimeValueCode(a) + "," + nativeRuntimeValueCode(b) + ")", nativeValue}, true
			}
			if a.kind == nativeVector || b.kind == nativeVector {
				e.usesVector = true
				return nativeExpression{"rVectorBinary(\"^\"," + nativeVectorCode(a) + "," + nativeVectorCode(b) + ")", nativeVector}, true
			}
			return nativeExpression{"math.Pow(" + a.code + "," + b.code + ")", nativeNumber}, a.kind == nativeNumber && b.kind == nativeNumber
		case "%%":
			return nativeExpression{"(" + a.code + "-" + b.code + "*math.Floor(" + a.code + "/" + b.code + "))", nativeNumber}, a.kind == nativeNumber && b.kind == nativeNumber
		case "%/%":
			return nativeExpression{"math.Floor(" + a.code + "/" + b.code + ")", nativeNumber}, a.kind == nativeNumber && b.kind == nativeNumber
		case "==", "!=", "<", "<=", ">", ">=":
			if a.kind == nativeValue || b.kind == nativeValue {
				e.usesRuntime = true
				return nativeExpression{"runtime.MustBinary(" + strconv.Quote(fn.Name) + "," + nativeRuntimeValueCode(a) + "," + nativeRuntimeValueCode(b) + ")", nativeValue}, true
			}
			if a.kind == nativeVector || b.kind == nativeVector {
				return nativeExpression{}, false
			}
			return nativeExpression{"(" + a.code + fn.Name + b.code + ")", nativeLogical}, a.kind == b.kind
		case "&&", "||":
			return nativeExpression{"(" + a.code + fn.Name + b.code + ")", nativeLogical}, a.kind == nativeLogical && b.kind == nativeLogical
		}
	}
	return nativeExpression{}, false
}
func nativePrint(v nativeExpression) string {
	switch v.kind {
	case nativeNumber:
		return "fmt.Println(rNumber(" + v.code + "))"
	case nativeLogical:
		return "fmt.Println(map[bool]string{true:\"TRUE\",false:\"FALSE\"}[" + v.code + "])"
	case nativeString:
		return "fmt.Println(strconv.Quote(" + v.code + "))"
	case nativeVector:
		return "rPrintVector(" + v.code + ")"
	case nativeValue:
		return "fmt.Println(" + v.code + ".String())"
	}
	panic("invalid native kind")
}

func nativeVectorCode(value nativeExpression) string {
	if value.kind == nativeVector {
		return value.code
	}
	return "rScalarVector(" + value.code + ")"
}

func nativeRuntimeValueCode(value nativeExpression) string {
	switch value.kind {
	case nativeValue:
		return value.code
	case nativeVector:
		return "runtime.NumericVector(" + value.code + "...)"
	case nativeNumber:
		return "runtime.NumericVector(" + value.code + ")"
	case nativeLogical:
		return "runtime.LogicalScalar(" + value.code + ")"
	case nativeString:
		return "runtime.CharacterScalar(" + value.code + ")"
	default:
		return value.code
	}
}

func optionalComma(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return ","
}

func nativeOperator(name string) bool {
	switch name {
	case "+", "-", "*", "/", "^", "%%", "%/%", "==", "!=", "<", "<=", ">", ">=", "&&", "||", "!", "(", "c":
		return true
	default:
		return false
	}
}
func goName(s string) string {
	var b strings.Builder
	b.WriteString("r_")
	for _, r := range s {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_' {
			b.WriteRune(r)
		} else {
			fmt.Fprintf(&b, "_u%x_", r)
		}
	}
	return b.String()
}
