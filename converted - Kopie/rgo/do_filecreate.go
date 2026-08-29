package rgo

// do_filecreate translates batch_092.c (src/main/platform.c); R primitive(s): file.create.
func do_filecreate(call, op, args, env Value) Value {
	return filePrimitive("do_filecreate", args)
}
