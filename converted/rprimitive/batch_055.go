package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_getVarsFromFrame(call, op, args, env Value) Value {
	RT.Call("checkArity", op, args)
	return RT.Call("R_getVarsFromFrame", RT.Call("CAR", args), RT.Call("CADR", args), RT.Call("CADDR", args))
}
