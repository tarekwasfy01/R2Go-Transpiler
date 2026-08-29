package runtime

import (
	"fmt"
	"time"
)

func init() {
	for _, e := range []string{"do_D2POSIXlt", "do_POSIXlt2D", "do_asPOSIXct", "do_formatPOSIXlt", "do_balancePOSIXlt"} {
		registerLoweringKernel(e, "0", kernelDateTime)
	}
	registerLoweringKernel("do_asPOSIXlt", "0", kernelDateTime)
	registerLoweringKernel("do_balancePOSIXlt", "1", kernelDateTime)
}
func kernelDateTime(c *Context, f *LoweringFrame) error {
	if len(f.Arguments) == 0 {
		return fmt.Errorf("%s: missing time value", f.Plan.Name)
	}
	v, e := frameValue(c, f, 0)
	if e != nil {
		return e
	}
	switch f.Plan.CEntry {
	case "do_asPOSIXct", "do_POSIXlt2D":
		if l, ok := v.(*List); ok && len(l.Data) > 0 {
			sec, _ := numbers(l.Data[0])
			min, _ := numbers(l.Data[1])
			hour, _ := numbers(l.Data[2])
			day, _ := numbers(l.Data[3])
			mon, _ := numbers(l.Data[4])
			year, _ := numbers(l.Data[5])
			out := &DoubleVector{Data: make([]float64, len(year.Data)), Attr: map[string]Value{"class": &CharacterVector{Data: []string{"POSIXct", "POSIXt"}}}}
			for i := range out.Data {
				out.Data[i] = float64(time.Date(int(year.Data[i])+1900, time.Month(int(mon.Data[i])+1), int(day.Data[i]), int(hour.Data[i]), int(min.Data[i]), int(sec.Data[i]), 0, time.Local).Unix())
			}
			f.Result = out
		} else if text, ok := v.(*CharacterVector); ok {
			out := &DoubleVector{Data: make([]float64, len(text.Data)), Attr: map[string]Value{"class": &CharacterVector{Data: []string{"POSIXct", "POSIXt"}}}}
			for i, s := range text.Data {
				t, er := time.ParseInLocation("2006-01-02", s, time.Local)
				if er != nil {
					return er
				}
				out.Data[i] = float64(t.Unix())
			}
			f.Result = out
		} else {
			f.Result = v
		}
	case "do_D2POSIXlt", "do_asPOSIXlt", "do_balancePOSIXlt":
		if l, ok := v.(*List); ok {
			// balancePOSIXlt receives the POSIXlt list it is meant to normalize.
			// A complete calendar object is already normalized by this Pure-Go
			// representation, so preserve it rather than coercing it to numbers.
			if l.Attr == nil {
				l.Attr = map[string]Value{}
			}
			l.Attr["class"] = &CharacterVector{Data: []string{"POSIXlt", "POSIXt"}}
			f.Result = l
			return nil
		}
		if text, ok := v.(*CharacterVector); ok {
			out := &List{Names: []string{"sec", "min", "hour", "mday", "mon", "year"}, Data: make([]Value, 6)}
			fields := make([][]float64, 6)
			for _, item := range text.Data {
				t, parseErr := time.ParseInLocation("2006-01-02", item, time.Local)
				if parseErr != nil {
					return parseErr
				}
				parts := []float64{float64(t.Second()), float64(t.Minute()), float64(t.Hour()), float64(t.Day()), float64(t.Month() - 1), float64(t.Year() - 1900)}
				for i := range fields {
					fields[i] = append(fields[i], parts[i])
				}
			}
			for i := range fields {
				out.Data[i] = &DoubleVector{Data: fields[i]}
			}
			out.Attr = map[string]Value{"class": &CharacterVector{Data: []string{"POSIXlt", "POSIXt"}}}
			f.Result = out
			return nil
		}
		n, e := numbers(v)
		if e != nil {
			return e
		}
		out := &List{Names: []string{"sec", "min", "hour", "mday", "mon", "year"}, Data: make([]Value, 6)}
		fields := make([][]float64, 6)
		for i, x := range n.Data {
			t := time.Unix(int64(x), 0).Local()
			parts := []float64{float64(t.Second()), float64(t.Minute()), float64(t.Hour()), float64(t.Day()), float64(t.Month() - 1), float64(t.Year() - 1900)}
			for j := range fields {
				fields[j] = append(fields[j], parts[j])
			}
			_ = i
		}
		for i := range fields {
			out.Data[i] = &DoubleVector{Data: fields[i]}
		}
		out.Attr = map[string]Value{"class": &CharacterVector{Data: []string{"POSIXlt", "POSIXt"}}}
		f.Result = out
	case "do_formatPOSIXlt":
		if l, ok := v.(*List); ok && len(l.Data) > 0 {
			f.Result = &CharacterVector{Data: []string{l.Data[0].String()}}
		} else {
			f.Result = &CharacterVector{Data: []string{v.String()}}
		}
	}
	return nil
}
