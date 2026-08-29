package rgo

// do_External translates batch_005.c (src/main/dotcode.c); R primitive(s): .External, .External2.
func do_External(call, op, args, env Value) Value {
	return unsupported("do_External", "native dynamic loading is unavailable without cgo; this is the defined Pure-Go error path")
}
