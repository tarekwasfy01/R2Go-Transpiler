package rgo

// do_filesymlink translates batch_098.c (src/main/platform.c); R primitive(s): file.symlink.
func do_filesymlink(call, op, args, env Value) Value {
	return filePrimitive("do_filesymlink", args)
}
