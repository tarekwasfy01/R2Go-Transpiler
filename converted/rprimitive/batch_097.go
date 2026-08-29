package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_saveToConn(call, op, args, env Value) Value {
	var (
		s       Value
		t       Value
		source  Value
		list    Value
		tmp     Value
		ascii   Value
		wasopen Value
		len     Value
		j       Value
		version Value
		ep      Value
		con     Value
		out     Value
		type_v  Value
		magic   Value
		cntxt   Value
		mode    Value
	)
	Assign(LocalRef(&out), RT.NewObject())
	Assign(LocalRef(&magic), RT.NewArray(RT.Const("int", "6")))
	RT.Call("checkArity", op, args)
	if RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", RT.Call("CAR", args)), RT.Symbol("STRSXP"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"first argument must be a character vector\"")))
	}
	Assign(LocalRef(&list), RT.Call("CAR", args))
	Assign(LocalRef(&con), RT.Call("getConnection", RT.Call("asInteger", RT.Call("CADR", args))))
	Assign(LocalRef(&ascii), RT.Call("asBool2", RT.Call("CADDR", args), call))
	if RT.Truth(RT.Binary("==", RT.Call("CADDDR", args), RT.Symbol("R_NilValue"))) {
		Assign(LocalRef(&version), RT.Call("defaultSaveVersion"))
	} else {
		Assign(LocalRef(&version), RT.Call("asInteger", RT.Call("CADDDR", args)))
	}
	if RT.Truth(func() Value {
		if RT.Truth(RT.Binary("==", version, RT.Symbol("NA_INTEGER"))) {
			return true
		}
		return RT.Truth(RT.Binary("<=", version, RT.Const("int", "0")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"version\""))
	}
	if RT.Truth(RT.Binary("<", version, RT.Const("int", "2"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"cannot save to connections in version %d format\"")), version)
	}
	Assign(LocalRef(&source), RT.Call("CAR", RT.Call("nthcdr", args, RT.Const("int", "4"))))
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Binary("!=", source, RT.Symbol("R_NilValue"))) {
			return false
		}
		return RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", source), RT.Symbol("ENVSXP")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"environment\""))
	}
	Assign(LocalRef(&ep), RT.Call("asLogical", RT.Call("CAR", RT.Call("nthcdr", args, RT.Const("int", "5")))))
	if RT.Truth(RT.Binary("==", ep, RT.Symbol("NA_LOGICAL"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"eval.promises\""))
	}
	Assign(LocalRef(&wasopen), RT.Field(con, "isopen"))
	if RT.Truth(RT.Unary("!", wasopen)) {
		Assign(LocalRef(&mode), RT.NewArray(RT.Const("int", "5")))
		RT.Call("strcpy", mode, RT.Field(con, "mode"))
		RT.Call("strcpy", RT.Field(con, "mode"), RT.Const("string", "\"wb\""))
		if RT.Truth(RT.Unary("!", RT.CallIndirect(RT.Field(con, "open"), con))) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"cannot open the connection\"")))
		}
		RT.Call("strcpy", RT.Field(con, "mode"), mode)
		RT.Call("begincontext", LocalRef(&cntxt), RT.Symbol("CTXT_CCODE"), RT.Symbol("R_NilValue"), RT.Symbol("R_BaseEnv"), RT.Symbol("R_BaseEnv"), RT.Symbol("R_NilValue"), RT.Symbol("R_NilValue"))
		RT.AssignField(cntxt, "cend", RT.SymbolRef("con_cleanup"))
		RT.AssignField(cntxt, "cenddata", con)
	}
	if RT.Truth(RT.Unary("!", RT.Field(con, "canwrite"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"connection not open for writing\"")))
	}
	RT.Call("strcpy", magic, RT.Const("string", "\"RD??\\n\""))
	if RT.Truth(ascii) {
		RT.AssignIndex(magic, RT.Const("int", "2"), RT.Const("char", "'A'"))
		Assign(LocalRef(&type_v), func() Value {
			if RT.Truth(RT.Binary("==", ascii, RT.Symbol("NA_LOGICAL"))) {
				return RT.Symbol("R_pstream_asciihex_format")
			}
			return RT.Symbol("R_pstream_ascii_format")
		}())
	} else {
		if RT.Truth(RT.Field(con, "text")) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"cannot save XDR format to a text-mode connection\"")))
		}
		RT.AssignIndex(magic, RT.Const("int", "2"), RT.Const("char", "'X'"))
		Assign(LocalRef(&type_v), RT.Symbol("R_pstream_xdr_format"))
	}
	RT.AssignIndex(magic, RT.Const("int", "3"), RT.Cast("char", RT.Binary("+", RT.Const("char", "'0'"), version)))
	if RT.Truth(RT.Field(con, "text")) {
		RT.Call("Rconn_printf", con, RT.Const("string", "\"%s\""), magic)
	} else {
		Assign(LocalRef(&len), RT.Call("strlen", magic))
		if RT.Truth(RT.Binary("!=", len, RT.CallIndirect(RT.Field(con, "write"), magic, RT.Const("int", "1"), len, con))) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"error writing to connection\"")))
		}
	}
	RT.Call("R_InitConnOutPStream", LocalRef(&out), con, type_v, version, RT.Symbol("NULL"), RT.Symbol("NULL"))
	Assign(LocalRef(&len), RT.Call("length", list))
	RT.Call("PROTECT", Assign(LocalRef(&s), RT.Call("allocList", len)))
	Assign(LocalRef(&t), s)
	for Assign(LocalRef(&j), RT.Const("int", "0")); RT.Truth(RT.Binary("<", j, len)); RT.Sequence(RT.Inc(LocalRef(&j), 1, true), Assign(LocalRef(&t), RT.Call("CDR", t))) {
		RT.Call("SET_TAG", t, RT.Call("installTrChar", RT.Call("STRING_ELT", list, j)))
		RT.Call("SETCAR", t, RT.Call("R_findVar", RT.Call("TAG", t), source))
		Assign(LocalRef(&tmp), RT.Call("R_findVar", RT.Call("TAG", t), source))
		if RT.Truth(RT.Binary("==", tmp, RT.Symbol("R_UnboundValue"))) {
			RT.Call("R_ObjectNotFoundError", RT.Call("TAG", t), RT.Symbol("R_CurrentExpression"), RT.Symbol("NULL"))
		}
		if RT.Truth(func() Value {
			if !RT.Truth(ep) {
				return false
			}
			return RT.Truth(RT.Binary("==", RT.Call("TYPEOF", tmp), RT.Symbol("PROMSXP")))
		}()) {
			RT.Call("PROTECT", tmp)
			Assign(LocalRef(&tmp), RT.Call("eval", tmp, source))
			RT.Call("UNPROTECT", RT.Const("int", "1"))
		}
		RT.Call("SETCAR", t, tmp)
	}
	RT.Call("R_Serialize", s, LocalRef(&out))
	if RT.Truth(RT.Unary("!", wasopen)) {
		RT.CallIndirect(RT.Field(con, "close"), con)
	}
	RT.Call("UNPROTECT", RT.Const("int", "1"))
	return RT.Symbol("R_NilValue")
}
