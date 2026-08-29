package rgo

// do_getconnection translates batch_118.c (src/main/connections.c); R primitive(s): getConnection.
func do_getconnection(call, op, args, env Value) Value {
	return connectionPrimitive("do_getconnection", args)
}
