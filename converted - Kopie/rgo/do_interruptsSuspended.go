package rgo

// do_interruptsSuspended translates batch_131.c (src/main/errors.c); R primitive(s): interruptsSuspended.
func do_interruptsSuspended(call, op, args, env Value) Value {
	return systemPrimitive("do_interruptsSuspended", args)
}
