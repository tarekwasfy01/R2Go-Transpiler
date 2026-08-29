package runtime

import (
	"fmt"
	"r2go/syntax"
)

// SetGlobal and MustGlobal are the readable code-generation boundary for
// variables shared by matrix-lowered and compatibility-lowered blocks.
func SetGlobal(ctx *Context, name string, value Value) Value {
	if ctx == nil {
		panic("r2go: nil runtime context")
	}
	ctx.Global.Set(name, value)
	return value
}

func MustGlobal(ctx *Context, name string) Value {
	if ctx == nil {
		panic("r2go: nil runtime context")
	}
	value, err := ctx.Global.Get(ctx, name)
	if err != nil {
		panic(err)
	}
	return value
}

// Call invokes either a matrix primitive, a Base-R evaluator family, or a
// user closure already stored in the shared context. Arguments are already
// Go runtime values; no R source or serialized syntax is involved.
func Call(ctx *Context, name string, arguments ...PrimitiveArgument) (Value, error) {
	if ctx == nil {
		ctx = NewContext()
	}
	env := NewEnvironment(ctx.Global)
	args := make([]syntax.Argument, len(arguments))
	for i, argument := range arguments {
		if argument.Omitted {
			args[i] = syntax.Argument{Name: argument.Name}
			continue
		}
		if (name == "$" || name == "@") && i == 1 {
			if member, ok := argument.Value.(*CharacterVector); ok && len(member.Data) == 1 {
				args[i] = syntax.Argument{Name: argument.Name, Value: &syntax.Symbol{Name: member.Data[0]}}
				continue
			}
		}
		symbol := fmt.Sprintf(".r2go_value_%d", i)
		env.Set(symbol, argument.Value)
		args[i] = syntax.Argument{Name: argument.Name, Value: &syntax.Symbol{Name: symbol}}
	}
	return ctx.evalCall(&syntax.Call{Function: &syntax.Symbol{Name: name}, Arguments: args}, env)
}

func MustCall(ctx *Context, name string, arguments ...PrimitiveArgument) Value {
	value, err := Call(ctx, name, arguments...)
	if err != nil {
		panic(err)
	}
	return value
}

func MissingValue(kind string) Value {
	switch kind {
	case "integer":
		return &IntegerVector{Data: []int64{0}, Missing: []bool{true}}
	case "double", "real":
		return &DoubleVector{Data: []float64{NAReal()}, Missing: []bool{true}}
	case "complex":
		return &ComplexVector{Data: []complex128{0}, Missing: []bool{true}}
	case "character":
		return &CharacterVector{Data: []string{""}, Missing: []bool{true}}
	default:
		return &LogicalVector{Data: []Logical{NA}}
	}
}

func MustTrue(value Value) bool {
	result, err := IsTrue(value)
	if err != nil {
		panic(err)
	}
	return result
}

// The replacement helpers expose the generated replacement-policy matrix as
// readable operations while keeping copy-on-write updates in the runtime.
func MustSetMember(ctx *Context, variable, member string, value Value) Value {
	descriptor, ok := ExecutionPlanByName["$<-"]
	if !ok || descriptor.ReplacementPolicy != "member" {
		panic("r2go: member replacement matrix is unavailable")
	}
	object := cloneValue(MustGlobal(ctx, variable))
	if err := replaceDollar(object, member, value); err != nil {
		panic(err)
	}
	return SetGlobal(ctx, variable, object)
}

func MustSetDimensionNames(ctx *Context, variable, axis string, value Value) Value {
	object := cloneValue(MustGlobal(ctx, variable))
	if err := setDimensionNamesComponent(object, axis, value); err != nil {
		panic(err)
	}
	return SetGlobal(ctx, variable, object)
}
