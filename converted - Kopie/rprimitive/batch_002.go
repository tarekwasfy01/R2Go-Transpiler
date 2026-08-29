package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_Externalgr(call, op, args, env Value) Value {
	var (
		retval Value
		dd     Value
		record Value
	)
	Assign(LocalRef(&dd), RT.Call("GEcurrentDevice"))
	Assign(LocalRef(&record), RT.Field(dd, "recordGraphics"))
	RT.AssignField(dd, "recordGraphics", RT.Symbol("false"))
	RT.Call("PROTECT", Assign(LocalRef(&retval), do_External(call, op, args, env)))
	RT.AssignField(dd, "recordGraphics", record)
	if RT.Truth(RT.Call("GErecording", call, dd)) {
		if RT.Truth(RT.Unary("!", RT.Call("GEcheckState", dd))) {
			RT.Call("errorcall", call, RT.Call("_", RT.Const("string", "\"invalid graphics state\"")))
		}
		RT.Call("R_args_enable_refcnt", args)
		RT.Call("GErecordGraphicOperation", op, args, dd)
	}
	RT.Call("check_retval", call, retval)
	RT.Call("UNPROTECT", RT.Const("int", "1"))
	return retval
}
