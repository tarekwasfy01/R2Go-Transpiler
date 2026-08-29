package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_save(call, op, args, env Value) Value {
	var (
		s       Value
		t       Value
		source  Value
		tmp     Value
		len     Value
		j       Value
		version Value
		ep      Value
		fp      Value
		cntxt   Value
		cfile   Value
	)
	RT.Call("checkArity", op, args)
	if RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", RT.Call("CAR", args)), RT.Symbol("STRSXP"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"first argument must be a character vector\"")))
	}
	if RT.Truth(RT.Unary("!", RT.Call("isValidStringF", RT.Call("CADR", args)))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"'file' must be non-empty string\"")))
	}
	if RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", RT.Call("CADDR", args)), RT.Symbol("LGLSXP"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"'ascii' must be logical\"")))
	}
	if RT.Truth(RT.Binary("==", RT.Call("CADDDR", args), RT.Symbol("R_NilValue"))) {
		Assign(LocalRef(&version), RT.Call("defaultSaveVersion"))
	} else {
		Assign(LocalRef(&version), RT.Call("asInteger", RT.Call("CADDDR", args)))
	}
	if RT.Truth(func() Value {
		if RT.Truth(RT.Binary("==", version, RT.Symbol("NA_INTEGER"))) {
			return true
		}
		return RT.Truth(RT.Binary("<=", version, RT.Const("int", "0")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"version\""))
	}
	Assign(LocalRef(&source), RT.Call("CAR", RT.Call("nthcdr", args, RT.Const("int", "4"))))
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Binary("!=", source, RT.Symbol("R_NilValue"))) {
			return false
		}
		return RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", source), RT.Symbol("ENVSXP")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"environment\""))
	}
	Assign(LocalRef(&ep), RT.Call("asLogical", RT.Call("CAR", RT.Call("nthcdr", args, RT.Const("int", "5")))))
	if RT.Truth(RT.Binary("==", ep, RT.Symbol("NA_LOGICAL"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"eval.promises\""))
	}
	Assign(LocalRef(&fp), RT.Call("RC_fopen", RT.Call("STRING_ELT", RT.Call("CADR", args), RT.Const("int", "0")), RT.Const("string", "\"wb\""), RT.Symbol("TRUE")))
	if RT.Truth(RT.Unary("!", fp)) {
		Assign(LocalRef(&cfile), RT.Call("CHAR", RT.Call("STRING_ELT", RT.Call("CADR", args), RT.Const("int", "0"))))
		RT.Call("error", RT.Call("_", RT.Const("string", "\"cannot open file '%s': %s\"")), cfile, RT.Call("strerror", RT.Symbol("errno")))
	}
	RT.Call("begincontext", LocalRef(&cntxt), RT.Symbol("CTXT_CCODE"), RT.Symbol("R_NilValue"), RT.Symbol("R_BaseEnv"), RT.Symbol("R_BaseEnv"), RT.Symbol("R_NilValue"), RT.Symbol("R_NilValue"))
	RT.AssignField(cntxt, "cend", RT.SymbolRef("saveload_cleanup"))
	RT.AssignField(cntxt, "cenddata", fp)
	Assign(LocalRef(&len), RT.Call("length", RT.Call("CAR", args)))
	RT.Call("PROTECT", Assign(LocalRef(&s), RT.Call("allocList", len)))
	Assign(LocalRef(&t), s)
	for Assign(LocalRef(&j), RT.Const("int", "0")); RT.Truth(RT.Binary("<", j, len)); RT.Sequence(RT.Inc(LocalRef(&j), 1, true), Assign(LocalRef(&t), RT.Call("CDR", t))) {
		RT.Call("SET_TAG", t, RT.Call("installTrChar", RT.Call("STRING_ELT", RT.Call("CAR", args), j)))
		Assign(LocalRef(&tmp), RT.Call("R_findVar", RT.Call("TAG", t), source))
		if RT.Truth(RT.Binary("==", tmp, RT.Symbol("R_UnboundValue"))) {
			RT.Call("R_ObjectNotFoundError", RT.Call("TAG", t), RT.Symbol("R_CurrentExpression"), RT.Symbol("NULL"))
		}
		if RT.Truth(func() Value {
			if !RT.Truth(ep) {
				return false
			}
			return RT.Truth(RT.Binary("==", RT.Call("TYPEOF", tmp), RT.Symbol("PROMSXP")))
		}()) {
			RT.Call("PROTECT", tmp)
			Assign(LocalRef(&tmp), RT.Call("eval", tmp, source))
			RT.Call("UNPROTECT", RT.Const("int", "1"))
		}
		RT.Call("SETCAR", t, tmp)
	}
	RT.Call("R_SaveToFileV", s, fp, RT.Index(RT.Call("INTEGER", RT.Call("CADDR", args)), RT.Const("int", "0")), version)
	RT.Call("UNPROTECT", RT.Const("int", "1"))
	RT.Call("endcontext", LocalRef(&cntxt))
	RT.Call("fclose", fp)
	return RT.Symbol("R_NilValue")
}
