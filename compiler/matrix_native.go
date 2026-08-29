package compiler

import (
	"strconv"
	"strings"

	rruntime "r2go/runtime"
	"r2go/syntax"
)

// matrixNativeEmitter lowers a complete top-level R expression into readable
// Go calls over the generated primitive matrix. All values live in the shared
// runtime context, so these blocks can safely surround compatibility blocks.
type matrixNativeEmitter struct {
	environment string
	inFunction  bool
	temporary   *int
}

func newMatrixNativeEmitter() matrixNativeEmitter {
	counter := 0
	return matrixNativeEmitter{temporary: &counter}
}

func (e matrixNativeEmitter) local() bool { return e.environment != "" }

func (e matrixNativeEmitter) nextTemporary() string {
	if e.temporary == nil {
		counter := 0
		e.temporary = &counter
	}
	(*e.temporary)++
	return "r2goItem" + strconv.Itoa(*e.temporary)
}

func (e matrixNativeEmitter) statement(expression syntax.Expr) (string, bool) {
	switch node := expression.(type) {
	case *syntax.Block:
		parts := make([]string, len(node.Expressions))
		for i, child := range node.Expressions {
			code, ok := e.statement(child)
			if !ok {
				return "", false
			}
			parts[i] = code
		}
		return strings.Join(parts, "\n"), true
	case *syntax.If:
		condition, ok := e.expression(node.Condition)
		if !ok {
			return "", false
		}
		thenCode, ok := e.statement(node.Then)
		if !ok {
			return "", false
		}
		code := "if runtime.MustTrue(" + condition + ") {\n" + indentMatrixCode(thenCode) + "\n}"
		if node.Else != nil {
			elseCode, ok := e.statement(node.Else)
			if !ok {
				return "", false
			}
			code += " else {\n" + indentMatrixCode(elseCode) + "\n}"
		}
		return code, true
	case *syntax.While:
		condition, ok := e.expression(node.Condition)
		if !ok {
			return "", false
		}
		body, ok := e.statement(node.Body)
		if !ok {
			return "", false
		}
		return "for runtime.MustTrue(" + condition + ") {\n" + indentMatrixCode(body) + "\n}", true
	case *syntax.Repeat:
		body, ok := e.statement(node.Body)
		if !ok {
			return "", false
		}
		return "for {\n" + indentMatrixCode(body) + "\n}", true
	case *syntax.For:
		sequence, ok := e.expression(node.Sequence)
		if !ok {
			return "", false
		}
		body, ok := e.statement(node.Body)
		if !ok {
			return "", false
		}
		item := e.nextTemporary()
		set := "runtime.SetGlobal(ctx, " + strconv.Quote(node.Variable) + ", " + item + ")"
		if e.local() {
			set = "runtime.SetLocal(" + e.environment + ", " + strconv.Quote(node.Variable) + ", " + item + ")"
		}
		return "for _, " + item + " := range runtime.Elements(" + sequence + ") {\n\t" + set + "\n" + indentMatrixCode(body) + "\n}", true
	}
	if call, ok := expression.(*syntax.Call); ok {
		if symbol, ok := call.Function.(*syntax.Symbol); ok && (symbol.Name == "<-" || symbol.Name == "=") {
			if len(call.Arguments) != 2 {
				return "", false
			}
			value, ok := e.expression(call.Arguments[1].Value)
			if !ok {
				return "", false
			}
			switch target := call.Arguments[0].Value.(type) {
			case *syntax.Symbol:
				if e.local() {
					return "runtime.SetLocal(" + e.environment + ", " + strconv.Quote(target.Name) + ", " + value + ")", true
				}
				return "runtime.SetGlobal(ctx, " + strconv.Quote(target.Name) + ", " + value + ")", true
			case *syntax.Call:
				name, ok := target.Function.(*syntax.Symbol)
				if !ok || len(target.Arguments) == 0 {
					return "", false
				}
				root, ok := target.Arguments[0].Value.(*syntax.Symbol)
				if !ok {
					return "", false
				}
				if name.Name == "$" && len(target.Arguments) == 2 {
					member, ok := target.Arguments[1].Value.(*syntax.Symbol)
					if !ok {
						return "", false
					}
					if e.local() {
						return "runtime.MustSetMemberIn(ctx, " + e.environment + ", " + strconv.Quote(root.Name) + ", " + strconv.Quote(member.Name) + ", " + value + ")", true
					}
					return "runtime.MustSetMember(ctx, " + strconv.Quote(root.Name) + ", " + strconv.Quote(member.Name) + ", " + value + ")", true
				}
				if (name.Name == "rownames" || name.Name == "colnames") && len(target.Arguments) == 1 {
					if e.local() {
						return "runtime.MustSetDimensionNamesIn(ctx, " + e.environment + ", " + strconv.Quote(root.Name) + ", " + strconv.Quote(name.Name) + ", " + value + ")", true
					}
					return "runtime.MustSetDimensionNames(ctx, " + strconv.Quote(root.Name) + ", " + strconv.Quote(name.Name) + ", " + value + ")", true
				}
				if (name.Name == "[" || name.Name == "[[") && len(target.Arguments) >= 2 {
					indices := make([]string, len(target.Arguments)-1)
					for i, argument := range target.Arguments[1:] {
						if argument.Value == nil {
							indices[i] = "nil"
							continue
						}
						index, indexOK := e.expression(argument.Value)
						if !indexOK {
							return "", false
						}
						indices[i] = index
					}
					environment := "nil"
					if e.local() {
						environment = e.environment
					}
					return "runtime.MustReplace(ctx, " + environment + ", " + strconv.Quote(root.Name) + ", []runtime.Value{" + strings.Join(indices, ", ") + "}, " + value + ")", true
				}
			}
			return "", false
		}
		if symbol, ok := call.Function.(*syntax.Symbol); ok {
			switch symbol.Name {
			case "return":
				if !e.inFunction || len(call.Arguments) > 1 {
					return "", false
				}
				value := "runtime.NullValue"
				if len(call.Arguments) == 1 {
					var ok bool
					value, ok = e.expression(call.Arguments[0].Value)
					if !ok {
						return "", false
					}
				}
				return "return " + value + ", nil", true
			case "break":
				return "break", len(call.Arguments) == 0
			case "next":
				return "continue", len(call.Arguments) == 0
			}
		}
	}
	value, ok := e.expression(expression)
	if !ok {
		return "", false
	}
	return "_ = " + value, true
}

func (e matrixNativeEmitter) expression(expression syntax.Expr) (string, bool) {
	switch value := expression.(type) {
	case *syntax.Function:
		return e.function(value)
	case *syntax.Literal:
		switch value.Kind {
		case syntax.NullLiteral:
			return "runtime.NullValue", true
		case syntax.StringLiteral:
			return "runtime.CharacterScalar(" + strconv.Quote(value.Text) + ")", true
		case syntax.LogicalLiteral:
			logical := value.Text == "TRUE" || value.Text == "T"
			return "runtime.LogicalScalar(" + strconv.FormatBool(logical) + ")", true
		case syntax.NALiteral:
			kind := "logical"
			switch value.Text {
			case "NA_integer_":
				kind = "integer"
			case "NA_real_":
				kind = "real"
			case "NA_complex_":
				kind = "complex"
			case "NA_character_":
				kind = "character"
			}
			return "runtime.MissingValue(" + strconv.Quote(kind) + ")", true
		case syntax.NumberLiteral:
			if strings.HasSuffix(value.Text, "i") {
				return "", false
			}
			number := strings.TrimSuffix(value.Text, "L")
			if number == "Inf" || number == "NaN" {
				return "", false
			}
			return "runtime.NumericVector(" + number + ")", true
		}
	case *syntax.Symbol:
		if e.local() {
			return "runtime.MustLookup(ctx, " + e.environment + ", " + strconv.Quote(value.Name) + ")", true
		}
		return "runtime.MustGlobal(ctx, " + strconv.Quote(value.Name) + ")", true
	case *syntax.Call:
		function, ok := value.Function.(*syntax.Symbol)
		if !ok || function.Name == "<-" || function.Name == "=" || function.Name == "<<-" || function.Name == "return" || function.Name == "break" || function.Name == "next" {
			return "", false
		}
		arguments := make([]string, len(value.Arguments))
		for i, argument := range value.Arguments {
			if argument.Value == nil {
				arguments[i] = "runtime.OmittedPrimitiveArg()"
				continue
			}
			converted := ""
			convertedOK := false
			// R member operators capture their right operand as a name rather
			// than evaluating it as a variable.
			if (function.Name == "$" || function.Name == "@") && i == 1 {
				if member, symbolic := argument.Value.(*syntax.Symbol); symbolic {
					converted = "runtime.CharacterScalar(" + strconv.Quote(member.Name) + ")"
					convertedOK = true
				}
			}
			if !convertedOK {
				converted, convertedOK = e.expression(argument.Value)
			}
			if !convertedOK {
				return "", false
			}
			if argument.Name == "" {
				arguments[i] = "runtime.PrimitiveArg(" + converted + ")"
			} else {
				arguments[i] = "runtime.NamedPrimitiveArg(" + strconv.Quote(argument.Name) + ", " + converted + ")"
			}
		}
		caller := "runtime.MustCall"
		if rruntime.PrimitiveKnown(function.Name) {
			caller = "runtime.MustCallPrimitive"
		}
		prefix := "ctx, " + strconv.Quote(function.Name)
		if e.local() {
			caller = "runtime.MustCallIn"
			prefix = "ctx, " + e.environment + ", " + strconv.Quote(function.Name)
		}
		comma := ""
		if len(arguments) > 0 {
			comma = ", "
		}
		return caller + "(" + prefix + comma + strings.Join(arguments, ", ") + ")", true
	}
	return "", false
}

func (e matrixNativeEmitter) function(function *syntax.Function) (string, bool) {
	child := e
	child.environment = "env"
	child.inFunction = true
	parameters := make([]string, len(function.Parameters))
	for i, parameter := range function.Parameters {
		if parameter.Name == "..." {
			return "", false
		}
		entry := "{Name: " + strconv.Quote(parameter.Name)
		if parameter.Default != nil {
			defaultCode, ok := child.expression(parameter.Default)
			if !ok {
				return "", false
			}
			entry += ", Default: func(ctx *runtime.Context, env *runtime.Environment) (runtime.Value, error) { return " + defaultCode + ", nil }"
		}
		parameters[i] = entry + "}"
	}
	body, ok := child.functionBody(function.Body)
	if !ok {
		return "", false
	}
	parent := "ctx.Global"
	if e.local() {
		parent = e.environment
	}
	return "runtime.NewNativeFunction(" + parent + ", []runtime.NativeParameter{" + strings.Join(parameters, ", ") + "}, func(ctx *runtime.Context, env *runtime.Environment) (runtime.Value, error) {\n" + indentMatrixCode(body) + "\n})", true
}

func (e matrixNativeEmitter) functionBody(body syntax.Expr) (string, bool) {
	if block, ok := body.(*syntax.Block); ok {
		if len(block.Expressions) == 0 {
			return "return runtime.NullValue, nil", true
		}
		parts := make([]string, 0, len(block.Expressions))
		for _, expression := range block.Expressions[:len(block.Expressions)-1] {
			code, ok := e.statement(expression)
			if !ok {
				return "", false
			}
			parts = append(parts, code)
		}
		last, ok := e.returnExpression(block.Expressions[len(block.Expressions)-1])
		if !ok {
			return "", false
		}
		parts = append(parts, last)
		return strings.Join(parts, "\n"), true
	}
	return e.returnExpression(body)
}

func (e matrixNativeEmitter) returnExpression(expression syntax.Expr) (string, bool) {
	if call, ok := expression.(*syntax.Call); ok {
		if symbol, ok := call.Function.(*syntax.Symbol); ok && symbol.Name == "return" {
			return e.statement(call)
		}
	}
	if branch, ok := expression.(*syntax.If); ok {
		condition, ok := e.expression(branch.Condition)
		if !ok {
			return "", false
		}
		thenCode, ok := e.returnExpression(branch.Then)
		if !ok {
			return "", false
		}
		elseCode := "return runtime.NullValue, nil"
		if branch.Else != nil {
			elseCode, ok = e.returnExpression(branch.Else)
			if !ok {
				return "", false
			}
		}
		return "if runtime.MustTrue(" + condition + ") {\n" + indentMatrixCode(thenCode) + "\n} else {\n" + indentMatrixCode(elseCode) + "\n}", true
	}
	value, ok := e.expression(expression)
	if !ok {
		return "", false
	}
	return "return " + value + ", nil", true
}

func indentMatrixCode(code string) string {
	return "\t" + strings.ReplaceAll(code, "\n", "\n\t")
}
