package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_getNSValue(call, op, args, rho Value) Value {
	var (
		ns       Value
		name     Value
		exported Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&ns), RT.Call("CAR", args))
	Assign(LocalRef(&name), RT.Call("CADR", args))
	Assign(LocalRef(&exported), RT.Call("asLogical", RT.Call("CADDR", args)))
	return RT.Call("R_getNSValue", RT.Symbol("R_NilValue"), ns, name, exported)
}
