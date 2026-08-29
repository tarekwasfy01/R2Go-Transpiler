package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_getRegisteredRoutines(call, op, args, env Value) Value {
	RT.Call("checkArity", op, args)
	return RT.Call("R_getRegisteredRoutines", RT.Call("CAR", args))
}
