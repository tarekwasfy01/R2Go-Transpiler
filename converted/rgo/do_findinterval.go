package rgo

// do_findinterval translates batch_099.c (src/main/util.c); R primitive(s): findInterval.
func do_findinterval(call, op, args, env Value) Value {
	return findIntervalImpl(arg(args, 1), arg(args, 0), arg(args, 4), arg(args, 2), arg(args, 3))
}
