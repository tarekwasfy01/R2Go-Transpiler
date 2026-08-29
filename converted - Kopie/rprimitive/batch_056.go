package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_getconst(call, op, args, env Value) Value {
	var (
		constBuf Value
		ans      Value
		i        Value
		n        Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&constBuf), RT.Call("CAR", args))
	Assign(LocalRef(&n), RT.Call("asInteger", RT.Call("CADR", args)))
	if RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", constBuf), RT.Symbol("VECSXP"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"constant buffer must be a generic vector\"")))
	}
	if RT.Truth(func() Value {
		if RT.Truth(RT.Binary("<", n, RT.Const("int", "0"))) {
			return true
		}
		return RT.Truth(RT.Binary(">", n, RT.Call("LENGTH", constBuf)))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"bad constant count\"")))
	}
	Assign(LocalRef(&ans), RT.Call("allocVector", RT.Symbol("VECSXP"), n))
	for Assign(LocalRef(&i), RT.Const("int", "0")); RT.Truth(RT.Binary("<", i, n)); RT.Inc(LocalRef(&i), 1, true) {
		RT.Call("SET_VECTOR_ELT", ans, i, RT.Call("VECTOR_ELT", constBuf, i))
	}
	return ans
}
