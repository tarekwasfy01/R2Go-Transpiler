package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_scan(call, op, args, rho Value) Value {
	var (
		ans        Value
		file       Value
		sep        Value
		what       Value
		stripwhite Value
		dec        Value
		quotes     Value
		comstr     Value
		c          Value
		flush      Value
		fill       Value
		blskip     Value
		multiline  Value
		escapes    Value
		skipNul    Value
		nmax       Value
		nlines     Value
		nskip      Value
		p          Value
		encoding   Value
		cntxt      Value
		data       Value
		sc         Value
		dc         Value
		ii         Value
		i          Value
		j          Value
		line       Value
	)
	Assign(LocalRef(&data), RT.List(RT.Symbol("NULL"), RT.Const("int", "0"), RT.Const("int", "0"), RT.Const("char", "'.'"), RT.Symbol("NULL"), RT.Symbol("NO_COMCHAR"), RT.Const("int", "0"), RT.Symbol("NULL"), RT.Symbol("false"), RT.Symbol("false"), RT.Const("int", "0"), RT.Symbol("false"), RT.Symbol("false"), RT.Symbol("false"), RT.Symbol("false"), RT.Symbol("false"), RT.List(RT.Symbol("false"))))
	RT.AssignField(data, "NAstrings", RT.Symbol("R_NilValue"))
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&file), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&what), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&nmax), RT.Call("asXLength", RT.Call("CAR", args)))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&sep), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&dec), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&quotes), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&nskip), RT.Call("asXLength", RT.Call("CAR", args)))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&nlines), RT.Call("asXLength", RT.Call("CAR", args)))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	RT.AssignField(data, "NAstrings", RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&flush), RT.Call("asLogical", RT.Call("CAR", args)))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&fill), RT.Call("asLogical", RT.Call("CAR", args)))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&stripwhite), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	RT.AssignField(data, "quiet", RT.Call("asLogical", RT.Call("CAR", args)))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&blskip), RT.Call("asLogical", RT.Call("CAR", args)))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&multiline), RT.Call("asLogical", RT.Call("CAR", args)))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&comstr), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&escapes), RT.Call("asLogical", RT.Call("CAR", args)))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	if RT.Truth(func() Value {
		if RT.Truth(RT.Unary("!", RT.Call("isString", RT.Call("CAR", args)))) {
			return true
		}
		return RT.Truth(RT.Binary("!=", RT.Call("LENGTH", RT.Call("CAR", args)), RT.Const("int", "1")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"encoding\""))
	}
	Assign(LocalRef(&encoding), RT.Call("CHAR", RT.Call("STRING_ELT", RT.Call("CAR", args), RT.Const("int", "0"))))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	if RT.Truth(RT.Call("streql", encoding, RT.Const("string", "\"latin1\""))) {
		RT.AssignField(data, "isLatin1", RT.Symbol("true"))
	}
	if RT.Truth(RT.Call("streql", encoding, RT.Const("string", "\"UTF-8\""))) {
		RT.AssignField(data, "isUTF8", RT.Symbol("true"))
	}
	Assign(LocalRef(&skipNul), RT.Call("asLogical", RT.Call("CAR", args)))
	if RT.Truth(RT.Binary("==", RT.Field(data, "quiet"), RT.Symbol("NA_LOGICAL"))) {
		RT.AssignField(data, "quiet", RT.Const("int", "0"))
	}
	if RT.Truth(RT.Binary("==", blskip, RT.Symbol("NA_LOGICAL"))) {
		Assign(LocalRef(&blskip), RT.Const("int", "1"))
	}
	if RT.Truth(RT.Binary("==", multiline, RT.Symbol("NA_LOGICAL"))) {
		Assign(LocalRef(&multiline), RT.Const("int", "1"))
	}
	if RT.Truth(RT.Binary("<", nskip, RT.Const("int", "0"))) {
		Assign(LocalRef(&nskip), RT.Const("int", "0"))
	}
	if RT.Truth(RT.Binary("<", nlines, RT.Const("int", "0"))) {
		Assign(LocalRef(&nlines), RT.Const("int", "0"))
	}
	if RT.Truth(RT.Binary("<", nmax, RT.Const("int", "0"))) {
		Assign(LocalRef(&nmax), RT.Const("int", "0"))
	}
	if RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", stripwhite), RT.Symbol("LGLSXP"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"strip.white\""))
	}
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Binary("!=", RT.Call("xlength", stripwhite), RT.Const("int", "1"))) {
			return false
		}
		return RT.Truth(RT.Binary("!=", RT.Call("xlength", stripwhite), RT.Call("xlength", what)))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid 'strip.white' length\"")))
	}
	if RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", RT.Field(data, "NAstrings")), RT.Symbol("STRSXP"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"na.strings\""))
	}
	if RT.Truth(func() Value {
		if RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", comstr), RT.Symbol("STRSXP"))) {
			return true
		}
		return RT.Truth(RT.Binary("!=", RT.Call("length", comstr), RT.Const("int", "1")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"comment.char\""))
	}
	if RT.Truth(func() Value {
		if RT.Truth(RT.Call("isString", sep)) {
			return true
		}
		return RT.Truth(RT.Call("isNull", sep))
	}()) {
		if RT.Truth(RT.Binary("==", RT.Call("length", sep), RT.Const("int", "0"))) {
			RT.AssignField(data, "sepchar", RT.Const("int", "0"))
		} else {
			Assign(LocalRef(&sc), RT.Call("translateChar", RT.Call("STRING_ELT", sep, RT.Const("int", "0"))))
			if RT.Truth(RT.Binary(">", RT.Call("strlen", sc), RT.Const("int", "1"))) {
				RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid 'sep' value: must be one byte\"")))
			}
			RT.AssignField(data, "sepchar", RT.Cast("unsigned char", RT.Index(sc, RT.Const("int", "0"))))
		}
	} else {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"sep\""))
	}
	if RT.Truth(func() Value {
		if RT.Truth(RT.Call("isString", dec)) {
			return true
		}
		return RT.Truth(RT.Call("isNull", dec))
	}()) {
		if RT.Truth(RT.Binary("==", RT.Call("length", dec), RT.Const("int", "0"))) {
			RT.AssignField(data, "decchar", RT.Const("char", "'.'"))
		} else {
			Assign(LocalRef(&dc), RT.Call("translateChar", RT.Call("STRING_ELT", dec, RT.Const("int", "0"))))
			if RT.Truth(RT.Binary("!=", RT.Call("strlen", dc), RT.Const("int", "1"))) {
				RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid decimal separator: must be one byte\"")))
			}
			RT.AssignField(data, "decchar", RT.Index(dc, RT.Const("int", "0")))
		}
	} else {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid decimal separator\"")))
	}
	RT.Call("begincontext", LocalRef(&cntxt), RT.Symbol("CTXT_CCODE"), RT.Field(RT.Symbol("R_GlobalContext"), "call"), RT.Symbol("R_BaseEnv"), RT.Symbol("R_BaseEnv"), RT.Symbol("R_NilValue"), RT.Symbol("R_NilValue"))
	RT.AssignField(cntxt, "cend", RT.SymbolRef("scan_cleanup"))
	RT.AssignField(cntxt, "cenddata", LocalRef(&data))
	if RT.Truth(RT.Call("isString", quotes)) {
		Assign(LocalRef(&sc), RT.Call("translateChar", RT.Call("STRING_ELT", quotes, RT.Const("int", "0"))))
		if RT.Truth(RT.Call("strlen", sc)) {
			RT.AssignField(data, "quoteset", RT.Call("Rstrdup", sc))
		} else {
			RT.AssignField(data, "quoteset", RT.Const("string", "\"\""))
		}
	} else {
		if RT.Truth(RT.Call("isNull", quotes)) {
			RT.AssignField(data, "quoteset", RT.Const("string", "\"\""))
		} else {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid quote symbol set\"")))
		}
	}
	Assign(LocalRef(&p), RT.Call("translateChar", RT.Call("STRING_ELT", comstr, RT.Const("int", "0"))))
	RT.AssignField(data, "comchar", RT.Symbol("NO_COMCHAR"))
	if RT.Truth(RT.Binary(">", RT.Call("strlen", p), RT.Const("int", "1"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"comment.char\""))
	} else {
		if RT.Truth(RT.Binary("==", RT.Call("strlen", p), RT.Const("int", "1"))) {
			RT.AssignField(data, "comchar", RT.Cast("unsigned char", RT.Deref(p)))
		}
	}
	if RT.Truth(RT.Binary("==", escapes, RT.Symbol("NA_LOGICAL"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"allowEscapes\""))
	}
	RT.AssignField(data, "escapes", RT.Binary("!=", escapes, RT.Const("int", "0")))
	if RT.Truth(RT.Binary("==", skipNul, RT.Symbol("NA_LOGICAL"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"skipNul\""))
	}
	RT.AssignField(data, "skipNul", RT.Binary("!=", skipNul, RT.Const("int", "0")))
	Assign(LocalRef(&ii), RT.Call("asInteger", file))
	RT.AssignField(data, "con", RT.Call("getConnection", ii))
	if RT.Truth(RT.Binary("==", ii, RT.Const("int", "0"))) {
		RT.AssignField(data, "atStart", RT.Symbol("false"))
		RT.AssignField(data, "ttyflag", RT.Const("int", "1"))
	} else {
		RT.AssignField(data, "atStart", RT.Binary("==", nskip, RT.Const("int", "0")))
		RT.AssignField(data, "ttyflag", RT.Const("int", "0"))
		RT.AssignField(data, "wasopen", RT.Field(RT.Field(data, "con"), "isopen"))
		if RT.Truth(RT.Unary("!", RT.Field(data, "wasopen"))) {
			RT.AssignField(RT.Field(data, "con"), "UTF8out", RT.Symbol("true"))
			RT.Call("strcpy", RT.Field(RT.Field(data, "con"), "mode"), RT.Const("string", "\"r\""))
			if RT.Truth(RT.Unary("!", RT.CallIndirect(RT.Field(RT.Field(data, "con"), "open"), RT.Field(data, "con")))) {
				RT.Call("error", RT.Call("_", RT.Const("string", "\"cannot open the connection\"")))
			}
			if RT.Truth(RT.Unary("!", RT.Field(RT.Field(data, "con"), "canread"))) {
				RT.CallIndirect(RT.Field(RT.Field(data, "con"), "close"), RT.Field(data, "con"))
				RT.Call("error", RT.Call("_", RT.Const("string", "\"cannot read from this connection\"")))
			}
		} else {
			if RT.Truth(RT.Unary("!", RT.Field(RT.Field(data, "con"), "canread"))) {
				RT.Call("error", RT.Call("_", RT.Const("string", "\"cannot read from this connection\"")))
			}
		}
		for RT.Sequence(Assign(LocalRef(&i), RT.Const("int", "0")), Assign(LocalRef(&j), RT.Const("int", "10000"))); RT.Truth(RT.Binary("<", i, nskip)); RT.Inc(LocalRef(&i), 1, true) {
			for RT.Truth(RT.Const("int", "1")) {
				Assign(LocalRef(&c), RT.Call("scanchar", RT.Symbol("false"), LocalRef(&data)))
				if RT.Truth(RT.Unary("!", RT.Inc(LocalRef(&j), -1, true))) {
					RT.Call("R_CheckUserInterrupt")
					Assign(LocalRef(&j), RT.Const("int", "10000"))
				}
				if RT.Truth(func() Value {
					if RT.Truth(RT.Binary("==", c, RT.Const("char", "'\\n'"))) {
						return true
					}
					return RT.Truth(RT.Binary("==", c, RT.Symbol("R_EOF")))
				}()) {
					break
				}
			}
		}
	}
	Assign(LocalRef(&ans), RT.Symbol("R_NilValue"))
	RT.AssignField(data, "save", RT.Const("int", "0"))
	switch RT.Key(RT.Call("TYPEOF", what)) {
	case RT.Key(RT.Symbol("LGLSXP")), RT.Key(RT.Symbol("INTSXP")), RT.Key(RT.Symbol("REALSXP")), RT.Key(RT.Symbol("CPLXSXP")), RT.Key(RT.Symbol("STRSXP")), RT.Key(RT.Symbol("RAWSXP")):
		Assign(LocalRef(&ans), RT.Call("scanVector", RT.Call("TYPEOF", what), nmax, nlines, flush, stripwhite, blskip, LocalRef(&data)))
		break
	case RT.Key(RT.Symbol("VECSXP")):
		Assign(LocalRef(&ans), RT.Call("scanFrame", what, nmax, nlines, flush, fill, stripwhite, blskip, multiline, LocalRef(&data)))
		break
	default:
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"what\""))
	}
	RT.Call("PROTECT", ans)
	RT.Call("endcontext", LocalRef(&cntxt))
	if RT.Truth(func() Value {
		if !RT.Truth(func() Value {
			if !RT.Truth(RT.Field(data, "save")) {
				return false
			}
			return RT.Truth(RT.Unary("!", RT.Field(data, "ttyflag")))
		}()) {
			return false
		}
		return RT.Truth(RT.Field(data, "wasopen"))
	}()) {
		Assign(LocalRef(&line), RT.Const("string", "\" \""))
		RT.AssignIndex(line, RT.Const("int", "0"), RT.Cast("char", RT.Field(data, "save")))
		RT.Call("con_pushback", RT.Field(data, "con"), RT.Symbol("false"), line)
	}
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Unary("!", RT.Field(data, "ttyflag"))) {
			return false
		}
		return RT.Truth(RT.Unary("!", RT.Field(data, "wasopen")))
	}()) {
		RT.CallIndirect(RT.Field(RT.Field(data, "con"), "close"), RT.Field(data, "con"))
	}
	if RT.Truth(RT.Index(RT.Field(data, "quoteset"), RT.Const("int", "0"))) {
		RT.Call("free", RT.Field(data, "quoteset"))
	}
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Unary("!", skipNul)) {
			return false
		}
		return RT.Truth(RT.Field(data, "embedWarn"))
	}()) {
		RT.Call("warning", RT.Call("_", RT.Const("string", "\"embedded nul(s) found in input\"")))
	}
	RT.Call("UNPROTECT", RT.Const("int", "1"))
	return ans
}
