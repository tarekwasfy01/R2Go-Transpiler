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

type NativeParameter struct {
	Name    string
	Default func(*Context, *Environment) (Value, error)
}

func NewNativeFunction(parent *Environment, parameters []NativeParameter, body func(*Context, *Environment) (Value, error)) Value {
	formals := make([]syntax.Parameter, len(parameters))
	defaults := make(map[string]func(*Context, *Environment) (Value, error))
	for i, parameter := range parameters {
		formals[i] = syntax.Parameter{Name: parameter.Name}
		if parameter.Default != nil {
			defaults[parameter.Name] = parameter.Default
		}
	}
	return &Closure{Parameters: formals, Env: parent, NativeBody: body, Defaults: defaults}
}

func SetLocal(env *Environment, name string, value Value) Value {
	if env == nil {
		panic("r2go: nil local environment")
	}
	env.Set(name, value)
	return value
}

func MustLookup(ctx *Context, env *Environment, name string) Value {
	if env == nil {
		return MustGlobal(ctx, name)
	}
	value, err := env.Get(ctx, name)
	if err != nil {
		panic(err)
	}
	return value
}

// Call invokes either a matrix primitive, a Base-R evaluator family, or a
// user closure already stored in the shared context. Arguments are already
// Go runtime values; no R source or serialized syntax is involved.
func Call(ctx *Context, name string, arguments ...PrimitiveArgument) (Value, error) {
	return CallIn(ctx, nil, name, arguments...)
}

func CallIn(ctx *Context, parent *Environment, name string, arguments ...PrimitiveArgument) (Value, error) {
	if ctx == nil {
		ctx = NewContext()
	}
	if parent == nil {
		parent = ctx.Global
	}
	env := NewEnvironment(parent)
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

func MustCallIn(ctx *Context, env *Environment, name string, arguments ...PrimitiveArgument) Value {
	value, err := CallIn(ctx, env, name, arguments...)
	if err != nil {
		panic(err)
	}
	return value
}

func Elements(value Value) []Value { return elements(value) }

func MustReplace(ctx *Context, env *Environment, variable string, indices []Value, replacement Value) Value {
	object := cloneValue(MustLookup(ctx, env, variable))
	var positions []int
	var err error
	if dims, matrix := dimensions(object); matrix && len(indices) > 1 {
		if len(indices) != len(dims) {
			panic("incorrect number of dimensions")
		}
		axes := make([][]int, len(dims))
		dimnames, _ := Attributes(object)["dimnames"].(*List)
		for axis, index := range indices {
			if index == nil {
				axes[axis] = rangePositions(dims[axis])
				continue
			}
			var names []string
			if dimnames != nil && axis < len(dimnames.Data) {
				if labels, ok := dimnames.Data[axis].(*CharacterVector); ok {
					names = labels.Data
				}
			}
			axes[axis], err = subsetPositionsDimension(dims[axis], index, names)
			if err != nil {
				panic(err)
			}
		}
		positions = cartesianColumnMajor(axes, dims)
	} else {
		if len(indices) != 1 || indices[0] == nil {
			panic("replacement requires one subscript")
		}
		positions, err = subsetPositions(object, indices[0])
		if err != nil {
			panic(err)
		}
	}
	object, err = replacePositions(object, positions, replacement)
	if err != nil {
		panic(err)
	}
	if env == nil {
		return SetGlobal(ctx, variable, object)
	}
	return SetLocal(env, variable, object)
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
	return MustSetMemberIn(ctx, nil, variable, member, value)
}

func MustSetMemberIn(ctx *Context, env *Environment, variable, member string, value Value) Value {
	descriptor, ok := ExecutionPlanByName["$<-"]
	if !ok || descriptor.ReplacementPolicy != "member" {
		panic("r2go: member replacement matrix is unavailable")
	}
	object := cloneValue(MustLookup(ctx, env, variable))
	if err := replaceDollar(object, member, value); err != nil {
		panic(err)
	}
	if env == nil {
		return SetGlobal(ctx, variable, object)
	}
	return SetLocal(env, variable, object)
}

func MustSetDimensionNames(ctx *Context, variable, axis string, value Value) Value {
	return MustSetDimensionNamesIn(ctx, nil, variable, axis, value)
}

func MustSetDimensionNamesIn(ctx *Context, env *Environment, variable, axis string, value Value) Value {
	object := cloneValue(MustLookup(ctx, env, variable))
	if err := setDimensionNamesComponent(object, axis, value); err != nil {
		panic(err)
	}
	if env == nil {
		return SetGlobal(ctx, variable, object)
	}
	return SetLocal(env, variable, object)
}
