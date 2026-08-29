package runtime

import (
	"fmt"
	"math"
	"r2go/syntax"
)

type unaryScalarFunction func(float64, float64) float64
type binaryScalarFunction func(float64, float64) float64
type bitwiseScalarFunction func(uint32, uint32) (uint32, bool)

func coordinate(entry, offset string) string   { return entry + ":" + offset }
func planCoordinate(plan ExecutionPlan) string { return coordinate(plan.CEntry, plan.Offset) }

var unaryFunctionVector = map[string]unaryScalarFunction{
	coordinate("do_Math2", "10001"): func(x, digits float64) float64 {
		scale := math.Pow(10, digits)
		return math.RoundToEven(x*scale) / scale
	},
	coordinate("do_Math2", "10004"): func(x, digits float64) float64 {
		if x == 0 {
			return 0
		}
		places := digits - 1 - math.Floor(math.Log10(math.Abs(x)))
		scale := math.Pow(10, places)
		return math.RoundToEven(x*scale) / scale
	},
	coordinate("do_log", "10003"): func(x, base float64) float64 {
		if base != 0 {
			return math.Log(x) / math.Log(base)
		}
		return math.Log(x)
	},
	coordinate("do_log1arg", "10010"): func(x, _ float64) float64 { return math.Log10(x) },
	coordinate("do_log1arg", "10002"): func(x, _ float64) float64 { return math.Log2(x) },
	coordinate("do_abs", "6"):         func(x, _ float64) float64 { return math.Abs(x) },
	coordinate("do_trunc", "5"):       func(x, _ float64) float64 { return math.Trunc(x) },
	coordinate("do_math1", "1"):       func(x, _ float64) float64 { return math.Floor(x) }, coordinate("do_math1", "2"): func(x, _ float64) float64 { return math.Ceil(x) },
	coordinate("do_math1", "3"): func(x, _ float64) float64 { return math.Sqrt(x) }, coordinate("do_math1", "4"): func(x, _ float64) float64 {
		if x > 0 {
			return 1
		}
		if x < 0 {
			return -1
		}
		return 0
	},
	coordinate("do_math1", "10"): func(x, _ float64) float64 { return math.Exp(x) }, coordinate("do_math1", "11"): func(x, _ float64) float64 { return math.Expm1(x) }, coordinate("do_math1", "12"): func(x, _ float64) float64 { return math.Log1p(x) },
	coordinate("do_math1", "20"): func(x, _ float64) float64 { return math.Cos(x) }, coordinate("do_math1", "21"): func(x, _ float64) float64 { return math.Sin(x) }, coordinate("do_math1", "22"): func(x, _ float64) float64 { return math.Tan(x) },
	coordinate("do_math1", "23"): func(x, _ float64) float64 { return math.Acos(x) }, coordinate("do_math1", "24"): func(x, _ float64) float64 { return math.Asin(x) }, coordinate("do_math1", "25"): func(x, _ float64) float64 { return math.Atan(x) },
	coordinate("do_math1", "30"): func(x, _ float64) float64 { return math.Cosh(x) }, coordinate("do_math1", "31"): func(x, _ float64) float64 { return math.Sinh(x) }, coordinate("do_math1", "32"): func(x, _ float64) float64 { return math.Tanh(x) },
	coordinate("do_math1", "33"): func(x, _ float64) float64 { return math.Acosh(x) }, coordinate("do_math1", "34"): func(x, _ float64) float64 { return math.Asinh(x) }, coordinate("do_math1", "35"): func(x, _ float64) float64 { return math.Atanh(x) },
	coordinate("do_math1", "40"): func(x, _ float64) float64 { y, _ := math.Lgamma(x); return y }, coordinate("do_math1", "41"): func(x, _ float64) float64 { return math.Gamma(x) }, coordinate("do_math1", "42"): func(x, _ float64) float64 { return digamma(x) }, coordinate("do_math1", "43"): func(x, _ float64) float64 { return trigamma(x) },
	coordinate("do_math1", "47"): func(x, _ float64) float64 { return math.Cos(math.Pi * x) }, coordinate("do_math1", "48"): func(x, _ float64) float64 { return math.Sin(math.Pi * x) }, coordinate("do_math1", "49"): func(x, _ float64) float64 { return math.Tan(math.Pi * x) },
}

var binaryFunctionVector = map[string]binaryScalarFunction{
	coordinate("do_math2", "0"):  math.Atan2,
	coordinate("do_math2", "24"): func(x, order float64) float64 { return besselJ(order, x) },
	coordinate("do_math2", "25"): func(x, order float64) float64 { return besselY(order, x) },
	coordinate("do_math2", "26"): func(x, order float64) float64 { return polyGamma(x, int(order)) },
	coordinate("do_math2", "2"): func(x, y float64) float64 {
		a, _ := math.Lgamma(x)
		b, _ := math.Lgamma(y)
		c, _ := math.Lgamma(x + y)
		return a + b - c
	},
	coordinate("do_math2", "3"): func(x, y float64) float64 {
		a, _ := math.Lgamma(x)
		b, _ := math.Lgamma(y)
		c, _ := math.Lgamma(x + y)
		return math.Exp(a + b - c)
	},
	coordinate("do_math2", "4"): func(x, y float64) float64 {
		a, _ := math.Lgamma(x + 1)
		b, _ := math.Lgamma(y + 1)
		c, _ := math.Lgamma(x - y + 1)
		return a - b - c
	},
	coordinate("do_math2", "5"): func(x, y float64) float64 {
		a, _ := math.Lgamma(x + 1)
		b, _ := math.Lgamma(y + 1)
		c, _ := math.Lgamma(x - y + 1)
		return RoundZero(math.Exp(a - b - c))
	},
}

var operatorFunctionVector = map[string]string{
	coordinate("do_arith", "PLUSOP"): "+", coordinate("do_arith", "MINUSOP"): "-", coordinate("do_arith", "TIMESOP"): "*", coordinate("do_arith", "DIVOP"): "/", coordinate("do_arith", "POWOP"): "^", coordinate("do_arith", "MODOP"): "%%", coordinate("do_arith", "IDIVOP"): "%/%",
	coordinate("do_relop", "EQOP"): "==", coordinate("do_relop", "NEOP"): "!=", coordinate("do_relop", "LTOP"): "<", coordinate("do_relop", "LEOP"): "<=", coordinate("do_relop", "GEOP"): ">=", coordinate("do_relop", "GTOP"): ">",
	coordinate("do_logic", "1"): "&", coordinate("do_logic", "2"): "|", coordinate("do_logic", "3"): "!",
}

var bitwiseFunctionVector = map[string]bitwiseScalarFunction{
	coordinate("do_bitwise", "1"): func(a, b uint32) (uint32, bool) { return a & b, true }, coordinate("do_bitwise", "2"): func(a, _ uint32) (uint32, bool) { return ^a, true }, coordinate("do_bitwise", "3"): func(a, b uint32) (uint32, bool) { return a | b, true }, coordinate("do_bitwise", "4"): func(a, b uint32) (uint32, bool) { return a ^ b, true },
	coordinate("do_bitwise", "5"): func(a, b uint32) (uint32, bool) {
		if b > 31 {
			return 0, false
		}
		return a << b, true
	}, coordinate("do_bitwise", "6"): func(a, b uint32) (uint32, bool) {
		if b > 31 {
			return 0, false
		}
		return a >> b, true
	},
}

func functionVectorAvailable(plan ExecutionPlan) bool {
	switch plan.Opcode {
	case "VECTOR_UNARY":
		return unaryFunctionVector[planCoordinate(plan)] != nil
	case "NUMERIC_BINARY":
		return binaryFunctionVector[planCoordinate(plan)] != nil
	case "OPS_BINARY":
		_, ok := operatorFunctionVector[planCoordinate(plan)]
		return ok
	case "BITWISE":
		return bitwiseFunctionVector[planCoordinate(plan)] != nil
	case "VECTOR_SCAN":
		return plan.CEntry == "do_cum"
	case "VECTOR_REDUCE":
		return plan.CEntry == "do_summary" || plan.CEntry == "do_logic3" || plan.CEntry == "do_pmin" || plan.CEntry == "do_colsum"
	case "TYPE_PREDICATE":
		return plan.CEntry == "do_is"
	}
	return true
}

func (c *Context) executeUnaryVector(plan ExecutionPlan, args []syntax.Argument, env *Environment) (Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("%s expects an argument", plan.Name)
	}
	v, e := c.Eval(args[0].Value, env)
	if e != nil {
		return nil, e
	}
	x, e := numbers(v)
	if e != nil {
		return nil, e
	}
	parameter := 0.0
	if plan.CEntry == "do_log" {
		parameter = math.E
	}
	if plan.CEntry == "do_Math2" && plan.Offset == "10004" {
		parameter = 6
	}
	if len(args) > 1 {
		p, e := c.Eval(args[1].Value, env)
		if e != nil {
			return nil, e
		}
		pn, e := numbers(p)
		if e != nil {
			return nil, e
		}
		parameter = pn.Data[0]
	}
	fn := unaryFunctionVector[planCoordinate(plan)]
	if plan.CEntry == "do_abs" {
		if integers, ok := v.(*IntegerVector); ok {
			out := &IntegerVector{Data: append([]int64(nil), integers.Data...), Missing: append([]bool(nil), integers.Missing...)}
			for i := range out.Data {
				if i < len(out.Missing) && out.Missing[i] {
					continue
				}
				if out.Data[i] < 0 {
					out.Data[i] = -out.Data[i]
				}
			}
			return inheritUnaryAttributes(out, v), nil
		}
	}
	out := &DoubleVector{Data: make([]float64, len(x.Data)), Missing: append([]bool(nil), x.Missing...)}
	for i, n := range x.Data {
		if missingAt(x, i) {
			out.Data[i] = NAReal()
		} else {
			out.Data[i] = fn(n, parameter)
		}
	}
	return inheritUnaryAttributes(out, v), nil
}

func (c *Context) executeBinaryVector(plan ExecutionPlan, args []syntax.Argument, env *Environment) (Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("%s expects two arguments", plan.Name)
	}
	a, e := c.Eval(args[0].Value, env)
	if e != nil {
		return nil, e
	}
	b, e := c.Eval(args[1].Value, env)
	if e != nil {
		return nil, e
	}
	x, e := numbers(a)
	if e != nil {
		return nil, e
	}
	y, e := numbers(b)
	if e != nil {
		return nil, e
	}
	if len(x.Data) == 0 || len(y.Data) == 0 {
		return &DoubleVector{}, nil
	}
	n := max(len(x.Data), len(y.Data))
	out := &DoubleVector{Data: make([]float64, n), Missing: make([]bool, n)}
	fn := binaryFunctionVector[planCoordinate(plan)]
	for i := 0; i < n; i++ {
		xi, yi := i%len(x.Data), i%len(y.Data)
		if missingAt(x, xi) || missingAt(y, yi) {
			out.Missing[i] = true
			out.Data[i] = NAReal()
		} else {
			out.Data[i] = fn(x.Data[xi], y.Data[yi])
		}
	}
	return out, nil
}

func (c *Context) executeBitwiseVector(plan ExecutionPlan, args []syntax.Argument, env *Environment) (Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("%s expects arguments", plan.Name)
	}
	a, e := c.Eval(args[0].Value, env)
	if e != nil {
		return nil, e
	}
	x, e := numbers(a)
	if e != nil {
		return nil, e
	}
	var y *DoubleVector
	if len(args) > 1 {
		b, e := c.Eval(args[1].Value, env)
		if e != nil {
			return nil, e
		}
		y, e = numbers(b)
		if e != nil {
			return nil, e
		}
	} else {
		y = &DoubleVector{Data: []float64{0}}
	}
	if len(x.Data) == 0 || len(y.Data) == 0 {
		return &IntegerVector{}, nil
	}
	n := max(len(x.Data), len(y.Data))
	out := &IntegerVector{Data: make([]int64, n), Missing: make([]bool, n)}
	fn := bitwiseFunctionVector[planCoordinate(plan)]
	for i := 0; i < n; i++ {
		xi, yi := i%len(x.Data), i%len(y.Data)
		if missingAt(x, xi) || missingAt(y, yi) {
			out.Missing[i] = true
			continue
		}
		value, ok := fn(uint32(int64(x.Data[xi])), uint32(int64(y.Data[yi])))
		if !ok {
			out.Missing[i] = true
		} else {
			out.Data[i] = int64(value)
		}
	}
	return out, nil
}

func (c *Context) executeScanVector(plan ExecutionPlan, args []syntax.Argument, env *Environment) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("%s expects one argument", plan.Name)
	}
	v, e := c.Eval(args[0].Value, env)
	if e != nil {
		return nil, e
	}
	x, e := numbers(v)
	if e != nil {
		return nil, e
	}
	out := &DoubleVector{Data: make([]float64, len(x.Data)), Missing: make([]bool, len(x.Data))}
	if plan.Offset == "5" {
		// cumvar is a state-vector operation. Welford's recurrence gives the
		// complete prefix variance in one matrix pass and avoids rescanning every
		// prefix. R has no na.rm argument here, so NA poisons later prefixes.
		count, mean, m2 := 0.0, 0.0, 0.0
		poisoned := false
		for i, n := range x.Data {
			if poisoned || missingAt(x, i) {
				poisoned = true
				out.Missing[i] = true
				out.Data[i] = NAReal()
				continue
			}
			count++
			delta := n - mean
			mean += delta / count
			m2 += delta * (n - mean)
			if count < 2 {
				out.Missing[i] = true
				out.Data[i] = NAReal()
			} else {
				out.Data[i] = m2 / (count - 1)
			}
		}
		return out, nil
	}
	state := 0.0
	if plan.Offset == "2" {
		state = 1
	}
	for i, n := range x.Data {
		if missingAt(x, i) {
			out.Missing[i] = true
			out.Data[i] = NAReal()
			continue
		}
		if i == 0 && (plan.Offset == "3" || plan.Offset == "4") {
			state = n
		} else {
			switch plan.Offset {
			case "1":
				state += n
			case "2":
				state *= n
			case "3":
				if n > state {
					state = n
				}
			case "4":
				if n < state {
					state = n
				}
			}
		}
		out.Data[i] = state
	}
	if (v.Kind() == IntegerKind || v.Kind() == LogicalKind) && (plan.Offset == "1" || plan.Offset == "3" || plan.Offset == "4") {
		integers := &IntegerVector{Data: make([]int64, len(out.Data)), Missing: append([]bool(nil), out.Missing...)}
		for i, value := range out.Data {
			if i < len(integers.Missing) && integers.Missing[i] {
				continue
			}
			integers.Data[i] = int64(value)
		}
		return inheritUnaryAttributes(integers, v), nil
	}
	return out, nil
}

func (c *Context) executeReduceVector(plan ExecutionPlan, args []syntax.Argument, env *Environment) (Value, error) {
	if plan.CEntry == "do_pmin" {
		return c.executeParallelExtrema(plan, args, env)
	}
	if plan.CEntry == "do_colsum" {
		return c.executeMarginReduce(plan, args, env)
	}
	if plan.CEntry == "do_logic3" {
		wantAll := plan.Offset == "1"
		seenNA := false
		naRemove := false
		for _, arg := range args {
			if arg.Name != "na.rm" {
				continue
			}
			v, e := c.Eval(arg.Value, env)
			if e != nil {
				return nil, e
			}
			naRemove, e = IsTrue(v)
			if e != nil {
				return nil, e
			}
		}
		for _, arg := range args {
			if arg.Name == "na.rm" {
				continue
			}
			if arg.Name != "" {
				continue
			}
			v, e := c.Eval(arg.Value, env)
			if e != nil {
				return nil, e
			}
			for _, element := range elements(v) {
				truth, e := IsTrue(element)
				if e != nil {
					if scalarMissing(element) {
						if !naRemove {
							seenNA = true
						}
						continue
					}
					return nil, e
				}
				if wantAll && !truth {
					return boolValue(false), nil
				}
				if !wantAll && truth {
					return boolValue(true), nil
				}
			}
		}
		if seenNA {
			return &LogicalVector{Data: []Logical{NA}}, nil
		}
		return boolValue(wantAll), nil
	}
	naRemove := false
	values := []float64{}
	missing := false
	integerInputs := true
	for _, arg := range args {
		if arg.Name != "na.rm" {
			continue
		}
		v, e := c.Eval(arg.Value, env)
		if e != nil {
			return nil, e
		}
		naRemove, e = IsTrue(v)
		if e != nil {
			return nil, e
		}
	}
	for _, arg := range args {
		if arg.Name == "na.rm" {
			continue
		}
		v, e := c.Eval(arg.Value, env)
		if e != nil {
			return nil, e
		}
		if v.Kind() != IntegerKind && v.Kind() != LogicalKind {
			integerInputs = false
		}
		n, e := numbers(v)
		if e != nil {
			return nil, e
		}
		for i, x := range n.Data {
			if missingAt(n, i) {
				if !naRemove {
					missing = true
				}
				continue
			}
			values = append(values, x)
		}
	}
	if missing {
		if (plan.Offset == "0" || plan.Offset == "2" || plan.Offset == "3") && integerInputs {
			return &IntegerVector{Data: []int64{0}, Missing: []bool{true}}, nil
		}
		return &DoubleVector{Data: []float64{NAReal()}, Missing: []bool{true}}, nil
	}
	result := 0.0
	switch plan.Offset {
	case "0":
		for _, x := range values {
			result += x
		}
		if integerInputs && result >= math.MinInt64 && result <= math.MaxInt64 && result == math.Trunc(result) {
			return &IntegerVector{Data: []int64{int64(result)}}, nil
		}
	case "1":
		if len(values) == 0 {
			return &DoubleVector{Data: []float64{math.NaN()}}, nil
		}
		for _, x := range values {
			result += x
		}
		result /= float64(len(values))
	case "2":
		if len(values) == 0 {
			result = math.Inf(1)
		} else {
			result = values[0]
			for _, x := range values[1:] {
				if x < result {
					result = x
				}
			}
		}
	case "3":
		if len(values) == 0 {
			result = math.Inf(-1)
		} else {
			result = values[0]
			for _, x := range values[1:] {
				if x > result {
					result = x
				}
			}
		}
	case "4":
		result = 1
		for _, x := range values {
			result *= x
		}
	}
	if integerInputs && (plan.Offset == "2" || plan.Offset == "3") {
		return &IntegerVector{Data: []int64{int64(result)}}, nil
	}
	return &DoubleVector{Data: []float64{result}}, nil
}

func (c *Context) executeParallelExtrema(plan ExecutionPlan, args []syntax.Argument, env *Environment) (Value, error) {
	naRemove := false
	vectors := []*DoubleVector{}
	length := 0
	for _, arg := range args {
		if arg.Name == "na.rm" {
			v, e := c.Eval(arg.Value, env)
			if e != nil {
				return nil, e
			}
			naRemove, e = IsTrue(v)
			if e != nil {
				return nil, e
			}
			continue
		}
		v, e := c.Eval(arg.Value, env)
		if e != nil {
			return nil, e
		}
		n, e := numbers(v)
		if e != nil {
			return nil, e
		}
		vectors = append(vectors, n)
		if len(n.Data) > length {
			length = len(n.Data)
		}
	}
	out := &DoubleVector{Data: make([]float64, length), Missing: make([]bool, length)}
	for i := 0; i < length; i++ {
		set := false
		value := 0.0
		for _, v := range vectors {
			if len(v.Data) == 0 {
				continue
			}
			j := i % len(v.Data)
			if missingAt(v, j) {
				if !naRemove {
					out.Missing[i] = true
					out.Data[i] = NAReal()
					break
				}
				continue
			}
			if !set {
				value = v.Data[j]
				set = true
			} else if plan.Offset == "0" && v.Data[j] < value {
				value = v.Data[j]
			} else if plan.Offset == "1" && v.Data[j] > value {
				value = v.Data[j]
			}
		}
		if !out.Missing[i] {
			if set {
				out.Data[i] = value
			} else {
				out.Missing[i] = true
				out.Data[i] = NAReal()
			}
		}
	}
	return out, nil
}

func (c *Context) executeMarginReduce(plan ExecutionPlan, args []syntax.Argument, env *Environment) (Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("margin reduction expects data")
	}
	v, e := c.Eval(args[0].Value, env)
	if e != nil {
		return nil, e
	}
	x, e := numbers(v)
	if e != nil {
		return nil, e
	}
	dims, ok := dimensions(v)
	if !ok || len(dims) != 2 {
		if len(args) < 3 {
			return nil, fmt.Errorf("matrix dimensions are required")
		}
		nr, e := c.Eval(args[1].Value, env)
		if e != nil {
			return nil, e
		}
		nc, e := c.Eval(args[2].Value, env)
		if e != nil {
			return nil, e
		}
		r, er := scalarInt(nr)
		co, ec := scalarInt(nc)
		if er != nil || ec != nil {
			return nil, fmt.Errorf("invalid matrix dimensions")
		}
		dims = []int{r, co}
	}
	naRemove := false
	for _, arg := range args {
		if arg.Name == "na.rm" {
			a, e := c.Eval(arg.Value, env)
			if e != nil {
				return nil, e
			}
			naRemove, e = IsTrue(a)
			if e != nil {
				return nil, e
			}
		}
	}
	rows, cols := dims[0], dims[1]
	byColumn := plan.Offset == "0" || plan.Offset == "1"
	means := plan.Offset == "1" || plan.Offset == "3"
	size := rows
	if !byColumn {
		size = cols
	}
	count := cols
	if byColumn {
		count = rows
	}
	outLen := rows
	if byColumn {
		outLen = cols
	}
	out := &DoubleVector{Data: make([]float64, outLen), Missing: make([]bool, outLen)}
	for outer := 0; outer < outLen; outer++ {
		sum := 0.0
		n := 0
		for inner := 0; inner < count; inner++ {
			row, col := outer, inner
			if byColumn {
				row, col = inner, outer
			}
			index := row + rows*col
			if missingAt(x, index) {
				if !naRemove {
					out.Missing[outer] = true
					out.Data[outer] = NAReal()
					break
				}
				continue
			}
			sum += x.Data[index]
			n++
		}
		if !out.Missing[outer] {
			if means {
				if n == 0 {
					out.Data[outer] = math.NaN()
				} else {
					out.Data[outer] = sum / float64(n)
				}
			} else {
				out.Data[outer] = sum
			}
		}
	}
	_ = size
	return out, nil
}

func (c *Context) executeTypePredicate(plan ExecutionPlan, args []syntax.Argument, env *Environment) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("%s expects one argument", plan.Name)
	}
	v, e := c.Eval(args[0].Value, env)
	if e != nil {
		return nil, e
	}
	result := false
	switch plan.Offset {
	case "NILSXP":
		_, result = v.(Null)
	case "LGLSXP":
		_, result = v.(*LogicalVector)
	case "INTSXP":
		_, result = v.(*IntegerVector)
	case "REALSXP":
		_, result = v.(*DoubleVector)
	case "CPLXSXP":
		_, result = v.(*ComplexVector)
	case "STRSXP":
		_, result = v.(*CharacterVector)
	case "ENVSXP":
		_, result = v.(*EnvironmentValue)
	case "VECSXP":
		_, result = v.(*List)
	case "RAWSXP":
		_, result = v.(*RawVector)
	case "50":
		result = len(Attributes(v)) > 0
	case "100":
		_, result = v.(*IntegerVector)
		if !result {
			_, result = v.(*DoubleVector)
		}
	case "101":
		d, ok := dimensions(v)
		result = ok && len(d) == 2
	case "102":
		_, result = dimensions(v)
	case "200":
		switch v.(type) {
		case *RawVector, *LogicalVector, *IntegerVector, *DoubleVector, *ComplexVector, *CharacterVector:
			result = true
		}
	case "201":
		switch v.(type) {
		case *List, *EnvironmentValue:
			result = true
		}
	case "300":
		_, result = v.(*Language)
	case "301":
		switch v.(type) {
		case *Language, *Formula:
			result = true
		}
	case "302":
		_, result = v.(*Closure)
	}
	return boolValue(result), nil
}
