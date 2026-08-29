package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_envirgets(call, op, args, rho Value) Value {
	var (
		s   Value
		env Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&s), RT.Call("CAR", args))
	Assign(LocalRef(&env), RT.Call("CADR", args))
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Binary("==", RT.Call("TYPEOF", s), RT.Symbol("CLOSXP"))) {
			return false
		}
		return RT.Truth(func() Value {
			if RT.Truth(func() Value {
				if RT.Truth(RT.Call("isEnvironment", env)) {
					return true
				}
				return RT.Truth(RT.Call("isEnvironment", Assign(LocalRef(&env), RT.Call("simple_as_environment", env))))
			}()) {
				return true
			}
			return RT.Truth(RT.Call("isNull", env))
		}())
	}()) {
		if RT.Truth(RT.Call("isNull", env)) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"use of NULL environment is defunct\"")))
		}
		if RT.Truth(func() Value {
			if RT.Truth(RT.Call("MAYBE_SHARED", s)) {
				return true
			}
			return RT.Truth(func() Value {
				if !RT.Truth(RT.Unary("!", RT.Call("IS_ASSIGNMENT_CALL", call))) {
					return false
				}
				return RT.Truth(RT.Call("MAYBE_REFERENCED", s))
			}())
		}()) {
			Assign(LocalRef(&s), RT.Call("duplicate", s))
		}
		if RT.Truth(RT.Binary("==", RT.Call("TYPEOF", RT.Call("BODY", s)), RT.Symbol("BCODESXP"))) {
			RT.Call("SET_BODY", s, RT.Call("R_ClosureExpr", s))
		}
		RT.Call("SET_CLOENV", s, env)
	} else {
		if RT.Truth(func() Value {
			if RT.Truth(func() Value {
				if RT.Truth(RT.Call("isNull", env)) {
					return true
				}
				return RT.Truth(RT.Call("isEnvironment", env))
			}()) {
				return true
			}
			return RT.Truth(RT.Call("isEnvironment", Assign(LocalRef(&env), RT.Call("simple_as_environment", env))))
		}()) {
			if RT.Truth(func() Value {
				if !RT.Truth(RT.Unary("!", RT.Call("isNull", env))) {
					return false
				}
				return RT.Truth(RT.Call("isPrimitive", s))
			}()) {
				RT.Call("warning", RT.Call("_", RT.Const("string", "\"setting environment(<primitive function>) is not possible and trying it is deprecated\"")))
			} else {
				RT.Call("setAttrib", s, RT.Symbol("R_DotEnvSymbol"), env)
			}
		} else {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"replacement object is not an environment\"")))
		}
	}
	return s
}
