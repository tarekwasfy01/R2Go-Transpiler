package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_regNS(call, op, args, rho Value) Value {
	var (
		name Value
		val  Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&name), RT.Call("checkNSname", call, RT.Call("CAR", args)))
	Assign(LocalRef(&val), RT.Call("CADR", args))
	if RT.Truth(RT.Binary("!=", RT.Call("R_findVarInFrame", RT.Symbol("R_NamespaceRegistry"), name), RT.Symbol("R_UnboundValue"))) {
		RT.Call("errorcall", call, RT.Call("_", RT.Const("string", "\"namespace already registered\"")))
	}
	RT.Call("defineVar", name, val, RT.Symbol("R_NamespaceRegistry"))
	return RT.Symbol("R_NilValue")
}
