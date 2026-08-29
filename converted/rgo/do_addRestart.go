package rgo

// do_addRestart translates batch_014.c (src/main/errors.c); R primitive(s): .addRestart.
func do_addRestart(call, op, args, env Value) Value {
	return unsupported("do_addRestart", "requires the evaluator/context stack; isolated batch corpus does not define enough evaluator semantics")
}
