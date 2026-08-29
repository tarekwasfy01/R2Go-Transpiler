package rgo

// do_remove translates batch_217.c (src/main/envir.c); R primitive(s): remove.
func do_remove(call, op, args, env Value) Value {
	return envPrimitive("do_remove", args, env)
}
