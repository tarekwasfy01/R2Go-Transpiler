package rgo

// do_sumconnection translates batch_253.c (src/main/connections.c); R primitive(s): summary.connection.
func do_sumconnection(call, op, args, env Value) Value {
	return connectionPrimitive("do_sumconnection", args)
}
