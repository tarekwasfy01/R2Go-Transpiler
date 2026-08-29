package runtime

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"r2go/syntax"
	"sort"
)

var serializationMagic = []byte{'R', '2', 'G', 'O', 'S', 1}

type wireAttr struct {
	Name  string
	Value wireValue
}
type wireValue struct {
	Type       Kind
	Raw        []byte
	Logical    []int8
	Integer    []int64
	Double     []float64
	Real, Imag []float64
	Character  []string
	Missing    []bool
	List       []wireValue
	Names      []string
	Attrs      []wireAttr
}

func (c *Context) serializationBuiltin(name string, args []syntax.Argument, env *Environment) (Value, error) {
	switch name {
	case "serialize":
		if len(args) < 1 {
			return nil, fmt.Errorf("serialize expects object")
		}
		v, e := c.Eval(args[0].Value, env)
		if e != nil {
			return nil, e
		}
		w, e := toWire(v)
		if e != nil {
			return nil, e
		}
		var b bytes.Buffer
		b.Write(serializationMagic)
		e = gob.NewEncoder(&b).Encode(w)
		if e != nil {
			return nil, e
		}
		return &RawVector{Data: b.Bytes()}, nil
	case "unserialize":
		if len(args) != 1 {
			return nil, fmt.Errorf("unserialize expects raw data")
		}
		v, e := c.Eval(args[0].Value, env)
		if e != nil {
			return nil, e
		}
		raw, ok := v.(*RawVector)
		if !ok || len(raw.Data) < len(serializationMagic) || !bytes.Equal(raw.Data[:len(serializationMagic)], serializationMagic) {
			return nil, fmt.Errorf("unknown input format")
		}
		var w wireValue
		e = gob.NewDecoder(bytes.NewReader(raw.Data[len(serializationMagic):])).Decode(&w)
		if e != nil {
			return nil, e
		}
		return fromWire(w)
	}
	return nil, fmt.Errorf("unknown serialization function")
}
func attrsToWire(v Value) ([]wireAttr, error) {
	attrs := Attributes(v)
	keys := make([]string, 0, len(attrs))
	for k := range attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]wireAttr, 0, len(keys))
	for _, k := range keys {
		x, e := toWire(attrs[k])
		if e != nil {
			return nil, e
		}
		out = append(out, wireAttr{Name: k, Value: x})
	}
	return out, nil
}
func toWire(v Value) (wireValue, error) {
	w := wireValue{Type: v.Kind()}
	var e error
	switch x := v.(type) {
	case Null:
		return w, nil
	case *RawVector:
		w.Raw = append([]byte(nil), x.Data...)
	case *LogicalVector:
		w.Logical = make([]int8, len(x.Data))
		for i, n := range x.Data {
			w.Logical[i] = int8(n)
		}
	case *IntegerVector:
		w.Integer = append([]int64(nil), x.Data...)
		w.Missing = append([]bool(nil), x.Missing...)
	case *DoubleVector:
		w.Double = append([]float64(nil), x.Data...)
		w.Missing = append([]bool(nil), x.Missing...)
	case *ComplexVector:
		w.Real = make([]float64, len(x.Data))
		w.Imag = make([]float64, len(x.Data))
		for i, z := range x.Data {
			w.Real[i], w.Imag[i] = real(z), imag(z)
		}
		w.Missing = append([]bool(nil), x.Missing...)
	case *CharacterVector:
		w.Character = append([]string(nil), x.Data...)
		w.Missing = append([]bool(nil), x.Missing...)
	case *List:
		w.Names = append([]string(nil), x.Names...)
		w.List = make([]wireValue, len(x.Data))
		for i, item := range x.Data {
			w.List[i], e = toWire(item)
			if e != nil {
				return w, e
			}
		}
	default:
		return w, fmt.Errorf("serialization of %s is not implemented", v.Kind())
	}
	w.Attrs, e = attrsToWire(v)
	return w, e
}
func attrsFromWire(v Value, attrs []wireAttr) error {
	m := map[string]Value{}
	for _, a := range attrs {
		x, e := fromWire(a.Value)
		if e != nil {
			return e
		}
		m[a.Name] = x
	}
	if len(m) == 0 {
		return nil
	}
	return setAttributes(v, m)
}
func fromWire(w wireValue) (Value, error) {
	var v Value
	switch w.Type {
	case NullKind:
		v = NullValue
	case RawKind:
		v = &RawVector{Data: append([]byte(nil), w.Raw...)}
	case LogicalKind:
		x := &LogicalVector{Data: make([]Logical, len(w.Logical))}
		for i, n := range w.Logical {
			x.Data[i] = Logical(n)
		}
		v = x
	case IntegerKind:
		v = &IntegerVector{Data: append([]int64(nil), w.Integer...), Missing: append([]bool(nil), w.Missing...)}
	case DoubleKind:
		v = &DoubleVector{Data: append([]float64(nil), w.Double...), Missing: append([]bool(nil), w.Missing...)}
	case ComplexKind:
		x := &ComplexVector{Data: make([]complex128, len(w.Real)), Missing: append([]bool(nil), w.Missing...)}
		for i := range x.Data {
			x.Data[i] = complex(w.Real[i], w.Imag[i])
		}
		v = x
	case CharacterKind:
		v = &CharacterVector{Data: append([]string(nil), w.Character...), Missing: append([]bool(nil), w.Missing...)}
	case ListKind:
		x := &List{Data: make([]Value, len(w.List)), Names: append([]string(nil), w.Names...)}
		for i, item := range w.List {
			y, e := fromWire(item)
			if e != nil {
				return nil, e
			}
			x.Data[i] = y
		}
		v = x
	default:
		return nil, fmt.Errorf("unsupported serialized kind %s", w.Type)
	}
	if _, ok := v.(Null); !ok {
		if e := attrsFromWire(v, w.Attrs); e != nil {
			return nil, e
		}
	}
	return v, nil
}
