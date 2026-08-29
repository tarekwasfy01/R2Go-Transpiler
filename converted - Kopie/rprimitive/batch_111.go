package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_tailcall(call, op, args, rho Value) Value {
	var (
		expr    Value
		env     Value
		formals Value
		api     Value
		mask    Value
		target  Value
		fun     Value
		val     Value
	)
	if RT.Truth(RT.Binary("==", RT.Call("PRIMVAL", op), RT.Const("int", "0"))) {
		Assign(LocalRef(&formals), RT.Symbol("NULL"))
		if RT.Truth(RT.Binary("==", formals, RT.Symbol("NULL"))) {
			Assign(LocalRef(&formals), RT.Call("allocFormalsList2", RT.Call("install", RT.Const("string", "\"expr\"")), RT.Call("install", RT.Const("string", "\"envir\""))))
		}
		RT.Call("PROTECT_WITH_INDEX", Assign(LocalRef(&args), RT.Call("matchArgs_NR", formals, args, call)), LocalRef(&api))
		RT.Call("REPROTECT", Assign(LocalRef(&args), RT.Call("evalListKeepMissing", args, rho)), api)
		Assign(LocalRef(&expr), RT.Call("CAR", args))
		if RT.Truth(RT.Binary("==", expr, RT.Symbol("R_MissingArg"))) {
			RT.Call("R_MissingArgError", RT.Call("install", RT.Const("string", "\"expr\"")), RT.Call("getLexicalCall", rho), RT.Const("string", "\"tailcallError\""))
		}
		if RT.Truth(func() Value {
			if !RT.Truth(RT.Binary("==", RT.Call("TYPEOF", expr), RT.Symbol("EXPRSXP"))) {
				return false
			}
			return RT.Truth(RT.Binary("==", RT.Call("XLENGTH", expr), RT.Const("int", "1")))
		}()) {
			Assign(LocalRef(&expr), RT.Call("VECTOR_ELT", expr, RT.Const("int", "0")))
		}
		if RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", expr), RT.Symbol("LANGSXP"))) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"\\\"expr\\\" must be a call expression\"")))
		}
		Assign(LocalRef(&env), RT.Call("CADR", args))
		if RT.Truth(RT.Binary("==", env, RT.Symbol("R_MissingArg"))) {
			Assign(LocalRef(&env), rho)
		}
		RT.Call("UNPROTECT", RT.Const("int", "1"))
	} else {
		if RT.Truth(func() Value {
			if RT.Truth(RT.Binary("==", args, RT.Symbol("R_NilValue"))) {
				return true
			}
			return RT.Truth(RT.Binary("==", RT.Call("CAR", args), RT.Symbol("R_MissingArg")))
		}()) {
			RT.Call("R_MissingArgError", RT.Call("install", RT.Const("string", "\"FUN\"")), RT.Call("getLexicalCall", rho), RT.Const("string", "\"tailcallRecError\""))
		}
		Assign(LocalRef(&expr), RT.Call("LCONS", RT.Call("CAR", args), RT.Call("CDR", args)))
		Assign(LocalRef(&env), rho)
	}
	RT.Call("PROTECT", expr)
	RT.Call("PROTECT", env)
	Assign(LocalRef(&mask), RT.Binary("|", RT.Symbol("CTXT_BROWSER"), RT.Symbol("CTXT_FUNCTION")))
	Assign(LocalRef(&target), RT.Call("getTailcallTarget", rho, mask))
	if RT.Truth(RT.Binary("!=", target, RT.Symbol("NULL"))) {
		Assign(LocalRef(&fun), RT.Call("CAR", expr))
		if RT.Truth(func() Value {
			if !RT.Truth(RT.Binary("==", RT.Call("TYPEOF", fun), RT.Symbol("STRSXP"))) {
				return false
			}
			return RT.Truth(RT.Binary("==", RT.Call("XLENGTH", fun), RT.Const("int", "1")))
		}()) {
			Assign(LocalRef(&fun), RT.Call("installTrChar", RT.Call("STRING_ELT", fun, RT.Const("int", "0"))))
		}
		if RT.Truth(RT.Binary("==", RT.Call("TYPEOF", fun), RT.Symbol("SYMSXP"))) {
			Assign(LocalRef(&fun), RT.Call("findFun3", fun, env, call))
		} else {
			Assign(LocalRef(&fun), RT.Call("eval", fun, env))
		}
		RT.Call("PROTECT", fun)
		Assign(LocalRef(&val), RT.Call("allocVector", RT.Symbol("VECSXP"), RT.Const("int", "4")))
		RT.Call("UNPROTECT", RT.Const("int", "1"))
		RT.Call("SET_VECTOR_ELT", val, RT.Const("int", "0"), RT.Symbol("R_exec_token"))
		RT.Call("SET_VECTOR_ELT", val, RT.Const("int", "1"), expr)
		RT.Call("SET_VECTOR_ELT", val, RT.Const("int", "2"), env)
		RT.Call("SET_VECTOR_ELT", val, RT.Const("int", "3"), fun)
		RT.Call("R_jumpctxt", target, mask, val)
	} else {
		Assign(LocalRef(&val), RT.Call("eval", expr, rho))
		RT.Call("UNPROTECT", RT.Const("int", "2"))
		return val
	}
	return RT.Symbol("R_NilValue")
}
