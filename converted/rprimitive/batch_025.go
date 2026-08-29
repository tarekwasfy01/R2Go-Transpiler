package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_delayed(call, op, args, rho Value) Value {
	var (
		name Value
		expr Value
		eenv Value
		aenv Value
	)
	Assign(LocalRef(&name), RT.Symbol("R_NilValue"))
	RT.Call("checkArity", op, args)
	if RT.Truth(func() Value {
		if RT.Truth(RT.Unary("!", RT.Call("isString", RT.Call("CAR", args)))) {
			return true
		}
		return RT.Truth(RT.Binary("==", RT.Call("LENGTH", RT.Call("CAR", args)), RT.Const("int", "0")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid first argument\"")))
	} else {
		Assign(LocalRef(&name), RT.Call("installTrChar", RT.Call("STRING_ELT", RT.Call("CAR", args), RT.Const("int", "0"))))
	}
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&expr), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&eenv), RT.Call("CAR", args))
	if RT.Truth(RT.Call("isNull", eenv)) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"use of NULL environment is defunct\"")))
		Assign(LocalRef(&eenv), RT.Symbol("R_BaseEnv"))
	} else {
		if RT.Truth(RT.Unary("!", RT.Call("isEnvironment", eenv))) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"eval.env\""))
		}
	}
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&aenv), RT.Call("CAR", args))
	if RT.Truth(RT.Call("isNull", aenv)) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"use of NULL environment is defunct\"")))
		Assign(LocalRef(&aenv), RT.Symbol("R_BaseEnv"))
	} else {
		if RT.Truth(RT.Unary("!", RT.Call("isEnvironment", aenv))) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"assign.env\""))
		}
	}
	RT.Call("defineVar", name, RT.Call("mkPROMISE", expr, eenv), aenv)
	return RT.Symbol("R_NilValue")
}
