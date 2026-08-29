package rgo

// do_filerename translates batch_096.c (src/main/platform.c); R primitive(s): file.rename.
func do_filerename(call, op, args, env Value) Value {
	return filePrimitive("do_filerename", args)
}
