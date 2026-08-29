package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_remove(call, op, args, rho Value) Value {
	var (
		name      Value
		envarg    Value
		ginherits Value
		i         Value
		done      Value
		tsym      Value
		hashcode  Value
		tenv      Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&name), RT.Call("CAR", args))
	if RT.Truth(RT.Binary("==", RT.Call("TYPEOF", name), RT.Symbol("NILSXP"))) {
		return RT.Symbol("R_NilValue")
	}
	if RT.Truth(RT.Unary("!", RT.Call("isString", name))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid first argument\"")))
	}
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&envarg), RT.Call("CAR", args))
	if RT.Truth(RT.Binary("==", RT.Call("TYPEOF", envarg), RT.Symbol("NILSXP"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"use of NULL environment is defunct\"")))
	}
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", envarg), RT.Symbol("ENVSXP"))) {
			return false
		}
		return RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", Assign(LocalRef(&envarg), RT.Call("simple_as_environment", envarg))), RT.Symbol("ENVSXP")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"envir\""))
	}
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&ginherits), RT.Call("asLogical", RT.Call("CAR", args)))
	if RT.Truth(RT.Binary("==", ginherits, RT.Symbol("NA_LOGICAL"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"inherits\""))
	}
	for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0")), Assign(LocalRef(&done), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, RT.Call("LENGTH", name))); RT.Inc(LocalRef(&i), 1, true) {
		Assign(LocalRef(&tsym), RT.Call("installTrChar", RT.Call("STRING_ELT", name, i)))
		if RT.Truth(RT.Unary("!", RT.Call("HASHASH", RT.Call("PRINTNAME", tsym)))) {
			Assign(LocalRef(&hashcode), RT.Call("R_Newhashpjw", RT.Call("CHAR", RT.Call("PRINTNAME", tsym))))
		} else {
			Assign(LocalRef(&hashcode), RT.Call("HASHVALUE", RT.Call("PRINTNAME", tsym)))
		}
		Assign(LocalRef(&tenv), envarg)
		for RT.Truth(RT.Binary("!=", tenv, RT.Symbol("R_EmptyEnv"))) {
			Assign(LocalRef(&done), RT.Call("RemoveVariable", tsym, hashcode, tenv))
			if RT.Truth(func() Value {
				if RT.Truth(done) {
					return true
				}
				return RT.Truth(RT.Unary("!", ginherits))
			}()) {
				break
			}
			Assign(LocalRef(&tenv), RT.Call("CDR", tenv))
		}
		if RT.Truth(RT.Unary("!", done)) {
			RT.Call("warning", RT.Call("_", RT.Const("string", "\"object '%s' not found\"")), RT.Call("EncodeChar", RT.Call("PRINTNAME", tsym)))
		}
	}
	return RT.Symbol("R_NilValue")
}
