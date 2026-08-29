package rgo

// do_flush translates batch_100.c (src/main/connections.c); R primitive(s): flush.
func do_flush(call, op, args, env Value) Value {
	return connectionPrimitive("do_flush", args)
}
