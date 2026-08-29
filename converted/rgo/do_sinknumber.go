package rgo

// do_sinknumber translates batch_240.c (src/main/connections.c); R primitive(s): sink.number.
func do_sinknumber(call, op, args, env Value) Value {
	return connectionPrimitive("do_sinknumber", args)
}
