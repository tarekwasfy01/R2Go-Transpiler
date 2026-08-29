package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_disassemble(call, op, args, rho Value) Value {
	var (
		code Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&code), RT.Call("CAR", args))
	if RT.Truth(RT.Unary("!", RT.Call("isByteCode", code))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"argument is not a byte code object\"")))
	}
	return RT.Call("disassemble", code)
}
