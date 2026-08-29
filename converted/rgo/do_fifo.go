package rgo

// do_fifo translates batch_087.c (src/main/connections.c); R primitive(s): fifo.
func do_fifo(call, op, args, env Value) Value {
	return connectionPrimitive("do_fifo", args)
}
