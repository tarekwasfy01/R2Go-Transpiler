package rgo

// do_regexec translates batch_215.c (src/main/grep.c); R primitive(s): regexec.
func do_regexec(call, op, args, env Value) Value {
	return aregexecImpl(arg(args, 0), arg(args, 1))
}
