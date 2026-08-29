package runtime

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"r2go/syntax"
	"strconv"
	"strings"
	"time"
)

type Context struct {
	Global   *Environment
	Options  map[string]Value
	Output   io.Writer
	Warnings []*ConditionValue
	RNG      *rand.Rand
	Started  time.Time
	Locale   string
	Host     Host
}

func NewContext() *Context {
	c := NewContextWithHost(LocalHost{})
	return c
}
func NewContextWithHost(host Host) *Context {
	if host == nil {
		host = LocalHost{}
	}
	c := &Context{Output: host.Stdout(), Options: map[string]Value{}, RNG: rand.New(rand.NewSource(1)), Started: host.Now(), Locale: "C", Host: host}
	c.Global = NewEnvironment(nil)
	return c
}
func (c *Context) EvalProgram(p *syntax.Program) (Value, error) {
	v := NullValue
	var err error
	for _, e := range p.Expressions {
		v, err = c.Eval(e, c.Global)
		if err != nil {
			return nil, err
		}
	}
	return v, nil
}

type control struct {
	kind  string
	value Value
}

func (x *control) Error() string { return x.kind }

func (c *Context) Eval(expr syntax.Expr, env *Environment) (Value, error) {
	switch e := expr.(type) {
	case *syntax.Literal:
		return evalLiteral(e)
	case *syntax.Symbol:
		return env.Get(c, e.Name)
	case *syntax.Block:
		v := NullValue
		for _, item := range e.Expressions {
			x, err := c.Eval(item, env)
			if err != nil {
				return nil, err
			}
			v = x
		}
		return v, nil
	case *syntax.Function:
		return &Closure{Parameters: e.Parameters, Body: e.Body, Env: env}, nil
	case *syntax.If:
		cond, err := c.Eval(e.Condition, env)
		if err != nil {
			return nil, err
		}
		ok, err := IsTrue(cond)
		if err != nil {
			return nil, err
		}
		if ok {
			return c.Eval(e.Then, env)
		}
		if e.Else != nil {
			return c.Eval(e.Else, env)
		}
		return NullValue, nil
	case *syntax.While:
		for {
			v, err := c.Eval(e.Condition, env)
			if err != nil {
				return nil, err
			}
			ok, err := IsTrue(v)
			if err != nil {
				return nil, err
			}
			if !ok {
				return NullValue, nil
			}
			_, err = c.Eval(e.Body, env)
			if x, ok := err.(*control); ok {
				if x.kind == "break" {
					return NullValue, nil
				}
				if x.kind == "next" {
					continue
				}
			}
			if err != nil {
				return nil, err
			}
		}
	case *syntax.Repeat:
		for {
			_, err := c.Eval(e.Body, env)
			if x, ok := err.(*control); ok {
				if x.kind == "break" {
					return NullValue, nil
				}
				if x.kind == "next" {
					continue
				}
			}
			if err != nil {
				return nil, err
			}
		}
	case *syntax.For:
		seq, err := c.Eval(e.Sequence, env)
		if err != nil {
			return nil, err
		}
		for _, item := range elements(seq) {
			env.Set(e.Variable, item)
			_, err = c.Eval(e.Body, env)
			if x, ok := err.(*control); ok {
				if x.kind == "break" {
					return NullValue, nil
				}
				if x.kind == "next" {
					continue
				}
			}
			if err != nil {
				return nil, err
			}
		}
		return NullValue, nil
	case *syntax.Call:
		return c.evalCall(e, env)
	default:
		return nil, fmt.Errorf("unsupported AST node %T", expr)
	}
}

func evalLiteral(e *syntax.Literal) (Value, error) {
	switch e.Kind {
	case syntax.NullLiteral:
		return NullValue, nil
	case syntax.StringLiteral:
		return &CharacterVector{Data: []string{e.Text}}, nil
	case syntax.LogicalLiteral:
		v := False
		if e.Text == "TRUE" || e.Text == "T" {
			v = True
		}
		return &LogicalVector{Data: []Logical{v}}, nil
	case syntax.NALiteral:
		if e.Text == "NA_character_" {
			return &CharacterVector{Data: []string{""}, Missing: []bool{true}}, nil
		}
		if e.Text == "NA_integer_" {
			return &IntegerVector{Data: []int64{0}, Missing: []bool{true}}, nil
		}
		if e.Text == "NA_real_" {
			return &DoubleVector{Data: []float64{NAReal()}, Missing: []bool{true}}, nil
		}
		if e.Text == "NA_complex_" {
			return &ComplexVector{Data: []complex128{0}, Missing: []bool{true}}, nil
		}
		return &LogicalVector{Data: []Logical{NA}}, nil
	case syntax.NumberLiteral:
		if e.Text == "Inf" {
			return &DoubleVector{Data: []float64{math.Inf(1)}}, nil
		}
		if e.Text == "NaN" {
			return &DoubleVector{Data: []float64{math.NaN()}}, nil
		}
		if strings.HasSuffix(e.Text, "i") {
			s := strings.TrimSuffix(e.Text, "i")
			n, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return nil, err
			}
			return &ComplexVector{Data: []complex128{complex(0, n)}}, nil
		}
		if strings.HasSuffix(e.Text, "L") {
			s := strings.TrimSuffix(e.Text, "L")
			n, err := strconv.ParseInt(s, 0, 64)
			if err != nil {
				return nil, err
			}
			return &IntegerVector{Data: []int64{n}}, nil
		}
		s := strings.TrimSuffix(e.Text, "i")
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		return &DoubleVector{Data: []float64{n}}, nil
	}
	return nil, fmt.Errorf("unknown literal")
}

func (c *Context) evalCall(call *syntax.Call, env *Environment) (Value, error) {
	// Namespace-qualified calls are language objects in R. Route the selected
	// function through the normal dispatcher so builtins and closures share the
	// same call semantics (for example base::sum(...)).
	if qualified, ok := call.Function.(*syntax.Call); ok {
		if operator, ok := qualified.Function.(*syntax.Symbol); ok && (operator.Name == "::" || operator.Name == ":::") && len(qualified.Arguments) == 2 {
			if name, ok := qualified.Arguments[1].Value.(*syntax.Symbol); ok {
				return c.evalCall(&syntax.Call{Function: &syntax.Symbol{Name: name.Name, At: name.At}, Arguments: call.Arguments, At: call.At}, env)
			}
		}
	}
	if s, ok := call.Function.(*syntax.Symbol); ok {
		if s.Name == "print" || s.Name == "summary" || s.Name == "as.character" {
			if v, done, e := c.dispatchS3(s.Name, call.Arguments, env); done {
				return v, e
			}
		}
		// These Base-R language helpers need evaluator context or lazy arguments.
		// Route them before the generated GNU-R primitive matrix so translated C
		// control flow cannot replace their host semantics.
		switch s.Name {
		case "stopifnot":
			return c.stopIfNot(call.Arguments, env)
		case "ls":
			return c.listEnvironment(call.Arguments, env)
		case "new.env":
			return c.newEnvironment(call.Arguments, env)
		case "rownames", "colnames":
			return c.dimensionNamesBuiltin(s.Name, call.Arguments, env)
		case "apply":
			return c.matrixApplyBuiltin(call.Arguments, env)
		case "split", "lapply", "sapply":
			return c.vectorBuiltin(s.Name, call.Arguments, env)
		case "order":
			return c.moreBuiltin(s.Name, call.Arguments, env)
		case "median", "var", "sd", "quantile":
			return c.statisticsBuiltin(s.Name, call.Arguments, env)
		}
		if v, handled, err := c.dispatchPrimitive(s.Name, call.Arguments, env); handled {
			return v, err
		}
		switch s.Name {
		case "(":
			return c.Eval(call.Arguments[0].Value, env)
		case "<-", "=", "<<-":
			return c.assign(s.Name, call, env)
		case "return":
			v := NullValue
			var err error
			if len(call.Arguments) > 0 {
				v, err = c.Eval(call.Arguments[0].Value, env)
			}
			if err != nil {
				return nil, err
			}
			return nil, &control{kind: "return", value: v}
		case "break", "next":
			return nil, &control{kind: s.Name}
		case "quote":
			if len(call.Arguments) != 1 {
				return nil, fmt.Errorf("quote expects one argument")
			}
			return &Language{Expr: call.Arguments[0].Value}, nil
		case "missing":
			if len(call.Arguments) != 1 {
				return nil, fmt.Errorf("missing expects one argument")
			}
			symbol, ok := call.Arguments[0].Value.(*syntax.Symbol)
			if !ok {
				return nil, fmt.Errorf("invalid use of missing")
			}
			_, binding, ok := env.Find(symbol.Name)
			if !ok {
				return nil, fmt.Errorf("missing can only be used for arguments")
			}
			status, ok := binding.(interface{ IsMissing() bool })
			return boolValue(ok && status.IsMissing()), nil
		case "stop":
			return nil, c.makeCondition("error", call.Arguments, env, call)
		case "warning":
			condition := c.makeCondition("warning", call.Arguments, env, call)
			c.Warnings = append(c.Warnings, condition)
			return NullValue, nil
		case "tryCatch":
			return c.tryCatch(call.Arguments, env)
		case "UseMethod":
			return c.useMethod(call.Arguments, env)
		case "NextMethod":
			return c.nextMethod(call.Arguments, env)
		case "%in%", ":", "&&", "||", "!":
			return c.operator(s.Name, call.Arguments, env)
		case "c", "list", "length", "print", "typeof", "is.null", "is.na", "complete.cases", "is.data.frame", "is.factor", "is.raw", "inherits", "as.logical", "as.integer", "as.double", "as.numeric", "as.complex", "as.raw", "as.character", "as.vector", "numeric", "integer", "logical", "character", "complex", "raw", "rawToChar", "charToRaw", "Re", "Im", "Mod", "Arg", "Conj", "sum", "prod", "min", "max", "mean", "seq_along", "names", "levels", "dim", "nrow", "ncol", "attr", "attributes", "class", "unclass", "structure", "eval", "conditionMessage":
			return c.builtin(s.Name, call.Arguments, env)
		case "matrix":
			return c.matrixBuiltin(call.Arguments, env)
		case "t":
			return c.matrixTranspose(call.Arguments, env)
		case "array":
			return c.arrayBuiltin(call.Arguments, env)
		case "data.frame":
			return c.dataFrameBuiltin(call.Arguments, env)
		case "factor":
			return c.factorBuiltin(call.Arguments, env)
		case "seq", "seq.int", "rep", "any", "all", "which", "unique", "duplicated", "match", "paste", "paste0", "table", "ifelse", "lapply", "sapply":
			return c.vectorBuiltin(s.Name, call.Arguments, env)
		case "rev", "head", "tail", "sort", "order", "rank", "diff", "pmin", "pmax", "nchar", "tolower", "toupper", "trimws", "substr", "substring", "startsWith", "endsWith", "strsplit":
			return c.moreBuiltin(s.Name, call.Arguments, env)
		case "getwd", "setwd", "file.exists", "dir.exists", "dir.create", "basename", "dirname", "normalizePath", "readLines", "writeLines", "list.files", "read.csv", "read.csv2", "write.csv", "write.csv2":
			return c.ioBuiltin(s.Name, call.Arguments, env)
		case "serialize", "unserialize":
			return c.serializationBuiltin(s.Name, call.Arguments, env)
		}
	}
	fn, err := c.Eval(call.Function, env)
	if err != nil {
		return nil, err
	}
	closure, ok := fn.(*Closure)
	if !ok {
		return nil, fmt.Errorf("attempt to apply non-function of type %s", fn.Kind())
	}
	return c.callClosure(closure, call.Arguments, env)
}

func (c *Context) stopIfNot(args []syntax.Argument, env *Environment) (Value, error) {
	for i, arg := range args {
		if arg.Value == nil {
			return nil, fmt.Errorf("argument %d is empty", i+1)
		}
		value, err := c.Eval(arg.Value, env)
		if err != nil {
			return nil, err
		}
		ok, err := IsTrue(value)
		if err != nil || !ok {
			return nil, fmt.Errorf("%s is not TRUE", deparseExpr(arg.Value))
		}
	}
	return NullValue, nil
}

func (c *Context) listEnvironment(args []syntax.Argument, env *Environment) (Value, error) {
	target := env
	allNames := false
	for i, arg := range args {
		switch arg.Name {
		case "all.names":
			value, err := c.Eval(arg.Value, env)
			if err != nil {
				return nil, err
			}
			allNames, err = IsTrue(value)
			if err != nil {
				return nil, err
			}
		case "envir":
			value, err := c.Eval(arg.Value, env)
			if err != nil {
				return nil, err
			}
			var ok bool
			target, ok = value.(*Environment)
			if !ok {
				return nil, fmt.Errorf("invalid 'envir' argument")
			}
		case "sorted":
			// Environment.Names is deterministic and sorted already.
		case "name", "pos":
			return nil, fmt.Errorf("character and numeric search positions are not available without an attached search path")
		case "":
			if i == 0 {
				value, err := c.Eval(arg.Value, env)
				if err != nil {
					return nil, err
				}
				var ok bool
				target, ok = value.(*Environment)
				if !ok {
					return nil, fmt.Errorf("invalid environment")
				}
			}
		}
	}
	return &CharacterVector{Data: target.Names(allNames)}, nil
}

func (c *Context) newEnvironment(args []syntax.Argument, env *Environment) (Value, error) {
	parent := env
	for i, arg := range args {
		// new.env(hash = TRUE, parent = parent.frame(), size = 29L)
		// has parent as its second formal argument.
		if arg.Name != "parent" && !(arg.Name == "" && i == 1) {
			continue
		}
		value, err := c.Eval(arg.Value, env)
		if err != nil {
			return nil, err
		}
		var ok bool
		parent, ok = value.(*Environment)
		if !ok {
			return nil, fmt.Errorf("'parent' must be an environment")
		}
	}
	return NewEnvironment(parent), nil
}

func (c *Context) assign(op string, call *syntax.Call, env *Environment) (Value, error) {
	if len(call.Arguments) != 2 {
		return nil, fmt.Errorf("assignment expects two operands")
	}
	v, err := c.Eval(call.Arguments[1].Value, env)
	if err != nil {
		return nil, err
	}
	if s, ok := call.Arguments[0].Value.(*syntax.Symbol); ok {
		if op == "<<-" {
			env.SuperSet(s.Name, v)
		} else {
			env.Set(s.Name, v)
		}
		return v, nil
	}
	target, ok := call.Arguments[0].Value.(*syntax.Call)
	if !ok {
		return nil, fmt.Errorf("invalid assignment target")
	}
	if err := c.assignCall(op, target, v, env); err != nil {
		return nil, err
	}
	return v, nil
}

func (c *Context) assignCall(op string, target *syntax.Call, replacement Value, env *Environment) error {
	fn, ok := target.Function.(*syntax.Symbol)
	if !ok {
		return fmt.Errorf("invalid replacement function")
	}
	if len(target.Arguments) == 0 {
		return fmt.Errorf("invalid assignment target")
	}
	if fn.Name == "rownames" || fn.Name == "colnames" {
		if len(target.Arguments) != 1 {
			return fmt.Errorf("%s replacement expects one object", fn.Name)
		}
		parent := target.Arguments[0].Value
		original, err := c.Eval(parent, env)
		if err != nil {
			return err
		}
		object := cloneValue(original)
		if err := setDimensionNamesComponent(object, fn.Name, replacement); err != nil {
			return err
		}
		return c.storeReplacement(op, parent, object, env)
	}
	descriptor, known := ExecutionPlanByName[fn.Name+"<-"]
	if !known || descriptor.Opcode != "REPLACE" {
		return fmt.Errorf("replacement function %s<- is not described by the translation matrix", fn.Name)
	}
	parent := target.Arguments[0].Value
	original, err := c.Eval(parent, env)
	if err != nil {
		return err
	}
	object := cloneValue(original)
	switch descriptor.ReplacementPolicy {
	case "subscript":
		var positions []int
		if len(target.Arguments) > 2 {
			dims, has := dimensions(object)
			if !has || len(dims) != len(target.Arguments)-1 {
				return fmt.Errorf("incorrect number of dimensions")
			}
			indices := make([][]int, len(dims))
			for axis := range dims {
				a := target.Arguments[axis+1]
				if a.Value == nil {
					indices[axis] = rangePositions(dims[axis])
					continue
				}
				iv, e := c.Eval(a.Value, env)
				if e != nil {
					return e
				}
				indices[axis], e = subsetPositionsLength(dims[axis], iv)
				if e != nil {
					return e
				}
			}
			positions = cartesianColumnMajor(indices, dims)
		} else {
			if len(target.Arguments) != 2 || target.Arguments[1].Value == nil {
				return fmt.Errorf("replacement requires one subscript")
			}
			index, e := c.Eval(target.Arguments[1].Value, env)
			if e != nil {
				return e
			}
			positions, e = subsetPositions(object, index)
			if e != nil {
				return e
			}
		}
		if descriptor.CEntry == "do_subassign2" && len(positions) != 1 {
			return fmt.Errorf("attempt to select more than one element")
		}
		object, err = replacePositions(object, positions, replacement)
		if err != nil {
			return err
		}
	case "member":
		if len(target.Arguments) != 2 {
			return fmt.Errorf("$ replacement requires a name")
		}
		name, ok := target.Arguments[1].Value.(*syntax.Symbol)
		if !ok {
			return fmt.Errorf("invalid $ name")
		}
		if e := replaceDollar(object, name.Name, replacement); e != nil {
			return e
		}
	case "attribute-key":
		if len(target.Arguments) != 1 {
			return fmt.Errorf("%s expects one object", descriptor.Name)
		}
		key := strings.TrimSuffix(descriptor.Name, "<-")
		if key == "oldClass" {
			key = "class"
		}
		if key == "names" {
			err = replaceNames(object, replacement)
		} else {
			err = setAttribute(object, key, replacement)
		}
		if err != nil {
			return err
		}
	case "attribute-named":
		if len(target.Arguments) != 2 {
			return fmt.Errorf("attr replacement expects object and name")
		}
		nameValue, e := c.Eval(target.Arguments[1].Value, env)
		if e != nil {
			return e
		}
		name, ok := nameValue.(*CharacterVector)
		if !ok || len(name.Data) != 1 {
			return fmt.Errorf("attribute name must be one string")
		}
		if e := setAttribute(object, name.Data[0], replacement); e != nil {
			return e
		}
	case "attribute-map":
		if len(target.Arguments) != 1 {
			return fmt.Errorf("attributes replacement expects one object")
		}
		if err = replaceAttributeMap(object, replacement); err != nil {
			return err
		}
	default:
		return fmt.Errorf("replacement policy %s is not implemented", descriptor.ReplacementPolicy)
	}
	return c.storeReplacement(op, parent, object, env)
}

func (c *Context) storeReplacement(op string, target syntax.Expr, value Value, env *Environment) error {
	if root, ok := target.(*syntax.Symbol); ok {
		if op == "<<-" {
			env.SuperSet(root.Name, value)
		} else {
			env.Set(root.Name, value)
		}
		return nil
	}
	call, ok := target.(*syntax.Call)
	if !ok {
		return fmt.Errorf("invalid nested replacement target")
	}
	return c.assignCall(op, call, value, env)
}

func (c *Context) callClosure(fn *Closure, args []syntax.Argument, caller *Environment) (Value, error) {
	actuals := make([]ActualArgument, 0, len(args))
	for _, arg := range args {
		if symbol, ok := arg.Value.(*syntax.Symbol); ok && symbol.Name == "..." {
			_, binding, found := caller.Find("...")
			if !found {
				return nil, fmt.Errorf("'...' used in an incorrect context")
			}
			dots, ok := binding.(*DotsBinding)
			if !ok {
				return nil, fmt.Errorf("invalid ... binding")
			}
			actuals = append(actuals, dots.Arguments...)
			continue
		}
		actuals = append(actuals, ActualArgument{Argument: arg, Env: caller})
	}
	return c.callClosureActual(fn, actuals)
}

func (c *Context) callClosureActual(fn *Closure, args []ActualArgument) (Value, error) {
	frame := NewEnvironment(fn.Env)
	matched, dots, err := matchArguments(fn.Parameters, args)
	if err != nil {
		return nil, err
	}
	for i, parameter := range fn.Parameters {
		if parameter.Name == "..." {
			frame.Bind("...", &DotsBinding{Arguments: dots})
			continue
		}
		actual := matched[i]
		if actual != nil && actual.HasValue {
			frame.Bind(parameter.Name, &Eager{Value: actual.EagerValue})
		} else if actual != nil && actual.Argument.Value != nil {
			frame.Bind(parameter.Name, &Promise{Expr: actual.Argument.Value, Env: actual.Env, Supplied: true})
		} else if nativeDefault := fn.Defaults[parameter.Name]; nativeDefault != nil {
			value, defaultErr := nativeDefault(c, frame)
			if defaultErr != nil {
				return nil, defaultErr
			}
			frame.Bind(parameter.Name, &Eager{Value: value})
		} else if parameter.Default != nil {
			frame.Bind(parameter.Name, &Promise{Expr: parameter.Default, Env: frame, Supplied: false})
		} else {
			frame.Bind(parameter.Name, &MissingBinding{Name: parameter.Name})
		}
	}
	for _, parameter := range fn.Parameters {
		if parameter.Name == "..." {
			continue
		}
		if binding, ok := frame.bindings[parameter.Name]; ok {
			frame.Bind(".S3Object", binding)
		}
		break
	}
	if fn.NativeBody != nil {
		return fn.NativeBody(c, frame)
	}
	v, err := c.Eval(fn.Body, frame)
	if x, ok := err.(*control); ok && x.kind == "return" {
		return x.value, nil
	}
	return v, err
}

func matchArguments(parameters []syntax.Parameter, args []ActualArgument) ([]*ActualArgument, []ActualArgument, error) {
	matched := make([]*ActualArgument, len(parameters))
	used := make([]bool, len(args))
	dotsIndex := len(parameters)
	for i, parameter := range parameters {
		if parameter.Name == "..." {
			dotsIndex = i
			break
		}
	}

	// Pass 1: exact names, including formals following ... .
	for ai := range args {
		if args[ai].Argument.Name == "" {
			continue
		}
		for pi := range parameters {
			if parameters[pi].Name == "..." || parameters[pi].Name != args[ai].Argument.Name {
				continue
			}
			if matched[pi] != nil {
				return nil, nil, fmt.Errorf("formal argument %q matched by multiple actual arguments", parameters[pi].Name)
			}
			matched[pi], used[ai] = &args[ai], true
			break
		}
	}

	// Pass 2: unique partial names only before ... .
	for ai := range args {
		if used[ai] || args[ai].Argument.Name == "" {
			continue
		}
		candidate := -1
		for pi := 0; pi < dotsIndex; pi++ {
			if matched[pi] == nil && strings.HasPrefix(parameters[pi].Name, args[ai].Argument.Name) {
				if candidate >= 0 {
					return nil, nil, fmt.Errorf("argument %q matches multiple formal arguments", args[ai].Argument.Name)
				}
				candidate = pi
			}
		}
		if candidate >= 0 {
			matched[candidate], used[ai] = &args[ai], true
		}
	}

	// Pass 3: unnamed arguments positionally match only before ... .
	nextFormal := 0
	for ai := range args {
		if used[ai] || args[ai].Argument.Name != "" {
			continue
		}
		for nextFormal < dotsIndex && matched[nextFormal] != nil {
			nextFormal++
		}
		if nextFormal < dotsIndex {
			matched[nextFormal], used[ai] = &args[ai], true
			nextFormal++
		}
	}

	var dots []ActualArgument
	for ai := range args {
		if used[ai] {
			continue
		}
		if dotsIndex < len(parameters) {
			dots = append(dots, args[ai])
			continue
		}
		label := args[ai].Argument.Name
		if label == "" {
			label = fmt.Sprintf("position %d", ai+1)
		}
		return nil, nil, fmt.Errorf("unused argument (%s)", label)
	}
	return matched, dots, nil
}

func elements(v Value) []Value {
	switch x := v.(type) {
	case *RawVector:
		o := make([]Value, len(x.Data))
		for i, b := range x.Data {
			o[i] = &RawVector{Data: []byte{b}}
		}
		return o
	case *DoubleVector:
		o := make([]Value, len(x.Data))
		for i, n := range x.Data {
			o[i] = &DoubleVector{Data: []float64{n}, Missing: []bool{i < len(x.Missing) && x.Missing[i]}}
		}
		return o
	case *LogicalVector:
		o := make([]Value, len(x.Data))
		for i, n := range x.Data {
			o[i] = &LogicalVector{Data: []Logical{n}}
		}
		return o
	case *IntegerVector:
		o := make([]Value, len(x.Data))
		for i, n := range x.Data {
			o[i] = &IntegerVector{Data: []int64{n}, Missing: []bool{i < len(x.Missing) && x.Missing[i]}}
		}
		return o
	case *CharacterVector:
		o := make([]Value, len(x.Data))
		for i, n := range x.Data {
			o[i] = &CharacterVector{Data: []string{n}, Missing: []bool{i < len(x.Missing) && x.Missing[i]}}
		}
		return o
	case *ComplexVector:
		o := make([]Value, len(x.Data))
		for i, z := range x.Data {
			o[i] = &ComplexVector{Data: []complex128{z}, Missing: []bool{i < len(x.Missing) && x.Missing[i]}}
		}
		return o
	case *List:
		return x.Data
	default:
		return []Value{v}
	}
}
