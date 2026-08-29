package rgo

// do_polyroot translates batch_187.c (src/main/complex.c); R primitive(s): polyroot.
func do_polyroot(call, op, args, env Value) Value {
	return polyrootImpl(arg(args, 0))
}
