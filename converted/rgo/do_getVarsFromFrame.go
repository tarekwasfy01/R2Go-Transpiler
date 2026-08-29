package rgo

// do_getVarsFromFrame translates batch_116.c (src/main/serialize.c); R primitive(s): getVarsFromFrame.
func do_getVarsFromFrame(call, op, args, env Value) Value {
	return serializePrimitive("do_getVarsFromFrame", args)
}
