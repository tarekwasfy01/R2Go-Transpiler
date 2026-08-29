package rgo

// do_truncate translates batch_273.c (src/main/connections.c); R primitive(s): truncate.
func do_truncate(call, op, args, env Value) Value {
	return connectionPrimitive("do_truncate", args)
}
