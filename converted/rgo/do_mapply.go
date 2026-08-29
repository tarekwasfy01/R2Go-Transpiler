package rgo

// do_mapply translates batch_158.c (src/main/mapply.c); R primitive(s): mapply.
func do_mapply(call, op, args, env Value) Value {
	return mapplyImpl(arg(args, 0), arg(args, 1), arg(args, 2), env)
}
