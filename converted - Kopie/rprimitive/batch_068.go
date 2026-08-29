package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_lazyLoadDBflush(call, op, args, env Value) Value {
	var (
		i     Value
		cfile Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&cfile), RT.Call("translateCharFP", RT.Call("STRING_ELT", RT.Call("CAR", args), RT.Const("int", "0"))))
	for Assign(LocalRef(&i), RT.Const("int", "0")); RT.Truth(RT.Binary("<", i, RT.Symbol("used"))); RT.Inc(LocalRef(&i), 1, true) {
		if RT.Truth(func() Value {
			if !RT.Truth(RT.Binary("!=", RT.Index(RT.Symbol("names"), i), RT.Symbol("NULL"))) {
				return false
			}
			return RT.Truth(RT.Binary("==", RT.Call("strcmp", cfile, RT.Index(RT.Symbol("names"), i)), RT.Const("int", "0")))
		}()) {
			RT.Call("free", RT.Index(RT.Symbol("names"), i))
			RT.AssignIndex(RT.Symbol("names"), i, RT.Symbol("NULL"))
			RT.Call("free", RT.Index(RT.Symbol("ptr"), i))
			break
		}
	}
	return RT.Symbol("R_NilValue")
}
