package rgo

// do_mkcode translates batch_171.c (src/main/eval.c); R primitive(s): mkCode.
func do_mkcode(call, op, args, env Value) Value {
	return unsupported("do_mkcode", "requires the evaluator/context stack; isolated batch corpus does not define enough evaluator semantics")
}
