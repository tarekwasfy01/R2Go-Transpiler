package rgo

// do_ls translates batch_154.c (src/main/envir.c); R primitive(s): ls.
func do_ls(call, op, args, env Value) Value {
	return envPrimitive("do_ls", args, env)
}
