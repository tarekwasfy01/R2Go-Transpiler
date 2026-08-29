package rgo

// do_CDotsNames translates batch_004.c (src/main/envir.c); R primitive(s): dotsNames.
func do_CDotsNames(call, op, args, env Value) Value {
	return Strings(names(args)...)
}
