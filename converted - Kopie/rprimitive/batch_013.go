package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_bcprofstart(call, op, args, env Value) Value {
	var (
		itv       Value
		interval  Value
		dinterval Value
		i         Value
	)
	Assign(LocalRef(&itv), RT.NewObject())
	Assign(LocalRef(&dinterval), RT.Const("double", "0.02"))
	RT.Call("checkArity", op, args)
	if RT.Truth(RT.Symbol("R_Profiling")) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"profile timer in use\"")))
	}
	if RT.Truth(RT.Symbol("bc_profiling")) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"already byte code profiling\"")))
	}
	Assign(LocalRef(&interval), RT.Binary("+", RT.Binary("*", RT.Const("double", "1e6"), dinterval), RT.Const("double", "0.5")))
	RT.AssignSymbol("current_opcode", RT.Symbol("NO_CURRENT_OPCODE"))
	for Assign(LocalRef(&i), RT.Const("int", "0")); RT.Truth(RT.Binary("<", i, RT.Symbol("OPCOUNT"))); RT.Inc(LocalRef(&i), 1, true) {
		RT.AssignIndex(RT.Symbol("opcode_counts"), i, RT.Const("int", "0"))
	}
	RT.Call("signal", RT.Symbol("SIGPROF"), RT.Symbol("dobcprof"))
	RT.AssignField(RT.Field(itv, "it_interval"), "tv_sec", RT.Binary("/", interval, RT.Const("int", "1000000")))
	RT.AssignField(RT.Field(itv, "it_interval"), "tv_usec", RT.Call("suseconds_t", RT.Binary("-", interval, RT.Binary("*", RT.Field(RT.Field(itv, "it_interval"), "tv_sec"), RT.Const("int", "1000000")))))
	RT.AssignField(RT.Field(itv, "it_value"), "tv_sec", RT.Binary("/", interval, RT.Const("int", "1000000")))
	RT.AssignField(RT.Field(itv, "it_value"), "tv_usec", RT.Call("suseconds_t", RT.Binary("-", interval, RT.Binary("*", RT.Field(RT.Field(itv, "it_value"), "tv_sec"), RT.Const("int", "1000000")))))
	if RT.Truth(RT.Binary("==", RT.Call("setitimer", RT.Symbol("ITIMER_PROF"), LocalRef(&itv), RT.Symbol("NULL")), RT.Unary("-", RT.Const("int", "1")))) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"setting profile timer failed\"")))
	}
	RT.AssignSymbol("bc_profiling", RT.Symbol("TRUE"))
	return RT.Symbol("R_NilValue")
}
