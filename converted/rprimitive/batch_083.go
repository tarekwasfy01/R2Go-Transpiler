package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_mmap_file(call, op, args, env Value) Value {
	var (
		file    Value
		stype   Value
		sptrOK  Value
		swrtOK  Value
		sserOK  Value
		type_v  Value
		typestr Value
		ptrOK   Value
		wrtOK   Value
		serOK   Value
	)
	Assign(LocalRef(&file), RT.Call("CAR", args))
	Assign(LocalRef(&stype), RT.Call("CADR", args))
	Assign(LocalRef(&sptrOK), RT.Call("CADDR", args))
	Assign(LocalRef(&swrtOK), RT.Call("CADDDR", args))
	Assign(LocalRef(&sserOK), RT.Call("CAD4R", args))
	Assign(LocalRef(&type_v), RT.Symbol("REALSXP"))
	if RT.Truth(RT.Binary("!=", stype, RT.Symbol("R_NilValue"))) {
		Assign(LocalRef(&typestr), RT.Call("CHAR", RT.Call("asChar", stype)))
		if RT.Truth(RT.Binary("==", RT.Call("strcmp", typestr, RT.Const("string", "\"double\"")), RT.Const("int", "0"))) {
			Assign(LocalRef(&type_v), RT.Symbol("REALSXP"))
		} else {
			if RT.Truth(func() Value {
				if RT.Truth(RT.Binary("==", RT.Call("strcmp", typestr, RT.Const("string", "\"integer\"")), RT.Const("int", "0"))) {
					return true
				}
				return RT.Truth(RT.Binary("==", RT.Call("strcmp", typestr, RT.Const("string", "\"int\"")), RT.Const("int", "0")))
			}()) {
				Assign(LocalRef(&type_v), RT.Symbol("INTSXP"))
			} else {
				RT.Call("error", RT.Const("string", "\"type '%s' is not supported\""), typestr)
			}
		}
	}
	Assign(LocalRef(&ptrOK), func() Value {
		if RT.Truth(RT.Binary("==", sptrOK, RT.Symbol("R_NilValue"))) {
			return RT.Symbol("TRUE")
		}
		return RT.Call("asLogicalNA", sptrOK, RT.Symbol("FALSE"))
	}())
	Assign(LocalRef(&wrtOK), func() Value {
		if RT.Truth(RT.Binary("==", swrtOK, RT.Symbol("R_NilValue"))) {
			return RT.Symbol("FALSE")
		}
		return RT.Call("asLogicalNA", swrtOK, RT.Symbol("FALSE"))
	}())
	Assign(LocalRef(&serOK), func() Value {
		if RT.Truth(RT.Binary("==", sserOK, RT.Symbol("R_NilValue"))) {
			return RT.Symbol("FALSE")
		}
		return RT.Call("asLogicalNA", sserOK, RT.Symbol("FALSE"))
	}())
	if RT.Truth(func() Value {
		if RT.Truth(func() Value {
			if RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", file), RT.Symbol("STRSXP"))) {
				return true
			}
			return RT.Truth(RT.Binary("!=", RT.Call("LENGTH", file), RT.Const("int", "1")))
		}()) {
			return true
		}
		return RT.Truth(RT.Binary("==", file, RT.Symbol("NA_STRING")))
	}()) {
		RT.Call("error", RT.Const("string", "\"invalud 'file' argument\""))
	}
	return RT.Call("mmap_file", file, type_v, ptrOK, wrtOK, serOK, RT.Symbol("FALSE"))
}
