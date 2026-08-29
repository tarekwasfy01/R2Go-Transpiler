package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_putconst(call, op, args, env Value) Value {
	var (
		constBuf   Value
		x          Value
		i          Value
		constCount Value
		y          Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&constBuf), RT.Call("CAR", args))
	if RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", constBuf), RT.Symbol("VECSXP"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"constant buffer must be a generic vector\"")))
	}
	Assign(LocalRef(&constCount), RT.Call("asInteger", RT.Call("CADR", args)))
	if RT.Truth(func() Value {
		if RT.Truth(RT.Binary("<", constCount, RT.Const("int", "0"))) {
			return true
		}
		return RT.Truth(RT.Binary(">=", constCount, RT.Call("LENGTH", constBuf)))
	}()) {
		RT.Call("error", RT.Const("string", "\"bad constCount value\""))
	}
	Assign(LocalRef(&x), RT.Call("CADDR", args))
	for Assign(LocalRef(&i), RT.Const("int", "0")); RT.Truth(RT.Binary("<", i, constCount)); RT.Inc(LocalRef(&i), 1, true) {
		Assign(LocalRef(&y), RT.Call("VECTOR_ELT", constBuf, i))
		if RT.Truth(func() Value {
			if RT.Truth(RT.Binary("==", x, y)) {
				return true
			}
			return RT.Truth(RT.Call("R_compute_identical", x, y, RT.Const("int", "16")))
		}()) {
			return RT.Call("ScalarInteger", i)
		}
	}
	RT.Call("SET_VECTOR_ELT", constBuf, constCount, x)
	return RT.Call("ScalarInteger", constCount)
}
