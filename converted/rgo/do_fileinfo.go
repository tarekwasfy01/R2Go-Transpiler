package rgo

// do_fileinfo translates batch_093.c (src/main/platform.c); R primitive(s): file.info.
func do_fileinfo(call, op, args, env Value) Value {
	return filePrimitive("do_fileinfo", args)
}
