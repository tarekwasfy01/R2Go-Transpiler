package rgo

// do_mkjunction translates batch_172.c (src/main/platform.c); R primitive(s): mkjunction.
func do_mkjunction(call, op, args, env Value) Value {
	return filePrimitive("do_mkjunction", args)
}
