package rgo

// do_delayed translates batch_055.c (src/main/builtin.c); R primitive(s): delayedAssign.
func do_delayed(call, op, args, env Value) Value {
	return unsupported("do_delayed", "algorithm depends on GNU R internals not included in this batch corpus")
}
