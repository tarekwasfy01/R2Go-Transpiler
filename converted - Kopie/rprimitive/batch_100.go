package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_search(call, op, args, env Value) Value {
	var (
		ans  Value
		name Value
		t    Value
		i    Value
		n    Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&n), RT.Const("int", "2"))
	for Assign(LocalRef(&t), RT.Call("ENCLOS", RT.Symbol("R_GlobalEnv"))); RT.Truth(RT.Binary("!=", t, RT.Symbol("R_BaseEnv"))); Assign(LocalRef(&t), RT.Call("ENCLOS", t)) {
		RT.Inc(LocalRef(&n), 1, true)
	}
	RT.Call("PROTECT", Assign(LocalRef(&ans), RT.Call("allocVector", RT.Symbol("STRSXP"), n)))
	RT.Call("SET_STRING_ELT", ans, RT.Const("int", "0"), RT.Call("mkChar", RT.Const("string", "\".GlobalEnv\"")))
	RT.Call("SET_STRING_ELT", ans, RT.Binary("-", n, RT.Const("int", "1")), RT.Call("mkChar", RT.Const("string", "\"package:base\"")))
	Assign(LocalRef(&i), RT.Const("int", "1"))
	for Assign(LocalRef(&t), RT.Call("ENCLOS", RT.Symbol("R_GlobalEnv"))); RT.Truth(RT.Binary("!=", t, RT.Symbol("R_BaseEnv"))); Assign(LocalRef(&t), RT.Call("ENCLOS", t)) {
		Assign(LocalRef(&name), RT.Call("getAttrib", t, RT.Symbol("R_NameSymbol")))
		if RT.Truth(func() Value {
			if RT.Truth(RT.Unary("!", RT.Call("isString", name))) {
				return true
			}
			return RT.Truth(RT.Binary("<", RT.Call("length", name), RT.Const("int", "1")))
		}()) {
			RT.Call("SET_STRING_ELT", ans, i, RT.Call("mkChar", RT.Const("string", "\"(unknown)\"")))
		} else {
			RT.Call("SET_STRING_ELT", ans, i, RT.Call("STRING_ELT", name, RT.Const("int", "0")))
		}
		RT.Inc(LocalRef(&i), 1, true)
	}
	RT.Call("UNPROTECT", RT.Const("int", "1"))
	return ans
}
