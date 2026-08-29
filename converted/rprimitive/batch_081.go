package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_mkUnbound(call, op, args, rho Value) Value {
	var (
		sym Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&sym), RT.Call("CAR", args))
	if RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", sym), RT.Symbol("SYMSXP"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"not a symbol\"")))
	}
	if RT.Truth(RT.Call("FRAME_IS_LOCKED", RT.Symbol("R_BaseEnv"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"cannot remove bindings from a locked environment\"")))
	}
	if RT.Truth(RT.Call("R_BindingIsLocked", sym, RT.Symbol("R_BaseEnv"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"cannot unbind a locked binding\"")))
	}
	if RT.Truth(RT.Call("R_BindingIsActive", sym, RT.Symbol("R_BaseEnv"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"cannot unbind an active binding\"")))
	}
	RT.Call("SET_SYMVALUE", sym, RT.Symbol("R_UnboundValue"))
	return RT.Symbol("R_NilValue")
}
