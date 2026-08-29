package rgo

// do_islistfactor translates batch_137.c (src/main/apply.c); R primitive(s): islistfactor.
func do_islistfactor(call, op, args, env Value) Value {
	return isListFactorImpl(arg(args, 0), arg(args, 1))
}
