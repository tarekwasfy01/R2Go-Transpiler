package rgo

// do_seterrmessage translates batch_233.c (src/main/errors.c); R primitive(s): seterrmessage.
func do_seterrmessage(call, op, args, env Value) Value {
	return statePrimitive("do_seterrmessage", args)
}
