package rgo

// do_strptime translates batch_250.c (src/main/datetime.c); R primitive(s): strptime.
func do_strptime(call, op, args, env Value) Value {
	return strptimeImpl(arg(args, 0), arg(args, 1))
}
