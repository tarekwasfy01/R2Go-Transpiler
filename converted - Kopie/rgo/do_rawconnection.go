package rgo

// do_rawconnection translates batch_202.c (src/main/connections.c); R primitive(s): rawConnection.
func do_rawconnection(call, op, args, env Value) Value {
	return connectionPrimitive("do_rawconnection", args)
}
