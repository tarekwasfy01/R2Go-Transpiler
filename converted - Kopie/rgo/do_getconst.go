package rgo

// do_getconst translates batch_119.c (src/main/eval.c); R primitive(s): getconst.
func do_getconst(call, op, args, env Value) Value {
	return unsupported("do_getconst", "requires the evaluator/context stack; isolated batch corpus does not define enough evaluator semantics")
}
