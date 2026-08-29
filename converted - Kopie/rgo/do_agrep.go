package rgo

// do_agrep translates batch_018.c (src/main/agrep.c); R primitive(s): agrep, agrepl.
func do_agrep(call, op, args, env Value) Value {
	return agrepImpl(arg(args, 0), arg(args, 1), arg(args, 4), arg(args, 5))
}
