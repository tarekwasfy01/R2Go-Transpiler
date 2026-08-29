package runtime

import "fmt"

func init() {
	registerLoweringKernel("do_bind", "1", kernelBind)
	registerLoweringKernel("do_bind", "2", kernelBind)
}
func kernelBind(c *Context, f *LoweringFrame) error {
	if len(f.Arguments) == 0 {
		f.Result = &DoubleVector{}
		return nil
	}
	values := make([]*DoubleVector, len(f.Arguments))
	for i := range values {
		v, e := frameValue(c, f, i)
		if e != nil {
			return e
		}
		values[i], e = numbers(v)
		if e != nil {
			return e
		}
	}
	if f.Plan.Offset == "1" {
		rows := 0
		for _, v := range values {
			if len(v.Data) > rows {
				rows = len(v.Data)
			}
		}
		if rows == 0 {
			f.Result = &DoubleVector{}
			return nil
		}
		out := &DoubleVector{Data: make([]float64, rows*len(values)), Missing: make([]bool, rows*len(values))}
		for col, v := range values {
			for row := 0; row < rows; row++ {
				i := row % len(v.Data)
				j := row + rows*col
				out.Data[j] = v.Data[i]
				out.Missing[j] = missingAt(v, i)
			}
		}
		_ = setAttribute(out, "dim", &IntegerVector{Data: []int64{int64(rows), int64(len(values))}})
		f.Result = out
		return nil
	}
	cols := 0
	for _, v := range values {
		if len(v.Data) > cols {
			cols = len(v.Data)
		}
	}
	if cols == 0 {
		return fmt.Errorf("rbind: empty arguments")
	}
	out := &DoubleVector{Data: make([]float64, len(values)*cols), Missing: make([]bool, len(values)*cols)}
	for row, v := range values {
		for col := 0; col < cols; col++ {
			i := col % len(v.Data)
			j := row + len(values)*col
			out.Data[j] = v.Data[i]
			out.Missing[j] = missingAt(v, i)
		}
	}
	_ = setAttribute(out, "dim", &IntegerVector{Data: []int64{int64(len(values)), int64(cols)}})
	f.Result = out
	return nil
}
