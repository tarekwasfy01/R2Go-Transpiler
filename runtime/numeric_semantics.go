package runtime

import "math"

// GNU R distinguishes NA_real_ from ordinary NaN with a stable payload.
// Keeping the bit pattern available is necessary for serialization and a
// future R C-API facade even though vectors also carry an explicit bitmap.
const naRealBits uint64 = 0x7ff00000000007a2

func NAReal() float64              { return math.Float64frombits(naRealBits) }
func IsNAReal(x float64) bool      { return math.IsNaN(x) && math.Float64bits(x)&0xffffffff == 1954 }
func IsNaNButNotNA(x float64) bool { return math.IsNaN(x) && !IsNAReal(x) }

// RoundZero follows R's fround(x, 0) ties-to-even behavior.
func RoundZero(x float64) float64 { return math.RoundToEven(x) }
