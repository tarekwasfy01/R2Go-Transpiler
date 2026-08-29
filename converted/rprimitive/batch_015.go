package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_bindingType(call, op, args, rho Value) Value {
	var (
		sym Value
		env Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&sym), RT.Call("CAR", args))
	Assign(LocalRef(&env), RT.Call("CADR", args))
	switch RT.Key(RT.Call("R_GetBindingType", sym, env)) {
	case RT.Key(RT.Symbol("R_BindingTypeUnbound")):
		return RT.Call("mkString", RT.Const("string", "\"unbound\""))
	case RT.Key(RT.Symbol("R_BindingTypeValue")):
		return RT.Call("mkString", RT.Const("string", "\"value\""))
	case RT.Key(RT.Symbol("R_BindingTypeMissing")):
		return RT.Call("mkString", RT.Const("string", "\"missing\""))
	case RT.Key(RT.Symbol("R_BindingTypeDelayed")):
		return RT.Call("mkString", RT.Const("string", "\"delayed\""))
	case RT.Key(RT.Symbol("R_BindingTypeForced")):
		return RT.Call("mkString", RT.Const("string", "\"forced\""))
	case RT.Key(RT.Symbol("R_BindingTypeActive")):
		return RT.Call("mkString", RT.Const("string", "\"active\""))
	default:
		RT.Call("error", RT.Const("string", "\"unknown binding type; should not happen\""))
	}
	return RT.Symbol("R_NilValue")
}
