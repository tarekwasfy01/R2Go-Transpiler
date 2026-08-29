package rgo

// do_readEnviron translates batch_205.c (src/main/Renviron.c); R primitive(s): readRenviron.
func do_readEnviron(call, op, args, env Value) Value {
	return readEnvironImpl(arg(args, 0))
}
