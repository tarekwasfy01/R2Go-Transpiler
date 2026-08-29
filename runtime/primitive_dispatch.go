package runtime

import (
	"fmt"
	"math"
	converted "r2go/converted/rgo"
	rprimitive "r2go/converted/rprimitive"
	"r2go/syntax"
	"sort"
	"strings"
)

var executableOpcodes = map[string]bool{
	"VECTOR_UNARY": true, "NUMERIC_BINARY": true, "BITWISE": true,
	"VECTOR_SCAN": true, "TYPE_PREDICATE": true, "OPS_BINARY": true,
	"SUBSET": true, "COERCE": true, "MISSINGNESS": true,
	"NUMERIC_PREDICATE": true,
	"VECTOR_REDUCE":     true,
}

var executableAttributeOperations = map[string]bool{"get-key": true, "get-map": true, "get-named": true, "unclass": true}
var executableReplacementPolicies = map[string]bool{"subscript": true, "member": true, "attribute-key": true, "attribute-map": true, "attribute-named": true}
var executableMatrixOperations = map[string]bool{"construct-matrix": true, "construct-array": true, "length": true, "transpose": true, "coordinates": true, "drop": true, "product": true, "diagonal": true, "array-runtime": true}

func planExecutable(plan ExecutionPlan) bool {
	if loweringKernelAvailable(plan) {
		return true
	}
	if executableOpcodes[plan.Opcode] {
		return functionVectorAvailable(plan)
	}
	if plan.Opcode == "ATTRIBUTE" {
		return executableAttributeOperations[plan.AttributeOperation]
	}
	if plan.Opcode == "REPLACE" {
		return executableReplacementPolicies[plan.ReplacementPolicy]
	}
	if plan.Opcode == "MATRIX" {
		return executableMatrixOperations[plan.MatrixOperation]
	}
	return false
}

// ExecutionPlanExecutable exposes the calculated executor-capability vector to
// matrix/report generators. Runtime dispatch still uses the same predicate, so
// readiness reports cannot drift away from what generated programs execute.
func ExecutionPlanExecutable(name string) bool {
	plan, ok := ExecutionPlanByName[name]
	return ok && planExecutable(plan)
}

// PrimitiveExecutable is the single capability truth for matrices: either the
// generated plan path or an existing evaluator kernel can execute the entry.
func PrimitiveExecutable(name string) bool {
	if ExecutionPlanExecutable(name) || ImplementedNames[name] {
		return true
	}
	if manualPrimitiveRoutes[name] {
		return true
	}
	plan, ok := ExecutionPlanByName[name]
	if !ok {
		return false
	}
	// A translated control-flow entry is executable through the DynamicRuntime
	// bridge. Its unsupported native operations produce a typed Pure-Go error;
	// they never fall through to a hidden R or cgo runtime.
	if _, translated := rprimitive.PrimitiveTable[name]; translated {
		return true
	}
	return converted.Eligible(plan.CEntry)
}

func (c *Context) dispatchPrimitive(name string, args []syntax.Argument, env *Environment) (Value, bool, error) {
	plan, known := ExecutionPlanByName[name]
	if !known {
		return nil, false, nil
	}
	if !planExecutable(plan) {
		if value, handled, err := c.executeManualPrimitive(name, args, env); handled {
			return value, true, err
		}
		// Imported C-to-Go entries are reached only when the existing native
		// Pure-Go kernel matrix has no stronger implementation for the primitive.
		if value, handled, err := c.executeConvertedEntry(plan, args, env); handled {
			return value, true, err
		}
		if value, handled, err := c.executeTranslatedEntry(plan, args, env); handled {
			return value, true, err
		}
		// Legacy evaluator paths remain authoritative where present. Every other
		// primitive is passed through the generated 8-factor pure-Go matrix VM;
		// missing coordinate kernels are reported separately from structure.
		if !ImplementedNames[name] {
			if _, structured := LoweringRecipeByName[name]; structured {
				value, err := c.executeLoweringVM(plan, args, env)
				return value, true, err
			}
		}
		return nil, false, nil
	}
	var v Value
	var err error
	// A coordinate-specific family kernel is more precise than the broad
	// opcode class. This is essential for entries such as lengths, whose GNU R
	// table metadata says MATRIX while their semantics are a dedicated vector
	// operation.
	if loweringKernelAvailable(plan) {
		v, err = c.executeLoweringVM(plan, args, env)
		return v, true, err
	}
	switch plan.Opcode {
	case "VECTOR_UNARY":
		v, err = c.executeUnaryVector(plan, args, env)
	case "NUMERIC_BINARY":
		v, err = c.executeBinaryVector(plan, args, env)
	case "BITWISE":
		v, err = c.executeBitwiseVector(plan, args, env)
	case "VECTOR_SCAN":
		v, err = c.executeScanVector(plan, args, env)
	case "TYPE_PREDICATE":
		v, err = c.executeTypePredicate(plan, args, env)
	case "OPS_BINARY":
		v, err = c.operator(operatorFunctionVector[planCoordinate(plan)], args, env)
	case "VECTOR_REDUCE":
		v, err = c.executeReduceVector(plan, args, env)
	case "SUBSET":
		subsetOperator := "["
		if plan.CEntry == "do_subset2" || plan.CEntry == "do_subset2_dflt" {
			subsetOperator = "[["
		}
		if plan.CEntry == "do_subset3" {
			v, err = c.dollar(args, env)
		} else {
			v, err = c.subset(subsetOperator, args, env)
		}
	case "COERCE":
		v, err = c.builtin(name, args, env)
	case "MISSINGNESS":
		v, err = c.matrixMissingness(plan, args, env)
	case "NUMERIC_PREDICATE":
		v, err = c.matrixNumericPredicate(plan, args, env)
	case "ATTRIBUTE":
		v, err = c.matrixAttribute(plan, args, env)
	case "MATRIX":
		v, err = c.matrixKernel(plan, args, env)
	default:
		if loweringKernelAvailable(plan) {
			v, err = c.executeLoweringVM(plan, args, env)
		} else {
			return nil, false, nil
		}
	}
	return v, true, err
}

func (c *Context) matrixAttribute(descriptor ExecutionPlan, args []syntax.Argument, env *Environment) (Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("%s expects an object", descriptor.Name)
	}
	object, err := c.Eval(args[0].Value, env)
	if err != nil {
		return nil, err
	}
	switch descriptor.AttributeOperation {
	case "get-key":
		key := descriptor.Name
		if key == "oldClass" {
			if value, ok := Attributes(object)["class"]; ok {
				return value, nil
			}
			return NullValue, nil
		}
		if key == "names" {
			names := valueNames(object)
			if names == nil {
				return NullValue, nil
			}
			return &CharacterVector{Data: append([]string(nil), names...)}, nil
		}
		if value, ok := Attributes(object)[key]; ok {
			return value, nil
		}
		if key == "class" {
			return &CharacterVector{Data: []string{defaultClass(object)}}, nil
		}
		return NullValue, nil
	case "get-map":
		attrs := Attributes(object)
		if len(attrs) == 0 {
			return NullValue, nil
		}
		names := make([]string, 0, len(attrs))
		for name := range attrs {
			names = append(names, name)
		}
		sort.Strings(names)
		out := &List{Names: names, Data: make([]Value, len(names))}
		for i, name := range names {
			out.Data[i] = attrs[name]
		}
		return out, nil
	case "get-named":
		if len(args) < 2 {
			return nil, fmt.Errorf("attr expects object and name")
		}
		nameValue, err := c.Eval(args[1].Value, env)
		if err != nil {
			return nil, err
		}
		name, ok := nameValue.(*CharacterVector)
		if !ok || len(name.Data) != 1 {
			return nil, fmt.Errorf("attribute name must be one string")
		}
		if value, ok := Attributes(object)[name.Data[0]]; ok {
			return value, nil
		}
		return NullValue, nil
	case "unclass":
		out := cloneValue(object)
		if err := setAttribute(out, "class", NullValue); err != nil {
			return nil, err
		}
		return out, nil
	}
	return nil, fmt.Errorf("attribute operation %s is not executable", descriptor.AttributeOperation)
}

func (c *Context) matrixMissingness(descriptor ExecutionPlan, args []syntax.Argument, env *Environment) (Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("%s expects an argument", descriptor.Name)
	}
	v, err := c.Eval(args[0].Value, env)
	if err != nil {
		return nil, err
	}
	flags := make([]Logical, Length(v))
	for i, element := range elements(v) {
		missing := scalarMissing(element)
		switch x := element.(type) {
		case *DoubleVector:
			missing = missing || math.IsNaN(x.Data[0])
		case *ComplexVector:
			missing = missing || math.IsNaN(real(x.Data[0])) || math.IsNaN(imag(x.Data[0]))
		}
		if missing {
			flags[i] = True
		}
	}
	if descriptor.CEntry == "do_anyNA" {
		for _, flag := range flags {
			if flag == True {
				return boolValue(true), nil
			}
		}
		return boolValue(false), nil
	}
	return &LogicalVector{Data: flags}, nil
}

func (c *Context) matrixNumericPredicate(descriptor ExecutionPlan, args []syntax.Argument, env *Environment) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("%s expects one argument", descriptor.Name)
	}
	v, err := c.Eval(args[0].Value, env)
	if err != nil {
		return nil, err
	}
	out := &LogicalVector{Data: make([]Logical, Length(v))}
	op := strings.TrimPrefix(descriptor.CEntry, "do_is")
	for i, element := range elements(v) {
		result := false
		switch x := element.(type) {
		case *RawVector, *LogicalVector, *IntegerVector:
			result = op == "finite"
		case *DoubleVector:
			n := x.Data[0]
			if !scalarMissing(element) {
				switch op {
				case "nan":
					result = math.IsNaN(n)
				case "finite":
					result = !math.IsNaN(n) && !math.IsInf(n, 0)
				case "infinite":
					result = math.IsInf(n, 0)
				}
			}
		case *ComplexVector:
			z := x.Data[0]
			if !scalarMissing(element) {
				switch op {
				case "nan":
					result = math.IsNaN(real(z)) || math.IsNaN(imag(z))
				case "finite":
					result = !math.IsNaN(real(z)) && !math.IsNaN(imag(z)) && !math.IsInf(real(z), 0) && !math.IsInf(imag(z), 0)
				case "infinite":
					result = math.IsInf(real(z), 0) || math.IsInf(imag(z), 0)
				}
			}
		}
		if result {
			out.Data[i] = True
		}
	}
	return out, nil
}
