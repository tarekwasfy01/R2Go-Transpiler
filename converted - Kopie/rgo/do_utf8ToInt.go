package rgo

// do_utf8ToInt translates batch_282.c (src/main/raw.c); R primitive(s): utf8ToInt.
func do_utf8ToInt(call, op, args, env Value) Value {
	return utf8ToIntImpl(arg(args, 0))
}
