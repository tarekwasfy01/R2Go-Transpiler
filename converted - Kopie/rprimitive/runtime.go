package rprimitive

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// Value is the dynamic value carrier used at the GNU R compatibility boundary.
type Value any

// Ref models C address/dereference semantics without cgo.
type Ref interface {
	Get() Value
	Set(Value)
}

type localRef struct{ p *Value }

func (r localRef) Get() Value  { return *r.p }
func (r localRef) Set(v Value) { *r.p = v }
func LocalRef(p *Value) Ref    { return localRef{p: p} }
func Assign(r Ref, v Value) Value {
	if r != nil {
		r.Set(v)
	}
	return v
}

// Hook is one GNU R/C-runtime operation implemented by an embedding runtime.
type Hook func(args ...Value) Value

// RuntimeError models an R-style non-local error at the translated runtime
// boundary. Embedders can recover it and expose it as their normal evaluator
// error without mistaking it for an unimplemented hook.
type RuntimeError struct {
	Operation string
	Message   string
}

func (e RuntimeError) Error() string {
	if e.Operation == "" {
		return e.Message
	}
	if e.Message == "" {
		return e.Operation
	}
	return e.Operation + ": " + e.Message
}

// PluginBoundaryError is deliberately distinct from RuntimeError. Operations
// in this class require a native ABI boundary and cannot be implemented by the
// pure-Go compatibility layer without loading/calling native code.
type PluginBoundaryError struct {
	Operation string
	Detail    string
}

func (e PluginBoundaryError) Error() string {
	if e.Detail == "" {
		return fmt.Sprintf("%s requires the native R plugin/ABI boundary", e.Operation)
	}
	return fmt.Sprintf("%s requires the native R plugin/ABI boundary: %s", e.Operation, e.Detail)
}

// ExternalPointer is the pure-Go representation of R's EXTPTRSXP container.
// Addr is intentionally an opaque Value: the compatibility runtime never
// dereferences native addresses by itself.
type ExternalPointer struct {
	mu        sync.RWMutex
	Addr      Value
	Tag       Value
	Protected Value
	Finalizer Value
	OnExit    bool
}

func (p *ExternalPointer) Address() Value {
	if p == nil {
		return nil
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Addr
}

// DynamicRuntime is the interface consumed by all translated functions.
// It intentionally centralizes operations that came from GNU R's internal C API.
type DynamicRuntime interface {
	Symbol(name string) Value
	SymbolRef(name string) Ref
	AssignSymbol(name string, v Value) Value
	Const(kind, literal string) Value
	Call(name string, args ...Value) Value
	CallIndirect(fn Value, args ...Value) Value
	Binary(op string, a, b Value) Value
	Unary(op string, v Value) Value
	Truth(v Value) bool
	Key(v Value) string
	Cast(typeName string, v Value) Value
	Field(v Value, name string) Value
	FieldRef(v Value, name string) Ref
	AssignField(v Value, name string, x Value) Value
	Index(v, i Value) Value
	IndexRef(v, i Value) Ref
	AssignIndex(v, i, x Value) Value
	Deref(v Value) Value
	DerefRef(v Value) Ref
	AssignDeref(v, x Value) Value
	AssignRef(ref Ref, x Value) Value
	RefOf(v Value) Ref
	Inc(r Ref, delta int, post bool) Value
	SizeOf(v Value) Value
	SizeOfType(typeName string) Value
	TypeValue(typeName string) Value
	List(v ...Value) Value
	Sequence(v ...Value) Value
	NewArray(size Value) Value
	NewObject() Value
	UnsupportedExpr(text string) Value
	UnsupportedStmt(text string) Value
}

// Runtime is a reusable default implementation. Register Hooks for GNU R APIs.
type Runtime struct {
	mu      sync.RWMutex
	symbols map[string]Value
	hooks   map[string]Hook

	stateMu      sync.Mutex
	protectStack []Value
	preserved    map[string]int
	warnings     []string
}

func NewRuntime() *Runtime {
	r := &Runtime{
		symbols:   map[string]Value{},
		hooks:     map[string]Hook{},
		preserved: map[string]int{},
	}
	r.symbols["R_NilValue"] = nil
	r.symbols["NULL"] = nil
	r.symbols["TRUE"] = int64(1)
	r.symbols["FALSE"] = int64(0)
	r.symbols["true"] = true
	r.symbols["false"] = false
	return r
}

var RT DynamicRuntime = NewRuntime()

func SetRuntime(rt DynamicRuntime) {
	if rt == nil {
		RT = NewRuntime()
	} else {
		RT = rt
	}
}
func (r *Runtime) Register(name string, h Hook) { r.mu.Lock(); defer r.mu.Unlock(); r.hooks[name] = h }
func (r *Runtime) SetSymbol(name string, v Value) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.symbols[name] = v
}
func (r *Runtime) Symbol(name string) Value {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if v, ok := r.symbols[name]; ok {
		return v
	}
	return symbol{name}
}

type symbol struct{ Name string }
type symRef struct {
	r    *Runtime
	name string
}

func (s symRef) Get() Value                                { return s.r.Symbol(s.name) }
func (s symRef) Set(v Value)                               { s.r.SetSymbol(s.name, v) }
func (r *Runtime) SymbolRef(name string) Ref               { return symRef{r: r, name: name} }
func (r *Runtime) AssignSymbol(name string, v Value) Value { r.SetSymbol(name, v); return v }

func (r *Runtime) Const(kind, literal string) Value {
	switch kind {
	case "string":
		if v, e := strconv.Unquote(literal); e == nil {
			return v
		}
	case "char":
		if v, e := strconv.Unquote(literal); e == nil {
			rr := []rune(v)
			if len(rr) > 0 {
				return int64(rr[0])
			}
		}
	case "float", "double":
		s := strings.TrimRight(literal, "fFlL")
		if v, e := strconv.ParseFloat(s, 64); e == nil {
			return v
		}
	default:
		s := strings.TrimSpace(literal)
		s = strings.TrimRight(s, "uUlL")
		if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
			if v, e := strconv.ParseInt(s[2:], 16, 64); e == nil {
				return v
			}
		}
		if v, e := strconv.ParseInt(s, 0, 64); e == nil {
			return v
		}
		if v, e := strconv.ParseUint(s, 0, 64); e == nil {
			return int64(v)
		}
	}
	return literal
}

func (r *Runtime) Call(name string, args ...Value) Value {
	r.mu.RLock()
	h := r.hooks[name]
	r.mu.RUnlock()
	if h != nil {
		return h(args...)
	}
	// Small C-library/R-macro-neutral core; GNU R-specific operations should be hooks.
	switch name {
	case "_":
		if len(args) > 0 {
			return args[0]
		}
		return ""
	case "strlen":
		if len(args) > 0 {
			return int64(len(asString(args[0])))
		}
	case "strcmp":
		if len(args) >= 2 {
			a, b := asString(args[0]), asString(args[1])
			if a < b {
				return int64(-1)
			}
			if a > b {
				return int64(1)
			}
			return int64(0)
		}
	case "strncmp":
		if len(args) >= 3 {
			n := int(asInt(args[2]))
			a, b := asString(args[0]), asString(args[1])
			if len(a) > n {
				a = a[:n]
			}
			if len(b) > n {
				b = b[:n]
			}
			if a < b {
				return int64(-1)
			}
			if a > b {
				return int64(1)
			}
			return int64(0)
		}
	case "streql":
		if len(args) >= 2 {
			return asString(args[0]) == asString(args[1])
		}
	case "abs":
		if len(args) > 0 {
			v := asInt(args[0])
			if v < 0 {
				v = -v
			}
			return v
		}
	case "max":
		if len(args) >= 2 {
			a, b := asFloat(args[0]), asFloat(args[1])
			if a > b {
				return a
			}
			return b
		}
	case "min":
		if len(args) >= 2 {
			a, b := asFloat(args[0]), asFloat(args[1])
			if a < b {
				return a
			}
			return b
		}
	case "free":
		return nil
	}
	panic(fmt.Sprintf("GNU R runtime operation %q is not registered", name))
}
func (r *Runtime) CallIndirect(fn Value, args ...Value) Value {
	if h, ok := fn.(Hook); ok {
		return h(args...)
	}
	if s, ok := fn.(symbol); ok {
		return r.Call(s.Name, args...)
	}
	panic(fmt.Sprintf("indirect call target is not callable: %T", fn))
}
func (r *Runtime) Truth(v Value) bool {
	switch x := v.(type) {
	case nil:
		return false
	case bool:
		return x
	case int:
		return x != 0
	case int64:
		return x != 0
	case uint64:
		return x != 0
	case float64:
		return x != 0
	case string:
		return x != ""
	default:
		return true
	}
}
func (r *Runtime) Key(v Value) string { return fmt.Sprintf("%T:%v", v, v) }
func asInt(v Value) int64 {
	switch x := v.(type) {
	case int:
		return int64(x)
	case int64:
		return x
	case uint:
		return int64(x)
	case uint64:
		return int64(x)
	case float64:
		return int64(x)
	case bool:
		if x {
			return 1
		}
		return 0
	case string:
		if z, e := strconv.ParseInt(x, 0, 64); e == nil {
			return z
		}
	}
	return 0
}
func asFloat(v Value) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case float32:
		return float64(x)
	default:
		return float64(asInt(v))
	}
}
func asString(v Value) string {
	switch x := v.(type) {
	case string:
		return x
	case []byte:
		return string(x)
	case nil:
		return ""
	default:
		return fmt.Sprint(v)
	}
}
func (r *Runtime) Binary(op string, a, b Value) Value {
	switch op {
	case "==":
		return reflect.DeepEqual(a, b)
	case "!=":
		return !reflect.DeepEqual(a, b)
	case "&&":
		return r.Truth(a) && r.Truth(b)
	case "||":
		return r.Truth(a) || r.Truth(b)
	case "+":
		if _, ok := a.(string); ok {
			return asString(a) + asString(b)
		}
		if _, ok := b.(string); ok {
			return asString(a) + asString(b)
		}
		return asFloat(a) + asFloat(b)
	case "-":
		return asFloat(a) - asFloat(b)
	case "*":
		return asFloat(a) * asFloat(b)
	case "/":
		return asFloat(a) / asFloat(b)
	case "%":
		y := asInt(b)
		if y == 0 {
			return int64(0)
		}
		return asInt(a) % y
	case "<":
		return asFloat(a) < asFloat(b)
	case "<=":
		return asFloat(a) <= asFloat(b)
	case ">":
		return asFloat(a) > asFloat(b)
	case ">=":
		return asFloat(a) >= asFloat(b)
	case "&":
		return asInt(a) & asInt(b)
	case "|":
		return asInt(a) | asInt(b)
	case "^":
		return asInt(a) ^ asInt(b)
	case "<<":
		return asInt(a) << uint(asInt(b))
	case ">>":
		return asInt(a) >> uint(asInt(b))
	}
	return opaqueOp{op, a, b}
}

type opaqueOp struct {
	Op   string
	A, B Value
}

func (r *Runtime) Unary(op string, v Value) Value {
	switch op {
	case "!":
		return !r.Truth(v)
	case "+":
		return asFloat(v)
	case "-":
		return -asFloat(v)
	case "~":
		return ^asInt(v)
	}
	return opaqueUnary{op, v}
}

type opaqueUnary struct {
	Op string
	V  Value
}

func (r *Runtime) Cast(typeName string, v Value) Value {
	t := strings.ToLower(typeName)
	if strings.Contains(t, "char") && strings.Contains(t, "*") {
		return asString(v)
	}
	if strings.Contains(t, "double") || strings.Contains(t, "float") {
		return asFloat(v)
	}
	if strings.Contains(t, "int") || strings.Contains(t, "long") || strings.Contains(t, "boolean") || strings.Contains(t, "r_xlen_t") {
		return asInt(v)
	}
	return v
}

type object map[string]Value
type objRef struct {
	o object
	k string
}

func (x objRef) Get() Value  { return x.o[x.k] }
func (x objRef) Set(v Value) { x.o[x.k] = v }
func asObject(v Value) object {
	if o, ok := v.(object); ok {
		return o
	}
	return nil
}
func (r *Runtime) Field(v Value, name string) Value {
	if o := asObject(v); o != nil {
		return o[name]
	}
	if ref, ok := v.(Ref); ok {
		return r.Field(ref.Get(), name)
	}
	return nil
}
func (r *Runtime) FieldRef(v Value, name string) Ref {
	if o := asObject(v); o != nil {
		return objRef{o: o, k: name}
	}
	o := object{}
	return objRef{o: o, k: name}
}
func (r *Runtime) AssignField(v Value, name string, x Value) Value {
	r.FieldRef(v, name).Set(x)
	return x
}

type array []Value
type arrRef struct {
	a array
	i int
}

func (x arrRef) Get() Value {
	if x.i >= 0 && x.i < len(x.a) {
		return x.a[x.i]
	}
	return nil
}
func (x arrRef) Set(v Value) {
	if x.i >= 0 && x.i < len(x.a) {
		x.a[x.i] = v
	}
}
func (r *Runtime) Index(v, i Value) Value {
	if ref, ok := v.(Ref); ok {
		v = ref.Get()
	}
	n := int(asInt(i))
	switch a := v.(type) {
	case array:
		if n >= 0 && n < len(a) {
			return a[n]
		}
	case []Value:
		if n >= 0 && n < len(a) {
			return a[n]
		}
	case string:
		rr := []rune(a)
		if n >= 0 && n < len(rr) {
			return int64(rr[n])
		}
	case []byte:
		if n >= 0 && n < len(a) {
			return int64(a[n])
		}
	}
	return nil
}
func (r *Runtime) IndexRef(v, i Value) Ref {
	if ref, ok := v.(Ref); ok {
		v = ref.Get()
	}
	n := int(asInt(i))
	if a, ok := v.(array); ok {
		return arrRef{a: a, i: n}
	}
	a := array{}
	return arrRef{a: a, i: n}
}
func (r *Runtime) AssignIndex(v, i, x Value) Value { r.IndexRef(v, i).Set(x); return x }
func (r *Runtime) Deref(v Value) Value {
	if x, ok := v.(Ref); ok {
		return x.Get()
	}
	return v
}
func (r *Runtime) DerefRef(v Value) Ref {
	if x, ok := v.(Ref); ok {
		return x
	}
	return nil
}
func (r *Runtime) AssignDeref(v, x Value) Value {
	if q, ok := v.(Ref); ok {
		q.Set(x)
	}
	return x
}
func (r *Runtime) AssignRef(ref Ref, x Value) Value { return Assign(ref, x) }
func (r *Runtime) RefOf(v Value) Ref                { return nil }
func (r *Runtime) Inc(ref Ref, delta int, post bool) Value {
	if ref == nil {
		return nil
	}
	old := ref.Get()
	nv := r.Binary("+", old, int64(delta))
	ref.Set(nv)
	if post {
		return old
	}
	return nv
}
func (r *Runtime) SizeOf(v Value) Value {
	switch x := v.(type) {
	case string:
		return int64(len(x))
	case array:
		return int64(len(x))
	case []byte:
		return int64(len(x))
	default:
		return int64(unsafeLogicalSize(x))
	}
}
func unsafeLogicalSize(v Value) int {
	if v == nil {
		return 0
	}
	return int(reflect.TypeOf(v).Size())
}
func (r *Runtime) SizeOfType(typeName string) Value { return int64(0) }
func (r *Runtime) TypeValue(typeName string) Value  { return typeName }
func (r *Runtime) List(v ...Value) Value            { return array(v) }
func (r *Runtime) Sequence(v ...Value) Value {
	if len(v) == 0 {
		return nil
	}
	return v[len(v)-1]
}
func (r *Runtime) NewArray(size Value) Value {
	n := int(asInt(size))
	if n < 0 {
		n = 0
	}
	return make(array, n)
}
func (r *Runtime) NewObject() Value { return object{} }
func (r *Runtime) UnsupportedExpr(text string) Value {
	panic("unsupported translated expression: " + text)
}
func (r *Runtime) UnsupportedStmt(text string) Value {
	panic("unsupported translated statement: " + text)
}

// Warnings returns a snapshot of warnings emitted through warning/warningcall.
func (r *Runtime) Warnings() []string {
	r.stateMu.Lock()
	defer r.stateMu.Unlock()
	return append([]string(nil), r.warnings...)
}

func runtimeMessage(args []Value) string {
	if len(args) == 0 {
		return "R runtime error"
	}
	format := asString(args[0])
	if len(args) == 1 {
		return format
	}
	// Most R error/warning formats use printf verbs that fmt.Sprintf also
	// understands. Normalize the common C length modifiers first.
	format = strings.ReplaceAll(format, "%lld", "%d")
	format = strings.ReplaceAll(format, "%llu", "%d")
	format = strings.ReplaceAll(format, "%ld", "%d")
	format = strings.ReplaceAll(format, "%lu", "%d")
	vals := make([]any, len(args)-1)
	for i := 1; i < len(args); i++ {
		vals[i-1] = args[i]
	}
	return fmt.Sprintf(format, vals...)
}

func protectionIndex(v Value) int {
	if ref, ok := v.(Ref); ok {
		v = ref.Get()
	}
	return int(asInt(v))
}

func (r *Runtime) protect(v Value) Value {
	r.stateMu.Lock()
	r.protectStack = append(r.protectStack, v)
	r.stateMu.Unlock()
	return v
}

func (r *Runtime) unprotect(n int) {
	if n <= 0 {
		return
	}
	r.stateMu.Lock()
	defer r.stateMu.Unlock()
	if n >= len(r.protectStack) {
		r.protectStack = r.protectStack[:0]
		return
	}
	r.protectStack = r.protectStack[:len(r.protectStack)-n]
}

func (r *Runtime) preserve(v Value) {
	r.stateMu.Lock()
	r.preserved[r.Key(v)]++
	r.stateMu.Unlock()
}

func (r *Runtime) release(v Value) {
	r.stateMu.Lock()
	defer r.stateMu.Unlock()
	key := r.Key(v)
	if n := r.preserved[key]; n > 1 {
		r.preserved[key] = n - 1
	} else {
		delete(r.preserved, key)
	}
}

func pluginBoundary(name string) Hook {
	return func(args ...Value) Value {
		detail := ""
		if len(args) > 0 {
			detail = asString(args[0])
		}
		panic(PluginBoundaryError{Operation: name, Detail: detail})
	}
}

// PrimitiveOp carries the table offset normally read through R's PRIMVAL macro.
type PrimitiveOp struct {
	Primitive string
	Offset    int
	Base      Value
}

// InstallCoreCompatibility registers a few frequently used macro semantics.
func (r *Runtime) InstallCoreCompatibility() {
	r.Register("PRIMVAL", func(a ...Value) Value {
		if len(a) > 0 {
			if p, ok := a[0].(PrimitiveOp); ok {
				return int64(p.Offset)
			}
		}
		return int64(0)
	})
	r.Register("R_FINITE", func(a ...Value) Value {
		if len(a) > 0 {
			return !math.IsInf(asFloat(a[0]), 0) && !math.IsNaN(asFloat(a[0]))
		}
		return false
	})

	// R's protection API is a GC-root stack. Go's GC does not require those
	// roots, but retaining the stack preserves return values and the observable
	// PROTECT/REPROTECT contract used by translated control flow.
	r.Register("PROTECT", func(a ...Value) Value {
		if len(a) == 0 {
			return nil
		}
		return r.protect(a[0])
	})
	r.Register("UNPROTECT", func(a ...Value) Value {
		if len(a) > 0 {
			r.unprotect(int(asInt(a[0])))
		}
		return nil
	})
	r.Register("PROTECT_WITH_INDEX", func(a ...Value) Value {
		if len(a) == 0 {
			return nil
		}
		r.stateMu.Lock()
		r.protectStack = append(r.protectStack, a[0])
		idx := int64(len(r.protectStack) - 1)
		r.stateMu.Unlock()
		if len(a) > 1 {
			if ref, ok := a[1].(Ref); ok {
				ref.Set(idx)
			}
		}
		return a[0]
	})
	r.Register("REPROTECT", func(a ...Value) Value {
		if len(a) == 0 {
			return nil
		}
		if len(a) > 1 {
			idx := protectionIndex(a[1])
			r.stateMu.Lock()
			if idx >= 0 && idx < len(r.protectStack) {
				r.protectStack[idx] = a[0]
			}
			r.stateMu.Unlock()
		}
		return a[0]
	})
	r.Register("R_PreserveObject", func(a ...Value) Value {
		if len(a) > 0 {
			r.preserve(a[0])
		}
		return nil
	})
	r.Register("R_ReleaseObject", func(a ...Value) Value {
		if len(a) > 0 {
			r.release(a[0])
		}
		return nil
	})

	for _, name := range []string{"error", "errorcall", "errorcall_dflt", "Rf_error", "R_MissingArgError", "R_ObjectNotFoundError"} {
		op := name
		r.Register(name, func(a ...Value) Value {
			panic(RuntimeError{Operation: op, Message: runtimeMessage(a)})
		})
	}
	for _, name := range []string{"warning", "warningcall"} {
		r.Register(name, func(a ...Value) Value {
			r.stateMu.Lock()
			r.warnings = append(r.warnings, runtimeMessage(a))
			r.stateMu.Unlock()
			return nil
		})
	}

	r.Register("R_MakeExternalPtr", func(a ...Value) Value {
		p := &ExternalPointer{}
		if len(a) > 0 {
			p.Addr = a[0]
		}
		if len(a) > 1 {
			p.Tag = a[1]
		}
		if len(a) > 2 {
			p.Protected = a[2]
		}
		return p
	})
	r.Register("R_ExternalPtrAddr", func(a ...Value) Value {
		if len(a) > 0 {
			if p, ok := a[0].(*ExternalPointer); ok {
				return p.Address()
			}
		}
		return nil
	})
	r.Register("R_ExternalPtrTag", func(a ...Value) Value {
		if len(a) > 0 {
			if p, ok := a[0].(*ExternalPointer); ok {
				p.mu.RLock()
				defer p.mu.RUnlock()
				return p.Tag
			}
		}
		return nil
	})
	r.Register("R_ExternalPtrProtected", func(a ...Value) Value {
		if len(a) > 0 {
			if p, ok := a[0].(*ExternalPointer); ok {
				p.mu.RLock()
				defer p.mu.RUnlock()
				return p.Protected
			}
		}
		return nil
	})
	r.Register("R_SetExternalPtrAddr", func(a ...Value) Value {
		if len(a) > 1 {
			if p, ok := a[0].(*ExternalPointer); ok {
				p.mu.Lock()
				p.Addr = a[1]
				p.mu.Unlock()
			}
		}
		return nil
	})
	r.Register("R_SetExternalPtrTag", func(a ...Value) Value {
		if len(a) > 1 {
			if p, ok := a[0].(*ExternalPointer); ok {
				p.mu.Lock()
				p.Tag = a[1]
				p.mu.Unlock()
			}
		}
		return nil
	})
	r.Register("R_SetExternalPtrProtected", func(a ...Value) Value {
		if len(a) > 1 {
			if p, ok := a[0].(*ExternalPointer); ok {
				p.mu.Lock()
				p.Protected = a[1]
				p.mu.Unlock()
			}
		}
		return nil
	})
	r.Register("R_ClearExternalPtr", func(a ...Value) Value {
		if len(a) > 0 {
			if p, ok := a[0].(*ExternalPointer); ok {
				p.mu.Lock()
				p.Addr = nil
				p.mu.Unlock()
			}
		}
		return nil
	})
	r.Register("R_RegisterFinalizerEx", func(a ...Value) Value {
		if len(a) >= 2 {
			if p, ok := a[0].(*ExternalPointer); ok {
				p.mu.Lock()
				p.Finalizer = a[1]
				if len(a) > 2 {
					p.OnExit = r.Truth(a[2])
				}
				p.mu.Unlock()
			}
		}
		return nil
	})
	r.Register("R_RunWeakRefFinalizer", func(a ...Value) Value {
		if len(a) == 0 {
			return nil
		}
		p, ok := a[0].(*ExternalPointer)
		if !ok {
			return nil
		}
		p.mu.RLock()
		fn := p.Finalizer
		p.mu.RUnlock()
		if h, ok := fn.(Hook); ok {
			return h(p)
		}
		return nil
	})

	// The translated LAPACK initialization hook only initializes R's native
	// dispatch table. Pure-Go translated code has no such table to initialize.
	r.Register("La_Init", func(a ...Value) Value { return nil })

	// Never fabricate success for operations whose contract is to cross a
	// native ABI/plugin boundary.
	for _, name := range []string{
		".C", ".Call", ".External", ".Fortran", "dyn.load", "dyn.unload",
		"resolveNativeRoutine", "R_doDotCall", "checkNativeType",
		"GetFullDLLPath", "AddDLL", "DeleteDLL", "R_getDllTable", "Rf_MakeDLLInfo",
		"R_FindSymbol", "R_GetCCallable", "R_RegisterCCallable", "R_RegisterCFinalizerEx",
	} {
		r.Register(name, pluginBoundary(name))
	}
}
