package rgo

// do_lazyLoadDBflush translates batch_144.c (src/main/serialize.c); R primitive(s): lazyLoadDBflush.
func do_lazyLoadDBflush(call, op, args, env Value) Value {
	return serializePrimitive("do_lazyLoadDBflush", args)
}
