package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_load(call, op, args, env Value) Value {
	var (
		fname Value
		aenv  Value
		val   Value
		fp    Value
		cntxt Value
	)
	RT.Call("checkArity", op, args)
	if RT.Truth(RT.Unary("!", RT.Call("isValidString", Assign(LocalRef(&fname), RT.Call("CAR", args))))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"first argument must be a file name\"")))
	}
	Assign(LocalRef(&aenv), RT.Call("CADR", args))
	if RT.Truth(RT.Binary("==", RT.Call("TYPEOF", aenv), RT.Symbol("NILSXP"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"use of NULL environment is defunct\"")))
	} else {
		if RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", aenv), RT.Symbol("ENVSXP"))) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"envir\""))
		}
	}
	Assign(LocalRef(&fp), RT.Call("RC_fopen", RT.Call("STRING_ELT", fname, RT.Const("int", "0")), RT.Const("string", "\"rb\""), RT.Symbol("TRUE")))
	if RT.Truth(RT.Unary("!", fp)) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"unable to open file\"")))
	}
	RT.Call("begincontext", LocalRef(&cntxt), RT.Symbol("CTXT_CCODE"), RT.Symbol("R_NilValue"), RT.Symbol("R_BaseEnv"), RT.Symbol("R_BaseEnv"), RT.Symbol("R_NilValue"), RT.Symbol("R_NilValue"))
	RT.AssignField(cntxt, "cend", RT.SymbolRef("saveload_cleanup"))
	RT.AssignField(cntxt, "cenddata", fp)
	RT.Call("PROTECT", Assign(LocalRef(&val), RT.Call("R_LoadSavedData", fp, aenv)))
	RT.Call("endcontext", LocalRef(&cntxt))
	RT.Call("fclose", fp)
	RT.Call("UNPROTECT", RT.Const("int", "1"))
	return val
}
