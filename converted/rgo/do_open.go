package rgo

// do_open translates batch_181.c (src/main/connections.c); R primitive(s): open.
func do_open(call, op, args, env Value) Value {
	return connectionPrimitive("do_open", args)
}
