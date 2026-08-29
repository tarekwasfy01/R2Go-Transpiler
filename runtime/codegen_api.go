package runtime

import "r2go/syntax"

// NumericVector and the scalar constructors are the stable surface used by
// native R-to-Go output. Generated code contains explicit values and calls.
func NumericVector(values ...float64) Value {
	return &DoubleVector{Data: append([]float64(nil), values...)}
}

func LogicalScalar(value bool) Value { return boolValue(value) }

func CharacterScalar(value string) Value {
	return &CharacterVector{Data: []string{value}}
}

// Binary invokes one vectorized Pure-Go operator kernel. It does not parse,
// decode, or evaluate a stored R program.
func Binary(operator string, left, right Value) (Value, error) {
	ctx := NewContext()
	ctx.Global.Set(".r2go_left", left)
	ctx.Global.Set(".r2go_right", right)
	return ctx.operator(operator, []syntax.Argument{
		{Value: &syntax.Symbol{Name: ".r2go_left"}},
		{Value: &syntax.Symbol{Name: ".r2go_right"}},
	}, ctx.Global)
}

func MustBinary(operator string, left, right Value) Value {
	value, err := Binary(operator, left, right)
	if err != nil {
		panic(err)
	}
	return value
}
