package rgo

// do_gzcon translates batch_124.c (src/main/connections.c); R primitive(s): gzcon.
func do_gzcon(call, op, args, env Value) Value {
	return connectionPrimitive("do_gzcon", args)
}
