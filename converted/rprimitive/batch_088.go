package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_pos2env(call, op, args, rho Value) Value {
	var (
		env  Value
		pos  Value
		i    Value
		npos Value
	)
	RT.Call("checkArity", op, args)
	RT.Call("check1arg", args, call, RT.Const("string", "\"x\""))
	RT.Call("PROTECT", Assign(LocalRef(&pos), RT.Call("coerceVector", RT.Call("CAR", args), RT.Symbol("INTSXP"))))
	Assign(LocalRef(&npos), RT.Call("length", pos))
	if RT.Truth(RT.Binary("<=", npos, RT.Const("int", "0"))) {
		RT.Call("errorcall", call, RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"pos\""))
	}
	if RT.Truth(RT.Binary("==", npos, RT.Const("int", "1"))) {
		Assign(LocalRef(&env), RT.Call("pos2env", RT.Index(RT.Call("INTEGER", pos), RT.Const("int", "0")), call))
	} else {
		RT.Call("PROTECT", Assign(LocalRef(&env), RT.Call("allocVector", RT.Symbol("VECSXP"), npos)))
		for Assign(LocalRef(&i), RT.Const("int", "0")); RT.Truth(RT.Binary("<", i, npos)); RT.Inc(LocalRef(&i), 1, true) {
			RT.Call("SET_VECTOR_ELT", env, i, RT.Call("pos2env", RT.Index(RT.Call("INTEGER", pos), i), call))
		}
		RT.Call("UNPROTECT", RT.Const("int", "1"))
	}
	RT.Call("UNPROTECT", RT.Const("int", "1"))
	return env
}
