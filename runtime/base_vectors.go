package runtime

import (
	"fmt"
	"math"
	"r2go/syntax"
	"sort"
	"strings"
)

func (c *Context) vectorBuiltin(name string, args []syntax.Argument, env *Environment) (Value, error) {
	if name == "lapply" || name == "sapply" {
		return c.applyBuiltin(name, args, env)
	}
	vals := make([]Value, len(args))
	named := map[string]Value{}
	for i, a := range args {
		if a.Value == nil {
			return nil, fmt.Errorf("argument %d is empty", i+1)
		}
		v, e := c.Eval(a.Value, env)
		if e != nil {
			return nil, e
		}
		vals[i] = v
		if a.Name != "" {
			named[a.Name] = v
		}
	}
	switch name {
	case "seq", "seq.int":
		return sequence(vals, named)
	case "rep":
		return repeatVector(vals, named)
	case "any", "all":
		return logicalSummary(name, vals, named)
	case "which":
		if len(vals) != 1 {
			return nil, fmt.Errorf("which expects one argument")
		}
		x, e := logical(vals[0])
		if e != nil {
			return nil, e
		}
		o := &IntegerVector{}
		for i, v := range x.Data {
			if v == True {
				o.Data = append(o.Data, int64(i+1))
			}
		}
		return o, nil
	case "unique", "duplicated":
		if len(vals) != 1 {
			return nil, fmt.Errorf("%s expects one argument", name)
		}
		return uniqueOrDuplicated(vals[0], name == "duplicated")
	case "match":
		if len(vals) < 2 {
			return nil, fmt.Errorf("match expects x and table")
		}
		if len(vals) > 2 {
			nomatch, err := scalarInt(vals[2])
			if err == nil {
				return matchValuesNomatch(vals[0], vals[1], int64(nomatch)), nil
			}
		}
		return matchValues(vals[0], vals[1]), nil
	case "paste", "paste0":
		return pasteValues(name, vals, args, named)
	case "table":
		return tableValues(vals, args)
	case "ifelse":
		return ifelseValues(vals)
	case "split":
		return splitValues(vals)
	}
	return nil, fmt.Errorf("unknown vector builtin %s", name)
}

func splitValues(vals []Value) (Value, error) {
	if len(vals) < 2 {
		return nil, fmt.Errorf("split expects x and f")
	}
	x, factor := vals[0], vals[1]
	if Length(factor) == 0 {
		return &List{}, nil
	}
	var labels *CharacterVector
	var err error
	if coded, ok := factor.(*IntegerVector); ok && hasClass(coded, "factor") {
		labels, err = factorStrings(coded)
	} else {
		labels, err = stringsOf(factor)
	}
	if err != nil {
		return nil, err
	}
	groups := map[string][]int{}
	for i := 0; i < Length(x); i++ {
		j := i % len(labels.Data)
		if j < len(labels.Missing) && labels.Missing[j] {
			continue
		}
		groups[labels.Data[j]] = append(groups[labels.Data[j]], i)
	}
	names := make([]string, 0, len(groups))
	for name := range groups {
		names = append(names, name)
	}
	sort.Strings(names)
	out := &List{Names: append([]string(nil), names...)}
	for _, name := range names {
		out.Data = append(out.Data, takePositions(x, groups[name]))
	}
	return out, nil
}

func (c *Context) matrixApplyBuiltin(args []syntax.Argument, env *Environment) (Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("apply expects X, MARGIN and FUN")
	}
	x, err := c.Eval(args[0].Value, env)
	if err != nil {
		return nil, err
	}
	dims, ok := dimensions(x)
	if !ok || len(dims) != 2 {
		return nil, fmt.Errorf("dim(X) must have a positive length")
	}
	marginValue, err := c.Eval(args[1].Value, env)
	if err != nil {
		return nil, err
	}
	margin, err := scalarInt(marginValue)
	if err != nil || (margin != 1 && margin != 2) {
		return nil, fmt.Errorf("invalid MARGIN")
	}
	funValue, err := c.Eval(args[2].Value, env)
	if err != nil {
		return nil, err
	}
	fn, ok := funValue.(*Closure)
	if !ok {
		return nil, fmt.Errorf("FUN must be a function")
	}
	rows, cols := dims[0], dims[1]
	count := rows
	if margin == 2 {
		count = cols
	}
	results := make([]Value, 0, count)
	for outer := 0; outer < count; outer++ {
		var positions []int
		if margin == 1 {
			positions = make([]int, cols)
			for col := 0; col < cols; col++ {
				positions[col] = outer + rows*col
			}
		} else {
			positions = make([]int, rows)
			for row := 0; row < rows; row++ {
				positions[row] = row + rows*outer
			}
		}
		slice := takePositions(x, positions)
		actuals := []ActualArgument{{EagerValue: slice, HasValue: true}}
		for _, argument := range args[3:] {
			actuals = append(actuals, ActualArgument{Argument: argument, Env: env})
		}
		value, err := c.callClosureActual(fn, actuals)
		if err != nil {
			return nil, err
		}
		results = append(results, value)
	}
	for _, value := range results {
		if _, isList := value.(*List); isList || Length(value) != 1 {
			return &List{Data: results}, nil
		}
	}
	return combine(results)
}

func ifelseValues(vals []Value) (Value, error) {
	if len(vals) < 3 {
		return nil, fmt.Errorf("ifelse expects test, yes and no")
	}
	test, err := logical(vals[0])
	if err != nil {
		return nil, err
	}
	if len(test.Data) == 0 {
		return &LogicalVector{}, nil
	}
	if Length(vals[1]) == 0 || Length(vals[2]) == 0 {
		return nil, fmt.Errorf("replacement has length zero")
	}
	kind, err := CommonKindValues([]Value{vals[1], vals[2]})
	if err != nil {
		return nil, err
	}
	yes, err := CoerceTo(vals[1], kind)
	if err != nil {
		return nil, err
	}
	no, err := CoerceTo(vals[2], kind)
	if err != nil {
		return nil, err
	}
	selected := make([]Value, len(test.Data))
	for i, condition := range test.Data {
		switch condition {
		case True:
			selected[i] = elementAt(yes, i%Length(yes))
		case False:
			selected[i] = elementAt(no, i%Length(no))
		default:
			selected[i] = missingScalarForKind(kind)
		}
	}
	return combine(selected)
}

func missingScalarForKind(kind Kind) Value {
	switch kind {
	case LogicalKind:
		return &LogicalVector{Data: []Logical{NA}}
	case IntegerKind:
		return &IntegerVector{Data: []int64{0}, Missing: []bool{true}}
	case DoubleKind:
		return &DoubleVector{Data: []float64{NAReal()}, Missing: []bool{true}}
	case ComplexKind:
		return &ComplexVector{Data: []complex128{0}, Missing: []bool{true}}
	case CharacterKind:
		return &CharacterVector{Data: []string{""}, Missing: []bool{true}}
	default:
		return NullValue
	}
}

func tableValues(vals []Value, args []syntax.Argument) (Value, error) {
	if len(vals) == 0 {
		return &IntegerVector{}, nil
	}
	if len(vals) != 1 || len(args) != 1 || args[0].Name != "" {
		return nil, fmt.Errorf("multi-way table is not implemented")
	}
	counts := map[string]int64{}
	if factor, ok := vals[0].(*IntegerVector); ok && hasClass(factor, "factor") {
		levels, _ := Attributes(factor)["levels"].(*CharacterVector)
		for i, code := range factor.Data {
			if i < len(factor.Missing) && factor.Missing[i] {
				continue
			}
			level := int(code) - 1
			if levels != nil && level >= 0 && level < len(levels.Data) {
				counts[levels.Data[level]]++
			}
		}
		// GNU R retains unused factor levels as zero-count table cells.
		if levels != nil {
			for _, level := range levels.Data {
				if _, exists := counts[level]; !exists {
					counts[level] = 0
				}
			}
		}
	} else {
		for _, item := range elements(vals[0]) {
			if scalarMissing(item) {
				continue
			}
			counts[scalarText(item)]++
		}
	}
	labels := make([]string, 0, len(counts))
	for label := range counts {
		labels = append(labels, label)
	}
	sort.Strings(labels)
	out := &IntegerVector{Data: make([]int64, len(labels))}
	for i, label := range labels {
		out.Data[i] = counts[label]
	}
	if err := setAttribute(out, "dim", &IntegerVector{Data: []int64{int64(len(labels))}}); err != nil {
		return nil, err
	}
	if err := setAttribute(out, "dimnames", &List{Data: []Value{&CharacterVector{Data: labels}}}); err != nil {
		return nil, err
	}
	if err := setAttribute(out, "class", &CharacterVector{Data: []string{"table"}}); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Context) applyBuiltin(name string, args []syntax.Argument, env *Environment) (Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("%s expects X and FUN", name)
	}
	x, e := c.Eval(args[0].Value, env)
	if e != nil {
		return nil, e
	}
	fun, e := c.Eval(args[1].Value, env)
	if e != nil {
		return nil, e
	}
	fn, ok := fun.(*Closure)
	if !ok {
		return nil, fmt.Errorf("FUN must be a function")
	}
	out := &List{Names: append([]string(nil), valueNames(x)...)}
	for _, item := range elements(x) {
		actuals := make([]ActualArgument, 0, len(args)-1)
		actuals = append(actuals, ActualArgument{EagerValue: item, HasValue: true})
		for _, argument := range args[2:] {
			actuals = append(actuals, ActualArgument{Argument: argument, Env: env})
		}
		v, e := c.callClosureActual(fn, actuals)
		if e != nil {
			return nil, e
		}
		out.Data = append(out.Data, v)
	}
	if name == "lapply" {
		return out, nil
	}
	for _, v := range out.Data {
		if _, ok := v.(*List); ok {
			return out, nil
		}
		if Length(v) != 1 {
			return out, nil
		}
	}
	return combine(out.Data)
}

func sequence(vals []Value, named map[string]Value) (Value, error) {
	if along, ok := named["along.with"]; ok {
		n := Length(along)
		o := &IntegerVector{Data: make([]int64, n)}
		for i := range o.Data {
			o.Data[i] = int64(i + 1)
		}
		return o, nil
	}
	var from, to, by float64 = 1, 1, 1
	if len(vals) == 1 {
		n, e := numbers(vals[0])
		if e != nil || len(n.Data) == 0 {
			return nil, fmt.Errorf("invalid seq argument")
		}
		if len(n.Data) > 1 {
			o := &IntegerVector{Data: make([]int64, len(n.Data))}
			for i := range o.Data {
				o.Data[i] = int64(i + 1)
			}
			return o, nil
		}
		to = n.Data[0]
	}
	if len(vals) >= 2 {
		a, e := numbers(vals[0])
		if e != nil {
			return nil, e
		}
		b, e := numbers(vals[1])
		if e != nil {
			return nil, e
		}
		from, to = a.Data[0], b.Data[0]
	}
	if v, ok := named["from"]; ok {
		n, e := numbers(v)
		if e != nil {
			return nil, e
		}
		from = n.Data[0]
	}
	if v, ok := named["to"]; ok {
		n, e := numbers(v)
		if e != nil {
			return nil, e
		}
		to = n.Data[0]
	}
	if v, ok := named["by"]; ok {
		n, e := numbers(v)
		if e != nil {
			return nil, e
		}
		by = n.Data[0]
	}
	if v, ok := named["length.out"]; ok {
		lengthOut, e := scalarInt(v)
		if e != nil || lengthOut < 0 {
			return nil, fmt.Errorf("invalid 'length.out' argument")
		}
		if lengthOut == 0 {
			return &DoubleVector{}, nil
		}
		o := &DoubleVector{Data: make([]float64, lengthOut)}
		if lengthOut == 1 {
			o.Data[0] = from
			return o, nil
		}
		step := (to - from) / float64(lengthOut-1)
		for i := range o.Data {
			o.Data[i] = from + float64(i)*step
		}
		return o, nil
	}
	if by == 0 {
		return nil, fmt.Errorf("invalid '(to - from)/by'")
	}
	count := int(math.Floor((to-from)/by)) + 1
	if count < 0 {
		count = 0
	}
	o := &DoubleVector{Data: make([]float64, count)}
	for i := range o.Data {
		o.Data[i] = from + float64(i)*by
	}
	return o, nil
}
func repeatVector(vals []Value, named map[string]Value) (Value, error) {
	if len(vals) < 1 {
		return nil, fmt.Errorf("rep expects x")
	}
	times, each := 1, 1
	if len(vals) > 1 {
		times, _ = scalarInt(vals[1])
	}
	if v, ok := named["times"]; ok {
		times, _ = scalarInt(v)
	}
	if v, ok := named["each"]; ok {
		each, _ = scalarInt(v)
	}
	if times < 0 || each < 0 {
		return nil, fmt.Errorf("invalid 'times' argument")
	}
	src := elements(vals[0])
	positions := make([]int, 0, len(src)*times*each)
	for t := 0; t < times; t++ {
		for i := range src {
			for j := 0; j < each; j++ {
				positions = append(positions, i)
			}
		}
	}
	return takePositions(vals[0], positions), nil
}
func logicalSummary(name string, vals []Value, named map[string]Value) (Value, error) {
	naRm := false
	if v, ok := named["na.rm"]; ok {
		b, e := IsTrue(v)
		if e == nil {
			naRm = b
		}
	}
	sawNA := false
	for i, v := range vals {
		if i < len(vals) && argsNameValue(named, v, "na.rm") {
			continue
		}
		x, e := logical(v)
		if e != nil {
			return nil, e
		}
		for _, b := range x.Data {
			if b == NA {
				sawNA = true
				continue
			}
			if name == "any" && b == True {
				return boolValue(true), nil
			}
			if name == "all" && b == False {
				return boolValue(false), nil
			}
		}
	}
	if sawNA && !naRm {
		return &LogicalVector{Data: []Logical{NA}}, nil
	}
	return boolValue(name == "all"), nil
}
func argsNameValue(named map[string]Value, v Value, name string) bool {
	x, ok := named[name]
	return ok && x == v
}
func uniqueOrDuplicated(v Value, duplicates bool) (Value, error) {
	seen := map[string]bool{}
	pos := []int{}
	flags := make([]Logical, Length(v))
	for i, e := range elements(v) {
		k := MatchKey(e)
		dup := seen[k]
		if dup {
			flags[i] = True
		} else {
			seen[k] = true
			pos = append(pos, i)
		}
	}
	if duplicates {
		return &LogicalVector{Data: flags}, nil
	}
	return takePositions(v, pos), nil
}

func matchValues(x, table Value) Value {
	return matchValuesInternal(x, table, 0, true)
}

func matchValuesNomatch(x, table Value, nomatch int64) Value {
	return matchValuesInternal(x, table, nomatch, false)
}

func matchValuesInternal(x, table Value, nomatch int64, missingNomatch bool) Value {
	xx, tt, err := coercePairForMatching(x, table)
	if err == nil {
		x, table = xx, tt
	}
	lookup := map[string]int{}
	for i, e := range elements(table) {
		k := MatchKey(e)
		if _, ok := lookup[k]; !ok {
			lookup[k] = i + 1
		}
	}
	o := &IntegerVector{Data: make([]int64, Length(x)), Missing: make([]bool, Length(x))}
	for i, e := range elements(x) {
		if p, ok := lookup[MatchKey(e)]; ok {
			o.Data[i] = int64(p)
		} else {
			o.Data[i] = nomatch
			o.Missing[i] = missingNomatch
		}
	}
	return o
}

func pasteValues(name string, vals []Value, args []syntax.Argument, named map[string]Value) (Value, error) {
	sep := " "
	if name == "paste0" {
		sep = ""
	}
	if v, ok := named["sep"]; ok {
		sep = scalarText(elements(v)[0])
	}
	collapse := ""
	hasCollapse := false
	if v, ok := named["collapse"]; ok {
		if _, null := v.(Null); !null {
			collapse = scalarText(elements(v)[0])
			hasCollapse = true
		}
	}
	data := []Value{}
	for i, v := range vals {
		if args[i].Name == "sep" || args[i].Name == "collapse" {
			continue
		}
		data = append(data, v)
	}
	n := 0
	for _, v := range data {
		if Length(v) > n {
			n = Length(v)
		}
	}
	out := make([]string, n)
	for i := 0; i < n; i++ {
		parts := make([]string, len(data))
		for j, v := range data {
			es := elements(v)
			parts[j] = scalarText(es[i%len(es)])
		}
		out[i] = strings.Join(parts, sep)
	}
	if hasCollapse {
		return &CharacterVector{Data: []string{strings.Join(out, collapse)}}, nil
	}
	return &CharacterVector{Data: out}, nil
}

func (c *Context) mathBuiltin(name string, args []syntax.Argument, env *Environment) (Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("%s expects an argument", name)
	}
	v, e := c.Eval(args[0].Value, env)
	if e != nil {
		return nil, e
	}
	x, e := numbers(v)
	if e != nil {
		return nil, e
	}
	// GNU R's abs is type-stable for integer inputs, unlike the other Math
	// operations which produce doubles. Apply this before numeric coercion has
	// erased the original container kind.
	if name == "abs" {
		if integers, ok := v.(*IntegerVector); ok {
			o := &IntegerVector{Data: append([]int64(nil), integers.Data...), Missing: append([]bool(nil), integers.Missing...)}
			for i := range o.Data {
				if i < len(o.Missing) && o.Missing[i] {
					continue
				}
				if o.Data[i] < 0 {
					o.Data[i] = -o.Data[i]
				}
			}
			return inheritUnaryAttributes(o, v), nil
		}
	}
	o := &DoubleVector{Data: make([]float64, len(x.Data)), Missing: append([]bool(nil), x.Missing...)}
	digits := 0.0
	if len(args) > 1 {
		d, e := c.Eval(args[1].Value, env)
		if e != nil {
			return nil, e
		}
		dn, e := numbers(d)
		if e != nil {
			return nil, e
		}
		digits = dn.Data[0]
	}
	for i, n := range x.Data {
		if missingAt(x, i) {
			o.Data[i] = NAReal()
			continue
		}
		switch name {
		case "abs":
			o.Data[i] = math.Abs(n)
		case "sqrt":
			o.Data[i] = math.Sqrt(n)
		case "exp":
			o.Data[i] = math.Exp(n)
		case "expm1":
			o.Data[i] = math.Expm1(n)
		case "log":
			o.Data[i] = math.Log(n)
		case "log1p":
			o.Data[i] = math.Log1p(n)
		case "log10":
			o.Data[i] = math.Log10(n)
		case "sin":
			o.Data[i] = math.Sin(n)
		case "cos":
			o.Data[i] = math.Cos(n)
		case "tan":
			o.Data[i] = math.Tan(n)
		case "acos":
			o.Data[i] = math.Acos(n)
		case "asin":
			o.Data[i] = math.Asin(n)
		case "atan":
			o.Data[i] = math.Atan(n)
		case "cosh":
			o.Data[i] = math.Cosh(n)
		case "sinh":
			o.Data[i] = math.Sinh(n)
		case "tanh":
			o.Data[i] = math.Tanh(n)
		case "acosh":
			o.Data[i] = math.Acosh(n)
		case "asinh":
			o.Data[i] = math.Asinh(n)
		case "atanh":
			o.Data[i] = math.Atanh(n)
		case "gamma":
			o.Data[i] = math.Gamma(n)
		case "lgamma":
			o.Data[i], _ = math.Lgamma(n)
		case "digamma":
			o.Data[i] = digamma(n)
		case "trigamma":
			o.Data[i] = trigamma(n)
		case "cospi":
			o.Data[i] = math.Cos(math.Pi * n)
		case "sinpi":
			o.Data[i] = math.Sin(math.Pi * n)
		case "tanpi":
			o.Data[i] = math.Tan(math.Pi * n)
		case "sign":
			if n > 0 {
				o.Data[i] = 1
			} else if n < 0 {
				o.Data[i] = -1
			}
		case "floor":
			o.Data[i] = math.Floor(n)
		case "ceiling":
			o.Data[i] = math.Ceil(n)
		case "trunc":
			o.Data[i] = math.Trunc(n)
		case "round":
			o.Data[i] = RoundZero(n)
		case "signif":
			if n == 0 {
				o.Data[i] = 0
			} else {
				places := digits - 1 - math.Floor(math.Log10(math.Abs(n)))
				scale := math.Pow(10, places)
				o.Data[i] = RoundZero(n*scale) / scale
			}
		}
	}
	return inheritUnaryAttributes(o, v), nil
}
func digamma(x float64) float64 {
	result := 0.0
	for x < 8 {
		result -= 1 / x
		x++
	}
	inv := 1 / x
	inv2 := inv * inv
	return result + math.Log(x) - .5*inv - inv2*(1.0/12-inv2*(1.0/120-inv2/252))
}
func trigamma(x float64) float64 {
	result := 0.0
	for x < 8 {
		result += 1 / (x * x)
		x++
	}
	inv := 1 / x
	inv2 := inv * inv
	return result + inv + inv2/2 + inv2*inv/6 - inv2*inv2*inv/30 + inv2*inv2*inv2*inv/42
}
