package rgo

// do_adist translates batch_017.c (src/main/agrep.c); R primitive(s): adist.
func do_adist(call, op, args, env Value) Value {
	return adistImpl(arg(args, 0), arg(args, 1))
}
