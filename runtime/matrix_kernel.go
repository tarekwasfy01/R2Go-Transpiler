package runtime

import (
	"fmt"
	"r2go/syntax"
)

func (c *Context) matrixKernel(descriptor ExecutionPlan, args []syntax.Argument, env *Environment) (Value, error) {
	switch descriptor.MatrixOperation {
	case "construct-matrix":
		return c.matrixBuiltin(args, env)
	case "construct-array":
		return c.arrayBuiltin(args, env)
	case "length":
		if len(args) != 1 {
			return nil, fmt.Errorf("length expects one argument")
		}
		v, err := c.Eval(args[0].Value, env)
		if err != nil {
			return nil, err
		}
		return &IntegerVector{Data: []int64{int64(Length(v))}}, nil
	case "transpose":
		return c.matrixTranspose(args, env)
	case "coordinates":
		return c.matrixCoordinates(descriptor, args, env)
	case "drop":
		return c.matrixDrop(args, env)
	case "product":
		return c.matrixProduct(descriptor, args, env)
	case "diagonal":
		return c.matrixDiagonal(args, env)
	case "array-runtime":
		return c.matrixPermute(args, env)
	}
	return nil, fmt.Errorf("matrix operation %s is not executable", descriptor.MatrixOperation)
}

func (c *Context) matrixDiagonal(args []syntax.Argument, env *Environment) (Value, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("diag expects an argument")
	}
	x, err := c.Eval(args[0].Value, env)
	if err != nil {
		return nil, err
	}
	if dims, ok := dimensions(x); ok && len(dims) == 2 && len(args) == 1 {
		n := dims[0]
		if dims[1] < n {
			n = dims[1]
		}
		positions := make([]int, n)
		for i := 0; i < n; i++ {
			positions[i] = i + dims[0]*i
		}
		return takePositions(x, positions), nil
	}
	rows := Length(x)
	cols := rows
	if len(args) > 1 {
		v, err := c.Eval(args[1].Value, env)
		if err != nil {
			return nil, err
		}
		if rows, err = scalarInt(v); err != nil {
			return nil, err
		}
	}
	if len(args) > 2 {
		v, err := c.Eval(args[2].Value, env)
		if err != nil {
			return nil, err
		}
		if cols, err = scalarInt(v); err != nil {
			return nil, err
		}
	}
	if rows < 0 || cols < 0 {
		return nil, fmt.Errorf("diag: invalid dimensions")
	}
	values, err := numbers(x)
	if err != nil {
		return nil, err
	}
	if len(values.Data) == 0 {
		return nil, fmt.Errorf("diag: empty diagonal value")
	}
	out := &DoubleVector{Data: make([]float64, rows*cols), Missing: make([]bool, rows*cols)}
	for i := 0; i < rows && i < cols; i++ {
		out.Data[i+rows*i], out.Missing[i+rows*i] = values.Data[i%len(values.Data)], missingAt(values, i%len(values.Data))
	}
	if err := setAttribute(out, "dim", &IntegerVector{Data: []int64{int64(rows), int64(cols)}}); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Context) matrixPermute(args []syntax.Argument, env *Environment) (Value, error) {
	if len(args) < 1 || len(args) > 3 {
		return nil, fmt.Errorf("aperm expects an array and optional permutation")
	}
	x, err := c.Eval(args[0].Value, env)
	if err != nil {
		return nil, err
	}
	dims, ok := dimensions(x)
	if !ok || len(dims) == 0 {
		return cloneValue(x), nil
	}
	perm := make([]int, len(dims))
	for i := range perm {
		perm[i] = i
	}
	if len(args) > 1 {
		v, err := c.Eval(args[1].Value, env)
		if err != nil {
			return nil, err
		}
		p, err := numbers(v)
		if err != nil || len(p.Data) != len(dims) {
			return nil, fmt.Errorf("aperm: invalid permutation")
		}
		seen := make([]bool, len(dims))
		for i, item := range p.Data {
			perm[i] = int(item) - 1
			if perm[i] < 0 || perm[i] >= len(dims) || seen[perm[i]] {
				return nil, fmt.Errorf("aperm: invalid permutation")
			}
			seen[perm[i]] = true
		}
	}
	newDims := make([]int, len(dims))
	for i, axis := range perm {
		newDims[i] = dims[axis]
	}
	positions := make([]int, Length(x))
	for target := range positions {
		coordinate, q := make([]int, len(dims)), target
		for axis := range newDims {
			coordinate[axis], q = q%newDims[axis], q/newDims[axis]
		}
		sourceCoordinates := make([]int, len(dims))
		for targetAxis, sourceAxis := range perm {
			sourceCoordinates[sourceAxis] = coordinate[targetAxis]
		}
		source, stride := 0, 1
		for sourceAxis := range dims {
			source += sourceCoordinates[sourceAxis] * stride
			stride *= dims[sourceAxis]
		}
		positions[target] = source
	}
	out := takePositions(x, positions)
	dim := &IntegerVector{Data: make([]int64, len(newDims))}
	for i, n := range newDims {
		dim.Data[i] = int64(n)
	}
	if err := setAttribute(out, "dim", dim); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Context) matrixTranspose(args []syntax.Argument, env *Environment) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("transpose expects one argument")
	}
	v, err := c.Eval(args[0].Value, env)
	if err != nil {
		return nil, err
	}
	dims, ok := dimensions(v)
	if !ok {
		dims = []int{1, Length(v)}
	}
	if len(dims) != 2 {
		return nil, fmt.Errorf("argument is not a matrix")
	}
	rows, cols := dims[0], dims[1]
	positions := make([]int, 0, rows*cols)
	for newColumn := 0; newColumn < rows; newColumn++ {
		for newRow := 0; newRow < cols; newRow++ {
			positions = append(positions, newColumn+rows*newRow)
		}
	}
	out := takePositions(v, positions)
	_ = setAttribute(out, "dim", &IntegerVector{Data: []int64{int64(cols), int64(rows)}})
	return out, nil
}

func (c *Context) matrixCoordinates(descriptor ExecutionPlan, args []syntax.Argument, env *Environment) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("coordinate matrix expects dimensions")
	}
	v, err := c.Eval(args[0].Value, env)
	if err != nil {
		return nil, err
	}
	rows, cols := 0, 0
	if dims, ok := dimensions(v); ok && len(dims) == 2 {
		rows, cols = dims[0], dims[1]
	} else {
		d, err := numbers(v)
		if err != nil || len(d.Data) != 2 {
			return nil, fmt.Errorf("a two-element dimension vector is required")
		}
		rows, cols = int(d.Data[0]), int(d.Data[1])
	}
	out := &IntegerVector{Data: make([]int64, rows*cols)}
	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {
			value := row + 1
			if descriptor.Offset == "2" {
				value = col + 1
			}
			out.Data[row+rows*col] = int64(value)
		}
	}
	_ = setAttribute(out, "dim", &IntegerVector{Data: []int64{int64(rows), int64(cols)}})
	return out, nil
}

func (c *Context) matrixDrop(args []syntax.Argument, env *Environment) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("drop expects one argument")
	}
	v, err := c.Eval(args[0].Value, env)
	if err != nil {
		return nil, err
	}
	out := cloneValue(v)
	dims, ok := dimensions(out)
	if !ok {
		return out, nil
	}
	kept := []int{}
	for _, n := range dims {
		if n != 1 {
			kept = append(kept, n)
		}
	}
	if len(kept) <= 1 {
		_ = setAttribute(out, "dim", NullValue)
		_ = setAttribute(out, "dimnames", NullValue)
		return out, nil
	}
	d := &IntegerVector{Data: make([]int64, len(kept))}
	for i, n := range kept {
		d.Data[i] = int64(n)
	}
	if err := setAttribute(out, "dim", d); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Context) matrixProduct(descriptor ExecutionPlan, args []syntax.Argument, env *Environment) (Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, fmt.Errorf("matrix product expects one or two arguments")
	}
	a, err := c.Eval(args[0].Value, env)
	if err != nil {
		return nil, err
	}
	var b Value
	if len(args) == 2 {
		b, err = c.Eval(args[1].Value, env)
		if err != nil {
			return nil, err
		}
	}
	ad, aIsMatrix := dimensions(a)
	if aIsMatrix && len(ad) != 2 {
		return nil, fmt.Errorf("first argument is not a matrix")
	}
	transposeA, transposeB := false, false
	switch descriptor.Offset {
	case "1":
		transposeA = true
		if b == nil {
			b = a
		}
	case "2":
		transposeB = true
		if b == nil {
			b = a
		}
	default:
		if b == nil {
			return nil, fmt.Errorf("matrix multiplication requires two arguments")
		}
	}
	bd, bIsMatrix := dimensions(b)
	if bIsMatrix && len(bd) != 2 {
		return nil, fmt.Errorf("second argument is not a matrix")
	}
	// GNU R treats an undimensioned left vector as a 1 x n row matrix and an
	// undimensioned right vector as an n x 1 column matrix for %*%.
	if !aIsMatrix {
		if descriptor.Offset == "1" || descriptor.Offset == "2" {
			ad = []int{Length(a), 1}
		} else {
			ad = []int{1, Length(a)}
		}
	}
	if !bIsMatrix {
		bd = []int{Length(b), 1}
	}
	aRows, aCols := ad[0], ad[1]
	if transposeA {
		aRows, aCols = aCols, aRows
	}
	bRows, bCols := bd[0], bd[1]
	if transposeB {
		bRows, bCols = bCols, bRows
	}
	if aCols != bRows {
		return nil, fmt.Errorf("non-conformable arguments")
	}
	ax, err := numbers(a)
	if err != nil {
		return nil, err
	}
	bx, err := numbers(b)
	if err != nil {
		return nil, err
	}
	out := &DoubleVector{Data: make([]float64, aRows*bCols), Missing: make([]bool, aRows*bCols)}
	at := func(row, col int) (float64, bool) {
		if transposeA {
			row, col = col, row
		}
		i := row + ad[0]*col
		return ax.Data[i], missingAt(ax, i)
	}
	bt := func(row, col int) (float64, bool) {
		if transposeB {
			row, col = col, row
		}
		i := row + bd[0]*col
		return bx.Data[i], missingAt(bx, i)
	}
	for col := 0; col < bCols; col++ {
		for row := 0; row < aRows; row++ {
			index := row + aRows*col
			sum := 0.0
			for k := 0; k < aCols; k++ {
				av, am := at(row, k)
				bv, bm := bt(k, col)
				if am || bm {
					out.Missing[index] = true
					out.Data[index] = NAReal()
					break
				}
				sum += av * bv
			}
			if !out.Missing[index] {
				out.Data[index] = sum
			}
		}
	}
	if err := setAttribute(out, "dim", &IntegerVector{Data: []int64{int64(aRows), int64(bCols)}}); err != nil {
		return nil, err
	}
	return out, nil
}
