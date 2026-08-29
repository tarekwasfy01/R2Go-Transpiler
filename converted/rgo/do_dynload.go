package rgo

// do_dynload translates batch_076.c (src/main/Rdynload.c); R primitive(s): dyn.load.
func do_dynload(call, op, args, env Value) Value {
	return unsupported("do_dynload", "native dynamic loading is unavailable without cgo; this is the defined Pure-Go error path")
}
