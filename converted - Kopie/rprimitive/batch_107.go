package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_standardGeneric(call, op, args, env Value) Value {
	var (
		arg   Value
		value Value
		fdef  Value
		ptr   Value
	)
	Assign(LocalRef(&ptr), RT.Call("R_get_standardGeneric_ptr"))
	RT.Call("checkArity", op, args)
	RT.Call("check1arg", args, call, RT.Const("string", "\"f\""))
	if RT.Truth(RT.Unary("!", ptr)) {
		RT.Call("warningcall", call, RT.Call("_", RT.Const("string", "\"'standardGeneric' called without 'methods' dispatch enabled (will be ignored)\"")))
		RT.Call("R_set_standardGeneric_ptr", RT.Symbol("dispatchNonGeneric"), RT.Symbol("NULL"))
		Assign(LocalRef(&ptr), RT.Call("R_get_standardGeneric_ptr"))
	}
	Assign(LocalRef(&arg), RT.Call("CAR", args))
	if RT.Truth(RT.Unary("!", RT.Call("isValidStringF", arg))) {
		RT.Call("errorcall", call, RT.Call("_", RT.Const("string", "\"argument to 'standardGeneric' must be a non-empty character string\"")))
	}
	RT.Call("PROTECT", Assign(LocalRef(&fdef), RT.Call("get_this_generic", args)))
	if RT.Truth(RT.Call("isNull", fdef)) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"call to standardGeneric(\\\"%s\\\") apparently not from the body of that generic function\"")), RT.Call("translateChar", RT.Call("STRING_ELT", arg, RT.Const("int", "0"))))
	}
	Assign(LocalRef(&value), RT.CallIndirect(RT.Deref(ptr), arg, env, fdef))
	RT.Call("UNPROTECT", RT.Const("int", "1"))
	return value
}
