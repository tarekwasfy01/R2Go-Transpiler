package rgo

// do_bcclose translates batch_027.c (src/main/eval.c); R primitive(s): bcClose.
func do_bcclose(call, op, args, env Value) Value {
	return unsupported("do_bcclose", "requires the evaluator/context stack; isolated batch corpus does not define enough evaluator semantics")
}
