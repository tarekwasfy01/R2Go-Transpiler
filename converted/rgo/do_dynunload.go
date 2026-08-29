package rgo

// do_dynunload translates batch_077.c (src/main/Rdynload.c); R primitive(s): dyn.unload.
func do_dynunload(call, op, args, env Value) Value {
	return unsupported("do_dynunload", "native dynamic loading is unavailable without cgo; this is the defined Pure-Go error path")
}
