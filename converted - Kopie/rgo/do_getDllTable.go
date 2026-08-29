package rgo

// do_getDllTable translates batch_109.c (src/main/Rdynload.c); R primitive(s): getLoadedDLLs.
func do_getDllTable(call, op, args, env Value) Value {
	return unsupported("do_getDllTable", "native dynamic loading is unavailable without cgo; this is the defined Pure-Go error path")
}
