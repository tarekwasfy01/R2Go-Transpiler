//go:build windows

package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_fifo(call, op, args, env Value) Value {
	var (
		sfile Value
		sopen Value
		ans   Value
		class Value
		enc   Value
		file  Value
		open  Value
		ncon  Value
		block Value
		con   Value
	)
	Assign(LocalRef(&con), RT.Symbol("NULL"))
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&sfile), RT.Call("CAR", args))
	if RT.Truth(func() Value {
		if RT.Truth(func() Value {
			if RT.Truth(RT.Unary("!", RT.Call("isString", sfile))) {
				return true
			}
			return RT.Truth(RT.Binary("!=", RT.Call("LENGTH", sfile), RT.Const("int", "1")))
		}()) {
			return true
		}
		return RT.Truth(RT.Binary("==", RT.Call("STRING_ELT", sfile, RT.Const("int", "0")), RT.Symbol("NA_STRING")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"description\""))
	}
	if RT.Truth(RT.Binary(">", RT.Call("length", sfile), RT.Const("int", "1"))) {
		RT.Call("warning", RT.Call("_", RT.Const("string", "\"only first element of 'description' argument used\"")))
	}
	Assign(LocalRef(&file), RT.Call("translateCharFP", RT.Call("STRING_ELT", sfile, RT.Const("int", "0"))))
	Assign(LocalRef(&sopen), RT.Call("CADR", args))
	if RT.Truth(func() Value {
		if RT.Truth(RT.Unary("!", RT.Call("isString", sopen))) {
			return true
		}
		return RT.Truth(RT.Binary("!=", RT.Call("LENGTH", sopen), RT.Const("int", "1")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"open\""))
	}
	Assign(LocalRef(&block), RT.Call("asRbool", RT.Call("CADDR", args), call))
	Assign(LocalRef(&enc), RT.Call("CADDDR", args))
	if RT.Truth(func() Value {
		if RT.Truth(func() Value {
			if RT.Truth(RT.Unary("!", RT.Call("isString", enc))) {
				return true
			}
			return RT.Truth(RT.Binary("!=", RT.Call("LENGTH", enc), RT.Const("int", "1")))
		}()) {
			return true
		}
		return RT.Truth(RT.Binary(">", RT.Call("strlen", RT.Call("CHAR", RT.Call("STRING_ELT", enc, RT.Const("int", "0")))), RT.Const("int", "100")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"encoding\""))
	}
	Assign(LocalRef(&open), RT.Call("CHAR", RT.Call("STRING_ELT", sopen, RT.Const("int", "0"))))
	if RT.Truth(RT.Binary("==", RT.Call("strlen", file), RT.Const("int", "0"))) {
		if RT.Truth(RT.Unary("!", RT.Call("strlen", open))) {
			Assign(LocalRef(&open), RT.Const("string", "\"w+\""))
		}
		if RT.Truth(func() Value {
			if !RT.Truth(RT.Binary("!=", RT.Call("strcmp", open, RT.Const("string", "\"w+\"")), RT.Const("int", "0"))) {
				return false
			}
			return RT.Truth(RT.Binary("!=", RT.Call("strcmp", open, RT.Const("string", "\"w+b\"")), RT.Const("int", "0")))
		}()) {
			Assign(LocalRef(&open), RT.Const("string", "\"w+\""))
			RT.Call("warning", RT.Call("_", RT.Const("string", "\"fifo(\\\"\\\") only supports open = \\\"w+\\\" and open = \\\"w+b\\\": using the former\"")))
		}
	}
	Assign(LocalRef(&ncon), RT.Call("NextConnection"))
	Assign(LocalRef(&con), RT.AssignIndex(RT.Symbol("Connections"), ncon, RT.Call("newfifo", file, func() Value {
		if RT.Truth(RT.Call("strlen", open)) {
			return open
		}
		return RT.Const("string", "\"r\"")
	}())))
	RT.AssignField(con, "blocking", block)
	RT.Call("strncpy", RT.Field(con, "encname"), RT.Call("CHAR", RT.Call("STRING_ELT", enc, RT.Const("int", "0"))), RT.Const("int", "100"))
	RT.AssignIndex(RT.Field(con, "encname"), RT.Binary("-", RT.Const("int", "100"), RT.Const("int", "1")), RT.Const("char", "'\\0'"))
	RT.AssignField(con, "ex_ptr", RT.Call("PROTECT", RT.Call("R_MakeExternalPtr", RT.Field(con, "id"), RT.Call("install", RT.Const("string", "\"connection\"")), RT.Symbol("R_NilValue"))))
	if RT.Truth(RT.Call("strlen", open)) {
		RT.Call("checked_open", ncon)
	}
	RT.Call("PROTECT", Assign(LocalRef(&ans), RT.Call("ScalarInteger", ncon)))
	RT.Call("PROTECT", Assign(LocalRef(&class), RT.Call("allocVector", RT.Symbol("STRSXP"), RT.Const("int", "2"))))
	RT.Call("SET_STRING_ELT", class, RT.Const("int", "0"), RT.Call("mkChar", RT.Const("string", "\"fifo\"")))
	RT.Call("SET_STRING_ELT", class, RT.Const("int", "1"), RT.Call("mkChar", RT.Const("string", "\"connection\"")))
	RT.Call("classgets", ans, class)
	RT.Call("setAttrib", ans, RT.Symbol("R_ConnIdSymbol"), RT.Field(con, "ex_ptr"))
	RT.Call("R_RegisterCFinalizerEx", RT.Field(con, "ex_ptr"), RT.Symbol("conFinalizer"), RT.Symbol("FALSE"))
	RT.Call("UNPROTECT", RT.Const("int", "3"))
	return ans
}
