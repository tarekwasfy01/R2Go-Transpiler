package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_dynunload(call, op, args, env Value) Value {
	var (
		buf Value
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
	if RT.Truth(RT.Unary("!", RT.Call("DeleteDLL", buf))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"shared object '%s\\' was not loaded\"")), buf)
	}
	return RT.Symbol("R_NilValue")
}
