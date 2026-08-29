package rgo

// do_altrep_class translates batch_020.c (src/main/altrep.c); R primitive(s): altrep_class.
func do_altrep_class(call, op, args, env Value) Value {
	return s4Impl("do_altrep_class", args)
}
