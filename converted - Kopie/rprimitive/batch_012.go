package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_bcprofcounts(call, op, args, env Value) Value {
	var (
		val Value
		i   Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&val), RT.Call("allocVector", RT.Symbol("INTSXP"), RT.Symbol("OPCOUNT")))
	for Assign(LocalRef(&i), RT.Const("int", "0")); RT.Truth(RT.Binary("<", i, RT.Symbol("OPCOUNT"))); RT.Inc(LocalRef(&i), 1, true) {
		RT.AssignIndex(RT.Call("INTEGER", val), i, RT.Index(RT.Symbol("opcode_counts"), i))
	}
	return val
}
