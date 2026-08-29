package rgo

// do_dotcall translates batch_068.c (src/main/dotcode.c); R primitive(s): .Call.
func do_dotcall(call, op, args, env Value) Value {
	return unsupported("do_dotcall", "native dynamic loading is unavailable without cgo; this is the defined Pure-Go error path")
}
