package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_bodyCode(call, op, args, rho Value) Value {
	var (
		bc Value
	)
	RT.Call("checkArity", op, args)
	if RT.Truth(RT.Binary("==", RT.Call("TYPEOF", RT.Call("CAR", args)), RT.Symbol("CLOSXP"))) {
		Assign(LocalRef(&bc), RT.Call("BODY", RT.Call("CAR", args)))
		RT.Call("RAISE_NAMED", bc, RT.Call("NAMED", RT.Call("CAR", args)))
		return bc
	} else {
		return RT.Symbol("R_NilValue")
	}
}
