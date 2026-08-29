package runtime

import (
	"fmt"
	"math"
	"r2go/syntax"
	"strconv"
)

var manualPrimitiveRoutes = map[string]bool{
	".Primitive": true, "browser": true, "deparseRd": true, "format.info": true,
	"formatC": true, "parse": true, "parseLatex": true, "parse_Rd": true, "~": true,
	"mem.maxNSize": true, "mem.maxVSize": true,
}

// executeManualPrimitive supplies the nine entries for which GNU R has no
// extractable C body in R_FunTab.  Each route is explicit: implemented parts
// use the existing Pure-Go parser/value model; parser variants not present in
// the checkout return a precise boundary error rather than a fake parse.
func (c *Context) executeManualPrimitive(name string, args []syntax.Argument, env *Environment) (Value, bool, error) {
	if !manualPrimitiveRoutes[name] {
		return nil, false, nil
	}
	if name == "~" {
		if len(args) == 0 || args[0].Value == nil {
			return nil, true, fmt.Errorf("formula requires a left-hand expression")
		}
		return &Formula{Expr: args[0].Value, Env: env}, true, nil
	}
	values := make([]Value, len(args))
	for i, argument := range args {
		if argument.Value == nil {
			return nil, true, fmt.Errorf("missing argument %d", i+1)
		}
		value, err := c.Eval(argument.Value, env)
		if err != nil {
			return nil, true, err
		}
		values[i] = value
	}
	first := func() Value {
		if len(values) == 0 {
			return NullValue
		}
		return values[0]
	}
	switch name {
	case "browser":
		return NullValue, true, nil // non-interactive Pure-Go browser hook
	case "deparseRd":
		return &CharacterVector{Data: []string{first().String()}}, true, nil
	case "format.info":
		digits := int64(0)
		if len(values) > 1 {
			if n, err := scalarInt(values[1]); err == nil {
				digits = int64(n)
			}
		}
		return &IntegerVector{Data: []int64{int64(Length(first())), digits, 0}}, true, nil
	case "formatC":
		return formatCValue(first()), true, nil
	case "mem.maxNSize", "mem.maxVSize":
		return &DoubleVector{Data: []float64{math.Inf(1)}}, true, nil
	case "parse":
		source, ok := first().(*CharacterVector)
		if !ok || len(source.Data) == 0 {
			return nil, true, fmt.Errorf("parse expects a non-empty character vector")
		}
		program, err := syntax.Parse(source.Data[0])
		if err != nil {
			return nil, true, err
		}
		items := make([]Value, len(program.Expressions))
		for i, expression := range program.Expressions {
			items[i] = &Language{Expr: expression}
		}
		return &List{Data: items}, true, nil
	case ".Primitive", "parseLatex", "parse_Rd":
		return nil, true, fmt.Errorf("%s requires a GNU-R internal subsystem not available in the Pure-Go core", name)
	}
	return nil, true, fmt.Errorf("unhandled manual primitive %s", name)
}

func formatCValue(value Value) Value {
	switch v := value.(type) {
	case *IntegerVector:
		out := make([]string, len(v.Data))
		for i, item := range v.Data {
			out[i] = strconv.FormatInt(item, 10)
		}
		return &CharacterVector{Data: out, Missing: append([]bool(nil), v.Missing...)}
	case *DoubleVector:
		out := make([]string, len(v.Data))
		for i, item := range v.Data {
			out[i] = strconv.FormatFloat(item, 'g', -1, 64)
		}
		return &CharacterVector{Data: out, Missing: append([]bool(nil), v.Missing...)}
	default:
		return &CharacterVector{Data: []string{value.String()}}
	}
}
