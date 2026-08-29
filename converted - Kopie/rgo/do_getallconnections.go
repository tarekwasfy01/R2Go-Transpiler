package rgo

// do_getallconnections translates batch_117.c (src/main/connections.c); R primitive(s): getAllConnections.
func do_getallconnections(call, op, args, env Value) Value {
	return connectionPrimitive("do_getallconnections", args)
}
