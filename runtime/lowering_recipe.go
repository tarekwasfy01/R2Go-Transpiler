package runtime

// LoweringRecipe is a fixed-width row in the pure-Go lowering tensor. The
// micro-operation vector is generated for every GNU R primitive, including
// primitives whose algorithm kernel is not implemented yet.
type LoweringRecipe struct {
	PlanIndex   int
	Name        string
	RecipeClass uint16
	Route       string
	Wave        string
	MicroOps    [55]uint8
}

var LoweringRecipeByName = buildLoweringRecipeIndex()

func buildLoweringRecipeIndex() map[string]LoweringRecipe {
	result := make(map[string]LoweringRecipe, len(LoweringRecipeTable))
	for _, recipe := range LoweringRecipeTable {
		result[recipe.Name] = recipe
	}
	return result
}

// StructuralLoweringCoverage reports whether every primitive has a complete
// pure-Go execution-structure row. It intentionally does not claim that every
// coordinate-specific algorithm kernel already exists.
func StructuralLoweringCoverage() (covered, total int) {
	return len(LoweringRecipeByName), len(ExecutionPlanByName)
}
