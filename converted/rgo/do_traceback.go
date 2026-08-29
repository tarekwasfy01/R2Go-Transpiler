package rgo

// do_traceback translates batch_271.c (src/main/errors.c); R primitive(s): traceback.
func do_traceback(call, op, args, env Value) Value {
	return unsupported("do_traceback", "requires the evaluator/context stack; isolated batch corpus does not define enough evaluator semantics")
}
