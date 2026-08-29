package runtime

import (
	"fmt"
	"r2go/syntax"
)

func (c *Context) factorBuiltin(args []syntax.Argument, env *Environment) (Value, error) {
	if len(args) == 0 || args[0].Value == nil {
		return nil, fmt.Errorf("factor requires x")
	}
	x, e := c.Eval(args[0].Value, env)
	if e != nil {
		return nil, e
	}
	values := elements(x)
	var levels []string
	for _, a := range args[1:] {
		if a.Name == "levels" && a.Value != nil {
			v, e := c.Eval(a.Value, env)
			if e != nil {
				return nil, e
			}
			for _, item := range elements(v) {
				if !scalarMissing(item) {
					levels = append(levels, scalarText(item))
				}
			}
		}
	}
	if levels == nil {
		seen := map[string]bool{}
		for _, item := range values {
			if scalarMissing(item) {
				continue
			}
			s := scalarText(item)
			if !seen[s] {
				seen[s] = true
				levels = append(levels, s)
			}
		}
	}
	out := &IntegerVector{Data: make([]int64, len(values)), Missing: make([]bool, len(values))}
	for i, item := range values {
		if scalarMissing(item) {
			out.Missing[i] = true
			continue
		}
		s := scalarText(item)
		code := 0
		for j, l := range levels {
			if s == l {
				code = j + 1
				break
			}
		}
		if code == 0 {
			out.Missing[i] = true
		} else {
			out.Data[i] = int64(code)
		}
	}
	_ = setAttribute(out, "levels", &CharacterVector{Data: levels})
	_ = setAttribute(out, "class", &CharacterVector{Data: []string{"factor"}})
	return out, nil
}

func factorStrings(v *IntegerVector) (*CharacterVector, error) {
	l, ok := v.Attr["levels"].(*CharacterVector)
	if !ok {
		return nil, fmt.Errorf("malformed factor")
	}
	out := &CharacterVector{Data: make([]string, len(v.Data)), Missing: make([]bool, len(v.Data))}
	for i, code := range v.Data {
		if i < len(v.Missing) && v.Missing[i] {
			out.Missing[i] = true
			continue
		}
		j := int(code) - 1
		if j < 0 || j >= len(l.Data) {
			out.Missing[i] = true
		} else {
			out.Data[i] = l.Data[j]
		}
	}
	return out, nil
}
