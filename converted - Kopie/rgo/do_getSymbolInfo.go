package rgo

// do_getSymbolInfo translates batch_115.c (src/main/Rdynload.c); R primitive(s): getSymbolInfo.
func do_getSymbolInfo(call, op, args, env Value) Value {
	return unsupported("do_getSymbolInfo", "native dynamic loading is unavailable without cgo; this is the defined Pure-Go error path")
}
