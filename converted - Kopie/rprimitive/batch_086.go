package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_parentenvgets(call, op, args, rho Value) Value {
	var (
		env    Value
		parent Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&env), RT.Call("CAR", args))
	if RT.Truth(RT.Call("isNull", env)) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"use of NULL environment is defunct\"")))
		Assign(LocalRef(&env), RT.Symbol("R_BaseEnv"))
	} else {
		if RT.Truth(func() Value {
			if !RT.Truth(RT.Unary("!", RT.Call("isEnvironment", env))) {
				return false
			}
			return RT.Truth(RT.Unary("!", RT.Call("isEnvironment", Assign(LocalRef(&env), RT.Call("simple_as_environment", env)))))
		}()) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"argument is not an environment\"")))
		}
	}
	if RT.Truth(RT.Binary("==", env, RT.Symbol("R_EmptyEnv"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"can not set parent of the empty environment\"")))
	}
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Call("R_EnvironmentIsLocked", env)) {
			return false
		}
		return RT.Truth(RT.Call("R_IsNamespaceEnv", env))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"can not set the parent environment of a namespace\"")))
	}
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Call("R_EnvironmentIsLocked", env)) {
			return false
		}
		return RT.Truth(RT.Call("R_IsImportsEnv", env))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"can not set the parent environment of package imports\"")))
	}
	Assign(LocalRef(&parent), RT.Call("CADR", args))
	if RT.Truth(RT.Call("isNull", parent)) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"use of NULL environment is defunct\"")))
		Assign(LocalRef(&parent), RT.Symbol("R_BaseEnv"))
	} else {
		if RT.Truth(func() Value {
			if !RT.Truth(RT.Unary("!", RT.Call("isEnvironment", parent))) {
				return false
			}
			return RT.Truth(RT.Unary("!", RT.Call("isEnvironment", Assign(LocalRef(&parent), RT.Call("simple_as_environment", parent)))))
		}()) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"'parent' is not an environment\"")))
		}
	}
	RT.Call("SET_ENCLOS", env, parent)
	return RT.Call("CAR", args)
}
