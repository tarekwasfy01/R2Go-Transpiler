package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_attach(call, op, args, env Value) Value {
	var (
		name      Value
		s         Value
		t         Value
		x         Value
		pos       Value
		hsize     Value
		isSpecial Value
		p         Value
		loadenv   Value
		i         Value
		n         Value
		tb        Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&pos), RT.Call("asInteger", RT.Call("CADR", args)))
	if RT.Truth(RT.Binary("==", pos, RT.Symbol("NA_INTEGER"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"'pos' must be an integer\"")))
	}
	Assign(LocalRef(&name), RT.Call("CADDR", args))
	if RT.Truth(RT.Unary("!", RT.Call("isValidStringF", name))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"name\""))
	}
	Assign(LocalRef(&isSpecial), RT.Call("IS_USER_DATABASE", RT.Call("CAR", args)))
	if RT.Truth(RT.Unary("!", isSpecial)) {
		if RT.Truth(RT.Call("isNewList", RT.Call("CAR", args))) {
			RT.Call("SETCAR", args, RT.Call("VectorToPairList", RT.Call("CAR", args)))
			for Assign(LocalRef(&x), RT.Call("CAR", args)); RT.Truth(RT.Binary("!=", x, RT.Symbol("R_NilValue"))); Assign(LocalRef(&x), RT.Call("CDR", x)) {
				if RT.Truth(RT.Binary("==", RT.Call("TAG", x), RT.Symbol("R_NilValue"))) {
					RT.Call("error", RT.Call("_", RT.Const("string", "\"all elements of a list must be named\"")))
				}
			}
			RT.Call("PROTECT", Assign(LocalRef(&s), RT.Call("allocSExp", RT.Symbol("ENVSXP"))))
			RT.Call("SET_FRAME", s, RT.Call("shallow_duplicate", RT.Call("CAR", args)))
		} else {
			if RT.Truth(RT.Call("isEnvironment", RT.Call("CAR", args))) {
				Assign(LocalRef(&loadenv), RT.Call("CAR", args))
				RT.Call("PROTECT", Assign(LocalRef(&s), RT.Call("allocSExp", RT.Symbol("ENVSXP"))))
				if RT.Truth(RT.Binary("!=", RT.Call("HASHTAB", loadenv), RT.Symbol("R_NilValue"))) {
					Assign(LocalRef(&n), RT.Call("length", RT.Call("HASHTAB", loadenv)))
					for Assign(LocalRef(&i), RT.Const("int", "0")); RT.Truth(RT.Binary("<", i, n)); RT.Inc(LocalRef(&i), 1, true) {
						Assign(LocalRef(&p), RT.Call("VECTOR_ELT", RT.Call("HASHTAB", loadenv), i))
						for RT.Truth(RT.Binary("!=", p, RT.Symbol("R_NilValue"))) {
							RT.Call("set_attach_frame_value", p, s)
							Assign(LocalRef(&p), RT.Call("CDR", p))
						}
					}
				} else {
					for Assign(LocalRef(&p), RT.Call("FRAME", loadenv)); RT.Truth(RT.Binary("!=", p, RT.Symbol("R_NilValue"))); Assign(LocalRef(&p), RT.Call("CDR", p)) {
						RT.Call("set_attach_frame_value", p, s)
					}
				}
			} else {
				RT.Call("error", RT.Call("_", RT.Const("string", "\"'attach' only works for lists, data frames and environments\"")))
				Assign(LocalRef(&s), RT.Symbol("R_NilValue"))
			}
		}
		if RT.Truth(RT.Binary("<", RT.Call("length", s), RT.Symbol("HASHMINSIZE"))) {
			Assign(LocalRef(&hsize), RT.Symbol("HASHMINSIZE"))
		} else {
			Assign(LocalRef(&hsize), RT.Call("length", s))
		}
		RT.Call("SET_HASHTAB", s, RT.Call("R_NewHashTable", hsize))
		Assign(LocalRef(&s), RT.Call("R_HashFrame", s))
		for RT.Truth(RT.Call("R_HashSizeCheck", RT.Call("HASHTAB", s))) {
			RT.Call("SET_HASHTAB", s, RT.Call("R_HashResize", RT.Call("HASHTAB", s)))
		}
	} else {
		Assign(LocalRef(&tb), RT.Cast("R_ObjectTable *", RT.Call("R_ExternalPtrAddr", RT.Call("CAR", args))))
		if RT.Truth(RT.Field(tb, "onAttach")) {
			RT.CallIndirect(RT.Field(tb, "onAttach"), tb)
		}
		RT.Call("PROTECT", Assign(LocalRef(&s), RT.Call("allocSExp", RT.Symbol("ENVSXP"))))
		RT.Call("SET_HASHTAB", s, RT.Call("CAR", args))
		RT.Call("setAttrib", s, RT.Symbol("R_ClassSymbol"), RT.Call("getAttrib", RT.Call("HASHTAB", s), RT.Symbol("R_ClassSymbol")))
	}
	RT.Call("setAttrib", s, RT.Symbol("R_NameSymbol"), name)
	for Assign(LocalRef(&t), RT.Symbol("R_GlobalEnv")); RT.Truth(func() Value {
		if !RT.Truth(RT.Binary("!=", RT.Call("ENCLOS", t), RT.Symbol("R_BaseEnv"))) {
			return false
		}
		return RT.Truth(RT.Binary(">", pos, RT.Const("int", "2")))
	}()); Assign(LocalRef(&t), RT.Call("ENCLOS", t)) {
		RT.Inc(LocalRef(&pos), -1, true)
	}
	if RT.Truth(RT.Binary("==", RT.Call("ENCLOS", t), RT.Symbol("R_BaseEnv"))) {
		RT.Call("SET_ENCLOS", t, s)
		RT.Call("SET_ENCLOS", s, RT.Symbol("R_BaseEnv"))
	} else {
		Assign(LocalRef(&x), RT.Call("ENCLOS", t))
		RT.Call("SET_ENCLOS", t, s)
		RT.Call("SET_ENCLOS", s, x)
	}
	if RT.Truth(RT.Unary("!", isSpecial)) {
	} else {
	}
	RT.Call("UNPROTECT", RT.Const("int", "1"))
	return s
}
