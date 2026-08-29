package runtime

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// One object-family matrix supplies payload, attributes and container shape
// for every regular-expression primitive.  The surface functions only select
// a row (single/all matches, vector/list result, exact/approximate matching).
func init() {
	for _, off := range []string{"0", "1"} {
		registerLoweringKernel("do_regexpr", off, kernelMatchObjects)
		registerLoweringKernel("do_regexec", off, kernelMatchObjects)
		registerLoweringKernel("do_aregexec", off, kernelMatchObjects)
		registerLoweringKernel("do_agrep", off, kernelMatchObjects)
	}
	registerLoweringKernel("do_formatC", "0", kernelFormatCObject)
	registerLoweringKernel("do_ascall", "0", kernelAsCallObject)
}

func kernelAsCallObject(c *Context, f *LoweringFrame) error {
	v, e := frameValue(c, f, 0)
	if e != nil {
		return e
	}
	l, ok := v.(*List)
	if !ok || len(l.Data) == 0 {
		return fmt.Errorf("as.call expects a non-empty list")
	}
	parts := make([]string, len(l.Data)-1)
	for i, x := range l.Data[1:] {
		parts[i] = strings.Trim(x.String(), "\"")
	}
	name := strings.Trim(l.Data[0].String(), "\"")
	f.Result = &Language{Text: strconv.Quote(name) + "(" + strings.Join(parts, ", ") + ")"}
	return nil
}

func matchVector(pos, lengths []int64) *IntegerVector {
	v := &IntegerVector{Data: pos, Attr: map[string]Value{}}
	v.Attr["match.length"] = &IntegerVector{Data: lengths}
	v.Attr["index.type"] = &CharacterVector{Data: []string{"chars"}}
	v.Attr["useBytes"] = &LogicalVector{Data: []Logical{True}}
	return v
}

func kernelMatchObjects(c *Context, f *LoweringFrame) error {
	pattern, err := frameText(c, f, 0)
	if err != nil {
		return err
	}
	tv, err := frameValue(c, f, 1)
	if err != nil {
		return err
	}
	texts, err := characterData(tv)
	if err != nil {
		return err
	}
	if f.Plan.CEntry == "do_agrep" {
		logical := &LogicalVector{Data: make([]Logical, len(texts.Data))}
		indices := &IntegerVector{}
		limit := int(math.Max(1, math.Ceil(.1*float64(utf8.RuneCountInString(pattern)))))
		for i, text := range texts.Data {
			ok := strings.Contains(text, pattern) || editDistance([]rune(pattern), []rune(text)) <= limit
			logical.Data[i] = logicalFromBool(ok)
			if ok {
				indices.Data = append(indices.Data, int64(i+1))
			}
		}
		if f.Plan.Offset == "1" {
			f.Result = logical
		} else {
			f.Result = indices
		}
		return nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	result := &List{Data: make([]Value, len(texts.Data))}
	for i, text := range texts.Data {
		all := f.Plan.CEntry == "do_regexpr" && f.Plan.Offset == "1"
		matches := re.FindAllStringIndex(text, -1)
		if !all && len(matches) > 1 {
			matches = matches[:1]
		}
		pos, lens := []int64{}, []int64{}
		for _, m := range matches {
			pos = append(pos, int64(utf8.RuneCountInString(text[:m[0]])+1))
			lens = append(lens, int64(utf8.RuneCountInString(text[m[0]:m[1]])))
		}
		if len(pos) == 0 {
			pos = []int64{-1}
			lens = []int64{-1}
		}
		result.Data[i] = matchVector(pos, lens)
		if f.Plan.CEntry == "do_aregexec" {
			delete(result.Data[i].(*IntegerVector).Attr, "index.type")
			delete(result.Data[i].(*IntegerVector).Attr, "useBytes")
		}
	}
	if f.Plan.CEntry == "do_regexpr" && f.Plan.Offset == "0" {
		if len(result.Data) == 1 {
			f.Result = result.Data[0]
			return nil
		}
		pos, lens := make([]int64, len(result.Data)), make([]int64, len(result.Data))
		for i, v := range result.Data {
			q := v.(*IntegerVector)
			pos[i] = q.Data[0]
			lens[i] = q.Attr["match.length"].(*IntegerVector).Data[0]
		}
		f.Result = matchVector(pos, lens)
		return nil
	}
	f.Result = result
	return nil
}

func editDistance(a, b []rune) int {
	prev := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		cur := make([]int, len(b)+1)
		cur[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			cur[j] = min3(cur[j-1]+1, prev[j]+1, prev[j-1]+cost)
		}
		prev = cur
	}
	return prev[len(b)]
}
func min3(a, b, c int) int {
	if b < a {
		a = b
	}
	if c < a {
		a = c
	}
	return a
}

func kernelFormatCObject(c *Context, f *LoweringFrame) error {
	v, err := frameValue(c, f, 0)
	if err != nil {
		return err
	}
	x, err := numbers(v)
	if err != nil {
		return err
	}
	digits := 2
	if d, ok, e := frameNamed(c, f, "digits"); e == nil && ok {
		if n, e := scalarInt(d); e == nil {
			digits = n
		}
	}
	out := &CharacterVector{Data: make([]string, len(x.Data)), Missing: append([]bool(nil), x.Missing...)}
	for i, n := range x.Data {
		out.Data[i] = fmt.Sprintf("%*.*g", digits+1, digits, n)
	}
	f.Result = out
	return nil
}
