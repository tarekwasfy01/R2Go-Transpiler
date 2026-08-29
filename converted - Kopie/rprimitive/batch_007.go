package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_addTryHandlers(call, op, args, rho Value) Value {
	RT.Call("checkArity", op, args)
	if RT.Truth(func() Value {
		if RT.Truth(RT.Binary("==", RT.Symbol("R_GlobalContext"), RT.Symbol("R_ToplevelContext"))) {
			return true
		}
		return RT.Truth(RT.Unary("!", RT.Binary("&", RT.Field(RT.Symbol("R_GlobalContext"), "callflag"), RT.Symbol("CTXT_FUNCTION"))))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"not in a try context\"")))
	}
	RT.Call("SET_RESTART_BIT_ON", RT.Field(RT.Symbol("R_GlobalContext"), "callflag"))
	RT.Call("R_InsertRestartHandlers", RT.Symbol("R_GlobalContext"), RT.Const("string", "\"tryRestart\""))
	return RT.Symbol("R_NilValue")
}
