package runtime

import (
	"fmt"
	"r2go/syntax"
)

func hasClass(v Value, want string) bool {
	c, ok := Attributes(v)["class"].(*CharacterVector)
	if !ok {
		return false
	}
	for _, s := range c.Data {
		if s == want {
			return true
		}
	}
	return false
}

func (c *Context) dataFrameBuiltin(args []syntax.Argument, env *Environment) (Value, error) {
	columns := []Value{}
	names := []string{}
	rows := 0
	for i, a := range args {
		if a.Name == "stringsAsFactors" || a.Name == "check.names" {
			continue
		}
		if a.Value == nil {
			return nil, fmt.Errorf("argument is missing, with no default")
		}
		v, e := c.Eval(a.Value, env)
		if e != nil {
			return nil, e
		}
		if _, ok := v.(Null); ok {
			continue
		}
		n := Length(v)
		if n > rows {
			rows = n
		}
		columns = append(columns, v)
		name := a.Name
		if name == "" {
			name = fmt.Sprintf("X%d", i+1)
		}
		names = append(names, name)
	}
	for i, v := range columns {
		n := Length(v)
		if rows > 0 && (n == 0 || rows%n != 0) {
			return nil, fmt.Errorf("arguments imply differing number of rows: %d, %d", rows, n)
		}
		if n != rows {
			p := make([]int, rows)
			for j := range p {
				p[j] = j % n
			}
			columns[i] = takePositions(v, p)
		}
	}
	out := &List{Data: columns, Names: names}
	_ = setAttribute(out, "class", &CharacterVector{Data: []string{"data.frame"}})
	rn := &IntegerVector{Data: make([]int64, rows)}
	for i := range rn.Data {
		rn.Data[i] = int64(i + 1)
	}
	_ = setAttribute(out, "row.names", rn)
	return out, nil
}

func dataFrameDims(v Value) (int, int, bool) {
	x, ok := v.(*List)
	if !ok || !hasClass(v, "data.frame") {
		return 0, 0, false
	}
	rows := 0
	if len(x.Data) > 0 {
		rows = Length(x.Data[0])
	}
	return rows, len(x.Data), true
}

func (c *Context) dataFrameSubset(args []syntax.Argument, env *Environment) (Value, error) {
	dfValue, e := c.Eval(args[0].Value, env)
	if e != nil {
		return nil, e
	}
	df := dfValue.(*List)
	rows, cols, _ := dataFrameDims(df)
	rowPos := rangePositions(rows)
	colPos := rangePositions(cols)
	if len(args) > 1 && args[1].Value != nil {
		v, e := c.Eval(args[1].Value, env)
		if e != nil {
			return nil, e
		}
		rowPos, e = subsetPositionsLength(rows, v)
		if e != nil {
			return nil, e
		}
	}
	if len(args) > 2 && args[2].Value != nil {
		v, e := c.Eval(args[2].Value, env)
		if e != nil {
			return nil, e
		}
		colPos, e = subsetPositions(df, v)
		if e != nil {
			return nil, e
		}
	}
	if len(args) > 2 && len(colPos) == 1 && colPos[0] >= 0 && colPos[0] < len(df.Data) {
		return takePositions(df.Data[colPos[0]], rowPos), nil
	}
	out := &List{Data: make([]Value, len(colPos)), Names: make([]string, len(colPos))}
	for i, p := range colPos {
		if p < 0 || p >= len(df.Data) {
			return nil, fmt.Errorf("undefined columns selected")
		}
		out.Data[i] = takePositions(df.Data[p], rowPos)
		out.Names[i] = df.Names[p]
	}
	_ = setAttribute(out, "class", &CharacterVector{Data: []string{"data.frame"}})
	rn := &IntegerVector{Data: make([]int64, len(rowPos))}
	for i := range rn.Data {
		rn.Data[i] = int64(i + 1)
	}
	_ = setAttribute(out, "row.names", rn)
	return out, nil
}
