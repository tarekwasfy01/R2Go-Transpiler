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
type matrixNativeEmitter struct{}

func (matrixNativeEmitter) statement(expression syntax.Expr) (string, bool) {
	switch node := expression.(type) {
	case *syntax.Block:
		parts := make([]string, len(node.Expressions))
		for i, child := range node.Expressions {
			code, ok := (matrixNativeEmitter{}).statement(child)
			if !ok {
				return "", false
			}
			parts[i] = code
		}
		return strings.Join(parts, "\n"), true
	case *syntax.If:
		condition, ok := (matrixNativeEmitter{}).expression(node.Condition)
		if !ok {
			return "", false
		}
		thenCode, ok := (matrixNativeEmitter{}).statement(node.Then)
		if !ok {
			return "", false
		}
		code := "if runtime.MustTrue(" + condition + ") {\n" + indentMatrixCode(thenCode) + "\n}"
		if node.Else != nil {
			elseCode, ok := (matrixNativeEmitter{}).statement(node.Else)
			if !ok {
				return "", false
			}
			code += " else {\n" + indentMatrixCode(elseCode) + "\n}"
		}
		return code, true
	case *syntax.While:
		condition, ok := (matrixNativeEmitter{}).expression(node.Condition)
		if !ok {
			return "", false
		}
		body, ok := (matrixNativeEmitter{}).statement(node.Body)
		if !ok {
			return "", false
		}
		return "for runtime.MustTrue(" + condition + ") {\n" + indentMatrixCode(body) + "\n}", true
	case *syntax.Repeat:
		body, ok := (matrixNativeEmitter{}).statement(node.Body)
		if !ok {
			return "", false
		}
		return "for {\n" + indentMatrixCode(body) + "\n}", true
	}
	if call, ok := expression.(*syntax.Call); ok {
		if symbol, ok := call.Function.(*syntax.Symbol); ok && (symbol.Name == "<-" || symbol.Name == "=") {
			if len(call.Arguments) != 2 {
				return "", false
			}
			if _, closure := call.Arguments[1].Value.(*syntax.Function); closure {
				return "", false
			}
			value, ok := (matrixNativeEmitter{}).expression(call.Arguments[1].Value)
			if !ok {
				return "", false
			}
			switch target := call.Arguments[0].Value.(type) {
			case *syntax.Symbol:
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
					return "runtime.MustSetMember(ctx, " + strconv.Quote(root.Name) + ", " + strconv.Quote(member.Name) + ", " + value + ")", true
				}
				if (name.Name == "rownames" || name.Name == "colnames") && len(target.Arguments) == 1 {
					return "runtime.MustSetDimensionNames(ctx, " + strconv.Quote(root.Name) + ", " + strconv.Quote(name.Name) + ", " + value + ")", true
				}
			}
			return "", false
		}
		if symbol, ok := call.Function.(*syntax.Symbol); ok {
			switch symbol.Name {
			case "break":
				return "break", len(call.Arguments) == 0
			case "next":
				return "continue", len(call.Arguments) == 0
			}
		}
	}
	value, ok := (matrixNativeEmitter{}).expression(expression)
	if !ok {
		return "", false
	}
	return "_ = " + value, true
}

func (matrixNativeEmitter) expression(expression syntax.Expr) (string, bool) {
	switch value := expression.(type) {
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
				converted, convertedOK = (matrixNativeEmitter{}).expression(argument.Value)
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
		comma := ""
		if len(arguments) > 0 {
			comma = ", "
		}
		return caller + "(ctx, " + strconv.Quote(function.Name) + comma + strings.Join(arguments, ", ") + ")", true
	}
	return "", false
}

func indentMatrixCode(code string) string {
	return "\t" + strings.ReplaceAll(code, "\n", "\n\t")
}
