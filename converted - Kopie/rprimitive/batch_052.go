package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_getRegNS(call, op, args, rho Value) Value {
	var (
		name Value
		val  Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&name), RT.Call("checkNSname", call, RT.Call("PROTECT", RT.Call("coerceVector", RT.Call("CAR", args), RT.Symbol("SYMSXP")))))
	RT.Call("UNPROTECT", RT.Const("int", "1"))
	Assign(LocalRef(&val), RT.Call("R_findVarInFrame", RT.Symbol("R_NamespaceRegistry"), name))
	switch RT.Key(RT.Call("PRIMVAL", op)) {
	case RT.Key(RT.Const("int", "0")):
		if RT.Truth(RT.Binary("==", val, RT.Symbol("R_UnboundValue"))) {
			return RT.Symbol("R_NilValue")
		} else {
			return val
		}
	case RT.Key(RT.Const("int", "1")):
		return RT.Call("ScalarLogical", func() Value {
			if RT.Truth(RT.Binary("==", val, RT.Symbol("R_UnboundValue"))) {
				return RT.Symbol("FALSE")
			}
			return RT.Symbol("TRUE")
		}())
	default:
		RT.Call("error", RT.Call("_", RT.Const("string", "\"unknown op\"")))
	}
	return RT.Symbol("R_NilValue")
}
