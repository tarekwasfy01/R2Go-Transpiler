package rgo

// do_signalCondition translates batch_238.c (src/main/errors.c); R primitive(s): .signalCondition.
func do_signalCondition(call, op, args, env Value) Value {
	return unsupported("do_signalCondition", "requires the evaluator/context stack; isolated batch corpus does not define enough evaluator semantics")
}
