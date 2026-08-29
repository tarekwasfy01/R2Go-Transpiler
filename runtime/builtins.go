package runtime

import (
	"fmt"
	"math"
	"math/cmplx"
	"r2go/syntax"
	"strconv"
)

func (c *Context) builtin(name string, args []syntax.Argument, env *Environment) (Value, error) {
	if name == "eval" {
		if len(args) < 1 {
			return nil, fmt.Errorf("eval requires an expression")
		}
		v, e := c.Eval(args[0].Value, env)
		if e != nil {
			return nil, e
		}
		lang, ok := v.(*Language)
		if !ok {
			return v, nil
		}
		return c.Eval(lang.Expr, env)
	}
	if name == "structure" {
		if len(args) < 1 {
			return nil, fmt.Errorf("structure requires an object")
		}
		v, err := c.Eval(args[0].Value, env)
		if err != nil {
			return nil, err
		}
		out := cloneValue(v)
		for _, arg := range args[1:] {
			if arg.Name == "" {
				return nil, fmt.Errorf("all attributes must be named")
			}
			av, e := c.Eval(arg.Value, env)
			if e != nil {
				return nil, e
			}
			key := arg.Name
			if key == ".Names" {
				key = "names"
			}
			if e = setAttribute(out, key, av); e != nil {
				return nil, e
			}
		}
		return out, nil
	}
	vals := make([]Value, 0, len(args))
	expanded := make([]syntax.Argument, 0, len(args))
	for _, a := range args {
		if symbol, ok := a.Value.(*syntax.Symbol); ok && symbol.Name == "..." {
			_, binding, found := env.Find("...")
			if !found {
				return nil, fmt.Errorf("'...' used in an incorrect context")
			}
			dots, ok := binding.(*DotsBinding)
			if !ok {
				return nil, fmt.Errorf("invalid ... binding")
			}
			forced, e := dots.Force(c)
			if e != nil {
				return nil, e
			}
			list := forced.(*List)
			for i, v := range list.Data {
				name := ""
				if i < len(list.Names) {
					name = list.Names[i]
				}
				vals = append(vals, v)
				expanded = append(expanded, syntax.Argument{Name: name})
			}
			continue
		}
		v, e := c.Eval(a.Value, env)
		if e != nil {
			return nil, e
		}
		vals = append(vals, v)
		expanded = append(expanded, a)
	}
	args = expanded
	switch name {
	case "conditionMessage":
		if len(vals) != 1 {
			return nil, fmt.Errorf("conditionMessage expects one argument")
		}
		condition, ok := vals[0].(*ConditionValue)
		if !ok {
			return nil, fmt.Errorf("not a condition")
		}
		return &CharacterVector{Data: []string{condition.Message}}, nil
	case "length":
		if len(vals) != 1 {
			return nil, fmt.Errorf("length expects one argument")
		}
		return &DoubleVector{Data: []float64{float64(Length(vals[0]))}}, nil
	case "typeof":
		if len(vals) != 1 {
			return nil, fmt.Errorf("typeof expects one argument")
		}
		return &CharacterVector{Data: []string{string(vals[0].Kind())}}, nil
	case "is.null":
		if len(vals) != 1 {
			return nil, fmt.Errorf("is.null expects one argument")
		}
		_, ok := vals[0].(Null)
		return boolValue(ok), nil
	case "is.data.frame":
		if len(vals) != 1 {
			return nil, fmt.Errorf("is.data.frame expects one argument")
		}
		return boolValue(hasClass(vals[0], "data.frame")), nil
	case "is.factor":
		if len(vals) != 1 {
			return nil, fmt.Errorf("is.factor expects one argument")
		}
		return boolValue(hasClass(vals[0], "factor")), nil
	case "is.raw":
		if len(vals) != 1 {
			return nil, fmt.Errorf("is.raw expects one argument")
		}
		_, ok := vals[0].(*RawVector)
		return boolValue(ok), nil
	case "inherits":
		if len(vals) < 2 {
			return nil, fmt.Errorf("inherits expects object and class")
		}
		what, ok := vals[1].(*CharacterVector)
		if !ok {
			return nil, fmt.Errorf("'what' must be a character vector")
		}
		classes := classNames(vals[0])
		out := make([]Logical, len(what.Data))
		for i, w := range what.Data {
			for _, cl := range classes {
				if w == cl {
					out[i] = True
					break
				}
			}
		}
		return &LogicalVector{Data: out}, nil
	case "names":
		if len(vals) != 1 {
			return nil, fmt.Errorf("names expects one argument")
		}
		names := valueNames(vals[0])
		if names == nil {
			return NullValue, nil
		}
		return &CharacterVector{Data: append([]string(nil), names...)}, nil
	case "levels":
		if len(vals) != 1 {
			return nil, fmt.Errorf("levels expects one argument")
		}
		if v, ok := Attributes(vals[0])["levels"]; ok {
			return v, nil
		}
		return NullValue, nil
	case "dim":
		if len(vals) != 1 {
			return nil, fmt.Errorf("dim expects one argument")
		}
		d, ok := dimensions(vals[0])
		if !ok {
			return NullValue, nil
		}
		out := &IntegerVector{Data: make([]int64, len(d))}
		for i, n := range d {
			out.Data[i] = int64(n)
		}
		return out, nil
	case "nrow", "ncol":
		if len(vals) != 1 {
			return nil, fmt.Errorf("%s expects one argument", name)
		}
		if rows, cols, ok := dataFrameDims(vals[0]); ok {
			n := rows
			if name == "ncol" {
				n = cols
			}
			return &IntegerVector{Data: []int64{int64(n)}}, nil
		}
		d, ok := dimensions(vals[0])
		if !ok || len(d) < 2 {
			return NullValue, nil
		}
		i := 0
		if name == "ncol" {
			i = 1
		}
		return &IntegerVector{Data: []int64{int64(d[i])}}, nil
	case "attr":
		if len(vals) != 2 {
			return nil, fmt.Errorf("attr expects object and name")
		}
		n, ok := vals[1].(*CharacterVector)
		if !ok || len(n.Data) != 1 {
			return nil, fmt.Errorf("attribute name must be one string")
		}
		if v, ok := Attributes(vals[0])[n.Data[0]]; ok {
			return v, nil
		}
		return NullValue, nil
	case "attributes":
		if len(vals) != 1 {
			return nil, fmt.Errorf("attributes expects one argument")
		}
		attrs := Attributes(vals[0])
		if len(attrs) == 0 {
			return NullValue, nil
		}
		out := &List{}
		for k, v := range attrs {
			out.Names = append(out.Names, k)
			out.Data = append(out.Data, v)
		}
		return out, nil
	case "class":
		if len(vals) != 1 {
			return nil, fmt.Errorf("class expects one argument")
		}
		if v, ok := Attributes(vals[0])["class"]; ok {
			return v, nil
		}
		return &CharacterVector{Data: []string{defaultClass(vals[0])}}, nil
	case "unclass":
		if len(vals) != 1 {
			return nil, fmt.Errorf("unclass expects one argument")
		}
		out := cloneValue(vals[0])
		_ = setAttribute(out, "class", NullValue)
		return out, nil
	case "is.na":
		if len(vals) != 1 {
			return nil, fmt.Errorf("is.na expects one argument")
		}
		out := make([]Logical, Length(vals[0]))
		switch x := vals[0].(type) {
		case *LogicalVector:
			for i, v := range x.Data {
				if v == NA {
					out[i] = True
				}
			}
		case *IntegerVector:
			for i := range x.Data {
				if i < len(x.Missing) && x.Missing[i] {
					out[i] = True
				}
			}
		case *DoubleVector:
			for i, v := range x.Data {
				if i < len(x.Missing) && x.Missing[i] || math.IsNaN(v) {
					out[i] = True
				}
			}
		case *CharacterVector:
			for i := range x.Data {
				if i < len(x.Missing) && x.Missing[i] {
					out[i] = True
				}
			}
		case *ComplexVector:
			for i := range x.Data {
				if i < len(x.Missing) && x.Missing[i] {
					out[i] = True
				}
			}
		}
		return &LogicalVector{Data: out}, nil
	case "complete.cases":
		return completeCases(vals)
	case "as.logical":
		if len(vals) != 1 {
			return nil, fmt.Errorf("as.logical expects one argument")
		}
		return logical(vals[0])
	case "as.double", "as.numeric":
		if len(vals) != 1 {
			return nil, fmt.Errorf("%s expects one argument", name)
		}
		if name == "as.numeric" {
			return CoerceTo(vals[0], DoubleKind)
		}
		return numbers(vals[0])
	case "as.complex":
		if len(vals) != 1 {
			return nil, fmt.Errorf("as.complex expects one argument")
		}
		return complexNumbers(vals[0])
	case "as.raw":
		if len(vals) != 1 {
			return nil, fmt.Errorf("as.raw expects one argument")
		}
		n, e := numbers(vals[0])
		if e != nil {
			return nil, e
		}
		o := &RawVector{Data: make([]byte, len(n.Data))}
		for i, v := range n.Data {
			if !missingAt(n, i) && v >= 0 && v <= 255 {
				o.Data[i] = byte(v)
			}
		}
		return o, nil
	case "as.vector":
		if len(vals) < 1 || len(vals) > 2 {
			return nil, fmt.Errorf("as.vector expects object and optional mode")
		}
		return cloneValue(vals[0]), nil
	case "numeric", "integer", "logical", "character", "complex", "raw":
		n := 0
		if len(vals) > 0 {
			var err error
			n, err = scalarInt(vals[0])
			if err != nil {
				return nil, fmt.Errorf("invalid 'length' argument")
			}
		}
		if n < 0 {
			return nil, fmt.Errorf("invalid 'length' argument")
		}
		switch name {
		case "numeric":
			return &DoubleVector{Data: make([]float64, n)}, nil
		case "integer":
			return &IntegerVector{Data: make([]int64, n)}, nil
		case "logical":
			return &LogicalVector{Data: make([]Logical, n)}, nil
		case "character":
			return &CharacterVector{Data: make([]string, n)}, nil
		case "complex":
			return &ComplexVector{Data: make([]complex128, n)}, nil
		default:
			return &RawVector{Data: make([]byte, n)}, nil
		}
	case "rawToChar":
		if len(vals) != 1 {
			return nil, fmt.Errorf("rawToChar expects one argument")
		}
		x, ok := vals[0].(*RawVector)
		if !ok {
			return nil, fmt.Errorf("argument must be raw")
		}
		return &CharacterVector{Data: []string{string(x.Data)}}, nil
	case "charToRaw":
		if len(vals) != 1 {
			return nil, fmt.Errorf("charToRaw expects one argument")
		}
		x, ok := vals[0].(*CharacterVector)
		if !ok || len(x.Data) != 1 {
			return nil, fmt.Errorf("argument must be one string")
		}
		return &RawVector{Data: []byte(x.Data[0])}, nil
	case "as.integer":
		if len(vals) != 1 {
			return nil, fmt.Errorf("as.integer expects one argument")
		}
		n, err := numbers(vals[0])
		if err != nil {
			return nil, err
		}
		out := &IntegerVector{Data: make([]int64, len(n.Data)), Missing: append([]bool(nil), n.Missing...)}
		for i, v := range n.Data {
			if math.IsNaN(v) || math.IsInf(v, 0) {
				ensureMissing(&out.Missing, len(n.Data))
				out.Missing[i] = true
			} else {
				out.Data[i] = int64(v)
			}
		}
		return out, nil
	case "as.character":
		if len(vals) != 1 {
			return nil, fmt.Errorf("as.character expects one argument")
		}
		if f, ok := vals[0].(*IntegerVector); ok && hasClass(f, "factor") {
			return factorStrings(f)
		}
		out := &CharacterVector{}
		for _, v := range elements(vals[0]) {
			out.Data = append(out.Data, scalarText(v))
			out.Missing = append(out.Missing, scalarMissing(v))
		}
		return out, nil
	case "Re", "Im", "Mod", "Arg":
		if len(vals) != 1 {
			return nil, fmt.Errorf("%s expects one argument", name)
		}
		z, e := complexNumbers(vals[0])
		if e != nil {
			return nil, e
		}
		o := &DoubleVector{Data: make([]float64, len(z.Data)), Missing: append([]bool(nil), z.Missing...)}
		for i, v := range z.Data {
			switch name {
			case "Re":
				o.Data[i] = real(v)
			case "Im":
				o.Data[i] = imag(v)
			case "Mod":
				o.Data[i] = cmplx.Abs(v)
			case "Arg":
				o.Data[i] = cmplx.Phase(v)
			}
		}
		return o, nil
	case "Conj":
		if len(vals) != 1 {
			return nil, fmt.Errorf("Conj expects one argument")
		}
		// GNU R preserves real/logical/integer inputs; only complex values need
		// conjugation and a complex result container.
		if vals[0].Kind() != ComplexKind {
			return vals[0], nil
		}
		z, e := complexNumbers(vals[0])
		if e != nil {
			return nil, e
		}
		for i := range z.Data {
			z.Data[i] = cmplx.Conj(z.Data[i])
		}
		return z, nil
	case "seq_along":
		if len(vals) != 1 {
			return nil, fmt.Errorf("seq_along expects one argument")
		}
		out := &IntegerVector{Data: make([]int64, Length(vals[0]))}
		for i := range out.Data {
			out.Data[i] = int64(i + 1)
		}
		return out, nil
	case "sum", "prod", "min", "max", "mean":
		return numericSummary(name, vals)
	case "print":
		if len(vals) != 1 {
			return nil, fmt.Errorf("print expects one argument")
		}
		fmt.Fprintln(c.Output, vals[0].String())
		return vals[0], nil
	case "list":
		names := make([]string, len(args))
		for i, a := range args {
			names[i] = a.Name
		}
		return &List{Data: vals, Names: names}, nil
	case "c":
		out, err := combine(vals)
		if err != nil {
			return nil, err
		}
		names := make([]string, len(args))
		for i, argument := range args {
			names[i] = argument.Name
		}
		// Argument tags are names in R's c(a = 1, b = 2). Preserve them on
		// the combined vector rather than dropping them during coercion.
		if hasNames(names) {
			setCombinedNames(out, names)
		}
		return out, nil
	}
	return nil, fmt.Errorf("unknown builtin %s", name)
}

func hasNames(names []string) bool {
	for _, name := range names {
		if name != "" {
			return true
		}
	}
	return false
}
func setCombinedNames(value Value, names []string) {
	out := append([]string(nil), names...)
	switch v := value.(type) {
	case *List:
		v.Names = out
	case *RawVector:
		v.Attr = mapWithNames(v.Attr, out)
	case *LogicalVector:
		v.Attr = mapWithNames(v.Attr, out)
	case *IntegerVector:
		v.Attr = mapWithNames(v.Attr, out)
	case *DoubleVector:
		v.Attr = mapWithNames(v.Attr, out)
	case *ComplexVector:
		v.Attr = mapWithNames(v.Attr, out)
	case *CharacterVector:
		v.Attr = mapWithNames(v.Attr, out)
	}
}
func mapWithNames(attrs map[string]Value, names []string) map[string]Value {
	if attrs == nil {
		attrs = map[string]Value{}
	}
	attrs["names"] = &CharacterVector{Data: names}
	return attrs
}

func defaultClass(v Value) string {
	switch v.Kind() {
	case IntegerKind:
		return "integer"
	case DoubleKind:
		return "numeric"
	case LogicalKind:
		return "logical"
	case CharacterKind:
		return "character"
	case ListKind:
		return "list"
	case ComplexKind:
		return "complex"
	case RawKind:
		return "raw"
	}
	return string(v.Kind())
}
func combine(vals []Value) (Value, error) {
	if len(vals) == 0 {
		return NullValue, nil
	}
	kind, err := CommonKindValues(vals)
	if err != nil {
		return nil, err
	}
	if kind == ListKind {
		out := &List{}
		for _, v := range vals {
			if l, ok := v.(*List); ok {
				out.Data = append(out.Data, l.Data...)
			} else {
				out.Data = append(out.Data, elements(v)...)
			}
		}
		return out, nil
	}
	var pieces []Value
	for _, v := range vals {
		cv, err := CoerceTo(v, kind)
		if err != nil {
			return nil, err
		}
		pieces = append(pieces, cv)
	}
	switch kind {
	case RawKind:
		o := &RawVector{}
		for _, v := range pieces {
			o.Data = append(o.Data, v.(*RawVector).Data...)
		}
		return o, nil
	case LogicalKind:
		o := &LogicalVector{}
		for _, v := range pieces {
			o.Data = append(o.Data, v.(*LogicalVector).Data...)
		}
		return o, nil
	case IntegerKind:
		o := &IntegerVector{}
		for _, v := range pieces {
			x := v.(*IntegerVector)
			o.Data = append(o.Data, x.Data...)
			for i := range x.Data {
				o.Missing = append(o.Missing, i < len(x.Missing) && x.Missing[i])
			}
		}
		return o, nil
	case DoubleKind:
		o := &DoubleVector{}
		for _, v := range pieces {
			x := v.(*DoubleVector)
			o.Data = append(o.Data, x.Data...)
			for i := range x.Data {
				o.Missing = append(o.Missing, i < len(x.Missing) && x.Missing[i])
			}
		}
		return o, nil
	case ComplexKind:
		o := &ComplexVector{}
		for _, v := range pieces {
			x := v.(*ComplexVector)
			o.Data = append(o.Data, x.Data...)
			for i := range x.Data {
				o.Missing = append(o.Missing, i < len(x.Missing) && x.Missing[i])
			}
		}
		return o, nil
	case CharacterKind:
		o := &CharacterVector{}
		for _, v := range pieces {
			x := v.(*CharacterVector)
			o.Data = append(o.Data, x.Data...)
			for i := range x.Data {
				o.Missing = append(o.Missing, i < len(x.Missing) && x.Missing[i])
			}
		}
		return o, nil
	default:
		return nil, fmt.Errorf("unsupported combine kind %s", kind)
	}
}

func ensureMissing(p *[]bool, n int) {
	if len(*p) < n {
		*p = append(*p, make([]bool, n-len(*p))...)
	}
}
func scalarMissing(v Value) bool {
	switch x := v.(type) {
	case *LogicalVector:
		return x.Data[0] == NA
	case *IntegerVector:
		return len(x.Missing) > 0 && x.Missing[0]
	case *DoubleVector:
		return len(x.Missing) > 0 && x.Missing[0]
	case *CharacterVector:
		return len(x.Missing) > 0 && x.Missing[0]
	case *ComplexVector:
		return len(x.Missing) > 0 && x.Missing[0]
	}
	return false
}

// completeCases performs R's row-wise missingness reduction.  A data frame is
// a list of equal-length columns; a matrix contributes each of its columns to
// the same row mask.  This shared kernel also covers ordinary vector arguments.
func completeCases(values []Value) (Value, error) {
	rows := -1
	columns := make([]struct {
		value  Value
		rows   int
		stride int
	}, 0, len(values))
	addColumn := func(value Value, n, stride int) error {
		if rows < 0 {
			rows = n
		} else if rows != n {
			return fmt.Errorf("not all arguments have the same length")
		}
		columns = append(columns, struct {
			value  Value
			rows   int
			stride int
		}{value: value, rows: n, stride: stride})
		return nil
	}
	for _, value := range values {
		if frame, ok := value.(*List); ok && hasClass(frame, "data.frame") {
			for _, column := range frame.Data {
				if err := addColumn(column, Length(column), 1); err != nil {
					return nil, err
				}
			}
			continue
		}
		if dims, ok := dimensions(value); ok && len(dims) >= 2 {
			n := dims[0]
			if err := addColumn(value, n, n); err != nil {
				return nil, err
			}
			continue
		}
		if err := addColumn(value, Length(value), 1); err != nil {
			return nil, err
		}
	}
	if rows < 0 {
		return &LogicalVector{}, nil
	}
	result := &LogicalVector{Data: make([]Logical, rows)}
	for row := 0; row < rows; row++ {
		result.Data[row] = True
		for _, column := range columns {
			for index := row; index < Length(column.value); index += column.stride {
				if valueMissingAt(column.value, index) {
					result.Data[row] = False
					break
				}
				if column.stride == 1 {
					break
				}
			}
			if result.Data[row] == False {
				break
			}
		}
	}
	return result, nil
}

func valueMissingAt(value Value, index int) bool {
	switch vector := value.(type) {
	case *LogicalVector:
		return index >= len(vector.Data) || vector.Data[index] == NA
	case *IntegerVector:
		return index >= len(vector.Data) || index < len(vector.Missing) && vector.Missing[index]
	case *DoubleVector:
		return index >= len(vector.Data) || index < len(vector.Missing) && vector.Missing[index] || index < len(vector.Data) && math.IsNaN(vector.Data[index])
	case *CharacterVector:
		return index >= len(vector.Data) || index < len(vector.Missing) && vector.Missing[index]
	case *ComplexVector:
		return index >= len(vector.Data) || index < len(vector.Missing) && vector.Missing[index]
	case *List:
		return index >= len(vector.Data) || scalarMissing(vector.Data[index])
	default:
		return false
	}
}
func scalarText(v Value) string {
	if scalarMissing(v) {
		return ""
	}
	switch x := v.(type) {
	case *LogicalVector:
		if x.Data[0] == True {
			return "TRUE"
		}
		return "FALSE"
	case *IntegerVector:
		return strconv.FormatInt(x.Data[0], 10)
	case *DoubleVector:
		return strconv.FormatFloat(x.Data[0], 'g', -1, 64)
	case *CharacterVector:
		return x.Data[0]
	case *ComplexVector:
		return x.String()
	case *RawVector:
		return fmt.Sprintf("%02x", x.Data[0])
	}
	return v.String()
}
func numericSummary(name string, vals []Value) (Value, error) {
	var all []float64
	for _, v := range vals {
		n, err := numbers(v)
		if err != nil {
			return nil, err
		}
		for i, x := range n.Data {
			if missingAt(n, i) || math.IsNaN(x) {
				return &DoubleVector{Data: []float64{NAReal()}, Missing: []bool{true}}, nil
			}
			all = append(all, x)
		}
	}
	if len(all) == 0 {
		if name == "sum" {
			return &DoubleVector{Data: []float64{0}}, nil
		}
		if name == "prod" {
			return &DoubleVector{Data: []float64{1}}, nil
		}
		return &DoubleVector{Data: []float64{math.NaN()}}, nil
	}
	r := all[0]
	switch name {
	case "sum", "mean":
		r = 0
		for _, x := range all {
			r += x
		}
		if name == "mean" {
			r /= float64(len(all))
		}
	case "prod":
		r = 1
		for _, x := range all {
			r *= x
		}
	case "min":
		for _, x := range all[1:] {
			if x < r {
				r = x
			}
		}
	case "max":
		for _, x := range all[1:] {
			if x > r {
				r = x
			}
		}
	}
	return &DoubleVector{Data: []float64{r}}, nil
}
