package rgo

// do_addCondHands translates batch_012.c (src/main/errors.c); R primitive(s): .addCondHands.
func do_addCondHands(call, op, args, env Value) Value {
	return unsupported("do_addCondHands", "requires the evaluator/context stack; isolated batch corpus does not define enough evaluator semantics")
}
