package rgo

// do_filelink translates batch_094.c (src/main/platform.c); R primitive(s): file.link.
func do_filelink(call, op, args, env Value) Value {
	return filePrimitive("do_filelink", args)
}
