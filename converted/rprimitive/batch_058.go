package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_gzcon(call, op, args, rho Value) Value {
	var (
		ans         Value
		class       Value
		icon        Value
		level       Value
		allow       Value
		incon       Value
		new         Value
		m           Value
		mode        Value
		description Value
		text        Value
	)
	Assign(LocalRef(&incon), RT.Symbol("NULL"))
	Assign(LocalRef(&new), RT.Symbol("NULL"))
	Assign(LocalRef(&mode), RT.Symbol("NULL"))
	Assign(LocalRef(&description), RT.NewArray(RT.Const("int", "1000")))
	RT.Call("checkArity", op, args)
	if RT.Truth(RT.Unary("!", RT.Call("inherits", RT.Call("CAR", args), RT.Const("string", "\"connection\"")))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"'con' is not a connection\"")))
	}
	Assign(LocalRef(&incon), RT.Call("getConnection", Assign(LocalRef(&icon), RT.Call("asInteger", RT.Call("CAR", args)))))
	Assign(LocalRef(&level), RT.Call("asInteger", RT.Call("CADR", args)))
	if RT.Truth(func() Value {
		if RT.Truth(func() Value {
			if RT.Truth(RT.Binary("==", level, RT.Symbol("NA_INTEGER"))) {
				return true
			}
			return RT.Truth(RT.Binary("<", level, RT.Const("int", "0")))
		}()) {
			return true
		}
		return RT.Truth(RT.Binary(">", level, RT.Const("int", "9")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"'level' must be one of 0 ... 9\"")))
	}
	Assign(LocalRef(&allow), RT.Call("asLogical", RT.Call("CADDR", args)))
	if RT.Truth(RT.Binary("==", allow, RT.Symbol("NA_INTEGER"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"'allowNonCompression' must be TRUE or FALSE\"")))
	}
	Assign(LocalRef(&text), RT.Call("asLogical", RT.Call("CADDDR", args)))
	if RT.Truth(RT.Binary("==", text, RT.Symbol("NA_INTEGER"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"'text' must be TRUE or FALSE\"")))
	}
	if RT.Truth(RT.Field(incon, "isGzcon")) {
		RT.Call("warning", RT.Call("_", RT.Const("string", "\"this is already a 'gzcon' connection\"")))
		return RT.Call("CAR", args)
	}
	Assign(LocalRef(&m), RT.Field(incon, "mode"))
	if RT.Truth(func() Value {
		if RT.Truth(RT.Binary("==", RT.Call("strcmp", m, RT.Const("string", "\"r\"")), RT.Const("int", "0"))) {
			return true
		}
		return RT.Truth(RT.Binary("==", RT.Call("strncmp", m, RT.Const("string", "\"rb\""), RT.Const("int", "2")), RT.Const("int", "0")))
	}()) {
		Assign(LocalRef(&mode), RT.Const("string", "\"rb\""))
	} else {
		if RT.Truth(func() Value {
			if RT.Truth(RT.Binary("==", RT.Call("strcmp", m, RT.Const("string", "\"w\"")), RT.Const("int", "0"))) {
				return true
			}
			return RT.Truth(RT.Binary("==", RT.Call("strncmp", m, RT.Const("string", "\"wb\""), RT.Const("int", "2")), RT.Const("int", "0")))
		}()) {
			Assign(LocalRef(&mode), RT.Const("string", "\"wb\""))
		} else {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"can only use read- or write- binary connections\"")))
		}
	}
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Binary("==", RT.Call("strcmp", RT.Field(incon, "class"), RT.Const("string", "\"file\"")), RT.Const("int", "0"))) {
			return false
		}
		return RT.Truth(func() Value {
			if RT.Truth(RT.Binary("==", RT.Call("strcmp", m, RT.Const("string", "\"r\"")), RT.Const("int", "0"))) {
				return true
			}
			return RT.Truth(RT.Binary("==", RT.Call("strcmp", m, RT.Const("string", "\"w\"")), RT.Const("int", "0")))
		}())
	}()) {
		RT.Call("warning", RT.Call("_", RT.Const("string", "\"using a text-mode 'file' connection may not work correctly\"")))
	} else {
		if RT.Truth(func() Value {
			if !RT.Truth(RT.Binary("==", RT.Call("strcmp", RT.Field(incon, "class"), RT.Const("string", "\"textConnection\"")), RT.Const("int", "0"))) {
				return false
			}
			return RT.Truth(RT.Binary("==", RT.Call("strcmp", m, RT.Const("string", "\"w\"")), RT.Const("int", "0")))
		}()) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"cannot create a 'gzcon' connection from a writable textConnection; maybe use rawConnection\"")))
		}
	}
	Assign(LocalRef(&new), RT.Cast("Rconnection", RT.Call("malloc", RT.SizeOfType("struct Rconn"))))
	if RT.Truth(RT.Unary("!", new)) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"allocation of 'gzcon' connection failed\"")))
	}
	RT.AssignField(new, "class", RT.Cast("char *", RT.Call("malloc", RT.Binary("+", RT.Call("strlen", RT.Const("string", "\"gzcon\"")), RT.Const("int", "1")))))
	if RT.Truth(RT.Unary("!", RT.Field(new, "class"))) {
		RT.Call("free", new)
		RT.Call("error", RT.Call("_", RT.Const("string", "\"allocation of 'gzcon' connection failed\"")))
		Assign(LocalRef(&new), RT.Symbol("NULL"))
	}
	RT.Call("strcpy", RT.Field(new, "class"), RT.Const("string", "\"gzcon\""))
	RT.Call("Rsnprintf_mbcs", description, RT.Const("int", "1000"), RT.Const("string", "\"gzcon(%s)\""), RT.Field(incon, "description"))
	RT.AssignField(new, "description", RT.Cast("char *", RT.Call("malloc", RT.Binary("+", RT.Call("strlen", description), RT.Const("int", "1")))))
	if RT.Truth(RT.Unary("!", RT.Field(new, "description"))) {
		RT.Call("free", RT.Field(new, "class"))
		RT.Call("free", new)
		RT.Call("error", RT.Call("_", RT.Const("string", "\"allocation of 'gzcon' connection failed\"")))
		Assign(LocalRef(&new), RT.Symbol("NULL"))
	}
	RT.Call("init_con", new, description, RT.Symbol("CE_NATIVE"), mode)
	RT.AssignField(new, "text", RT.Cast("Rboolean", text))
	RT.AssignField(new, "isGzcon", RT.Symbol("TRUE"))
	RT.AssignField(new, "open", RT.SymbolRef("gzcon_open"))
	RT.AssignField(new, "close", RT.SymbolRef("gzcon_close"))
	RT.AssignField(new, "vfprintf", RT.SymbolRef("dummy_vfprintf"))
	RT.AssignField(new, "fgetc", RT.SymbolRef("gzcon_fgetc"))
	RT.AssignField(new, "read", RT.SymbolRef("gzcon_read"))
	RT.AssignField(new, "write", RT.SymbolRef("gzcon_write"))
	RT.AssignField(new, "private", RT.Cast("void *", RT.Call("malloc", RT.SizeOfType("struct gzconn"))))
	if RT.Truth(RT.Unary("!", RT.Field(new, "private"))) {
		RT.Call("free", RT.Field(new, "description"))
		RT.Call("free", RT.Field(new, "class"))
		RT.Call("free", new)
		RT.Call("error", RT.Call("_", RT.Const("string", "\"allocation of 'gzcon' connection failed\"")))
		Assign(LocalRef(&new), RT.Symbol("NULL"))
	}
	RT.AssignField(RT.Call("Rgzconn", RT.Field(new, "private")), "con", incon)
	RT.AssignField(RT.Call("Rgzconn", RT.Field(new, "private")), "cp", level)
	RT.AssignField(RT.Call("Rgzconn", RT.Field(new, "private")), "allow", RT.Cast("Rboolean", allow))
	RT.Call("R_PreserveObject", RT.Field(incon, "ex_ptr"))
	RT.AssignIndex(RT.Symbol("Connections"), icon, new)
	RT.Call("strncpy", RT.Field(new, "encname"), RT.Field(incon, "encname"), RT.Const("int", "100"))
	RT.AssignIndex(RT.Field(new, "encname"), RT.Binary("-", RT.Const("int", "100"), RT.Const("int", "1")), RT.Const("char", "'\\0'"))
	RT.AssignField(new, "ex_ptr", RT.Call("PROTECT", RT.Call("R_MakeExternalPtr", RT.Cast("void *", RT.Field(new, "id")), RT.Call("install", RT.Const("string", "\"connection\"")), RT.Symbol("R_NilValue"))))
	if RT.Truth(RT.Field(incon, "isopen")) {
		RT.CallIndirect(RT.Field(new, "open"), new)
	}
	RT.Call("PROTECT", Assign(LocalRef(&ans), RT.Call("ScalarInteger", icon)))
	RT.Call("PROTECT", Assign(LocalRef(&class), RT.Call("allocVector", RT.Symbol("STRSXP"), RT.Const("int", "2"))))
	RT.Call("SET_STRING_ELT", class, RT.Const("int", "0"), RT.Call("mkChar", RT.Const("string", "\"gzcon\"")))
	RT.Call("SET_STRING_ELT", class, RT.Const("int", "1"), RT.Call("mkChar", RT.Const("string", "\"connection\"")))
	RT.Call("classgets", ans, class)
	RT.Call("setAttrib", ans, RT.Symbol("R_ConnIdSymbol"), RT.Field(new, "ex_ptr"))
	RT.Call("UNPROTECT", RT.Const("int", "3"))
	return ans
}
