package rgo

// do_putconst translates batch_194.c (src/main/eval.c); R primitive(s): putconst.
func do_putconst(call, op, args, env Value) Value {
	return unsupported("do_putconst", "requires the evaluator/context stack; isolated batch corpus does not define enough evaluator semantics")
}
