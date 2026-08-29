package rgo

// do_dotCode translates batch_063.c (src/main/dotcode.c); R primitive(s): .C, .Fortran.
func do_dotCode(call, op, args, env Value) Value {
	return unsupported("do_dotCode", "native dynamic loading is unavailable without cgo; this is the defined Pure-Go error path")
}
