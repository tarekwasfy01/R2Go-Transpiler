package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_dynload(call, op, args, env Value) Value {
	var (
		buf  Value
		info Value
	)
	Assign(LocalRef(&buf), RT.NewArray(RT.Binary("*", RT.Const("int", "2"), RT.Symbol("R_PATH_MAX"))))
	RT.Call("checkArity", op, args)
	if RT.Truth(func() Value {
		if RT.Truth(RT.Unary("!", RT.Call("isString", RT.Call("CAR", args)))) {
			return true
		}
		return RT.Truth(RT.Binary("!=", RT.Call("LENGTH", RT.Call("CAR", args)), RT.Const("int", "1")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"character argument expected\"")))
	}
	RT.Call("GetFullDLLPath", call, buf, RT.SizeOf(buf), RT.Call("translateCharFP", RT.Call("STRING_ELT", RT.Call("CAR", args), RT.Const("int", "0"))))
	Assign(LocalRef(&info), RT.Call("AddDLL", buf, RT.Index(RT.Call("LOGICAL", RT.Call("CADR", args)), RT.Const("int", "0")), RT.Index(RT.Call("LOGICAL", RT.Call("CADDR", args)), RT.Const("int", "0")), RT.Call("translateCharFP", RT.Call("STRING_ELT", RT.Call("CADDDR", args), RT.Const("int", "0")))))
	if RT.Truth(RT.Unary("!", info)) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"unable to load shared object '%s':\\n  %s\"")), buf, RT.Symbol("DLLerror"))
	}
	return RT.Call("Rf_MakeDLLInfo", info)
}
