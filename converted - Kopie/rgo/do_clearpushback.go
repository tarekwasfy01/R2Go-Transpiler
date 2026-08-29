package rgo

// do_clearpushback translates batch_042.c (src/main/connections.c); R primitive(s): clearPushBack.
func do_clearpushback(call, op, args, env Value) Value {
	return connectionPrimitive("do_clearpushback", args)
}
