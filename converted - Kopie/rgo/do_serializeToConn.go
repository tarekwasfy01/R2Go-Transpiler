package rgo

// do_serializeToConn translates batch_229.c (src/main/serialize.c); R primitive(s): serializeToConn.
func do_serializeToConn(call, op, args, env Value) Value {
	return serializePrimitive("do_serializeToConn", args)
}
