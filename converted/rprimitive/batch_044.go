package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_envirName(call, op, args, rho Value) Value {
	var (
		env Value
		ans Value
		res Value
	)
	Assign(LocalRef(&env), RT.Call("CAR", args))
	Assign(LocalRef(&ans), RT.Call("mkString", RT.Const("string", "\"\"")))
	RT.Call("checkArity", op, args)
	RT.Call("PROTECT", ans)
	if RT.Truth(func() Value {
		if RT.Truth(RT.Binary("==", RT.Call("TYPEOF", env), RT.Symbol("ENVSXP"))) {
			return true
		}
		return RT.Truth(RT.Binary("==", RT.Call("TYPEOF", Assign(LocalRef(&env), RT.Call("simple_as_environment", env))), RT.Symbol("ENVSXP")))
	}()) {
		if RT.Truth(RT.Binary("==", env, RT.Symbol("R_GlobalEnv"))) {
			Assign(LocalRef(&ans), RT.Call("mkString", RT.Const("string", "\"R_GlobalEnv\"")))
		} else {
			if RT.Truth(RT.Binary("==", env, RT.Symbol("R_BaseEnv"))) {
				Assign(LocalRef(&ans), RT.Call("mkString", RT.Const("string", "\"base\"")))
			} else {
				if RT.Truth(RT.Binary("==", env, RT.Symbol("R_EmptyEnv"))) {
					Assign(LocalRef(&ans), RT.Call("mkString", RT.Const("string", "\"R_EmptyEnv\"")))
				} else {
					if RT.Truth(RT.Call("R_IsPackageEnv", env)) {
						Assign(LocalRef(&ans), RT.Call("ScalarString", RT.Call("STRING_ELT", RT.Call("R_PackageEnvName", env), RT.Const("int", "0"))))
					} else {
						if RT.Truth(RT.Call("R_IsNamespaceEnv", env)) {
							Assign(LocalRef(&ans), RT.Call("ScalarString", RT.Call("STRING_ELT", RT.Call("R_NamespaceEnvSpec", env), RT.Const("int", "0"))))
						} else {
							if RT.Truth(RT.Unary("!", RT.Call("isNull", Assign(LocalRef(&res), RT.Call("getAttrib", env, RT.Symbol("R_NameSymbol")))))) {
								Assign(LocalRef(&ans), res)
							}
						}
					}
				}
			}
		}
	}
	RT.Call("UNPROTECT", RT.Const("int", "1"))
	return ans
}
