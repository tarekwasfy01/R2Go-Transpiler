package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_bcprofstop(call, op, args, env Value) Value {
	var (
		itv Value
	)
	Assign(LocalRef(&itv), RT.NewObject())
	RT.Call("checkArity", op, args)
	if RT.Truth(RT.Unary("!", RT.Symbol("bc_profiling"))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"not byte code profiling\"")))
	}
	RT.AssignField(RT.Field(itv, "it_interval"), "tv_sec", RT.Const("int", "0"))
	RT.AssignField(RT.Field(itv, "it_interval"), "tv_usec", RT.Const("int", "0"))
	RT.AssignField(RT.Field(itv, "it_value"), "tv_sec", RT.Const("int", "0"))
	RT.AssignField(RT.Field(itv, "it_value"), "tv_usec", RT.Const("int", "0"))
	RT.Call("setitimer", RT.Symbol("ITIMER_PROF"), LocalRef(&itv), RT.Symbol("NULL"))
	RT.Call("signal", RT.Symbol("SIGPROF"), RT.Symbol("dobcprof_null"))
	RT.AssignSymbol("bc_profiling", RT.Symbol("FALSE"))
	return RT.Symbol("R_NilValue")
}
