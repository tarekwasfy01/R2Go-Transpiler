package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_dump(call, op, args, rho Value) Value {
	var (
		names      Value
		file       Value
		nobjs      Value
		source     Value
		opts       Value
		objs       Value
		o          Value
		nout       Value
		i          Value
		outnames   Value
		obj_name   Value
		tval       Value
		j          Value
		con        Value
		wasopen    Value
		cntxt      Value
		mode       Value
		havewarned Value
		res        Value
		s          Value
		extra      Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&names), RT.Call("CAR", args))
	Assign(LocalRef(&file), RT.Call("CADR", args))
	if RT.Truth(RT.Unary("!", RT.Call("inherits", file, RT.Const("string", "\"connection\"")))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"'file' must be a character string or connection\"")))
	}
	if RT.Truth(RT.Unary("!", RT.Call("isString", names))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"character arguments expected\"")))
	}
	Assign(LocalRef(&nobjs), RT.Call("length", names))
	if RT.Truth(func() Value {
		if RT.Truth(RT.Binary("<", nobjs, RT.Const("int", "1"))) {
			return true
		}
		return RT.Truth(RT.Binary("<", RT.Call("length", file), RT.Const("int", "1")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"zero-length argument\"")))
	}
	Assign(LocalRef(&source), RT.Call("CADDR", args))
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Binary("!=", source, RT.Symbol("R_NilValue"))) {
			return false
		}
		return RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", source), RT.Symbol("ENVSXP")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"envir\""))
	}
	Assign(LocalRef(&opts), RT.Call("asInteger", RT.Call("CADDDR", args)))
	if RT.Truth(func() Value {
		if RT.Truth(func() Value {
			if RT.Truth(RT.Binary("==", opts, RT.Symbol("NA_INTEGER"))) {
				return true
			}
			return RT.Truth(RT.Binary("<", opts, RT.Const("int", "0")))
		}()) {
			return true
		}
		return RT.Truth(RT.Binary(">", opts, RT.Const("int", "2048")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"'opts' should be small non-negative integer\"")))
	}
	if RT.Truth(RT.Unary("!", RT.Call("asLogical", RT.Call("CAD4R", args)))) {
		Assign(LocalRef(&opts), RT.Binary("|", opts, RT.Symbol("DELAYPROMISES")))
	}
	Assign(LocalRef(&o), RT.Call("PROTECT", Assign(LocalRef(&objs), RT.Call("allocList", nobjs))))
	Assign(LocalRef(&nout), RT.Const("int", "0"))
	for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, nobjs)); RT.Sequence(RT.Inc(LocalRef(&i), 1, true), Assign(LocalRef(&o), RT.Call("CDR", o))) {
		RT.Call("SET_TAG", o, RT.Call("installTrChar", RT.Call("STRING_ELT", names, i)))
		RT.Call("SETCAR", o, RT.Call("R_findVar", RT.Call("TAG", o), source))
		if RT.Truth(RT.Binary("==", RT.Call("CAR", o), RT.Symbol("R_UnboundValue"))) {
			RT.Call("warning", RT.Call("_", RT.Const("string", "\"object '%s' not found\"")), RT.Call("EncodeChar", RT.Call("PRINTNAME", RT.Call("TAG", o))))
		} else {
			RT.Inc(LocalRef(&nout), 1, true)
		}
	}
	Assign(LocalRef(&o), objs)
	Assign(LocalRef(&outnames), RT.Call("PROTECT", RT.Call("allocVector", RT.Symbol("STRSXP"), nout)))
	if RT.Truth(RT.Binary(">", nout, RT.Const("int", "0"))) {
		if RT.Truth(RT.Binary("==", RT.Index(RT.Call("INTEGER", file), RT.Const("int", "0")), RT.Const("int", "1"))) {
			for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0")), Assign(LocalRef(&nout), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, nobjs)); RT.Inc(LocalRef(&i), 1, true) {
				if RT.Truth(RT.Binary("==", RT.Call("CAR", o), RT.Symbol("R_UnboundValue"))) {
					continue
				}
				Assign(LocalRef(&obj_name), RT.Call("translateChar", RT.Call("STRING_ELT", names, i)))
				RT.Call("SET_STRING_ELT", outnames, RT.Inc(LocalRef(&nout), 1, true), RT.Call("STRING_ELT", names, i))
				if RT.Truth(RT.Call("isValidName", obj_name)) {
					RT.Call("Rprintf", RT.Const("string", "\"%s <-\\n\""), obj_name)
				} else {
					if RT.Truth(RT.Binary("&", opts, RT.Symbol("S_COMPAT"))) {
						RT.Call("Rprintf", RT.Const("string", "\"\\\"%s\\\" <-\\n\""), obj_name)
					} else {
						RT.Call("Rprintf", RT.Const("string", "\"`%s` <-\\n\""), obj_name)
					}
				}
				Assign(LocalRef(&tval), RT.Call("PROTECT", RT.Call("deparse1", RT.Call("CAR", o), RT.Const("int", "0"), opts)))
				for RT.Sequence(Assign(LocalRef(&j), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", j, RT.Call("LENGTH", tval))); RT.Inc(LocalRef(&j), 1, true) {
					RT.Call("Rprintf", RT.Const("string", "\"%s\\n\""), RT.Call("CHAR", RT.Call("STRING_ELT", tval, j)))
				}
				RT.Call("UNPROTECT", RT.Const("int", "1"))
				Assign(LocalRef(&o), RT.Call("CDR", o))
			}
		} else {
			Assign(LocalRef(&con), RT.Call("getConnection", RT.Index(RT.Call("INTEGER", file), RT.Const("int", "0"))))
			Assign(LocalRef(&wasopen), RT.Field(con, "isopen"))
			if RT.Truth(RT.Unary("!", wasopen)) {
				Assign(LocalRef(&mode), RT.NewArray(RT.Const("int", "5")))
				RT.Call("strcpy", mode, RT.Field(con, "mode"))
				RT.Call("strcpy", RT.Field(con, "mode"), RT.Const("string", "\"w\""))
				if RT.Truth(RT.Unary("!", RT.CallIndirect(RT.Field(con, "open"), con))) {
					RT.Call("error", RT.Call("_", RT.Const("string", "\"cannot open the connection\"")))
				}
				RT.Call("strcpy", RT.Field(con, "mode"), mode)
				RT.Call("begincontext", LocalRef(&cntxt), RT.Symbol("CTXT_CCODE"), RT.Symbol("R_NilValue"), RT.Symbol("R_BaseEnv"), RT.Symbol("R_BaseEnv"), RT.Symbol("R_NilValue"), RT.Symbol("R_NilValue"))
				RT.AssignField(cntxt, "cend", RT.SymbolRef("con_cleanup"))
				RT.AssignField(cntxt, "cenddata", con)
			}
			if RT.Truth(RT.Unary("!", RT.Field(con, "canwrite"))) {
				RT.Call("error", RT.Call("_", RT.Const("string", "\"cannot write to this connection\"")))
			}
			Assign(LocalRef(&havewarned), RT.Symbol("false"))
			for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0")), Assign(LocalRef(&nout), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", i, nobjs)); RT.Inc(LocalRef(&i), 1, true) {
				if RT.Truth(RT.Binary("==", RT.Call("CAR", o), RT.Symbol("R_UnboundValue"))) {
					continue
				}
				RT.Call("SET_STRING_ELT", outnames, RT.Inc(LocalRef(&nout), 1, true), RT.Call("STRING_ELT", names, i))
				Assign(LocalRef(&s), RT.Call("translateChar", RT.Call("STRING_ELT", names, i)))
				Assign(LocalRef(&extra), RT.Const("int", "6"))
				if RT.Truth(RT.Call("isValidName", s)) {
					Assign(LocalRef(&extra), RT.Const("int", "4"))
					Assign(LocalRef(&res), RT.Call("Rconn_printf", con, RT.Const("string", "\"%s <-\\n\""), s))
				} else {
					if RT.Truth(RT.Binary("&", opts, RT.Symbol("S_COMPAT"))) {
						Assign(LocalRef(&res), RT.Call("Rconn_printf", con, RT.Const("string", "\"\\\"%s\\\" <-\\n\""), s))
					} else {
						Assign(LocalRef(&res), RT.Call("Rconn_printf", con, RT.Const("string", "\"`%s` <-\\n\""), s))
					}
				}
				if RT.Truth(func() Value {
					if !RT.Truth(RT.Unary("!", havewarned)) {
						return false
					}
					return RT.Truth(RT.Binary("<", res, RT.Binary("+", RT.Call("strlen", s), extra)))
				}()) {
					RT.Call("warning", RT.Call("_", RT.Const("string", "\"wrote too few characters\"")))
				}
				Assign(LocalRef(&tval), RT.Call("PROTECT", RT.Call("deparse1", RT.Call("CAR", o), RT.Const("int", "0"), opts)))
				for RT.Sequence(Assign(LocalRef(&j), RT.Const("int", "0"))); RT.Truth(RT.Binary("<", j, RT.Call("LENGTH", tval))); RT.Inc(LocalRef(&j), 1, true) {
					Assign(LocalRef(&res), RT.Call("Rconn_printf", con, RT.Const("string", "\"%s\\n\""), RT.Call("CHAR", RT.Call("STRING_ELT", tval, j))))
					if RT.Truth(func() Value {
						if !RT.Truth(RT.Unary("!", havewarned)) {
							return false
						}
						return RT.Truth(RT.Binary("<", res, RT.Binary("+", RT.Call("strlen", RT.Call("CHAR", RT.Call("STRING_ELT", tval, j))), RT.Const("int", "1"))))
					}()) {
						RT.Call("warning", RT.Call("_", RT.Const("string", "\"wrote too few characters\"")))
						Assign(LocalRef(&havewarned), RT.Symbol("true"))
					}
				}
				RT.Call("UNPROTECT", RT.Const("int", "1"))
				Assign(LocalRef(&o), RT.Call("CDR", o))
			}
			if RT.Truth(RT.Unary("!", wasopen)) {
				RT.Call("endcontext", LocalRef(&cntxt))
				RT.CallIndirect(RT.Field(con, "close"), con)
			}
		}
	}
	RT.Call("UNPROTECT", RT.Const("int", "2"))
	return outnames
}
