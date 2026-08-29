package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_serversocket(call, op, args, rho Value) Value {
	var (
		ans   Value
		class Value
		ncon  Value
		port  Value
		con   Value
	)
	Assign(LocalRef(&con), RT.Symbol("NULL"))
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&port), RT.Call("asInteger", RT.Call("CAR", args)))
	if RT.Truth(func() Value {
		if RT.Truth(RT.Binary("==", port, RT.Symbol("NA_INTEGER"))) {
			return true
		}
		return RT.Truth(RT.Binary("<", port, RT.Const("int", "0")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"port\""))
	}
	Assign(LocalRef(&ncon), RT.Call("NextConnection"))
	Assign(LocalRef(&con), RT.Call("R_newservsock", port))
	RT.AssignIndex(RT.Symbol("Connections"), ncon, con)
	RT.AssignField(con, "ex_ptr", RT.Call("PROTECT", RT.Call("R_MakeExternalPtr", RT.Field(con, "id"), RT.Call("install", RT.Const("string", "\"connection\"")), RT.Symbol("R_NilValue"))))
	RT.Call("PROTECT", Assign(LocalRef(&ans), RT.Call("ScalarInteger", ncon)))
	RT.Call("PROTECT", Assign(LocalRef(&class), RT.Call("allocVector", RT.Symbol("STRSXP"), RT.Const("int", "2"))))
	RT.Call("SET_STRING_ELT", class, RT.Const("int", "0"), RT.Call("mkChar", RT.Const("string", "\"servsockconn\"")))
	RT.Call("SET_STRING_ELT", class, RT.Const("int", "1"), RT.Call("mkChar", RT.Const("string", "\"connection\"")))
	RT.Call("classgets", ans, class)
	RT.Call("setAttrib", ans, RT.Symbol("R_ConnIdSymbol"), RT.Field(con, "ex_ptr"))
	RT.Call("R_RegisterCFinalizerEx", RT.Field(con, "ex_ptr"), RT.Symbol("conFinalizer"), RT.Symbol("FALSE"))
	RT.Call("UNPROTECT", RT.Const("int", "3"))
	return ans
}
