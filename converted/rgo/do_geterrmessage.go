package rgo

// do_geterrmessage translates batch_120.c (src/main/errors.c); R primitive(s): geterrmessage.
func do_geterrmessage(call, op, args, env Value) Value {
	return statePrimitive("do_geterrmessage", args)
}
