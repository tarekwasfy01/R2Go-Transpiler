package rgo

// do_unz translates batch_280.c (src/main/connections.c); R primitive(s): unz.
func do_unz(call, op, args, env Value) Value {
	return connectionPrimitive("do_unz", args)
}
