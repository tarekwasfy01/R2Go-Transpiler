package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_savefile(call, op, args, env Value) Value {
	var (
		fp      Value
		version Value
	)
	RT.Call("checkArity", op, args)
	if RT.Truth(RT.Unary("!", RT.Call("isValidStringF", RT.Call("CADR", args)))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"'file' must be non-empty string\"")))
	}
	if RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", RT.Call("CADDR", args)), RT.Symbol("LGLSXP"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"'ascii' must be logical\"")))
	}
	if RT.Truth(RT.Binary("==", RT.Call("CADDDR", args), RT.Symbol("R_NilValue"))) {
		Assign(LocalRef(&version), RT.Call("defaultSaveVersion"))
	} else {
		Assign(LocalRef(&version), RT.Call("asInteger", RT.Call("CADDDR", args)))
	}
	if RT.Truth(func() Value {
		if RT.Truth(RT.Binary("==", version, RT.Symbol("NA_INTEGER"))) {
			return true
		}
		return RT.Truth(RT.Binary("<=", version, RT.Const("int", "0")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"version\""))
	}
	Assign(LocalRef(&fp), RT.Call("RC_fopen", RT.Call("STRING_ELT", RT.Call("CADR", args), RT.Const("int", "0")), RT.Const("string", "\"wb\""), RT.Symbol("TRUE")))
	if RT.Truth(RT.Unary("!", fp)) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"unable to open 'file'\"")))
	}
	RT.Call("R_SaveToFileV", RT.Call("CAR", args), fp, RT.Index(RT.Call("INTEGER", RT.Call("CADDR", args)), RT.Const("int", "0")), version)
	RT.Call("fclose", fp)
	return RT.Symbol("R_NilValue")
}
