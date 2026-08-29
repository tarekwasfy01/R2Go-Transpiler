package rgo

// do_memCompress translates batch_164.c (src/main/connections.c); R primitive(s): memCompress.
func do_memCompress(call, op, args, env Value) Value {
	return connectionPrimitive("do_memCompress", args)
}
