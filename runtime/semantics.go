package runtime

import (
	"fmt"
	"math"
	"strconv"
)

// RLength and RIndex deliberately model GNU R's R_xlen_t-style long-vector
// semantics. Existing APIs may still expose int while the runtime migrates.
type RLength int64
type RIndex int64

// Atomic coercion order mirrors the central GNU R vector promotion hierarchy.
// List is included as the universal recursive container target.
var coercionRank = map[Kind]int{
	RawKind:       0,
	LogicalKind:   1,
	IntegerKind:   2,
	DoubleKind:    3,
	ComplexKind:   4,
	CharacterKind: 5,
	ListKind:      6,
}

func IsVectorKind(k Kind) bool { _, ok := coercionRank[k]; return ok }

func CommonKind(a, b Kind) (Kind, error) {
	ra, oka := coercionRank[a]
	rb, okb := coercionRank[b]
	if !oka || !okb {
		return "", fmt.Errorf("no vector coercion rule for %s and %s", a, b)
	}
	if rb > ra {
		return b, nil
	}
	return a, nil
}

func CommonKindValues(vals []Value) (Kind, error) {
	if len(vals) == 0 {
		return LogicalKind, nil
	}
	k := vals[0].Kind()
	for _, v := range vals[1:] {
		var err error
		k, err = CommonKind(k, v.Kind())
		if err != nil {
			return "", err
		}
	}
	return k, nil
}

func Length64(v Value) RLength { return RLength(Length(v)) }

// CoerceTo is the single atomic-vector conversion gateway. Generic runtime
// operations use it instead of maintaining private conversion switches.
func CoerceTo(v Value, target Kind) (Value, error) {
	if v.Kind() == target {
		return cloneValue(v), nil
	}
	if target == ListKind {
		return &List{Data: append([]Value(nil), elements(v)...)}, nil
	}
	if !IsVectorKind(v.Kind()) || !IsVectorKind(target) {
		return nil, fmt.Errorf("cannot coerce %s to %s", v.Kind(), target)
	}
	src := elements(v)
	switch target {
	case CharacterKind:
		out := &CharacterVector{Data: make([]string, len(src)), Missing: make([]bool, len(src))}
		for i, e := range src {
			if scalarMissing(e) {
				out.Missing[i] = true
				continue
			}
			switch x := e.(type) {
			case *RawVector:
				out.Data[i] = fmt.Sprintf("%02x", x.Data[0])
			case *LogicalVector:
				if x.Data[0] == True {
					out.Data[i] = "TRUE"
				} else {
					out.Data[i] = "FALSE"
				}
			case *IntegerVector:
				out.Data[i] = strconv.FormatInt(x.Data[0], 10)
			case *DoubleVector:
				if math.IsNaN(x.Data[0]) {
					out.Data[i] = "NaN"
				} else {
					out.Data[i] = strconv.FormatFloat(x.Data[0], 'g', -1, 64)
				}
			case *ComplexVector:
				out.Data[i] = x.String()
			case *CharacterVector:
				out.Data[i] = x.Data[0]
			default:
				return nil, fmt.Errorf("cannot coerce %s to character", e.Kind())
			}
		}
		return out, nil
	case ComplexKind:
		out := &ComplexVector{Data: make([]complex128, len(src)), Missing: make([]bool, len(src))}
		for i, e := range src {
			if scalarMissing(e) {
				out.Missing[i] = true
				continue
			}
			switch x := e.(type) {
			case *RawVector:
				out.Data[i] = complex(float64(x.Data[0]), 0)
			case *LogicalVector:
				out.Data[i] = complex(float64(x.Data[0]), 0)
			case *IntegerVector:
				out.Data[i] = complex(float64(x.Data[0]), 0)
			case *DoubleVector:
				out.Data[i] = complex(x.Data[0], 0)
			case *ComplexVector:
				out.Data[i] = x.Data[0]
			default:
				return nil, fmt.Errorf("cannot coerce %s to complex", e.Kind())
			}
		}
		return out, nil
	case DoubleKind:
		out := &DoubleVector{Data: make([]float64, len(src)), Missing: make([]bool, len(src))}
		for i, e := range src {
			if scalarMissing(e) {
				out.Missing[i] = true
				out.Data[i] = NAReal()
				continue
			}
			switch x := e.(type) {
			case *RawVector:
				out.Data[i] = float64(x.Data[0])
			case *LogicalVector:
				out.Data[i] = float64(x.Data[0])
			case *IntegerVector:
				out.Data[i] = float64(x.Data[0])
			case *DoubleVector:
				out.Data[i] = x.Data[0]
			case *CharacterVector:
				parsed, err := strconv.ParseFloat(x.Data[0], 64)
				if err != nil {
					out.Missing[i] = true
					out.Data[i] = NAReal()
				} else {
					out.Data[i] = parsed
				}
			default:
				return nil, fmt.Errorf("cannot coerce %s to double", e.Kind())
			}
		}
		return out, nil
	case IntegerKind:
		out := &IntegerVector{Data: make([]int64, len(src)), Missing: make([]bool, len(src))}
		for i, e := range src {
			if scalarMissing(e) {
				out.Missing[i] = true
				continue
			}
			switch x := e.(type) {
			case *RawVector:
				out.Data[i] = int64(x.Data[0])
			case *LogicalVector:
				out.Data[i] = int64(x.Data[0])
			case *IntegerVector:
				out.Data[i] = x.Data[0]
			default:
				return nil, fmt.Errorf("cannot coerce %s to integer without narrowing", e.Kind())
			}
		}
		return out, nil
	case LogicalKind:
		out := &LogicalVector{Data: make([]Logical, len(src))}
		for i, e := range src {
			if scalarMissing(e) {
				out.Data[i] = NA
				continue
			}
			switch x := e.(type) {
			case *RawVector:
				if x.Data[0] != 0 {
					out.Data[i] = True
				}
			case *LogicalVector:
				out.Data[i] = x.Data[0]
			default:
				return nil, fmt.Errorf("cannot coerce %s to logical without narrowing", e.Kind())
			}
		}
		return out, nil
	case RawKind:
		out := &RawVector{Data: make([]byte, len(src))}
		for i, e := range src {
			x, ok := e.(*RawVector)
			if !ok || len(x.Data) == 0 {
				return nil, fmt.Errorf("cannot coerce %s to raw", e.Kind())
			}
			out.Data[i] = x.Data[0]
		}
		return out, nil
	}
	return nil, fmt.Errorf("unsupported coercion target %s", target)
}

// MatchKey models equality used by match/unique hashing. It intentionally
// distinguishes NA from ordinary NaN and preserves character missingness.
func MatchKey(v Value) string {
	if scalarMissing(v) {
		return string(v.Kind()) + ":NA"
	}
	switch x := v.(type) {
	case *RawVector:
		return "raw:" + strconv.Itoa(int(x.Data[0]))
	case *LogicalVector:
		return "logical:" + strconv.Itoa(int(x.Data[0]))
	case *IntegerVector:
		return "number:" + strconv.FormatInt(x.Data[0], 10)
	case *DoubleVector:
		if math.IsNaN(x.Data[0]) {
			return "number:NaN"
		}
		return "number:" + strconv.FormatFloat(x.Data[0], 'g', -1, 64)
	case *ComplexVector:
		z := x.Data[0]
		return "complex:" + strconv.FormatFloat(real(z), 'g', -1, 64) + ":" + strconv.FormatFloat(imag(z), 'g', -1, 64)
	case *CharacterVector:
		return "character:" + x.Data[0]
	default:
		return string(v.Kind()) + ":" + v.String()
	}
}

func coercePairForMatching(a, b Value) (Value, Value, error) {
	k, err := CommonKind(a.Kind(), b.Kind())
	if err != nil {
		return a, b, nil
	} // recursive/non-vector matching keeps identity-like keys for now
	aa, err := CoerceTo(a, k)
	if err != nil {
		return nil, nil, err
	}
	bb, err := CoerceTo(b, k)
	if err != nil {
		return nil, nil, err
	}
	return aa, bb, nil
}
