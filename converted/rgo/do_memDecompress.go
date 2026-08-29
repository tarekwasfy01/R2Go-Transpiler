package rgo

// do_memDecompress translates batch_165.c (src/main/connections.c); R primitive(s): memDecompress.
func do_memDecompress(call, op, args, env Value) Value {
	return connectionPrimitive("do_memDecompress", args)
}
