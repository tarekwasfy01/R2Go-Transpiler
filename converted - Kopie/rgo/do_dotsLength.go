package rgo

// do_dotsLength translates batch_072.c (src/main/envir.c); R primitive(s): ...length.
func do_dotsLength(call, op, args, env Value) Value {
	return Ints(int64(nargs(args)))
}
