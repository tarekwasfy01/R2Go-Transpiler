package rgo

// do_disassemble translates batch_061.c (src/main/eval.c); R primitive(s): disassemble.
func do_disassemble(call, op, args, env Value) Value {
	return unsupported("do_disassemble", "requires the evaluator/context stack; isolated batch corpus does not define enough evaluator semantics")
}
