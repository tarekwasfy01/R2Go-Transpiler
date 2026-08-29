package runtime

import (
	"fmt"
	"math"
	"math/cmplx"
	"sort"
	"strconv"
)

func init() {
	registerLoweringKernel("do_split", "0", kernelSplit)
	registerLoweringKernel("do_backsolve", "0", kernelBacksolve)
	registerLoweringKernel("do_polyroot", "0", kernelPolyroot)
	registerLoweringKernel("do_asplit", "0", kernelASplit)
	registerLoweringKernel("do_maxcol", "0", kernelMaxCol)
}

func kernelMaxCol(c *Context, frame *LoweringFrame) error {
	v, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	dims, ok := dimensions(v)
	if !ok || len(dims) != 2 {
		return fmt.Errorf("max.col expects a matrix")
	}
	x, err := numbers(v)
	if err != nil {
		return err
	}
	rows, cols := dims[0], dims[1]
	out := &IntegerVector{Data: make([]int64, rows)}
	for r := 0; r < rows; r++ {
		best := 0
		for col := 1; col < cols; col++ {
			if x.Data[col*rows+r] > x.Data[best*rows+r] {
				best = col
			}
		}
		out.Data[r] = int64(best + 1)
	}
	frame.Result = out
	return nil
}

func kernelASplit(c *Context, frame *LoweringFrame) error {
	value, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	dims, ok := dimensions(value)
	if !ok || len(dims) == 0 {
		return fmt.Errorf("asplit expects an array")
	}
	axisValue, err := frameValue(c, frame, 1)
	if err != nil {
		return err
	}
	axis, err := scalarInt(axisValue)
	if err != nil || axis < 1 || axis > len(dims) {
		return fmt.Errorf("asplit: invalid margin")
	}
	axis--
	stride := 1
	for i := 0; i < axis; i++ {
		stride *= dims[i]
	}
	block := stride * dims[axis]
	after := 1
	for i := axis + 1; i < len(dims); i++ {
		after *= dims[i]
	}
	out := &List{Data: make([]Value, dims[axis])}
	for part := 0; part < dims[axis]; part++ {
		positions := make([]int, 0, Length(value)/dims[axis])
		for group := 0; group < after; group++ {
			base := group*block + part*stride
			for j := 0; j < stride; j++ {
				positions = append(positions, base+j)
			}
		}
		out.Data[part] = takePositions(value, positions)
	}
	frame.Result = out
	return nil
}

func kernelPolyroot(c *Context, frame *LoweringFrame) error {
	value, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	coefficients, err := complexNumbers(value)
	if err != nil {
		return err
	}
	last := len(coefficients.Data) - 1
	for last >= 0 && coefficients.Data[last] == 0 {
		last--
	}
	if last <= 0 {
		frame.Result = &ComplexVector{}
		return nil
	}
	degree := last
	if degree == 1 {
		frame.Result = &ComplexVector{Data: []complex128{-coefficients.Data[0] / coefficients.Data[1]}}
		return nil
	}
	if degree == 2 {
		a, b, d := coefficients.Data[2], coefficients.Data[1], coefficients.Data[0]
		discriminant := cmplx.Sqrt(b*b - 4*a*d)
		roots := []complex128{(-b + discriminant) / (2 * a), (-b - discriminant) / (2 * a)}
		sort.SliceStable(roots, func(i, j int) bool { return imag(roots[i]) > imag(roots[j]) })
		frame.Result = &ComplexVector{Data: roots}
		return nil
	}
	lead := coefficients.Data[degree]
	radius := 1.0
	for i := 0; i < degree; i++ {
		if q := cmplx.Abs(coefficients.Data[i] / lead); q+1 > radius {
			radius = q + 1
		}
	}
	roots := make([]complex128, degree)
	for i := range roots {
		angle := 2 * math.Pi * float64(i) / float64(degree)
		roots[i] = complex(radius*math.Cos(angle), radius*math.Sin(angle))
	}
	for iteration := 0; iteration < 2000; iteration++ {
		maxDelta := 0.0
		next := append([]complex128(nil), roots...)
		for i, z := range roots {
			p := coefficients.Data[degree]
			for j := degree - 1; j >= 0; j-- {
				p = p*z + coefficients.Data[j]
			}
			denominator := complex(1, 0)
			for j, w := range roots {
				if i != j {
					denominator *= z - w
				}
			}
			if denominator == 0 {
				denominator = complex(1e-12, 1e-12)
			}
			delta := p / denominator
			next[i] = z - delta
			if d := cmplx.Abs(delta); d > maxDelta {
				maxDelta = d
			}
		}
		roots = next
		if maxDelta < 1e-13 {
			break
		}
	}
	sort.SliceStable(roots, func(i, j int) bool {
		if imag(roots[i]) == imag(roots[j]) {
			return real(roots[i]) < real(roots[j])
		}
		return imag(roots[i]) > imag(roots[j])
	})
	frame.Result = &ComplexVector{Data: roots}
	return nil
}

func kernelSplit(c *Context, frame *LoweringFrame) error {
	value, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	factor, err := frameValue(c, frame, 1)
	if err != nil {
		return err
	}
	if Length(factor) == 0 {
		return fmt.Errorf("grouping factor has length zero")
	}
	groups := map[string][]int{}
	order := []string{}
	factorElements := elements(factor)
	for i := 0; i < Length(value); i++ {
		item := factorElements[i%len(factorElements)]
		name := scalarText(item)
		if _, exists := groups[name]; !exists {
			order = append(order, name)
		}
		groups[name] = append(groups[name], i)
	}
	out := &List{Names: order, Data: make([]Value, len(order))}
	for i, name := range order {
		out.Data[i] = takePositions(value, groups[name])
	}
	frame.Result = out
	return nil
}

func kernelBacksolve(c *Context, frame *LoweringFrame) error {
	rValue, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	bValue, err := frameValue(c, frame, 1)
	if err != nil {
		return err
	}
	r, err := numbers(rValue)
	if err != nil {
		return err
	}
	b, err := numbers(bValue)
	if err != nil {
		return err
	}
	k := len(r.Data)
	if len(frame.Arguments) > 2 {
		v, e := frameValue(c, frame, 2)
		if e != nil {
			return e
		}
		k, e = scalarInt(v)
		if e != nil {
			return e
		}
	}
	if k <= 0 || len(r.Data) < k*k || len(b.Data) < k {
		return fmt.Errorf("invalid triangular system")
	}
	upper, transpose := true, false
	if len(frame.Arguments) > 3 {
		v, e := frameValue(c, frame, 3)
		if e == nil {
			upper = scalarLogical(v)
		}
	}
	if len(frame.Arguments) > 4 {
		v, e := frameValue(c, frame, 4)
		if e == nil {
			transpose = scalarLogical(v)
		}
	}
	rhs := 1
	if dims, ok := dimensions(bValue); ok && len(dims) == 2 && dims[0] == k {
		rhs = dims[1]
	}
	if len(b.Data) < k*rhs {
		return fmt.Errorf("invalid right-hand side")
	}
	x := append([]float64(nil), b.Data[:k*rhs]...)
	at := func(row, col int) float64 {
		if transpose {
			row, col = col, row
		}
		return r.Data[row+k*col]
	}
	for column := 0; column < rhs; column++ {
		base := column * k
		forward := !upper
		if transpose {
			forward = upper
		}
		if forward {
			for row := 0; row < k; row++ {
				s := x[base+row]
				for j := 0; j < row; j++ {
					s -= at(row, j) * x[base+j]
				}
				x[base+row] = s / at(row, row)
			}
		} else {
			for row := k - 1; row >= 0; row-- {
				s := x[base+row]
				for j := row + 1; j < k; j++ {
					s -= at(row, j) * x[base+j]
				}
				x[base+row] = s / at(row, row)
			}
		}
	}
	frame.Result = &DoubleVector{Data: x}
	return nil
}

func scalarLogical(value Value) bool {
	switch v := value.(type) {
	case *LogicalVector:
		return len(v.Data) != 0 && v.Data[0] == True
	case *IntegerVector:
		return len(v.Data) != 0 && v.Data[0] != 0
	case *DoubleVector:
		return len(v.Data) != 0 && v.Data[0] != 0
	case *CharacterVector:
		if len(v.Data) == 0 {
			return false
		}
		parsed, _ := strconv.ParseBool(v.Data[0])
		return parsed
	}
	return false
}
