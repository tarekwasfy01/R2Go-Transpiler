package rgo

// do_chartr translates batch_041.c (src/main/character.c); R primitive(s): chartr.
func do_chartr(call, op, args, env Value) Value {
	return chartrImpl(arg(args, 0), arg(args, 1), arg(args, 2))
}
