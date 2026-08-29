package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_isloaded(call, op, args, env Value) Value {
	var (
		sym    Value
		type_v Value
		pkg    Value
		val    Value
		nargs  Value
		symbol Value
	)
	Assign(LocalRef(&type_v), RT.Const("string", "\"\""))
	Assign(LocalRef(&pkg), RT.Const("string", "\"\""))
	Assign(LocalRef(&val), RT.Const("int", "1"))
	Assign(LocalRef(&nargs), RT.Call("length", args))
	Assign(LocalRef(&symbol), RT.List(RT.Symbol("R_ANY_SYM"), RT.List(RT.Symbol("NULL")), RT.Symbol("NULL")))
	if RT.Truth(RT.Binary("<", nargs, RT.Const("int", "1"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"no arguments supplied\"")))
	}
	if RT.Truth(RT.Binary(">", nargs, RT.Const("int", "3"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"too many arguments\"")))
	}
	if RT.Truth(RT.Unary("!", RT.Call("isValidString", RT.Call("CAR", args)))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"symbol\""))
	}
	Assign(LocalRef(&sym), RT.Call("translateChar", RT.Call("STRING_ELT", RT.Call("CAR", args), RT.Const("int", "0"))))
	if RT.Truth(RT.Binary(">=", nargs, RT.Const("int", "2"))) {
		if RT.Truth(RT.Unary("!", RT.Call("isValidString", RT.Call("CADR", args)))) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"PACKAGE\""))
		}
		Assign(LocalRef(&pkg), RT.Call("translateChar", RT.Call("STRING_ELT", RT.Call("CADR", args), RT.Const("int", "0"))))
	}
	if RT.Truth(RT.Binary(">=", nargs, RT.Const("int", "3"))) {
		if RT.Truth(RT.Unary("!", RT.Call("isValidString", RT.Call("CADDR", args)))) {
			RT.Call("error", RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"type\""))
		}
		Assign(LocalRef(&type_v), RT.Call("CHAR", RT.Call("STRING_ELT", RT.Call("CADDR", args), RT.Const("int", "0"))))
		if RT.Truth(RT.Binary("==", RT.Call("strcmp", type_v, RT.Const("string", "\"C\"")), RT.Const("int", "0"))) {
			RT.AssignField(symbol, "type", RT.Symbol("R_C_SYM"))
		} else {
			if RT.Truth(RT.Binary("==", RT.Call("strcmp", type_v, RT.Const("string", "\"Fortran\"")), RT.Const("int", "0"))) {
				RT.AssignField(symbol, "type", RT.Symbol("R_FORTRAN_SYM"))
			} else {
				if RT.Truth(RT.Binary("==", RT.Call("strcmp", type_v, RT.Const("string", "\"Call\"")), RT.Const("int", "0"))) {
					RT.AssignField(symbol, "type", RT.Symbol("R_CALL_SYM"))
				} else {
					if RT.Truth(RT.Binary("==", RT.Call("strcmp", type_v, RT.Const("string", "\"External\"")), RT.Const("int", "0"))) {
						RT.AssignField(symbol, "type", RT.Symbol("R_EXTERNAL_SYM"))
					}
				}
			}
		}
	}
	if RT.Truth(RT.Unary("!", RT.Call("R_FindSymbol", sym, pkg, LocalRef(&symbol)))) {
		Assign(LocalRef(&val), RT.Const("int", "0"))
	}
	return RT.Call("ScalarLogical", val)
}
