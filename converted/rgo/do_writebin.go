package rgo

// do_writebin translates batch_288.c (src/main/connections.c); R primitive(s): writeBin.
func do_writebin(call, op, args, env Value) Value {
	return connectionPrimitive("do_writebin", args)
}
