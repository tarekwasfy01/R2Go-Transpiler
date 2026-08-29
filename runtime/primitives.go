package runtime

// PrimitiveDescriptor is the lossless Go representation of one GNU R
// R_FunTab row.  It keeps parser, evaluation and dispatcher metadata together
// so all later implementation and coverage decisions can be matrix-driven.
type PrimitiveDescriptor struct {
	Index              int
	Name               string
	CEntry             string
	SourceFile         string
	Offset             string
	Eval               int
	Arity              int
	PPKind             string
	Precedence         string
	RightAssociative   bool
	SEXPType           string
	EvaluatesArguments bool
	Kernel             string
	EvaluationPolicy   string
	ArgumentPolicy     string
	CoercionPolicy     string
	VectorPolicy       string
	IndexPolicy        string
	AttributePolicy    string
	DispatchPolicy     string
	EffectPolicy       string
	ReplacementPolicy  string
	AttributeOperation string
	MatrixOperation    string
}

type PrimitiveFamily struct {
	CEntry string
	Names  []string
}

var (
	PrimitiveByName   map[string]PrimitiveDescriptor
	PrimitiveByEntry  map[string][]PrimitiveDescriptor
	PrimitiveFamilies []PrimitiveFamily
)

func init() {
	PrimitiveByName = make(map[string]PrimitiveDescriptor, len(PrimitiveTable))
	PrimitiveByEntry = make(map[string][]PrimitiveDescriptor)
	for _, descriptor := range PrimitiveTable {
		PrimitiveByName[descriptor.Name] = descriptor
		PrimitiveByEntry[descriptor.CEntry] = append(PrimitiveByEntry[descriptor.CEntry], descriptor)
	}
	PrimitiveFamilies = make([]PrimitiveFamily, 0, len(PrimitiveByEntry))
	for entry, descriptors := range PrimitiveByEntry {
		family := PrimitiveFamily{CEntry: entry, Names: make([]string, len(descriptors))}
		for i := range descriptors {
			family.Names[i] = descriptors[i].Name
		}
		PrimitiveFamilies = append(PrimitiveFamilies, family)
	}
}

// PrimitiveStatus separates GNU-R discovery from R2Go implementation.  Every
// GNU-R primitive is therefore reportable even before its dispatcher family
// has a Go implementation.
type PrimitiveStatus struct {
	PrimitiveDescriptor
	Implemented bool
}

func PrimitiveCoverageMatrix() []PrimitiveStatus {
	result := make([]PrimitiveStatus, len(PrimitiveTable))
	for i, descriptor := range PrimitiveTable {
		result[i] = PrimitiveStatus{PrimitiveDescriptor: descriptor, Implemented: PrimitiveExecutable(descriptor.Name)}
	}
	return result
}
