package rgo

// do_objsxp translates batch_179.c (src/main/objects.c); R primitive(s): objsxp.
func do_objsxp(call, op, args, env Value) Value {
	return Bool(arg(args, 0).Kind == Environment || arg(args, 0).Kind == Function)
}
