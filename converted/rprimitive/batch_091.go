package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_recall(call, op, args, rho Value) Value {
	var (
		cptr Value
		s    Value
		ans  Value
	)
	Assign(LocalRef(&cptr), RT.Symbol("R_GlobalContext"))
	for RT.Truth(RT.Binary("!=", cptr, RT.Symbol("NULL"))) {
		if RT.Truth(func() Value {
			if !RT.Truth(RT.Binary("==", RT.Field(cptr, "callflag"), RT.Symbol("CTXT_RETURN"))) {
				return false
			}
			return RT.Truth(RT.Binary("==", RT.Field(cptr, "cloenv"), rho))
		}()) {
			break
		}
		Assign(LocalRef(&cptr), RT.Field(cptr, "nextcontext"))
	}
	if RT.Truth(RT.Binary("!=", cptr, RT.Symbol("NULL"))) {
		Assign(LocalRef(&args), RT.Field(cptr, "promargs"))
	}
	Assign(LocalRef(&s), RT.Field(RT.Symbol("R_GlobalContext"), "sysparent"))
	for RT.Truth(RT.Binary("!=", cptr, RT.Symbol("NULL"))) {
		if RT.Truth(func() Value {
			if !RT.Truth(RT.Binary("==", RT.Field(cptr, "callflag"), RT.Symbol("CTXT_RETURN"))) {
				return false
			}
			return RT.Truth(RT.Binary("==", RT.Field(cptr, "cloenv"), s))
		}()) {
			break
		}
		Assign(LocalRef(&cptr), RT.Field(cptr, "nextcontext"))
	}
	if RT.Truth(RT.Binary("==", cptr, RT.Symbol("NULL"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"'Recall' called from outside a closure\"")))
	}
	if RT.Truth(RT.Binary("!=", RT.Field(cptr, "callfun"), RT.Symbol("R_NilValue"))) {
		RT.Call("PROTECT", Assign(LocalRef(&s), RT.Field(cptr, "callfun")))
	} else {
		if RT.Truth(RT.Binary("==", RT.Call("TYPEOF", RT.Call("CAR", RT.Field(cptr, "call"))), RT.Symbol("SYMSXP"))) {
			RT.Call("PROTECT", Assign(LocalRef(&s), RT.Call("findFun", RT.Call("CAR", RT.Field(cptr, "call")), RT.Field(cptr, "sysparent"))))
		} else {
			RT.Call("PROTECT", Assign(LocalRef(&s), RT.Call("eval", RT.Call("CAR", RT.Field(cptr, "call")), RT.Field(cptr, "sysparent"))))
		}
	}
	if RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", s), RT.Symbol("CLOSXP"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"'Recall' called from outside a closure\"")))
	}
	Assign(LocalRef(&ans), RT.Call("applyClosure", RT.Field(cptr, "call"), s, args, RT.Field(cptr, "sysparent"), RT.Symbol("R_NilValue"), RT.Symbol("TRUE")))
	RT.Call("UNPROTECT", RT.Const("int", "1"))
	return ans
}
