package rgo

// do_serversocket translates batch_230.c (src/main/connections.c); R primitive(s): serverSocket.
func do_serversocket(call, op, args, env Value) Value {
	return connectionPrimitive("do_serversocket", args)
}
