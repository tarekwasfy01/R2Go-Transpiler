package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_as_environment(call, op, args, rho Value) Value {
	var (
		arg       Value
		ans       Value
		dot_xData Value
		val       Value
	)
	Assign(LocalRef(&arg), RT.Call("CAR", args))
	RT.Call("checkArity", op, args)
	RT.Call("check1arg", args, call, RT.Const("string", "\"x\""))
	if RT.Truth(RT.Call("isEnvironment", arg)) {
		return arg
	}
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Call("isObject", arg)) {
			return false
		}
		return RT.Truth(RT.Call("DispatchOrEval", call, op, RT.Const("string", "\"as.environment\""), args, rho, LocalRef(&ans), RT.Const("int", "0"), RT.Const("int", "1")))
	}()) {
		return ans
	}
	switch RT.Key(RT.Call("TYPEOF", arg)) {
	case RT.Key(RT.Symbol("STRSXP")):
		return RT.Call("matchEnvir", call, RT.Call("translateChar", RT.Call("asChar", arg)))
	case RT.Key(RT.Symbol("REALSXP")), RT.Key(RT.Symbol("INTSXP")):
		return do_pos2env(call, op, args, rho)
	case RT.Key(RT.Symbol("NILSXP")):
		RT.Call("errorcall", call, RT.Call("_", RT.Const("string", "\"using 'as.environment(NULL)' is defunct\"")))
		return RT.Symbol("R_BaseEnv")
	case RT.Key(RT.Symbol("OBJSXP")):
		Assign(LocalRef(&dot_xData), RT.Call("R_getS4DataSlot", arg, RT.Symbol("ENVSXP")))
		if RT.Truth(RT.Unary("!", RT.Call("isEnvironment", dot_xData))) {
			RT.Call("errorcall", call, RT.Call("_", RT.Const("string", "\"S4 object does not extend class \\\"environment\\\"\"")))
		} else {
			return dot_xData
		}
		fallthrough
	case RT.Key(RT.Symbol("VECSXP")):
		RT.Call("PROTECT", Assign(LocalRef(&call), RT.Call("lang4", RT.Call("install", RT.Const("string", "\"list2env\"")), arg, RT.Symbol("R_NilValue"), RT.Symbol("R_EmptyEnv"))))
		Assign(LocalRef(&val), RT.Call("eval", call, rho))
		RT.Call("UNPROTECT", RT.Const("int", "1"))
		return val
	default:
		RT.Call("errorcall", call, RT.Call("_", RT.Const("string", "\"invalid object for 'as.environment'\"")))
		return RT.Symbol("R_NilValue")
	}
}
