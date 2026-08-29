package rgo

// do_aregexec translates batch_021.c (src/main/agrep.c); R primitive(s): aregexec.
func do_aregexec(call, op, args, env Value) Value {
	return aregexecImpl(arg(args, 0), arg(args, 1))
}
