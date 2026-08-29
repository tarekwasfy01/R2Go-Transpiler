package runtime

import (
	"fmt"
	"math"
	"r2go/syntax"
)

func init() {
	for _, entry := range []string{"do_random1", "do_random2", "do_random3", "do_sample", "do_sample2", "do_RNGkind", "do_setseed"} {
		registerLoweringKernel(entry, "0", kernelRandomFamily)
	}
	for _, offset := range []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13"} {
		registerLoweringKernel("do_random1", offset, kernelRandomFamily)
		registerLoweringKernel("do_random2", offset, kernelRandomFamily)
	}
}

func (c *Context) random() float64 { return c.RNG.Float64() }

func randomNumbers(c *Context, frame *LoweringFrame) (int, []*DoubleVector, error) {
	if len(frame.Arguments) == 0 {
		return 0, nil, fmt.Errorf("%s: missing n", frame.Plan.Name)
	}
	nv, err := frameValue(c, frame, 0)
	if err != nil {
		return 0, nil, err
	}
	n := Length(nv)
	if n == 1 {
		n, err = scalarInt(nv)
	}
	if err != nil || n < 0 {
		return 0, nil, fmt.Errorf("%s: invalid n", frame.Plan.Name)
	}
	params := make([]*DoubleVector, len(frame.Arguments)-1)
	for i := range params {
		v, err := frameValue(c, frame, i+1)
		if err != nil {
			return 0, nil, err
		}
		params[i], err = numbers(v)
		if err != nil || len(params[i].Data) == 0 {
			return 0, nil, fmt.Errorf("%s: invalid parameter", frame.Plan.Name)
		}
	}
	return n, params, nil
}

func rv(v *DoubleVector, i int) float64 { return v.Data[i%len(v.Data)] }
func (c *Context) normal() float64      { return c.RNG.NormFloat64() }
func (c *Context) gamma(shape float64) float64 {
	if shape <= 0 {
		return math.NaN()
	}
	if shape < 1 {
		return c.gamma(shape+1) * math.Pow(c.random(), 1/shape)
	}
	d := shape - 1.0/3
	z := 1.0 / math.Sqrt(9*d)
	for {
		x := c.normal()
		v := 1 + z*x
		if v <= 0 {
			continue
		}
		v *= v * v
		u := c.random()
		if u < 1-0.0331*x*x*x*x || math.Log(u) < 0.5*x*x+d*(1-v+math.Log(v)) {
			return d * v
		}
	}
}
func (c *Context) poisson(lambda float64) float64 {
	if lambda < 0 {
		return math.NaN()
	}
	if lambda > 30 {
		return math.Max(0, math.Round(lambda+math.Sqrt(lambda)*c.normal()))
	}
	l, k, p := math.Exp(-lambda), 0, 1.0
	for p > l {
		k++
		p *= c.random()
	}
	return float64(k - 1)
}

func kernelRandomFamily(c *Context, frame *LoweringFrame) error {
	if frame.Plan.CEntry == "do_setseed" {
		v, err := frameValue(c, frame, 0)
		if err != nil {
			return err
		}
		seed, err := scalarInt(v)
		if err != nil {
			return err
		}
		c.RNG.Seed(int64(seed))
		frame.Result = NullValue
		return nil
	}
	if frame.Plan.CEntry == "do_RNGkind" {
		frame.Result = &CharacterVector{Data: []string{"Mersenne-Twister", "Inversion", "Rejection"}}
		return nil
	}
	if frame.Plan.CEntry == "do_sample" || frame.Plan.CEntry == "do_sample2" {
		return kernelSample(c, frame)
	}
	n, p, err := randomNumbers(c, frame)
	if err != nil {
		return err
	}
	out := &DoubleVector{Data: make([]float64, n)}
	for i := range out.Data {
		out.Data[i] = c.randomDraw(frame.Plan, p, i)
	}
	frame.Result = out
	return nil
}

func (c *Context) randomDraw(plan ExecutionPlan, p []*DoubleVector, i int) float64 {
	a := func(j int) float64 { return rv(p[j], i) }
	switch plan.CEntry + ":" + plan.Offset {
	case "do_random1:0":
		return 2 * c.gamma(a(0)/2)
	case "do_random1:1":
		return -math.Log1p(-c.random()) * a(0)
	case "do_random1:2":
		return math.Floor(math.Log1p(-c.random()) / math.Log1p(-a(0)))
	case "do_random1:3":
		return c.poisson(a(0))
	case "do_random1:4":
		return c.normal() / math.Sqrt((2*c.gamma(a(0)/2))/a(0))
	case "do_random1:5":
		n := int(a(0))
		s := 0.0
		for j := 1; j <= n; j++ {
			if c.random() < .5 {
				s += float64(j)
			}
		}
		return s
	case "do_random2:0":
		x, y := c.gamma(a(0)), c.gamma(a(1))
		return x / (x + y)
	case "do_random2:1":
		n := int(a(0))
		s := 0.0
		for j := 0; j < n; j++ {
			if c.random() < a(1) {
				s++
			}
		}
		return s
	case "do_random2:2":
		return a(0) + a(1)*math.Tan(math.Pi*(c.random()-.5))
	case "do_random2:3":
		return (2 * c.gamma(a(0)/2) / a(0)) / (2 * c.gamma(a(1)/2) / a(1))
	case "do_random2:4":
		return c.gamma(a(0)) * a(1)
	case "do_random2:5":
		return math.Exp(a(0) + a(1)*c.normal())
	case "do_random2:6":
		u := c.random()
		return a(0) + a(1)*math.Log(u/(1-u))
	case "do_random2:7":
		return c.poisson(c.gamma(a(0)) * (1 - a(1)) / a(1))
	case "do_random2:8":
		return a(0) + a(1)*c.normal()
	case "do_random2:9":
		return a(0) + (a(1)-a(0))*c.random()
	case "do_random2:10":
		return a(1) * math.Pow(-math.Log1p(-c.random()), 1/a(0))
	case "do_random2:12":
		return 2*c.gamma(a(0)/2) + math.Pow(c.normal()+a(1), 2)
	case "do_random2:13":
		return c.poisson(c.gamma(a(0)) * a(1) / a(0))
	case "do_random3:0":
		good, bad, draw := int(a(0)), int(a(1)), int(a(2))
		hits := 0
		for j := 0; j < draw && good+bad > 0; j++ {
			if c.random() < float64(good)/float64(good+bad) {
				hits++
				good--
			} else {
				bad--
			}
		}
		return float64(hits)
	}
	return math.NaN()
}

func kernelSample(c *Context, frame *LoweringFrame) error {
	if len(frame.Arguments) < 2 {
		return fmt.Errorf("sample expects population and size")
	}
	population, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	sizeValue, err := frameValue(c, frame, 1)
	if err != nil {
		return err
	}
	size, err := scalarInt(sizeValue)
	if err != nil || size < 0 {
		return fmt.Errorf("sample: invalid size")
	}
	data := make([]int, 0, Length(population))
	for i := 0; i < Length(population); i++ {
		data = append(data, i)
	}
	if len(data) == 0 {
		return fmt.Errorf("sample: empty population")
	}
	positions := make([]int, size)
	for i := range positions {
		j := c.RNG.Intn(len(data))
		positions[i] = data[j]
		data = append(data[:j], data[j+1:]...)
		if len(data) == 0 && i+1 < size {
			return fmt.Errorf("sample: cannot take a sample larger than population")
		}
	}
	frame.Result = takePositions(population, positions)
	return nil
}

var _ = syntax.Argument{}
