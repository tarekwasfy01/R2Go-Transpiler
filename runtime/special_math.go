package runtime

import "math"

func init() {
	registerLoweringKernel("do_math3", "43", kernelBesselIK)
	registerLoweringKernel("do_math3", "44", kernelBesselIK)
}

func kernelBesselIK(c *Context, frame *LoweringFrame) error {
	n, params, err := vectorArguments(c, frame, 3)
	if err != nil {
		return err
	}
	out := &DoubleVector{Data: make([]float64, n), Missing: make([]bool, n)}
	for i := 0; i < n; i++ {
		x := params[0].Data[i%len(params[0].Data)]
		order := params[1].Data[i%len(params[1].Data)]
		scaled := params[2].Data[i%len(params[2].Data)] != 0
		if frame.Plan.Offset == "43" {
			out.Data[i] = besselI(order, x)
			if scaled {
				out.Data[i] *= math.Exp(-math.Abs(x))
			}
		} else {
			out.Data[i] = besselK(order, x)
			if scaled {
				out.Data[i] *= math.Exp(x)
			}
		}
	}
	frame.Result = out
	return nil
}

func vectorArguments(c *Context, frame *LoweringFrame, count int) (int, []*DoubleVector, error) {
	values := make([]*DoubleVector, count)
	n := 0
	for i := 0; i < count; i++ {
		value, err := frameValue(c, frame, i)
		if err != nil {
			return 0, nil, err
		}
		values[i], err = numbers(value)
		if err != nil {
			return 0, nil, err
		}
		if len(values[i].Data) == 0 {
			return 0, values, nil
		}
		if len(values[i].Data) > n {
			n = len(values[i].Data)
		}
	}
	return n, values, nil
}

func besselJ(order, x float64) float64 {
	if order == math.Trunc(order) && math.Abs(order) < 1<<20 {
		return math.Jn(int(order), x)
	}
	return besselJSeries(order, x)
}

func besselJSeries(order, x float64) float64 {
	g, _ := math.Lgamma(order + 1)
	term := math.Exp(order*math.Log(x/2) - g)
	sum := term
	for k := 1; k < 10000; k++ {
		term *= -(x * x / 4) / (float64(k) * (order + float64(k)))
		sum += term
		if math.Abs(term) <= math.Abs(sum)*1e-16 {
			break
		}
	}
	return sum
}

func besselY(order, x float64) float64 {
	if order == math.Trunc(order) && math.Abs(order) < 1<<20 {
		return math.Yn(int(order), x)
	}
	s := math.Sin(math.Pi * order)
	if s == 0 {
		return math.NaN()
	}
	return (math.Cos(math.Pi*order)*besselJSeries(order, x) - besselJSeries(-order, x)) / s
}

func besselI(order, x float64) float64 {
	g, _ := math.Lgamma(order + 1)
	term := math.Exp(order*math.Log(math.Abs(x)/2) - g)
	sum := term
	for k := 1; k < 10000; k++ {
		term *= (x * x / 4) / (float64(k) * (order + float64(k)))
		sum += term
		if math.Abs(term) <= math.Abs(sum)*1e-16 {
			break
		}
	}
	return sum
}

func besselK(order, x float64) float64 {
	if x <= 0 {
		return math.NaN()
	}
	// K_v(x) = integral_0^inf exp(-x cosh(t)) cosh(v t) dt.
	const steps = 4096
	const upper = 20.0
	h := upper / steps
	sum := 0.0
	for i := 0; i <= steps; i++ {
		t := float64(i) * h
		f := math.Exp(-x*math.Cosh(t)) * math.Cosh(order*t)
		weight := 2.0
		if i == 0 || i == steps {
			weight = 1
		} else if i%2 == 1 {
			weight = 4
		}
		sum += weight * f
	}
	return sum * h / 3
}

func polyGamma(x float64, order int) float64 {
	if order < 0 || x <= 0 {
		return math.NaN()
	}
	if order == 0 {
		return digamma(x)
	}
	if order == 1 {
		return trigamma(x)
	}
	sum := 0.0
	power := float64(order + 1)
	for k := 0; k < 1000000; k++ {
		term := math.Pow(x+float64(k), -power)
		sum += term
		if term < 1e-16 {
			break
		}
	}
	factorial := 1.0
	for i := 2; i <= order; i++ {
		factorial *= float64(i)
	}
	if order%2 == 0 {
		factorial = -factorial
	}
	return factorial * sum
}
