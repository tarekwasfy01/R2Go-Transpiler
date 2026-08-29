package runtime

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

func init() {
	for _, e := range []string{"do_abbrev", "do_strtrim", "do_strtoi", "do_strrep", "do_setencoding", "do_intToUtf8", "do_utf8ToInt", "do_validUTF8", "do_encoding", "do_encodeString", "do_chartr", "do_charmatch", "do_asCharacterFactor", "do_nchar"} {
		for _, offset := range []string{"0", "1"} {
			registerLoweringKernel(e, offset, kernelBulkCharacter)
		}
	}
}
func kernelBulkCharacter(c *Context, f *LoweringFrame) error {
	v, e := frameValue(c, f, 0)
	if e != nil {
		return e
	}
	switch f.Plan.CEntry {
	case "do_utf8ToInt":
		s := []rune(scalarText(v))
		out := &IntegerVector{Data: make([]int64, len(s))}
		for i, r := range s {
			out.Data[i] = int64(r)
		}
		f.Result = out
	case "do_intToUtf8":
		x, e := numbers(v)
		if e != nil {
			return e
		}
		multiple := false
		if len(f.Arguments) > 1 {
			if mv, er := frameValue(c, f, 1); er == nil {
				if mn, er := numbers(mv); er == nil && len(mn.Data) > 0 {
					multiple = mn.Data[0] != 0
				}
			}
		}
		if multiple {
			out := &CharacterVector{Data: make([]string, len(x.Data)), Missing: append([]bool(nil), x.Missing...)}
			for i, n := range x.Data {
				out.Data[i] = string(rune(int64(n)))
			}
			f.Result = out
			break
		}
		r := make([]rune, len(x.Data))
		for i, n := range x.Data {
			r[i] = rune(int64(n))
		}
		f.Result = &CharacterVector{Data: []string{string(r)}}
	case "do_validUTF8":
		x, e := characterData(v)
		if e != nil {
			return e
		}
		o := &LogicalVector{Data: make([]Logical, len(x.Data))}
		for i, s := range x.Data {
			o.Data[i] = logicalFromBool(utf8.ValidString(s))
		}
		f.Result = o
	case "do_encoding":
		x, e := characterData(v)
		if e != nil {
			return e
		}
		o := &CharacterVector{Data: make([]string, len(x.Data)), Missing: append([]bool(nil), x.Missing...)}
		for i := range x.Data {
			o.Data[i] = "unknown"
		}
		f.Result = o
	case "do_setencoding", "do_encodeString", "do_asCharacterFactor":
		x, e := characterData(v)
		if e != nil {
			return e
		}
		f.Result = &CharacterVector{Data: append([]string(nil), x.Data...), Missing: append([]bool(nil), x.Missing...)}
	case "do_strtrim", "do_abbrev":
		x, e := characterData(v)
		if e != nil {
			return e
		}
		width := 1
		if f.Plan.CEntry == "do_abbrev" {
			width = 4
		}
		if len(f.Arguments) > 1 {
			w, er := frameValue(c, f, 1)
			if er == nil {
				width, _ = scalarInt(w)
			}
		}
		o := &CharacterVector{Data: make([]string, len(x.Data)), Missing: append([]bool(nil), x.Missing...)}
		for i, s := range x.Data {
			r := []rune(s)
			if len(r) > width {
				r = r[:width]
			}
			o.Data[i] = string(r)
		}
		f.Result = o
	case "do_nchar":
		x, e := characterData(v)
		if e != nil {
			return e
		}
		o := &IntegerVector{Data: make([]int64, len(x.Data)), Missing: append([]bool(nil), x.Missing...)}
		for i, s := range x.Data {
			o.Data[i] = int64(len(s))
		}
		f.Result = o
	case "do_strtoi":
		x, e := characterData(v)
		if e != nil {
			return e
		}
		base := 10
		if len(f.Arguments) > 1 {
			b, er := frameValue(c, f, 1)
			if er == nil {
				base, _ = scalarInt(b)
			}
		}
		o := &IntegerVector{Data: make([]int64, len(x.Data)), Missing: make([]bool, len(x.Data))}
		for i, s := range x.Data {
			n, er := strconv.ParseInt(strings.TrimSpace(s), base, 64)
			if er != nil {
				o.Missing[i] = true
			} else {
				o.Data[i] = n
			}
		}
		f.Result = o
	case "do_strrep":
		x, e := characterData(v)
		if e != nil {
			return e
		}
		times, er := frameValue(c, f, 1)
		if er != nil {
			return er
		}
		n, er := numbers(times)
		if er != nil {
			return er
		}
		o := &CharacterVector{Data: make([]string, len(x.Data))}
		for i, s := range x.Data {
			o.Data[i] = strings.Repeat(s, int(n.Data[i%len(n.Data)]))
		}
		f.Result = o
	case "do_chartr":
		old, er := frameText(c, f, 0)
		if er != nil {
			return er
		}
		newv, er := frameText(c, f, 1)
		if er != nil {
			return er
		}
		textv, er := frameValue(c, f, 2)
		if er != nil {
			return er
		}
		x, er := characterData(textv)
		if er != nil {
			return er
		}
		o := &CharacterVector{Data: make([]string, len(x.Data))}
		for i, s := range x.Data {
			o.Data[i] = strings.NewReplacer(old, newv).Replace(s)
		}
		f.Result = o
	case "do_charmatch":
		if len(f.Arguments) < 2 {
			return fmt.Errorf("charmatch expects x and table")
		}
		x, er := characterData(v)
		if er != nil {
			return er
		}
		tv, er := frameValue(c, f, 1)
		if er != nil {
			return er
		}
		t, er := characterData(tv)
		if er != nil {
			return er
		}
		o := &IntegerVector{Data: make([]int64, len(x.Data))}
		for i, s := range x.Data {
			for j, q := range t.Data {
				if s == q {
					o.Data[i] = int64(j + 1)
					break
				}
			}
		}
		f.Result = o
	}
	return nil
}
