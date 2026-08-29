package rgo

// do_tryCatchHelper translates batch_274.c (src/main/errors.c); R primitive(s): C_tryCatchHelper.
func do_tryCatchHelper(call, op, args, env Value) Value {
	return unsupported("do_tryCatchHelper", "requires the evaluator/context stack; isolated batch corpus does not define enough evaluator semantics")
}
