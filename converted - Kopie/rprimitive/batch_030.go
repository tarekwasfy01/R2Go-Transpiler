package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_docall(call, op, args, rho Value) Value {
	var (
		c     Value
		fun   Value
		names Value
		envir Value
		i     Value
		n     Value
		str   Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&fun), RT.Call("CAR", args))
	Assign(LocalRef(&envir), RT.Call("CADDR", args))
	Assign(LocalRef(&args), RT.Call("CADR", args))
	if RT.Truth(RT.Unary("!", func() Value {
		if RT.Truth(RT.Call("isFunction", fun)) {
			return true
		}
		return RT.Truth(func() Value {
			if !RT.Truth(RT.Call("isString", fun)) {
				return false
			}
			return RT.Truth(RT.Binary("==", RT.Call("length", fun), RT.Const("int", "1")))
		}())
	}())) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"'what' must be a function or character string\"")))
	}
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Unary("!", RT.Call("isNull", args))) {
			return false
		}
		return RT.Truth(RT.Unary("!", RT.Call("isNewList", args)))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"'%s' must be a list\"")), RT.Const("string", "\"args\""))
	}
	if RT.Truth(RT.Unary("!", RT.Call("isEnvironment", envir))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"'envir' must be an environment\"")))
	}
	Assign(LocalRef(&n), RT.Call("length", args))
	RT.Call("PROTECT", Assign(LocalRef(&names), RT.Call("getAttrib", args, RT.Symbol("R_NamesSymbol"))))
	RT.Call("PROTECT", Assign(LocalRef(&c), Assign(LocalRef(&call), RT.Call("allocLang", RT.Binary("+", n, RT.Const("int", "1"))))))
	if RT.Truth(RT.Call("isString", fun)) {
		Assign(LocalRef(&str), RT.Call("translateChar", RT.Call("STRING_ELT", fun, RT.Const("int", "0"))))
		if RT.Truth(RT.Call("streql", str, RT.Const("string", "\".Internal\""))) {
			RT.Call("error", RT.Const("string", "\"illegal usage\""))
		}
		RT.Call("SETCAR", c, RT.Call("install", str))
	} else {
		if RT.Truth(func() Value {
			if !RT.Truth(RT.Binary("==", RT.Call("TYPEOF", fun), RT.Symbol("SPECIALSXP"))) {
				return false
			}
			return RT.Truth(RT.Call("streql", RT.Call("PRIMNAME", fun), RT.Const("string", "\".Internal\"")))
		}()) {
			RT.Call("error", RT.Const("string", "\"illegal usage\""))
		}
		RT.Call("SETCAR", c, fun)
	}
	Assign(LocalRef(&c), RT.Call("CDR", c))
	for Assign(LocalRef(&i), RT.Const("int", "0")); RT.Truth(RT.Binary("<", i, n)); RT.Inc(LocalRef(&i), 1, true) {
		RT.Call("SETCAR", c, RT.Call("mkPROMISE", RT.Call("VECTOR_ELT", args, i), rho))
		RT.Call("SET_PRVALUE", RT.Call("CAR", c), RT.Call("VECTOR_ELT", args, i))
		if RT.Truth(RT.Binary("!=", RT.Call("ItemName", names, RT.Cast("int", i)), RT.Symbol("R_NilValue"))) {
			RT.Call("SET_TAG", c, RT.Call("installTrChar", RT.Call("ItemName", names, i)))
		}
		Assign(LocalRef(&c), RT.Call("CDR", c))
	}
	Assign(LocalRef(&call), RT.Call("eval", call, envir))
	RT.Call("UNPROTECT", RT.Const("int", "2"))
	return call
}
