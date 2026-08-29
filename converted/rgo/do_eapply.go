package rgo

// do_eapply translates batch_078.c (src/main/envir.c); R primitive(s): eapply.
func do_eapply(call, op, args, env Value) Value {
	return envPrimitive("do_eapply", args, env)
}
