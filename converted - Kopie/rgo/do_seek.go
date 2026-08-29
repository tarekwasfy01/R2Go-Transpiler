package rgo

// do_seek translates batch_226.c (src/main/connections.c); R primitive(s): seek.
func do_seek(call, op, args, env Value) Value {
	return connectionPrimitive("do_seek", args)
}
