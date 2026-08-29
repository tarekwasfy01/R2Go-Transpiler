package runtime

import (
	"fmt"
	"r2go/syntax"
)

// PrimitiveArgument preserves R's optional argument names across generated Go
// calls into the central primitive matrix.
type PrimitiveArgument struct {
	Name    string
	Value   Value
	Omitted bool
}

func PrimitiveArg(value Value) PrimitiveArgument { return PrimitiveArgument{Value: value} }
func NamedPrimitiveArg(name string, value Value) PrimitiveArgument {
	return PrimitiveArgument{Name: name, Value: value}
}
func OmittedPrimitiveArg() PrimitiveArgument { return PrimitiveArgument{Omitted: true} }

func PrimitiveKnown(name string) bool {
	_, ok := PrimitiveByName[name]
	return ok
}

// CallPrimitive routes a generated, already-evaluated argument vector through
// the same matrix dispatcher used by the interpreter, without embedding IR.
func CallPrimitive(ctx *Context, name string, arguments ...PrimitiveArgument) (Value, error) {
	if ctx == nil {
		ctx = NewContext()
	}
	if !PrimitiveKnown(name) {
		return nil, fmt.Errorf("unknown R primitive %s", name)
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
		symbol := fmt.Sprintf(".r2go_arg_%d", i)
		env.Set(symbol, argument.Value)
		args[i] = syntax.Argument{Name: argument.Name, Value: &syntax.Symbol{Name: symbol}}
	}
	return ctx.evalCall(&syntax.Call{
		Function:  &syntax.Symbol{Name: name},
		Arguments: args,
	}, env)
}

func MustCallPrimitive(ctx *Context, name string, arguments ...PrimitiveArgument) Value {
	value, err := CallPrimitive(ctx, name, arguments...)
	if err != nil {
		panic(err)
	}
	return value
}
