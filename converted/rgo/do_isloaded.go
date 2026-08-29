package rgo

// do_isloaded translates batch_138.c (src/main/dotcode.c); R primitive(s): is.loaded.
func do_isloaded(call, op, args, env Value) Value {
	return unsupported("do_isloaded", "native dynamic loading is unavailable without cgo; this is the defined Pure-Go error path")
}
