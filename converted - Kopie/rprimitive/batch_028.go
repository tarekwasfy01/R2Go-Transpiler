package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_detach(call, op, args, env Value) Value {
	var (
		s         Value
		t         Value
		x         Value
		pos       Value
		n         Value
		isSpecial Value
		tb        Value
	)
	Assign(LocalRef(&isSpecial), RT.Symbol("FALSE"))
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&pos), RT.Call("asInteger", RT.Call("CAR", args)))
	for RT.Sequence(Assign(LocalRef(&n), RT.Const("int", "2")), Assign(LocalRef(&t), RT.Call("ENCLOS", RT.Symbol("R_GlobalEnv")))); RT.Truth(RT.Binary("!=", t, RT.Symbol("R_BaseEnv"))); Assign(LocalRef(&t), RT.Call("ENCLOS", t)) {
		RT.Inc(LocalRef(&n), 1, true)
	}
	if RT.Truth(RT.Binary("==", pos, n)) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"detaching \\\"package:base\\\" is not allowed\"")))
	}
	for Assign(LocalRef(&t), RT.Symbol("R_GlobalEnv")); RT.Truth(func() Value {
		if !RT.Truth(RT.Binary("!=", RT.Call("ENCLOS", t), RT.Symbol("R_BaseEnv"))) {
			return false
		}
		return RT.Truth(RT.Binary(">", pos, RT.Const("int", "2")))
	}()); Assign(LocalRef(&t), RT.Call("ENCLOS", t)) {
		RT.Inc(LocalRef(&pos), -1, true)
	}
	if RT.Truth(RT.Binary("!=", pos, RT.Const("int", "2"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"pos\""))
		Assign(LocalRef(&s), t)
	} else {
		RT.Call("PROTECT", Assign(LocalRef(&s), RT.Call("ENCLOS", t)))
		Assign(LocalRef(&x), RT.Call("ENCLOS", s))
		RT.Call("SET_ENCLOS", t, x)
		Assign(LocalRef(&isSpecial), RT.Call("IS_USER_DATABASE", s))
		if RT.Truth(isSpecial) {
			Assign(LocalRef(&tb), RT.Cast("R_ObjectTable *", RT.Call("R_ExternalPtrAddr", RT.Call("HASHTAB", s))))
			if RT.Truth(RT.Field(tb, "onDetach")) {
				RT.CallIndirect(RT.Field(tb, "onDetach"), tb)
			}
		}
		RT.Call("SET_ENCLOS", s, RT.Symbol("R_BaseEnv"))
	}
	RT.Call("UNPROTECT", RT.Const("int", "1"))
	return s
}
