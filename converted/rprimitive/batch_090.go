package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_readln(call, op, args, rho Value) Value {
	var (
		c      Value
		buffer Value
		bufp   Value
		ans    Value
		prompt Value
	)
	Assign(LocalRef(&buffer), RT.NewArray(RT.Symbol("MAXELTSIZE")))
	Assign(LocalRef(&bufp), buffer)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&prompt), RT.Call("CAR", args))
	if RT.Truth(RT.Binary("==", prompt, RT.Symbol("R_NilValue"))) {
		RT.AssignIndex(RT.Symbol("ConsolePrompt"), RT.Const("int", "0"), RT.Const("char", "'\\0'"))
		RT.Call("PROTECT", prompt)
	} else {
		RT.Call("PROTECT", Assign(LocalRef(&prompt), RT.Call("coerceVector", prompt, RT.Symbol("STRSXP"))))
		if RT.Truth(RT.Binary(">", RT.Call("length", prompt), RT.Const("int", "0"))) {
			RT.Call("strncpy", RT.Symbol("ConsolePrompt"), RT.Call("translateChar", RT.Call("STRING_ELT", prompt, RT.Const("int", "0"))), RT.Binary("-", RT.Symbol("CONSOLE_PROMPT_SIZE"), RT.Const("int", "1")))
			RT.AssignIndex(RT.Symbol("ConsolePrompt"), RT.Binary("-", RT.Symbol("CONSOLE_PROMPT_SIZE"), RT.Const("int", "1")), RT.Const("char", "'\\0'"))
		}
	}
	if RT.Truth(RT.Symbol("R_Interactive")) {
		for RT.Truth(func() Value {
			if RT.Truth(RT.Binary("==", Assign(LocalRef(&c), RT.Call("ConsoleGetchar")), RT.Const("char", "' '"))) {
				return true
			}
			return RT.Truth(RT.Binary("==", c, RT.Const("char", "'\\t'")))
		}()) {
		}
		if RT.Truth(func() Value {
			if !RT.Truth(RT.Binary("!=", c, RT.Const("char", "'\\n'"))) {
				return false
			}
			return RT.Truth(RT.Binary("!=", c, RT.Symbol("R_EOF")))
		}()) {
			RT.AssignDeref(RT.Inc(LocalRef(&bufp), 1, true), RT.Cast("char", c))
			for RT.Truth(func() Value {
				if !RT.Truth(RT.Binary("!=", Assign(LocalRef(&c), RT.Call("ConsoleGetchar")), RT.Const("char", "'\\n'"))) {
					return false
				}
				return RT.Truth(RT.Binary("!=", c, RT.Symbol("R_EOF")))
			}()) {
				if RT.Truth(RT.Binary(">=", bufp, RT.IndexRef(buffer, RT.Binary("-", RT.Symbol("MAXELTSIZE"), RT.Const("int", "2"))))) {
					continue
				}
				RT.AssignDeref(RT.Inc(LocalRef(&bufp), 1, true), RT.Cast("char", c))
			}
		}
		for RT.Truth(func() Value {
			if !RT.Truth(RT.Binary(">=", RT.Inc(LocalRef(&bufp), -1, false), buffer)) {
				return false
			}
			return RT.Truth(func() Value {
				if RT.Truth(RT.Binary("==", RT.Deref(bufp), RT.Const("char", "' '"))) {
					return true
				}
				return RT.Truth(RT.Binary("==", RT.Deref(bufp), RT.Const("char", "'\\t'")))
			}())
		}()) {
		}
		RT.AssignDeref(RT.Inc(LocalRef(&bufp), 1, false), RT.Const("char", "'\\0'"))
		RT.AssignIndex(RT.Symbol("ConsolePrompt"), RT.Const("int", "0"), RT.Const("char", "'\\0'"))
		Assign(LocalRef(&ans), RT.Call("mkString", buffer))
	} else {
		RT.Call("Rprintf", RT.Const("string", "\"%s\\n\""), RT.Symbol("ConsolePrompt"))
		Assign(LocalRef(&ans), RT.Call("mkString", RT.Const("string", "\"\"")))
	}
	RT.Call("UNPROTECT", RT.Const("int", "1"))
	return ans
}
