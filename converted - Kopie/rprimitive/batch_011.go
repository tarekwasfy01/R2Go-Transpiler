package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_bcclose(call, op, args, rho Value) Value {
	var (
		forms Value
		body  Value
		env   Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&forms), RT.Call("CAR", args))
	Assign(LocalRef(&body), RT.Call("CADR", args))
	Assign(LocalRef(&env), RT.Call("CADDR", args))
	RT.Call("CheckFormals", forms, RT.Const("string", "\"bcClose\""))
	if RT.Truth(RT.Unary("!", RT.Call("isByteCode", body))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid body\"")))
	}
	if RT.Truth(RT.Call("isNull", env)) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"use of NULL environment is defunct\"")))
		Assign(LocalRef(&env), RT.Symbol("R_BaseEnv"))
	} else {
		if RT.Truth(RT.Unary("!", RT.Call("isEnvironment", env))) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid environment\"")))
		}
	}
	return RT.Call("mkCLOSXP", forms, body, env)
}
