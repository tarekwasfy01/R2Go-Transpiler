package rgo

// do_pushbacklength translates batch_193.c (src/main/connections.c); R primitive(s): pushBackLength.
func do_pushbacklength(call, op, args, env Value) Value {
	return connectionPrimitive("do_pushbacklength", args)
}
