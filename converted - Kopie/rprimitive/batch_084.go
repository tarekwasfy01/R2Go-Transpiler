package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_munmap_file(call, op, args, env Value) Value {
	var (
		x    Value
		eptr Value
	)
	Assign(LocalRef(&x), RT.Call("CAR", args))
	if RT.Truth(RT.Unary("!", func() Value {
		if RT.Truth(RT.Call("R_altrep_inherits", x, RT.Symbol("mmap_integer_class"))) {
			return true
		}
		return RT.Truth(RT.Call("R_altrep_inherits", x, RT.Symbol("mmap_real_class")))
	}())) {
		RT.Call("error", RT.Const("string", "\"not a memory-mapped object\""))
	}
	Assign(LocalRef(&eptr), RT.Call("MMAP_EPTR", x))
	RT.AssignSymbol("errno", RT.Const("int", "0"))
	RT.Call("R_RunWeakRefFinalizer", RT.Call("R_ExternalPtrTag", eptr))
	if RT.Truth(RT.Symbol("errno")) {
		RT.Call("error", RT.Const("string", "\"munmap: %s\""), RT.Call("strerror", RT.Symbol("errno")))
	}
	return RT.Symbol("R_NilValue")
}
