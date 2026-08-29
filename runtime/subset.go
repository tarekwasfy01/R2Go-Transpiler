package runtime

import (
	"fmt"
	"math"

	"r2go/syntax"
)

func (c *Context) subset(op string, args []syntax.Argument, env *Environment) (Value, error) {
	if len(args) > 2 {
		object, err := c.Eval(args[0].Value, env)
		if err != nil {
			return nil, err
		}
		if hasClass(object, "data.frame") {
			return c.dataFrameSubset(args, env)
		}
	}
	if len(args) > 2 {
		return c.arraySubset(op, args, env)
	}
	if len(args) != 2 || args[1].Value == nil {
		return nil, fmt.Errorf("only one-dimensional non-missing subscripts are implemented")
	}
	object, err := c.Eval(args[0].Value, env)
	if err != nil {
		return nil, err
	}
	index, err := c.Eval(args[1].Value, env)
	if err != nil {
		return nil, err
	}
	positions, err := subsetPositions(object, index)
	if err != nil {
		return nil, err
	}
	if op == "[[" {
		if len(positions) != 1 || positions[0] < 0 || positions[0] >= Length(object) {
			return nil, fmt.Errorf("subscript out of bounds")
		}
		return elementAt(object, positions[0]), nil
	}
	return takePositions(object, positions), nil
}

func (c *Context) dollar(args []syntax.Argument, env *Environment) (Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("$ expects object and name")
	}
	object, err := c.Eval(args[0].Value, env)
	if err != nil {
		return nil, err
	}
	name, ok := args[1].Value.(*syntax.Symbol)
	if !ok {
		return nil, fmt.Errorf("invalid $ name")
	}
	if environment, ok := object.(*Environment); ok {
		binding, found := environment.Local(name.Name)
		if !found {
			return NullValue, nil
		}
		return binding.Force(c)
	}
	list, ok := object.(*List)
	if !ok {
		return nil, fmt.Errorf("$ operator is invalid for atomic vectors")
	}
	for i, candidate := range list.Names {
		if candidate == name.Name {
			return list.Data[i], nil
		}
	}
	return NullValue, nil
}

func subsetPositions(object Value, index Value) ([]int, error) {
	n := Length(object)
	switch x := index.(type) {
	case *IntegerVector:
		d := &DoubleVector{Data: make([]float64, len(x.Data)), Missing: append([]bool(nil), x.Missing...)}
		for i, v := range x.Data {
			d.Data[i] = float64(v)
		}
		return subsetPositions(object, d)
	case *DoubleVector:
		positive, negative := false, false
		for i, raw := range x.Data {
			if i < len(x.Missing) && x.Missing[i] {
				continue
			}
			if raw > 0 {
				positive = true
			}
			if raw < 0 {
				negative = true
			}
		}
		if positive && negative {
			return nil, fmt.Errorf("only 0's may be mixed with negative subscripts")
		}
		if negative {
			excluded := map[int]bool{}
			for _, raw := range x.Data {
				if raw < 0 {
					excluded[int(-raw)-1] = true
				}
			}
			out := make([]int, 0, n)
			for i := 0; i < n; i++ {
				if !excluded[i] {
					out = append(out, i)
				}
			}
			return out, nil
		}
		out := make([]int, 0, len(x.Data))
		for i, raw := range x.Data {
			if i < len(x.Missing) && x.Missing[i] {
				out = append(out, -1)
				continue
			}
			if raw == 0 {
				continue
			}
			if raw != math.Trunc(raw) {
				return nil, fmt.Errorf("invalid fractional subscript")
			}
			out = append(out, int(raw)-1)
		}
		return out, nil
	case *LogicalVector:
		if len(x.Data) == 0 {
			return []int{}, nil
		}
		limit := n
		if len(x.Data) > limit {
			limit = len(x.Data)
		}
		out := []int{}
		for i := 0; i < limit; i++ {
			switch x.Data[i%len(x.Data)] {
			case True:
				out = append(out, i)
			case NA:
				out = append(out, -1)
			}
		}
		return out, nil
	case *CharacterVector:
		names := valueNames(object)
		out := make([]int, len(x.Data))
		for i, want := range x.Data {
			out[i] = -1
			for j, got := range names {
				if got == want {
					out[i] = j
					break
				}
			}
		}
		return out, nil
	default:
		return nil, fmt.Errorf("invalid subscript type %s", index.Kind())
	}
}

func valueNames(v Value) []string {
	switch x := v.(type) {
	case *RawVector:
		if n, ok := x.Attr["names"].(*CharacterVector); ok {
			return n.Data
		}
	case *List:
		return x.Names
	case *DoubleVector:
		if n, ok := x.Attr["names"].(*CharacterVector); ok {
			return n.Data
		}
	case *IntegerVector:
		if n, ok := x.Attr["names"].(*CharacterVector); ok {
			return n.Data
		}
	case *LogicalVector:
		if n, ok := x.Attr["names"].(*CharacterVector); ok {
			return n.Data
		}
	case *CharacterVector:
		if n, ok := x.Attr["names"].(*CharacterVector); ok {
			return n.Data
		}
	case *ComplexVector:
		if n, ok := x.Attr["names"].(*CharacterVector); ok {
			return n.Data
		}
	}
	// GNU R exposes the sole dimnames component of one-dimensional arrays
	// (notably table objects) through names(). Matrices intentionally do not use
	// this projection.
	attrs := Attributes(v)
	if dim, ok := attrs["dim"].(*IntegerVector); ok && len(dim.Data) == 1 {
		if dimnames, ok := attrs["dimnames"].(*List); ok && len(dimnames.Data) == 1 {
			if names, ok := dimnames.Data[0].(*CharacterVector); ok {
				return names.Data
			}
		}
	}
	return nil
}

func elementAt(v Value, i int) Value {
	switch x := v.(type) {
	case *RawVector:
		return &RawVector{Data: []byte{x.Data[i]}}
	case *IntegerVector:
		return &IntegerVector{Data: []int64{x.Data[i]}, Missing: []bool{i < len(x.Missing) && x.Missing[i]}}
	case *DoubleVector:
		return &DoubleVector{Data: []float64{x.Data[i]}, Missing: []bool{i < len(x.Missing) && x.Missing[i]}}
	case *LogicalVector:
		return &LogicalVector{Data: []Logical{x.Data[i]}}
	case *CharacterVector:
		return &CharacterVector{Data: []string{x.Data[i]}, Missing: []bool{i < len(x.Missing) && x.Missing[i]}}
	case *ComplexVector:
		return &ComplexVector{Data: []complex128{x.Data[i]}, Missing: []bool{i < len(x.Missing) && x.Missing[i]}}
	case *List:
		return x.Data[i]
	}
	return NullValue
}

func replacePositions(v Value, positions []int, replacement Value) (Value, error) {
	if len(positions) == 0 {
		return v, nil
	}
	if Length(replacement) == 0 {
		return nil, fmt.Errorf("replacement has length zero")
	}
	for _, p := range positions {
		if p < 0 {
			return nil, fmt.Errorf("NAs are not allowed in subscripted assignments")
		}
	}

	// GNU R promotes the target vector according to the combined target/RHS
	// type before assignment. Centralizing that decision avoids per-operator
	// narrowing bugs and shares the same matrix as c()/matching/coercion.
	targetKind, err := CommonKind(v.Kind(), replacement.Kind())
	if err != nil {
		return nil, err
	}
	if targetKind != v.Kind() {
		attributes := Attributes(v)
		v, err = CoerceTo(v, targetKind)
		if err != nil {
			return nil, err
		}
		for name, attribute := range attributes {
			if err := setAttribute(v, name, cloneValue(attribute)); err != nil {
				return nil, err
			}
		}
	}
	replacement, err = CoerceTo(replacement, targetKind)
	if err != nil {
		return nil, err
	}
	values := elements(replacement)

	for i, p := range positions {
		r := values[i%len(values)]
		switch x := v.(type) {
		case *LogicalVector:
			growLogical(x, p+1)
			x.Data[p] = r.(*LogicalVector).Data[0]
		case *IntegerVector:
			y := r.(*IntegerVector)
			growInteger(x, p+1)
			x.Data[p] = y.Data[0]
			x.Missing[p] = len(y.Missing) > 0 && y.Missing[0]
		case *DoubleVector:
			y := r.(*DoubleVector)
			growDouble(x, p+1)
			x.Data[p] = y.Data[0]
			x.Missing[p] = len(y.Missing) > 0 && y.Missing[0]
		case *CharacterVector:
			y := r.(*CharacterVector)
			growCharacter(x, p+1)
			x.Data[p] = y.Data[0]
			x.Missing[p] = len(y.Missing) > 0 && y.Missing[0]
		case *ComplexVector:
			y := r.(*ComplexVector)
			for len(x.Data) < p+1 {
				x.Data = append(x.Data, 0)
				x.Missing = append(x.Missing, true)
			}
			ensureMissing(&x.Missing, len(x.Data))
			x.Data[p] = y.Data[0]
			x.Missing[p] = len(y.Missing) > 0 && y.Missing[0]
		case *RawVector:
			y := r.(*RawVector)
			for len(x.Data) < p+1 {
				x.Data = append(x.Data, 0)
			}
			x.Data[p] = y.Data[0]
		case *List:
			for len(x.Data) < p+1 {
				x.Data = append(x.Data, NullValue)
				x.Names = append(x.Names, "")
			}
			x.Data[p] = r
		default:
			return nil, fmt.Errorf("object of type %s is not subsettable", v.Kind())
		}
	}
	return v, nil
}

func growLogical(x *LogicalVector, n int) {
	for len(x.Data) < n {
		x.Data = append(x.Data, NA)
	}
}
func growInteger(x *IntegerVector, n int) {
	for len(x.Data) < n {
		x.Data = append(x.Data, 0)
		x.Missing = append(x.Missing, true)
	}
	ensureMissing(&x.Missing, len(x.Data))
}
func growDouble(x *DoubleVector, n int) {
	for len(x.Data) < n {
		x.Data = append(x.Data, NAReal())
		x.Missing = append(x.Missing, true)
	}
	ensureMissing(&x.Missing, len(x.Data))
}
func growCharacter(x *CharacterVector, n int) {
	for len(x.Data) < n {
		x.Data = append(x.Data, "")
		x.Missing = append(x.Missing, true)
	}
	ensureMissing(&x.Missing, len(x.Data))
}

func replaceDollar(v Value, name string, replacement Value) error {
	if environment, ok := v.(*Environment); ok {
		environment.Set(name, replacement)
		return nil
	}
	x, ok := v.(*List)
	if !ok {
		return fmt.Errorf("$ operator is invalid for atomic vectors")
	}
	for i, n := range x.Names {
		if n == name {
			x.Data[i] = replacement
			return nil
		}
	}
	x.Data = append(x.Data, replacement)
	x.Names = append(x.Names, name)
	return nil
}

func replaceNames(v Value, replacement Value) error {
	if _, null := replacement.(Null); null {
		if x, ok := v.(*List); ok {
			x.Names = nil
		}
		return setAttribute(v, "names", NullValue)
	}
	n, ok := replacement.(*CharacterVector)
	if !ok {
		return fmt.Errorf("names must be a character vector")
	}
	if len(n.Data) > Length(v) {
		return fmt.Errorf("'names' attribute must be the same length as the vector")
	}
	names := append([]string(nil), n.Data...)
	for len(names) < Length(v) {
		names = append(names, "")
	}
	if x, ok := v.(*List); ok {
		x.Names = names
	}
	return setAttribute(v, "names", &CharacterVector{Data: names})
}

func takePositions(v Value, positions []int) Value {
	switch x := v.(type) {
	case *RawVector:
		o := &RawVector{Data: make([]byte, len(positions))}
		for i, p := range positions {
			if p >= 0 && p < len(x.Data) {
				o.Data[i] = x.Data[p]
			}
		}
		return o
	case *IntegerVector:
		o := &IntegerVector{Data: make([]int64, len(positions)), Missing: make([]bool, len(positions))}
		for i, p := range positions {
			if p < 0 || p >= len(x.Data) {
				o.Missing[i] = true
			} else {
				o.Data[i] = x.Data[p]
				o.Missing[i] = p < len(x.Missing) && x.Missing[p]
			}
		}
		if hasClass(x, "factor") {
			o.Attr = map[string]Value{"class": x.Attr["class"], "levels": x.Attr["levels"]}
		}
		copySubsetNames(o, x, positions)
		return o
	case *DoubleVector:
		o := &DoubleVector{Data: make([]float64, len(positions)), Missing: make([]bool, len(positions))}
		for i, p := range positions {
			if p < 0 || p >= len(x.Data) {
				o.Missing[i] = true
				o.Data[i] = NAReal()
			} else {
				o.Data[i] = x.Data[p]
				o.Missing[i] = p < len(x.Missing) && x.Missing[p]
			}
		}
		copySubsetNames(o, x, positions)
		return o
	case *LogicalVector:
		o := &LogicalVector{Data: make([]Logical, len(positions))}
		for i, p := range positions {
			if p < 0 || p >= len(x.Data) {
				o.Data[i] = NA
			} else {
				o.Data[i] = x.Data[p]
			}
		}
		copySubsetNames(o, x, positions)
		return o
	case *CharacterVector:
		o := &CharacterVector{Data: make([]string, len(positions)), Missing: make([]bool, len(positions))}
		for i, p := range positions {
			if p < 0 || p >= len(x.Data) {
				o.Missing[i] = true
			} else {
				o.Data[i] = x.Data[p]
				o.Missing[i] = p < len(x.Missing) && x.Missing[p]
			}
		}
		copySubsetNames(o, x, positions)
		return o
	case *ComplexVector:
		o := &ComplexVector{Data: make([]complex128, len(positions)), Missing: make([]bool, len(positions))}
		for i, p := range positions {
			if p < 0 || p >= len(x.Data) {
				o.Missing[i] = true
			} else {
				o.Data[i] = x.Data[p]
				o.Missing[i] = p < len(x.Missing) && x.Missing[p]
			}
		}
		copySubsetNames(o, x, positions)
		return o
	case *List:
		o := &List{Data: make([]Value, len(positions)), Names: make([]string, len(positions))}
		for i, p := range positions {
			if p < 0 || p >= len(x.Data) {
				o.Data[i] = NullValue
			} else {
				o.Data[i] = x.Data[p]
				if p < len(x.Names) {
					o.Names[i] = x.Names[p]
				}
			}
		}
		return o
	}
	return NullValue
}

func copySubsetNames(out, source Value, positions []int) {
	names := valueNames(source)
	if len(names) == 0 {
		return
	}
	selected := make([]string, len(positions))
	for i, position := range positions {
		if position >= 0 && position < len(names) {
			selected[i] = names[position]
		}
	}
	_ = setAttribute(out, "names", &CharacterVector{Data: selected})
}
