package runtime

import (
	"fmt"
	"path/filepath"
	"r2go/syntax"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

// The algorithm vector selects these kernels by GNU R dispatcher coordinate,
// never by R surface spelling. Aliases sharing a C entry use one Go kernel.
func init() {
	registerLoweringKernel("do_seq_len", "0", kernelSeqLength)
	registerLoweringKernel("do_rep_int", "0", kernelRepeat)
	registerLoweringKernel("do_rep_len", "0", kernelRepeat)
	registerLoweringKernel("do_invisible", "0", kernelIdentity)
	registerLoweringKernel("do_identical", "0", kernelIdentical)
	registerLoweringKernel("do_makevector", "0", kernelMakeVector)
	registerLoweringKernel("do_nzchar", "1", kernelNZChar)
	registerLoweringKernel("do_unlist", "0", kernelUnlist)
	registerLoweringKernel("do_lengths", "0", kernelLengths)
	registerLoweringKernel("do_tabulate", "0", kernelTabulate)
	registerLoweringKernel("do_findinterval", "0", kernelFindInterval)
	registerLoweringKernel("do_assign", "0", kernelAssign)
	registerLoweringKernel("do_get", "0", kernelEnvironmentLookup)
	registerLoweringKernel("do_get", "1", kernelEnvironmentLookup)
	registerLoweringKernel("do_get", "2", kernelEnvironmentLookup)
	registerLoweringKernel("do_newenv", "0", kernelNewEnvironment)
	registerLoweringKernel("do_parentenv", "0", kernelParentEnvironment)
	registerLoweringKernel("do_envir", "0", kernelEnvironmentOf)
	registerLoweringKernel("do_emptyenv", "0", kernelEmptyEnvironment)
	registerLoweringKernel("do_baseenv", "0", kernelBaseEnvironment)
	registerLoweringKernel("do_globalenv", "0", kernelGlobalEnvironment)
	registerLoweringKernel("do_getOption", "0", kernelGetOption)
	registerLoweringKernel("do_gsub", "0", kernelTextFamily)
	registerLoweringKernel("do_gsub", "1", kernelTextFamily)
	registerLoweringKernel("do_strsplit", "1", kernelTextFamily)
	registerLoweringKernel("do_grep", "0", kernelTextFamily)
	registerLoweringKernel("do_grep", "1", kernelTextFamily)
	registerLoweringKernel("do_internal", "0", kernelInternal)
	registerLoweringKernel("do_switch", "0", kernelSwitch)
	registerLoweringKernel("do_options", "0", kernelOptions)
	registerLoweringKernel("do_sprintf", "1", kernelRuntimeTextFamily)
	registerLoweringKernel("do_cat", "0", kernelRuntimeTextFamily)
	registerLoweringKernel("do_vapply", "0", kernelVApply)
	registerLoweringKernel("do_format", "0", kernelFormat)
	registerLoweringKernel("do_deparse", "0", kernelDeparse)
	registerLoweringKernel("do_filepath", "0", kernelFilePath)
	registerLoweringKernel("do_colon2", "0", kernelNamespaceLookup)
	registerLoweringKernel("do_colon3", "0", kernelNamespaceLookup)
	registerLoweringKernel("do_getenv", "0", kernelEnvironmentService)
	registerLoweringKernel("do_setenv", "0", kernelEnvironmentService)
	registerLoweringKernel("do_unsetenv", "0", kernelEnvironmentService)
	registerLoweringKernel("do_gettext", "0", kernelTranslationService)
	registerLoweringKernel("do_ngettext", "0", kernelTranslationService)
}

// kernelEnvironmentService is the shared vector for the Sys.* environment
// primitives.  Values are deliberately evaluated once by the RIR evaluation
// factor and are recycled exactly along the name vector where applicable.
func kernelEnvironmentService(c *Context, frame *LoweringFrame) error {
	namesValue, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	names, err := characterData(namesValue)
	if err != nil {
		return err
	}
	switch frame.Plan.CEntry {
	case "do_getenv":
		unset := ""
		if len(frame.Arguments) > 1 {
			if value, err := frameValue(c, frame, 1); err == nil {
				unset = scalarText(value)
			}
		}
		out := &CharacterVector{Data: make([]string, len(names.Data)), Missing: append([]bool(nil), names.Missing...)}
		for i, name := range names.Data {
			if i < len(names.Missing) && names.Missing[i] {
				continue
			}
			out.Data[i] = c.Host.Getenv(name)
			if out.Data[i] == "" {
				out.Data[i] = unset
			}
		}
		frame.Result = out
		return nil
	case "do_unsetenv":
		for i, name := range names.Data {
			if i >= len(names.Missing) || !names.Missing[i] {
				_ = c.Host.Unsetenv(name)
			}
		}
	case "do_setenv":
		valuesValue, err := frameValue(c, frame, 1)
		if err != nil {
			return err
		}
		values, err := characterData(valuesValue)
		if err != nil || len(values.Data) == 0 {
			return fmt.Errorf("Sys.setenv: value must be a non-empty character vector")
		}
		for i, name := range names.Data {
			if i < len(names.Missing) && names.Missing[i] {
				continue
			}
			if err := c.Host.Setenv(name, values.Data[i%len(values.Data)]); err != nil {
				return err
			}
		}
	}
	frame.Result = &LogicalVector{Data: []Logical{True}}
	return nil
}

// GNU R obtains translations from its message catalog.  The pure-Go baseline
// has no external R catalog, so its deterministic default locale returns the
// supplied singular/plural strings while preserving ngettext's selector.
func kernelTranslationService(c *Context, frame *LoweringFrame) error {
	if frame.Plan.CEntry == "do_gettext" {
		message, err := frameValue(c, frame, 0)
		if err != nil {
			return err
		}
		return setTranslationResult(message, frame)
	}
	n, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	count, err := scalarInt(n)
	if err != nil {
		return err
	}
	index := 1
	if count != 1 {
		index = 2
	}
	message, err := frameValue(c, frame, index)
	if err != nil {
		return err
	}
	return setTranslationResult(message, frame)
}

func setTranslationResult(v Value, frame *LoweringFrame) error {
	text, err := characterData(v)
	if err != nil {
		return err
	}
	frame.Result = &CharacterVector{Data: append([]string(nil), text.Data...), Missing: append([]bool(nil), text.Missing...)}
	return nil
}

func registerLoweringKernel(entry, offset string, kernel loweringKernel) {
	loweringKernelVector[coordinate(entry, offset)] = kernel
}

func loweringKernelAvailable(plan ExecutionPlan) bool {
	return loweringKernelVector[planCoordinate(plan)] != nil
}

// CoordinateKernelAvailable reports whether the pure-Go algorithm coefficient
// for a primitive's generated matrix row is registered.
func CoordinateKernelAvailable(name string) bool {
	plan, ok := ExecutionPlanByName[name]
	return ok && loweringKernelAvailable(plan)
}

func frameValue(c *Context, frame *LoweringFrame, index int) (Value, error) {
	if index >= len(frame.Arguments) || frame.Arguments[index].Value == nil {
		return nil, fmt.Errorf("%s: missing argument %d", frame.Plan.Name, index+1)
	}
	if index < len(frame.Values) && frame.Values[index] != nil {
		return frame.Values[index], nil
	}
	return c.Eval(frame.Arguments[index].Value, frame.Env)
}

func kernelSeqLength(c *Context, frame *LoweringFrame) error {
	v, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	n, err := scalarInt(v)
	if err != nil || n < 0 {
		return fmt.Errorf("invalid sequence length")
	}
	out := &IntegerVector{Data: make([]int64, n)}
	for i := range out.Data {
		out.Data[i] = int64(i + 1)
	}
	frame.Result = out
	return nil
}

func kernelRepeat(c *Context, frame *LoweringFrame) error {
	value, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	countValue, err := frameValue(c, frame, 1)
	if err != nil {
		return err
	}
	counts, numberErr := numbers(countValue)
	if numberErr != nil || len(counts.Data) == 0 {
		return fmt.Errorf("invalid repeat count")
	}
	count := int(counts.Data[0])
	if count < 0 {
		return fmt.Errorf("invalid repeat count")
	}
	if frame.Plan.CEntry == "do_rep_len" {
		if Length(value) == 0 && count > 0 {
			return fmt.Errorf("cannot repeat empty value to positive length")
		}
		positions := make([]int, count)
		for i := range positions {
			positions[i] = i % Length(value)
		}
		frame.Result = takePositions(value, positions)
		return nil
	}
	if len(counts.Data) > 1 {
		if len(counts.Data) != Length(value) {
			return fmt.Errorf("invalid repeat count")
		}
		positions := make([]int, 0)
		for i, itemCount := range counts.Data {
			if missingAt(counts, i) || itemCount < 0 || itemCount != float64(int(itemCount)) {
				return fmt.Errorf("invalid repeat count")
			}
			for repeat := 0; repeat < int(itemCount); repeat++ {
				positions = append(positions, i)
			}
		}
		frame.Result = takePositions(value, positions)
		return nil
	}
	positions := make([]int, 0, Length(value)*count)
	for repeat := 0; repeat < count; repeat++ {
		for i := 0; i < Length(value); i++ {
			positions = append(positions, i)
		}
	}
	frame.Result = takePositions(value, positions)
	return nil
}

func kernelIdentity(c *Context, frame *LoweringFrame) error {
	v, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	frame.Result = v
	return nil
}

func kernelIdentical(c *Context, frame *LoweringFrame) error {
	a, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	b, err := frameValue(c, frame, 1)
	if err != nil {
		return err
	}
	frame.Result = boolValue(reflect.DeepEqual(a, b))
	return nil
}

func kernelMakeVector(c *Context, frame *LoweringFrame) error {
	modeValue, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	lengthValue, err := frameValue(c, frame, 1)
	if err != nil {
		return err
	}
	mode, ok := modeValue.(*CharacterVector)
	if !ok || len(mode.Data) != 1 {
		return fmt.Errorf("vector mode must be one string")
	}
	n, err := scalarInt(lengthValue)
	if err != nil || n < 0 {
		return fmt.Errorf("invalid vector length")
	}
	switch mode.Data[0] {
	case "logical":
		frame.Result = &LogicalVector{Data: make([]Logical, n)}
	case "integer":
		frame.Result = &IntegerVector{Data: make([]int64, n)}
	case "double", "numeric":
		frame.Result = &DoubleVector{Data: make([]float64, n)}
	case "complex":
		frame.Result = &ComplexVector{Data: make([]complex128, n)}
	case "character":
		frame.Result = &CharacterVector{Data: make([]string, n)}
	case "raw":
		frame.Result = &RawVector{Data: make([]byte, n)}
	case "list", "expression":
		out := &List{Data: make([]Value, n)}
		for i := range out.Data {
			out.Data[i] = NullValue
		}
		frame.Result = out
	default:
		return fmt.Errorf("unsupported vector mode %q", mode.Data[0])
	}
	return nil
}

func kernelNZChar(c *Context, frame *LoweringFrame) error {
	v, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	characters, err := CoerceTo(v, CharacterKind)
	if err != nil {
		return err
	}
	x := characters.(*CharacterVector)
	out := &LogicalVector{Data: make([]Logical, len(x.Data))}
	for i, text := range x.Data {
		if !(i < len(x.Missing) && x.Missing[i]) && text != "" {
			out.Data[i] = True
		}
	}
	frame.Result = out
	return nil
}

func kernelUnlist(c *Context, frame *LoweringFrame) error {
	v, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	flat := []Value{}
	var visit func(Value)
	visit = func(item Value) {
		if list, ok := item.(*List); ok {
			for _, child := range list.Data {
				visit(child)
			}
			return
		}
		flat = append(flat, item)
	}
	visit(v)
	result, err := combine(flat)
	if err != nil {
		return err
	}
	frame.Result = result
	return nil
}

func kernelLengths(c *Context, frame *LoweringFrame) error {
	v, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	out := &IntegerVector{Data: make([]int64, Length(v))}
	for i, item := range elements(v) {
		out.Data[i] = int64(Length(item))
	}
	frame.Result = out
	return nil
}

func kernelTabulate(c *Context, frame *LoweringFrame) error {
	v, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	x, err := numbers(v)
	if err != nil {
		return err
	}
	maxBin := 0
	for i, item := range x.Data {
		if !missingAt(x, i) && int(item) > maxBin {
			maxBin = int(item)
		}
	}
	if len(frame.Arguments) > 1 {
		n, e := frameValue(c, frame, 1)
		if e != nil {
			return e
		}
		maxBin, e = scalarInt(n)
		if e != nil {
			return e
		}
	}
	if maxBin < 0 {
		return fmt.Errorf("invalid 'nbins'")
	}
	out := &IntegerVector{Data: make([]int64, maxBin)}
	for i, item := range x.Data {
		bin := int(item)
		if !missingAt(x, i) && bin >= 1 && bin <= maxBin {
			out.Data[bin-1]++
		}
	}
	frame.Result = out
	return nil
}

func frameEnvironment(c *Context, frame *LoweringFrame, index int) (*Environment, error) {
	if index >= len(frame.Arguments) || frame.Arguments[index].Value == nil {
		return frame.Env, nil
	}
	v, err := frameValue(c, frame, index)
	if err != nil {
		return nil, err
	}
	if environment, ok := v.(*EnvironmentValue); ok {
		return environment.Env, nil
	}
	return nil, fmt.Errorf("environment argument is not an environment")
}

func frameText(c *Context, frame *LoweringFrame, index int) (string, error) {
	v, err := frameValue(c, frame, index)
	if err != nil {
		return "", err
	}
	x, ok := v.(*CharacterVector)
	if !ok || len(x.Data) != 1 || (len(x.Missing) > 0 && x.Missing[0]) {
		return "", fmt.Errorf("name must be one string")
	}
	return x.Data[0], nil
}

func kernelAssign(c *Context, frame *LoweringFrame) error {
	name, err := frameText(c, frame, 0)
	if err != nil {
		return err
	}
	value, err := frameValue(c, frame, 1)
	if err != nil {
		return err
	}
	environment, err := frameEnvironment(c, frame, 2)
	if err != nil {
		return err
	}
	environment.Set(name, value)
	frame.Result = value
	return nil
}

func kernelEnvironmentLookup(c *Context, frame *LoweringFrame) error {
	name, err := frameText(c, frame, 0)
	if err != nil {
		return err
	}
	environment, err := frameEnvironment(c, frame, 1)
	if err != nil {
		return err
	}
	_, binding, found := environment.Find(name)
	if frame.Plan.Offset == "0" {
		frame.Result = boolValue(found)
		return nil
	}
	if !found {
		if frame.Plan.Offset == "2" {
			frame.Result = NullValue
			return nil
		}
		return fmt.Errorf("object %q not found", name)
	}
	value, err := binding.Force(c)
	if err != nil {
		return err
	}
	frame.Result = value
	return nil
}

func kernelNewEnvironment(c *Context, frame *LoweringFrame) error {
	parent := frame.Env
	for i, argument := range frame.Arguments {
		if argument.Name == "parent" {
			var err error
			parent, err = frameEnvironment(c, frame, i)
			if err != nil {
				return err
			}
		}
	}
	frame.Result = &EnvironmentValue{Env: NewEnvironment(parent)}
	return nil
}

func kernelParentEnvironment(c *Context, frame *LoweringFrame) error {
	v, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	environment, ok := v.(*EnvironmentValue)
	if !ok {
		return fmt.Errorf("parent.env expects an environment")
	}
	if environment.Env.Parent == nil {
		frame.Result = &EnvironmentValue{Env: NewEnvironment(nil), Name: "emptyenv"}
	} else {
		frame.Result = &EnvironmentValue{Env: environment.Env.Parent}
	}
	return nil
}

func kernelEnvironmentOf(c *Context, frame *LoweringFrame) error {
	v, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	if closure, ok := v.(*Closure); ok {
		frame.Result = &EnvironmentValue{Env: closure.Env}
		return nil
	}
	if environment, ok := v.(*EnvironmentValue); ok {
		frame.Result = environment
		return nil
	}
	return fmt.Errorf("environment expects a closure or environment")
}

func kernelEmptyEnvironment(_ *Context, frame *LoweringFrame) error {
	frame.Result = &EnvironmentValue{Env: NewEnvironment(nil), Name: "emptyenv"}
	return nil
}
func kernelBaseEnvironment(c *Context, frame *LoweringFrame) error {
	frame.Result = &EnvironmentValue{Env: c.Global, Name: "base"}
	return nil
}
func kernelGlobalEnvironment(c *Context, frame *LoweringFrame) error {
	frame.Result = &EnvironmentValue{Env: c.Global, Name: "global"}
	return nil
}

func kernelGetOption(c *Context, frame *LoweringFrame) error {
	name, err := frameText(c, frame, 0)
	if err != nil {
		return err
	}
	if value, ok := c.Options[name]; ok {
		frame.Result = value
	} else {
		frame.Result = NullValue
	}
	return nil
}

func frameNamed(c *Context, frame *LoweringFrame, name string) (Value, bool, error) {
	for i, argument := range frame.Arguments {
		if argument.Name == name {
			value, err := frameValue(c, frame, i)
			return value, true, err
		}
	}
	return nil, false, nil
}

func frameFlag(c *Context, frame *LoweringFrame, name string) (bool, error) {
	v, found, err := frameNamed(c, frame, name)
	if err != nil || !found {
		return false, err
	}
	return IsTrue(v)
}

func characterData(v Value) (*CharacterVector, error) {
	text, err := CoerceTo(v, CharacterKind)
	if err != nil {
		return nil, err
	}
	return text.(*CharacterVector), nil
}

func kernelSubstituteRegex(c *Context, frame *LoweringFrame) error {
	patternValue, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	replacementValue, err := frameValue(c, frame, 1)
	if err != nil {
		return err
	}
	textValue, err := frameValue(c, frame, 2)
	if err != nil {
		return err
	}
	patterns, err := characterData(patternValue)
	if err != nil {
		return err
	}
	replacements, err := characterData(replacementValue)
	if err != nil {
		return err
	}
	texts, err := characterData(textValue)
	if err != nil {
		return err
	}
	if len(patterns.Data) == 0 || len(replacements.Data) == 0 || len(texts.Data) == 0 {
		frame.Result = &CharacterVector{}
		return nil
	}
	fixed, err := frameFlag(c, frame, "fixed")
	if err != nil {
		return err
	}
	ignoreCase, err := frameFlag(c, frame, "ignore.case")
	if err != nil {
		return err
	}
	global := frame.Plan.Offset == "1"
	out := &CharacterVector{Data: make([]string, len(texts.Data)), Missing: make([]bool, len(texts.Data))}
	for i, text := range texts.Data {
		if i < len(texts.Missing) && texts.Missing[i] {
			out.Missing[i] = true
			continue
		}
		pattern := patterns.Data[i%len(patterns.Data)]
		replacement := replacements.Data[i%len(replacements.Data)]
		if fixed {
			pattern = regexp.QuoteMeta(pattern)
		}
		if ignoreCase {
			pattern = "(?i)" + pattern
		}
		re, compileErr := regexp.Compile(pattern)
		if compileErr != nil {
			return compileErr
		}
		if global {
			out.Data[i] = re.ReplaceAllString(text, replacement)
		} else {
			out.Data[i] = re.ReplaceAllStringFunc(text, func(match string) string { return replacement })
		}
		if !global {
			location := re.FindStringIndex(text)
			if location != nil {
				out.Data[i] = text[:location[0]] + replacement + text[location[1]:]
			}
		}
	}
	frame.Result = out
	return nil
}

func kernelSplitRegex(c *Context, frame *LoweringFrame) error {
	textValue, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	splitValue, err := frameValue(c, frame, 1)
	if err != nil {
		return err
	}
	texts, err := characterData(textValue)
	if err != nil {
		return err
	}
	splits, err := characterData(splitValue)
	if err != nil {
		return err
	}
	if len(splits.Data) == 0 {
		return fmt.Errorf("split must not be empty")
	}
	fixed, err := frameFlag(c, frame, "fixed")
	if err != nil {
		return err
	}
	out := &List{Data: make([]Value, len(texts.Data)), Names: append([]string(nil), valueNames(texts)...)}
	for i, text := range texts.Data {
		if i < len(texts.Missing) && texts.Missing[i] {
			out.Data[i] = &CharacterVector{Missing: []bool{true}, Data: []string{""}}
			continue
		}
		separator := splits.Data[i%len(splits.Data)]
		if fixed {
			separator = regexp.QuoteMeta(separator)
		}
		re, compileErr := regexp.Compile(separator)
		if compileErr != nil {
			return compileErr
		}
		out.Data[i] = &CharacterVector{Data: re.Split(text, -1)}
	}
	frame.Result = out
	return nil
}

// All regex/text primitives share the same vectorized pure-Go family kernel.
// C-entry and offset are coefficients in the generated plan, not individual
// evaluator branches in the dispatcher.
func kernelTextFamily(c *Context, frame *LoweringFrame) error {
	switch frame.Plan.CEntry {
	case "do_gsub":
		return kernelSubstituteRegex(c, frame)
	case "do_strsplit":
		return kernelSplitRegex(c, frame)
	case "do_grep":
		return kernelGrep(c, frame)
	}
	return fmt.Errorf("text coordinate %s is not implemented", planCoordinate(frame.Plan))
}

// .Internal is an R language boundary, not an external runtime call. The
// nested language object is evaluated by the same pure-Go plan dispatcher.
func kernelInternal(c *Context, frame *LoweringFrame) error {
	if len(frame.Arguments) != 1 || frame.Arguments[0].Value == nil {
		return fmt.Errorf(".Internal expects one language call")
	}
	value, err := c.Eval(frame.Arguments[0].Value, frame.Env)
	if err != nil {
		return err
	}
	frame.Result = value
	return nil
}

func kernelSwitch(c *Context, frame *LoweringFrame) error {
	if len(frame.Arguments) == 0 {
		return fmt.Errorf("switch expects an expression")
	}
	selector, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	selected := -1
	if text, ok := selector.(*CharacterVector); ok && len(text.Data) == 1 && !(len(text.Missing) > 0 && text.Missing[0]) {
		for i, argument := range frame.Arguments[1:] {
			if argument.Name == text.Data[0] {
				selected = i + 1
				break
			}
		}
	} else if index, indexErr := scalarInt(selector); indexErr == nil && index > 0 && index < len(frame.Arguments) {
		selected = index
	}
	if selected < 0 {
		for i, argument := range frame.Arguments[1:] {
			if argument.Name == "" && argument.Value != nil {
				selected = i + 1
			}
		}
	}
	if selected < 0 || frame.Arguments[selected].Value == nil {
		frame.Result = NullValue
		return nil
	}
	value, err := c.Eval(frame.Arguments[selected].Value, frame.Env)
	if err != nil {
		return err
	}
	frame.Result = value
	return nil
}

func kernelOptions(c *Context, frame *LoweringFrame) error {
	if len(frame.Arguments) == 0 {
		names := make([]string, 0, len(c.Options))
		for name := range c.Options {
			names = append(names, name)
		}
		sort.Strings(names)
		out := &List{Names: names, Data: make([]Value, len(names))}
		for i, name := range names {
			out.Data[i] = c.Options[name]
		}
		frame.Result = out
		return nil
	}
	previous := &List{}
	for i, argument := range frame.Arguments {
		if argument.Name == "" {
			continue
		}
		value, err := frameValue(c, frame, i)
		if err != nil {
			return err
		}
		previous.Names = append(previous.Names, argument.Name)
		if old, found := c.Options[argument.Name]; found {
			previous.Data = append(previous.Data, old)
		} else {
			previous.Data = append(previous.Data, NullValue)
		}
		c.Options[argument.Name] = value
	}
	frame.Result = previous
	return nil
}

func kernelRuntimeTextFamily(c *Context, frame *LoweringFrame) error {
	switch frame.Plan.CEntry {
	case "do_sprintf":
		return kernelSprintf(c, frame)
	case "do_cat":
		return kernelCat(c, frame)
	}
	return fmt.Errorf("runtime text coordinate %s is not implemented", planCoordinate(frame.Plan))
}

func scalarInterface(value Value) any {
	switch x := value.(type) {
	case *CharacterVector:
		if len(x.Data) > 0 {
			return x.Data[0]
		}
	case *IntegerVector:
		if len(x.Data) > 0 {
			return x.Data[0]
		}
	case *DoubleVector:
		if len(x.Data) > 0 {
			if x.Data[0] == float64(int64(x.Data[0])) {
				return int64(x.Data[0])
			}
			return x.Data[0]
		}
	case *LogicalVector:
		if len(x.Data) > 0 {
			return x.Data[0] == True
		}
	case *ComplexVector:
		if len(x.Data) > 0 {
			return x.Data[0]
		}
	}
	return value.String()
}

func kernelFindInterval(c *Context, frame *LoweringFrame) error {
	if len(frame.Arguments) < 2 {
		return fmt.Errorf("findInterval expects x and vec")
	}
	xv, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	vecv, err := frameValue(c, frame, 1)
	if err != nil {
		return err
	}
	x, err := numbers(xv)
	if err != nil {
		return err
	}
	vec, err := numbers(vecv)
	if err != nil {
		return err
	}
	out := &IntegerVector{Data: make([]int64, len(x.Data)), Missing: append([]bool(nil), x.Missing...)}
	for i, value := range x.Data {
		if missingAt(x, i) {
			continue
		}
		// default R interval convention: vec[j] <= x < vec[j+1].
		out.Data[i] = int64(sort.Search(len(vec.Data), func(j int) bool { return vec.Data[j] > value }))
	}
	frame.Result = out
	return nil
}

func kernelSprintf(c *Context, frame *LoweringFrame) error {
	if len(frame.Arguments) == 0 {
		return fmt.Errorf("sprintf expects a format")
	}
	formatValue, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	formats, err := characterData(formatValue)
	if err != nil {
		return err
	}
	if len(formats.Data) == 0 {
		frame.Result = &CharacterVector{}
		return nil
	}
	values := make([]Value, len(frame.Arguments)-1)
	for i := range values {
		value, e := frameValue(c, frame, i+1)
		if e != nil {
			return e
		}
		values[i] = value
	}
	length := len(formats.Data)
	for _, value := range values {
		if Length(value) > length {
			length = Length(value)
		}
	}
	out := &CharacterVector{Data: make([]string, length), Missing: make([]bool, length)}
	for i := 0; i < length; i++ {
		args := make([]any, len(values))
		for j, value := range values {
			pieces := elements(value)
			if len(pieces) == 0 {
				args[j] = ""
			} else {
				args[j] = scalarInterface(pieces[i%len(pieces)])
			}
		}
		out.Data[i] = fmt.Sprintf(formats.Data[i%len(formats.Data)], args...)
	}
	frame.Result = out
	return nil
}

func kernelCat(c *Context, frame *LoweringFrame) error {
	sep := " "
	if value, found, err := frameNamed(c, frame, "sep"); err != nil {
		return err
	} else if found {
		if text, ok := value.(*CharacterVector); ok && len(text.Data) > 0 {
			sep = text.Data[0]
		}
	}
	parts := []string{}
	for i, argument := range frame.Arguments {
		if argument.Name == "sep" || argument.Name == "file" || argument.Name == "fill" || argument.Name == "labels" || argument.Name == "append" {
			continue
		}
		value, err := frameValue(c, frame, i)
		if err != nil {
			return err
		}
		for _, piece := range elements(value) {
			parts = append(parts, scalarText(piece))
		}
	}
	if c.Output != nil {
		_, _ = fmt.Fprint(c.Output, strings.Join(parts, sep))
	}
	frame.Result = NullValue
	return nil
}

func kernelVApply(c *Context, frame *LoweringFrame) error {
	if len(frame.Arguments) < 3 {
		return fmt.Errorf("vapply expects X, FUN and FUN.VALUE")
	}
	x, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	function, err := frameValue(c, frame, 1)
	if err != nil {
		return err
	}
	template, err := frameValue(c, frame, 2)
	if err != nil {
		return err
	}
	closure, ok := function.(*Closure)
	if !ok {
		return fmt.Errorf("vapply FUN must be a closure")
	}
	results := make([]Value, 0, Length(x))
	expected := Length(template)
	for _, item := range elements(x) {
		value, callErr := c.callClosureWithValue(closure, item)
		if callErr != nil {
			return callErr
		}
		if Length(value) != expected {
			return fmt.Errorf("vapply result has length %d, expected %d", Length(value), expected)
		}
		coerced, coerceErr := CoerceTo(value, template.Kind())
		if coerceErr != nil {
			return coerceErr
		}
		results = append(results, coerced)
	}
	result, err := combine(results)
	if err != nil {
		return err
	}
	frame.Result = result
	return nil
}

func kernelFormat(c *Context, frame *LoweringFrame) error {
	v, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	out, err := CoerceTo(v, CharacterKind)
	if err != nil {
		return err
	}
	frame.Result = out
	return nil
}
func kernelDeparse(c *Context, frame *LoweringFrame) error {
	v, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	frame.Result = &CharacterVector{Data: []string{v.String()}}
	return nil
}

func kernelFilePath(c *Context, frame *LoweringFrame) error {
	parts := [][]string{}
	length := 1
	for i, argument := range frame.Arguments {
		if argument.Name == "fsep" {
			continue
		}
		value, err := frameValue(c, frame, i)
		if err != nil {
			return err
		}
		text, err := characterData(value)
		if err != nil {
			return err
		}
		if len(text.Data) == 0 {
			frame.Result = &CharacterVector{}
			return nil
		}
		parts = append(parts, text.Data)
		if len(text.Data) > length {
			length = len(text.Data)
		}
	}
	out := &CharacterVector{Data: make([]string, length)}
	for i := 0; i < length; i++ {
		fields := make([]string, len(parts))
		for j, part := range parts {
			fields[j] = part[i%len(part)]
		}
		out.Data[i] = filepath.Join(fields...)
	}
	frame.Result = out
	return nil
}

func kernelNamespaceLookup(c *Context, frame *LoweringFrame) error {
	if len(frame.Arguments) < 2 {
		return fmt.Errorf("namespace lookup expects package and name")
	}
	name := ""
	if symbol, ok := frame.Arguments[1].Value.(*syntax.Symbol); ok {
		name = symbol.Name
	} else {
		value, err := frameValue(c, frame, 1)
		if err != nil {
			return err
		}
		text, ok := value.(*CharacterVector)
		if !ok || len(text.Data) != 1 {
			return fmt.Errorf("namespace name must be a symbol or string")
		}
		name = text.Data[0]
	}
	value, err := c.Global.Get(c, name)
	if err != nil {
		return err
	}
	frame.Result = value
	return nil
}

func kernelGrep(c *Context, frame *LoweringFrame) error {
	patternValue, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	textValue, err := frameValue(c, frame, 1)
	if err != nil {
		return err
	}
	patterns, err := characterData(patternValue)
	if err != nil {
		return err
	}
	texts, err := characterData(textValue)
	if err != nil {
		return err
	}
	if len(patterns.Data) != 1 {
		return fmt.Errorf("grep pattern must have length one")
	}
	fixed, err := frameFlag(c, frame, "fixed")
	if err != nil {
		return err
	}
	ignoreCase, err := frameFlag(c, frame, "ignore.case")
	if err != nil {
		return err
	}
	invert, err := frameFlag(c, frame, "invert")
	if err != nil {
		return err
	}
	pattern := patterns.Data[0]
	if fixed {
		pattern = regexp.QuoteMeta(pattern)
	}
	if ignoreCase {
		pattern = "(?i)" + pattern
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	if frame.Plan.Offset == "1" {
		out := &LogicalVector{Data: make([]Logical, len(texts.Data))}
		for i, text := range texts.Data {
			match := re.MatchString(text)
			if invert {
				match = !match
			}
			if match {
				out.Data[i] = True
			}
		}
		frame.Result = out
		return nil
	}
	value, err := frameFlag(c, frame, "value")
	if err != nil {
		return err
	}
	positions := []int{}
	matches := []string{}
	for i, text := range texts.Data {
		match := re.MatchString(text)
		if invert {
			match = !match
		}
		if match {
			positions = append(positions, i)
			matches = append(matches, text)
		}
	}
	if value {
		frame.Result = &CharacterVector{Data: matches}
	} else {
		out := &IntegerVector{Data: make([]int64, len(positions))}
		for i, position := range positions {
			out.Data[i] = int64(position + 1)
		}
		frame.Result = out
	}
	return nil
}
