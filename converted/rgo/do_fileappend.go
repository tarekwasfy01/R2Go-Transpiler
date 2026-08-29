package rgo

// do_fileappend translates batch_089.c (src/main/platform.c); R primitive(s): file.append.
func do_fileappend(call, op, args, env Value) Value {
	return filePrimitive("do_fileappend", args)
}
