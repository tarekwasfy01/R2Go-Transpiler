package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_tryWrap(call, op, args, env Value) Value {
	var (
		x Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&x), RT.Call("CAR", args))
	return RT.Call("R_tryWrap", x)
}
