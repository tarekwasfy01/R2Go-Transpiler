package rgo

// do_substrgets translates batch_252.c (src/main/character.c); R primitive(s): substr<-.
func do_substrgets(call, op, args, env Value) Value {
	return substrGetsImpl(arg(args, 0), arg(args, 1), arg(args, 2), arg(args, 3))
}
