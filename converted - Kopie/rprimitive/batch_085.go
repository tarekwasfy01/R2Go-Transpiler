package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_onexit(call, op, args, rho Value) Value {
	var (
		ctxt              Value
		code              Value
		oldcode           Value
		argList           Value
		addit             Value
		after             Value
		do_onexit_formals Value
		codelist          Value
	)
	Assign(LocalRef(&addit), RT.Symbol("FALSE"))
	Assign(LocalRef(&after), RT.Symbol("TRUE"))
	Assign(LocalRef(&do_onexit_formals), RT.Symbol("NULL"))
	RT.Call("checkArity", op, args)
	if RT.Truth(RT.Binary("==", do_onexit_formals, RT.Symbol("NULL"))) {
		Assign(LocalRef(&do_onexit_formals), RT.Call("allocFormalsList3", RT.Call("install", RT.Const("string", "\"expr\"")), RT.Call("install", RT.Const("string", "\"add\"")), RT.Call("install", RT.Const("string", "\"after\""))))
	}
	RT.Call("PROTECT", Assign(LocalRef(&argList), RT.Call("matchArgs_NR", do_onexit_formals, args, call)))
	if RT.Truth(RT.Binary("==", RT.Call("CAR", argList), RT.Symbol("R_MissingArg"))) {
		Assign(LocalRef(&code), RT.Symbol("R_NilValue"))
	} else {
		Assign(LocalRef(&code), RT.Call("CAR", argList))
	}
	if RT.Truth(RT.Binary("!=", RT.Call("CADR", argList), RT.Symbol("R_MissingArg"))) {
		Assign(LocalRef(&addit), RT.Call("asLogical", RT.Call("PROTECT", RT.Call("eval", RT.Call("CADR", argList), rho))))
		RT.Call("UNPROTECT", RT.Const("int", "1"))
		if RT.Truth(RT.Binary("==", addit, RT.Symbol("NA_INTEGER"))) {
			RT.Call("errorcall", call, RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"add\""))
		}
	}
	if RT.Truth(RT.Binary("!=", RT.Call("CADDR", argList), RT.Symbol("R_MissingArg"))) {
		Assign(LocalRef(&after), RT.Call("asLogical", RT.Call("PROTECT", RT.Call("eval", RT.Call("CADDR", argList), rho))))
		RT.Call("UNPROTECT", RT.Const("int", "1"))
		if RT.Truth(RT.Binary("==", after, RT.Symbol("NA_INTEGER"))) {
			RT.Call("errorcall", call, RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"lifo\""))
		}
	}
	Assign(LocalRef(&ctxt), RT.Symbol("R_GlobalContext"))
	for RT.Truth(func() Value {
		if !RT.Truth(RT.Binary("!=", ctxt, RT.Symbol("R_ToplevelContext"))) {
			return false
		}
		return RT.Truth(RT.Unary("!", func() Value {
			if !RT.Truth(RT.Binary("&", RT.Field(ctxt, "callflag"), RT.Symbol("CTXT_FUNCTION"))) {
				return false
			}
			return RT.Truth(RT.Binary("==", RT.Field(ctxt, "cloenv"), rho))
		}()))
	}()) {
		Assign(LocalRef(&ctxt), RT.Field(ctxt, "nextcontext"))
	}
	if RT.Truth(RT.Binary("&", RT.Field(ctxt, "callflag"), RT.Symbol("CTXT_FUNCTION"))) {
		if RT.Truth(func() Value {
			if !RT.Truth(RT.Binary("==", code, RT.Symbol("R_NilValue"))) {
				return false
			}
			return RT.Truth(RT.Unary("!", addit))
		}()) {
			RT.AssignField(ctxt, "conexit", RT.Symbol("R_NilValue"))
		} else {
			Assign(LocalRef(&oldcode), RT.Field(ctxt, "conexit"))
			if RT.Truth(func() Value {
				if RT.Truth(RT.Binary("==", oldcode, RT.Symbol("R_NilValue"))) {
					return true
				}
				return RT.Truth(RT.Unary("!", addit))
			}()) {
				RT.AssignField(ctxt, "conexit", RT.Call("CONS", code, RT.Symbol("R_NilValue")))
			} else {
				if RT.Truth(after) {
					Assign(LocalRef(&codelist), RT.Call("PROTECT", RT.Call("CONS", code, RT.Symbol("R_NilValue"))))
					RT.AssignField(ctxt, "conexit", RT.Call("listAppend", RT.Call("shallow_duplicate", oldcode), codelist))
					RT.Call("UNPROTECT", RT.Const("int", "1"))
				} else {
					RT.AssignField(ctxt, "conexit", RT.Call("CONS", code, oldcode))
				}
			}
		}
	}
	RT.Call("UNPROTECT", RT.Const("int", "1"))
	return RT.Symbol("R_NilValue")
}
