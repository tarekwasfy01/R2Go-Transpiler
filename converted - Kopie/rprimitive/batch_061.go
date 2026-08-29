package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_inspect(call, op, args, env Value) Value {
	var (
		obj  Value
		deep Value
		pvec Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&obj), RT.Call("CAR", args))
	Assign(LocalRef(&deep), RT.Unary("-", RT.Const("int", "1")))
	Assign(LocalRef(&pvec), RT.Const("int", "5"))
	if RT.Truth(RT.Binary("!=", RT.Call("CDR", args), RT.Symbol("R_NilValue"))) {
		Assign(LocalRef(&deep), RT.Call("asInteger", RT.Call("CADR", args)))
		if RT.Truth(RT.Binary("!=", RT.Call("CDDR", args), RT.Symbol("R_NilValue"))) {
			Assign(LocalRef(&pvec), RT.Call("asInteger", RT.Call("CADDR", args)))
		}
	}
	RT.Call("inspect_tree", RT.Const("int", "0"), RT.Call("CAR", args), deep, pvec)
	return obj
}
