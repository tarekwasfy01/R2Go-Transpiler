package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_lazyLoadDBfetch(call, op, args, env Value) Value {
	var (
		key        Value
		file       Value
		compsxp    Value
		hook       Value
		vpi        Value
		compressed Value
		err        Value
		val        Value
	)
	Assign(LocalRef(&err), RT.Symbol("FALSE"))
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&key), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&file), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&compsxp), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&hook), RT.Call("CAR", args))
	Assign(LocalRef(&compressed), RT.Call("asInteger", compsxp))
	RT.Call("PROTECT_WITH_INDEX", Assign(LocalRef(&val), RT.Call("readRawFromFile", file, key)), LocalRef(&vpi))
	if RT.Truth(RT.Binary("==", compressed, RT.Const("int", "3"))) {
		RT.Call("REPROTECT", Assign(LocalRef(&val), RT.Call("R_decompress3", val, LocalRef(&err))), vpi)
	} else {
		if RT.Truth(RT.Binary("==", compressed, RT.Const("int", "2"))) {
			RT.Call("REPROTECT", Assign(LocalRef(&val), RT.Call("R_decompress2", val, LocalRef(&err))), vpi)
		} else {
			if RT.Truth(compressed) {
				RT.Call("REPROTECT", Assign(LocalRef(&val), RT.Call("R_decompress1", val, LocalRef(&err))), vpi)
			}
		}
	}
	if RT.Truth(err) {
		RT.Call("error", RT.Const("string", "\"lazy-load database '%s' is corrupt\""), RT.Call("translateChar", RT.Call("STRING_ELT", file, RT.Const("int", "0"))))
	}
	Assign(LocalRef(&val), RT.Call("R_unserialize", val, hook))
	if RT.Truth(RT.Binary("==", RT.Call("TYPEOF", val), RT.Symbol("PROMSXP"))) {
		RT.Call("REPROTECT", val, vpi)
		Assign(LocalRef(&val), RT.Call("eval", val, RT.Symbol("R_GlobalEnv")))
		RT.Call("ENSURE_NAMEDMAX", val)
	}
	RT.Call("UNPROTECT", RT.Const("int", "1"))
	return val
}
