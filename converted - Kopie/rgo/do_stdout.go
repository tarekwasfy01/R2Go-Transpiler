package rgo

// do_stdout translates batch_248.c (src/main/connections.c); R primitive(s): stdout.
func do_stdout(call, op, args, env Value) Value {
	return connectionPrimitive("do_stdout", args)
}
