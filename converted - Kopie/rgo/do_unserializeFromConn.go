package rgo

// do_unserializeFromConn translates batch_278.c (src/main/serialize.c); R primitive(s): serializeInfoFromConn, unserializeFromConn.
func do_unserializeFromConn(call, op, args, env Value) Value {
	return serializePrimitive("do_unserializeFromConn", args)
}
