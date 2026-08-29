package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_sysbrowser(call, op, args, rho Value) Value {
	var (
		rval     Value
		cptr     Value
		prevcptr Value
		n        Value
	)
	Assign(LocalRef(&rval), RT.Symbol("R_NilValue"))
	Assign(LocalRef(&prevcptr), RT.Symbol("NULL"))
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&n), RT.Call("asInteger", RT.Call("CAR", args)))
	if RT.Truth(RT.Binary("<", n, RT.Const("int", "1"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"number of contexts must be positive\"")))
	}
	Assign(LocalRef(&cptr), RT.Symbol("R_GlobalContext"))
	for RT.Truth(RT.Binary("!=", cptr, RT.Symbol("R_ToplevelContext"))) {
		if RT.Truth(RT.Binary("==", RT.Field(cptr, "callflag"), RT.Symbol("CTXT_BROWSER"))) {
			break
		}
		Assign(LocalRef(&cptr), RT.Field(cptr, "nextcontext"))
	}
	if RT.Truth(RT.Unary("!", RT.Binary("==", RT.Field(cptr, "callflag"), RT.Symbol("CTXT_BROWSER")))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"no browser context to query\"")))
	}
	switch RT.Key(RT.Call("PRIMVAL", op)) {
	case RT.Key(RT.Const("int", "1")), RT.Key(RT.Const("int", "2")):
		if RT.Truth(RT.Binary(">", n, RT.Const("int", "1"))) {
			for RT.Truth(func() Value {
				if !RT.Truth(RT.Binary("!=", cptr, RT.Symbol("R_ToplevelContext"))) {
					return false
				}
				return RT.Truth(RT.Binary(">", n, RT.Const("int", "0")))
			}()) {
				if RT.Truth(RT.Binary("==", RT.Field(cptr, "callflag"), RT.Symbol("CTXT_BROWSER"))) {
					RT.Inc(LocalRef(&n), -1, true)
					break
				}
				Assign(LocalRef(&cptr), RT.Field(cptr, "nextcontext"))
			}
		}
		if RT.Truth(RT.Unary("!", RT.Binary("==", RT.Field(cptr, "callflag"), RT.Symbol("CTXT_BROWSER")))) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"not that many calls to browser are active\"")))
		}
		if RT.Truth(RT.Binary("==", RT.Call("PRIMVAL", op), RT.Const("int", "1"))) {
			Assign(LocalRef(&rval), RT.Call("CAR", RT.Field(cptr, "promargs")))
		} else {
			Assign(LocalRef(&rval), RT.Call("CADR", RT.Field(cptr, "promargs")))
		}
		break
	case RT.Key(RT.Const("int", "3")):
		for RT.Truth(func() Value {
			if !RT.Truth(RT.Binary("!=", cptr, RT.Symbol("R_ToplevelContext"))) {
				return false
			}
			return RT.Truth(RT.Binary(">", n, RT.Const("int", "0")))
		}()) {
			if RT.Truth(RT.Binary("&", RT.Field(cptr, "callflag"), RT.Symbol("CTXT_FUNCTION"))) {
				RT.Inc(LocalRef(&n), -1, true)
			}
			Assign(LocalRef(&prevcptr), cptr)
			Assign(LocalRef(&cptr), RT.Field(cptr, "nextcontext"))
		}
		if RT.Truth(RT.Unary("!", RT.Binary("&", RT.Field(cptr, "callflag"), RT.Symbol("CTXT_FUNCTION")))) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"not that many functions on the call stack\"")))
		}
		if RT.Truth(func() Value {
			if !RT.Truth(prevcptr) {
				return false
			}
			return RT.Truth(RT.Binary("==", RT.Field(prevcptr, "srcref"), RT.Symbol("R_InBCInterpreter")))
		}()) {
			if RT.Truth(func() Value {
				if !RT.Truth(RT.Binary("==", RT.Call("TYPEOF", RT.Field(cptr, "callfun")), RT.Symbol("CLOSXP"))) {
					return false
				}
				return RT.Truth(RT.Binary("==", RT.Call("TYPEOF", RT.Call("BODY", RT.Field(cptr, "callfun"))), RT.Symbol("BCODESXP")))
			}()) {
				RT.Call("warning", RT.Call("_", RT.Const("string", "\"debug flag in compiled function has no effect\"")))
			} else {
				RT.Call("warning", RT.Call("_", RT.Const("string", "\"debug will apply when function leaves compiled code\"")))
			}
		}
		RT.Call("SET_RDEBUG", RT.Field(cptr, "cloenv"), RT.Const("int", "1"))
		break
	}
	return rval
}
