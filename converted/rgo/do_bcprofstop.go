package rgo

// do_bcprofstop translates batch_030.c (src/main/eval.c); R primitive(s): bcprofstop.
func do_bcprofstop(call, op, args, env Value) Value {
	return unsupported("do_bcprofstop", "requires the evaluator/context stack; isolated batch corpus does not define enough evaluator semantics")
}
