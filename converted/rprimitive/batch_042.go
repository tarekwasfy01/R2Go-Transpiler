package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_env2list(call, op, args, rho Value) Value {
	var (
		env      Value
		ans      Value
		names    Value
		k        Value
		all      Value
		xdata    Value
		sort_nms Value
		sind     Value
		indx     Value
		i        Value
		ans2     Value
		names2   Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&env), RT.Call("CAR", args))
	if RT.Truth(RT.Call("ISNULL", env)) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"use of NULL environment is defunct\"")))
	}
	if RT.Truth(RT.Unary("!", RT.Call("isEnvironment", env))) {
		if RT.Truth(func() Value {
			if !RT.Truth(func() Value {
				if !RT.Truth(RT.Call("IS_S4_OBJECT", env)) {
					return false
				}
				return RT.Truth(RT.Binary("==", RT.Call("TYPEOF", env), RT.Symbol("OBJSXP")))
			}()) {
				return false
			}
			return RT.Truth(RT.Binary("!=", Assign(LocalRef(&xdata), RT.Call("R_getS4DataSlot", env, RT.Symbol("ENVSXP"))), RT.Symbol("R_NilValue")))
		}()) {
			Assign(LocalRef(&env), xdata)
		} else {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"argument must be an environment\"")))
		}
	}
	Assign(LocalRef(&all), RT.Call("asLogical", RT.Call("CADR", args)))
	if RT.Truth(RT.Binary("==", all, RT.Symbol("NA_LOGICAL"))) {
		Assign(LocalRef(&all), RT.Const("int", "0"))
	}
	Assign(LocalRef(&sort_nms), RT.Call("asLogical", RT.Call("CADDR", args)))
	if RT.Truth(RT.Binary("==", sort_nms, RT.Symbol("NA_LOGICAL"))) {
		Assign(LocalRef(&sort_nms), RT.Const("int", "0"))
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
	RT.Call("PROTECT", Assign(LocalRef(&names), RT.Call("allocVector", RT.Symbol("STRSXP"), k)))
	RT.Call("PROTECT", Assign(LocalRef(&ans), RT.Call("allocVector", RT.Symbol("VECSXP"), k)))
	Assign(LocalRef(&k), RT.Const("int", "0"))
	if RT.Truth(func() Value {
		if RT.Truth(RT.Binary("==", env, RT.Symbol("R_BaseEnv"))) {
			return true
		}
		return RT.Truth(RT.Binary("==", env, RT.Symbol("R_BaseNamespace")))
	}()) {
		RT.Call("BuiltinValues", all, RT.Const("int", "0"), ans, LocalRef(&k))
	} else {
		if RT.Truth(RT.Binary("!=", RT.Call("HASHTAB", env), RT.Symbol("R_NilValue"))) {
			RT.Call("HashTableValues", RT.Call("HASHTAB", env), all, ans, LocalRef(&k))
		} else {
			RT.Call("FrameValues", RT.Call("FRAME", env), all, ans, LocalRef(&k))
		}
	}
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
	if RT.Truth(RT.Binary("==", k, RT.Const("int", "0"))) {
		RT.Call("UNPROTECT", RT.Const("int", "2"))
		return ans
	}
	if RT.Truth(sort_nms) {
		Assign(LocalRef(&sind), RT.Call("PROTECT", RT.Call("allocVector", RT.Symbol("INTSXP"), k)))
		Assign(LocalRef(&indx), RT.Call("INTEGER", sind))
		for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, k)); RT.Inc(LocalRef(&i), 1, true) {
			RT.AssignIndex(indx, i, i)
		}
		RT.Call("orderVector1", indx, k, names, RT.Symbol("true"), RT.Symbol("false"), RT.Symbol("R_NilValue"))
		Assign(LocalRef(&ans2), RT.Call("PROTECT", RT.Call("allocVector", RT.Symbol("VECSXP"), k)))
		Assign(LocalRef(&names2), RT.Call("PROTECT", RT.Call("allocVector", RT.Symbol("STRSXP"), k)))
		for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, k)); RT.Inc(LocalRef(&i), 1, true) {
			RT.Call("SET_STRING_ELT", names2, i, RT.Call("STRING_ELT", names, RT.Index(indx, i)))
			RT.Call("SET_VECTOR_ELT", ans2, i, RT.Call("VECTOR_ELT", ans, RT.Index(indx, i)))
		}
		RT.Call("setAttrib", ans2, RT.Symbol("R_NamesSymbol"), names2)
		RT.Call("UNPROTECT", RT.Const("int", "5"))
		return ans2
	} else {
		RT.Call("setAttrib", ans, RT.Symbol("R_NamesSymbol"), names)
		RT.Call("UNPROTECT", RT.Const("int", "2"))
		return ans
	}
}
