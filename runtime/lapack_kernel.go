package runtime

import (
	"fmt"
	"math"
)

func init() {
	for _, offset := range []string{"6", "7", "8", "61", "81", "100", "101", "200", "201", "302", "1000", "1001"} {
		registerLoweringKernel("do_lapack", offset, kernelLapack)
	}
}

func realSquare(c *Context, frame *LoweringFrame, arg int) ([]float64, int, error) {
	v, err := frameValue(c, frame, arg)
	if err != nil {
		return nil, 0, err
	}
	d, ok := dimensions(v)
	if !ok || len(d) != 2 || d[0] != d[1] {
		return nil, 0, fmt.Errorf("%s: square numeric matrix required", frame.Plan.Name)
	}
	x, err := numbers(v)
	if err != nil {
		return nil, 0, err
	}
	return append([]float64(nil), x.Data...), d[0], nil
}
func lu(a []float64, n int) ([]float64, []int, int, error) {
	p := make([]int, n)
	for i := range p {
		p[i] = i
	}
	sign := 1
	for k := 0; k < n; k++ {
		q := k
		for i := k + 1; i < n; i++ {
			if math.Abs(a[i+n*k]) > math.Abs(a[q+n*k]) {
				q = i
			}
		}
		if a[q+n*k] == 0 {
			return nil, nil, 0, fmt.Errorf("singular matrix")
		}
		if q != k {
			for j := 0; j < n; j++ {
				a[k+n*j], a[q+n*j] = a[q+n*j], a[k+n*j]
			}
			p[k], p[q] = p[q], p[k]
			sign = -sign
		}
		for i := k + 1; i < n; i++ {
			a[i+n*k] /= a[k+n*k]
			for j := k + 1; j < n; j++ {
				a[i+n*j] -= a[i+n*k] * a[k+n*j]
			}
		}
	}
	return a, p, sign, nil
}
func solveLU(a []float64, n int, b []float64, m int) ([]float64, error) {
	u, p, _, err := lu(a, n)
	if err != nil {
		return nil, err
	}
	x := make([]float64, n*m)
	for j := 0; j < m; j++ {
		for i := 0; i < n; i++ {
			x[i+n*j] = b[p[i]+n*j]
			for k := 0; k < i; k++ {
				x[i+n*j] -= u[i+n*k] * x[k+n*j]
			}
		}
		for i := n - 1; i >= 0; i-- {
			for k := i + 1; k < n; k++ {
				x[i+n*j] -= u[i+n*k] * x[k+n*j]
			}
			x[i+n*j] /= u[i+n*i]
		}
	}
	return x, nil
}
func matrixValue(data []float64, r, c int) Value {
	v := &DoubleVector{Data: data}
	_ = setAttribute(v, "dim", &IntegerVector{Data: []int64{int64(r), int64(c)}})
	return v
}
func kernelLapack(c *Context, frame *LoweringFrame) error {
	if frame.Plan.Offset == "1000" {
		frame.Result = &CharacterVector{Data: []string{"pure-go-linalg"}}
		return nil
	}
	if frame.Plan.Offset == "1001" {
		frame.Result = &CharacterVector{Data: []string{"pure-go"}}
		return nil
	}
	if frame.Plan.Offset == "61" {
		v, err := frameValue(c, frame, 0)
		if err != nil {
			return err
		}
		dims, ok := dimensions(v)
		if !ok || len(dims) != 2 {
			return fmt.Errorf("La_zlange: complex matrix required")
		}
		z, err := complexNumbers(v)
		if err != nil {
			return err
		}
		sum := 0.0
		for _, value := range z.Data {
			sum += real(value)*real(value) + imag(value)*imag(value)
		}
		frame.Result = &DoubleVector{Data: []float64{math.Sqrt(sum)}}
		return nil
	}
	a, n, err := realSquare(c, frame, 0)
	if err != nil {
		return err
	}
	switch frame.Plan.Offset {
	case "6":
		sum := 0.0
		for _, x := range a {
			sum += x * x
		}
		frame.Result = &DoubleVector{Data: []float64{math.Sqrt(sum)}}
	case "7", "8", "81":
		inverse, e := solveLU(a, n, identity(n), n)
		if e != nil {
			return e
		}
		norm := func(x []float64) float64 {
			s := 0.0
			for _, value := range x {
				s += value * value
			}
			return math.Sqrt(s)
		}
		frame.Result = &DoubleVector{Data: []float64{1 / (norm(a) * norm(inverse))}}
	case "100":
		b := make([]float64, n*n)
		for i := 0; i < n; i++ {
			b[i+n*i] = 1
		}
		if len(frame.Arguments) > 1 {
			v, e := frameValue(c, frame, 1)
			if e != nil {
				return e
			}
			d, ok := dimensions(v)
			if !ok || len(d) != 2 || d[0] != n {
				return fmt.Errorf("La_solve: non-conformable rhs")
			}
			x, e := numbers(v)
			if e != nil {
				return e
			}
			b = append([]float64(nil), x.Data...)
			r := d[1]
			out, e := solveLU(a, n, b, r)
			if e != nil {
				return e
			}
			frame.Result = matrixValue(out, n, r)
			return nil
		}
		out, e := solveLU(a, n, b, n)
		if e != nil {
			return e
		}
		frame.Result = matrixValue(out, n, n)
	case "101":
		r := append([]float64(nil), a...)
		for k := 0; k < n; k++ {
			for j := k + 1; j < n; j++ {
				dot := 0.0
				for i := 0; i < n; i++ {
					dot += r[i+n*k] * r[i+n*j]
				}
				norm := 0.0
				for i := 0; i < n; i++ {
					norm += r[i+n*k] * r[i+n*k]
				}
				if norm == 0 {
					continue
				}
				coef := dot / norm
				for i := 0; i < n; i++ {
					r[i+n*j] -= coef * r[i+n*k]
				}
			}
		}
		pivot := make([]int64, n)
		for i := range pivot {
			pivot[i] = int64(i + 1)
		}
		frame.Result = &List{Names: []string{"qr", "rank", "qraux", "pivot"}, Data: []Value{matrixValue(r, n, n), &IntegerVector{Data: []int64{int64(n)}}, &DoubleVector{Data: make([]float64, n)}, &IntegerVector{Data: pivot}}}
	case "200":
		l := make([]float64, n*n)
		for i := 0; i < n; i++ {
			for j := 0; j <= i; j++ {
				s := a[i+n*j]
				for k := 0; k < j; k++ {
					s -= l[i+n*k] * l[j+n*k]
				}
				if i == j {
					if s <= 0 {
						return fmt.Errorf("La_chol: matrix not positive definite")
					}
					l[i+n*j] = math.Sqrt(s)
				} else {
					l[i+n*j] = s / l[j+n*j]
				}
			}
		}
		frame.Result = matrixValue(l, n, n)
	case "201":
		out, e := solveLU(a, n, identity(n), n)
		if e != nil {
			return e
		}
		frame.Result = matrixValue(out, n, n)
	case "302":
		u, _, sign, e := lu(a, n)
		if e != nil {
			return e
		}
		det := float64(sign)
		for i := 0; i < n; i++ {
			det *= u[i+n*i]
		}
		frame.Result = &List{Names: []string{"modulus", "sign"}, Data: []Value{&DoubleVector{Data: []float64{math.Abs(det)}}, &IntegerVector{Data: []int64{int64(sign)}}}}
	}
	return nil
}
func identity(n int) []float64 {
	x := make([]float64, n*n)
	for i := 0; i < n; i++ {
		x[i+n*i] = 1
	}
	return x
}
