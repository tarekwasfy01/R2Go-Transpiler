package runtime

import (
	"fmt"
	"r2go/syntax"
)

func dimensions(v Value) ([]int, bool) {
	a, ok := Attributes(v)["dim"]
	if !ok {
		return nil, false
	}
	switch x := a.(type) {
	case *IntegerVector:
		o := make([]int, len(x.Data))
		for i, n := range x.Data {
			o[i] = int(n)
		}
		return o, true
	case *DoubleVector:
		o := make([]int, len(x.Data))
		for i, n := range x.Data {
			o[i] = int(n)
		}
		return o, true
	}
	return nil, false
}

func makeArray(data Value, dims []int) (Value, error) {
	product := 1
	for _, n := range dims {
		if n < 0 {
			return nil, fmt.Errorf("negative length vectors are not allowed")
		}
		product *= n
	}
	if product > 0 && Length(data) == 0 {
		return nil, fmt.Errorf("'data' must be of a vector type")
	}
	positions := make([]int, product)
	for i := range positions {
		positions[i] = i % Length(data)
	}
	out := takePositions(data, positions)
	di := &IntegerVector{Data: make([]int64, len(dims))}
	for i, n := range dims {
		di.Data[i] = int64(n)
	}
	if err := setAttribute(out, "dim", di); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Context) matrixBuiltin(args []syntax.Argument, env *Environment) (Value, error) {
	var data Value = &LogicalVector{}
	if len(args) > 0 && args[0].Value != nil {
		v, e := c.Eval(args[0].Value, env)
		if e != nil {
			return nil, e
		}
		data = v
	}

	nrow, ncol := 0, 0
	nrowSet, ncolSet := false, false
	byrow := false
	var dimnames Value
	for i, a := range args {
		if i == 0 && a.Name == "" {
			continue
		}
		if a.Value == nil {
			continue
		}
		v, e := c.Eval(a.Value, env)
		if e != nil {
			return nil, e
		}
		switch a.Name {
		case "byrow":
			byrow = scalarLogical(v)
			continue
		case "dimnames":
			dimnames = v
			continue
		case "nrow":
			n, e := scalarInt(v)
			if e != nil {
				return nil, e
			}
			nrow, nrowSet = n, true
			continue
		case "ncol":
			n, e := scalarInt(v)
			if e != nil {
				return nil, e
			}
			ncol, ncolSet = n, true
			continue
		}
		if a.Name != "" {
			return nil, fmt.Errorf("unused matrix argument %q", a.Name)
		}
		n, e := scalarInt(v)
		if e != nil {
			return nil, e
		}
		if i == 1 {
			nrow, nrowSet = n, true
		} else if i == 2 {
			ncol, ncolSet = n, true
		} else {
			return nil, fmt.Errorf("unused matrix argument")
		}
	}

	dataLen := Length(data)
	ceilDiv := func(n, d int) int {
		if d <= 0 || n == 0 {
			return 0
		}
		return (n + d - 1) / d
	}
	switch {
	case !nrowSet && !ncolSet:
		nrow, ncol = dataLen, 1
	case nrowSet && !ncolSet:
		if nrow < 0 {
			return nil, fmt.Errorf("invalid 'nrow' value")
		}
		ncol = ceilDiv(dataLen, nrow)
	case !nrowSet && ncolSet:
		if ncol < 0 {
			return nil, fmt.Errorf("invalid 'ncol' value")
		}
		nrow = ceilDiv(dataLen, ncol)
	}
	if nrow < 0 || ncol < 0 {
		return nil, fmt.Errorf("invalid matrix dimensions")
	}

	var out Value
	var err error
	if byrow && nrow*ncol > 0 && dataLen > 0 {
		positions := make([]int, nrow*ncol)
		for col := 0; col < ncol; col++ {
			for row := 0; row < nrow; row++ {
				positions[row+nrow*col] = (row*ncol + col) % dataLen
			}
		}
		out = takePositions(data, positions)
		di := &IntegerVector{Data: []int64{int64(nrow), int64(ncol)}}
		err = setAttribute(out, "dim", di)
	} else {
		out, err = makeArray(data, []int{nrow, ncol})
	}
	if err != nil {
		return nil, err
	}
	if dimnames != nil {
		if err := setAttribute(out, "dimnames", dimnames); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (c *Context) arrayBuiltin(args []syntax.Argument, env *Environment) (Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("array requires data and dim")
	}
	data, e := c.Eval(args[0].Value, env)
	if e != nil {
		return nil, e
	}
	var dimArg *syntax.Argument
	for i := 1; i < len(args); i++ {
		if args[i].Name == "dim" || args[i].Name == "" {
			dimArg = &args[i]
			break
		}
	}
	if dimArg == nil || dimArg.Value == nil {
		return nil, fmt.Errorf("'dim' must be specified")
	}
	d, e := c.Eval(dimArg.Value, env)
	if e != nil {
		return nil, e
	}
	nums, e := numbers(d)
	if e != nil {
		return nil, e
	}
	dims := make([]int, len(nums.Data))
	for i, n := range nums.Data {
		if n < 0 || n != float64(int(n)) {
			return nil, fmt.Errorf("invalid 'dim' argument")
		}
		dims[i] = int(n)
	}
	return makeArray(data, dims)
}

func (c *Context) dimensionNamesBuiltin(name string, args []syntax.Argument, env *Environment) (Value, error) {
	if len(args) != 1 || args[0].Value == nil {
		return nil, fmt.Errorf("%s expects one argument", name)
	}
	value, err := c.Eval(args[0].Value, env)
	if err != nil {
		return nil, err
	}
	dims, ok := dimensions(value)
	if !ok || len(dims) < 2 {
		return NullValue, nil
	}
	attrs := Attributes(value)
	dimnames, ok := attrs["dimnames"].(*List)
	if !ok || len(dimnames.Data) < 2 {
		return NullValue, nil
	}
	axis := 0
	if name == "colnames" {
		axis = 1
	}
	if axis >= len(dimnames.Data) {
		return NullValue, nil
	}
	return dimnames.Data[axis], nil
}

func setDimensionNamesComponent(value Value, name string, replacement Value) error {
	dims, ok := dimensions(value)
	if !ok || len(dims) < 2 {
		return fmt.Errorf("attempt to set '%s' on an object with less than two dimensions", name)
	}
	axis := 0
	if name == "colnames" {
		axis = 1
	}
	items := make([]Value, len(dims))
	for i := range items {
		items[i] = NullValue
	}
	if current, ok := Attributes(value)["dimnames"].(*List); ok {
		copy(items, current.Data)
	}
	if _, isNull := replacement.(Null); !isNull {
		labels, ok := replacement.(*CharacterVector)
		if !ok {
			return fmt.Errorf("invalid '%s' value", name)
		}
		if len(labels.Data) != dims[axis] {
			return fmt.Errorf("length of '%s' [%d] must match extent [%d]", name, len(labels.Data), dims[axis])
		}
	}
	items[axis] = replacement
	return setAttribute(value, "dimnames", &List{Data: items})
}

func scalarInt(v Value) (int, error) {
	n, e := numbers(v)
	if e != nil || len(n.Data) == 0 || missingAt(n, 0) {
		return 0, fmt.Errorf("expected one non-missing number")
	}
	return int(n.Data[0]), nil
}

func (c *Context) arraySubset(op string, args []syntax.Argument, env *Environment) (Value, error) {
	object, e := c.Eval(args[0].Value, env)
	if e != nil {
		return nil, e
	}
	dims, ok := dimensions(object)
	if !ok || len(dims) != len(args)-1 {
		return nil, fmt.Errorf("incorrect number of dimensions")
	}
	indices := make([][]int, len(dims))
	dimnames, _ := Attributes(object)["dimnames"].(*List)
	for axis := range dims {
		a := args[axis+1]
		if a.Value == nil {
			indices[axis] = rangePositions(dims[axis])
			continue
		}
		iv, e := c.Eval(a.Value, env)
		if e != nil {
			return nil, e
		}
		var axisNames []string
		if dimnames != nil && axis < len(dimnames.Data) {
			if names, ok := dimnames.Data[axis].(*CharacterVector); ok {
				axisNames = names.Data
			}
		}
		indices[axis], e = subsetPositionsDimension(dims[axis], iv, axisNames)
		if e != nil {
			return nil, e
		}
	}
	positions := cartesianColumnMajor(indices, dims)
	if op == "[[" {
		if len(positions) != 1 || positions[0] < 0 {
			return nil, fmt.Errorf("subscript out of bounds")
		}
		return elementAt(object, positions[0]), nil
	}
	out := takePositions(object, positions)
	kept := []int{}
	for _, idx := range indices {
		if len(idx) != 1 {
			kept = append(kept, len(idx))
		}
	}
	if len(kept) > 1 {
		di := &IntegerVector{Data: make([]int64, len(kept))}
		for i, n := range kept {
			di.Data[i] = int64(n)
		}
		_ = setAttribute(out, "dim", di)
	}
	return out, nil
}
func subsetPositionsLength(n int, index Value) ([]int, error) {
	return subsetPositions(&DoubleVector{Data: make([]float64, n)}, index)
}
func subsetPositionsDimension(n int, index Value, names []string) ([]int, error) {
	axis := &DoubleVector{Data: make([]float64, n)}
	if len(names) > 0 {
		axis.Attr = map[string]Value{"names": &CharacterVector{Data: append([]string(nil), names...)}}
	}
	return subsetPositions(axis, index)
}
func rangePositions(n int) []int {
	o := make([]int, n)
	for i := range o {
		o[i] = i
	}
	return o
}
func cartesianColumnMajor(indices [][]int, dims []int) []int {
	if len(indices) == 0 {
		return nil
	}
	out := []int{0}
	stride := 1
	for axis, idx := range indices {
		next := make([]int, 0, len(out)*len(idx))
		for _, i := range idx {
			for _, base := range out {
				if i < 0 {
					next = append(next, -1)
				} else if base < 0 {
					next = append(next, -1)
				} else {
					next = append(next, base+i*stride)
				}
			}
		}
		out = next
		stride *= dims[axis]
	}
	return out
}
