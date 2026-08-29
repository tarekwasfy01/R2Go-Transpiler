package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_is_builtin_internal(call, op, args, rho Value) Value {
	var (
		symbol Value
		i      Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&symbol), RT.Call("CAR", args))
	if RT.Truth(RT.Unary("!", RT.Call("isSymbol", symbol))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid symbol\"")))
	}
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Binary("!=", Assign(LocalRef(&i), RT.Call("INTERNAL", symbol)), RT.Symbol("R_NilValue"))) {
			return false
		}
		return RT.Truth(RT.Binary("==", RT.Call("TYPEOF", i), RT.Symbol("BUILTINSXP")))
	}()) {
		return RT.Symbol("R_TrueValue")
	} else {
		return RT.Symbol("R_FalseValue")
	}
}
