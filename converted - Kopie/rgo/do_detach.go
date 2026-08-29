package rgo

// do_detach translates batch_058.c (src/main/envir.c); R primitive(s): detach.
func do_detach(call, op, args, env Value) Value {
	return envPrimitive("do_detach", args, env)
}
