package cir

import "fmt"

// Host supplies the runtime-specific meaning of C calls and storage. CIR stays
// independent of R's SEXP representation and can therefore be tested without
// an embedded R runtime.
type Host interface {
	Call(name string) error
	Store(name string) error
	Branch(kind string) error
}

// Execute interprets structural CIR instructions. Later expression-lowering
// passes attach operands and branch targets; this VM already guarantees that
// every source-level call/control/effect is routed through one Pure-Go host.
func Execute(program Program, host Host) error {
	for pc := 0; pc < len(program.Instructions); pc++ {
		instruction := program.Instructions[pc]
		switch instruction.Op {
		case Call:
			if err := host.Call(instruction.A); err != nil {
				return fmt.Errorf("%s call %s: %w", program.Name, instruction.A, err)
			}
		case Store:
			if err := host.Store(instruction.A); err != nil {
				return fmt.Errorf("%s store: %w", program.Name, err)
			}
		case Branch:
			if err := host.Branch(instruction.A); err != nil {
				return fmt.Errorf("%s branch %s: %w", program.Name, instruction.A, err)
			}
		case Return:
			return nil
		}
	}
	return nil
}
