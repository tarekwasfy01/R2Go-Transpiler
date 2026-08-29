package rgo

// do_dfltStop translates batch_059.c (src/main/errors.c); R primitive(s): .dfltStop.
func do_dfltStop(call, op, args, env Value) Value {
	return ErrValue("%s", formatValue(arg(args, 0)))
}
