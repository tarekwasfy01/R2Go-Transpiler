package runtime

import (
	"fmt"
	"math"
	"math/cmplx"
	"r2go/syntax"
)

func (c *Context) operator(op string, args []syntax.Argument, env *Environment) (Value, error) {
	// R's scalar logical operators are genuinely lazy.
	if op == "&&" || op == "||" {
		if len(args) != 2 {
			return nil, fmt.Errorf("operator %s expects two operands", op)
		}
		left, e := c.Eval(args[0].Value, env)
		if e != nil {
			return nil, e
		}
		a, e := IsTrue(left)
		if e != nil {
			return nil, e
		}
		if op == "&&" && !a {
			return boolValue(false), nil
		}
		if op == "||" && a {
			return boolValue(true), nil
		}
		right, e := c.Eval(args[1].Value, env)
		if e != nil {
			return nil, e
		}
		b, e := IsTrue(right)
		if e != nil {
			return nil, e
		}
		if op == "&&" {
			return boolValue(a && b), nil
		}
		return boolValue(a || b), nil
	}
	vals := make([]Value, len(args))
	for i, a := range args {
		v, e := c.Eval(a.Value, env)
		if e != nil {
			return nil, e
		}
		vals[i] = v
	}
	if op == "!" {
		if len(vals) != 1 {
			return nil, fmt.Errorf("invalid unary !")
		}
		v, err := logical(vals[0])
		if err != nil {
			return nil, err
		}
		for i := range v.Data {
			if v.Data[i] != NA {
				if v.Data[i] == True {
					v.Data[i] = False
				} else {
					v.Data[i] = True
				}
			}
		}
		return v, nil
	}
	if len(vals) == 1 && (op == "+" || op == "-") && vals[0].Kind() == ComplexKind {
		x, _ := complexNumbers(vals[0])
		if op == "-" {
			for i := range x.Data {
				x.Data[i] = -x.Data[i]
			}
		}
		return x, nil
	}
	if len(vals) == 1 && (op == "+" || op == "-") {
		v, err := numbers(vals[0])
		if err != nil {
			return nil, err
		}
		if op == "-" {
			for i := range v.Data {
				v.Data[i] = -v.Data[i]
			}
		}
		return v, nil
	}
	if len(vals) != 2 {
		return nil, fmt.Errorf("operator %s expects two operands", op)
	}
	if vals[0].Kind() == ComplexKind || vals[1].Kind() == ComplexKind {
		return c.complexOperator(op, vals)
	}
	if op == "%in%" {
		matched := matchValues(vals[0], vals[1]).(*IntegerVector)
		out := &LogicalVector{Data: make([]Logical, len(matched.Data))}
		for i := range out.Data {
			if i >= len(matched.Missing) || !matched.Missing[i] {
				out.Data[i] = True
			}
		}
		return out, nil
	}
	a, err := numbers(vals[0])
	if err != nil {
		return nil, err
	}
	b, err := numbers(vals[1])
	if err != nil {
		return nil, err
	}
	n := max(len(a.Data), len(b.Data))
	if n == 0 {
		return &DoubleVector{}, nil
	}
	c.warnFractionalRecycling(len(a.Data), len(b.Data))
	if op == ":" {
		if len(a.Data) != 1 || len(b.Data) != 1 {
			return nil, fmt.Errorf(": requires scalar operands")
		}
		from, to := a.Data[0], b.Data[0]
		step := 1.0
		if from > to {
			step = -1
		}
		n = int(math.Abs(to-from)) + 1
		o := make([]float64, n)
		for i := range o {
			o[i] = from + float64(i)*step
		}
		return &DoubleVector{Data: o}, nil
	}
	cmp := op == "==" || op == "!=" || op == "<" || op == "<=" || op == ">" || op == ">=" || op == "&" || op == "|"
	if cmp {
		o := make([]Logical, n)
		for i := 0; i < n; i++ {
			if missingAt(a, i%len(a.Data)) || missingAt(b, i%len(b.Data)) || math.IsNaN(a.Data[i%len(a.Data)]) || math.IsNaN(b.Data[i%len(b.Data)]) {
				o[i] = NA
				continue
			}
			x, y := a.Data[i%len(a.Data)], b.Data[i%len(b.Data)]
			z := false
			switch op {
			case "==":
				z = x == y
			case "!=":
				z = x != y
			case "<":
				z = x < y
			case "<=":
				z = x <= y
			case ">":
				z = x > y
			case ">=":
				z = x >= y
			case "&":
				z = x != 0 && y != 0
			case "|":
				z = x != 0 || y != 0
			}
			if z {
				o[i] = True
			}
		}
		return &LogicalVector{Data: o}, nil
	}
	o := make([]float64, n)
	missing := make([]bool, n)
	for i := range o {
		if missingAt(a, i%len(a.Data)) || missingAt(b, i%len(b.Data)) {
			missing[i] = true
			o[i] = NAReal()
			continue
		}
		x, y := a.Data[i%len(a.Data)], b.Data[i%len(b.Data)]
		switch op {
		case "+":
			o[i] = x + y
		case "-":
			o[i] = x - y
		case "*":
			o[i] = x * y
		case "/":
			o[i] = x / y
		case "^":
			o[i] = math.Pow(x, y)
		case "%%":
			o[i] = x - y*math.Floor(x/y)
		case "%/%":
			o[i] = math.Floor(x / y)
		default:
			return nil, fmt.Errorf("unsupported operator %s", op)
		}
	}
	return &DoubleVector{Data: o, Missing: missing}, nil
}

func (c *Context) warnFractionalRecycling(a, b int) {
	shorter, longer := min(a, b), max(a, b)
	if shorter == 0 || longer%shorter == 0 {
		return
	}
	c.Warnings = append(c.Warnings, &ConditionValue{
		Classes: []string{"simpleWarning", "warning", "condition"},
		Message: "longer object length is not a multiple of shorter object length",
	})
}
func complexNumbers(v Value) (*ComplexVector, error) {
	switch x := v.(type) {
	case *ComplexVector:
		return &ComplexVector{Data: append([]complex128(nil), x.Data...), Missing: append([]bool(nil), x.Missing...)}, nil
	default:
		n, e := numbers(v)
		if e != nil {
			return nil, e
		}
		o := &ComplexVector{Data: make([]complex128, len(n.Data)), Missing: append([]bool(nil), n.Missing...)}
		for i, v := range n.Data {
			o.Data[i] = complex(v, 0)
		}
		return o, nil
	}
}
func (c *Context) complexOperator(op string, vals []Value) (Value, error) {
	a, e := complexNumbers(vals[0])
	if e != nil {
		return nil, e
	}
	b, e := complexNumbers(vals[1])
	if e != nil {
		return nil, e
	}
	n := max(len(a.Data), len(b.Data))
	c.warnFractionalRecycling(len(a.Data), len(b.Data))
	if op == "==" || op == "!=" {
		o := &LogicalVector{Data: make([]Logical, n)}
		for i := range o.Data {
			ai, bi := i%len(a.Data), i%len(b.Data)
			if ai < len(a.Missing) && a.Missing[ai] || bi < len(b.Missing) && b.Missing[bi] {
				o.Data[i] = NA
				continue
			}
			equal := a.Data[ai] == b.Data[bi]
			if op == "!=" {
				equal = !equal
			}
			if equal {
				o.Data[i] = True
			}
		}
		return o, nil
	}
	o := &ComplexVector{Data: make([]complex128, n), Missing: make([]bool, n)}
	for i := range o.Data {
		ai, bi := i%len(a.Data), i%len(b.Data)
		if ai < len(a.Missing) && a.Missing[ai] || bi < len(b.Missing) && b.Missing[bi] {
			o.Missing[i] = true
			continue
		}
		x, y := a.Data[ai], b.Data[bi]
		switch op {
		case "+":
			o.Data[i] = x + y
		case "-":
			o.Data[i] = x - y
		case "*":
			o.Data[i] = x * y
		case "/":
			o.Data[i] = x / y
		case "^":
			o.Data[i] = cmplx.Pow(x, y)
		default:
			return nil, fmt.Errorf("invalid operation %s for complex values", op)
		}
	}
	return o, nil
}
func missingAt(v *DoubleVector, i int) bool { return i < len(v.Missing) && v.Missing[i] }
func numbers(v Value) (*DoubleVector, error) {
	switch x := v.(type) {
	case *DoubleVector:
		o := &DoubleVector{Data: append([]float64(nil), x.Data...), Missing: append([]bool(nil), x.Missing...)}
		return o, nil
	case *LogicalVector:
		o := &DoubleVector{Data: make([]float64, len(x.Data)), Missing: make([]bool, len(x.Data))}
		for i, n := range x.Data {
			if n == NA {
				o.Missing[i] = true
			} else {
				o.Data[i] = float64(n)
			}
		}
		return o, nil
	case *IntegerVector:
		o := &DoubleVector{Data: make([]float64, len(x.Data)), Missing: append([]bool(nil), x.Missing...)}
		for i, n := range x.Data {
			o.Data[i] = float64(n)
		}
		return o, nil
	default:
		return nil, fmt.Errorf("non-numeric argument of type %s", v.Kind())
	}
}
func logical(v Value) (*LogicalVector, error) {
	switch x := v.(type) {
	case *LogicalVector:
		return &LogicalVector{Data: append([]Logical(nil), x.Data...)}, nil
	case *DoubleVector:
		o := &LogicalVector{Data: make([]Logical, len(x.Data))}
		for i, n := range x.Data {
			if i < len(x.Missing) && x.Missing[i] || math.IsNaN(n) {
				o.Data[i] = NA
			} else if n != 0 {
				o.Data[i] = True
			}
		}
		return o, nil
	case *IntegerVector:
		o := &LogicalVector{Data: make([]Logical, len(x.Data))}
		for i, n := range x.Data {
			if i < len(x.Missing) && x.Missing[i] {
				o.Data[i] = NA
			} else if n != 0 {
				o.Data[i] = True
			}
		}
		return o, nil
	default:
		return nil, fmt.Errorf("cannot coerce %s to logical", v.Kind())
	}
}
func boolValue(v bool) Value {
	if v {
		return &LogicalVector{Data: []Logical{True}}
	}
	return &LogicalVector{Data: []Logical{False}}
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
