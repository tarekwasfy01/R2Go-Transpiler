package rgo

// do_regexpr translates batch_216.c (src/main/grep.c); R primitive(s): gregexpr, regexpr.
func do_regexpr(call, op, args, env Value) Value {
	return regexprImpl(arg(args, 0), arg(args, 1), false)
}
