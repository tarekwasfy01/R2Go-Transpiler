package rgo

// do_rapply translates batch_199.c (src/main/apply.c); R primitive(s): rapply.
func do_rapply(call, op, args, env Value) Value {
	return rapplyImpl(arg(args, 0), arg(args, 1), arg(args, 2), arg(args, 3), env)
}
