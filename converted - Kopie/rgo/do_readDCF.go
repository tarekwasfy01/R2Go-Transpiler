package rgo

// do_readDCF translates batch_204.c (src/main/dcf.c); R primitive(s): readDCF.
func do_readDCF(call, op, args, env Value) Value {
	return readDCFImpl(arg(args, 0))
}
