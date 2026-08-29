package rprimitive

// Code generated from GNU R C source by the bundled c2go translation pipeline.
// The control flow is translated to Go; GNU R runtime operations are routed through RT.

func do_signalCondition(call, op, args, rho Value) Value {
	var (
		list     Value
		cond     Value
		msg      Value
		ecall    Value
		oldstack Value
		entry    Value
		h        Value
		msgstr   Value
		hcall    Value
	)
	RT.Call("checkArity", op, args)
	Assign(LocalRef(&cond), RT.Call("CAR", args))
	Assign(LocalRef(&msg), RT.Call("CADR", args))
	Assign(LocalRef(&ecall), RT.Call("CADDR", args))
	RT.Call("PROTECT", Assign(LocalRef(&oldstack), RT.Symbol("R_HandlerStack")))
	for RT.Truth(RT.Binary("!=", Assign(LocalRef(&list), RT.Call("findConditionHandler", cond)), RT.Symbol("R_NilValue"))) {
		Assign(LocalRef(&entry), RT.Call("CAR", list))
		RT.AssignSymbol("R_HandlerStack", RT.Call("CDR", list))
		if RT.Truth(RT.Call("IS_CALLING_ENTRY", entry)) {
			Assign(LocalRef(&h), RT.Call("ENTRY_HANDLER", entry))
			if RT.Truth(RT.Binary("==", h, RT.Symbol("R_RestartToken"))) {
				Assign(LocalRef(&msgstr), RT.Symbol("NULL"))
				if RT.Truth(func() Value {
					if !RT.Truth(RT.Binary("==", RT.Call("TYPEOF", msg), RT.Symbol("STRSXP"))) {
						return false
					}
					return RT.Truth(RT.Binary(">", RT.Call("LENGTH", msg), RT.Const("int", "0")))
				}()) {
					Assign(LocalRef(&msgstr), RT.Call("translateChar", RT.Call("STRING_ELT", msg, RT.Const("int", "0"))))
				} else {
					RT.Call("error", RT.Call("_", RT.Const("string", "\"error message not a string\"")))
				}
				RT.Call("errorcall_dflt", ecall, RT.Const("string", "\"%s\""), msgstr)
			} else {
				Assign(LocalRef(&hcall), RT.Call("LCONS", h, RT.Call("LCONS", cond, RT.Symbol("R_NilValue"))))
				RT.Call("PROTECT", hcall)
				RT.Call("eval", hcall, RT.Symbol("R_GlobalEnv"))
				RT.Call("UNPROTECT", RT.Const("int", "1"))
			}
		} else {
			RT.Call("gotoExitingHandler", cond, ecall, entry)
		}
	}
	RT.AssignSymbol("R_HandlerStack", oldstack)
	RT.Call("UNPROTECT", RT.Const("int", "1"))
	return RT.Symbol("R_NilValue")
}
