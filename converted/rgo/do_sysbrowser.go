package rgo

// do_sysbrowser translates batch_255.c (src/main/context.c); R primitive(s): browserCondition, browserSetDebug, browserText.
func do_sysbrowser(call, op, args, env Value) Value {
	return unsupported("do_sysbrowser", "requires the evaluator/context stack; isolated batch corpus does not define enough evaluator semantics")
}
