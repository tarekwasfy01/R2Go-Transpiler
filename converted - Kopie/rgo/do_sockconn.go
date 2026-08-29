package rgo

// do_sockconn translates batch_241.c (src/main/connections.c); R primitive(s): socketAccept, socketConnection.
func do_sockconn(call, op, args, env Value) Value {
	return connectionPrimitive("do_sockconn", args)
}
