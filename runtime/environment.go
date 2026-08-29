package runtime

import (
	"fmt"
	"r2go/syntax"
	"sort"
	"sync"
)

type Binding interface{ Force(*Context) (Value, error) }
type Eager struct{ Value Value }

func (e *Eager) Force(*Context) (Value, error) { return e.Value, nil }

type Promise struct {
	Expr syntax.Expr
	Env  *Environment
	// Supplied distinguishes an actual argument from a default expression.
	// R's missing() observes this distinction without forcing the promise.
	Supplied bool
	mu       sync.Mutex
	forced   bool
	forcing  bool
	value    Value
	err      error
}

func (p *Promise) IsMissing() bool { return !p.Supplied }

type MissingBinding struct{ Name string }

func (m *MissingBinding) Force(*Context) (Value, error) {
	return nil, fmt.Errorf("argument %q is missing, with no default", m.Name)
}
func (*MissingBinding) IsMissing() bool { return true }

type DotsBinding struct {
	Arguments []ActualArgument
}

type ActualArgument struct {
	Argument syntax.Argument
	Env      *Environment
	// EagerValue lets higher-order functions feed an already evaluated item
	// through the same matcher as lazy R expressions in additional arguments.
	EagerValue Value
	HasValue   bool
}

func (d *DotsBinding) Force(ctx *Context) (Value, error) {
	values := make([]Value, len(d.Arguments))
	names := make([]string, len(d.Arguments))
	for i, actual := range d.Arguments {
		arg := actual.Argument
		names[i] = arg.Name
		if actual.HasValue {
			values[i] = actual.EagerValue
			continue
		}
		if arg.Value == nil {
			return nil, fmt.Errorf("argument %d is empty", i+1)
		}
		value, err := ctx.Eval(arg.Value, actual.Env)
		if err != nil {
			return nil, err
		}
		values[i] = value
	}
	return &List{Data: values, Names: names}, nil
}
func (*DotsBinding) IsMissing() bool { return false }

func (p *Promise) Force(ctx *Context) (Value, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.forced {
		return p.value, p.err
	}
	if p.forcing {
		return nil, fmt.Errorf("promise already under evaluation")
	}
	p.forcing = true
	p.value, p.err = ctx.Eval(p.Expr, p.Env)
	p.forcing = false
	p.forced = true
	return p.value, p.err
}

type Environment struct {
	Parent   *Environment
	bindings map[string]Binding
}

// Environments can flow through the translated C compatibility bridge as
// ordinary host values while retaining their binding identity.
func (*Environment) Kind() Kind     { return EnvironmentKind }
func (*Environment) String() string { return "<environment>" }

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{Parent: parent, bindings: map[string]Binding{}}
}
func (e *Environment) Bind(name string, b Binding) { e.bindings[name] = b }
func (e *Environment) Set(name string, v Value)    { e.Bind(name, &Eager{Value: v}) }
func (e *Environment) Local(name string) (Binding, bool) {
	binding, ok := e.bindings[name]
	return binding, ok
}
func (e *Environment) Find(name string) (*Environment, Binding, bool) {
	for x := e; x != nil; x = x.Parent {
		if b, ok := x.bindings[name]; ok {
			return x, b, true
		}
	}
	return nil, nil, false
}
func (e *Environment) Get(ctx *Context, name string) (Value, error) {
	_, b, ok := e.Find(name)
	if !ok {
		return nil, fmt.Errorf("object %q not found", name)
	}
	return b.Force(ctx)
}
func (e *Environment) SuperSet(name string, v Value) {
	if owner, _, ok := e.Parent.Find(name); ok {
		owner.Set(name, v)
	} else {
		root := e
		for root.Parent != nil {
			root = root.Parent
		}
		root.Set(name, v)
	}
}
func (e *Environment) Names(all bool) []string {
	o := make([]string, 0, len(e.bindings))
	for n := range e.bindings {
		if all || len(n) == 0 || n[0] != '.' {
			o = append(o, n)
		}
	}
	sort.Strings(o)
	return o
}
func (e *Environment) Remove(name string) bool {
	if _, ok := e.bindings[name]; ok {
		delete(e.bindings, name)
		return true
	}
	return false
}
