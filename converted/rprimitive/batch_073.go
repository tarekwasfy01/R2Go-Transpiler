package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_loadfile(call, op, args, env Value) Value {
	var (
		file Value
		s    Value
		fp   Value
	)
	RT.Call("checkArity", op, args)
	RT.Call("PROTECT", Assign(LocalRef(&file), RT.Call("coerceVector", RT.Call("CAR", args), RT.Symbol("STRSXP"))))
	if RT.Truth(RT.Unary("!", RT.Call("isValidStringF", file))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"bad file name\"")))
	}
	Assign(LocalRef(&fp), RT.Call("RC_fopen", RT.Call("STRING_ELT", file, RT.Const("int", "0")), RT.Const("string", "\"rb\""), RT.Symbol("TRUE")))
	if RT.Truth(RT.Unary("!", fp)) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"unable to open 'file'\"")))
	}
	Assign(LocalRef(&s), RT.Call("R_LoadFromFile", fp, RT.Const("int", "0")))
	RT.Call("fclose", fp)
	RT.Call("UNPROTECT", RT.Const("int", "1"))
	return s
}
