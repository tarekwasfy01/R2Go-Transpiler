package rgo

// do_declare translates batch_054.c (src/main/eval.c); R primitive(s): declare.
func do_declare(call, op, args, env Value) Value {
	return unsupported("do_declare", "requires the evaluator/context stack; isolated batch corpus does not define enough evaluator semantics")
}
