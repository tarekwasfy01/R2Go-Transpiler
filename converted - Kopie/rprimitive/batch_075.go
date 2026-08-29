package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_lockEnv(call, op, args, rho Value) Value {
	var (
		frame    Value
		bindings Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&frame), RT.Call("CAR", args))
	Assign(LocalRef(&bindings), RT.Call("asRbool", RT.Call("CADR", args), call))
	RT.Call("R_LockEnvironment", frame, bindings)
	return RT.Symbol("R_NilValue")
}
