package rgo

// do_pipe translates batch_186.c (src/main/connections.c); R primitive(s): pipe.
func do_pipe(call, op, args, env Value) Value {
	return connectionPrimitive("do_pipe", args)
}
