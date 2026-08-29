package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_regFinaliz(call, op, args, rho Value) Value {
	var (
		onexit Value
	)
	RT.Call("checkArity", op, args)
	if RT.Truth(func() Value {
		if !RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", RT.Call("CAR", args)), RT.Symbol("ENVSXP"))) {
			return false
		}
		return RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", RT.Call("CAR", args)), RT.Symbol("EXTPTRSXP")))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"first argument must be environment or external pointer\"")))
	}
	if RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", RT.Call("CADR", args)), RT.Symbol("CLOSXP"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"second argument must be a function\"")))
	}
	Assign(LocalRef(&onexit), RT.Call("asLogical", RT.Call("CADDR", args)))
	if RT.Truth(RT.Binary("==", onexit, RT.Symbol("NA_LOGICAL"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"third argument must be 'TRUE' or 'FALSE'\"")))
	}
	RT.Call("R_RegisterFinalizerEx", RT.Call("CAR", args), RT.Call("CADR", args), RT.Cast("Rboolean", onexit))
	return RT.Symbol("R_NilValue")
}
