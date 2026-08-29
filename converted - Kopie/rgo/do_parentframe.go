package rgo

// do_parentframe translates batch_184.c (src/main/context.c); R primitive(s): parent.frame.
func do_parentframe(call, op, args, env Value) Value {
	if env.Kind == Environment && env.Env != nil && env.Env.Parent != nil {
		return Value{Kind: Environment, Env: env.Env.Parent}
	}
	return Nil
}
