package rgo

// do_CDotsLength translates batch_003.c (src/main/envir.c); R primitive(s): dotsLength.
func do_CDotsLength(call, op, args, env Value) Value {
	return Ints(int64(nargs(args)))
}
