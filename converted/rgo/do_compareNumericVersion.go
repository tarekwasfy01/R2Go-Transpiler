package rgo

// do_compareNumericVersion translates batch_045.c (src/main/util.c); R primitive(s): compareNumericVersion.
func do_compareNumericVersion(call, op, args, env Value) Value {
	return compareNumericVersionImpl(arg(args, 0), arg(args, 1))
}
