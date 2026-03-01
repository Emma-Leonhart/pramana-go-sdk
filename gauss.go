package pramana

import (
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/google/uuid"
)

// PramanaNamespace is the shared UUID v5 namespace for all Pramana entities.
var PramanaNamespace = uuid.MustParse("a6613321-e9f6-4348-8f8b-29d2a3c86349")

// Gauss represents a Gaussian rational number: a/b + (c/d)i
// where a, b, c, d are arbitrary-precision integers and fractions
// are always normalized to lowest terms with positive denominators.
type Gauss struct {
	a *big.Int // real numerator
	b *big.Int // real denominator (always positive)
	c *big.Int // imaginary numerator
	d *big.Int // imaginary denominator (always positive)
}

// normalize reduces fractions to lowest terms with positive denominators.
func normalize(num, den *big.Int) (*big.Int, *big.Int) {
	n := copyBigInt(num)
	dd := copyBigInt(den)

	if dd.Sign() == 0 {
		panic("pramana: division by zero in fraction")
	}

	// Zero numerator normalizes to 0/1
	if n.Sign() == 0 {
		return big.NewInt(0), big.NewInt(1)
	}

	// Make denominator positive
	if dd.Sign() < 0 {
		n.Neg(n)
		dd.Neg(dd)
	}

	// Reduce by GCD
	g := gcd(n, dd)
	n.Div(n, g)
	dd.Div(dd, g)

	return n, dd
}

// NewGauss creates a new Gaussian rational a/b + (c/d)i.
// The fraction is automatically normalized.
func NewGauss(a, b, c, d int64) *Gauss {
	return NewGaussBig(big.NewInt(a), big.NewInt(b), big.NewInt(c), big.NewInt(d))
}

// NewGaussBig creates a new Gaussian rational from big.Int values.
func NewGaussBig(a, b, c, d *big.Int) *Gauss {
	na, nb := normalize(a, b)
	nc, nd := normalize(c, d)
	return &Gauss{a: na, b: nb, c: nc, d: nd}
}

// NewGaussInt creates a Gaussian rational from two integers (real + imag*i).
func NewGaussInt(real, imag int64) *Gauss {
	return NewGauss(real, 1, imag, 1)
}

// NewGaussReal creates a Gaussian rational from a single real integer.
func NewGaussReal(real int64) *Gauss {
	return NewGauss(real, 1, 0, 1)
}

// NewGaussFromBigInts creates a Gaussian rational from two big.Int values (real + imag*i).
func NewGaussFromBigInts(real, imag *big.Int) *Gauss {
	return NewGaussBig(real, big.NewInt(1), imag, big.NewInt(1))
}

// GaussZero returns the zero value.
func GaussZero() *Gauss { return NewGaussReal(0) }

// GaussOne returns the value 1.
func GaussOne() *Gauss { return NewGaussReal(1) }

// GaussMinusOne returns the value -1.
func GaussMinusOne() *Gauss { return NewGaussReal(-1) }

// GaussI returns the imaginary unit i.
func GaussI() *Gauss { return NewGaussInt(0, 1) }

// Eye returns the imaginary unit i (alias for GaussI).
func GaussEye() *Gauss { return GaussI() }

// GaussUnits returns the four Gaussian units: [1, -1, i, -i].
func GaussUnits() []*Gauss {
	return []*Gauss{GaussOne(), GaussMinusOne(), GaussI(), NewGaussInt(0, -1)}
}

// --- Component Accessors ---

// A returns the real numerator.
func (g *Gauss) A() *big.Int { return copyBigInt(g.a) }

// B returns the real denominator (always positive).
func (g *Gauss) B() *big.Int { return copyBigInt(g.b) }

// C returns the imaginary numerator.
func (g *Gauss) C() *big.Int { return copyBigInt(g.c) }

// D returns the imaginary denominator (always positive).
func (g *Gauss) D() *big.Int { return copyBigInt(g.d) }

// --- Properties ---

// IsReal returns true if the imaginary part is zero.
func (g *Gauss) IsReal() bool {
	return g.c.Sign() == 0
}

// IsPurelyImaginary returns true if the real part is zero and imaginary is nonzero.
func (g *Gauss) IsPurelyImaginary() bool {
	return g.a.Sign() == 0 && g.c.Sign() != 0
}

// IsZero returns true if both parts are zero.
func (g *Gauss) IsZero() bool {
	return g.a.Sign() == 0 && g.c.Sign() == 0
}

// IsOne returns true if the value equals 1.
func (g *Gauss) IsOne() bool {
	return g.a.Cmp(big.NewInt(1)) == 0 && g.b.Cmp(big.NewInt(1)) == 0 && g.c.Sign() == 0
}

// IsInteger returns true if the value is a real integer (b=1, c=0).
func (g *Gauss) IsInteger() bool {
	return g.b.Cmp(big.NewInt(1)) == 0 && g.c.Sign() == 0
}

// IsGaussianInteger returns true if both denominators are 1.
func (g *Gauss) IsGaussianInteger() bool {
	return g.b.Cmp(big.NewInt(1)) == 0 && g.d.Cmp(big.NewInt(1)) == 0
}

// IsNegative returns true if the value is real and negative.
func (g *Gauss) IsNegative() bool {
	return g.IsReal() && g.a.Sign() < 0
}

// IsPositive returns true if the value is real and positive.
func (g *Gauss) IsPositive() bool {
	return g.IsReal() && g.a.Sign() > 0
}

// --- Derived Values ---

// Conjugate returns the complex conjugate (a/b - c/d*i).
func (g *Gauss) Conjugate() *Gauss {
	negC := new(big.Int).Neg(g.c)
	return NewGaussBig(copyBigInt(g.a), copyBigInt(g.b), negC, copyBigInt(g.d))
}

// RealPart returns the real part as a Gauss value.
func (g *Gauss) RealPart() *Gauss {
	return NewGaussBig(copyBigInt(g.a), copyBigInt(g.b), big.NewInt(0), big.NewInt(1))
}

// ImaginaryPart returns the imaginary coefficient as a Gauss value.
func (g *Gauss) ImaginaryPart() *Gauss {
	return NewGaussBig(copyBigInt(g.c), copyBigInt(g.d), big.NewInt(0), big.NewInt(1))
}

// Norm returns the squared magnitude |z|² = (a/b)² + (c/d)² as a Gauss value.
func (g *Gauss) Norm() *Gauss {
	// (a/b)² = a²/b²
	a2 := new(big.Int).Mul(g.a, g.a)
	b2 := new(big.Int).Mul(g.b, g.b)
	// (c/d)² = c²/d²
	c2 := new(big.Int).Mul(g.c, g.c)
	d2 := new(big.Int).Mul(g.d, g.d)
	// a²/b² + c²/d² = (a²*d² + c²*b²) / (b²*d²)
	num1 := new(big.Int).Mul(a2, d2)
	num2 := new(big.Int).Mul(c2, b2)
	num := new(big.Int).Add(num1, num2)
	den := new(big.Int).Mul(b2, d2)
	return NewGaussBig(num, den, big.NewInt(0), big.NewInt(1))
}

// MagnitudeSquared is an alias for Norm.
func (g *Gauss) MagnitudeSquared() *Gauss {
	return g.Norm()
}

// Magnitude returns |z| as a float64 approximation.
func (g *Gauss) Magnitude() float64 {
	re := g.RealFloat64()
	im := g.ImagFloat64()
	return math.Sqrt(re*re + im*im)
}

// Phase returns the argument (angle in radians) as a float64.
func (g *Gauss) Phase() float64 {
	return math.Atan2(g.ImagFloat64(), g.RealFloat64())
}

// Reciprocal returns 1/z = conjugate / |z|².
func (g *Gauss) Reciprocal() *Gauss {
	return GaussOne().Div(g)
}

// Inverse is an alias for Reciprocal.
func (g *Gauss) Inverse() *Gauss {
	return g.Reciprocal()
}

// --- Arithmetic ---

// Add returns g + other.
func (g *Gauss) Add(other *Gauss) *Gauss {
	// Real: a/b + e/f = (a*f + e*b) / (b*f)
	realNum := new(big.Int).Add(
		new(big.Int).Mul(g.a, other.b),
		new(big.Int).Mul(other.a, g.b),
	)
	realDen := new(big.Int).Mul(g.b, other.b)
	// Imag: c/d + g/h = (c*h + g*d) / (d*h)
	imagNum := new(big.Int).Add(
		new(big.Int).Mul(g.c, other.d),
		new(big.Int).Mul(other.c, g.d),
	)
	imagDen := new(big.Int).Mul(g.d, other.d)
	return NewGaussBig(realNum, realDen, imagNum, imagDen)
}

// Sub returns g - other.
func (g *Gauss) Sub(other *Gauss) *Gauss {
	return g.Add(other.Neg())
}

// Neg returns -g.
func (g *Gauss) Neg() *Gauss {
	return NewGaussBig(
		new(big.Int).Neg(g.a), copyBigInt(g.b),
		new(big.Int).Neg(g.c), copyBigInt(g.d),
	)
}

// Mul returns g * other.
func (g *Gauss) Mul(other *Gauss) *Gauss {
	// (a/b + c/d*i) * (e/f + g/h*i)
	// Real = a*e/(b*f) - c*g/(d*h)
	// Imag = a*g/(b*h) + c*e/(d*f)
	ae := new(big.Int).Mul(g.a, other.a)
	bf := new(big.Int).Mul(g.b, other.b)
	cg := new(big.Int).Mul(g.c, other.c)
	dh := new(big.Int).Mul(g.d, other.d)
	ag := new(big.Int).Mul(g.a, other.c)
	bh := new(big.Int).Mul(g.b, other.d)
	ce := new(big.Int).Mul(g.c, other.a)
	df := new(big.Int).Mul(g.d, other.b)

	// Real: ae/bf - cg/dh = (ae*dh - cg*bf) / (bf*dh)
	realNum := new(big.Int).Sub(
		new(big.Int).Mul(ae, dh),
		new(big.Int).Mul(cg, bf),
	)
	realDen := new(big.Int).Mul(bf, dh)

	// Imag: ag/bh + ce/df = (ag*df + ce*bh) / (bh*df)
	imagNum := new(big.Int).Add(
		new(big.Int).Mul(ag, df),
		new(big.Int).Mul(ce, bh),
	)
	imagDen := new(big.Int).Mul(bh, df)

	return NewGaussBig(realNum, realDen, imagNum, imagDen)
}

// Div returns g / other. Panics if other is zero.
func (g *Gauss) Div(other *Gauss) *Gauss {
	if other.IsZero() {
		panic("pramana: division by zero")
	}
	// g/other = g * conj(other) / |other|²
	conj := other.Conjugate()
	num := g.Mul(conj)
	normSq := other.Norm() // This is real-valued

	// Divide both real and imag parts by normSq (which is normNum/normDen, real)
	// real: (num.a/num.b) / (norm.a/norm.b) = (num.a * norm.b) / (num.b * norm.a)
	realNum := new(big.Int).Mul(num.a, normSq.b)
	realDen := new(big.Int).Mul(num.b, normSq.a)
	imagNum := new(big.Int).Mul(num.c, normSq.b)
	imagDen := new(big.Int).Mul(num.d, normSq.a)

	return NewGaussBig(realNum, realDen, imagNum, imagDen)
}

// Mod returns the modulo for real Gaussian rationals. Panics if either is not real.
func (g *Gauss) Mod(other *Gauss) *Gauss {
	if !g.IsReal() || !other.IsReal() {
		panic("pramana: modulo only defined for real values")
	}
	if other.IsZero() {
		panic("pramana: modulo by zero")
	}
	// a/b mod e/f: q = floor((a/b) / (e/f)), result = a/b - q * e/f
	quot := g.Div(other)
	// Floor the real part
	floored := GaussFloor(quot)
	return g.Sub(floored.Mul(other))
}

// Pow returns g raised to the power n.
func GaussPow(base *Gauss, n int) *Gauss {
	if n == 0 {
		return GaussOne()
	}
	if n < 0 {
		return GaussPow(base.Reciprocal(), -n)
	}
	result := GaussOne()
	b := base
	exp := n
	for exp > 0 {
		if exp%2 == 1 {
			result = result.Mul(b)
		}
		b = b.Mul(b)
		exp /= 2
	}
	return result
}

// Inc returns g + 1 (increments real part).
func (g *Gauss) Inc() *Gauss {
	return g.Add(GaussOne())
}

// Dec returns g - 1 (decrements real part).
func (g *Gauss) Dec() *Gauss {
	return g.Sub(GaussOne())
}

// --- Comparison ---

// Equal returns true if g and other represent the same value.
func (g *Gauss) Equal(other *Gauss) bool {
	return g.a.Cmp(other.a) == 0 && g.b.Cmp(other.b) == 0 &&
		g.c.Cmp(other.c) == 0 && g.d.Cmp(other.d) == 0
}

// Cmp compares two Gauss values. Compares real parts first, then imaginary.
// Returns -1, 0, or 1.
func (g *Gauss) Cmp(other *Gauss) int {
	// Compare real parts: a/b vs e/f => a*f vs e*b
	lhs := new(big.Int).Mul(g.a, other.b)
	rhs := new(big.Int).Mul(other.a, g.b)
	c := lhs.Cmp(rhs)
	if c != 0 {
		return c
	}
	// Compare imaginary parts: c/d vs g/h => c*h vs g*d
	lhs = new(big.Int).Mul(g.c, other.d)
	rhs = new(big.Int).Mul(other.c, g.d)
	return lhs.Cmp(rhs)
}

// Lt returns true if g < other.
func (g *Gauss) Lt(other *Gauss) bool { return g.Cmp(other) < 0 }

// Gt returns true if g > other.
func (g *Gauss) Gt(other *Gauss) bool { return g.Cmp(other) > 0 }

// Lte returns true if g <= other.
func (g *Gauss) Lte(other *Gauss) bool { return g.Cmp(other) <= 0 }

// Gte returns true if g >= other.
func (g *Gauss) Gte(other *Gauss) bool { return g.Cmp(other) >= 0 }

// --- Static Functions ---

// GaussMin returns the minimum of two Gauss values.
func GaussMin(a, b *Gauss) *Gauss {
	if a.Lte(b) {
		return a
	}
	return b
}

// GaussMax returns the maximum of two Gauss values.
func GaussMax(a, b *Gauss) *Gauss {
	if a.Gte(b) {
		return a
	}
	return b
}

// GaussClamp clamps value to the range [min, max].
func GaussClamp(value, min, max *Gauss) *Gauss {
	if value.Lt(min) {
		return min
	}
	if value.Gt(max) {
		return max
	}
	return value
}

// GaussAbs returns the absolute value for real Gauss values.
func GaussAbs(g *Gauss) *Gauss {
	if !g.IsReal() {
		panic("pramana: Abs only defined for real values")
	}
	if g.IsNegative() {
		return g.Neg()
	}
	return g
}

// GaussSign returns -1, 0, or 1 for real Gauss values.
func GaussSign(g *Gauss) int {
	if !g.IsReal() {
		panic("pramana: Sign only defined for real values")
	}
	return g.a.Sign()
}

// GaussFloor floors both real and imaginary parts to nearest integer toward negative infinity.
func GaussFloor(g *Gauss) *Gauss {
	realFloor := new(big.Int).Div(g.a, g.b)
	if g.a.Sign() < 0 {
		rem := new(big.Int).Mod(g.a, g.b)
		if rem.Sign() != 0 {
			realFloor.Sub(realFloor, big.NewInt(1))
		}
	}
	imagFloor := new(big.Int).Div(g.c, g.d)
	if g.c.Sign() < 0 {
		rem := new(big.Int).Mod(g.c, g.d)
		if rem.Sign() != 0 {
			imagFloor.Sub(imagFloor, big.NewInt(1))
		}
	}
	return NewGaussFromBigInts(realFloor, imagFloor)
}

// GaussCeiling ceils both real and imaginary parts to nearest integer toward positive infinity.
func GaussCeiling(g *Gauss) *Gauss {
	realCeil := new(big.Int).Div(g.a, g.b)
	rem := new(big.Int).Mod(g.a, g.b)
	if rem.Sign() != 0 && g.a.Sign() > 0 {
		realCeil.Add(realCeil, big.NewInt(1))
	}
	imagCeil := new(big.Int).Div(g.c, g.d)
	rem = new(big.Int).Mod(g.c, g.d)
	if rem.Sign() != 0 && g.c.Sign() > 0 {
		imagCeil.Add(imagCeil, big.NewInt(1))
	}
	return NewGaussFromBigInts(realCeil, imagCeil)
}

// GaussTruncate truncates both parts toward zero.
func GaussTruncate(g *Gauss) *Gauss {
	realTrunc := new(big.Int).Quo(g.a, g.b)
	imagTrunc := new(big.Int).Quo(g.c, g.d)
	return NewGaussFromBigInts(realTrunc, imagTrunc)
}

// --- Associates ---

// Associates returns the three non-trivial associates (multiplied by -1, i, -i).
func (g *Gauss) Associates() []*Gauss {
	return []*Gauss{
		g.Neg(),
		g.Mul(GaussI()),
		g.Mul(NewGaussInt(0, -1)),
	}
}

// IsAssociate returns true if g is an associate of other.
func (g *Gauss) IsAssociate(other *Gauss) bool {
	if g.Equal(other) {
		return true
	}
	for _, assoc := range g.Associates() {
		if assoc.Equal(other) {
			return true
		}
	}
	return false
}

// --- Conversion ---

// RealFloat64 returns the real part as a float64 approximation.
func (g *Gauss) RealFloat64() float64 {
	num := new(big.Float).SetInt(g.a)
	den := new(big.Float).SetInt(g.b)
	result, _ := new(big.Float).Quo(num, den).Float64()
	return result
}

// ImagFloat64 returns the imaginary part as a float64 approximation.
func (g *Gauss) ImagFloat64() float64 {
	num := new(big.Float).SetInt(g.c)
	den := new(big.Float).SetInt(g.d)
	result, _ := new(big.Float).Quo(num, den).Float64()
	return result
}

// ToComplex128 returns the value as a complex128 approximation.
func (g *Gauss) ToComplex128() complex128 {
	return complex(g.RealFloat64(), g.ImagFloat64())
}

// ToPolar returns (magnitude, phase) as float64 approximations.
func (g *Gauss) ToPolar() (float64, float64) {
	return g.Magnitude(), g.Phase()
}

// GaussFromPolar creates a Gauss from polar coordinates (approximate via float64).
func GaussFromPolar(magnitude, phase float64) *Gauss {
	return GaussFromFloat64(magnitude*math.Cos(phase), magnitude*math.Sin(phase))
}

// GaussFromFloat64 creates a Gauss from two float64 values.
// Uses a max denominator of 1,000,000 for rational approximation.
func GaussFromFloat64(real, imag float64) *Gauss {
	rr := new(big.Rat).SetFloat64(real)
	ri := new(big.Rat).SetFloat64(imag)
	if rr == nil {
		rr = new(big.Rat)
	}
	if ri == nil {
		ri = new(big.Rat)
	}
	return NewGaussBig(rr.Num(), rr.Denom(), ri.Num(), ri.Denom())
}

// --- String Representations ---

// String returns a human-readable string representation.
func (g *Gauss) String() string {
	realStr := fractionString(g.a, g.b)
	imagStr := fractionString(g.c, g.d)

	if g.IsZero() {
		return "0"
	}
	if g.IsReal() {
		return realStr
	}
	if g.IsPurelyImaginary() {
		if g.c.Cmp(big.NewInt(1)) == 0 && g.d.Cmp(big.NewInt(1)) == 0 {
			return "i"
		}
		if g.c.Cmp(big.NewInt(-1)) == 0 && g.d.Cmp(big.NewInt(1)) == 0 {
			return "-i"
		}
		return imagStr + "i"
	}

	sign := " + "
	if g.c.Sign() < 0 {
		sign = " - "
		imagStr = fractionString(new(big.Int).Abs(g.c), g.d)
	}
	if imagStr == "1" {
		imagStr = ""
	}
	return realStr + sign + imagStr + "i"
}

// fractionString formats a fraction a/b as a string.
func fractionString(num, den *big.Int) string {
	if den.Cmp(big.NewInt(1)) == 0 {
		return num.String()
	}
	return fmt.Sprintf("%s/%s", num.String(), den.String())
}

// ToRawString returns the raw format "<a,b,c,d>".
func (g *Gauss) ToRawString() string {
	return fmt.Sprintf("<%s,%s,%s,%s>", g.a.String(), g.b.String(), g.c.String(), g.d.String())
}

// ToDecimalString returns a decimal approximation with the given precision.
func (g *Gauss) ToDecimalString(precision int) string {
	re := g.RealFloat64()
	im := g.ImagFloat64()
	format := fmt.Sprintf("%%.%df", precision)
	reStr := fmt.Sprintf(format, re)
	if im == 0 {
		return reStr
	}
	if re == 0 {
		return fmt.Sprintf(format+"i", im)
	}
	sign := " + "
	if im < 0 {
		sign = " - "
		im = -im
	}
	return reStr + sign + fmt.Sprintf(format+"i", im)
}

// ToImproperFractionString returns the improper fraction notation.
func (g *Gauss) ToImproperFractionString() string {
	realStr := fractionString(g.a, g.b)
	if g.IsReal() {
		return realStr
	}
	imagStr := fractionString(g.c, g.d)
	if g.IsPurelyImaginary() {
		return imagStr + "i"
	}
	sign := " + "
	if g.c.Sign() < 0 {
		sign = " - "
		imagStr = fractionString(new(big.Int).Abs(g.c), g.d)
	}
	return realStr + sign + imagStr + "i"
}

// ToMixedString returns the mixed fraction notation (e.g., "3 & 1/2 + 1/4i").
func (g *Gauss) ToMixedString() string {
	realMixed := mixedFraction(g.a, g.b)
	if g.IsReal() {
		return realMixed
	}
	imagMixed := mixedFraction(g.c, g.d)
	if g.IsPurelyImaginary() {
		return imagMixed + "i"
	}
	sign := " + "
	if g.c.Sign() < 0 {
		sign = " - "
		imagMixed = mixedFraction(new(big.Int).Abs(g.c), g.d)
	}
	return realMixed + sign + imagMixed + "i"
}

// mixedFraction converts a fraction to mixed notation.
func mixedFraction(num, den *big.Int) string {
	if den.Cmp(big.NewInt(1)) == 0 {
		return num.String()
	}
	absNum := new(big.Int).Abs(num)
	whole := new(big.Int).Div(absNum, den)
	rem := new(big.Int).Mod(absNum, den)
	sign := ""
	if num.Sign() < 0 {
		sign = "-"
	}
	if whole.Sign() == 0 {
		return fmt.Sprintf("%s%s/%s", sign, rem.String(), den.String())
	}
	if rem.Sign() == 0 {
		return fmt.Sprintf("%s%s", sign, whole.String())
	}
	return fmt.Sprintf("%s%s & %s/%s", sign, whole.String(), rem.String(), den.String())
}

// --- Parsing ---

// GaussParse parses a Gauss from "a,b,c,d" format.
func GaussParse(s string) (*Gauss, error) {
	parts := strings.Split(strings.TrimSpace(s), ",")
	if len(parts) != 4 {
		return nil, fmt.Errorf("pramana: expected 4 comma-separated values, got %d", len(parts))
	}
	a, ok := new(big.Int).SetString(strings.TrimSpace(parts[0]), 10)
	if !ok {
		return nil, fmt.Errorf("pramana: invalid integer %q", parts[0])
	}
	b, ok := new(big.Int).SetString(strings.TrimSpace(parts[1]), 10)
	if !ok {
		return nil, fmt.Errorf("pramana: invalid integer %q", parts[1])
	}
	c, ok := new(big.Int).SetString(strings.TrimSpace(parts[2]), 10)
	if !ok {
		return nil, fmt.Errorf("pramana: invalid integer %q", parts[2])
	}
	d, ok := new(big.Int).SetString(strings.TrimSpace(parts[3]), 10)
	if !ok {
		return nil, fmt.Errorf("pramana: invalid integer %q", parts[3])
	}
	return NewGaussBig(a, b, c, d), nil
}

// GaussFromPramana parses from "pra:num:a,b,c,d" format.
func GaussFromPramana(s string) (*Gauss, error) {
	prefix := "pra:num:"
	if !strings.HasPrefix(s, prefix) {
		return nil, fmt.Errorf("pramana: expected prefix %q", prefix)
	}
	return GaussParse(s[len(prefix):])
}

// --- Pramana Identity ---

// PramanaKey returns the canonical key: "a,b,c,d".
func (g *Gauss) PramanaKey() string {
	return fmt.Sprintf("%s,%s,%s,%s", g.a.String(), g.b.String(), g.c.String(), g.d.String())
}

// PramanaID returns the UUID v5 for this value.
func (g *Gauss) PramanaID() uuid.UUID {
	return uuid.NewSHA1(PramanaNamespace, []byte(g.PramanaKey()))
}

// PramanaLabel returns the Pramana label: "pra:num:a,b,c,d".
func (g *Gauss) PramanaLabel() string {
	return "pra:num:" + g.PramanaKey()
}

// PramanaURL returns the entity URL with the UUID.
func (g *Gauss) PramanaURL() string {
	return "https://pramana.dev/entity/" + g.PramanaID().String()
}

// PramanaHashURL returns the entity URL with the UUID (same as PramanaURL).
func (g *Gauss) PramanaHashURL() string {
	return g.PramanaURL()
}
