package rgo

// do_topenv translates batch_268.c (src/main/envir.c); R primitive(s): topenv.
func do_topenv(call, op, args, env Value) Value {
	for env.Kind == Environment && env.Env != nil && env.Env.Parent != nil {
		env = Value{Kind: Environment, Env: env.Env.Parent}
	}
	return env
}
