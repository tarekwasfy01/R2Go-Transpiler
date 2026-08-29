package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_tryCatchHelper(call, op, args, env Value) Value {
	var (
		eptr Value
		sw   Value
		cond Value
		ptcd Value
		val  Value
	)
	Assign(LocalRef(&eptr), RT.Call("CAR", args))
	Assign(LocalRef(&sw), RT.Call("CADR", args))
	Assign(LocalRef(&cond), RT.Call("CADDR", args))
	if RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", eptr), RT.Symbol("EXTPTRSXP"))) {
		RT.Call("error", RT.Const("string", "\"not an external pointer\""))
	}
	Assign(LocalRef(&ptcd), RT.Call("R_ExternalPtrAddr", RT.Call("CAR", args)))
	switch RT.Key(RT.Call("asInteger", sw)) {
	case RT.Key(RT.Const("int", "0")):
		if RT.Truth(RT.Field(ptcd, "suspended")) {
			return RT.CallIndirect(RT.Field(ptcd, "body"), RT.Field(ptcd, "bdata"))
		} else {
			RT.AssignSymbol("R_interrupts_suspended", RT.Symbol("FALSE"))
			Assign(LocalRef(&val), RT.CallIndirect(RT.Field(ptcd, "body"), RT.Field(ptcd, "bdata")))
			RT.AssignSymbol("R_interrupts_suspended", RT.Symbol("TRUE"))
			return val
		}
	case RT.Key(RT.Const("int", "1")):
		if RT.Truth(RT.Binary("!=", RT.Field(ptcd, "handler"), RT.Symbol("NULL"))) {
			return RT.CallIndirect(RT.Field(ptcd, "handler"), cond, RT.Field(ptcd, "hdata"))
		} else {
			return RT.Symbol("R_NilValue")
		}
	case RT.Key(RT.Const("int", "2")):
		if RT.Truth(RT.Binary("!=", RT.Field(ptcd, "finally"), RT.Symbol("NULL"))) {
			RT.CallIndirect(RT.Field(ptcd, "finally"), RT.Field(ptcd, "fdata"))
		}
		return RT.Symbol("R_NilValue")
	default:
		return RT.Symbol("R_NilValue")
	}
}
