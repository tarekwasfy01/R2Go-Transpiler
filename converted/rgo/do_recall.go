package rgo

// do_recall translates batch_210.c (src/main/eval.c); R primitive(s): Recall.
func do_recall(call, op, args, env Value) Value {
	return unsupported("do_recall", "requires the evaluator/context stack; isolated batch corpus does not define enough evaluator semantics")
}
