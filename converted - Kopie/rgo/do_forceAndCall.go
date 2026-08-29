package rgo

// do_forceAndCall translates batch_101.c (src/main/eval.c); R primitive(s): forceAndCall.
func do_forceAndCall(call, op, args, env Value) Value {
	f := arg(args, 1)
	if f.Kind != Function || f.Fn == nil {
		return ErrValue("attempt to apply non-function")
	}
	av := arg(args, 2)
	if av.Kind != List {
		av = Lists(av)
	}
	v, e := f.Fn(av.V, env.Env)
	if e != nil {
		return ErrValue("%v", e)
	}
	return v
}
