package rgo

// do_shortRowNames translates batch_237.c (src/main/attrib.c); R primitive(s): shortRowNames.
func do_shortRowNames(call, op, args, env Value) Value {
	return shortRowNamesImpl(arg(args, 0), arg(args, 1))
}
