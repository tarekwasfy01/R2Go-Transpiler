package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_serializeToConn(call, op, args, env Value) Value {
	var (
		object  Value
		fun     Value
		ascii   Value
		wasopen Value
		version Value
		con     Value
		out     Value
		type_v  Value
		hook    Value
		cntxt   Value
		mode    Value
	)
	Assign(LocalRef(&out), RT.NewObject())
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&object), RT.Call("CAR", args))
	Assign(LocalRef(&con), RT.Call("getConnection", RT.Call("asInteger", RT.Call("CADR", args))))
	Assign(LocalRef(&ascii), RT.Call("asRbool", RT.Call("CADDR", args), call))
	if RT.Truth(RT.Binary("==", ascii, RT.Symbol("NA_LOGICAL"))) {
		Assign(LocalRef(&type_v), RT.Symbol("R_pstream_asciihex_format"))
	} else {
		if RT.Truth(ascii) {
			Assign(LocalRef(&type_v), RT.Symbol("R_pstream_ascii_format"))
		} else {
			Assign(LocalRef(&type_v), RT.Symbol("R_pstream_xdr_format"))
		}
	}
	if RT.Truth(RT.Binary("==", RT.Call("CADDDR", args), RT.Symbol("R_NilValue"))) {
		Assign(LocalRef(&version), RT.Call("defaultSerializeVersion"))
	} else {
		Assign(LocalRef(&version), RT.Call("asInteger", RT.Call("CADDDR", args)))
	}
	if RT.Truth(func() Value {
		if RT.Truth(RT.Binary("==", version, RT.Symbol("NA_INTEGER"))) {
			return true
		}
		return RT.Truth(RT.Binary("<=", version, RT.Const("int", "0")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"bad version value\"")))
	}
	if RT.Truth(RT.Binary("<", version, RT.Const("int", "2"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"cannot save to connections in version %d format\"")), version)
	}
	Assign(LocalRef(&fun), RT.Call("CAD4R", args))
	Assign(LocalRef(&hook), func() Value {
		if RT.Truth(RT.Binary("!=", fun, RT.Symbol("R_NilValue"))) {
			return RT.Symbol("CallHook")
		}
		return RT.Symbol("NULL")
	}())
	Assign(LocalRef(&wasopen), RT.Field(con, "isopen"))
	if RT.Truth(RT.Unary("!", wasopen)) {
		Assign(LocalRef(&mode), RT.NewArray(RT.Const("int", "5")))
		RT.Call("strcpy", mode, RT.Field(con, "mode"))
		RT.Call("strcpy", RT.Field(con, "mode"), func() Value {
			if RT.Truth(ascii) {
				return RT.Const("string", "\"w\"")
			}
			return RT.Const("string", "\"wb\"")
		}())
		if RT.Truth(RT.Unary("!", RT.CallIndirect(RT.Field(con, "open"), con))) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"cannot open the connection\"")))
		}
		RT.Call("strcpy", RT.Field(con, "mode"), mode)
		RT.Call("begincontext", LocalRef(&cntxt), RT.Symbol("CTXT_CCODE"), RT.Symbol("R_NilValue"), RT.Symbol("R_BaseEnv"), RT.Symbol("R_BaseEnv"), RT.Symbol("R_NilValue"), RT.Symbol("R_NilValue"))
		RT.AssignField(cntxt, "cend", RT.SymbolRef("con_cleanup"))
		RT.AssignField(cntxt, "cenddata", con)
	}
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Unary("!", ascii)) {
			return false
		}
		return RT.Truth(RT.Field(con, "text"))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"binary-mode connection required for ascii=false\"")))
	}
	if RT.Truth(RT.Unary("!", RT.Field(con, "canwrite"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"connection not open for writing\"")))
	}
	RT.Call("R_InitConnOutPStream", LocalRef(&out), con, type_v, version, hook, fun)
	RT.Call("R_Serialize", object, LocalRef(&out))
	if RT.Truth(RT.Unary("!", wasopen)) {
		RT.Call("endcontext", LocalRef(&cntxt))
		RT.CallIndirect(RT.Field(con, "close"), con)
	}
	return RT.Symbol("R_NilValue")
}
