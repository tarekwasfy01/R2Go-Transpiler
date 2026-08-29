package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_recordGraphics(call, op, args, env Value) Value {
	var (
		x         Value
		evalenv   Value
		retval    Value
		dd        Value
		record    Value
		code      Value
		list      Value
		parentenv Value
		xptr      Value
	)
	Assign(LocalRef(&dd), RT.Call("GEcurrentDevice"))
	Assign(LocalRef(&record), RT.Field(dd, "recordGraphics"))
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&code), RT.Call("CAR", args))
	Assign(LocalRef(&list), RT.Call("CADR", args))
	Assign(LocalRef(&parentenv), RT.Call("CADDR", args))
	if RT.Truth(RT.Unary("!", RT.Call("isLanguage", code))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"'expr' argument must be an expression\"")))
	}
	if RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", list), RT.Symbol("VECSXP"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"'list' argument must be a list\"")))
	}
	if RT.Truth(RT.Call("isNull", parentenv)) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"use of NULL environment is defunct\"")))
		Assign(LocalRef(&parentenv), RT.Symbol("R_BaseEnv"))
	} else {
		if RT.Truth(RT.Unary("!", RT.Call("isEnvironment", parentenv))) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"'env' argument must be an environment\"")))
		}
	}
	RT.Call("PROTECT", Assign(LocalRef(&x), RT.Call("VectorToPairList", list)))
	for RT.Sequence(Assign(LocalRef(&xptr), x)); RT.Truth(RT.Binary("!=", xptr, RT.Symbol("R_NilValue"))); Assign(LocalRef(&xptr), RT.Call("CDR", xptr)) {
		RT.Call("ENSURE_NAMEDMAX", RT.Call("CAR", xptr))
	}
	RT.Call("PROTECT", Assign(LocalRef(&evalenv), RT.Call("NewEnvironment", RT.Symbol("R_NilValue"), x, parentenv)))
	RT.AssignField(dd, "recordGraphics", RT.Symbol("FALSE"))
	RT.Call("PROTECT", Assign(LocalRef(&retval), RT.Call("Rf_eval_with_gd", code, evalenv, dd)))
	RT.AssignField(dd, "recordGraphics", record)
	if RT.Truth(RT.Call("GErecording", call, dd)) {
		if RT.Truth(RT.Unary("!", RT.Call("GEcheckState", dd))) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid graphics state\"")))
		}
		RT.Call("GErecordGraphicOperation", op, args, dd)
	}
	RT.Call("UNPROTECT", RT.Const("int", "3"))
	return retval
}
