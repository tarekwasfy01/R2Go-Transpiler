package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_asfunction(call, op, args, rho Value) Value {
	var (
		arglist Value
		envir   Value
		n       Value
		names   Value
		pargs   Value
		i       Value
		body    Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&arglist), RT.Call("CAR", args))
	if RT.Truth(RT.Unary("!", RT.Call("isNewList", arglist))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"list argument expected\"")))
	}
	Assign(LocalRef(&envir), RT.Call("CADR", args))
	if RT.Truth(RT.Call("isNull", envir)) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"use of NULL environment is defunct\"")))
		Assign(LocalRef(&envir), RT.Symbol("R_BaseEnv"))
	} else {
		if RT.Truth(RT.Unary("!", RT.Call("isEnvironment", envir))) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid environment\"")))
		}
	}
	Assign(LocalRef(&n), RT.Call("length", arglist))
	if RT.Truth(RT.Binary("<", n, RT.Const("int", "1"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"argument must have length at least 1\"")))
	}
	Assign(LocalRef(&names), RT.Call("PROTECT", RT.Call("getAttrib", arglist, RT.Symbol("R_NamesSymbol"))))
	Assign(LocalRef(&pargs), RT.Call("PROTECT", Assign(LocalRef(&args), RT.Call("allocList", RT.Binary("-", n, RT.Const("int", "1"))))))
	for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, RT.Binary("-", n, RT.Const("int", "1")))); RT.Inc(LocalRef(&i), 1, true) {
		RT.Call("SETCAR", pargs, RT.Call("VECTOR_ELT", arglist, i))
		if RT.Truth(func() Value {
			if !RT.Truth(RT.Binary("!=", names, RT.Symbol("R_NilValue"))) {
				return false
			}
			return RT.Truth(RT.Binary("!=", RT.Deref(RT.Call("CHAR", RT.Call("STRING_ELT", names, i))), RT.Const("char", "'\\0'")))
		}()) {
			RT.Call("SET_TAG", pargs, RT.Call("installTrChar", RT.Call("STRING_ELT", names, i)))
		} else {
			RT.Call("SET_TAG", pargs, RT.Symbol("R_NilValue"))
		}
		Assign(LocalRef(&pargs), RT.Call("CDR", pargs))
	}
	RT.Call("CheckFormals", args, RT.Const("string", "\"as.function\""))
	Assign(LocalRef(&body), RT.Call("PROTECT", RT.Call("VECTOR_ELT", arglist, RT.Binary("-", n, RT.Const("int", "1")))))
	if RT.Truth(func() Value {
		if RT.Truth(func() Value {
			if RT.Truth(func() Value {
				if RT.Truth(func() Value {
					if RT.Truth(func() Value {
						if RT.Truth(RT.Call("isList", body)) {
							return true
						}
						return RT.Truth(RT.Call("isLanguage", body))
					}()) {
						return true
					}
					return RT.Truth(RT.Call("isSymbol", body))
				}()) {
					return true
				}
				return RT.Truth(RT.Call("isExpression", body))
			}()) {
				return true
			}
			return RT.Truth(RT.Call("isVector", body))
		}()) {
			return true
		}
		return RT.Truth(RT.Call("isByteCode", body))
	}()) {
		Assign(LocalRef(&args), RT.Call("mkCLOSXP", args, body, envir))
	} else {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid body for function\"")))
	}
	RT.Call("UNPROTECT", RT.Const("int", "3"))
	return args
}
