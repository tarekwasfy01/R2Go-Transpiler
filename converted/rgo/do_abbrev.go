package rgo

// do_abbrev translates batch_010.c (src/main/character.c); R primitive(s): abbreviate.
func do_abbrev(call, op, args, env Value) Value {
	return abbreviateImpl(arg(args, 0), arg(args, 1), arg(args, 2), arg(args, 3))
}
