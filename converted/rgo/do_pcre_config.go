package rgo

// do_pcre_config translates batch_185.c (src/main/grep.c); R primitive(s): pcre_config.
func do_pcre_config(call, op, args, env Value) Value {
	return setNames(Lists(Bool(true), Bool(true), Bool(false)), []string{"UTF-8", "Unicode properties", "JIT"})
}
