package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_External(call, op, args, env Value) Value {
	var (
		ofun   Value
		retval Value
		symbol Value
		vmax   Value
		buf    Value
		nargs  Value
		fun    Value
	)
	Assign(LocalRef(&ofun), RT.Symbol("NULL"))
	Assign(LocalRef(&symbol), RT.List(RT.Symbol("R_EXTERNAL_SYM"), RT.List(RT.Symbol("NULL")), RT.Symbol("NULL")))
	Assign(LocalRef(&vmax), RT.Call("vmaxget"))
	Assign(LocalRef(&buf), RT.NewArray(RT.Symbol("MaxSymbolBytes")))
	if RT.Truth(RT.Binary("<", RT.Call("length", args), RT.Const("int", "1"))) {
		RT.Call("errorcall", call, RT.Call("_", RT.Const("string", "\"'.NAME' is missing\"")))
	}
	RT.Call("check1arg2", args, call, RT.Const("string", "\".NAME\""))
	Assign(LocalRef(&args), RT.Call("resolveNativeRoutine", args, LocalRef(&ofun), LocalRef(&symbol), buf, RT.Symbol("NULL"), RT.Symbol("NULL"), call, env))
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Field(RT.Field(symbol, "symbol"), "external")) {
			return false
		}
		return RT.Truth(RT.Binary(">", RT.Field(RT.Field(RT.Field(symbol, "symbol"), "external"), "numArgs"), RT.Unary("-", RT.Const("int", "1"))))
	}()) {
		Assign(LocalRef(&nargs), RT.Binary("-", RT.Call("length", args), RT.Const("int", "1")))
		if RT.Truth(RT.Binary("!=", RT.Field(RT.Field(RT.Field(symbol, "symbol"), "external"), "numArgs"), nargs)) {
			RT.Call("errorcall", call, RT.Call("_", RT.Const("string", "\"Incorrect number of arguments (%d), expecting %d for '%s'\"")), nargs, RT.Field(RT.Field(RT.Field(symbol, "symbol"), "external"), "numArgs"), buf)
		}
	}
	RT.Call("R_args_enable_refcnt", args)
	if RT.Truth(RT.Binary("==", RT.Call("PRIMVAL", op), RT.Const("int", "1"))) {
		Assign(LocalRef(&fun), RT.Cast("R_ExternalRoutine2", ofun))
		Assign(LocalRef(&retval), RT.Call("fun", call, op, args, env))
	} else {
		Assign(LocalRef(&fun), RT.Cast("R_ExternalRoutine", ofun))
		Assign(LocalRef(&retval), RT.Call("fun", args))
	}
	RT.Call("R_try_clear_args_refcnt", args)
	RT.Call("vmaxset", vmax)
	return RT.Call("check_retval", call, retval)
}
