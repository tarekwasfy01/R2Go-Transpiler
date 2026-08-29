package rgo

// do_curlVersion translates batch_052.c (src/main/internet.c); R primitive(s): curlVersion.
func do_curlVersion(call, op, args, env Value) Value {
	return unsupported("do_curlVersion", "algorithm depends on GNU R internals not included in this batch corpus")
}
