package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_eapply(call, op, args, rho Value) Value {
	var (
		env     Value
		ans     Value
		R_fcall Value
		FUN     Value
		tmp     Value
		tmp2    Value
		ind     Value
		i       Value
		k       Value
		k2      Value
		all     Value
		useNms  Value
		Xsym    Value
		isym    Value
		names   Value
	)
	RT.Call("checkArity", op, args)
	RT.Call("PROTECT", Assign(LocalRef(&env), RT.Call("eval", RT.Call("CAR", args), rho)))
	if RT.Truth(RT.Call("ISNULL", env)) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"use of NULL environment is defunct\"")))
	}
	if RT.Truth(RT.Unary("!", RT.Call("isEnvironment", env))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"argument must be an environment\"")))
	}
	Assign(LocalRef(&FUN), RT.Call("CADR", args))
	if RT.Truth(RT.Unary("!", RT.Call("isSymbol", FUN))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"arguments must be symbolic\"")))
	}
	Assign(LocalRef(&all), RT.Call("asLogical", RT.Call("PROTECT", RT.Call("eval", RT.Call("CADDR", args), rho))))
	RT.Call("UNPROTECT", RT.Const("int", "1"))
	if RT.Truth(RT.Binary("==", all, RT.Symbol("NA_LOGICAL"))) {
		Assign(LocalRef(&all), RT.Const("int", "0"))
	}
	Assign(LocalRef(&useNms), RT.Call("asLogical", RT.Call("PROTECT", RT.Call("eval", RT.Call("CADDDR", args), rho))))
	RT.Call("UNPROTECT", RT.Const("int", "1"))
	if RT.Truth(RT.Binary("==", useNms, RT.Symbol("NA_LOGICAL"))) {
		Assign(LocalRef(&useNms), RT.Const("int", "0"))
	}
	if RT.Truth(func() Value {
		if RT.Truth(RT.Binary("==", env, RT.Symbol("R_BaseEnv"))) {
			return true
		}
		return RT.Truth(RT.Binary("==", env, RT.Symbol("R_BaseNamespace")))
	}()) {
		Assign(LocalRef(&k), RT.Call("BuiltinSize", all, RT.Const("int", "0")))
	} else {
		if RT.Truth(RT.Binary("!=", RT.Call("HASHTAB", env), RT.Symbol("R_NilValue"))) {
			Assign(LocalRef(&k), RT.Call("HashTableSize", RT.Call("HASHTAB", env), all))
		} else {
			Assign(LocalRef(&k), RT.Call("FrameSize", RT.Call("FRAME", env), all))
		}
	}
	RT.Call("PROTECT", Assign(LocalRef(&ans), RT.Call("allocVector", RT.Symbol("VECSXP"), k)))
	RT.Call("PROTECT", Assign(LocalRef(&tmp2), RT.Call("allocVector", RT.Symbol("VECSXP"), k)))
	Assign(LocalRef(&k2), RT.Const("int", "0"))
	if RT.Truth(func() Value {
		if RT.Truth(RT.Binary("==", env, RT.Symbol("R_BaseEnv"))) {
			return true
		}
		return RT.Truth(RT.Binary("==", env, RT.Symbol("R_BaseNamespace")))
	}()) {
		RT.Call("BuiltinValues", all, RT.Const("int", "0"), tmp2, LocalRef(&k2))
	} else {
		if RT.Truth(RT.Binary("!=", RT.Call("HASHTAB", env), RT.Symbol("R_NilValue"))) {
			RT.Call("HashTableValues", RT.Call("HASHTAB", env), all, tmp2, LocalRef(&k2))
		} else {
			RT.Call("FrameValues", RT.Call("FRAME", env), all, tmp2, LocalRef(&k2))
		}
	}
	Assign(LocalRef(&Xsym), RT.Call("install", RT.Const("string", "\"X\"")))
	Assign(LocalRef(&isym), RT.Call("install", RT.Const("string", "\"i\"")))
	RT.Call("PROTECT", Assign(LocalRef(&ind), RT.Call("allocVector", RT.Symbol("INTSXP"), RT.Const("int", "1"))))
	RT.Call("PROTECT", Assign(LocalRef(&tmp), RT.Call("LCONS", RT.Symbol("R_Bracket2Symbol"), RT.Call("LCONS", Xsym, RT.Call("LCONS", isym, RT.Symbol("R_NilValue"))))))
	RT.Call("PROTECT", Assign(LocalRef(&R_fcall), RT.Call("LCONS", FUN, RT.Call("LCONS", tmp, RT.Call("LCONS", RT.Symbol("R_DotsSymbol"), RT.Symbol("R_NilValue"))))))
	RT.Call("defineVar", Xsym, tmp2, rho)
	RT.Call("INCREMENT_NAMED", tmp2)
	RT.Call("defineVar", isym, ind, rho)
	RT.Call("INCREMENT_NAMED", ind)
	for Assign(LocalRef(&i), RT.Const("int", "0")); RT.Truth(RT.Binary("<", i, k2)); RT.Inc(LocalRef(&i), 1, true) {
		RT.AssignIndex(RT.Call("INTEGER", ind), RT.Const("int", "0"), RT.Binary("+", i, RT.Const("int", "1")))
		Assign(LocalRef(&tmp), RT.Call("R_forceAndCall", R_fcall, RT.Const("int", "1"), rho))
		if RT.Truth(RT.Call("MAYBE_REFERENCED", tmp)) {
			Assign(LocalRef(&tmp), RT.Call("lazy_duplicate", tmp))
		}
		RT.Call("SET_VECTOR_ELT", ans, i, tmp)
	}
	if RT.Truth(useNms) {
		RT.Call("PROTECT", Assign(LocalRef(&names), RT.Call("allocVector", RT.Symbol("STRSXP"), k)))
		Assign(LocalRef(&k), RT.Const("int", "0"))
		if RT.Truth(func() Value {
			if RT.Truth(RT.Binary("==", env, RT.Symbol("R_BaseEnv"))) {
				return true
			}
			return RT.Truth(RT.Binary("==", env, RT.Symbol("R_BaseNamespace")))
		}()) {
			RT.Call("BuiltinNames", all, RT.Const("int", "0"), names, LocalRef(&k))
		} else {
			if RT.Truth(RT.Binary("!=", RT.Call("HASHTAB", env), RT.Symbol("R_NilValue"))) {
				RT.Call("HashTableNames", RT.Call("HASHTAB", env), all, names, LocalRef(&k))
			} else {
				RT.Call("FrameNames", RT.Call("FRAME", env), all, names, LocalRef(&k))
			}
		}
		RT.Call("setAttrib", ans, RT.Symbol("R_NamesSymbol"), names)
		RT.Call("UNPROTECT", RT.Const("int", "1"))
	}
	RT.Call("UNPROTECT", RT.Const("int", "6"))
	return ans
}
