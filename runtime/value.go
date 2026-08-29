package runtime

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"r2go/syntax"
)

func valueHasClass(v Value, wanted string) bool {
	c, ok := Attributes(v)["class"].(*CharacterVector)
	if !ok {
		return false
	}
	for _, x := range c.Data {
		if x == wanted {
			return true
		}
	}
	return false
}

func deparseExpr(e syntax.Expr) string {
	switch x := e.(type) {
	case *syntax.Literal:
		return x.Text
	case *syntax.Symbol:
		return x.Name
	case *syntax.Call:
		args := make([]string, len(x.Arguments))
		for i, a := range x.Arguments {
			args[i] = deparseExpr(a.Value)
		}
		if s, ok := x.Function.(*syntax.Symbol); ok && len(args) == 2 && strings.Contains("+-*/^", s.Name) {
			return args[0] + " " + s.Name + " " + args[1]
		}
		return deparseExpr(x.Function) + "(" + strings.Join(args, ", ") + ")"
	}
	return fmt.Sprintf("<language %T>", e)
}

type Kind string

const (
	NullKind        Kind = "NULL"
	RawKind         Kind = "raw"
	LogicalKind     Kind = "logical"
	IntegerKind     Kind = "integer"
	DoubleKind      Kind = "double"
	ComplexKind     Kind = "complex"
	CharacterKind   Kind = "character"
	ListKind        Kind = "list"
	ClosureKind     Kind = "closure"
	LanguageKind    Kind = "language"
	ConditionKind   Kind = "condition"
	EnvironmentKind Kind = "environment"
	FormulaKind     Kind = "formula"
)

type Logical int8

const (
	NA    Logical = -1
	False Logical = 0
	True  Logical = 1
)

// Logical is also used as a scalar by translated GNU R C entry points.
func (Logical) Kind() Kind { return LogicalKind }
func (v Logical) String() string {
	if v == NA {
		return "NA"
	}
	if v == True {
		return "TRUE"
	}
	return "FALSE"
}

type Value interface {
	Kind() Kind
	String() string
}
type Null struct{}

func (Null) Kind() Kind     { return NullKind }
func (Null) String() string { return "NULL" }

var NullValue Value = Null{}

type RawVector struct {
	Data []byte
	Attr map[string]Value
}

func (*RawVector) Kind() Kind { return RawKind }
func (v *RawVector) String() string {
	if len(v.Data) == 0 {
		return "raw(0)"
	}
	p := make([]string, len(v.Data))
	for i, b := range v.Data {
		p[i] = fmt.Sprintf("%02x", b)
	}
	base := strings.Join(p, " ")
	if names, ok := v.Attr["names"].(*CharacterVector); ok && len(names.Data) == len(v.Data) {
		return strings.Join(names.Data, " ") + "\n" + base
	}
	return base
}

type LogicalVector struct {
	Data []Logical
	Attr map[string]Value
}

type IntegerVector struct {
	Data    []int64
	Missing []bool
	Attr    map[string]Value
}

func (*IntegerVector) Kind() Kind { return IntegerKind }
func (v *IntegerVector) String() string {
	if v.Attr != nil {
		if dim, ok := v.Attr["dim"].(*IntegerVector); ok && len(dim.Data) == 2 {
			rows, cols := int(dim.Data[0]), int(dim.Data[1])
			if rows >= 0 && cols >= 0 && rows*cols == len(v.Data) {
				lines := make([]string, rows)
				for row := 0; row < rows; row++ {
					items := make([]string, cols)
					for col := 0; col < cols; col++ {
						i := row + rows*col
						if i < len(v.Missing) && v.Missing[i] {
							items[col] = "NA"
						} else {
							items[col] = strconv.FormatInt(v.Data[i], 10)
						}
					}
					lines[row] = strings.Join(items, " ")
				}
				return strings.Join(lines, "\n")
			}
		}
	}
	if len(v.Data) == 0 {
		return "integer(0)"
	}
	p := make([]string, len(v.Data))
	for i, x := range v.Data {
		if i < len(v.Missing) && v.Missing[i] {
			p[i] = "NA"
		} else {
			p[i] = strconv.FormatInt(x, 10) + "L"
		}
	}
	base := strings.Join(p, " ")
	// Attribute rendering is centralized here so every producer (regex,
	// factors, dimensions, imported code) exposes the same observable object.
	for _, name := range []string{"match.length", "index.type", "useBytes"} {
		if attr, ok := v.Attr[name]; ok {
			base += "\nattr(,\"" + name + "\")\n" + attr.String()
		}
	}
	return base
}

func (*LogicalVector) Kind() Kind { return LogicalKind }
func (v *LogicalVector) String() string {
	if len(v.Data) == 0 {
		return "logical(0)"
	}
	p := make([]string, len(v.Data))
	for i, x := range v.Data {
		if x == NA {
			p[i] = "NA"
		} else if x == True {
			p[i] = "TRUE"
		} else {
			p[i] = "FALSE"
		}
	}
	return strings.Join(p, " ")
}

type DoubleVector struct {
	Data    []float64
	Missing []bool
	Attr    map[string]Value
}

type ComplexVector struct {
	Data    []complex128
	Missing []bool
	Attr    map[string]Value
}

func (*ComplexVector) Kind() Kind { return ComplexKind }
func (v *ComplexVector) String() string {
	if len(v.Data) == 0 {
		return "complex(0)"
	}
	p := make([]string, len(v.Data))
	for i, z := range v.Data {
		if i < len(v.Missing) && v.Missing[i] {
			p[i] = "NA"
			continue
		}
		re, im := real(z), imag(z)
		if im == 0 {
			im = 0 // canonicalize negative zero before adding the explicit sign
		}
		sign := "+"
		if im < 0 {
			sign = "-"
			im = -im
		}
		p[i] = strconv.FormatFloat(re, 'g', -1, 64) + sign + strconv.FormatFloat(im, 'g', -1, 64) + "i"
	}
	return strings.Join(p, " ")
}

func (*DoubleVector) Kind() Kind { return DoubleKind }
func (v *DoubleVector) String() string {
	if valueHasClass(v, "POSIXct") {
		p := make([]string, len(v.Data))
		for i, x := range v.Data {
			p[i] = strconv.Quote(time.Unix(int64(x), 0).Local().Format("2006-01-02 MST"))
		}
		return strings.Join(p, " ")
	}
	if len(v.Data) == 0 {
		return "numeric(0)"
	}
	if v.Attr != nil {
		if dim, ok := v.Attr["dim"].(*IntegerVector); ok && len(dim.Data) == 2 {
			rows, cols := int(dim.Data[0]), int(dim.Data[1])
			if rows >= 0 && cols >= 0 && rows*cols == len(v.Data) {
				lines := make([]string, rows)
				for row := 0; row < rows; row++ {
					items := make([]string, cols)
					for col := 0; col < cols; col++ {
						i := row + rows*col
						if i < len(v.Missing) && v.Missing[i] {
							items[col] = "NA"
						} else if math.IsNaN(v.Data[i]) {
							items[col] = "NaN"
						} else {
							items[col] = strconv.FormatFloat(v.Data[i], 'g', -1, 64)
						}
					}
					lines[row] = strings.Join(items, " ")
				}
				return strings.Join(lines, "\n")
			}
		}
	}
	p := make([]string, len(v.Data))
	for i, x := range v.Data {
		if i < len(v.Missing) && v.Missing[i] {
			p[i] = "NA"
		} else if math.IsNaN(x) {
			p[i] = "NaN"
		} else {
			p[i] = strconv.FormatFloat(x, 'g', -1, 64)
		}
	}
	return strings.Join(p, " ")
}

type CharacterVector struct {
	Data    []string
	Missing []bool
	Attr    map[string]Value
}

func (*CharacterVector) Kind() Kind { return CharacterKind }
func (v *CharacterVector) String() string {
	if len(v.Data) == 0 {
		return "character(0)"
	}
	p := make([]string, len(v.Data))
	for i, x := range v.Data {
		if i < len(v.Missing) && v.Missing[i] {
			p[i] = "NA"
		} else {
			p[i] = strconv.Quote(x)
		}
	}
	base := strings.Join(p, " ")
	if names, ok := v.Attr["names"].(*CharacterVector); ok && len(names.Data) == len(v.Data) {
		return strings.Join(names.Data, " ") + "\n" + base
	}
	return base
}

type List struct {
	Data  []Value
	Names []string
	Attr  map[string]Value
}

func (*List) Kind() Kind { return ListKind }
func (v *List) String() string {
	if valueHasClass(v, "POSIXlt") && len(v.Data) >= 6 {
		sec, _ := numbers(v.Data[0])
		min, _ := numbers(v.Data[1])
		hour, _ := numbers(v.Data[2])
		day, _ := numbers(v.Data[3])
		mon, _ := numbers(v.Data[4])
		year, _ := numbers(v.Data[5])
		n := len(year.Data)
		p := make([]string, n)
		for i := 0; i < n; i++ {
			t := time.Date(int(year.Data[i])+1900, time.Month(int(mon.Data[i])+1), int(day.Data[i]), int(hour.Data[i]), int(min.Data[i]), int(sec.Data[i]), 0, time.Local)
			p[i] = strconv.Quote(t.Format("2006-01-02 MST"))
		}
		return strings.Join(p, " ")
	}
	if valueHasClass(v, "expression") {
		p := make([]string, len(v.Data))
		for i, x := range v.Data {
			p[i] = x.String()
		}
		return "expression(" + strings.Join(p, ", ") + ")"
	}
	p := make([]string, len(v.Data))
	for i, x := range v.Data {
		n := "[[" + strconv.Itoa(i+1) + "]]"
		if i < len(v.Names) && v.Names[i] != "" {
			name := v.Names[i]
			if name[0] >= '0' && name[0] <= '9' {
				name = "`" + name + "`"
			}
			n = "$" + name
		}
		p[i] = n + "\n" + x.String()
	}
	return strings.Join(p, "\n\n")
}

type Language struct {
	Expr syntax.Expr
	Text string
}

func (*Language) Kind() Kind { return LanguageKind }
func (v *Language) String() string {
	if v.Text != "" {
		return v.Text
	}
	return deparseExpr(v.Expr)
}

type EnvironmentValue struct {
	Env  *Environment
	Name string
}

func (*EnvironmentValue) Kind() Kind { return EnvironmentKind }
func (v *EnvironmentValue) String() string {
	if v.Name != "" {
		return "<environment: " + v.Name + ">"
	}
	return fmt.Sprintf("<environment: %p>", v.Env)
}

type Formula struct {
	Expr syntax.Expr
	Env  *Environment
	Attr map[string]Value
}

func (*Formula) Kind() Kind       { return FormulaKind }
func (v *Formula) String() string { return fmt.Sprintf("<formula %T>", v.Expr) }

type Closure struct {
	Parameters []syntax.Parameter
	Body       syntax.Expr
	Env        *Environment
	NativeBody func(*Context, *Environment) (Value, error)
	Defaults   map[string]func(*Context, *Environment) (Value, error)
}

type ConditionValue struct {
	Classes []string
	Message string
	Call    syntax.Expr
}

func (*ConditionValue) Kind() Kind       { return ConditionKind }
func (v *ConditionValue) String() string { return v.Message }

func (*Closure) Kind() Kind     { return ClosureKind }
func (*Closure) String() string { return "<closure>" }

func Length(v Value) int {
	switch x := v.(type) {
	case Null:
		return 0
	case *LogicalVector:
		return len(x.Data)
	case *RawVector:
		return len(x.Data)
	case *IntegerVector:
		return len(x.Data)
	case *DoubleVector:
		return len(x.Data)
	case *ComplexVector:
		return len(x.Data)
	case *CharacterVector:
		return len(x.Data)
	case *List:
		return len(x.Data)
	case *EnvironmentValue:
		return 1
	case *Formula:
		return 3
	default:
		return 1
	}
}

func IsTrue(v Value) (bool, error) {
	if x, ok := v.(*LogicalVector); ok {
		if len(x.Data) != 1 {
			return false, fmt.Errorf("condition has length %d", len(x.Data))
		}
		if x.Data[0] == NA {
			return false, fmt.Errorf("missing value where TRUE/FALSE needed")
		}
		return x.Data[0] == True, nil
	}
	if x, ok := v.(*DoubleVector); ok {
		if len(x.Data) != 1 {
			return false, fmt.Errorf("condition has length %d", len(x.Data))
		}
		if len(x.Missing) > 0 && x.Missing[0] || math.IsNaN(x.Data[0]) {
			return false, fmt.Errorf("missing value where TRUE/FALSE needed")
		}
		return x.Data[0] != 0, nil
	}
	if x, ok := v.(*IntegerVector); ok {
		if len(x.Data) != 1 {
			return false, fmt.Errorf("condition has length %d", len(x.Data))
		}
		if len(x.Missing) > 0 && x.Missing[0] {
			return false, fmt.Errorf("missing value where TRUE/FALSE needed")
		}
		return x.Data[0] != 0, nil
	}
	return false, fmt.Errorf("argument is not interpretable as logical: %s", v.Kind())
}
