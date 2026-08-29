package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_lockBnd(call, op, args, rho Value) Value {
	var (
		sym Value
		env Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&sym), RT.Call("CAR", args))
	Assign(LocalRef(&env), RT.Call("CADR", args))
	switch RT.Key(RT.Call("PRIMVAL", op)) {
	case RT.Key(RT.Const("int", "0")):
		RT.Call("R_LockBinding", sym, env)
		break
	case RT.Key(RT.Const("int", "1")):
		RT.Call("R_unLockBinding", sym, env)
		break
	default:
		RT.Call("error", RT.Call("_", RT.Const("string", "\"unknown op\"")))
	}
	return RT.Symbol("R_NilValue")
}
