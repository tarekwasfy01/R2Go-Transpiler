package runtime

import (
	"fmt"
	"r2go/syntax"
)

func init() {
	for _, offset := range []string{"0", "1"} {
		registerLoweringKernel("do_first_min", offset, kernelFirstExtrema)
		registerLoweringKernel("do_enc2", offset, kernelEncoding)
		registerLoweringKernel("do_encoding", offset, kernelEncoding)
		registerLoweringKernel("do_str2lang", offset, kernelStringLanguage)
	}
	for _, offset := range []string{"1", "2"} {
		registerLoweringKernel("R_do_data_class", offset, kernelDataClass)
	}
}

func kernelFirstExtrema(c *Context, f *LoweringFrame) error {
	v, err := frameValue(c, f, 0)
	if err != nil {
		return err
	}
	x, err := numbers(v)
	if err != nil {
		return err
	}
	best := -1
	for i, value := range x.Data {
		if missingAt(x, i) {
			continue
		}
		if best < 0 || f.Plan.Offset == "0" && value < x.Data[best] || f.Plan.Offset == "1" && value > x.Data[best] {
			best = i
		}
	}
	if best < 0 {
		f.Result = &IntegerVector{}
	} else {
		f.Result = &IntegerVector{Data: []int64{int64(best + 1)}}
	}
	return nil
}
func kernelEncoding(c *Context, f *LoweringFrame) error {
	v, e := frameValue(c, f, 0)
	if e != nil {
		return e
	}
	x, e := characterData(v)
	if e != nil {
		return e
	}
	// R distinguishes a string's bytes from its declared encoding.  Go strings
	// are UTF-8 bytes without R's per-element encoding tag; their R-compatible
	// tag is therefore "unknown" until explicitly converted by enc2utf8.
	if f.Plan.CEntry == "do_encoding" {
		out := &CharacterVector{Data: make([]string, len(x.Data)), Missing: append([]bool(nil), x.Missing...)}
		for i := range out.Data {
			out.Data[i] = "unknown"
		}
		f.Result = out
		return nil
	}
	f.Result = &CharacterVector{Data: append([]string(nil), x.Data...), Missing: append([]bool(nil), x.Missing...)}
	return nil
}
func kernelDataClass(c *Context, f *LoweringFrame) error {
	v, e := frameValue(c, f, 0)
	if e != nil {
		return e
	}
	f.Result = &CharacterVector{Data: classNames(v)}
	return nil
}
func kernelStringLanguage(c *Context, f *LoweringFrame) error {
	v, e := frameValue(c, f, 0)
	if e != nil {
		return e
	}
	texts, e := characterData(v)
	if e != nil {
		return e
	}
	expressions := make([]Value, 0)
	for _, text := range texts.Data {
		p, e := syntax.Parse(text)
		if e != nil {
			return fmt.Errorf("%s: %w", f.Plan.Name, e)
		}
		for _, expr := range p.Expressions {
			expressions = append(expressions, &Language{Expr: expr, Text: deparseExpr(expr)})
		}
	}
	if f.Plan.Offset == "0" {
		if len(expressions) != 1 {
			return fmt.Errorf("str2lang: exactly one expression required")
		}
		f.Result = expressions[0]
	} else {
		f.Result = &List{Data: expressions, Attr: map[string]Value{"class": &CharacterVector{Data: []string{"expression"}}}}
	}
	return nil
}
