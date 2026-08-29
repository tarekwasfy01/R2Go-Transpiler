package runtime

import (
	"fmt"
	"math"
	"r2go/syntax"
	"sort"
	"strings"
)

func evalArguments(c *Context, args []syntax.Argument, env *Environment) ([]Value, map[string]Value, error) {
	vals := make([]Value, len(args))
	named := map[string]Value{}
	for i, a := range args {
		if a.Value == nil {
			return nil, nil, fmt.Errorf("argument %d is empty", i+1)
		}
		v, e := c.Eval(a.Value, env)
		if e != nil {
			return nil, nil, e
		}
		vals[i] = v
		if a.Name != "" {
			named[a.Name] = v
		}
	}
	return vals, named, nil
}

func (c *Context) moreBuiltin(name string, args []syntax.Argument, env *Environment) (Value, error) {
	vals, named, e := evalArguments(c, args, env)
	if e != nil {
		return nil, e
	}
	if len(vals) == 0 {
		return nil, fmt.Errorf("%s requires an argument", name)
	}
	switch name {
	case "rev":
		p := make([]int, Length(vals[0]))
		for i := range p {
			p[i] = len(p) - 1 - i
		}
		return takePositions(vals[0], p), nil
	case "head", "tail":
		n := 6
		if len(vals) > 1 {
			n, _ = scalarInt(vals[1])
		}
		if v, ok := named["n"]; ok {
			n, _ = scalarInt(v)
		}
		length := Length(vals[0])
		if n < 0 {
			n = length + n
		}
		if n < 0 {
			n = 0
		}
		if n > length {
			n = length
		}
		p := make([]int, n)
		for i := range p {
			if name == "head" {
				p[i] = i
			} else {
				p[i] = length - n + i
			}
		}
		return takePositions(vals[0], p), nil
	case "sort":
		decreasing := false
		if value, ok := named["decreasing"]; ok {
			decreasing = scalarLogical(value)
		}
		var naLast *bool
		if value, ok := named["na.last"]; ok {
			if logicalValue, ok := value.(*LogicalVector); ok && len(logicalValue.Data) > 0 && logicalValue.Data[0] != NA {
				placement := logicalValue.Data[0] == True
				naLast = &placement
			}
		}
		return sortValue(vals[0], decreasing, naLast)
	case "order":
		decreasing := false
		if value, ok := named["decreasing"]; ok {
			decreasing = scalarLogical(value)
		}
		naLast := true
		if value, ok := named["na.last"]; ok {
			if logicalValue, ok := value.(*LogicalVector); ok && len(logicalValue.Data) > 0 && logicalValue.Data[0] == NA {
				naLast = true
			} else {
				naLast = scalarLogical(value)
			}
		}
		return orderValueOptions(vals[0], decreasing, naLast)
	case "rank":
		ord, e := orderValue(vals[0])
		if e != nil {
			return nil, e
		}
		r := &DoubleVector{Data: make([]float64, Length(vals[0]))}
		for i, p := range ord.(*IntegerVector).Data {
			r.Data[int(p)-1] = float64(i + 1)
		}
		return r, nil
	case "cumsum", "cumprod", "cummin", "cummax":
		return cumulative(name, vals[0])
	case "diff":
		lag := 1
		if len(vals) > 1 {
			lag, _ = scalarInt(vals[1])
		}
		return difference(vals[0], lag)
	case "pmin", "pmax":
		return parallelExtrema(name, vals)
	case "nchar":
		return characterLengths(vals[0])
	case "tolower", "toupper", "trimws":
		return transformStrings(name, vals[0])
	case "substr", "substring":
		if len(vals) < 3 {
			return nil, fmt.Errorf("%s expects x, start, stop", name)
		}
		return substringValues(vals[0], vals[1], vals[2])
	case "startsWith", "endsWith":
		if len(vals) < 2 {
			return nil, fmt.Errorf("%s expects x and prefix/suffix", name)
		}
		return stringAffix(name, vals[0], vals[1])
	case "strsplit":
		if len(vals) < 2 {
			return nil, fmt.Errorf("strsplit expects x and split")
		}
		return splitStrings(vals[0], vals[1])
	}
	return nil, fmt.Errorf("unknown base builtin %s", name)
}

func sortValue(v Value, decreasing bool, naLast *bool) (Value, error) {
	placement := true
	if naLast != nil {
		placement = *naLast
	}
	p, e := orderPositionsOptions(v, decreasing, placement)
	if e != nil {
		return nil, e
	}
	// GNU R's sort default is na.last=NA, which removes missing values. A
	// concrete TRUE/FALSE retains them at the requested end.
	if naLast == nil {
		filtered := p[:0]
		for _, position := range p {
			if !valueMissingAt(v, position) {
				filtered = append(filtered, position)
			}
		}
		p = filtered
	}
	return takePositions(v, p), nil
}
func orderValue(v Value) (Value, error) {
	return orderValueOptions(v, false, true)
}

func orderValueOptions(v Value, decreasing, naLast bool) (Value, error) {
	p, e := orderPositionsOptions(v, decreasing, naLast)
	if e != nil {
		return nil, e
	}
	o := &IntegerVector{Data: make([]int64, len(p))}
	for i, n := range p {
		o.Data[i] = int64(n + 1)
	}
	return o, nil
}

func orderPositions(v Value, decreasing bool) ([]int, error) {
	return orderPositionsOptions(v, decreasing, true)
}

func orderPositionsOptions(v Value, decreasing, naLast bool) ([]int, error) {
	p := rangePositions(Length(v))
	missing := func(index int) bool { return valueMissingAt(v, index) }
	lessMissing := func(a, b int) (bool, bool) {
		am, bm := missing(a), missing(b)
		if !am && !bm {
			return false, false
		}
		if am && bm {
			return false, true
		}
		if naLast {
			return !am && bm, true
		}
		return am && !bm, true
	}
	switch x := v.(type) {
	case *CharacterVector:
		sort.SliceStable(p, func(i, j int) bool {
			if result, handled := lessMissing(p[i], p[j]); handled {
				return result
			}
			a, b := x.Data[p[i]], x.Data[p[j]]
			if decreasing {
				return a > b
			}
			return a < b
		})
	case *DoubleVector:
		sort.SliceStable(p, func(i, j int) bool {
			if result, handled := lessMissing(p[i], p[j]); handled {
				return result
			}
			a, b := x.Data[p[i]], x.Data[p[j]]
			if decreasing {
				return a > b
			}
			return a < b
		})
	case *IntegerVector:
		sort.SliceStable(p, func(i, j int) bool {
			if result, handled := lessMissing(p[i], p[j]); handled {
				return result
			}
			a, b := x.Data[p[i]], x.Data[p[j]]
			if decreasing {
				return a > b
			}
			return a < b
		})
	case *LogicalVector:
		sort.SliceStable(p, func(i, j int) bool {
			if result, handled := lessMissing(p[i], p[j]); handled {
				return result
			}
			if decreasing {
				return x.Data[p[i]] > x.Data[p[j]]
			}
			return x.Data[p[i]] < x.Data[p[j]]
		})
	default:
		return nil, fmt.Errorf("unimplemented type in sort")
	}
	return p, nil
}
func (c *Context) statisticsBuiltin(name string, args []syntax.Argument, env *Environment) (Value, error) {
	vals, named, err := evalArguments(c, args, env)
	if err != nil {
		return nil, err
	}
	if len(vals) == 0 {
		return nil, fmt.Errorf("%s requires an argument", name)
	}
	x, err := numbers(vals[0])
	if err != nil {
		return nil, err
	}
	naRm := false
	if value, ok := named["na.rm"]; ok {
		naRm = scalarLogical(value)
	}
	clean := make([]float64, 0, len(x.Data))
	for i, value := range x.Data {
		if missingAt(x, i) || math.IsNaN(value) {
			if !naRm {
				return &DoubleVector{Data: []float64{NAReal()}, Missing: []bool{true}}, nil
			}
			continue
		}
		clean = append(clean, value)
	}
	sort.Float64s(clean)
	if name == "median" {
		if len(clean) == 0 {
			return &DoubleVector{Data: []float64{NAReal()}, Missing: []bool{true}}, nil
		}
		mid := len(clean) / 2
		if len(clean)%2 == 1 {
			return &DoubleVector{Data: []float64{clean[mid]}}, nil
		}
		return &DoubleVector{Data: []float64{(clean[mid-1] + clean[mid]) / 2}}, nil
	}
	if name == "var" || name == "sd" {
		if len(clean) < 2 {
			return &DoubleVector{Data: []float64{NAReal()}, Missing: []bool{true}}, nil
		}
		mean := 0.0
		for _, value := range clean {
			mean += value
		}
		mean /= float64(len(clean))
		sumSquares := 0.0
		for _, value := range clean {
			d := value - mean
			sumSquares += d * d
		}
		result := sumSquares / float64(len(clean)-1)
		if name == "sd" {
			result = math.Sqrt(result)
		}
		return &DoubleVector{Data: []float64{result}}, nil
	}
	probs := []float64{0, .25, .5, .75, 1}
	if len(vals) > 1 && args[1].Name == "" {
		p, e := numbers(vals[1])
		if e != nil {
			return nil, e
		}
		probs = p.Data
	}
	out := &DoubleVector{Data: make([]float64, len(probs)), Missing: make([]bool, len(probs))}
	labels := make([]string, len(probs))
	for i, probability := range probs {
		if probability < 0 || probability > 1 || math.IsNaN(probability) {
			return nil, fmt.Errorf("'probs' outside [0,1]")
		}
		labels[i] = fmt.Sprintf("%g%%", probability*100)
		if len(clean) == 0 {
			out.Missing[i] = true
			out.Data[i] = NAReal()
			continue
		}
		h := float64(len(clean)-1) * probability
		lo, hi := int(math.Floor(h)), int(math.Ceil(h))
		if lo == hi {
			out.Data[i] = clean[lo]
		} else {
			out.Data[i] = clean[lo] + (h-float64(lo))*(clean[hi]-clean[lo])
		}
	}
	_ = setAttribute(out, "names", &CharacterVector{Data: labels})
	return out, nil
}

func cumulative(name string, v Value) (Value, error) {
	x, e := numbers(v)
	if e != nil {
		return nil, e
	}
	o := &DoubleVector{Data: make([]float64, len(x.Data)), Missing: make([]bool, len(x.Data))}
	acc := 0.0
	if name == "cumprod" {
		acc = 1
	}
	for i, n := range x.Data {
		if missingAt(x, i) {
			for j := i; j < len(o.Data); j++ {
				o.Data[j] = NAReal()
				o.Missing[j] = true
			}
			break
		}
		switch name {
		case "cumsum":
			acc += n
		case "cumprod":
			acc *= n
		case "cummin":
			if i == 0 || n < acc {
				acc = n
			}
		case "cummax":
			if i == 0 || n > acc {
				acc = n
			}
		}
		o.Data[i] = acc
	}
	return o, nil
}
func difference(v Value, lag int) (Value, error) {
	x, e := numbers(v)
	if e != nil {
		return nil, e
	}
	if lag < 1 {
		return nil, fmt.Errorf("'lag' must be a positive integer")
	}
	if lag >= len(x.Data) {
		return &DoubleVector{}, nil
	}
	o := &DoubleVector{Data: make([]float64, len(x.Data)-lag), Missing: make([]bool, len(x.Data)-lag)}
	for i := range o.Data {
		if missingAt(x, i) || missingAt(x, i+lag) {
			o.Missing[i] = true
			o.Data[i] = NAReal()
		} else {
			o.Data[i] = x.Data[i+lag] - x.Data[i]
		}
	}
	return o, nil
}
func parallelExtrema(name string, vals []Value) (Value, error) {
	nums := make([]*DoubleVector, len(vals))
	n := 0
	for i, v := range vals {
		x, e := numbers(v)
		if e != nil {
			return nil, e
		}
		nums[i] = x
		if len(x.Data) > n {
			n = len(x.Data)
		}
	}
	o := &DoubleVector{Data: make([]float64, n), Missing: make([]bool, n)}
	for i := 0; i < n; i++ {
		first := true
		for _, x := range nums {
			j := i % len(x.Data)
			if missingAt(x, j) {
				o.Missing[i] = true
				o.Data[i] = NAReal()
				break
			}
			if first {
				o.Data[i] = x.Data[j]
				first = false
			} else if name == "pmin" && x.Data[j] < o.Data[i] || name == "pmax" && x.Data[j] > o.Data[i] {
				o.Data[i] = x.Data[j]
			}
		}
	}
	return o, nil
}
func stringsOf(v Value) (*CharacterVector, error) {
	if x, ok := v.(*CharacterVector); ok {
		return x, nil
	}
	out := &CharacterVector{}
	for _, e := range elements(v) {
		out.Data = append(out.Data, scalarText(e))
		out.Missing = append(out.Missing, scalarMissing(e))
	}
	return out, nil
}
func characterLengths(v Value) (Value, error) {
	x, e := stringsOf(v)
	if e != nil {
		return nil, e
	}
	o := &IntegerVector{Data: make([]int64, len(x.Data)), Missing: append([]bool(nil), x.Missing...)}
	for i, s := range x.Data {
		o.Data[i] = int64(len(s))
	}
	return o, nil
}
func transformStrings(name string, v Value) (Value, error) {
	x, e := stringsOf(v)
	if e != nil {
		return nil, e
	}
	o := &CharacterVector{Data: make([]string, len(x.Data)), Missing: append([]bool(nil), x.Missing...)}
	for i, s := range x.Data {
		switch name {
		case "tolower":
			o.Data[i] = strings.ToLower(s)
		case "toupper":
			o.Data[i] = strings.ToUpper(s)
		case "trimws":
			o.Data[i] = strings.TrimSpace(s)
		}
	}
	return o, nil
}
func substringValues(v, start, stop Value) (Value, error) {
	x, _ := stringsOf(v)
	a, e := numbers(start)
	if e != nil {
		return nil, e
	}
	b, e := numbers(stop)
	if e != nil {
		return nil, e
	}
	o := &CharacterVector{Data: make([]string, len(x.Data)), Missing: append([]bool(nil), x.Missing...)}
	for i, s := range x.Data {
		r := []rune(s)
		from, to := int(a.Data[i%len(a.Data)])-1, int(b.Data[i%len(b.Data)])
		if from < 0 {
			from = 0
		}
		if to > len(r) {
			to = len(r)
		}
		if from > to || from >= len(r) {
			o.Data[i] = ""
		} else {
			o.Data[i] = string(r[from:to])
		}
	}
	return o, nil
}
func stringAffix(name string, a, b Value) (Value, error) {
	x, _ := stringsOf(a)
	y, _ := stringsOf(b)
	n := max(len(x.Data), len(y.Data))
	o := &LogicalVector{Data: make([]Logical, n)}
	for i := range o.Data {
		ok := strings.HasPrefix(x.Data[i%len(x.Data)], y.Data[i%len(y.Data)])
		if name == "endsWith" {
			ok = strings.HasSuffix(x.Data[i%len(x.Data)], y.Data[i%len(y.Data)])
		}
		if ok {
			o.Data[i] = True
		}
	}
	return o, nil
}
func splitStrings(v, sep Value) (Value, error) {
	x, _ := stringsOf(v)
	s, _ := stringsOf(sep)
	out := &List{Data: make([]Value, len(x.Data))}
	for i, text := range x.Data {
		parts := strings.Split(text, s.Data[i%len(s.Data)])
		out.Data[i] = &CharacterVector{Data: parts}
	}
	return out, nil
}
