package rgo

// do_fileremove translates batch_095.c (src/main/platform.c); R primitive(s): file.remove.
func do_fileremove(call, op, args, env Value) Value {
	return filePrimitive("do_fileremove", args)
}
