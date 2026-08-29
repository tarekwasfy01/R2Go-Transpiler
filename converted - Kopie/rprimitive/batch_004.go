package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_addCondHands(call, op, args, rho Value) Value {
	var (
		classes   Value
		handlers  Value
		parentenv Value
		target    Value
		oldstack  Value
		newstack  Value
		result    Value
		calling   Value
		i         Value
		n         Value
		osi       Value
		klass     Value
		handler   Value
		entry     Value
	)
	if RT.Truth(RT.Binary("==", RT.Symbol("R_HandlerResultToken"), RT.Symbol("NULL"))) {
		RT.AssignSymbol("R_HandlerResultToken", RT.Call("allocVector", RT.Symbol("VECSXP"), RT.Const("int", "1")))
		RT.Call("R_PreserveObject", RT.Symbol("R_HandlerResultToken"))
	}
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&classes), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&handlers), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&parentenv), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&target), RT.Call("CAR", args))
	Assign(LocalRef(&args), RT.Call("CDR", args))
	Assign(LocalRef(&calling), RT.Call("asLogical", RT.Call("CAR", args)))
	if RT.Truth(func() Value {
		if RT.Truth(RT.Binary("==", classes, RT.Symbol("R_NilValue"))) {
			return true
		}
		return RT.Truth(RT.Binary("==", handlers, RT.Symbol("R_NilValue")))
	}()) {
		return RT.Symbol("R_HandlerStack")
	}
	if RT.Truth(func() Value {
		if RT.Truth(func() Value {
			if RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", classes), RT.Symbol("STRSXP"))) {
				return true
			}
			return RT.Truth(RT.Binary("!=", RT.Call("TYPEOF", handlers), RT.Symbol("VECSXP")))
		}()) {
			return true
		}
		return RT.Truth(RT.Binary("!=", RT.Call("LENGTH", classes), RT.Call("LENGTH", handlers)))
	}()) {
		RT.Call("error", RT.Call("_", RT.Const("string", "\"bad handler data\"")))
	}
	Assign(LocalRef(&n), RT.Call("LENGTH", handlers))
	Assign(LocalRef(&oldstack), RT.Symbol("R_HandlerStack"))
	RT.Call("PROTECT", Assign(LocalRef(&result), RT.Call("allocVector", RT.Symbol("VECSXP"), RT.Symbol("RESULT_SIZE"))))
	RT.Call("SET_VECTOR_ELT", result, RT.Binary("-", RT.Symbol("RESULT_SIZE"), RT.Const("int", "1")), RT.Symbol("R_HandlerResultToken"))
	RT.Call("PROTECT_WITH_INDEX", Assign(LocalRef(&newstack), oldstack), LocalRef(&osi))
	for Assign(LocalRef(&i), RT.Binary("-", n, RT.Const("int", "1"))); RT.Truth(RT.Binary(">=", i, RT.Const("int", "0"))); RT.Inc(LocalRef(&i), -1, true) {
		Assign(LocalRef(&klass), RT.Call("STRING_ELT", classes, i))
		Assign(LocalRef(&handler), RT.Call("VECTOR_ELT", handlers, i))
		Assign(LocalRef(&entry), RT.Call("mkHandlerEntry", klass, parentenv, handler, target, result, calling))
		RT.Call("REPROTECT", Assign(LocalRef(&newstack), RT.Call("CONS", entry, newstack)), osi)
	}
	RT.AssignSymbol("R_HandlerStack", newstack)
	RT.Call("UNPROTECT", RT.Const("int", "2"))
	return oldstack
}
