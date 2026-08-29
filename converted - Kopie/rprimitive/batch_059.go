package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_gzfile(call, op, args, env Value) Value {
	var (
		sfile    Value
		sopen    Value
		ans      Value
		class    Value
		enc      Value
		file     Value
		open     Value
		ncon     Value
		compress Value
		con      Value
		type_v   Value
		subtype  Value
		ct       Value
	)
	Assign(LocalRef(&compress), RT.Const("int", "9"))
	Assign(LocalRef(&con), RT.Symbol("NULL"))
	Assign(LocalRef(&type_v), RT.Call("PRIMVAL", op))
	Assign(LocalRef(&subtype), RT.Const("int", "0"))
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&sfile), RT.Call("CAR", args))
	if RT.Truth(func() Value {
		if RT.Truth(func() Value {
			if RT.Truth(RT.Unary("!", RT.Call("isString", sfile))) {
				return true
			}
			return RT.Truth(RT.Binary("!=", RT.Call("LENGTH", sfile), RT.Const("int", "1")))
		}()) {
			return true
		}
		return RT.Truth(RT.Binary("==", RT.Call("STRING_ELT", sfile, RT.Const("int", "0")), RT.Symbol("NA_STRING")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"description\""))
	}
	if RT.Truth(RT.Binary(">", RT.Call("LENGTH", sfile), RT.Const("int", "1"))) {
		RT.Call("warning", RT.Call("_", RT.Const("string", "\"only first element of 'description' argument used\"")))
	}
	Assign(LocalRef(&file), RT.Call("translateCharFP", RT.Call("STRING_ELT", sfile, RT.Const("int", "0"))))
	Assign(LocalRef(&sopen), RT.Call("CADR", args))
	if RT.Truth(func() Value {
		if RT.Truth(RT.Unary("!", RT.Call("isString", sopen))) {
			return true
		}
		return RT.Truth(RT.Binary("!=", RT.Call("LENGTH", sopen), RT.Const("int", "1")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"open\""))
	}
	Assign(LocalRef(&enc), RT.Call("CADDR", args))
	if RT.Truth(func() Value {
		if RT.Truth(func() Value {
			if RT.Truth(RT.Unary("!", RT.Call("isString", enc))) {
				return true
			}
			return RT.Truth(RT.Binary("!=", RT.Call("LENGTH", enc), RT.Const("int", "1")))
		}()) {
			return true
		}
		return RT.Truth(RT.Binary(">", RT.Call("strlen", RT.Call("CHAR", RT.Call("STRING_ELT", enc, RT.Const("int", "0")))), RT.Const("int", "100")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"encoding\""))
	}
	if RT.Truth(RT.Binary("<", type_v, RT.Const("int", "2"))) {
		Assign(LocalRef(&compress), RT.Call("asInteger", RT.Call("CADDDR", args)))
		if RT.Truth(func() Value {
			if RT.Truth(func() Value {
				if RT.Truth(RT.Binary("==", compress, RT.Symbol("NA_LOGICAL"))) {
					return true
				}
				return RT.Truth(RT.Binary("<", compress, RT.Const("int", "0")))
			}()) {
				return true
			}
			return RT.Truth(RT.Binary(">", compress, RT.Const("int", "9")))
		}()) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"compress\""))
		}
	}
	if RT.Truth(RT.Binary("==", type_v, RT.Const("int", "2"))) {
		Assign(LocalRef(&compress), RT.Call("asInteger", RT.Call("CADDDR", args)))
		if RT.Truth(func() Value {
			if RT.Truth(RT.Binary("==", compress, RT.Symbol("NA_LOGICAL"))) {
				return true
			}
			return RT.Truth(RT.Binary(">", RT.Call("abs", compress), RT.Const("int", "9")))
		}()) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"compress\""))
		}
	}
	if RT.Truth(RT.Binary("==", type_v, RT.Const("int", "3"))) {
		Assign(LocalRef(&compress), RT.Call("asInteger", RT.Call("CADDDR", args)))
		if RT.Truth(func() Value {
			if RT.Truth(RT.Binary("==", compress, RT.Symbol("NA_LOGICAL"))) {
				return true
			}
			return RT.Truth(RT.Binary(">", RT.Call("abs", compress), RT.Const("int", "19")))
		}()) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"compress\""))
		}
	}
	Assign(LocalRef(&open), RT.Call("CHAR", RT.Call("STRING_ELT", sopen, RT.Const("int", "0"))))
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Binary("==", type_v, RT.Const("int", "0"))) {
			return false
		}
		return RT.Truth(func() Value {
			if RT.Truth(RT.Unary("!", RT.Index(open, RT.Const("int", "0")))) {
				return true
			}
			return RT.Truth(RT.Binary("==", RT.Index(open, RT.Const("int", "0")), RT.Const("char", "'r'")))
		}())
	}()) {
		Assign(LocalRef(&ct), RT.Call("comp_type_from_file", RT.Call("R_ExpandFileName", file), RT.Symbol("FALSE"), LocalRef(&subtype)))
		switch RT.Key(ct) {
		case RT.Key(RT.Symbol("COMP_GZ")), RT.Key(RT.Symbol("COMP_UNKNOWN")):
			Assign(LocalRef(&type_v), RT.Const("int", "0"))
			break
		case RT.Key(RT.Symbol("COMP_BZ")):
			Assign(LocalRef(&type_v), RT.Const("int", "1"))
			break
		case RT.Key(RT.Symbol("COMP_XZ")):
			Assign(LocalRef(&type_v), RT.Const("int", "2"))
			break
		case RT.Key(RT.Symbol("COMP_ZSTD")):
			Assign(LocalRef(&type_v), RT.Const("int", "3"))
			break
		}
	}
	switch RT.Key(type_v) {
	case RT.Key(RT.Const("int", "0")):
		Assign(LocalRef(&con), RT.Call("newgzfile", file, func() Value {
			if RT.Truth(RT.Call("strlen", open)) {
				return open
			}
			return RT.Const("string", "\"rb\"")
		}(), compress))
		break
	case RT.Key(RT.Const("int", "1")):
		Assign(LocalRef(&con), RT.Call("newbzfile", file, func() Value {
			if RT.Truth(RT.Call("strlen", open)) {
				return open
			}
			return RT.Const("string", "\"rb\"")
		}(), compress))
		break
	case RT.Key(RT.Const("int", "2")):
		Assign(LocalRef(&con), RT.Call("newxzfile", file, func() Value {
			if RT.Truth(RT.Call("strlen", open)) {
				return open
			}
			return RT.Const("string", "\"rb\"")
		}(), subtype, compress))
		break
	case RT.Key(RT.Const("int", "3")):
		Assign(LocalRef(&con), RT.Call("newzstdfile", file, func() Value {
			if RT.Truth(RT.Call("strlen", open)) {
				return open
			}
			return RT.Const("string", "\"rb\"")
		}(), compress))
		break
	}
	Assign(LocalRef(&ncon), RT.Call("NextConnection"))
	RT.AssignIndex(RT.Symbol("Connections"), ncon, con)
	RT.AssignField(con, "blocking", RT.Symbol("TRUE"))
	RT.Call("strncpy", RT.Field(con, "encname"), RT.Call("CHAR", RT.Call("STRING_ELT", enc, RT.Const("int", "0"))), RT.Const("int", "100"))
	RT.AssignIndex(RT.Field(con, "encname"), RT.Binary("-", RT.Const("int", "100"), RT.Const("int", "1")), RT.Const("char", "'\\0'"))
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Index(RT.Field(con, "encname"), RT.Const("int", "0"))) {
			return false
		}
		return RT.Truth(RT.Unary("!", RT.Call("streql", RT.Field(con, "encname"), RT.Const("string", "\"native.enc\""))))
	}()) {
		RT.AssignField(con, "canseek", RT.Const("int", "0"))
	}
	RT.AssignField(con, "ex_ptr", RT.Call("PROTECT", RT.Call("R_MakeExternalPtr", RT.Field(con, "id"), RT.Call("install", RT.Const("string", "\"connection\"")), RT.Symbol("R_NilValue"))))
	if RT.Truth(RT.Call("strlen", open)) {
		RT.Call("checked_open", ncon)
	}
	RT.Call("PROTECT", Assign(LocalRef(&ans), RT.Call("ScalarInteger", ncon)))
	RT.Call("PROTECT", Assign(LocalRef(&class), RT.Call("allocVector", RT.Symbol("STRSXP"), RT.Const("int", "2"))))
	switch RT.Key(type_v) {
	case RT.Key(RT.Const("int", "0")):
		RT.Call("SET_STRING_ELT", class, RT.Const("int", "0"), RT.Call("mkChar", RT.Const("string", "\"gzfile\"")))
		break
	case RT.Key(RT.Const("int", "1")):
		RT.Call("SET_STRING_ELT", class, RT.Const("int", "0"), RT.Call("mkChar", RT.Const("string", "\"bzfile\"")))
		break
	case RT.Key(RT.Const("int", "2")):
		RT.Call("SET_STRING_ELT", class, RT.Const("int", "0"), RT.Call("mkChar", RT.Const("string", "\"xzfile\"")))
		break
	case RT.Key(RT.Const("int", "3")):
		RT.Call("SET_STRING_ELT", class, RT.Const("int", "0"), RT.Call("mkChar", RT.Const("string", "\"zstdfile\"")))
		break
	}
	RT.Call("SET_STRING_ELT", class, RT.Const("int", "1"), RT.Call("mkChar", RT.Const("string", "\"connection\"")))
	RT.Call("classgets", ans, class)
	RT.Call("setAttrib", ans, RT.Symbol("R_ConnIdSymbol"), RT.Field(con, "ex_ptr"))
	RT.Call("R_RegisterCFinalizerEx", RT.Field(con, "ex_ptr"), RT.Symbol("conFinalizer"), RT.Symbol("FALSE"))
	RT.Call("UNPROTECT", RT.Const("int", "3"))
	return ans
}
