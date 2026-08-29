package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_lazyLoadDBinsertValue(call, op, args, env Value) Value {
	var (
		value   Value
		file    Value
		ascii   Value
		compsxp Value
		hook    Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&value), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&file), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&ascii), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&compsxp), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&hook), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	return RT.Call("R_lazyLoadDBinsertValue", value, file, ascii, compsxp, hook)
}
