package rgo

// do_ICUset translates batch_008.c (src/main/util.c); R primitive(s): icuSetCollate.
func do_ICUset(call, op, args, env Value) Value {
	return statePrimitive("do_ICUset", args)
}
