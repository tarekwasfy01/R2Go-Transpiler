package rgo

// do_dotcallgr translates batch_069.c (src/main/dotcode.c); R primitive(s): .Call.graphics.
func do_dotcallgr(call, op, args, env Value) Value {
	return unsupported("do_dotcallgr", "native dynamic loading is unavailable without cgo; this is the defined Pure-Go error path")
}
