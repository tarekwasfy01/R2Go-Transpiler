package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_bndIsLocked(call, op, args, rho Value) Value {
	var (
		sym Value
		env Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&sym), RT.Call("CAR", args))
	Assign(LocalRef(&env), RT.Call("CADR", args))
	return RT.Call("ScalarLogical", RT.Call("R_BindingIsLocked", sym, env))
}
