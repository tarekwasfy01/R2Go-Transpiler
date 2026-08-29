package rgo

// do_sink translates batch_239.c (src/main/connections.c); R primitive(s): sink.
func do_sink(call, op, args, env Value) Value {
	return connectionPrimitive("do_sink", args)
}
