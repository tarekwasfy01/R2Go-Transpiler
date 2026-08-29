package rgo

// do_tracemem translates batch_272.c (src/main/debug.c); R primitive(s): tracemem.
func do_tracemem(call, op, args, env Value) Value {
	return addressImpl(arg(args, 0))
}
