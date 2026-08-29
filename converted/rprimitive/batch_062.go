package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_invokeRestart(call, op, args, rho Value) Value {
	RT.Call("checkArity", op, args)
	RT.Call("CHECK_RESTART", RT.Call("CAR", args))
	RT.Call("invokeRestart", RT.Call("CAR", args), RT.Call("CADR", args))
	return RT.Symbol("R_NilValue")
}
