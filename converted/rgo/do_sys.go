package rgo

// do_sys translates batch_254.c (src/main/context.c); R primitive(s): sys.call, sys.calls, sys.frame, sys.frames, sys.function, sys.nframe, sys.on.exit, sys.parent, sys.parents.
func do_sys(call, op, args, env Value) Value {
	return unsupported("do_sys", "requires the evaluator/context stack; isolated batch corpus does not define enough evaluator semantics")
}
