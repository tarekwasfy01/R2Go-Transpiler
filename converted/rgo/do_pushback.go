package rgo

// do_pushback translates batch_192.c (src/main/connections.c); R primitive(s): pushBack.
func do_pushback(call, op, args, env Value) Value {
	return connectionPrimitive("do_pushback", args)
}
