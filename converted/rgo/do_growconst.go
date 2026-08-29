package rgo

// do_growconst translates batch_123.c (src/main/eval.c); R primitive(s): growconst.
func do_growconst(call, op, args, env Value) Value {
	return unsupported("do_growconst", "requires the evaluator/context stack; isolated batch corpus does not define enough evaluator semantics")
}
