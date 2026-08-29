package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_importIntoEnv(call, op, args, rho Value) Value {
	var (
		impenv   Value
		impnames Value
		expenv   Value
		expnames Value
		impsym   Value
		expsym   Value
		val      Value
		i        Value
		n        Value
		binding  Value
		env      Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&impenv), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&impnames), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&expenv), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&expnames), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	if RT.Truth(RT.Binary("==", RT.Call("TYPEOF", impenv), RT.Symbol("NILSXP"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"use of NULL environment is defunct\"")))
	}
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", impenv), RT.Symbol("ENVSXP"))) {
			return false
		}
		return RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", Assign(LocalRef(&impenv), RT.Call("simple_as_environment", impenv))), RT.Symbol("ENVSXP")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"bad import environment argument\"")))
	}
	if RT.Truth(RT.Binary("==", RT.Call("TYPEOF", expenv), RT.Symbol("NILSXP"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"use of NULL environment is defunct\"")))
	}
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", expenv), RT.Symbol("ENVSXP"))) {
			return false
		}
		return RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", Assign(LocalRef(&expenv), RT.Call("simple_as_environment", expenv))), RT.Symbol("ENVSXP")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"bad export environment argument\"")))
	}
	if RT.Truth(func() Value {
		if RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", impnames), RT.Symbol("STRSXP"))) {
			return true
		}
		return RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", expnames), RT.Symbol("STRSXP")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"names\""))
	}
	if RT.Truth(RT.Binary("!=", RT.Call("LENGTH", impnames), RT.Call("LENGTH", expnames))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"length of import and export names must match\"")))
	}
	Assign(LocalRef(&n), RT.Call("LENGTH", impnames))
	for Assign(LocalRef(&i), RT.Const("int", "0")); RT.Truth(RT.Binary("<", i, n)); RT.Inc(LocalRef(&i), 1, true) {
		Assign(LocalRef(&impsym), RT.Call("installTrChar", RT.Call("STRING_ELT", impnames, i)))
		Assign(LocalRef(&expsym), RT.Call("installTrChar", RT.Call("STRING_ELT", expnames, i)))
		Assign(LocalRef(&binding), RT.Symbol("R_NilValue"))
		for RT.Sequence(Assign(LocalRef(&env), expenv)); RT.Truth(func() Value {
			if !RT.Truth(RT.Binary("!=", env, RT.Symbol("R_EmptyEnv"))) {
				return false
			}
			return RT.Truth(RT.Binary("==", binding, RT.Symbol("R_NilValue")))
		}()); Assign(LocalRef(&env), RT.Call("ENCLOS", env)) {
			if RT.Truth(RT.Binary("==", env, RT.Symbol("R_BaseNamespace"))) {
				if RT.Truth(RT.Binary("!=", RT.Call("SYMVALUE", expsym), RT.Symbol("R_UnboundValue"))) {
					Assign(LocalRef(&binding), expsym)
				}
			} else {
				Assign(LocalRef(&binding), RT.Call("findVarLocInFrame", env, expsym, RT.Symbol("NULL")))
			}
		}
		if RT.Truth(RT.Binary("==", binding, RT.Symbol("R_NilValue"))) {
			Assign(LocalRef(&binding), expsym)
		}
		if RT.Truth(RT.Binary("==", RT.Call("TYPEOF", binding), RT.Symbol("SYMSXP"))) {
			if RT.Truth(RT.Binary("==", RT.Call("SYMVALUE", expsym), RT.Symbol("R_UnboundValue"))) {
				RT.Call("error", RT.Call("_", RT.Const("string", "\"exported symbol '%s' has no value\"")), RT.Call("CHAR", RT.Call("PRINTNAME", expsym)))
			}
			Assign(LocalRef(&val), RT.Call("SYMVALUE", expsym))
		} else {
			Assign(LocalRef(&val), RT.Call("CAR", binding))
		}
		if RT.Truth(RT.Call("IS_ACTIVE_BINDING", binding)) {
			RT.Call("R_MakeActiveBinding", impsym, val, impenv)
		} else {
			if RT.Truth(func() Value {
				if RT.Truth(RT.Binary("==", impenv, RT.Symbol("R_BaseNamespace"))) {
					return true
				}
				return RT.Truth(RT.Binary("==", impenv, RT.Symbol("R_BaseEnv")))
			}()) {
				RT.Call("gsetVar", impsym, val, impenv)
			} else {
				RT.Call("defineVar", impsym, val, impenv)
			}
		}
	}
	return RT.Symbol("R_NilValue")
}
