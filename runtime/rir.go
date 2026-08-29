package runtime

// RIR (R Intermediate Representation) is the executable form of a primitive
// matrix row.  It intentionally has a tiny, fixed instruction alphabet: the
// instruction stream is assembled from the lowering and algorithm vectors,
// while all state lives in LoweringFrame.  Thus generators can emit data rather
// than one Go function per R primitive.
type RIROp uint8

const (
	RIREvaluate RIROp = iota + 1
	RIRArguments
	RIRTypes
	RIRShape
	RIRAttributes
	RIRDispatch
	RIRDynamic
	RIRBackend
)

const rirWidth = 8

// rirStageIncidence maps the 55 lowering micro-operation columns to the eight
// executable RIR stages.  Program assembly is the integer product
// `recipe.MicroOps × rirStageIncidence`; this is the point where the matrix
// determines executable control flow rather than merely a report column.
var rirStageIncidence = [55][rirWidth]uint8{
	// evaluation / promises
	{1}, {1}, {1}, {1},
	// argument matching
	{0, 1}, {0, 1}, {0, 1}, {0, 1}, {0, 1},
	// type and coercion
	{0, 0, 1}, {0, 0, 1}, {0, 0, 1}, {0, 0, 1}, {0, 0, 1}, {0, 0, 1}, {0, 0, 1},
	// shape / recycling
	{0, 0, 0, 1}, {0, 0, 0, 1}, {0, 0, 0, 1}, {0, 0, 0, 1}, {0, 0, 0, 1},
	// attributes
	{0, 0, 0, 0, 1}, {0, 0, 0, 0, 1}, {0, 0, 0, 0, 1}, {0, 0, 0, 0, 1},
	{0, 0, 0, 0, 1}, {0, 0, 0, 0, 1}, {0, 0, 0, 0, 1}, {0, 0, 0, 0, 1},
	// dispatch
	{0, 0, 0, 0, 0, 1}, {0, 0, 0, 0, 0, 1}, {0, 0, 0, 0, 0, 1}, {0, 0, 0, 0, 0, 1},
	{0, 0, 0, 0, 0, 1}, {0, 0, 0, 0, 0, 1},
	// dynamic effects
	{0, 0, 0, 0, 0, 0, 1}, {0, 0, 0, 0, 0, 0, 1}, {0, 0, 0, 0, 0, 0, 1}, {0, 0, 0, 0, 0, 0, 1},
	{0, 0, 0, 0, 0, 0, 1}, {0, 0, 0, 0, 0, 0, 1}, {0, 0, 0, 0, 0, 0, 1}, {0, 0, 0, 0, 0, 0, 1},
	// backend mechanisms
	{0, 0, 0, 0, 0, 0, 0, 1}, {0, 0, 0, 0, 0, 0, 0, 1}, {0, 0, 0, 0, 0, 0, 0, 1},
	{0, 0, 0, 0, 0, 0, 0, 1}, {0, 0, 0, 0, 0, 0, 0, 1}, {0, 0, 0, 0, 0, 0, 0, 1},
	{0, 0, 0, 0, 0, 0, 0, 1}, {0, 0, 0, 0, 0, 0, 0, 1}, {0, 0, 0, 0, 0, 0, 0, 1},
	{0, 0, 0, 0, 0, 0, 0, 1}, {0, 0, 0, 0, 0, 0, 0, 1}, {0, 0, 0, 0, 0, 0, 0, 1},
}

// RIRProgram is a fixed-size operation vector.  Zero entries are inert; this
// permits a complete primitive x operation matrix without variable control
// flow or generated per-primitive Go code.
type RIRProgram struct {
	PlanIndex int
	Name      string
	Ops       [rirWidth]RIROp
}

// RIRProgramForPrimitive materializes a program from the three generated
// source matrices.  The lowering row supplies the semantic route, and the
// GNU-R algorithm row is required as a provenance guard.  The latter becomes
// operational as its individual C operations gain RIR opcodes; it is not
// treated as proof that a coordinate algorithm is implemented.
func RIRProgramForPrimitive(name string) (RIRProgram, bool) {
	plan, ok := ExecutionPlanByName[name]
	if !ok {
		return RIRProgram{}, false
	}
	recipe, ok := LoweringRecipeByName[name]
	if !ok {
		return RIRProgram{}, false
	}
	if _, ok := AlgorithmRecipeByName[name]; !ok {
		return RIRProgram{}, false
	}
	var product [rirWidth]uint16
	for micro, enabled := range recipe.MicroOps {
		if enabled == 0 {
			continue
		}
		for stage, coefficient := range rirStageIncidence[micro] {
			product[stage] += uint16(enabled) * uint16(coefficient)
		}
	}
	var ops [rirWidth]RIROp
	for stage, coefficient := range product {
		if coefficient != 0 {
			ops[stage] = RIROp(stage + 1)
		}
	}
	// Every R primitive has an algorithm coordinate.  A backend instruction is
	// therefore mandatory even when the structural row contains no backend bit.
	ops[rirWidth-1] = RIRBackend
	return RIRProgram{PlanIndex: plan.Index, Name: name, Ops: ops}, true
}

// RIRProgramCoverage is deliberately structural: it establishes that every
// R_FunTab primitive has a data-driven execution program, not algorithmic
// equivalence to GNU R.
func RIRProgramCoverage() (covered, total int) {
	total = len(ExecutionPlanByName)
	for name := range ExecutionPlanByName {
		if _, ok := RIRProgramForPrimitive(name); ok {
			covered++
		}
	}
	return covered, total
}
