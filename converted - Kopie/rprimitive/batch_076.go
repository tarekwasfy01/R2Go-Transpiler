package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_ls(call, op, args, rho Value) Value {
	var (
		tb       Value
		env      Value
		all      Value
		sort_nms Value
	)
	RT.Call("checkArity", op, args)
	if RT.Truth(RT.Call("IS_USER_DATABASE", RT.Call("CAR", args))) {
		Assign(LocalRef(&tb), RT.Cast("R_ObjectTable *", RT.Call("R_ExternalPtrAddr", RT.Call("HASHTAB", RT.Call("CAR", args)))))
		return RT.CallIndirect(RT.Field(tb, "objects"), tb)
	}
	Assign(LocalRef(&env), RT.Call("CAR", args))
	Assign(LocalRef(&all), RT.Call("asLogical", RT.Call("CADR", args)))
	if RT.Truth(RT.Binary("==", all, RT.Symbol("NA_LOGICAL"))) {
		Assign(LocalRef(&all), RT.Const("int", "0"))
	}
	Assign(LocalRef(&sort_nms), RT.Call("asLogical", RT.Call("CADDR", args)))
	if RT.Truth(RT.Binary("==", sort_nms, RT.Symbol("NA_LOGICAL"))) {
		Assign(LocalRef(&sort_nms), RT.Const("int", "0"))
	}
	return RT.Call("R_lsInternal3", env, RT.Cast("Rboolean", all), RT.Cast("Rboolean", sort_nms))
}
