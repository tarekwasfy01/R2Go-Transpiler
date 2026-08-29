package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_unserializeFromConn(call, op, args, env Value) Value {
	var (
		in      Value
		con     Value
		fun     Value
		ans     Value
		hook    Value
		wasopen Value
		cntxt   Value
		mode    Value
	)
	Assign(LocalRef(&in), RT.NewObject())
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&con), RT.Call("getConnection", RT.Call("asInteger", RT.Call("CAR", args))))
	Assign(LocalRef(&wasopen), RT.Field(con, "isopen"))
	if RT.Truth(RT.Unary("!", wasopen)) {
		Assign(LocalRef(&mode), RT.NewArray(RT.Const("int", "5")))
		RT.Call("strcpy", mode, RT.Field(con, "mode"))
		RT.Call("strcpy", RT.Field(con, "mode"), RT.Const("string", "\"rb\""))
		if RT.Truth(RT.Unary("!", RT.CallIndirect(RT.Field(con, "open"), con))) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"cannot open the connection\"")))
		}
		RT.Call("strcpy", RT.Field(con, "mode"), mode)
		RT.Call("begincontext", LocalRef(&cntxt), RT.Symbol("CTXT_CCODE"), RT.Symbol("R_NilValue"), RT.Symbol("R_BaseEnv"), RT.Symbol("R_BaseEnv"), RT.Symbol("R_NilValue"), RT.Symbol("R_NilValue"))
		RT.AssignField(cntxt, "cend", RT.SymbolRef("con_cleanup"))
		RT.AssignField(cntxt, "cenddata", con)
	}
	if RT.Truth(RT.Unary("!", RT.Field(con, "canread"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"connection not open for reading\"")))
	}
	Assign(LocalRef(&fun), func() Value {
		if RT.Truth(RT.Binary("==", RT.Call("PRIMVAL", op), RT.Const("int", "0"))) {
			return RT.Call("CADR", args)
		}
		return RT.Symbol("R_NilValue")
	}())
	Assign(LocalRef(&hook), func() Value {
		if RT.Truth(RT.Binary("!=", fun, RT.Symbol("R_NilValue"))) {
			return RT.Symbol("CallHook")
		}
		return RT.Symbol("NULL")
	}())
	RT.Call("R_InitConnInPStream", LocalRef(&in), con, RT.Symbol("R_pstream_any_format"), hook, fun)
	Assign(LocalRef(&ans), func() Value {
		if RT.Truth(RT.Binary("==", RT.Call("PRIMVAL", op), RT.Const("int", "0"))) {
			return RT.Call("R_Unserialize", LocalRef(&in))
		}
		return RT.Call("R_SerializeInfo", LocalRef(&in))
	}())
	if RT.Truth(RT.Unary("!", wasopen)) {
		RT.Call("PROTECT", ans)
		RT.Call("endcontext", LocalRef(&cntxt))
		RT.CallIndirect(RT.Field(con, "close"), con)
		RT.Call("UNPROTECT", RT.Const("int", "1"))
	}
	return RT.Call("checkNotPromise", ans)
}
