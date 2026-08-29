package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_dotDelayedExpr(call, op, args, rho Value) Value {
	var (
		i   Value
		env Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&i), RT.Call("asInteger", RT.Call("CAR", args)))
	Assign(LocalRef(&env), RT.Call("resolveDotsEnv", RT.Call("CADR", args), RT.Call("asLogical", RT.Call("CADDR", args))))
	return RT.Call("R_DotDelayedExpression", i, env)
}
