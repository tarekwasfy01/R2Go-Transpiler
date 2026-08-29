package rgo

// do_close translates batch_043.c (src/main/connections.c); R primitive(s): close.
func do_close(call, op, args, env Value) Value {
	return connectionPrimitive("do_close", args)
}
