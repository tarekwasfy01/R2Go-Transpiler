package rgo

// do_ICUget translates batch_007.c (src/main/util.c); R primitive(s): icuGetCollate.
func do_ICUget(call, op, args, env Value) Value {
	return statePrimitive("do_ICUget", args)
}
