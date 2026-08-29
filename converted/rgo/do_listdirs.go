package rgo

// do_listdirs translates batch_148.c (src/main/platform.c); R primitive(s): list.dirs.
func do_listdirs(call, op, args, env Value) Value {
	return filePrimitive("do_listdirs", args)
}
