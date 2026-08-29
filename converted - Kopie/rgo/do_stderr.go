package rgo

// do_stderr translates batch_246.c (src/main/connections.c); R primitive(s): stderr.
func do_stderr(call, op, args, env Value) Value {
	return connectionPrimitive("do_stderr", args)
}
