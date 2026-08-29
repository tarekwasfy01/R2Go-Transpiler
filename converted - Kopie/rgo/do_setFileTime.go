package rgo

// do_setFileTime translates batch_231.c (src/main/platform.c); R primitive(s): setFileTime.
func do_setFileTime(call, op, args, env Value) Value {
	return filePrimitive("do_setFileTime", args)
}
