package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_sockconn(call, op, args, env Value) Value {
	var (
		scmd     Value
		sopen    Value
		ans      Value
		class    Value
		enc      Value
		host     Value
		open     Value
		ncon     Value
		port     Value
		server   Value
		timeout  Value
		serverfd Value
		options  Value
		blocking Value
		con      Value
		scon     Value
		sOpts    Value
		i        Value
		n        Value
		opt      Value
	)
	Assign(LocalRef(&options), RT.Const("int", "0"))
	Assign(LocalRef(&con), RT.Symbol("NULL"))
	Assign(LocalRef(&scon), RT.Symbol("NULL"))
	RT.Call("checkArity", op, args)
	if RT.Truth(RT.Binary("==", RT.Call("PRIMVAL", op), RT.Const("int", "0"))) {
		Assign(LocalRef(&scmd), RT.Call("CAR", args))
		if RT.Truth(func() Value {
			if RT.Truth(RT.Unary("!", RT.Call("isString", scmd))) {
				return true
			}
			return RT.Truth(RT.Binary("!=", RT.Call("LENGTH", scmd), RT.Const("int", "1")))
		}()) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"host\""))
		}
		Assign(LocalRef(&host), RT.Call("translateCharFP", RT.Call("STRING_ELT", scmd, RT.Const("int", "0"))))
		Assign(LocalRef(&args), RT.Call("CDR", args))
		Assign(LocalRef(&port), RT.Call("asInteger", RT.Call("CAR", args)))
		if RT.Truth(func() Value {
			if RT.Truth(RT.Binary("==", port, RT.Symbol("NA_INTEGER"))) {
				return true
			}
			return RT.Truth(RT.Binary("<", port, RT.Const("int", "0")))
		}()) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"port\""))
		}
		Assign(LocalRef(&args), RT.Call("CDR", args))
		Assign(LocalRef(&server), RT.Call("asLogical", RT.Call("CAR", args)))
		if RT.Truth(RT.Binary("==", server, RT.Symbol("NA_LOGICAL"))) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"server\""))
		}
		Assign(LocalRef(&serverfd), RT.Unary("-", RT.Const("int", "1")))
	} else {
		Assign(LocalRef(&scon), RT.Field(RT.Call("getConnectionCheck", RT.Call("CAR", args), RT.Const("string", "\"servsockconn\""), RT.Const("string", "\"socket\"")), "private"))
		Assign(LocalRef(&port), RT.Field(scon, "port"))
		Assign(LocalRef(&server), RT.Const("int", "1"))
		Assign(LocalRef(&host), RT.Const("string", "\"localhost\""))
		Assign(LocalRef(&serverfd), RT.Field(scon, "fd"))
	}
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&blocking), RT.Call("asRbool", RT.Call("CAR", args), call))
	if RT.Truth(RT.Binary("==", blocking, RT.Symbol("NA_LOGICAL"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"blocking\""))
	}
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&sopen), RT.Call("CAR", args))
	if RT.Truth(func() Value {
		if RT.Truth(RT.Unary("!", RT.Call("isString", sopen))) {
			return true
		}
		return RT.Truth(RT.Binary("!=", RT.Call("LENGTH", sopen), RT.Const("int", "1")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"open\""))
	}
	Assign(LocalRef(&open), RT.Call("CHAR", RT.Call("STRING_ELT", sopen, RT.Const("int", "0"))))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&enc), RT.Call("CAR", args))
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
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&timeout), RT.Call("asInteger", RT.Call("CAR", args)))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	if RT.Truth(RT.Call("isString", RT.Call("CAR", args))) {
		Assign(LocalRef(&sOpts), RT.Call("CAR", args))
		Assign(LocalRef(&i), RT.Const("int", "0"))
		Assign(LocalRef(&n), RT.Call("LENGTH", sOpts))
		for RT.Truth(RT.Binary("<", i, n)) {
			Assign(LocalRef(&opt), RT.Call("CHAR", RT.Call("STRING_ELT", sOpts, i)))
			if RT.Truth(RT.Unary("!", RT.Call("strcmp", RT.Const("string", "\"no-delay\""), opt))) {
				Assign(LocalRef(&options), RT.Binary("|", options, RT.Symbol("RSC_SET_TCP_NODELAY")))
			}
			RT.Inc(LocalRef(&i), 1, true)
		}
	}
	Assign(LocalRef(&ncon), RT.Call("NextConnection"))
	Assign(LocalRef(&con), RT.Call("R_newsock", host, port, server, serverfd, open, timeout, options))
	RT.AssignIndex(RT.Symbol("Connections"), ncon, con)
	RT.AssignField(con, "blocking", blocking)
	RT.Call("strncpy", RT.Field(con, "encname"), RT.Call("CHAR", RT.Call("STRING_ELT", enc, RT.Const("int", "0"))), RT.Const("int", "100"))
	RT.AssignIndex(RT.Field(con, "encname"), RT.Binary("-", RT.Const("int", "100"), RT.Const("int", "1")), RT.Const("char", "'\\0'"))
	RT.AssignField(con, "ex_ptr", RT.Call("PROTECT", RT.Call("R_MakeExternalPtr", RT.Field(con, "id"), RT.Call("install", RT.Const("string", "\"connection\"")), RT.Symbol("R_NilValue"))))
	if RT.Truth(RT.Call("strlen", open)) {
		RT.Call("checked_open", ncon)
	}
	RT.Call("PROTECT", Assign(LocalRef(&ans), RT.Call("ScalarInteger", ncon)))
	RT.Call("PROTECT", Assign(LocalRef(&class), RT.Call("allocVector", RT.Symbol("STRSXP"), RT.Const("int", "2"))))
	RT.Call("SET_STRING_ELT", class, RT.Const("int", "0"), RT.Call("mkChar", RT.Const("string", "\"sockconn\"")))
	RT.Call("SET_STRING_ELT", class, RT.Const("int", "1"), RT.Call("mkChar", RT.Const("string", "\"connection\"")))
	RT.Call("classgets", ans, class)
	RT.Call("setAttrib", ans, RT.Symbol("R_ConnIdSymbol"), RT.Field(con, "ex_ptr"))
	RT.Call("R_RegisterCFinalizerEx", RT.Field(con, "ex_ptr"), RT.Symbol("conFinalizer"), RT.Symbol("FALSE"))
	RT.Call("UNPROTECT", RT.Const("int", "3"))
	return ans
}
