package rgo

// do_filecopy translates batch_091.c (src/main/platform.c); R primitive(s): file.copy.
func do_filecopy(call, op, args, env Value) Value {
	return filePrimitive("do_filecopy", args)
}
