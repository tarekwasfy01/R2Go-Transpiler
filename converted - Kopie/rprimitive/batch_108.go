package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_substitute(call, op, args, rho Value) Value {
	var (
		argList               Value
		env                   Value
		s                     Value
		t                     Value
		do_substitute_formals Value
	)
	Assign(LocalRef(&do_substitute_formals), RT.Symbol("NULL"))
	if RT.Truth(RT.Binary("==", do_substitute_formals, RT.Symbol("NULL"))) {
		Assign(LocalRef(&do_substitute_formals), RT.Call("allocFormalsList2", RT.Call("install", RT.Const("string", "\"expr\"")), RT.Call("install", RT.Const("string", "\"env\""))))
	}
	RT.Call("PROTECT", Assign(LocalRef(&argList), RT.Call("matchArgs_NR", do_substitute_formals, args, call)))
	if RT.Truth(RT.Binary("==", RT.Call("CADR", argList), RT.Symbol("R_MissingArg"))) {
		Assign(LocalRef(&env), rho)
	} else {
		Assign(LocalRef(&env), RT.Call("eval", RT.Call("CADR", argList), rho))
	}
	if RT.Truth(RT.Binary("==", env, RT.Symbol("R_GlobalEnv"))) {
		Assign(LocalRef(&env), RT.Symbol("R_NilValue"))
	} else {
		if RT.Truth(RT.Binary("==", RT.Call("TYPEOF", env), RT.Symbol("VECSXP"))) {
			Assign(LocalRef(&env), RT.Call("NewEnvironment", RT.Symbol("R_NilValue"), RT.Call("VectorToPairList", env), RT.Symbol("R_BaseEnv")))
		} else {
			if RT.Truth(RT.Binary("==", RT.Call("TYPEOF", env), RT.Symbol("LISTSXP"))) {
				Assign(LocalRef(&env), RT.Call("NewEnvironment", RT.Symbol("R_NilValue"), RT.Call("duplicate", env), RT.Symbol("R_BaseEnv")))
			}
		}
	}
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Binary("!=", env, RT.Symbol("R_NilValue"))) {
			return false
		}
		return RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", env), RT.Symbol("ENVSXP")))
	}()) {
		RT.Call("errorcall", call, RT.Call("_", RT.Const("string", "\"invalid environment specified\"")))
	}
	RT.Call("PROTECT", env)
	RT.Call("PROTECT", Assign(LocalRef(&t), RT.Call("CONS", RT.Call("duplicate", RT.Call("CAR", argList)), RT.Symbol("R_NilValue"))))
	Assign(LocalRef(&s), RT.Call("substituteList", t, env))
	RT.Call("UNPROTECT", RT.Const("int", "3"))
	return RT.Call("CAR", s)
}
