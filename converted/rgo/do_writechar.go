package rgo

// do_writechar translates batch_289.c (src/main/connections.c); R primitive(s): writeChar.
func do_writechar(call, op, args, env Value) Value {
	return connectionPrimitive("do_writechar", args)
}
