package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_dotcall(call, op, args, env Value) Value {
	var (
		ofun     Value
		retval   Value
		cargs    Value
		pargs    Value
		symbol   Value
		nargs    Value
		vmax     Value
		buf      Value
		nprotect Value
		cargscp  Value
		i        Value
		constsOK Value
	)
	Assign(LocalRef(&ofun), RT.Symbol("NULL"))
	Assign(LocalRef(&cargs), RT.NewArray(RT.Symbol("MAX_ARGS")))
	Assign(LocalRef(&symbol), RT.List(RT.Symbol("R_CALL_SYM"), RT.List(RT.Symbol("NULL")), RT.Symbol("NULL")))
	Assign(LocalRef(&vmax), RT.Call("vmaxget"))
	Assign(LocalRef(&buf), RT.NewArray(RT.Symbol("MaxSymbolBytes")))
	Assign(LocalRef(&nprotect), RT.Const("int", "0"))
	if RT.Truth(RT.Binary("<", RT.Call("length", args), RT.Const("int", "1"))) {
		RT.Call("errorcall", call, RT.Call("_", RT.Const("string", "\"'.NAME' is missing\"")))
	}
	RT.Call("check1arg2", args, call, RT.Const("string", "\".NAME\""))
	Assign(LocalRef(&args), RT.Call("resolveNativeRoutine", args, LocalRef(&ofun), LocalRef(&symbol), buf, RT.Symbol("NULL"), RT.Symbol("NULL"), call, env))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	for RT.Sequence(Assign(LocalRef(&nargs), RT.Const("int", "0")), Assign(LocalRef(&pargs), args)); RT.Truth(RT.Binary("!=", pargs, RT.Symbol("R_NilValue"))); Assign(LocalRef(&pargs), RT.Call("CDR", pargs)) {
		if RT.Truth(RT.Binary("==", nargs, RT.Symbol("MAX_ARGS"))) {
			RT.Call("errorcall", call, RT.Call("_", RT.Const("string", "\"too many arguments in foreign function call\"")))
		}
		RT.AssignIndex(cargs, nargs, RT.Call("CAR", pargs))
		RT.Inc(LocalRef(&nargs), 1, true)
	}
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Field(RT.Field(symbol, "symbol"), "call")) {
			return false
		}
		return RT.Truth(RT.Binary(">", RT.Field(RT.Field(RT.Field(symbol, "symbol"), "call"), "numArgs"), RT.Unary("-", RT.Const("int", "1"))))
	}()) {
		if RT.Truth(RT.Binary("!=", RT.Field(RT.Field(RT.Field(symbol, "symbol"), "call"), "numArgs"), nargs)) {
			RT.Call("errorcall", call, RT.Call("_", RT.Const("string", "\"Incorrect number of arguments (%d), expecting %d for '%s'\"")), nargs, RT.Field(RT.Field(RT.Field(symbol, "symbol"), "call"), "numArgs"), buf)
		}
	}
	if RT.Truth(RT.Binary("<", RT.Symbol("R_check_constants"), RT.Const("int", "4"))) {
		Assign(LocalRef(&retval), RT.Call("R_doDotCall", ofun, nargs, cargs, call))
	} else {
		Assign(LocalRef(&cargscp), RT.Cast("SEXP *", RT.Call("R_alloc", nargs, RT.SizeOfType("SEXP"))))
		for Assign(LocalRef(&i), RT.Const("int", "0")); RT.Truth(RT.Binary("<", i, nargs)); RT.Inc(LocalRef(&i), 1, true) {
			RT.AssignIndex(cargscp, i, RT.Call("PROTECT", RT.Call("duplicate", RT.Index(cargs, i))))
			RT.Inc(LocalRef(&nprotect), 1, true)
		}
		Assign(LocalRef(&retval), RT.Call("PROTECT", RT.Call("R_doDotCall", ofun, nargs, cargs, call)))
		RT.Inc(LocalRef(&nprotect), 1, true)
		Assign(LocalRef(&constsOK), RT.Symbol("true"))
		for Assign(LocalRef(&i), RT.Const("int", "0")); RT.Truth(func() Value {
			if !RT.Truth(constsOK) {
				return false
			}
			return RT.Truth(RT.Binary("<", i, nargs))
		}()); RT.Inc(LocalRef(&i), 1, true) {
			if RT.Truth(func() Value {
				if !RT.Truth(RT.Unary("!", RT.Call("R_compute_identical", RT.Index(cargs, i), RT.Index(cargscp, i), RT.Const("int", "39")))) {
					return false
				}
				return RT.Truth(RT.Unary("!", RT.Call("R_checkConstants", RT.Symbol("FALSE"))))
			}()) {
				Assign(LocalRef(&constsOK), RT.Symbol("false"))
			}
		}
		if RT.Truth(RT.Unary("!", constsOK)) {
			RT.Call("REprintf", RT.Const("string", "\"ERROR: detected compiler constant(s) modification after .Call invocation of function %s from library %s (%s).\\n\""), buf, func() Value {
				if RT.Truth(RT.Field(symbol, "dll")) {
					return RT.Field(RT.Field(symbol, "dll"), "name")
				}
				return RT.Const("string", "\"unknown\"")
			}(), func() Value {
				if RT.Truth(RT.Field(symbol, "dll")) {
					return RT.Field(RT.Field(symbol, "dll"), "path")
				}
				return RT.Const("string", "\"unknown\"")
			}())
			for Assign(LocalRef(&i), RT.Const("int", "0")); RT.Truth(RT.Binary("<", i, nargs)); RT.Inc(LocalRef(&i), 1, true) {
				if RT.Truth(RT.Unary("!", RT.Call("R_compute_identical", RT.Index(cargs, i), RT.Index(cargscp, i), RT.Const("int", "39")))) {
					RT.Call("REprintf", RT.Const("string", "\"NOTE: .Call function %s modified its argument (number %d, type %s, length %d)\\n\""), buf, RT.Binary("+", i, RT.Const("int", "1")), RT.Call("CHAR", RT.Call("type2str", RT.Call("TYPEOF", RT.Index(cargscp, i)))), RT.Call("length", RT.Index(cargscp, i)))
				}
			}
			RT.Call("R_Suicide", RT.Const("string", "\"compiler constants were modified (in .Call?)!\\n\""))
		}
		RT.Call("UNPROTECT", nprotect)
	}
	RT.Call("vmaxset", vmax)
	return retval
}
