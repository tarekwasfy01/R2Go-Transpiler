package rgo

// do_readbin translates batch_206.c (src/main/connections.c); R primitive(s): readBin.
func do_readbin(call, op, args, env Value) Value {
	return connectionPrimitive("do_readbin", args)
}
