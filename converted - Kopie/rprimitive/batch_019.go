package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_builtins(call, op, args, rho Value) Value {
	var (
		ans    Value
		intern Value
		nelts  Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&intern), RT.Call("asLogical", RT.Call("CAR", args)))
	if RT.Truth(RT.Binary("==", intern, RT.Symbol("NA_INTEGER"))) {
		Assign(LocalRef(&intern), RT.Const("int", "0"))
	}
	Assign(LocalRef(&nelts), RT.Call("BuiltinSize", RT.Const("int", "1"), intern))
	RT.Call("PROTECT", Assign(LocalRef(&ans), RT.Call("allocVector", RT.Symbol("STRSXP"), nelts)))
	Assign(LocalRef(&nelts), RT.Const("int", "0"))
	RT.Call("BuiltinNames", RT.Const("int", "1"), intern, ans, LocalRef(&nelts))
	RT.Call("sortVector", ans, RT.Symbol("TRUE"))
	RT.Call("UNPROTECT", RT.Const("int", "1"))
	return ans
}
