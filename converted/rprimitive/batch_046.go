package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_envprofile(call, op, args, rho Value) Value {
	var (
		env Value
		ans Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&ans), RT.Symbol("R_NilValue"))
	Assign(LocalRef(&env), RT.Call("CAR", args))
	if RT.Truth(RT.Call("isEnvironment", env)) {
		if RT.Truth(RT.Call("IS_HASHED", env)) {
			Assign(LocalRef(&ans), RT.Call("R_HashProfile", RT.Call("HASHTAB", env)))
		}
	} else {
		RT.Call("error", RT.Const("string", "\"argument must be a hashed environment\""))
	}
	return ans
}
