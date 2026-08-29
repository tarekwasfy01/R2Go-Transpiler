package runtime

import "fmt"

func cloneAttributeMap(src map[string]Value) map[string]Value {
	if src == nil {
		return nil
	}
	out := make(map[string]Value, len(src))
	for name, value := range src {
		out[name] = value
	}
	return out
}

// inheritUnaryAttributes applies the Math/Math2 shape rule in one place.
// Element-wise unary operations retain names, dimensions and dimnames even
// when their result container changes from integer to double.
func inheritUnaryAttributes(dst, src Value) Value {
	if err := setAttributes(dst, cloneAttributeMap(Attributes(src))); err != nil {
		panic(err)
	}
	return dst
}

// Attributes returns the mutable attribute map owned by an R value. NULL and
// non-object runtime values deliberately have no attribute storage.
func Attributes(v Value) map[string]Value {
	switch x := v.(type) {
	case *RawVector:
		return x.Attr
	case *LogicalVector:
		return x.Attr
	case *IntegerVector:
		return x.Attr
	case *DoubleVector:
		return x.Attr
	case *ComplexVector:
		return x.Attr
	case *CharacterVector:
		return x.Attr
	case *List:
		return x.Attr
	case *Formula:
		return x.Attr
	}
	return nil
}

func setAttributes(v Value, attrs map[string]Value) error {
	switch x := v.(type) {
	case *RawVector:
		x.Attr = attrs
	case *LogicalVector:
		x.Attr = attrs
	case *IntegerVector:
		x.Attr = attrs
	case *DoubleVector:
		x.Attr = attrs
	case *ComplexVector:
		x.Attr = attrs
	case *CharacterVector:
		x.Attr = attrs
	case *List:
		x.Attr = attrs
	case *Formula:
		x.Attr = attrs
	default:
		return fmt.Errorf("cannot set attributes on %s", v.Kind())
	}
	return nil
}

func setAttribute(v Value, name string, value Value) error {
	normalized, err := validateAttribute(v, name, value)
	if err != nil {
		return err
	}
	value = normalized
	attrs := Attributes(v)
	if attrs == nil {
		attrs = map[string]Value{}
		if err := setAttributes(v, attrs); err != nil {
			return err
		}
	}
	if _, null := value.(Null); null {
		delete(attrs, name)
	} else {
		attrs[name] = value
	}
	return nil
}

// validateAttribute is the shared attribute-rule matrix used by direct
// setters, replacement functions, arrays, data frames and factors.
func validateAttribute(v Value, name string, value Value) (Value, error) {
	if _, null := value.(Null); null {
		return value, nil
	}
	switch name {
	case "names":
		n, ok := value.(*CharacterVector)
		if !ok {
			return nil, fmt.Errorf("'names' attribute must be a character vector")
		}
		if len(n.Data) != Length(v) {
			return nil, fmt.Errorf("'names' attribute [%d] must be the same length as the vector [%d]", len(n.Data), Length(v))
		}
	case "dim":
		d, err := numbers(value)
		if err != nil {
			return nil, fmt.Errorf("'dim' attribute must be numeric")
		}
		product := 1
		normalized := &IntegerVector{Data: make([]int64, len(d.Data))}
		for i, n := range d.Data {
			if missingAt(d, i) || n < 0 || n != float64(int(n)) {
				return nil, fmt.Errorf("invalid 'dim' attribute")
			}
			normalized.Data[i] = int64(n)
			product *= int(n)
		}
		if product != Length(v) {
			return nil, fmt.Errorf("dims [product %d] do not match the length of object [%d]", product, Length(v))
		}
		value = normalized
	case "dimnames":
		dims, ok := dimensions(v)
		if !ok {
			return nil, fmt.Errorf("'dimnames' applied to non-array")
		}
		items, ok := value.(*List)
		if !ok || len(items.Data) != len(dims) {
			return nil, fmt.Errorf("length of 'dimnames' [%d] must match number of dimensions [%d]", Length(value), len(dims))
		}
		for axis, item := range items.Data {
			if _, null := item.(Null); null {
				continue
			}
			labels, ok := item.(*CharacterVector)
			if !ok || len(labels.Data) != dims[axis] {
				return nil, fmt.Errorf("length of 'dimnames' component %d does not match extent %d", axis+1, dims[axis])
			}
		}
	case "class", "levels", "comment":
		if _, ok := value.(*CharacterVector); !ok {
			return nil, fmt.Errorf("'%s' attribute must be character", name)
		}
	}
	return value, nil
}

func replaceAttributeMap(v Value, replacement Value) error {
	if _, null := replacement.(Null); null {
		return setAttributes(v, nil)
	}
	items, ok := replacement.(*List)
	if !ok || len(items.Names) != len(items.Data) {
		return fmt.Errorf("attributes must be a named list")
	}
	if err := setAttributes(v, nil); err != nil {
		return err
	}
	for i, name := range items.Names {
		if name == "" {
			return fmt.Errorf("all attributes must be named")
		}
		if name == "names" {
			if err := replaceNames(v, items.Data[i]); err != nil {
				return err
			}
			continue
		}
		if err := setAttribute(v, name, items.Data[i]); err != nil {
			return err
		}
	}
	return nil
}

func cloneValue(v Value) Value {
	switch x := v.(type) {
	case *RawVector:
		return &RawVector{Data: append([]byte(nil), x.Data...), Attr: cloneAttributeMap(x.Attr)}
	case *LogicalVector:
		return &LogicalVector{Data: append([]Logical(nil), x.Data...), Attr: cloneAttributeMap(x.Attr)}
	case *IntegerVector:
		return &IntegerVector{Data: append([]int64(nil), x.Data...), Missing: append([]bool(nil), x.Missing...), Attr: cloneAttributeMap(x.Attr)}
	case *DoubleVector:
		return &DoubleVector{Data: append([]float64(nil), x.Data...), Missing: append([]bool(nil), x.Missing...), Attr: cloneAttributeMap(x.Attr)}
	case *ComplexVector:
		return &ComplexVector{Data: append([]complex128(nil), x.Data...), Missing: append([]bool(nil), x.Missing...), Attr: cloneAttributeMap(x.Attr)}
	case *CharacterVector:
		return &CharacterVector{Data: append([]string(nil), x.Data...), Missing: append([]bool(nil), x.Missing...), Attr: cloneAttributeMap(x.Attr)}
	case *List:
		return &List{Data: append([]Value(nil), x.Data...), Names: append([]string(nil), x.Names...), Attr: cloneAttributeMap(x.Attr)}
	case *Formula:
		return &Formula{Expr: x.Expr, Env: x.Env, Attr: cloneAttributeMap(x.Attr)}
	default:
		return v
	}
}
