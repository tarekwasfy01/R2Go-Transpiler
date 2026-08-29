package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_list2env(call, op, args, rho Value) Value {
	var (
		x     Value
		xnms  Value
		envir Value
		n     Value
		i     Value
		name  Value
	)
	RT.Call("checkArity", op, args)
	if RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", RT.Call("CAR", args)), RT.Symbol("VECSXP"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"first argument must be a named list\"")))
	}
	Assign(LocalRef(&x), RT.Call("CAR", args))
	Assign(LocalRef(&n), RT.Call("LENGTH", x))
	Assign(LocalRef(&xnms), RT.Call("getAttrib", x, RT.Symbol("R_NamesSymbol")))
	RT.Call("PROTECT", xnms)
	if RT.Truth(func() Value {
		if !RT.Truth(n) {
			return false
		}
		return RT.Truth(func() Value {
			if RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", xnms), RT.Symbol("STRSXP"))) {
				return true
			}
			return RT.Truth(RT.Binary("!=", RT.Call("LENGTH", xnms), n))
		}())
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"names(x) must be a character vector of the same length as x\"")))
	}
	Assign(LocalRef(&envir), RT.Call("CADR", args))
	if RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", envir), RT.Symbol("ENVSXP"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"'envir' argument must be an environment\"")))
	}
	for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, n)); RT.Inc(LocalRef(&i), 1, true) {
		Assign(LocalRef(&name), RT.Call("installTrChar", RT.Call("STRING_ELT", xnms, i)))
		RT.Call("defineVar", name, RT.Call("lazy_duplicate", RT.Call("VECTOR_ELT", x, i)), envir)
	}
	RT.Call("UNPROTECT", RT.Const("int", "1"))
	return envir
}
