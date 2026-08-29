package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_unregNS(call, op, args, rho Value) Value {
	var (
		name     Value
		hashcode Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&name), RT.Call("checkNSname", call, RT.Call("CAR", args)))
	if RT.Truth(RT.Binary("==", RT.Call("R_findVarInFrame", RT.Symbol("R_NamespaceRegistry"), name), RT.Symbol("R_UnboundValue"))) {
		RT.Call("errorcall", call, RT.Call("_", RT.Const("string", "\"namespace not registered\"")))
	}
	if RT.Truth(RT.Unary("!", RT.Call("HASHASH", RT.Call("PRINTNAME", name)))) {
		Assign(LocalRef(&hashcode), RT.Call("R_Newhashpjw", RT.Call("CHAR", RT.Call("PRINTNAME", name))))
	} else {
		Assign(LocalRef(&hashcode), RT.Call("HASHVALUE", RT.Call("PRINTNAME", name)))
	}
	RT.Call("RemoveVariable", name, hashcode, RT.Symbol("R_NamespaceRegistry"))
	return RT.Symbol("R_NilValue")
}
