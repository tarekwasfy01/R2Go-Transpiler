package runtime

import (
	"fmt"
	"math"
)

func init() {
	for _, e := range []string{"do_intToBits", "do_numToBits", "do_rawToBits", "do_packBits", "do_rawShift"} {
		registerLoweringKernel(e, "0", kernelBitPack)
		registerLoweringKernel(e, "1", kernelBitPack)
	}
}
func kernelBitPack(c *Context, f *LoweringFrame) error {
	v, e := frameValue(c, f, 0)
	if e != nil {
		return e
	}
	switch f.Plan.CEntry {
	case "do_intToBits", "do_numToBits":
		x, e := numbers(v)
		if e != nil {
			return e
		}
		width := 32
		if f.Plan.CEntry == "do_numToBits" {
			width = 64
		}
		out := &RawVector{Data: make([]byte, len(x.Data)*width)}
		for i, n := range x.Data {
			u := uint64(uint32(int64(n)))
			if f.Plan.CEntry == "do_numToBits" {
				u = math.Float64bits(n)
			}
			for b := 0; b < width; b++ {
				out.Data[i*width+b] = byte((u >> b) & 1)
			}
		}
		f.Result = out
	case "do_rawToBits":
		x, ok := v.(*RawVector)
		if !ok {
			return fmt.Errorf("rawToBits expects raw")
		}
		out := &RawVector{Data: make([]byte, len(x.Data)*8)}
		for i, u := range x.Data {
			for b := 0; b < 8; b++ {
				out.Data[i*8+b] = (u >> b) & 1
			}
		}
		f.Result = out
	case "do_packBits":
		x, ok := v.(*RawVector)
		if !ok {
			return fmt.Errorf("packBits expects raw bits")
		}
		out := &RawVector{Data: make([]byte, (len(x.Data)+7)/8)}
		for i, b := range x.Data {
			if b != 0 {
				out.Data[i/8] |= 1 << uint(i%8)
			}
		}
		f.Result = out
	case "do_rawShift":
		x, ok := v.(*RawVector)
		if !ok {
			return fmt.Errorf("rawShift expects raw")
		}
		s, e := frameValue(c, f, 1)
		if e != nil {
			return e
		}
		shift, e := scalarInt(s)
		if e != nil {
			return e
		}
		out := &RawVector{Data: make([]byte, len(x.Data))}
		for i, b := range x.Data {
			if shift >= 0 {
				out.Data[i] = b << uint(shift)
			} else {
				out.Data[i] = b >> uint(-shift)
			}
		}
		f.Result = out
	}
	return nil
}
