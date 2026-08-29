package rgo

// do_debug translates batch_053.c (src/main/debug.c); R primitive(s): debug, debugonce, isdebugged, undebug.
func do_debug(call, op, args, env Value) Value {
	return unsupported("do_debug", "requires the evaluator/context stack; isolated batch corpus does not define enough evaluator semantics")
}
