package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_dotType(call, op, args, rho Value) Value {
	var (
		i   Value
		env Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&i), RT.Call("asInteger", RT.Call("CAR", args)))
	Assign(LocalRef(&env), RT.Call("resolveDotsEnv", RT.Call("CADR", args), RT.Call("asLogical", RT.Call("CADDR", args))))
	switch RT.Key(RT.Call("R_GetDotType", i, env)) {
	case RT.Key(RT.Symbol("R_DotTypeValue")):
		return RT.Call("mkString", RT.Const("string", "\"value\""))
	case RT.Key(RT.Symbol("R_DotTypeMissing")):
		return RT.Call("mkString", RT.Const("string", "\"missing\""))
	case RT.Key(RT.Symbol("R_DotTypeDelayed")):
		return RT.Call("mkString", RT.Const("string", "\"delayed\""))
	case RT.Key(RT.Symbol("R_DotTypeForced")):
		return RT.Call("mkString", RT.Const("string", "\"forced\""))
	default:
		RT.Call("error", RT.Const("string", "\"unknown dot type; should not happen\""))
	}
	return RT.Symbol("R_NilValue")
}
