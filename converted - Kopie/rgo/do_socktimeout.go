package rgo

// do_socktimeout translates batch_243.c (src/main/connections.c); R primitive(s): socketTimeout.
func do_socktimeout(call, op, args, env Value) Value {
	return connectionPrimitive("do_socktimeout", args)
}
