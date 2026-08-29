package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_mget(call, op, args, rho Value) Value {
	var (
		ans        Value
		env        Value
		x          Value
		mode       Value
		ifnotfound Value
		ginherits  Value
		nvals      Value
		nmode      Value
		nifnfnd    Value
		i          Value
		wants_S4   Value
		modestr    Value
		gmode      Value
		nf         Value
		ans_i      Value
	)
	Assign(LocalRef(&ginherits), RT.Const("int", "0"))
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&x), RT.Call("CAR", args))
	Assign(LocalRef(&nvals), RT.Call("length", x))
	if RT.Truth(RT.Unary("!", RT.Call("isString", x))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid first argument\"")))
	}
	for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, nvals)); RT.Inc(LocalRef(&i), 1, true) {
		if RT.Truth(func() Value {
			if RT.Truth(RT.Call("isNull", RT.Call("STRING_ELT", x, i))) {
				return true
			}
			return RT.Truth(RT.Unary("!", RT.Index(RT.Call("CHAR", RT.Call("STRING_ELT", x, RT.Const("int", "0"))), RT.Const("int", "0"))))
		}()) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid name in position %d\"")), RT.Binary("+", i, RT.Const("int", "1")))
		}
	}
	Assign(LocalRef(&env), RT.Call("CADR", args))
	if RT.Truth(RT.Call("ISNULL", env)) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"use of NULL environment is defunct\"")))
	} else {
		if RT.Truth(RT.Unary("!", RT.Call("isEnvironment", env))) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"second argument must be an environment\"")))
		}
	}
	Assign(LocalRef(&mode), RT.Call("CADDR", args))
	Assign(LocalRef(&nmode), RT.Call("length", mode))
	if RT.Truth(RT.Unary("!", RT.Call("isString", mode))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"mode\""))
	}
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Binary("!=", nmode, nvals)) {
			return false
		}
		return RT.Truth(RT.Binary("!=", nmode, RT.Const("int", "1")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"wrong length for '%s' argument\"")), RT.Const("string", "\"mode\""))
	}
	RT.Call("PROTECT", Assign(LocalRef(&ifnotfound), RT.Call("coerceVector", RT.Call("CADDDR", args), RT.Symbol("VECSXP"))))
	Assign(LocalRef(&nifnfnd), RT.Call("length", ifnotfound))
	if RT.Truth(RT.Unary("!", RT.Call("isVector", ifnotfound))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"ifnotfound\""))
	}
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Binary("!=", nifnfnd, nvals)) {
			return false
		}
		return RT.Truth(RT.Binary("!=", nifnfnd, RT.Const("int", "1")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"wrong length for '%s' argument\"")), RT.Const("string", "\"ifnotfound\""))
	}
	Assign(LocalRef(&ginherits), RT.Call("asLogical", RT.Call("CAD4R", args)))
	if RT.Truth(RT.Binary("==", ginherits, RT.Symbol("NA_LOGICAL"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"inherits\""))
	}
	RT.Call("PROTECT", Assign(LocalRef(&ans), RT.Call("allocVector", RT.Symbol("VECSXP"), nvals)))
	for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, nvals)); RT.Inc(LocalRef(&i), 1, true) {
		Assign(LocalRef(&wants_S4), RT.Symbol("FALSE"))
		Assign(LocalRef(&modestr), RT.Call("CHAR", RT.Call("STRING_ELT", RT.Call("CADDR", args), RT.Binary("%", i, nmode))))
		Assign(LocalRef(&gmode), RT.Call("str2mode", modestr, LocalRef(&wants_S4)))
		Assign(LocalRef(&nf), RT.Call("VECTOR_ELT", ifnotfound, RT.Binary("%", i, nifnfnd)))
		Assign(LocalRef(&ans_i), RT.Call("gfind", RT.Call("translateChar", RT.Call("STRING_ELT", x, RT.Binary("%", i, nvals))), env, gmode, wants_S4, nf, ginherits, rho))
		RT.Call("SET_VECTOR_ELT", ans, i, RT.Call("lazy_duplicate", ans_i))
	}
	RT.Call("setAttrib", ans, RT.Symbol("R_NamesSymbol"), RT.Call("lazy_duplicate", x))
	RT.Call("UNPROTECT", RT.Const("int", "2"))
	return ans
}
