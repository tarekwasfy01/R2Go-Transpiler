package rgo

// do_bindtextdomain translates batch_034.c (src/main/errors.c); R primitive(s): bindtextdomain.
func do_bindtextdomain(call, op, args, env Value) Value {
	if nargs(args) < 2 {
		return Strings("")
	}
	return arg(args, 1)
}
