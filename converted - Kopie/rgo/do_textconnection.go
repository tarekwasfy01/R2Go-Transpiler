package rgo

// do_textconnection translates batch_266.c (src/main/connections.c); R primitive(s): textConnection.
func do_textconnection(call, op, args, env Value) Value {
	return connectionPrimitive("do_textconnection", args)
}
