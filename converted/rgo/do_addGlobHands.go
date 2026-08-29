package rgo

// do_addGlobHands translates batch_013.c (src/main/errors.c); R primitive(s): .addGlobHands.
func do_addGlobHands(call, op, args, env Value) Value {
	return unsupported("do_addGlobHands", "requires the evaluator/context stack; isolated batch corpus does not define enough evaluator semantics")
}
