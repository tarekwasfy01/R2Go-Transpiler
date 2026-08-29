package rgo

// do_gzfile translates batch_125.c (src/main/connections.c); R primitive(s): bzfile, gzfile, xzfile, zstdfile.
func do_gzfile(call, op, args, env Value) Value {
	return connectionPrimitive("do_gzfile", args)
}
