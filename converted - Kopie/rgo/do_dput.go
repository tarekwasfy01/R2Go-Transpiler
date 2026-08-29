package rgo

// do_dput translates batch_074.c (src/main/deparse.c); R primitive(s): dput.
func do_dput(call, op, args, env Value) Value {
	return dputImpl(arg(args, 0), arg(args, 1))
}
