package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_mkActiveBnd(call, op, args, rho Value) Value {
	var (
		sym Value
		fun Value
		env Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&sym), RT.Call("CAR", args))
	Assign(LocalRef(&fun), RT.Call("CADR", args))
	Assign(LocalRef(&env), RT.Call("CADDR", args))
	RT.Call("R_MakeActiveBinding", sym, fun, env)
	return RT.Symbol("R_NilValue")
}
