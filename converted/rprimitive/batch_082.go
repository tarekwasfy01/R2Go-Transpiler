package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_mkcode(call, op, args, rho Value) Value {
	var (
		bytes  Value
		consts Value
		ans    Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&bytes), RT.Call("CAR", args))
	Assign(LocalRef(&consts), RT.Call("CADR", args))
	Assign(LocalRef(&ans), RT.Call("PROTECT", RT.Call("CONS", RT.Call("R_bcEncode", bytes), consts)))
	RT.Call("SET_TYPEOF", ans, RT.Symbol("BCODESXP"))
	RT.Call("R_registerBC", bytes, ans)
	RT.Call("UNPROTECT", RT.Const("int", "1"))
	return ans
}
