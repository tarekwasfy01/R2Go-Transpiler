package rgo

// do_dfltWarn translates batch_060.c (src/main/errors.c); R primitive(s): .dfltWarn.
func do_dfltWarn(call, op, args, env Value) Value {
	runtimeState.mu.Lock()
	runtimeState.lastError = formatValue(arg(args, 0))
	runtimeState.mu.Unlock()
	return Nil
}
