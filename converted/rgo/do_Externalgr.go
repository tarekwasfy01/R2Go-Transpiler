package rgo

// do_Externalgr translates batch_006.c (src/main/dotcode.c); R primitive(s): .External.graphics.
func do_Externalgr(call, op, args, env Value) Value {
	return unsupported("do_Externalgr", "native dynamic loading is unavailable without cgo; this is the defined Pure-Go error path")
}
