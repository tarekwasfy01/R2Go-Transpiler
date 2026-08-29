package rgo

// do_dotsNames translates batch_073.c (src/main/envir.c); R primitive(s): ...names.
func do_dotsNames(call, op, args, env Value) Value {
	return Strings(names(args)...)
}
