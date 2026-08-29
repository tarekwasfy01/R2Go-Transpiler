package rgo

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"math/cmplx"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
)

type Kind uint8

const (
	Null Kind = iota
	Logical
	Integer
	Double
	Complex
	String
	Raw
	List
	Symbol
	Environment
	Function
	Error
	Connection
)

type NAState uint8

const (
	Present NAState = iota
	NA
)

type Value struct {
	Kind   Kind
	L      []int8
	I      []int64
	D      []float64
	Z      []complex128
	S      []string
	B      []byte
	V      []Value
	NA     []bool
	Attr   map[string]Value
	Name   string
	Env    *Env
	Fn     func([]Value, *Env) (Value, error)
	Err    error
	Conn   *Conn
	shared bool
}

type Env struct {
	mu     sync.RWMutex
	Parent *Env
	Name   string
	Values map[string]Value
	Locked bool
	Active map[string]func() Value
}
type Conn struct {
	mu       sync.Mutex
	R        io.Reader
	W        io.Writer
	C        io.Closer
	Name     string
	Open     bool
	Text     bool
	Pushback []string
}

var Nil = Value{Kind: Null}

func ErrValue(format string, a ...any) Value {
	return Value{Kind: Error, Err: fmt.Errorf(format, a...)}
}
func isErr(v Value) bool { return v.Kind == Error }
func Bool(x bool) Value {
	if x {
		return Value{Kind: Logical, L: []int8{1}}
	}
	return Value{Kind: Logical, L: []int8{0}}
}
func Ints(x ...int64) Value      { return Value{Kind: Integer, I: append([]int64(nil), x...)} }
func Doubles(x ...float64) Value { return Value{Kind: Double, D: append([]float64(nil), x...)} }
func Strings(x ...string) Value  { return Value{Kind: String, S: append([]string(nil), x...)} }
func Raws(x ...byte) Value       { return Value{Kind: Raw, B: append([]byte(nil), x...)} }
func Lists(x ...Value) Value     { return Value{Kind: List, V: append([]Value(nil), x...)} }
func Sym(s string) Value         { return Value{Kind: Symbol, Name: s} }
func NewEnv(parent *Env, name string) Value {
	return Value{Kind: Environment, Env: &Env{Parent: parent, Name: name, Values: map[string]Value{}, Active: map[string]func() Value{}}}
}

func clone(v Value) Value {
	v.shared = false
	v.L = append([]int8(nil), v.L...)
	v.I = append([]int64(nil), v.I...)
	v.D = append([]float64(nil), v.D...)
	v.Z = append([]complex128(nil), v.Z...)
	v.S = append([]string(nil), v.S...)
	v.B = append([]byte(nil), v.B...)
	v.V = append([]Value(nil), v.V...)
	v.NA = append([]bool(nil), v.NA...)
	if v.Attr != nil {
		m := map[string]Value{}
		for k, x := range v.Attr {
			m[k] = x
		}
		v.Attr = m
	}
	return v
}
func writable(v Value) Value {
	if v.shared {
		return clone(v)
	}
	return v
}
func withAttr(v Value, k string, x Value) Value {
	v = writable(v)
	if v.Attr == nil {
		v.Attr = map[string]Value{}
	}
	v.Attr[k] = x
	return v
}
func attr(v Value, k string) Value {
	if v.Attr == nil {
		return Nil
	}
	return v.Attr[k]
}
func copyAttrs(dst, src Value) Value {
	dst = writable(dst)
	if src.Attr != nil {
		dst.Attr = map[string]Value{}
		for k, v := range src.Attr {
			dst.Attr[k] = v
		}
	}
	return dst
}
func length(v Value) int {
	switch v.Kind {
	case Logical:
		return len(v.L)
	case Integer:
		return len(v.I)
	case Double:
		return len(v.D)
	case Complex:
		return len(v.Z)
	case String:
		return len(v.S)
	case Raw:
		return len(v.B)
	case List:
		return len(v.V)
	case Null:
		return 0
	default:
		return 1
	}
}
func arg(args Value, i int) Value {
	if args.Kind != List || i < 0 || i >= len(args.V) {
		return Nil
	}
	return args.V[i]
}
func nargs(args Value) int {
	if args.Kind == List {
		return len(args.V)
	}
	if args.Kind == Null {
		return 0
	}
	return 1
}
func asBool(v Value) (bool, error) {
	if len(v.NA) > 0 && v.NA[0] {
		return false, errors.New("missing value where TRUE/FALSE needed")
	}
	switch v.Kind {
	case Logical:
		if len(v.L) == 0 {
			return false, errors.New("length zero")
		}
		return v.L[0] != 0, nil
	case Integer:
		if len(v.I) == 0 {
			return false, errors.New("length zero")
		}
		return v.I[0] != 0, nil
	case Double:
		if len(v.D) == 0 {
			return false, errors.New("length zero")
		}
		return v.D[0] != 0, nil
	case String:
		if len(v.S) == 0 {
			return false, errors.New("length zero")
		}
		b, e := strconv.ParseBool(v.S[0])
		return b, e
	}
	return false, fmt.Errorf("cannot coerce %v to logical", v.Kind)
}
func asInt(v Value) (int64, error) {
	switch v.Kind {
	case Integer:
		if len(v.I) > 0 {
			return v.I[0], nil
		}
	case Double:
		if len(v.D) > 0 {
			return int64(v.D[0]), nil
		}
	case Logical:
		if len(v.L) > 0 {
			return int64(v.L[0]), nil
		}
	case String:
		if len(v.S) > 0 {
			return strconv.ParseInt(v.S[0], 10, 64)
		}
	}
	return 0, errors.New("cannot coerce to integer")
}
func asFloat(v Value) (float64, error) {
	switch v.Kind {
	case Double:
		if len(v.D) > 0 {
			return v.D[0], nil
		}
	case Integer:
		if len(v.I) > 0 {
			return float64(v.I[0]), nil
		}
	case Logical:
		if len(v.L) > 0 {
			return float64(v.L[0]), nil
		}
	case String:
		if len(v.S) > 0 {
			return strconv.ParseFloat(v.S[0], 64)
		}
	}
	return math.NaN(), errors.New("cannot coerce to double")
}
func asString(v Value) (string, error) {
	switch v.Kind {
	case String:
		if len(v.S) > 0 {
			return v.S[0], nil
		}
	case Symbol:
		return v.Name, nil
	case Integer:
		if len(v.I) > 0 {
			return strconv.FormatInt(v.I[0], 10), nil
		}
	case Double:
		if len(v.D) > 0 {
			return strconv.FormatFloat(v.D[0], 'g', -1, 64), nil
		}
	case Logical:
		if len(v.L) > 0 {
			if v.L[0] != 0 {
				return "TRUE", nil
			}
			return "FALSE", nil
		}
	}
	return "", errors.New("cannot coerce to string")
}
func names(v Value) []string {
	a := attr(v, "names")
	if a.Kind == String {
		return a.S
	}
	return nil
}
func setNames(v Value, n []string) Value { return withAttr(v, "names", Strings(n...)) }
func preserveShape(dst, src Value) Value {
	for _, k := range []string{"names", "dim", "dimnames"} {
		if x := attr(src, k); x.Kind != Null {
			dst = withAttr(dst, k, x)
		}
	}
	return dst
}

func envGet(e *Env, k string) (Value, bool) {
	for e != nil {
		e.mu.RLock()
		if f, ok := e.Active[k]; ok {
			e.mu.RUnlock()
			return f(), true
		}
		v, ok := e.Values[k]
		e.mu.RUnlock()
		if ok {
			return v, true
		}
		e = e.Parent
	}
	return Nil, false
}
func envSet(e *Env, k string, v Value) error {
	if e == nil {
		return errors.New("NULL environment")
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.Locked {
		return errors.New("cannot add bindings to a locked environment")
	}
	e.Values[k] = v
	return nil
}
func envRemove(e *Env, k string) error {
	if e == nil {
		return errors.New("NULL environment")
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.Locked {
		return errors.New("cannot remove bindings from a locked environment")
	}
	delete(e.Values, k)
	return nil
}

func recycleLen(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func naAt(v Value, i int) bool { return len(v.NA) > 0 && v.NA[i%len(v.NA)] }
func vectorFloat(v Value, i int) (float64, bool) {
	if length(v) == 0 {
		return 0, true
	}
	if naAt(v, i) {
		return 0, true
	}
	switch v.Kind {
	case Double:
		return v.D[i%len(v.D)], false
	case Integer:
		return float64(v.I[i%len(v.I)]), false
	case Logical:
		return float64(v.L[i%len(v.L)]), false
	}
	return math.NaN(), false
}
func vectorString(v Value, i int) (string, bool) {
	if length(v) == 0 {
		return "", true
	}
	if naAt(v, i) {
		return "", true
	}
	switch v.Kind {
	case String:
		return v.S[i%len(v.S)], false
	case Symbol:
		return v.Name, false
	}
	s, e := asString(v)
	return s, e != nil
}

func bincodeImpl(x, br, right, lowest Value) Value {
	r, _ := asBool(right)
	low, _ := asBool(lowest)
	if br.Kind != Double && br.Kind != Integer {
		return ErrValue("'breaks' must be numeric")
	}
	nb := length(br)
	if nb < 2 {
		return ErrValue("invalid 'breaks' argument")
	}
	bs := make([]float64, nb)
	for i := range bs {
		bs[i], _ = vectorFloat(br, i)
	}
	n := length(x)
	out := Value{Kind: Integer, I: make([]int64, n), NA: make([]bool, n)}
	for i := 0; i < n; i++ {
		z, na := vectorFloat(x, i)
		if na || math.IsNaN(z) {
			out.NA[i] = true
			continue
		}
		j := sort.Search(len(bs), func(k int) bool {
			if r {
				return z <= bs[k]
			}
			return z < bs[k]
		})
		code := j
		if !r {
			code = j
		}
		if j == 0 {
			if low && z == bs[0] {
				code = 1
			} else {
				out.NA[i] = true
				continue
			}
		}
		if j >= nb {
			if low && z == bs[nb-1] {
				code = nb - 1
			} else {
				out.NA[i] = true
				continue
			}
		}
		if r {
			code = j
		}
		out.I[i] = int64(code)
	}
	return out
}
func tabulateImpl(bin, nb Value) Value {
	n, err := asInt(nb)
	if err != nil || n < 0 {
		return ErrValue("invalid 'nbins'")
	}
	out := Ints(make([]int64, int(n))...)
	for i := 0; i < length(bin); i++ {
		x, err := asIntAt(bin, i)
		if err == nil && x > 0 && x <= n {
			out.I[x-1]++
		}
	}
	return out
}
func asIntAt(v Value, i int) (int64, error) {
	if naAt(v, i) {
		return 0, errors.New("NA")
	}
	switch v.Kind {
	case Integer:
		return v.I[i%len(v.I)], nil
	case Double:
		return int64(v.D[i%len(v.D)]), nil
	case Logical:
		return int64(v.L[i%len(v.L)]), nil
	}
	return 0, errors.New("not integer")
}
func findIntervalImpl(x, vec, leftOpen, rightmost, allInside Value) Value {
	lo, _ := asBool(leftOpen)
	rm, _ := asBool(rightmost)
	ai, _ := asBool(allInside)
	m := length(vec)
	br := make([]float64, m)
	for i := range br {
		br[i], _ = vectorFloat(vec, i)
	}
	out := Value{Kind: Integer, I: make([]int64, length(x)), NA: make([]bool, length(x))}
	for i := 0; i < length(x); i++ {
		z, na := vectorFloat(x, i)
		if na || math.IsNaN(z) {
			out.NA[i] = true
			continue
		}
		var j int
		if lo {
			j = sort.Search(m, func(k int) bool { return br[k] >= z })
		} else {
			j = sort.Search(m, func(k int) bool { return br[k] > z })
		}
		if rm && z == br[m-1] {
			j = m - 1
		}
		if ai {
			if j < 1 {
				j = 1
			}
			if j >= m {
				j = m - 1
			}
		}
		out.I[i] = int64(j)
	}
	return out
}
func lengthGetsImpl(x, nv Value) Value {
	n, err := asInt(nv)
	if err != nil || n < 0 {
		return ErrValue("invalid value")
	}
	N := int(n)
	out := clone(x)
	switch x.Kind {
	case Logical:
		out.L = resizeInt8(x.L, N)
	case Integer:
		out.I = resizeInt64(x.I, N)
	case Double:
		out.D = resizeF64(x.D, N)
	case Complex:
		out.Z = resizeC128(x.Z, N)
	case String:
		out.S = resizeString(x.S, N)
	case Raw:
		out.B = resizeByte(x.B, N)
	case List:
		out.V = resizeValue(x.V, N)
	default:
		return ErrValue("invalid value")
	}
	old := length(x)
	if N > old && x.Kind != Raw {
		out.NA = make([]bool, N)
		copy(out.NA, x.NA)
		for i := old; i < N; i++ {
			out.NA[i] = true
		}
	}
	if nm := names(x); nm != nil {
		nn := make([]string, N)
		copy(nn, nm)
		out = setNames(out, nn)
	}
	if out.Attr != nil {
		delete(out.Attr, "dim")
		delete(out.Attr, "dimnames")
	}
	return out
}
func resizeInt8(a []int8, n int) []int8             { b := make([]int8, n); copy(b, a); return b }
func resizeInt64(a []int64, n int) []int64          { b := make([]int64, n); copy(b, a); return b }
func resizeF64(a []float64, n int) []float64        { b := make([]float64, n); copy(b, a); return b }
func resizeC128(a []complex128, n int) []complex128 { b := make([]complex128, n); copy(b, a); return b }
func resizeString(a []string, n int) []string       { b := make([]string, n); copy(b, a); return b }
func resizeByte(a []byte, n int) []byte             { b := make([]byte, n); copy(b, a); return b }
func resizeValue(a []Value, n int) []Value          { b := make([]Value, n); copy(b, a); return b }

func intToBitsImpl(x Value) Value {
	n := length(x)
	out := Value{Kind: Raw, B: make([]byte, n*32)}
	for i := 0; i < n; i++ {
		v, _ := asIntAt(x, i)
		u := uint32(v)
		for b := 0; b < 32; b++ {
			out.B[i*32+b] = byte((u >> b) & 1)
		}
	}
	return out
}
func numToBitsImpl(x Value) Value {
	n := length(x)
	out := Value{Kind: Raw, B: make([]byte, n*64)}
	for i := 0; i < n; i++ {
		f, _ := vectorFloat(x, i)
		u := math.Float64bits(f)
		for b := 0; b < 64; b++ {
			out.B[i*64+b] = byte((u >> b) & 1)
		}
	}
	return out
}
func numToIntsImpl(x Value) Value {
	n := length(x)
	out := Value{Kind: Integer, I: make([]int64, n*2)}
	for i := 0; i < n; i++ {
		f, _ := vectorFloat(x, i)
		u := math.Float64bits(f)
		out.I[2*i] = int64(int32(u))
		out.I[2*i+1] = int64(int32(u >> 32))
	}
	return out
}
func rawToBitsImpl(x Value) Value {
	if x.Kind != Raw {
		return ErrValue("argument must be raw")
	}
	o := Value{Kind: Raw, B: make([]byte, len(x.B)*8)}
	for i, v := range x.B {
		for b := 0; b < 8; b++ {
			o.B[i*8+b] = (v >> b) & 1
		}
	}
	return o
}
func rawShiftImpl(x, n Value) Value {
	if x.Kind != Raw {
		return ErrValue("argument must be raw")
	}
	s, e := asInt(n)
	if e != nil || s < -8 || s > 8 {
		return ErrValue("shift must be between -8 and 8")
	}
	o := Value{Kind: Raw, B: make([]byte, len(x.B))}
	for i, v := range x.B {
		if s >= 0 {
			o.B[i] = byte(uint16(v) << uint(s))
		} else {
			o.B[i] = byte(uint16(v) >> uint(-s))
		}
	}
	return o
}
func packBitsImpl(x, typ Value) Value {
	t, _ := asString(typ)
	n := length(x)
	bits := make([]byte, n)
	for i := 0; i < n; i++ {
		v, _ := asIntAt(x, i)
		bits[i] = byte(v & 1)
	}
	switch t {
	case "raw":
		if n%8 != 0 {
			return ErrValue("must be a multiple of 8")
		}
		o := make([]byte, n/8)
		for i, b := range bits {
			o[i/8] |= b << uint(i%8)
		}
		return Raws(o...)
	case "integer":
		if n%32 != 0 {
			return ErrValue("must be a multiple of 32")
		}
		o := make([]int64, n/32)
		for i, b := range bits {
			o[i/32] |= int64(b) << uint(i%32)
		}
		return Ints(o...)
	case "double":
		if n%64 != 0 {
			return ErrValue("must be a multiple of 64")
		}
		o := make([]float64, n/64)
		for j := range o {
			var u uint64
			for i := 0; i < 64; i++ {
				u |= uint64(bits[j*64+i]) << uint(i)
			}
			o[j] = math.Float64frombits(u)
		}
		return Doubles(o...)
	}
	return ErrValue("invalid 'type'")
}

func splitImpl(x, f Value) Value {
	if f.Kind != Integer && f.Kind != String {
		return ErrValue("grouping factor must be integer or string")
	}
	n := length(x)
	if n == 0 {
		return Lists()
	}
	groups := map[string][]int{}
	order := []string{}
	for i := 0; i < n; i++ {
		var k string
		if f.Kind == Integer {
			v, _ := asIntAt(f, i)
			k = strconv.FormatInt(v, 10)
		} else {
			k, _ = vectorString(f, i)
		}
		if _, ok := groups[k]; !ok {
			order = append(order, k)
		}
		groups[k] = append(groups[k], i)
	}
	out := make([]Value, len(order))
	for j, k := range order {
		out[j] = subsetIndices(x, groups[k])
	}
	return setNames(Lists(out...), order)
}
func subsetIndices(x Value, idx []int) Value {
	o := Value{Kind: x.Kind, Attr: map[string]Value{}}
	switch x.Kind {
	case Integer:
		o.I = make([]int64, len(idx))
		for j, i := range idx {
			o.I[j] = x.I[i%len(x.I)]
		}
	case Double:
		o.D = make([]float64, len(idx))
		for j, i := range idx {
			o.D[j] = x.D[i%len(x.D)]
		}
	case Logical:
		o.L = make([]int8, len(idx))
		for j, i := range idx {
			o.L[j] = x.L[i%len(x.L)]
		}
	case String:
		o.S = make([]string, len(idx))
		for j, i := range idx {
			o.S[j] = x.S[i%len(x.S)]
		}
	case Raw:
		o.B = make([]byte, len(idx))
		for j, i := range idx {
			o.B[j] = x.B[i%len(x.B)]
		}
	case List:
		o.V = make([]Value, len(idx))
		for j, i := range idx {
			o.V[j] = x.V[i%len(x.V)]
		}
	}
	if len(x.NA) > 0 {
		o.NA = make([]bool, len(idx))
		for j, i := range idx {
			o.NA[j] = naAt(x, i)
		}
	}
	if nm := names(x); nm != nil {
		ns := make([]string, len(idx))
		for j, i := range idx {
			if i < len(nm) {
				ns[j] = nm[i]
			}
		}
		o = setNames(o, ns)
	}
	return o
}
func qsortImpl(x Value, decreasing, naLast Value) Value {
	dec, _ := asBool(decreasing)
	o := clone(x)
	switch o.Kind {
	case Integer:
		sort.SliceStable(o.I, func(i, j int) bool {
			if dec {
				return o.I[i] > o.I[j]
			}
			return o.I[i] < o.I[j]
		})
	case Double:
		sort.SliceStable(o.D, func(i, j int) bool {
			a, b := o.D[i], o.D[j]
			if math.IsNaN(a) {
				return false
			}
			if math.IsNaN(b) {
				return true
			}
			if dec {
				return a > b
			}
			return a < b
		})
	case String:
		sort.SliceStable(o.S, func(i, j int) bool {
			if dec {
				return o.S[i] > o.S[j]
			}
			return o.S[i] < o.S[j]
		})
	default:
		return ErrValue("unimplemented type in sort")
	}
	return o
}
func rangeImpl(x Value, naRm Value) Value {
	rm, _ := asBool(naRm)
	if length(x) == 0 {
		return Doubles(math.Inf(1), math.Inf(-1))
	}
	mn := math.Inf(1)
	mx := math.Inf(-1)
	seen := false
	for i := 0; i < length(x); i++ {
		v, na := vectorFloat(x, i)
		if na || math.IsNaN(v) {
			if !rm {
				return Value{Kind: Double, D: []float64{0, 0}, NA: []bool{true, true}}
			}
			continue
		}
		seen = true
		if v < mn {
			mn = v
		}
		if v > mx {
			mx = v
		}
	}
	if !seen {
		return Doubles(math.Inf(1), math.Inf(-1))
	}
	return Doubles(mn, mx)
}
func sequenceImpl(nvec, from, by Value) Value {
	n := length(nvec)
	total := 0
	ns := make([]int, n)
	for i := 0; i < n; i++ {
		z, e := asIntAt(nvec, i)
		if e != nil || z < 0 {
			return ErrValue("invalid 'lengths' argument")
		}
		ns[i] = int(z)
		total += int(z)
	}
	out := make([]float64, 0, total)
	for i, k := range ns {
		f := 1.0
		b := 1.0
		if length(from) > 0 {
			f, _ = vectorFloat(from, i)
		}
		if length(by) > 0 {
			b, _ = vectorFloat(by, i)
		}
		for j := 0; j < k; j++ {
			out = append(out, f+float64(j)*b)
		}
	}
	return Doubles(out...)
}
func mergeImpl(xinds, yinds, allx, ally Value) Value {
	ax, _ := asBool(allx)
	ay, _ := asBool(ally)
	mx := map[int64][]int64{}
	my := map[int64][]int64{}
	for i := 0; i < length(xinds); i++ {
		k, _ := asIntAt(xinds, i)
		mx[k] = append(mx[k], int64(i+1))
	}
	for i := 0; i < length(yinds); i++ {
		k, _ := asIntAt(yinds, i)
		my[k] = append(my[k], int64(i+1))
	}
	keys := make([]int64, 0)
	for k := range mx {
		if _, ok := my[k]; ok {
			keys = append(keys, k)
		}
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	var xi, yi, xa, ya []int64
	for _, k := range keys {
		for _, a := range mx[k] {
			for _, b := range my[k] {
				xi = append(xi, a)
				yi = append(yi, b)
			}
		}
	}
	if ax {
		for k, a := range mx {
			if _, ok := my[k]; !ok {
				xa = append(xa, a...)
			}
		}
	}
	if ay {
		for k, a := range my {
			if _, ok := mx[k]; !ok {
				ya = append(ya, a...)
			}
		}
	}
	sort.Slice(xa, func(i, j int) bool { return xa[i] < xa[j] })
	sort.Slice(ya, func(i, j int) bool { return ya[i] < ya[j] })
	return Lists(Ints(xi...), Ints(yi...), Ints(xa...), Ints(ya...))
}

func chartrImpl(old, newv, x Value) Value {
	o, _ := asString(old)
	n, _ := asString(newv)
	or := []rune(o)
	nr := []rune(n)
	if len(or) != len(nr) {
		return ErrValue("'old' is longer than 'new'")
	}
	m := map[rune]rune{}
	for i, r := range or {
		m[r] = nr[i]
	}
	out := clone(x)
	if x.Kind != String {
		return ErrValue("invalid 'x'")
	}
	for i, s := range out.S {
		out.S[i] = strings.Map(func(r rune) rune {
			if q, ok := m[r]; ok {
				return q
			}
			return r
		}, s)
	}
	return out
}
func substrGetsImpl(x, start, stop, val Value) Value {
	if x.Kind != String {
		return ErrValue("replacing in non-character object")
	}
	o := clone(x)
	for i, s := range o.S {
		a, _ := asIntAt(start, i)
		b, _ := asIntAt(stop, i)
		rep, _ := vectorString(val, i)
		rs := []rune(s)
		rr := []rune(rep)
		if a < 1 {
			a = 1
		}
		if b > int64(len(rs)) {
			b = int64(len(rs))
		}
		if a <= b && a <= int64(len(rs)) {
			max := int(b - a + 1)
			if len(rr) > max {
				rr = rr[:max]
			}
			o.S[i] = string(rs[:a-1]) + string(rr) + string(rs[b:])
		}
	}
	return preserveShape(o, x)
}
func intToUtf8Impl(x, multiple, allowSurrogate Value) Value {
	mul, _ := asBool(multiple)
	allow, _ := asBool(allowSurrogate)
	conv := func(i int) (string, bool) {
		v, e := asIntAt(x, i)
		if e != nil {
			return "", true
		}
		if v < 0 || v > utf8.MaxRune || (!allow && v >= 0xD800 && v <= 0xDFFF) {
			return "", true
		}
		return string(rune(v)), false
	}
	if mul {
		o := Value{Kind: String, S: make([]string, length(x)), NA: make([]bool, length(x))}
		for i := 0; i < length(x); i++ {
			o.S[i], o.NA[i] = conv(i)
		}
		return o
	}
	var b strings.Builder
	for i := 0; i < length(x); i++ {
		s, na := conv(i)
		if na {
			return Value{Kind: String, S: []string{""}, NA: []bool{true}}
		}
		b.WriteString(s)
	}
	return Strings(b.String())
}
func utf8ToIntImpl(x Value) Value {
	if x.Kind != String || len(x.S) != 1 {
		return ErrValue("argument must be a character vector of length 1")
	}
	rs := []rune(x.S[0])
	o := make([]int64, len(rs))
	for i, r := range rs {
		o[i] = int64(r)
	}
	return Ints(o...)
}
func validUTF8Impl(x Value) Value {
	if x.Kind != String {
		return ErrValue("argument must be character")
	}
	o := Value{Kind: Logical, L: make([]int8, len(x.S)), NA: append([]bool(nil), x.NA...)}
	for i, s := range x.S {
		if utf8.ValidString(s) {
			o.L[i] = 1
		}
	}
	return o
}
func encodeStringImpl(x, quote, naencode, useBytes, justify, width Value) Value {
	if x.Kind != String {
		return ErrValue("'x' must be a character vector")
	}
	q, _ := asString(quote)
	o := clone(x)
	for i, s := range o.S {
		if naAt(x, i) {
			o.S[i] = "NA"
			continue
		}
		if q != "" {
			s = strings.ReplaceAll(s, "\\", "\\\\")
			s = strings.ReplaceAll(s, q, "\\"+q)
			s = q + s + q
		}
		o.S[i] = s
	}
	return o
}
func compareNumericVersionImpl(a, b Value) Value {
	n := recycleLen(length(a), length(b))
	o := Value{Kind: Integer, I: make([]int64, n), NA: make([]bool, n)}
	for i := 0; i < n; i++ {
		x, nax := vectorString(a, i)
		y, nay := vectorString(b, i)
		if nax || nay {
			o.NA[i] = true
			continue
		}
		o.I[i] = int64(compareVersion(x, y))
	}
	return o
}
func compareVersion(a, b string) int {
	pa := strings.FieldsFunc(a, func(r rune) bool { return r == '.' || r == '-' })
	pb := strings.FieldsFunc(b, func(r rune) bool { return r == '.' || r == '-' })
	n := len(pa)
	if len(pb) > n {
		n = len(pb)
	}
	for i := 0; i < n; i++ {
		var x, y int64
		if i < len(pa) {
			x, _ = strconv.ParseInt(pa[i], 10, 64)
		}
		if i < len(pb) {
			y, _ = strconv.ParseInt(pb[i], 10, 64)
		}
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
	}
	return 0
}
func makeNamesImpl(x, allow Value) Value {
	if x.Kind != String {
		return ErrValue("invalid names")
	}
	allowUnderscore, _ := asBool(allow)
	o := clone(x)
	for i, s := range o.S {
		r := []rune(s)
		var b strings.Builder
		if len(r) == 0 || !(unicode.IsLetter(r[0]) || r[0] == '.') {
			b.WriteByte('X')
		}
		for _, c := range r {
			ok := unicode.IsLetter(c) || unicode.IsDigit(c) || c == '.' || (allowUnderscore && c == '_')
			if ok {
				b.WriteRune(c)
			} else {
				b.WriteByte('.')
			}
		}
		z := b.String()
		switch z {
		case "if", "else", "repeat", "while", "function", "for", "in", "next", "break", "TRUE", "FALSE", "NULL", "Inf", "NaN", "NA", "NA_integer_", "NA_real_", "NA_complex_", "NA_character_":
			z += "."
		}
		o.S[i] = z
	}
	return o
}
func grepRawImpl(pattern, x, offset, ignoreCase, fixed, value, all Value) Value {
	p, _ := asString(pattern)
	if x.Kind != Raw {
		return ErrValue("'x' must be raw")
	}
	off, _ := asInt(offset)
	if off < 1 {
		off = 1
	}
	ic, _ := asBool(ignoreCase)
	fx, _ := asBool(fixed)
	val, _ := asBool(value)
	al, _ := asBool(all)
	hay := string(x.B[off-1:])
	var locs [][]int
	if fx {
		needle := p
		hh := hay
		if ic {
			needle = strings.ToLower(needle)
			hh = strings.ToLower(hh)
		}
		for pos := 0; ; {
			j := strings.Index(hh[pos:], needle)
			if j < 0 {
				break
			}
			j += pos
			locs = append(locs, []int{j, j + len(needle)})
			pos = j + 1
			if !al {
				break
			}
		}
	} else {
		flags := ""
		if ic {
			flags = "(?i)"
		}
		re, e := regexp.Compile(flags + p)
		if e != nil {
			return ErrValue("invalid regular expression: %v", e)
		}
		locs = re.FindAllStringIndex(hay, -1)
		if !al && len(locs) > 1 {
			locs = locs[:1]
		}
	}
	if val {
		if len(locs) == 0 {
			return Raws()
		}
		var b []byte
		for _, q := range locs {
			b = append(b, []byte(hay[q[0]:q[1]])...)
		}
		return Raws(b...)
	}
	o := make([]int64, len(locs))
	for i, q := range locs {
		o[i] = int64(q[0]) + off
	}
	return Ints(o...)
}

func math2Impl(x, y Value, which int) Value {
	n := recycleLen(length(x), length(y))
	o := Value{Kind: Double, D: make([]float64, n), NA: make([]bool, n)}
	for i := 0; i < n; i++ {
		a, na := vectorFloat(x, i)
		b, nb := vectorFloat(y, i)
		if na || nb {
			o.NA[i] = true
			continue
		}
		switch which {
		case 0:
			o.D[i] = math.Atan2(a, b)
		case 1:
			o.D[i] = math.Pow(a, b)
		case 2:
			o.D[i] = math.Log(a) / math.Log(b)
		case 24:
			o.D[i] = math.Jn(int(b), a)
		case 25:
			o.D[i] = math.Yn(int(b), a)
		default:
			o.D[i] = math.NaN()
		}
	}
	return preserveShape(o, x)
}
func math3Impl(a, b, c Value, which int) Value {
	n := recycleLen(recycleLen(length(a), length(b)), length(c))
	o := Value{Kind: Double, D: make([]float64, n), NA: make([]bool, n)}
	for i := 0; i < n; i++ {
		x, nx := vectorFloat(a, i)
		y, ny := vectorFloat(b, i)
		z, nz := vectorFloat(c, i)
		if nx || ny || nz {
			o.NA[i] = true
			continue
		}
		switch which {
		case 0:
			o.D[i] = x + y + z
		case 43:
			o.D[i] = math.NaN()
		case 44:
			o.D[i] = math.NaN()
		default:
			o.D[i] = math.NaN()
		}
	}
	return preserveShape(o, a)
}
func complexImpl(real, imag Value) Value {
	n := recycleLen(length(real), length(imag))
	o := Value{Kind: Complex, Z: make([]complex128, n), NA: make([]bool, n)}
	for i := 0; i < n; i++ {
		r, nr := vectorFloat(real, i)
		im, ni := vectorFloat(imag, i)
		if nr || ni {
			o.NA[i] = true
		} else {
			o.Z[i] = complex(r, im)
		}
	}
	return o
}
func polyrootImpl(z Value) Value {
	if z.Kind != Complex && z.Kind != Double && z.Kind != Integer {
		return ErrValue("invalid polynomial coefficient")
	}
	n := length(z) - 1
	if n <= 0 {
		return Value{Kind: Complex}
	}
	coef := make([]complex128, n+1)
	for i := range coef {
		switch z.Kind {
		case Complex:
			coef[i] = z.Z[i]
		default:
			f, _ := vectorFloat(z, i)
			coef[i] = complex(f, 0)
		}
	}
	for n > 0 && coef[n] == 0 {
		n--
	}
	if n == 0 {
		return Value{Kind: Complex}
	}
	roots := make([]complex128, n)
	radius := 1.0
	for k := 0; k < n; k++ {
		roots[k] = cmplx.Rect(radius, 2*math.Pi*float64(k)/float64(n))
	}
	for iter := 0; iter < 200; iter++ {
		maxd := 0.0
		for i := range roots {
			p := coef[n]
			for j := n - 1; j >= 0; j-- {
				p = p*roots[i] + coef[j]
			}
			den := complex(1, 0)
			for j := range roots {
				if i != j {
					den *= roots[i] - roots[j]
				}
			}
			if den == 0 {
				continue
			}
			d := p / den
			roots[i] -= d
			if cmplx.Abs(d) > maxd {
				maxd = cmplx.Abs(d)
			}
		}
		if maxd < 1e-12 {
			break
		}
	}
	return Value{Kind: Complex, Z: roots}
}

func filePrimitive(entry string, args Value) Value {
	p := func(i int) string { s, _ := asString(arg(args, i)); return s }
	switch entry {
	case "do_filecreate":
		paths := arg(args, 0)
		o := Value{Kind: Logical, L: make([]int8, length(paths))}
		for i := 0; i < length(paths); i++ {
			s, _ := vectorString(paths, i)
			f, e := os.OpenFile(s, os.O_CREATE|os.O_WRONLY, 0666)
			if e == nil {
				o.L[i] = 1
				f.Close()
			}
		}
		return o
	case "do_fileremove":
		paths := arg(args, 0)
		o := Value{Kind: Logical, L: make([]int8, length(paths))}
		for i := 0; i < length(paths); i++ {
			s, _ := vectorString(paths, i)
			if os.Remove(s) == nil {
				o.L[i] = 1
			}
		}
		return o
	case "do_filerename":
		return Bool(os.Rename(p(0), p(1)) == nil)
	case "do_fileappend":
		dst, src := p(0), p(1)
		in, e := os.Open(src)
		if e != nil {
			return Bool(false)
		}
		defer in.Close()
		out, e := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if e != nil {
			return Bool(false)
		}
		defer out.Close()
		_, e = io.Copy(out, in)
		return Bool(e == nil)
	case "do_fileaccess":
		paths := arg(args, 0)
		o := Value{Kind: Integer, I: make([]int64, length(paths))}
		for i := 0; i < length(paths); i++ {
			s, _ := vectorString(paths, i)
			_, e := os.Stat(s)
			if e != nil {
				o.I[i] = -1
			}
		}
		return o
	case "do_fileinfo":
		paths := arg(args, 0)
		rows := make([]Value, length(paths))
		for i := 0; i < length(paths); i++ {
			s, _ := vectorString(paths, i)
			st, e := os.Stat(s)
			if e != nil {
				rows[i] = Lists(Value{Kind: Double, NA: []bool{true}}, Value{Kind: Logical, NA: []bool{true}})
				continue
			}
			rows[i] = Lists(Doubles(float64(st.Size())), Bool(st.IsDir()), Doubles(float64(st.ModTime().Unix())))
		}
		return Lists(rows...)
	case "do_filelink", "do_filesymlink":
		e := os.Symlink(p(0), p(1))
		return Bool(e == nil)
	case "do_readlink":
		s, e := os.Readlink(p(0))
		if e != nil {
			return Strings("")
		}
		return Strings(s)
	case "do_setFileTime":
		t, _ := asFloat(arg(args, 1))
		tm := time.Unix(int64(t), 0)
		return Bool(os.Chtimes(p(0), tm, tm) == nil)
	case "do_mkjunction":
		return Bool(os.Symlink(p(0), p(1)) == nil)
	case "do_listdirs":
		root := p(0)
		var out []string
		filepath.WalkDir(root, func(path string, d os.DirEntry, e error) error {
			if e == nil && d.IsDir() {
				out = append(out, path)
			}
			return nil
		})
		return Strings(out...)
	}
	return ErrValue("%s: unsupported filesystem operation", entry)
}
func systemPrimitive(entry string, args Value) Value {
	switch entry {
	case "do_sysgetpid":
		return Ints(int64(os.Getpid()))
	case "do_systime":
		return Doubles(float64(time.Now().UnixNano()) / 1e9)
	case "do_tempdir":
		return Strings(os.TempDir())
	case "do_tempfile":
		pat, _ := asString(arg(args, 0))
		dir := os.TempDir()
		if d, e := asString(arg(args, 1)); e == nil && d != "" {
			dir = d
		}
		f, e := os.CreateTemp(dir, pat)
		if e != nil {
			return ErrValue("tempfile: %v", e)
		}
		name := f.Name()
		f.Close()
		os.Remove(name)
		return Strings(name)
	case "do_commandArgs":
		return Strings(os.Args...)
	case "do_syswhich":
		cmds := arg(args, 0)
		o := Value{Kind: String, S: make([]string, length(cmds))}
		for i := 0; i < length(cmds); i++ {
			c, _ := vectorString(cmds, i)
			if strings.ContainsRune(c, os.PathSeparator) {
				if _, e := os.Stat(c); e == nil {
					o.S[i] = c
				}
				continue
			}
			for _, d := range filepath.SplitList(os.Getenv("PATH")) {
				p := filepath.Join(d, c)
				if _, e := os.Stat(p); e == nil {
					o.S[i] = p
					break
				}
			}
		}
		return o
	case "do_sysinfo":
		host, _ := os.Hostname()
		wd, _ := os.Getwd()
		return setNames(Strings(runtimeGOOS(), host, os.Getenv("USER"), wd), []string{"sysname", "nodename", "user", "working_dir"})
	}
	return ErrValue("%s: operation not available in the Pure-Go runtime", entry)
}
func runtimeGOOS() string { return os.Getenv("GOOS") }

func connectionFromFile(path, mode string) Value {
	var f *os.File
	var e error
	if strings.Contains(mode, "w") {
		f, e = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	} else if strings.Contains(mode, "a") {
		f, e = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	} else {
		f, e = os.Open(path)
	}
	if e != nil {
		return ErrValue("connection: %v", e)
	}
	return Value{Kind: Connection, Conn: &Conn{R: f, W: f, C: f, Name: path, Open: true, Text: !strings.Contains(mode, "b")}}
}
func connectionPrimitive(entry string, args Value) Value {
	switch entry {
	case "do_close":
		c := arg(args, 0)
		if c.Kind != Connection || c.Conn == nil {
			return ErrValue("invalid connection")
		}
		c.Conn.mu.Lock()
		defer c.Conn.mu.Unlock()
		if c.Conn.C != nil {
			if e := c.Conn.C.Close(); e != nil {
				return ErrValue("close: %v", e)
			}
		}
		c.Conn.Open = false
		return Nil
	case "do_isopen":
		c := arg(args, 0)
		return Bool(c.Kind == Connection && c.Conn != nil && c.Conn.Open)
	case "do_isseekable":
		c := arg(args, 0)
		if c.Kind != Connection || c.Conn == nil {
			return Bool(false)
		}
		_, ok := c.Conn.R.(io.Seeker)
		return Bool(ok)
	case "do_flush":
		c := arg(args, 0)
		if c.Kind != Connection || c.Conn == nil {
			return ErrValue("invalid connection")
		}
		if f, ok := c.Conn.W.(*bufio.Writer); ok {
			return errorOrNil(f.Flush())
		}
		return Nil
	case "do_stdin":
		return Value{Kind: Connection, Conn: &Conn{R: os.Stdin, Name: "stdin", Open: true, Text: true}}
	case "do_stdout":
		return Value{Kind: Connection, Conn: &Conn{W: os.Stdout, Name: "stdout", Open: true, Text: true}}
	case "do_stderr":
		return Value{Kind: Connection, Conn: &Conn{W: os.Stderr, Name: "stderr", Open: true, Text: true}}
	case "do_fifo", "do_pipe", "do_serversocket", "do_sockconn", "do_url", "do_unz", "do_gzfile", "do_gzcon":
		return ErrValue("%s requires a transport/compression backend not represented by this batch corpus", entry)
	case "do_rawconnection":
		buf := bytes.NewBuffer(append([]byte(nil), arg(args, 1).B...))
		return Value{Kind: Connection, Conn: &Conn{R: buf, W: buf, Name: "rawConnection", Open: true}}
	case "do_rawconvalue":
		c := arg(args, 0)
		if c.Kind != Connection || c.Conn == nil {
			return ErrValue("invalid connection")
		}
		if b, ok := c.Conn.R.(*bytes.Buffer); ok {
			return Raws(b.Bytes()...)
		}
		return ErrValue("not a raw connection")
	case "do_textconnection":
		buf := bytes.NewBufferString("")
		return Value{Kind: Connection, Conn: &Conn{R: buf, W: buf, Name: "textConnection", Open: true, Text: true}}
	}
	return ErrValue("%s: connection operation not implemented by available Go interfaces", entry)
}
func errorOrNil(e error) Value {
	if e != nil {
		return ErrValue("%v", e)
	}
	return Nil
}

func envPrimitive(entry string, args Value, env Value) Value {
	switch entry {
	case "do_envirName":
		e := arg(args, 0)
		if e.Kind != Environment || e.Env == nil {
			return Strings("")
		}
		return Strings(e.Env.Name)
	case "do_envIsLocked":
		e := arg(args, 0)
		return Bool(e.Kind == Environment && e.Env != nil && e.Env.Locked)
	case "do_lockEnv":
		e := arg(args, 0)
		if e.Kind != Environment || e.Env == nil {
			return ErrValue("not an environment")
		}
		e.Env.mu.Lock()
		e.Env.Locked = true
		e.Env.mu.Unlock()
		return Nil
	case "do_ls":
		e := arg(args, 0)
		if e.Kind != Environment || e.Env == nil {
			e = env
		}
		if e.Kind != Environment || e.Env == nil {
			return ErrValue("invalid environment")
		}
		e.Env.mu.RLock()
		ks := make([]string, 0, len(e.Env.Values))
		for k := range e.Env.Values {
			ks = append(ks, k)
		}
		e.Env.mu.RUnlock()
		sort.Strings(ks)
		return Strings(ks...)
	case "do_remove":
		namesv := arg(args, 0)
		e := arg(args, 1)
		if e.Kind != Environment {
			e = env
		}
		for i := 0; i < length(namesv); i++ {
			k, _ := vectorString(namesv, i)
			if er := envRemove(e.Env, k); er != nil {
				return ErrValue("%v", er)
			}
		}
		return Nil
	case "do_env2list":
		e := arg(args, 0)
		if e.Kind != Environment || e.Env == nil {
			return ErrValue("not an environment")
		}
		e.Env.mu.RLock()
		ks := make([]string, 0, len(e.Env.Values))
		for k := range e.Env.Values {
			ks = append(ks, k)
		}
		sort.Strings(ks)
		vs := make([]Value, len(ks))
		for i, k := range ks {
			vs[i] = e.Env.Values[k]
		}
		e.Env.mu.RUnlock()
		return setNames(Lists(vs...), ks)
	case "do_list2env":
		x := arg(args, 0)
		e := arg(args, 1)
		if e.Kind != Environment || e.Env == nil {
			return ErrValue("invalid environment")
		}
		nm := names(x)
		if len(nm) != len(x.V) {
			return ErrValue("names(x) must be a character vector of the same length as x")
		}
		for i, k := range nm {
			if er := envSet(e.Env, k, x.V[i]); er != nil {
				return ErrValue("%v", er)
			}
		}
		return e
	case "do_parentenvgets":
		e := arg(args, 0)
		p := arg(args, 1)
		if e.Kind != Environment || p.Kind != Environment {
			return ErrValue("invalid environment")
		}
		e.Env.mu.Lock()
		e.Env.Parent = p.Env
		e.Env.mu.Unlock()
		return e
	case "do_pos2env":
		return env
	case "do_as_environment":
		x := arg(args, 0)
		if x.Kind == Environment {
			return x
		}
		if x.Kind == List {
			e := NewEnv(nil, "")
			nm := names(x)
			for i, v := range x.V {
				if i < len(nm) {
					envSet(e.Env, nm[i], v)
				}
			}
			return e
		}
		return ErrValue("invalid object for 'as.environment'")
	}
	return ErrValue("%s: environment feature requires evaluator state not present in isolated batch", entry)
}

func serializePrimitive(entry string, args Value) Value { // deterministic private encoding used only inside this Pure-Go runtime
	switch entry {
	case "do_serialize":
		var b bytes.Buffer
		if e := encodeValue(&b, arg(args, 0)); e != nil {
			return ErrValue("serialize: %v", e)
		}
		return Raws(b.Bytes()...)
	case "do_unserializeFromConn":
		return ErrValue("unserialize from connection requires stream framing from GNU R format")
	}
	return ErrValue("%s: GNU R serialization wire format is not implemented", entry)
}
func encodeValue(w io.Writer, v Value) error {
	if e := binary.Write(w, binary.LittleEndian, uint8(v.Kind)); e != nil {
		return e
	}
	n := uint64(length(v))
	if e := binary.Write(w, binary.LittleEndian, n); e != nil {
		return e
	}
	switch v.Kind {
	case Integer:
		for _, x := range v.I {
			binary.Write(w, binary.LittleEndian, x)
		}
	case Double:
		for _, x := range v.D {
			binary.Write(w, binary.LittleEndian, x)
		}
	case Raw:
		w.Write(v.B)
	case String:
		for _, s := range v.S {
			binary.Write(w, binary.LittleEndian, uint64(len(s)))
			io.WriteString(w, s)
		}
	default:
		return fmt.Errorf("kind %d not serializable", v.Kind)
	}
	return nil
}

func unsupported(entry, reason string) Value { return ErrValue("%s: %s", entry, reason) }

func boolInt(b bool) int64 {
	if b {
		return 1
	}
	return 0
}
func addressImpl(v Value) Value { return Strings(fmt.Sprintf("%p", &v)) }
func coerceMode(v Value, mode string) Value {
	n := length(v)
	switch mode {
	case "double", "numeric":
		o := Value{Kind: Double, D: make([]float64, n), NA: append([]bool(nil), v.NA...)}
		for i := 0; i < n; i++ {
			o.D[i], _ = vectorFloat(v, i)
		}
		return copyAttrs(o, v)
	case "integer":
		o := Value{Kind: Integer, I: make([]int64, n), NA: append([]bool(nil), v.NA...)}
		for i := 0; i < n; i++ {
			o.I[i], _ = asIntAt(v, i)
		}
		return copyAttrs(o, v)
	case "logical":
		o := Value{Kind: Logical, L: make([]int8, n), NA: append([]bool(nil), v.NA...)}
		for i := 0; i < n; i++ {
			x, _ := asIntAt(v, i)
			if x != 0 {
				o.L[i] = 1
			}
		}
		return copyAttrs(o, v)
	case "character":
		o := Value{Kind: String, S: make([]string, n), NA: append([]bool(nil), v.NA...)}
		for i := 0; i < n; i++ {
			o.S[i], _ = vectorString(v, i)
		}
		return copyAttrs(o, v)
	case "raw":
		o := Value{Kind: Raw, B: make([]byte, n)}
		for i := 0; i < n; i++ {
			x, _ := asIntAt(v, i)
			o.B[i] = byte(x)
		}
		return copyAttrs(o, v)
	case "complex":
		o := Value{Kind: Complex, Z: make([]complex128, n), NA: append([]bool(nil), v.NA...)}
		for i := 0; i < n; i++ {
			x, _ := vectorFloat(v, i)
			o.Z[i] = complex(x, 0)
		}
		return copyAttrs(o, v)
	}
	return ErrValue("invalid value for 'storage.mode': %s", mode)
}
func isVectorImpl(x, mode Value) Value {
	m, _ := asString(mode)
	ok := x.Kind == Logical || x.Kind == Integer || x.Kind == Double || x.Kind == Complex || x.Kind == String || x.Kind == Raw || x.Kind == List
	if m != "any" && m != "" {
		names := map[string]Kind{"logical": Logical, "integer": Integer, "double": Double, "numeric": Double, "complex": Complex, "character": String, "raw": Raw, "list": List}
		k, found := names[m]
		ok = ok && found && x.Kind == k
	}
	if x.Attr != nil {
		for k := range x.Attr {
			if k != "names" {
				ok = false
				break
			}
		}
	}
	return Bool(ok)
}
func shortRowNamesImpl(x, typ Value) Value {
	rn := attr(x, "row.names")
	t, _ := asInt(typ)
	if rn.Kind == Null {
		return Nil
	}
	if t == 0 {
		return rn
	}
	n := length(rn)
	if rn.Kind == Integer && len(rn.I) == 2 && len(rn.NA) > 0 && rn.NA[0] {
		n = int(-rn.I[1])
	}
	if t == 1 {
		return Ints(int64(n))
	}
	return Ints(int64(n))
}
func isListFactorImpl(x, recursive Value) Value {
	rec, _ := asBool(recursive)
	if x.Kind != List {
		return Bool(false)
	}
	var check func(Value) bool
	check = func(v Value) bool {
		if c := attr(v, "class"); c.Kind == String {
			for _, s := range c.S {
				if s == "factor" {
					return true
				}
			}
		}
		if rec && v.Kind == List {
			for _, q := range v.V {
				if !check(q) {
					return false
				}
			}
			return true
		}
		return false
	}
	for _, v := range x.V {
		if !check(v) {
			return Bool(false)
		}
	}
	return Bool(true)
}
func abbreviateImpl(x, minlength, useClasses, dot Value) Value {
	if x.Kind != String {
		return ErrValue("'x' must be a character vector")
	}
	ml, _ := asInt(minlength)
	if ml < 1 {
		ml = 1
	}
	o := clone(x)
	for i, s := range o.S {
		r := []rune(strings.TrimSpace(s))
		if int64(len(r)) <= ml {
			continue
		} // remove interior vowels first, then truncate
		first := r[0]
		rest := make([]rune, 0, len(r)-1)
		for _, c := range r[1:] {
			if !strings.ContainsRune("aeiouAEIOU", c) {
				rest = append(rest, c)
			}
		}
		z := append([]rune{first}, rest...)
		if int64(len(z)) < ml {
			z = r
		}
		if int64(len(z)) > ml {
			z = z[:ml]
		}
		o.S[i] = string(z)
	}
	return preserveShape(o, x)
}
func crc64Impl(x Value) Value {
	var data []byte
	if x.Kind == Raw {
		data = x.B
	} else if x.Kind == String && len(x.S) > 0 {
		data = []byte(x.S[0])
	} else {
		return ErrValue("invalid input")
	}
	const poly uint64 = 0xC96C5795D7870F42
	var crc uint64 = ^uint64(0)
	for _, b := range data {
		crc ^= uint64(b)
		for i := 0; i < 8; i++ {
			if crc&1 != 0 {
				crc = (crc >> 1) ^ poly
			} else {
				crc >>= 1
			}
		}
	}
	crc = ^crc
	return Strings(fmt.Sprintf("%016x", crc))
}

var runtimeState = struct {
	mu         sync.Mutex
	methods    bool
	jit        int64
	gcInfo     bool
	lastError  string
	icuCollate string
	trace      bool
}{methods: true, icuCollate: "C"}

func statePrimitive(entry string, args Value) Value {
	runtimeState.mu.Lock()
	defer runtimeState.mu.Unlock()
	switch entry {
	case "do_S4on":
		return Bool(runtimeState.methods)
	case "do_ICUget":
		return Strings(runtimeState.icuCollate)
	case "do_ICUset":
		old := runtimeState.icuCollate
		if s, e := asString(arg(args, 0)); e == nil {
			runtimeState.icuCollate = s
		}
		return Strings(old)
	case "do_enablejit":
		old := runtimeState.jit
		if n, e := asInt(arg(args, 0)); e == nil {
			runtimeState.jit = n
		}
		return Ints(old)
	case "do_gcinfo":
		old := runtimeState.gcInfo
		if nargs(args) > 0 {
			if b, e := asBool(arg(args, 0)); e == nil {
				runtimeState.gcInfo = b
			}
		}
		return Bool(old)
	case "do_geterrmessage":
		return Strings(runtimeState.lastError)
	case "do_seterrmessage":
		s, e := asString(arg(args, 0))
		if e != nil {
			return ErrValue("invalid error message")
		}
		runtimeState.lastError = s
		return Nil
	case "do_traceOnOff":
		old := runtimeState.trace
		if nargs(args) > 0 {
			b, _ := asBool(arg(args, 0))
			runtimeState.trace = b
		}
		return Bool(old)
	}
	return ErrValue("%s: invalid runtime-state primitive", entry)
}

func gcPrimitive(entry string, args Value) Value {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	switch entry {
	case "do_gc":
		runtime.GC()
		runtime.ReadMemStats(&m)
		return setNames(Doubles(float64(m.HeapAlloc)/8, float64(m.HeapSys)/8), []string{"used", "gc trigger"})
	case "do_gctime":
		return Doubles(float64(m.PauseTotalNs) / 1e9)
	case "do_memoryprofile":
		return setNames(Doubles(float64(m.Alloc), float64(m.TotalAlloc), float64(m.Sys), float64(m.NumGC)), []string{"Alloc", "TotalAlloc", "Sys", "NumGC"})
	case "do_gctorture", "do_gctorture2":
		return Bool(false)
	}
	return ErrValue("%s: invalid GC primitive", entry)
}

func allNamesImpl(x Value, functions, unique, maxNames Value) Value {
	incFn, _ := asBool(functions)
	uniq, _ := asBool(unique)
	max, _ := asInt(maxNames)
	if max <= 0 {
		max = math.MaxInt32
	}
	out := []string{}
	seen := map[string]bool{}
	var walk func(Value, bool)
	walk = func(v Value, head bool) {
		if int64(len(out)) >= max {
			return
		}
		switch v.Kind {
		case Symbol:
			if head && !incFn {
				return
			}
			if !uniq || !seen[v.Name] {
				out = append(out, v.Name)
				seen[v.Name] = true
			}
		case List:
			for i, q := range v.V {
				walk(q, i == 0)
			}
		}
	}
	walk(x, false)
	return Strings(out...)
}

func levenshtein(a, b []rune) int {
	if len(a) == 0 {
		return len(b)
	}
	prev := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}
	for i, ca := range a {
		cur := make([]int, len(b)+1)
		cur[0] = i + 1
		for j, cb := range b {
			c := 0
			if ca != cb {
				c = 1
			}
			x := prev[j+1] + 1
			if cur[j]+1 < x {
				x = cur[j] + 1
			}
			if prev[j]+c < x {
				x = prev[j] + c
			}
			cur[j+1] = x
		}
		prev = cur
	}
	return prev[len(b)]
}
func adistImpl(x, y Value) Value {
	if x.Kind != String || y.Kind != String {
		return ErrValue("'x' and 'y' must be character")
	}
	nr, nc := len(x.S), len(y.S)
	o := Value{Kind: Integer, I: make([]int64, nr*nc)}
	for j, b := range y.S {
		for i, a := range x.S {
			o.I[i+nr*j] = int64(levenshtein([]rune(a), []rune(b)))
		}
	}
	o = withAttr(o, "dim", Ints(int64(nr), int64(nc)))
	return o
}
func agrepImpl(pattern, x, maxDist, value Value) Value {
	p, e := asString(pattern)
	if e != nil || x.Kind != String {
		return ErrValue("invalid pattern or text")
	}
	md := 1
	if f, e := asFloat(maxDist); e == nil {
		if f < 1 {
			md = int(math.Ceil(f * float64(len([]rune(p)))))
		} else {
			md = int(f)
		}
	}
	var idx []int64
	var vals []string
	for i, s := range x.S {
		if levenshtein([]rune(p), []rune(s)) <= md {
			idx = append(idx, int64(i+1))
			vals = append(vals, s)
		}
	}
	v, _ := asBool(value)
	if v {
		return Strings(vals...)
	}
	return Ints(idx...)
}
func aregexecImpl(pattern, text Value) Value {
	p, e := asString(pattern)
	if e != nil || text.Kind != String {
		return ErrValue("invalid input")
	}
	re, e := regexp.Compile(p)
	if e != nil {
		return ErrValue("invalid regular expression: %v", e)
	}
	out := make([]Value, len(text.S))
	for i, s := range text.S {
		q := re.FindStringSubmatchIndex(s)
		if q == nil {
			out[i] = Ints(-1)
			continue
		}
		z := make([]int64, len(q)/2)
		for j := 0; j < len(q); j += 2 {
			if q[j] < 0 {
				z[j/2] = -1
			} else {
				z[j/2] = int64(q[j] + 1)
			}
		}
		out[i] = Ints(z...)
	}
	return Lists(out...)
}
func regexprImpl(pattern, text Value, global bool) Value {
	p, e := asString(pattern)
	if e != nil || text.Kind != String {
		return ErrValue("invalid input")
	}
	re, e := regexp.Compile(p)
	if e != nil {
		return ErrValue("invalid regular expression: %v", e)
	}
	if global {
		out := make([]Value, len(text.S))
		for i, s := range text.S {
			qs := re.FindAllStringIndex(s, -1)
			z := make([]int64, len(qs))
			for j, q := range qs {
				z[j] = int64(q[0] + 1)
			}
			out[i] = Ints(z...)
		}
		return Lists(out...)
	}
	o := Value{Kind: Integer, I: make([]int64, len(text.S))}
	for i, s := range text.S {
		q := re.FindStringIndex(s)
		if q == nil {
			o.I[i] = -1
		} else {
			o.I[i] = int64(q[0] + 1)
		}
	}
	return o
}

func readDCFImpl(path Value) Value {
	p, e := asString(path)
	if e != nil {
		return ErrValue("invalid file")
	}
	b, e := os.ReadFile(p)
	if e != nil {
		return ErrValue("readDCF: %v", e)
	}
	records := strings.Split(strings.ReplaceAll(string(b), "\r\n", "\n"), "\n\n")
	cols := []string{}
	colset := map[string]bool{}
	parsed := []map[string]string{}
	for _, rec := range records {
		m := map[string]string{}
		last := ""
		for _, ln := range strings.Split(rec, "\n") {
			if strings.TrimSpace(ln) == "" {
				continue
			}
			if (strings.HasPrefix(ln, " ") || strings.HasPrefix(ln, "\t")) && last != "" {
				m[last] += "\n" + strings.TrimSpace(ln)
				continue
			}
			i := strings.IndexByte(ln, ':')
			if i < 1 {
				continue
			}
			k := ln[:i]
			v := strings.TrimSpace(ln[i+1:])
			m[k] = v
			last = k
			if !colset[k] {
				colset[k] = true
				cols = append(cols, k)
			}
		}
		if len(m) > 0 {
			parsed = append(parsed, m)
		}
	}
	vals := make([]string, len(parsed)*len(cols))
	for j, c := range cols {
		for i, m := range parsed {
			vals[i+len(parsed)*j] = m[c]
		}
	}
	v := Strings(vals...)
	v = withAttr(v, "dim", Ints(int64(len(parsed)), int64(len(cols))))
	v = withAttr(v, "dimnames", Lists(Nil, Strings(cols...)))
	return v
}
func readEnvironImpl(path Value) Value {
	p, e := asString(path)
	if e != nil {
		return Bool(false)
	}
	f, e := os.Open(p)
	if e != nil {
		return Bool(false)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	ok := true
	for sc.Scan() {
		ln := strings.TrimSpace(sc.Text())
		if ln == "" || strings.HasPrefix(ln, "#") {
			continue
		}
		i := strings.IndexByte(ln, '=')
		if i < 1 {
			continue
		}
		k := strings.TrimSpace(ln[:i])
		v := strings.TrimSpace(ln[i+1:])
		v = os.ExpandEnv(v)
		if e := os.Setenv(k, v); e != nil {
			ok = false
		}
	}
	return Bool(ok && sc.Err() == nil)
}

func prettyImpl(lo, up, n Value) Value {
	a, e1 := asFloat(lo)
	b, e2 := asFloat(up)
	nn, e3 := asInt(n)
	if e1 != nil || e2 != nil || e3 != nil || nn < 1 {
		return ErrValue("invalid pretty arguments")
	}
	if a > b {
		a, b = b, a
	}
	span := b - a
	if span == 0 {
		span = math.Abs(a)
		if span == 0 {
			span = 1
		}
	}
	raw := span / float64(nn)
	p := math.Pow(10, math.Floor(math.Log10(raw)))
	q := raw / p
	step := p
	if q >= 5 {
		step = 5 * p
	} else if q >= 2 {
		step = 2 * p
	}
	start := math.Floor(a/step) * step
	end := math.Ceil(b/step) * step
	return Doubles(start, end, step)
}

func systemImpl(command Value, intern Value) Value {
	cmd, e := asString(command)
	if e != nil {
		return ErrValue("invalid command")
	}
	capture, _ := asBool(intern)
	c := exec.Command("/bin/sh", "-c", cmd)
	if capture {
		b, e := c.Output()
		if e != nil {
			return ErrValue("system: %v", e)
		}
		s := strings.TrimSuffix(string(b), "\n")
		if s == "" {
			return Strings()
		}
		return Strings(strings.Split(s, "\n")...)
	}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	if e := c.Run(); e != nil {
		if ee, ok := e.(*exec.ExitError); ok {
			return Ints(int64(ee.ExitCode()))
		}
		return ErrValue("system: %v", e)
	}
	return Ints(0)
}
func strptimeImpl(x, format Value) Value {
	if x.Kind != String || format.Kind != String {
		return ErrValue("invalid time input")
	}
	o := Value{Kind: Double, D: make([]float64, len(x.S)), NA: make([]bool, len(x.S))}
	for i, s := range x.S {
		f, _ := vectorString(format, i)
		layout := rTimeLayout(f)
		t, e := time.ParseInLocation(layout, s, time.Local)
		if e != nil {
			o.NA[i] = true
		} else {
			o.D[i] = float64(t.Unix())
		}
	}
	o = withAttr(o, "class", Strings("POSIXct", "POSIXt"))
	return o
}
func rTimeLayout(f string) string {
	r := strings.NewReplacer("%Y", "2006", "%y", "06", "%m", "01", "%d", "02", "%H", "15", "%M", "04", "%S", "05", "%z", "-0700", "%Z", "MST", "%b", "Jan", "%B", "January", "%a", "Mon", "%A", "Monday")
	return r.Replace(f)
}

func mapplyImpl(fun, dots, more Value, env Value) Value {
	if fun.Kind != Function || fun.Fn == nil {
		return ErrValue("'FUN' is not a function")
	}
	if dots.Kind != List {
		return ErrValue("'dots' must be a list")
	}
	n := 0
	for _, v := range dots.V {
		if length(v) > n {
			n = length(v)
		}
	}
	out := make([]Value, n)
	for i := 0; i < n; i++ {
		av := make([]Value, 0, len(dots.V)+length(more))
		for _, v := range dots.V {
			av = append(av, subsetIndices(v, []int{i % maxInt(1, length(v))}))
		}
		if more.Kind == List {
			av = append(av, more.V...)
		}
		z, e := fun.Fn(av, env.Env)
		if e != nil {
			return ErrValue("mapply: %v", e)
		}
		out[i] = z
	}
	return Lists(out...)
}
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func rapplyImpl(x, fun, classes, defhow, env Value) Value {
	if x.Kind != List || fun.Kind != Function || fun.Fn == nil {
		return ErrValue("invalid rapply arguments")
	}
	o := clone(x)
	var walk func(Value) Value
	walk = func(v Value) Value {
		if v.Kind == List {
			q := clone(v)
			for i, z := range q.V {
				q.V[i] = walk(z)
			}
			return q
		}
		z, e := fun.Fn([]Value{v}, env.Env)
		if e != nil {
			return ErrValue("rapply: %v", e)
		}
		return z
	}
	return walk(o)
}

func printDefaultImpl(x Value) Value { fmt.Fprintln(os.Stdout, formatValue(x)); return x }
func formatValue(v Value) string {
	if v.Kind == Error && v.Err != nil {
		return "Error: " + v.Err.Error()
	}
	switch v.Kind {
	case Null:
		return "NULL"
	case String:
		return fmt.Sprintf("%q", v.S)
	case Integer:
		return fmt.Sprint(v.I)
	case Double:
		return fmt.Sprint(v.D)
	case Logical:
		return fmt.Sprint(v.L)
	case Raw:
		return fmt.Sprintf("%x", v.B)
	case Complex:
		return fmt.Sprint(v.Z)
	case Symbol:
		return v.Name
	case List:
		parts := make([]string, len(v.V))
		for i, x := range v.V {
			parts[i] = formatValue(x)
		}
		return "list(" + strings.Join(parts, ", ") + ")"
	}
	return fmt.Sprintf("<%d>", v.Kind)
}
func dputImpl(x, file Value) Value {
	s := formatValue(x)
	if file.Kind == Connection && file.Conn != nil && file.Conn.W != nil {
		_, e := io.WriteString(file.Conn.W, s+"\n")
		if e != nil {
			return ErrValue("dput: %v", e)
		}
		return x
	}
	if p, e := asString(file); e == nil && p != "" {
		if e := os.WriteFile(p, []byte(s+"\n"), 0666); e != nil {
			return ErrValue("dput: %v", e)
		}
	}
	return x
}

func s4Impl(entry string, args Value) Value {
	switch entry {
	case "do_setS4Object":
		x := clone(arg(args, 0))
		flag, _ := asBool(arg(args, 1))
		if flag {
			x = withAttr(x, ".S4", Bool(true))
		} else if x.Attr != nil {
			delete(x.Attr, ".S4")
		}
		return x
	case "do_altrep_class":
		return attr(arg(args, 0), "class")
	case "do_AT":
		obj := arg(args, 0)
		slot, _ := asString(arg(args, 1))
		if v := attr(obj, slot); v.Kind != Null {
			return v
		}
		return ErrValue("no slot of name %q", slot)
	}
	return ErrValue("%s: invalid S4 primitive", entry)
}
