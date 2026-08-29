package runtime

// AlgorithmRecipe is generated from the transitive GNU R C call graph. Its
// operation vector describes which pure-Go mechanisms a primitive's complete
// reachable implementation requires; it contains no R runtime dependency.
type AlgorithmRecipe struct {
	PlanIndex        int
	Name             string
	CEntry           string
	ReachableHelpers uint16
	Operations       [48]uint8
}

var AlgorithmRecipeByName = buildAlgorithmRecipeIndex()

func buildAlgorithmRecipeIndex() map[string]AlgorithmRecipe {
	result := make(map[string]AlgorithmRecipe, len(AlgorithmRecipeTable))
	for _, recipe := range AlgorithmRecipeTable {
		result[recipe.Name] = recipe
	}
	return result
}

func AlgorithmRecipeCoverage() (covered, total int) {
	return len(AlgorithmRecipeByName), len(ExecutionPlanByName)
}
