package rgo

// do_makenames translates batch_157.c (src/main/character.c); R primitive(s): make.names.
func do_makenames(call, op, args, env Value) Value {
	return makeNamesImpl(arg(args, 0), arg(args, 1))
}
