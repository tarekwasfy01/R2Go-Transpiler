package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_curlVersion(call, op, args, rho Value) Value {
	RT.Call("checkArity", op, args)
	if RT.Truth(RT.Unary("!", RT.Symbol("initialized"))) {
		RT.Call("internet_Init")
	}
	if RT.Truth(RT.Binary(">", RT.Symbol("initialized"), RT.Const("int", "0"))) {
		return RT.CallIndirect(RT.Deref(RT.Field(RT.Symbol("ptr"), "curlVersion")), call, op, args, rho)
	} else {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"internet routines cannot be loaded\"")))
		return RT.Symbol("R_NilValue")
	}
}
