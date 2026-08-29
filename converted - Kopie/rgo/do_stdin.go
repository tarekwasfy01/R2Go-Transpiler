package rgo

// do_stdin translates batch_247.c (src/main/connections.c); R primitive(s): stdin.
func do_stdin(call, op, args, env Value) Value {
	return connectionPrimitive("do_stdin", args)
}
