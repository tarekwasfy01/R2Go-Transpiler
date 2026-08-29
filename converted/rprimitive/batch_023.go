package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_debug(call, op, args, rho Value) Value {
	var (
		ans Value
		s   Value
	)
	Assign(LocalRef(&ans), RT.Symbol("R_NilValue"))
	RT.Call("checkArity", op, args)
	if RT.Truth(RT.Call("isValidString", RT.Call("CAR", args))) {
		RT.Call("PROTECT", Assign(LocalRef(&s), RT.Call("installTrChar", RT.Call("STRING_ELT", RT.Call("CAR", args), RT.Const("int", "0")))))
		RT.Call("SETCAR", args, RT.Call("findFun", s, rho))
		RT.Call("UNPROTECT", RT.Const("int", "1"))
	}
	if RT.Truth(func() Value {
		if !RT.Truth(func() Value {
			if !RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", RT.Call("CAR", args)), RT.Symbol("CLOSXP"))) {
				return false
			}
			return RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", RT.Call("CAR", args)), RT.Symbol("SPECIALSXP")))
		}()) {
			return false
		}
		return RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", RT.Call("CAR", args)), RT.Symbol("BUILTINSXP")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"argument must be a function\"")))
	}
	switch RT.Key(RT.Call("PRIMVAL", op)) {
	case RT.Key(RT.Const("int", "0")):
		RT.Call("SET_RDEBUG", RT.Call("CAR", args), RT.Const("int", "1"))
		break
	case RT.Key(RT.Const("int", "1")):
		if RT.Truth(RT.Binary("!=", RT.Call("RDEBUG", RT.Call("CAR", args)), RT.Const("int", "1"))) {
			RT.Call("warning", RT.Const("string", "\"argument is not being debugged\""))
		}
		RT.Call("SET_RDEBUG", RT.Call("CAR", args), RT.Const("int", "0"))
		break
	case RT.Key(RT.Const("int", "2")):
		Assign(LocalRef(&ans), RT.Call("ScalarLogical", RT.Call("RDEBUG", RT.Call("CAR", args))))
		break
	case RT.Key(RT.Const("int", "3")):
		RT.Call("SET_RSTEP", RT.Call("CAR", args), RT.Const("int", "1"))
		break
	}
	return ans
}
