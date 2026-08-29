package rgo

// do_fileaccess translates batch_088.c (src/main/platform.c); R primitive(s): file.access.
func do_fileaccess(call, op, args, env Value) Value {
	return filePrimitive("do_fileaccess", args)
}
