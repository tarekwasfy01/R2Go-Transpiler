package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_addGlobHands(call, op, args, rho Value) Value {
	var (
		oldstk Value
		cptr   Value
	)
	Assign(LocalRef(&oldstk), RT.Field(RT.Symbol("R_ToplevelContext"), "handlerstack"))
	for RT.Sequence(Assign(LocalRef(&cptr), RT.Symbol("R_GlobalContext"))); RT.Truth(RT.Binary("!=", cptr, RT.Symbol("R_ToplevelContext"))); Assign(LocalRef(&cptr), RT.Field(cptr, "nextcontext")) {
		if RT.Truth(RT.Binary("!=", RT.Field(cptr, "handlerstack"), oldstk)) {
			RT.Call("error", RT.Const("string", "\"should not be called with handlers on the stack\""))
		}
	}
	RT.AssignSymbol("R_HandlerStack", RT.Symbol("R_NilValue"))
	do_addCondHands(call, op, args, rho)
	for RT.Sequence(Assign(LocalRef(&cptr), RT.Symbol("R_GlobalContext"))); RT.Truth(RT.Binary("!=", cptr, RT.Symbol("R_ToplevelContext"))); Assign(LocalRef(&cptr), RT.Field(cptr, "nextcontext")) {
		if RT.Truth(RT.Binary("==", RT.Field(cptr, "handlerstack"), oldstk)) {
			RT.AssignField(cptr, "handlerstack", RT.Symbol("R_HandlerStack"))
		} else {
			RT.Call("error", RT.Const("string", "\"should not be called with handlers on the stack\""))
		}
	}
	RT.AssignField(RT.Symbol("R_ToplevelContext"), "handlerstack", RT.Symbol("R_HandlerStack"))
	return RT.Symbol("R_NilValue")
}
