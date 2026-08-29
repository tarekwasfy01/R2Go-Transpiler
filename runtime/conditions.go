package runtime

import (
	"fmt"
	"strings"

	"r2go/syntax"
)

func (c *Context) makeCondition(class string, args []syntax.Argument, env *Environment, call syntax.Expr) *ConditionValue {
	parts := make([]string, 0, len(args))
	for _, arg := range args {
		if arg.Value == nil {
			continue
		}
		value, err := c.Eval(arg.Value, env)
		if err != nil {
			parts = append(parts, err.Error())
			continue
		}
		if text, ok := value.(*CharacterVector); ok {
			parts = append(parts, text.Data...)
			continue
		}
		parts = append(parts, value.String())
	}
	classes := []string{class, "condition"}
	if class == "error" {
		classes = []string{"simpleError", "error", "condition"}
	}
	if class == "warning" {
		classes = []string{"simpleWarning", "warning", "condition"}
	}
	return &ConditionValue{Classes: classes, Message: strings.Join(parts, ""), Call: call}
}

func (c *Context) tryCatch(args []syntax.Argument, env *Environment) (Value, error) {
	if len(args) == 0 || args[0].Value == nil {
		return nil, fmt.Errorf("tryCatch requires an expression")
	}
	var errorHandler syntax.Expr
	var finally syntax.Expr
	for _, arg := range args[1:] {
		switch arg.Name {
		case "error":
			errorHandler = arg.Value
		case "finally":
			finally = arg.Value
		}
	}
	value, caught := c.Eval(args[0].Value, env)
	if finally != nil {
		if _, err := c.Eval(finally, env); err != nil {
			return nil, err
		}
	}
	if caught == nil {
		return value, nil
	}
	condition, ok := caught.(*ConditionValue)
	if !ok || errorHandler == nil {
		return nil, caught
	}
	handlerValue, err := c.Eval(errorHandler, env)
	if err != nil {
		return nil, err
	}
	handler, ok := handlerValue.(*Closure)
	if !ok {
		return nil, fmt.Errorf("error handler is not a function")
	}
	return c.callClosureWithValue(handler, condition)
}

func (c *Context) callClosureWithValue(fn *Closure, value Value) (Value, error) {
	frame := NewEnvironment(fn.Env)
	bound := false
	for _, parameter := range fn.Parameters {
		if parameter.Name == "..." {
			frame.Bind("...", &DotsBinding{})
			continue
		}
		if !bound {
			frame.Set(parameter.Name, value)
			bound = true
			continue
		}
		if parameter.Default != nil {
			frame.Bind(parameter.Name, &Promise{Expr: parameter.Default, Env: frame})
		} else {
			frame.Bind(parameter.Name, &MissingBinding{Name: parameter.Name})
		}
	}
	result, err := c.Eval(fn.Body, frame)
	if signal, ok := err.(*control); ok && signal.kind == "return" {
		return signal.value, nil
	}
	return result, err
}

func (v *ConditionValue) Error() string { return v.Message }
