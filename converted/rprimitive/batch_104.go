package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_shellexec(call, op, args, env Value) Value {
	var (
		file Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&file), RT.Call("CAR", args))
	if RT.Truth(func() Value {
		if RT.Truth(RT.Unary("!", RT.Call("isString", file))) {
			return true
		}
		return RT.Truth(RT.Binary("!=", RT.Call("length", file), RT.Const("int", "1")))
	}()) {
		RT.Call("errorcall", call, RT.Call("_", RT.Const("string", "\"invalid '%s' argument\"")), RT.Const("string", "\"file\""))
	}
	RT.Call("internal_shellexecW", RT.Call("filenameToWchar", RT.Call("STRING_ELT", file, RT.Const("int", "0")), RT.Symbol("FALSE")), RT.Symbol("FALSE"))
	return RT.Symbol("R_NilValue")
}
