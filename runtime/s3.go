package runtime

import (
	"fmt"
	"r2go/syntax"
)

func classNames(v Value) []string {
	if c, ok := Attributes(v)["class"].(*CharacterVector); ok {
		return append([]string(nil), c.Data...)
	}
	return []string{defaultClass(v)}
}
func (c *Context) findS3MethodFrom(generic string, object Value, env *Environment, start int) (*Closure, int, []string, bool, error) {
	classes := append(classNames(object), "default")
	for i := start; i < len(classes); i++ {
		name := generic + "." + classes[i]
		_, b, ok := env.Find(name)
		if !ok {
			continue
		}
		v, e := b.Force(c)
		if e != nil {
			return nil, 0, nil, false, e
		}
		fn, ok := v.(*Closure)
		if !ok {
			return nil, 0, nil, false, fmt.Errorf("S3 method %s is not a function", name)
		}
		return fn, i, classes, true, nil
	}
	return nil, 0, classes, false, nil
}
func (c *Context) useMethod(args []syntax.Argument, env *Environment) (Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("UseMethod requires a generic name")
	}
	g, e := c.Eval(args[0].Value, env)
	if e != nil {
		return nil, e
	}
	text, ok := g.(*CharacterVector)
	if !ok || len(text.Data) != 1 {
		return nil, fmt.Errorf("generic name must be a string")
	}
	var object Value
	if len(args) > 1 {
		object, e = c.Eval(args[1].Value, env)
	} else {
		object, e = env.Get(c, ".S3Object")
	}
	if e != nil {
		return nil, e
	}
	method, index, classes, found, e := c.findS3MethodFrom(text.Data[0], object, env, 0)
	if e != nil {
		return nil, e
	}
	if !found {
		return nil, fmt.Errorf("no applicable method for %q applied to an object of class %q", text.Data[0], classNames(object)[0])
	}
	var dots []ActualArgument
	if _, b, ok := env.Find("..."); ok {
		if d, ok := b.(*DotsBinding); ok {
			dots = d.Arguments
		}
	}
	return c.callS3Method(method, object, text.Data[0], classes, index, dots)
}
func (c *Context) dispatchS3(generic string, args []syntax.Argument, env *Environment) (Value, bool, error) {
	if len(args) == 0 || args[0].Value == nil {
		return nil, false, nil
	}
	object, e := c.Eval(args[0].Value, env)
	if e != nil {
		return nil, true, e
	}
	method, index, classes, ok, e := c.findS3MethodFrom(generic, object, env, 0)
	if e != nil || !ok {
		return nil, ok, e
	}
	dots := make([]ActualArgument, 0, len(args)-1)
	for _, a := range args[1:] {
		dots = append(dots, ActualArgument{Argument: a, Env: env})
	}
	v, e := c.callS3Method(method, object, generic, classes, index, dots)
	return v, true, e
}

func (c *Context) callS3Method(fn *Closure, object Value, generic string, classes []string, index int, dots []ActualArgument) (Value, error) {
	frame := NewEnvironment(fn.Env)
	bound := false
	for _, p := range fn.Parameters {
		if p.Name == "..." {
			frame.Bind("...", &DotsBinding{Arguments: dots})
			continue
		}
		if !bound {
			frame.Set(p.Name, object)
			bound = true
		} else if p.Default != nil {
			frame.Bind(p.Name, &Promise{Expr: p.Default, Env: frame})
		} else {
			frame.Bind(p.Name, &MissingBinding{Name: p.Name})
		}
	}
	frame.Set(".S3Object", object)
	frame.Set(".Generic", &CharacterVector{Data: []string{generic}})
	frame.Set(".Class", &CharacterVector{Data: classes})
	frame.Set(".S3Index", &IntegerVector{Data: []int64{int64(index)}})
	v, e := c.Eval(fn.Body, frame)
	if x, ok := e.(*control); ok && x.kind == "return" {
		return x.value, nil
	}
	return v, e
}

func (c *Context) nextMethod(args []syntax.Argument, env *Environment) (Value, error) {
	g, e := env.Get(c, ".Generic")
	if e != nil {
		return nil, fmt.Errorf("NextMethod called outside a method")
	}
	generic := g.(*CharacterVector).Data[0]
	object, e := env.Get(c, ".S3Object")
	if e != nil {
		return nil, e
	}
	iv, e := env.Get(c, ".S3Index")
	if e != nil {
		return nil, e
	}
	start := int(iv.(*IntegerVector).Data[0]) + 1
	method, index, classes, ok, e := c.findS3MethodFrom(generic, object, env, start)
	if e != nil {
		return nil, e
	}
	if !ok {
		return nil, fmt.Errorf("no next method available")
	}
	var dots []ActualArgument
	if _, b, ok := env.Find("..."); ok {
		if d, ok := b.(*DotsBinding); ok {
			dots = d.Arguments
		}
	}
	return c.callS3Method(method, object, generic, classes, index, dots)
}
