package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_sys(call, op, args, rho Value) Value {
	var (
		i       Value
		n       Value
		nframe  Value
		rval    Value
		t       Value
		cptr    Value
		conexit Value
	)
	Assign(LocalRef(&n), RT.Unary("-", RT.Const("int", "1")))
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&t), RT.Field(RT.Symbol("R_GlobalContext"), "sysparent"))
	Assign(LocalRef(&cptr), RT.Call("getLexicalContext", t))
	if RT.Truth(RT.Binary("==", RT.Call("length", args), RT.Const("int", "1"))) {
		Assign(LocalRef(&n), RT.Call("asInteger", RT.Call("CAR", args)))
	}
	switch RT.Key(RT.Call("PRIMVAL", op)) {
	case RT.Key(RT.Const("int", "1")):
		if RT.Truth(RT.Binary("==", n, RT.Symbol("NA_INTEGER"))) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"n\""))
		}
		Assign(LocalRef(&i), Assign(LocalRef(&nframe), RT.Call("framedepth", cptr)))
		for RT.Truth(RT.Binary(">", RT.Inc(LocalRef(&n), -1, true), RT.Const("int", "0"))) {
			Assign(LocalRef(&i), RT.Call("R_sysparent", RT.Binary("+", RT.Binary("-", nframe, i), RT.Const("int", "1")), cptr))
		}
		return RT.Call("ScalarInteger", i)
	case RT.Key(RT.Const("int", "2")):
		if RT.Truth(RT.Binary("==", n, RT.Symbol("NA_INTEGER"))) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"which\""))
		}
		return RT.Call("R_syscall", n, cptr)
	case RT.Key(RT.Const("int", "3")):
		if RT.Truth(RT.Binary("==", n, RT.Symbol("NA_INTEGER"))) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"which\""))
		}
		return RT.Call("R_sysframe", n, cptr)
	case RT.Key(RT.Const("int", "4")):
		return RT.Call("ScalarInteger", RT.Call("framedepth", cptr))
	case RT.Key(RT.Const("int", "5")):
		Assign(LocalRef(&nframe), RT.Call("framedepth", cptr))
		RT.Call("PROTECT", Assign(LocalRef(&rval), RT.Call("allocList", nframe)))
		Assign(LocalRef(&t), rval)
		for Assign(LocalRef(&i), RT.Const("int", "1")); RT.Truth(RT.Binary("<=", i, nframe)); RT.Sequence(RT.Inc(LocalRef(&i), 1, true), Assign(LocalRef(&t), RT.Call("CDR", t))) {
			RT.Call("SETCAR", t, RT.Call("R_syscall", i, cptr))
		}
		RT.Call("UNPROTECT", RT.Const("int", "1"))
		return rval
	case RT.Key(RT.Const("int", "6")):
		Assign(LocalRef(&nframe), RT.Call("framedepth", cptr))
		RT.Call("PROTECT", Assign(LocalRef(&rval), RT.Call("allocList", nframe)))
		Assign(LocalRef(&t), rval)
		for Assign(LocalRef(&i), RT.Const("int", "1")); RT.Truth(RT.Binary("<=", i, nframe)); RT.Sequence(RT.Inc(LocalRef(&i), 1, true), Assign(LocalRef(&t), RT.Call("CDR", t))) {
			RT.Call("SETCAR", t, RT.Call("R_sysframe", i, cptr))
		}
		RT.Call("UNPROTECT", RT.Const("int", "1"))
		return rval
	case RT.Key(RT.Const("int", "7")):
		Assign(LocalRef(&conexit), RT.Field(cptr, "conexit"))
		if RT.Truth(RT.Binary("==", conexit, RT.Symbol("R_NilValue"))) {
			return RT.Symbol("R_NilValue")
		} else {
			if RT.Truth(RT.Binary("==", RT.Call("CDR", conexit), RT.Symbol("R_NilValue"))) {
				return RT.Call("CAR", conexit)
			} else {
				return RT.Call("LCONS", RT.Symbol("R_BraceSymbol"), conexit)
			}
		}
	case RT.Key(RT.Const("int", "8")):
		Assign(LocalRef(&nframe), RT.Call("framedepth", cptr))
		Assign(LocalRef(&rval), RT.Call("allocVector", RT.Symbol("INTSXP"), nframe))
		for Assign(LocalRef(&i), RT.Const("int", "0")); RT.Truth(RT.Binary("<", i, nframe)); RT.Inc(LocalRef(&i), 1, true) {
			RT.AssignIndex(RT.Call("INTEGER", rval), i, RT.Call("R_sysparent", RT.Binary("-", nframe, i), cptr))
		}
		return rval
	case RT.Key(RT.Const("int", "9")):
		if RT.Truth(RT.Binary("==", n, RT.Symbol("NA_INTEGER"))) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' value\"")), RT.Const("string", "\"which\""))
		}
		return RT.Call("R_sysfunction", n, cptr)
	default:
		RT.Call("error", RT.Call("_", RT.Const("string", "\"internal error in 'do_sys'\"")))
		return RT.Symbol("R_NilValue")
	}
}
