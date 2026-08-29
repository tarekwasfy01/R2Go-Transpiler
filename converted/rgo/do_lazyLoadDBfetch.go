package rgo

// do_lazyLoadDBfetch translates batch_143.c (src/main/serialize.c); R primitive(s): lazyLoadDBfetch.
func do_lazyLoadDBfetch(call, op, args, env Value) Value {
	return serializePrimitive("do_lazyLoadDBfetch", args)
}
