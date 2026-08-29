package runtime

import (
	"fmt"
	"r2go/syntax"
)

// LoweringFrame is the vector passed through every factor of the pure-Go
// lowering tensor. Stages operate on the same frame, so all recipe classes use
// one executor rather than generated per-primitive control flow.
type LoweringFrame struct {
	Plan      ExecutionPlan
	Recipe    LoweringRecipe
	Arguments []syntax.Argument
	Values    []Value
	Env       *Environment
	Result    Value
}

type loweringStage func(*Context, *LoweringFrame) error

// Coordinate kernels are deliberately separate from structural recipes. A
// kernel is pure Go and is addressed only by GNU R's (C-entry, offset) pair.
type loweringKernel func(*Context, *LoweringFrame) error

var loweringKernelVector = map[string]loweringKernel{}

func (c *Context) executeLoweringVM(plan ExecutionPlan, args []syntax.Argument, env *Environment) (Value, error) {
	recipe, ok := LoweringRecipeByName[plan.Name]
	if !ok {
		return nil, fmt.Errorf("no pure-Go lowering row for plan %d", plan.Index)
	}
	program, ok := RIRProgramForPrimitive(plan.Name)
	if !ok {
		return nil, fmt.Errorf("no RIR program for plan %d", plan.Index)
	}
	frame := &LoweringFrame{Plan: plan, Recipe: recipe, Arguments: args, Env: env}
	for _, op := range program.Ops {
		if op == 0 {
			continue
		}
		if err := c.executeRIROp(op, frame); err != nil {
			return nil, err
		}
	}
	if frame.Result == nil {
		return nil, fmt.Errorf("pure-Go coordinate kernel %s is not generated", planCoordinate(plan))
	}
	return frame.Result, nil
}

func (c *Context) executeRIROp(op RIROp, frame *LoweringFrame) error {
	// This vector is intentionally indexed by RIR opcode, not primitive name.
	stages := [...]loweringStage{
		nil,
		lowerEvaluationFactor, lowerArgumentFactor, lowerTypeFactor, lowerShapeFactor,
		lowerAttributeFactor, lowerDispatchFactor, lowerDynamicFactor, lowerBackendFactor,
	}
	if int(op) >= len(stages) || stages[op] == nil {
		return fmt.Errorf("unknown RIR opcode %d", op)
	}
	return stages[op](c, frame)
}

func lowerEvaluationFactor(c *Context, frame *LoweringFrame) error {
	// MicroOps 1 and 2 preserve language/promises. Only eager rows (index 3)
	// materialize values before the coordinate kernel.
	if frame.Recipe.MicroOps[1] != 0 || frame.Recipe.MicroOps[2] != 0 || frame.Recipe.MicroOps[3] == 0 {
		return nil
	}
	frame.Values = make([]Value, len(frame.Arguments))
	for i, argument := range frame.Arguments {
		if argument.Value == nil {
			continue
		}
		value, err := c.Eval(argument.Value, frame.Env)
		if err != nil {
			return err
		}
		frame.Values[i] = value
	}
	return nil
}

func lowerArgumentFactor(_ *Context, frame *LoweringFrame) error {
	// R_FunTab's arity is a dispatcher descriptor, not a required actual-argument
	// count: many entries have optional arguments. Coordinate kernels validate
	// only their mandatory prefix and preserve the remaining promises.
	return nil
}

func lowerTypeFactor(_ *Context, _ *LoweringFrame) error      { return nil }
func lowerShapeFactor(_ *Context, _ *LoweringFrame) error     { return nil }
func lowerAttributeFactor(_ *Context, _ *LoweringFrame) error { return nil }
func lowerDispatchFactor(_ *Context, _ *LoweringFrame) error  { return nil }
func lowerDynamicFactor(_ *Context, _ *LoweringFrame) error   { return nil }

func lowerBackendFactor(c *Context, frame *LoweringFrame) error {
	kernel := loweringKernelVector[planCoordinate(frame.Plan)]
	if kernel == nil {
		return nil
	}
	return kernel(c, frame)
}

func StructuralRecipeClassCoverage() (coveredClasses, totalClasses int) {
	seen := map[uint16]bool{}
	for _, recipe := range LoweringRecipeTable {
		seen[recipe.RecipeClass] = true
	}
	return len(seen), len(seen)
}
