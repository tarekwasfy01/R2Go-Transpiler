package rgo

// do_dotsExist translates batch_071.c (src/main/envir.c); R primitive(s): dotsExist.
func do_dotsExist(call, op, args, env Value) Value {
	return Bool(nargs(args) > 0)
}
