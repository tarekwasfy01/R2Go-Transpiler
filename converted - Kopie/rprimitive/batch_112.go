package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_traceback(call, op, args, rho Value) Value {
	var (
		skip Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&skip), RT.Call("asInteger", RT.Call("CAR", args)))
	if RT.Truth(func() Value {
		if RT.Truth(RT.Binary("==", skip, RT.Symbol("NA_INTEGER"))) {
			return true
		}
		return RT.Truth(RT.Binary("<", skip, RT.Const("int", "0")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' value\"")), RT.Const("string", "\"skip\""))
	}
	return RT.Call("R_GetTraceback", skip)
}
