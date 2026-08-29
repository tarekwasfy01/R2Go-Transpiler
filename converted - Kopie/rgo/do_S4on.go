package rgo

// do_S4on translates batch_009.c (src/main/objects.c); R primitive(s): .isMethodsDispatchOn.
func do_S4on(call, op, args, env Value) Value {
	return statePrimitive("do_S4on", args)
}
