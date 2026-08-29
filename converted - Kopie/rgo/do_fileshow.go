package rgo

// do_fileshow translates batch_097.c (src/main/platform.c); R primitive(s): file.show.
func do_fileshow(call, op, args, env Value) Value {
	return filePrimitive("do_fileshow", args)
}
