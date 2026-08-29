package runtime

// ExecutionPlan is the product of the primitive-kernel incidence matrix and
// the orthogonal semantic-function matrix. Arrays are fixed-width vectors so
// the executor can compose policies without name-based dispatch.
type ExecutionPlan struct {
	Index                                                  int
	Name, CEntry, Offset, Kernel, Opcode                   string
	Arity                                                  int
	Evaluation                                             [3]uint8
	Types                                                  [6]uint8
	Shape                                                  [8]uint8
	Attributes                                             [8]uint8
	Dispatch                                               [4]uint8
	Effects                                                [6]uint8
	Backend                                                [5]uint8
	ReplacementPolicy, AttributeOperation, MatrixOperation string
}

var ExecutionPlanByName = buildExecutionPlanIndex()

func buildExecutionPlanIndex() map[string]ExecutionPlan {
	result := make(map[string]ExecutionPlan, len(ExecutionPlanTable))
	for _, plan := range ExecutionPlanTable {
		result[plan.Name] = plan
	}
	return result
}
