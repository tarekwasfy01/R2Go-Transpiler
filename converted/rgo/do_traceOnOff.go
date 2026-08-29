package rgo

// do_traceOnOff translates batch_270.c (src/main/debug.c); R primitive(s): debugOnOff, traceOnOff.
func do_traceOnOff(call, op, args, env Value) Value {
	return statePrimitive("do_traceOnOff", args)
}
