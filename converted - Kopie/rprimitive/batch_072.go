package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_loadFromConn2(call, op, args, env Value) Value {
	var (
		in                    Value
		con                   Value
		aenv                  Value
		res                   Value
		buf                   Value
		count                 Value
		wasopen               Value
		cntxt                 Value
		mode                  Value
		old_InitReadItemDepth Value
		old_ReadItemDepth     Value
	)
	Assign(LocalRef(&in), RT.NewObject())
	Assign(LocalRef(&aenv), RT.Symbol("R_NilValue"))
	Assign(LocalRef(&res), RT.Symbol("R_NilValue"))
	Assign(LocalRef(&buf), RT.NewArray(RT.Const("int", "6")))
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&con), RT.Call("getConnection", RT.Call("asInteger", RT.Call("CAR", args))))
	Assign(LocalRef(&wasopen), RT.Field(con, "isopen"))
	if RT.Truth(RT.Unary("!", wasopen)) {
		Assign(LocalRef(&mode), RT.NewArray(RT.Const("int", "5")))
		RT.Call("strcpy", mode, RT.Field(con, "mode"))
		RT.Call("strcpy", RT.Field(con, "mode"), RT.Const("string", "\"rb\""))
		if RT.Truth(RT.Unary("!", RT.CallIndirect(RT.Field(con, "open"), con))) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"cannot open the connection\"")))
		}
		RT.Call("strcpy", RT.Field(con, "mode"), mode)
		RT.Call("begincontext", LocalRef(&cntxt), RT.Symbol("CTXT_CCODE"), RT.Symbol("R_NilValue"), RT.Symbol("R_BaseEnv"), RT.Symbol("R_BaseEnv"), RT.Symbol("R_NilValue"), RT.Symbol("R_NilValue"))
		RT.AssignField(cntxt, "cend", RT.SymbolRef("con_cleanup"))
		RT.AssignField(cntxt, "cenddata", con)
	}
	if RT.Truth(RT.Unary("!", RT.Field(con, "canread"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"connection not open for reading\"")))
	}
	if RT.Truth(RT.Field(con, "text")) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"can only load() from a binary connection\"")))
	}
	if RT.Truth(RT.Binary("==", RT.Call("PRIMVAL", op), RT.Const("int", "0"))) {
		Assign(LocalRef(&aenv), RT.Call("CADR", args))
		if RT.Truth(RT.Binary("==", RT.Call("TYPEOF", aenv), RT.Symbol("NILSXP"))) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"use of NULL environment is defunct\"")))
		} else {
			if RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", aenv), RT.Symbol("ENVSXP"))) {
				RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"envir\""))
			}
		}
	}
	RT.Call("memset", buf, RT.Const("int", "0"), RT.Const("int", "6"))
	Assign(LocalRef(&count), RT.CallIndirect(RT.Field(con, "read"), buf, RT.SizeOfType("char"), RT.Const("int", "5"), con))
	if RT.Truth(RT.Binary("==", count, RT.Const("int", "0"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"no input is available\"")))
	}
	if RT.Truth(func() Value {
		if RT.Truth(func() Value {
			if RT.Truth(func() Value {
				if RT.Truth(func() Value {
					if RT.Truth(func() Value {
						if RT.Truth(RT.Binary("==", RT.Call("strncmp", RT.Cast("char *", buf), RT.Const("string", "\"RDA2\\n\""), RT.Const("int", "5")), RT.Const("int", "0"))) {
							return true
						}
						return RT.Truth(RT.Binary("==", RT.Call("strncmp", RT.Cast("char *", buf), RT.Const("string", "\"RDB2\\n\""), RT.Const("int", "5")), RT.Const("int", "0")))
					}()) {
						return true
					}
					return RT.Truth(RT.Binary("==", RT.Call("strncmp", RT.Cast("char *", buf), RT.Const("string", "\"RDX2\\n\""), RT.Const("int", "5")), RT.Const("int", "0")))
				}()) {
					return true
				}
				return RT.Truth(RT.Binary("==", RT.Call("strncmp", RT.Cast("char *", buf), RT.Const("string", "\"RDA3\\n\""), RT.Const("int", "5")), RT.Const("int", "0")))
			}()) {
				return true
			}
			return RT.Truth(RT.Binary("==", RT.Call("strncmp", RT.Cast("char *", buf), RT.Const("string", "\"RDB3\\n\""), RT.Const("int", "5")), RT.Const("int", "0")))
		}()) {
			return true
		}
		return RT.Truth(RT.Binary("==", RT.Call("strncmp", RT.Cast("char *", buf), RT.Const("string", "\"RDX3\\n\""), RT.Const("int", "5")), RT.Const("int", "0")))
	}()) {
		RT.Call("R_InitConnInPStream", LocalRef(&in), con, RT.Symbol("R_pstream_any_format"), RT.Symbol("NULL"), RT.Symbol("NULL"))
		if RT.Truth(RT.Binary("==", RT.Call("PRIMVAL", op), RT.Const("int", "0"))) {
			Assign(LocalRef(&old_InitReadItemDepth), RT.Symbol("R_InitReadItemDepth"))
			Assign(LocalRef(&old_ReadItemDepth), RT.Symbol("R_ReadItemDepth"))
			RT.AssignSymbol("R_InitReadItemDepth", RT.AssignSymbol("R_ReadItemDepth", RT.Unary("-", RT.Call("asInteger", RT.Call("CADDR", args)))))
			Assign(LocalRef(&res), RT.Call("RestoreToEnv", RT.Call("R_Unserialize", LocalRef(&in)), aenv))
			RT.AssignSymbol("R_InitReadItemDepth", old_InitReadItemDepth)
			RT.AssignSymbol("R_ReadItemDepth", old_ReadItemDepth)
		} else {
			Assign(LocalRef(&res), RT.Call("R_SerializeInfo", LocalRef(&in)))
		}
		if RT.Truth(RT.Unary("!", wasopen)) {
			RT.Call("PROTECT", res)
			RT.Call("endcontext", LocalRef(&cntxt))
			RT.CallIndirect(RT.Field(con, "close"), con)
			RT.Call("UNPROTECT", RT.Const("int", "1"))
		}
	} else {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"the input does not start with a magic number compatible with loading from a connection\"")))
	}
	return res
}
