package rgo

// do_bndIsLocked translates batch_036.c (src/main/envir.c); R primitive(s): bindingIsLocked.
func do_bndIsLocked(call, op, args, env Value) Value {
	return envPrimitive("do_bndIsLocked", args, env)
}
