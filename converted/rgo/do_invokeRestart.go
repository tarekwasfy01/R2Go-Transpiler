package rgo

// do_invokeRestart translates batch_132.c (src/main/errors.c); R primitive(s): .invokeRestart.
func do_invokeRestart(call, op, args, env Value) Value {
	return unsupported("do_invokeRestart", "requires the evaluator/context stack; isolated batch corpus does not define enough evaluator semantics")
}
