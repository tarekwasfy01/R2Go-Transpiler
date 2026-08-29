package rgo

// do_bndIsActive translates batch_035.c (src/main/envir.c); R primitive(s): bindingIsActive.
func do_bndIsActive(call, op, args, env Value) Value {
	return envPrimitive("do_bndIsActive", args, env)
}
