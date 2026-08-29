package rgo

// do_serialize translates batch_228.c (src/main/serialize.c); R primitive(s): serializeb.
func do_serialize(call, op, args, env Value) Value {
	return serializePrimitive("do_serialize", args)
}
