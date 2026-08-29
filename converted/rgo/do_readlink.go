package rgo

// do_readlink translates batch_208.c (src/main/platform.c); R primitive(s): Sys.readlink.
func do_readlink(call, op, args, env Value) Value {
	return filePrimitive("do_readlink", args)
}
