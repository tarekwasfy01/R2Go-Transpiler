package rgo

// do_getRegisteredRoutines translates batch_113.c (src/main/Rdynload.c); R primitive(s): getRegisteredRoutines.
func do_getRegisteredRoutines(call, op, args, env Value) Value {
	return unsupported("do_getRegisteredRoutines", "native dynamic loading is unavailable without cgo; this is the defined Pure-Go error path")
}
