package rgo

// do_addTryHandlers translates batch_015.c (src/main/errors.c); R primitive(s): .addTryHandlers.
func do_addTryHandlers(call, op, args, env Value) Value {
	return unsupported("do_addTryHandlers", "requires the evaluator/context stack; isolated batch corpus does not define enough evaluator semantics")
}
