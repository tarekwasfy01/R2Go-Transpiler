package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_matchcall(call, op, args, env Value) Value {
	var (
		formals Value
		actuals Value
		rlist   Value
		funcall Value
		f       Value
		b       Value
		rval    Value
		sysp    Value
		t1      Value
		t2      Value
		tail    Value
		expdots Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&funcall), RT.Call("CADR", args))
	if RT.Truth(RT.Binary("==", RT.Call("TYPEOF", funcall), RT.Symbol("EXPRSXP"))) {
		Assign(LocalRef(&funcall), RT.Call("VECTOR_ELT", funcall, RT.Const("int", "0")))
	}
	if RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", funcall), RT.Symbol("LANGSXP"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"call\""))
	}
	Assign(LocalRef(&b), RT.Call("CAR", args))
	if RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", b), RT.Symbol("CLOSXP"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"definition\""))
	}
	Assign(LocalRef(&sysp), RT.Call("CAR", RT.Call("CDDDR", args)))
	if RT.Truth(RT.Unary("!", RT.Call("isEnvironment", sysp))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"'envir' must be an environment\"")))
	}
	Assign(LocalRef(&expdots), RT.Call("asLogical", RT.Call("CAR", RT.Call("CDDR", args))))
	if RT.Truth(RT.Binary("==", expdots, RT.Symbol("NA_LOGICAL"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"expand.dots\""))
	}
	Assign(LocalRef(&formals), RT.Call("FORMALS", b))
	RT.Call("PROTECT", Assign(LocalRef(&actuals), RT.Call("shallow_duplicate", RT.Call("CDR", funcall))))
	Assign(LocalRef(&t2), RT.Symbol("R_MissingArg"))
	for Assign(LocalRef(&t1), actuals); RT.Truth(RT.Binary("!=", t1, RT.Symbol("R_NilValue"))); Assign(LocalRef(&t1), RT.Call("CDR", t1)) {
		if RT.Truth(RT.Binary("==", RT.Call("CAR", t1), RT.Symbol("R_DotsSymbol"))) {
			Assign(LocalRef(&t2), RT.Call("subDots", sysp))
			break
		}
	}
	if RT.Truth(RT.Binary("!=", t2, RT.Symbol("R_MissingArg"))) {
		if RT.Truth(RT.Binary("==", RT.Call("CAR", actuals), RT.Symbol("R_DotsSymbol"))) {
			RT.Call("UNPROTECT", RT.Const("int", "1"))
			Assign(LocalRef(&actuals), RT.Call("listAppend", t2, RT.Call("CDR", actuals)))
			RT.Call("PROTECT", actuals)
		} else {
			for Assign(LocalRef(&t1), actuals); RT.Truth(RT.Binary("!=", t1, RT.Symbol("R_NilValue"))); Assign(LocalRef(&t1), RT.Call("CDR", t1)) {
				if RT.Truth(RT.Binary("==", RT.Call("CADR", t1), RT.Symbol("R_DotsSymbol"))) {
					Assign(LocalRef(&tail), RT.Call("CDDR", t1))
					RT.Call("SETCDR", t1, t2)
					RT.Call("listAppend", actuals, tail)
					break
				}
			}
		}
	} else {
		if RT.Truth(RT.Binary("==", RT.Call("CAR", actuals), RT.Symbol("R_DotsSymbol"))) {
			RT.Call("UNPROTECT", RT.Const("int", "1"))
			Assign(LocalRef(&actuals), RT.Call("CDR", actuals))
			RT.Call("PROTECT", actuals)
		} else {
			for Assign(LocalRef(&t1), actuals); RT.Truth(RT.Binary("!=", t1, RT.Symbol("R_NilValue"))); Assign(LocalRef(&t1), RT.Call("CDR", t1)) {
				if RT.Truth(RT.Binary("==", RT.Call("CADR", t1), RT.Symbol("R_DotsSymbol"))) {
					Assign(LocalRef(&tail), RT.Call("CDDR", t1))
					RT.Call("SETCDR", t1, tail)
					break
				}
			}
		}
	}
	Assign(LocalRef(&rlist), RT.Call("matchArgs_RC", formals, actuals, call))
	for RT.Sequence(Assign(LocalRef(&f), formals), Assign(LocalRef(&b), rlist)); RT.Truth(RT.Binary("!=", b, RT.Symbol("R_NilValue"))); RT.Sequence(Assign(LocalRef(&b), RT.Call("CDR", b)), Assign(LocalRef(&f), RT.Call("CDR", f))) {
		RT.Call("SET_TAG", b, RT.Call("TAG", f))
	}
	RT.Call("PROTECT", Assign(LocalRef(&rlist), RT.Call("ExpandDots", rlist, expdots)))
	Assign(LocalRef(&rlist), RT.Call("StripUnmatched", rlist))
	RT.Call("PROTECT", Assign(LocalRef(&rval), RT.Call("allocSExp", RT.Symbol("LANGSXP"))))
	RT.Call("SETCAR", rval, RT.Call("lazy_duplicate", RT.Call("CAR", funcall)))
	RT.Call("SETCDR", rval, rlist)
	RT.Call("UNPROTECT", RT.Const("int", "3"))
	return rval
}
