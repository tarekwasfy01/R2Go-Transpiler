package rgo

// do_attach translates batch_026.c (src/main/envir.c); R primitive(s): attach.
func do_attach(call, op, args, env Value) Value {
	return envPrimitive("do_attach", args, env)
}
