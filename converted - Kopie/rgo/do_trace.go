package rgo

// do_trace translates batch_269.c (src/main/debug.c); R primitive(s): .primTrace, .primUntrace.
func do_trace(call, op, args, env Value) Value {
	return arg(args, 0)
}
