package rgo

// do_enablejit translates batch_079.c (src/main/eval.c); R primitive(s): enableJIT.
func do_enablejit(call, op, args, env Value) Value {
	return statePrimitive("do_enablejit", args)
}
