package runtime

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"r2go/syntax"
	"regexp"
	"strconv"
	"strings"
)

func (c *Context) ioBuiltin(name string, args []syntax.Argument, env *Environment) (Value, error) {
	vals, named, e := evalArguments(c, args, env)
	if e != nil {
		return nil, e
	}
	switch name {
	case "getwd":
		wd, e := c.Host.Getwd()
		if e != nil {
			return nil, e
		}
		return &CharacterVector{Data: []string{wd}}, nil
	case "setwd":
		if len(vals) != 1 {
			return nil, fmt.Errorf("setwd expects dir")
		}
		old, e := c.Host.Getwd()
		if e != nil {
			return nil, e
		}
		if e = c.Host.Chdir(firstString(vals[0])); e != nil {
			return nil, e
		}
		return &CharacterVector{Data: []string{old}}, nil
	case "file.exists", "dir.exists":
		if len(vals) < 1 {
			return nil, fmt.Errorf("%s expects paths", name)
		}
		x, _ := stringsOf(vals[0])
		out := &LogicalVector{Data: make([]Logical, len(x.Data))}
		for i, p := range x.Data {
			exists, isDir := c.Host.Exists(p)
			ok := exists
			if name == "dir.exists" {
				ok = ok && isDir
			}
			if ok {
				out.Data[i] = True
			}
		}
		return out, nil
	case "dir.create":
		if len(vals) < 1 {
			return nil, fmt.Errorf("dir.create expects path")
		}
		recursive := false
		if v, ok := named["recursive"]; ok {
			recursive, _ = IsTrue(v)
		}
		err := c.Host.Mkdir(firstString(vals[0]), recursive)
		return boolValue(err == nil), nil
	case "basename", "dirname", "normalizePath":
		x, _ := stringsOf(vals[0])
		out := &CharacterVector{Data: make([]string, len(x.Data))}
		for i, p := range x.Data {
			switch name {
			case "basename":
				out.Data[i] = filepath.Base(p)
			case "dirname":
				out.Data[i] = filepath.Dir(p)
			case "normalizePath":
				q, e := filepath.Abs(p)
				if e != nil {
					return nil, e
				}
				out.Data[i] = filepath.Clean(q)
			}
		}
		return out, nil
	case "readLines":
		return readLinesValue(c, vals, named)
	case "writeLines":
		return writeLinesValue(c, vals)
	case "list.files":
		return listFilesValue(vals, named)
	case "read.csv", "read.csv2":
		return readCSVValue(c, name, vals, named)
	case "write.csv", "write.csv2":
		return writeCSVValue(c, name, vals, named)
	}
	return nil, fmt.Errorf("unknown I/O builtin %s", name)
}
func firstString(v Value) string {
	x, _ := stringsOf(v)
	if len(x.Data) == 0 {
		return ""
	}
	return x.Data[0]
}
func readLinesValue(c *Context, vals []Value, named map[string]Value) (Value, error) {
	if len(vals) < 1 {
		return nil, fmt.Errorf("readLines expects con")
	}
	b, e := c.Host.ReadFile(firstString(vals[0]))
	if e != nil {
		return nil, e
	}
	limit := -1
	if v, ok := named["n"]; ok {
		limit, _ = scalarInt(v)
	}
	scan := bufio.NewScanner(strings.NewReader(string(b)))
	out := &CharacterVector{}
	for scan.Scan() {
		if limit >= 0 && len(out.Data) >= limit {
			break
		}
		out.Data = append(out.Data, scan.Text())
	}
	return out, scan.Err()
}
func writeLinesValue(c *Context, vals []Value) (Value, error) {
	if len(vals) < 2 {
		return nil, fmt.Errorf("writeLines expects text and con")
	}
	x, _ := stringsOf(vals[0])
	var text strings.Builder
	for _, s := range x.Data {
		text.WriteString(s + "\n")
	}
	return NullValue, c.Host.WriteFile(firstString(vals[1]), []byte(text.String()))
}
func listFilesValue(vals []Value, named map[string]Value) (Value, error) {
	path := "."
	if len(vals) > 0 {
		path = firstString(vals[0])
	}
	pattern := ""
	if v, ok := named["pattern"]; ok {
		pattern = firstString(v)
	}
	var re *regexp.Regexp
	var e error
	if pattern != "" {
		re, e = regexp.Compile(pattern)
		if e != nil {
			return nil, e
		}
	}
	recursive := false
	if v, ok := named["recursive"]; ok {
		recursive, _ = IsTrue(v)
	}
	full := false
	if v, ok := named["full.names"]; ok {
		full, _ = IsTrue(v)
	}
	out := &CharacterVector{}
	if recursive {
		e = filepath.Walk(path, func(p string, info os.FileInfo, e error) error {
			if e != nil {
				return e
			}
			if p == path {
				return nil
			}
			rel, _ := filepath.Rel(path, p)
			if re == nil || re.MatchString(info.Name()) {
				if full {
					out.Data = append(out.Data, p)
				} else {
					out.Data = append(out.Data, rel)
				}
			}
			return nil
		})
	} else {
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if re == nil || re.MatchString(entry.Name()) {
				if full {
					out.Data = append(out.Data, filepath.Join(path, entry.Name()))
				} else {
					out.Data = append(out.Data, entry.Name())
				}
			}
		}
	}
	return out, e
}
func readCSVValue(c *Context, name string, vals []Value, named map[string]Value) (Value, error) {
	if len(vals) < 1 {
		return nil, fmt.Errorf("read.csv expects file")
	}
	b, e := c.Host.ReadFile(firstString(vals[0]))
	if e != nil {
		return nil, e
	}
	r := csv.NewReader(strings.NewReader(string(b)))
	if name == "read.csv2" {
		r.Comma = ';'
	}
	records, e := r.ReadAll()
	if e != nil {
		return nil, e
	}
	if len(records) == 0 {
		return &List{Attr: map[string]Value{"class": &CharacterVector{Data: []string{"data.frame"}}}}, nil
	}
	header := records[0]
	rows := records[1:]
	cols := make([]Value, len(header))
	for j := range header {
		raw := make([]string, len(rows))
		for i := range rows {
			if j < len(rows[i]) {
				raw[i] = rows[i][j]
			}
		}
		cols[j] = inferColumn(raw)
	}
	out := &List{Data: cols, Names: header}
	_ = setAttribute(out, "class", &CharacterVector{Data: []string{"data.frame"}})
	rn := &IntegerVector{Data: make([]int64, len(rows))}
	for i := range rn.Data {
		rn.Data[i] = int64(i + 1)
	}
	_ = setAttribute(out, "row.names", rn)
	return out, nil
}
func inferColumn(raw []string) Value {
	nums := make([]float64, len(raw))
	allNum := true
	for i, s := range raw {
		n, e := strconv.ParseFloat(s, 64)
		if e != nil {
			allNum = false
			break
		}
		nums[i] = n
	}
	if allNum {
		return &DoubleVector{Data: nums}
	}
	logicals := make([]Logical, len(raw))
	allLog := true
	for i, s := range raw {
		switch strings.ToUpper(s) {
		case "TRUE", "T":
			logicals[i] = True
		case "FALSE", "F":
		default:
			allLog = false
		}
		if !allLog {
			break
		}
	}
	if allLog {
		return &LogicalVector{Data: logicals}
	}
	return &CharacterVector{Data: raw}
}
func writeCSVValue(c *Context, name string, vals []Value, named map[string]Value) (Value, error) {
	if len(vals) < 2 {
		return nil, fmt.Errorf("write.csv expects x and file")
	}
	df, ok := vals[0].(*List)
	if !ok || !hasClass(df, "data.frame") {
		return nil, fmt.Errorf("x must be a data frame")
	}
	var text strings.Builder
	w := csv.NewWriter(&text)
	var err error
	if name == "write.csv2" {
		w.Comma = ';'
	}
	if err = w.Write(df.Names); err != nil {
		return nil, err
	}
	rows, _, _ := dataFrameDims(df)
	for i := 0; i < rows; i++ {
		row := make([]string, len(df.Data))
		for j, col := range df.Data {
			row[j] = scalarText(elementAt(col, i))
		}
		if err = w.Write(row); err != nil {
			return nil, err
		}
	}
	w.Flush()
	if err = w.Error(); err != nil {
		return nil, err
	}
	return NullValue, c.Host.WriteFile(firstString(vals[1]), []byte(text.String()))
}
