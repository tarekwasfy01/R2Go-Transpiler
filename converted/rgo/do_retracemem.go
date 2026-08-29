package rgo

// do_retracemem translates batch_219.c (src/main/debug.c); R primitive(s): retracemem.
func do_retracemem(call, op, args, env Value) Value {
	return addressImpl(arg(args, 0))
}
