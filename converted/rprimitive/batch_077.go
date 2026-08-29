package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_makelazy(call, op, args, rho Value) Value {
	var (
		names  Value
		values Value
		val    Value
		expr   Value
		eenv   Value
		aenv   Value
		expr0  Value
		i      Value
		name   Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&names), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	if RT.Truth(RT.Unary("!", RT.Call("isString", names))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid first argument\"")))
	}
	Assign(LocalRef(&values), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&expr), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&eenv), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	if RT.Truth(RT.Unary("!", RT.Call("isEnvironment", eenv))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"eval.env\""))
	}
	Assign(LocalRef(&aenv), RT.Call("CAR", args))
	if RT.Truth(RT.Unary("!", RT.Call("isEnvironment", aenv))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"assign.env\""))
	}
	for Assign(LocalRef(&i), RT.Const("int", "0")); RT.Truth(RT.Binary("<", i, RT.Call("XLENGTH", names))); RT.Inc(LocalRef(&i), 1, true) {
		Assign(LocalRef(&name), RT.Call("installTrChar", RT.Call("STRING_ELT", names, i)))
		RT.Call("PROTECT", Assign(LocalRef(&val), RT.Call("eval", RT.Call("VECTOR_ELT", values, i), eenv)))
		RT.Call("PROTECT", Assign(LocalRef(&expr0), RT.Call("duplicate", expr)))
		RT.Call("SETCAR", RT.Call("CDR", expr0), val)
		RT.Call("defineVar", name, RT.Call("mkPROMISE", expr0, eenv), aenv)
		RT.Call("UNPROTECT", RT.Const("int", "2"))
	}
	return RT.Symbol("R_NilValue")
}
