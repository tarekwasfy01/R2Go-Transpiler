package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_wrap_meta(call, op, args, env Value) Value {
	var (
		x     Value
		srt   Value
		no_na Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&x), RT.Call("CAR", args))
	Assign(LocalRef(&srt), RT.Call("asInteger", RT.Call("CADR", args)))
	Assign(LocalRef(&no_na), RT.Call("asInteger", RT.Call("CADDR", args)))
	return RT.Call("wrap_meta", x, srt, no_na)
}
