package runtime

import (
	"fmt"
	converted "r2go/converted/rgo"
	"r2go/syntax"
)

// executeConvertedEntry is the sole boundary between the existing evaluator
// and the independently converted GNU-R C entry-point package.  The manifest
// gate prevents native-ABI and explicitly partial conversions from becoming
// silently executable primitives.
func (c *Context) executeConvertedEntry(plan ExecutionPlan, args []syntax.Argument, env *Environment) (Value, bool, error) {
	if !converted.Eligible(plan.CEntry) {
		return nil, false, nil
	}
	values := make([]converted.Value, len(args))
	for i, argument := range args {
		if argument.Value == nil {
			return nil, true, fmt.Errorf("missing argument %d", i+1)
		}
		value, err := c.Eval(argument.Value, env)
		if err != nil {
			return nil, true, err
		}
		values[i], err = c.toConvertedValue(value)
		if err != nil {
			return nil, true, err
		}
	}
	convertedEnv, err := c.toConvertedEnvironment(env)
	if err != nil {
		return nil, true, err
	}
	output, found := converted.Dispatch(plan.CEntry, converted.Nil, converted.Nil, converted.Lists(values...), convertedEnv)
	if !found {
		return nil, false, nil
	}
	value, err := c.fromConvertedValue(output)
	return value, true, err
}

func (c *Context) toConvertedEnvironment(env *Environment) (converted.Value, error) {
	if env == nil {
		return converted.Nil, nil
	}
	parent, err := c.toConvertedEnvironment(env.Parent)
	if err != nil {
		return converted.Nil, err
	}
	result := converted.NewEnv(nil, "")
	if parent.Kind == converted.Environment {
		result.Env.Parent = parent.Env
	}
	for _, name := range env.Names(true) {
		value, err := env.Get(c, name)
		if err != nil {
			return converted.Nil, err
		}
		convertedValue, err := c.toConvertedValue(value)
		if err != nil {
			return converted.Nil, err
		}
		result.Env.Values[name] = convertedValue
	}
	return result, nil
}

func (c *Context) toConvertedValue(value Value) (converted.Value, error) {
	var result converted.Value
	switch v := value.(type) {
	case nil, Null:
		return converted.Nil, nil
	case *RawVector:
		result = converted.Raws(v.Data...)
	case *LogicalVector:
		data, missing := make([]int8, len(v.Data)), make([]bool, len(v.Data))
		for i, item := range v.Data {
			if item == NA {
				missing[i] = true
			} else {
				data[i] = int8(item)
			}
		}
		result = converted.Value{Kind: converted.Logical, L: data, NA: missing}
	case *IntegerVector:
		result = converted.Value{Kind: converted.Integer, I: append([]int64(nil), v.Data...), NA: append([]bool(nil), v.Missing...)}
	case *DoubleVector:
		result = converted.Value{Kind: converted.Double, D: append([]float64(nil), v.Data...), NA: append([]bool(nil), v.Missing...)}
	case *ComplexVector:
		result = converted.Value{Kind: converted.Complex, Z: append([]complex128(nil), v.Data...), NA: append([]bool(nil), v.Missing...)}
	case *CharacterVector:
		result = converted.Value{Kind: converted.String, S: append([]string(nil), v.Data...), NA: append([]bool(nil), v.Missing...)}
	case *List:
		items := make([]converted.Value, len(v.Data))
		for i, item := range v.Data {
			var err error
			items[i], err = c.toConvertedValue(item)
			if err != nil {
				return converted.Nil, err
			}
		}
		result = converted.Lists(items...)
		if len(v.Names) != 0 {
			result.Attr = map[string]converted.Value{"names": converted.Strings(v.Names...)}
		}
	case *EnvironmentValue:
		return c.toConvertedEnvironment(v.Env)
	default:
		return converted.Nil, fmt.Errorf("converted entry does not support %s values", value.Kind())
	}
	attrs, err := c.toConvertedAttributes(attributesOf(value))
	if err != nil {
		return converted.Nil, err
	}
	for key, attr := range attrs {
		if result.Attr == nil {
			result.Attr = map[string]converted.Value{}
		}
		result.Attr[key] = attr
	}
	return result, nil
}

func (c *Context) toConvertedAttributes(attrs map[string]Value) (map[string]converted.Value, error) {
	result := map[string]converted.Value{}
	for key, value := range attrs {
		convertedValue, err := c.toConvertedValue(value)
		if err != nil {
			return nil, err
		}
		result[key] = convertedValue
	}
	return result, nil
}

func (c *Context) fromConvertedValue(value converted.Value) (Value, error) {
	if value.Kind == converted.Error {
		if value.Err != nil {
			return nil, value.Err
		}
		return nil, fmt.Errorf("converted entry returned an error")
	}
	var result Value
	switch value.Kind {
	case converted.Null:
		return NullValue, nil
	case converted.Raw:
		result = &RawVector{Data: append([]byte(nil), value.B...)}
	case converted.Logical:
		data := make([]Logical, len(value.L))
		for i, item := range value.L {
			if i < len(value.NA) && value.NA[i] {
				data[i] = NA
			} else {
				data[i] = Logical(item)
			}
		}
		result = &LogicalVector{Data: data}
	case converted.Integer:
		result = &IntegerVector{Data: append([]int64(nil), value.I...), Missing: append([]bool(nil), value.NA...)}
	case converted.Double:
		result = &DoubleVector{Data: append([]float64(nil), value.D...), Missing: append([]bool(nil), value.NA...)}
	case converted.Complex:
		result = &ComplexVector{Data: append([]complex128(nil), value.Z...), Missing: append([]bool(nil), value.NA...)}
	case converted.String:
		result = &CharacterVector{Data: append([]string(nil), value.S...), Missing: append([]bool(nil), value.NA...)}
	case converted.List:
		items := make([]Value, len(value.V))
		for i, item := range value.V {
			var err error
			items[i], err = c.fromConvertedValue(item)
			if err != nil {
				return nil, err
			}
		}
		result = &List{Data: items}
	case converted.Environment:
		return &EnvironmentValue{Name: value.Env.Name}, nil
	default:
		return nil, fmt.Errorf("converted entry returned unsupported value kind %d", value.Kind)
	}
	attrs, err := c.fromConvertedAttributes(value.Attr)
	if err != nil {
		return nil, err
	}
	applyAttributes(result, attrs)
	return result, nil
}

func (c *Context) fromConvertedAttributes(attrs map[string]converted.Value) (map[string]Value, error) {
	result := map[string]Value{}
	for key, value := range attrs {
		runtimeValue, err := c.fromConvertedValue(value)
		if err != nil {
			return nil, err
		}
		result[key] = runtimeValue
	}
	return result, nil
}

func attributesOf(value Value) map[string]Value {
	switch v := value.(type) {
	case *RawVector:
		return v.Attr
	case *LogicalVector:
		return v.Attr
	case *IntegerVector:
		return v.Attr
	case *DoubleVector:
		return v.Attr
	case *ComplexVector:
		return v.Attr
	case *CharacterVector:
		return v.Attr
	case *List:
		return v.Attr
	case *Formula:
		return v.Attr
	default:
		return nil
	}
}

func applyAttributes(value Value, attrs map[string]Value) {
	switch v := value.(type) {
	case *RawVector:
		v.Attr = attrs
	case *LogicalVector:
		v.Attr = attrs
	case *IntegerVector:
		v.Attr = attrs
	case *DoubleVector:
		v.Attr = attrs
	case *ComplexVector:
		v.Attr = attrs
	case *CharacterVector:
		v.Attr = attrs
	case *List:
		v.Attr = attrs
	}
}
