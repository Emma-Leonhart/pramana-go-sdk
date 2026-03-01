package pramana

import (
	"fmt"
	"math/big"
	"math/rand"
	"strings"

	"github.com/google/uuid"
)

// Gint represents a Gaussian integer (Z[i]): real + imag*i
// where both components are arbitrary-precision integers.
type Gint struct {
	real *big.Int
	imag *big.Int
}

// NewGint creates a new Gaussian integer real + imag*i.
func NewGint(real, imag int64) *Gint {
	return &Gint{real: big.NewInt(real), imag: big.NewInt(imag)}
}

// NewGintReal creates a Gaussian integer with zero imaginary part.
func NewGintReal(real int64) *Gint {
	return NewGint(real, 0)
}

// NewGintBig creates a Gaussian integer from big.Int values.
func NewGintBig(real, imag *big.Int) *Gint {
	return &Gint{real: copyBigInt(real), imag: copyBigInt(imag)}
}

// NewGintFromArray creates a Gint from a two-element slice [real, imag].
func NewGintFromArray(arr []*big.Int) *Gint {
	if len(arr) < 2 {
		panic("pramana: array must have at least 2 elements")
	}
	return NewGintBig(arr[0], arr[1])
}

// --- Constants ---

// GintZero returns the zero value.
func GintZero() *Gint { return NewGint(0, 0) }

// GintOne returns the value 1.
func GintOne() *Gint { return NewGint(1, 0) }

// GintMinusOne returns the value -1.
func GintMinusOne() *Gint { return NewGint(-1, 0) }

// GintI returns the imaginary unit i.
func GintI() *Gint { return NewGint(0, 1) }

// GintEye returns the imaginary unit i (alias for GintI).
func GintEye() *Gint { return GintI() }

// GintUnits returns the four Gaussian units: [1, -1, i, -i].
func GintUnits() []*Gint {
	return []*Gint{GintOne(), GintMinusOne(), GintI(), NewGint(0, -1)}
}

// GintTwo returns 1+i (useful as the even norm detector in Z[i]).
func GintTwo() *Gint { return NewGint(1, 1) }

// GintRandom returns a random Gaussian integer within the given bounds.
func GintRandom(re1, re2, im1, im2 int64) *Gint {
	re := re1 + rand.Int63n(re2-re1+1)
	im := im1 + rand.Int63n(im2-im1+1)
	return NewGint(re, im)
}

// --- Accessors ---

// Real returns the real component.
func (g *Gint) Real() *big.Int { return copyBigInt(g.real) }

// Imag returns the imaginary component.
func (g *Gint) Imag() *big.Int { return copyBigInt(g.imag) }

// --- Properties ---

// IsReal returns true if the imaginary part is zero.
func (g *Gint) IsReal() bool {
	return g.imag.Sign() == 0
}

// IsPurelyImaginary returns true if the real part is zero and imaginary is nonzero.
func (g *Gint) IsPurelyImaginary() bool {
	return g.real.Sign() == 0 && g.imag.Sign() != 0
}

// IsZero returns true if both components are zero.
func (g *Gint) IsZero() bool {
	return g.real.Sign() == 0 && g.imag.Sign() == 0
}

// IsOne returns true if the value equals 1.
func (g *Gint) IsOne() bool {
	return g.real.Cmp(big.NewInt(1)) == 0 && g.imag.Sign() == 0
}

// IsUnit returns true if the Gaussian integer is one of {1, -1, i, -i}.
func (g *Gint) IsUnit() bool {
	return g.Norm().Cmp(big.NewInt(1)) == 0
}

// IsPositive returns true if the value is real and positive.
func (g *Gint) IsPositive() bool {
	return g.IsReal() && g.real.Sign() > 0
}

// IsNegative returns true if the value is real and negative.
func (g *Gint) IsNegative() bool {
	return g.IsReal() && g.real.Sign() < 0
}

// --- Derived Values ---

// Conjugate returns the complex conjugate (real - imag*i).
func (g *Gint) Conjugate() *Gint {
	return NewGintBig(g.real, new(big.Int).Neg(g.imag))
}

// Norm returns real² + imag² (the squared magnitude).
func (g *Gint) Norm() *big.Int {
	r2 := new(big.Int).Mul(g.real, g.real)
	i2 := new(big.Int).Mul(g.imag, g.imag)
	return r2.Add(r2, i2)
}

// --- Arithmetic ---

// Add returns g + other.
func (g *Gint) Add(other *Gint) *Gint {
	return NewGintBig(
		new(big.Int).Add(g.real, other.real),
		new(big.Int).Add(g.imag, other.imag),
	)
}

// Sub returns g - other.
func (g *Gint) Sub(other *Gint) *Gint {
	return NewGintBig(
		new(big.Int).Sub(g.real, other.real),
		new(big.Int).Sub(g.imag, other.imag),
	)
}

// Neg returns -g.
func (g *Gint) Neg() *Gint {
	return NewGintBig(
		new(big.Int).Neg(g.real),
		new(big.Int).Neg(g.imag),
	)
}

// Mul returns g * other.
func (g *Gint) Mul(other *Gint) *Gint {
	// (a + bi)(c + di) = (ac - bd) + (ad + bc)i
	ac := new(big.Int).Mul(g.real, other.real)
	bd := new(big.Int).Mul(g.imag, other.imag)
	ad := new(big.Int).Mul(g.real, other.imag)
	bc := new(big.Int).Mul(g.imag, other.real)
	return NewGintBig(
		new(big.Int).Sub(ac, bd),
		new(big.Int).Add(ad, bc),
	)
}

// DivExact performs exact division, returning a Gauss (Gaussian rational).
func (g *Gint) DivExact(other *Gint) *Gauss {
	return g.ToGauss().Div(other.ToGauss())
}

// FloorDiv returns the Gaussian integer quotient using rounding division.
func (g *Gint) FloorDiv(other *Gint) *Gint {
	q, _ := GintModifiedDivmod(g, other)
	return q
}

// ModG returns the Gaussian integer remainder from modified division.
func (g *Gint) ModG(other *Gint) *Gint {
	_, r := GintModifiedDivmod(g, other)
	return r
}

// GintPow returns base raised to the power n.
func GintPow(base *Gint, n int) *Gint {
	if n == 0 {
		return GintOne()
	}
	if n < 0 {
		panic("pramana: negative exponent not supported for Gint; use DivExact for Gauss result")
	}
	result := GintOne()
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

// Inc returns g + 1.
func (g *Gint) Inc() *Gint {
	return g.Add(GintOne())
}

// Dec returns g - 1.
func (g *Gint) Dec() *Gint {
	return g.Sub(GintOne())
}

// --- Comparison ---

// Equal returns true if g and other have the same components.
func (g *Gint) Equal(other *Gint) bool {
	return g.real.Cmp(other.real) == 0 && g.imag.Cmp(other.imag) == 0
}

// Cmp compares two Gint values. Compares real parts first, then imaginary.
func (g *Gint) Cmp(other *Gint) int {
	c := g.real.Cmp(other.real)
	if c != 0 {
		return c
	}
	return g.imag.Cmp(other.imag)
}

// Lt returns true if g < other.
func (g *Gint) Lt(other *Gint) bool { return g.Cmp(other) < 0 }

// Gt returns true if g > other.
func (g *Gint) Gt(other *Gint) bool { return g.Cmp(other) > 0 }

// Lte returns true if g <= other.
func (g *Gint) Lte(other *Gint) bool { return g.Cmp(other) <= 0 }

// Gte returns true if g >= other.
func (g *Gint) Gte(other *Gint) bool { return g.Cmp(other) >= 0 }

// --- Associates ---

// Associates returns the three non-trivial associates (multiplied by -1, i, -i).
func (g *Gint) Associates() []*Gint {
	// i * (a + bi) = -b + ai
	// -i * (a + bi) = b - ai
	return []*Gint{
		g.Neg(),
		NewGintBig(new(big.Int).Neg(g.imag), copyBigInt(g.real)),
		NewGintBig(copyBigInt(g.imag), new(big.Int).Neg(g.real)),
	}
}

// IsAssociate returns true if g is an associate of other.
func (g *Gint) IsAssociate(other *Gint) bool {
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

// --- Number Theory ---

// GintModifiedDivmod performs division with rounding, returning (quotient, remainder)
// such that |remainder|² < |b|²/2.
func GintModifiedDivmod(a, b *Gint) (*Gint, *Gint) {
	if b.IsZero() {
		panic("pramana: division by zero")
	}
	// q = round(a * conj(b) / |b|²)
	conj := b.Conjugate()
	product := a.Mul(conj)
	normB := b.Norm()

	// Round real and imag parts
	qReal := roundDiv(product.real, normB)
	qImag := roundDiv(product.imag, normB)

	q := NewGintBig(qReal, qImag)
	r := a.Sub(b.Mul(q))
	return q, r
}

// roundDiv computes round(a/b) using rounding-half-away-from-zero.
func roundDiv(a, b *big.Int) *big.Int {
	// Use the formula: round(a/b) = floor((2*a + b) / (2*b))
	// This handles both positive and negative correctly
	two := big.NewInt(2)
	twoA := new(big.Int).Mul(two, a)
	twoB := new(big.Int).Mul(two, b)

	// Adjust for sign: if b < 0, we need to flip
	num := new(big.Int).Add(twoA, b)
	if b.Sign() < 0 {
		num = new(big.Int).Sub(twoA, b)
		twoB.Neg(twoB)
	}

	// Euclidean-style division
	result := new(big.Int)
	if num.Sign() >= 0 {
		result.Div(num, twoB)
	} else {
		// For negative numerators, adjust to get rounding behavior
		num.Sub(num, new(big.Int).Sub(twoB, big.NewInt(1)))
		result.Div(num, twoB)
	}
	return result
}

// GintGCD returns the greatest common divisor of a and b using the Euclidean algorithm.
func GintGCD(a, b *Gint) *Gint {
	x := a
	y := b
	for !y.IsZero() {
		_, r := GintModifiedDivmod(x, y)
		x = y
		y = r
	}
	return x
}

// GintXGCD returns (gcd, x, y) such that a*x + b*y = gcd (extended Euclidean algorithm).
func GintXGCD(alpha, beta *Gint) (*Gint, *Gint, *Gint) {
	if beta.IsZero() {
		return alpha, GintOne(), GintZero()
	}

	oldR, r := alpha, beta
	oldS, s := GintOne(), GintZero()
	oldT, t := GintZero(), GintOne()

	for !r.IsZero() {
		q, _ := GintModifiedDivmod(oldR, r)
		oldR, r = r, oldR.Sub(q.Mul(r))
		oldS, s = s, oldS.Sub(q.Mul(s))
		oldT, t = t, oldT.Sub(q.Mul(t))
	}

	return oldR, oldS, oldT
}

// GintIsRelativelyPrime returns true if the GCD of a and b is a unit.
func GintIsRelativelyPrime(a, b *Gint) bool {
	return GintGCD(a, b).IsUnit()
}

// GintIsGaussianPrime returns true if x is a Gaussian prime.
func GintIsGaussianPrime(x *Gint) bool {
	if x.IsZero() {
		return false
	}
	if x.IsUnit() {
		return false
	}

	// If both parts nonzero: prime iff norm is prime
	if x.real.Sign() != 0 && x.imag.Sign() != 0 {
		return IsPrime(x.Norm())
	}

	// If one part is zero: the nonzero part must be prime and ≡ 3 (mod 4)
	var p *big.Int
	if x.imag.Sign() == 0 {
		p = new(big.Int).Abs(x.real)
	} else {
		p = new(big.Int).Abs(x.imag)
	}

	if !IsPrime(p) {
		return false
	}

	four := big.NewInt(4)
	mod := new(big.Int).Mod(p, four)
	return mod.Cmp(big.NewInt(3)) == 0
}

// GintNormsDivide checks if the norm of the larger divides the norm of the smaller.
// Returns the quotient (larger/smaller) if it divides, or nil if not.
func GintNormsDivide(a, b *Gint) *big.Int {
	na := a.Norm()
	nb := b.Norm()
	var lg, sm *big.Int
	if na.Cmp(nb) >= 0 {
		lg, sm = na, nb
	} else {
		lg, sm = nb, na
	}
	if sm.Sign() == 0 {
		return nil
	}
	mod := new(big.Int).Mod(lg, sm)
	if mod.Sign() != 0 {
		return nil
	}
	return new(big.Int).Div(lg, sm)
}

// GintCongruentModulo checks if (a - b) is divisible by c.
// Returns (isCongruent, (a-b)/c).
func GintCongruentModulo(a, b, c *Gint) (bool, *Gauss) {
	diff := a.Sub(b)
	result := diff.DivExact(c)
	return result.IsGaussianInteger(), result
}

// --- Conversion ---

// ToGauss converts a Gint to a Gauss (Gaussian rational with denominators of 1).
func (g *Gint) ToGauss() *Gauss {
	return NewGaussFromBigInts(copyBigInt(g.real), copyBigInt(g.imag))
}

// GintFromGauss converts a Gauss to a Gint. Panics if not a Gaussian integer.
func GintFromGauss(g *Gauss) *Gint {
	if !g.IsGaussianInteger() {
		panic("pramana: Gauss value is not a Gaussian integer")
	}
	return NewGintBig(g.a, g.c)
}

// ToArray returns [real, imag] as a slice of *big.Int.
func (g *Gint) ToArray() []*big.Int {
	return []*big.Int{copyBigInt(g.real), copyBigInt(g.imag)}
}

// --- String Representations ---

// String returns a human-readable string representation.
func (g *Gint) String() string {
	if g.IsZero() {
		return "0"
	}
	if g.IsReal() {
		return g.real.String()
	}
	if g.IsPurelyImaginary() {
		if g.imag.Cmp(big.NewInt(1)) == 0 {
			return "i"
		}
		if g.imag.Cmp(big.NewInt(-1)) == 0 {
			return "-i"
		}
		return g.imag.String() + "i"
	}

	sign := " + "
	imagPart := g.imag.String()
	if g.imag.Sign() < 0 {
		sign = " - "
		imagPart = new(big.Int).Abs(g.imag).String()
	}
	if imagPart == "1" {
		imagPart = ""
	}
	return g.real.String() + sign + imagPart + "i"
}

// ToRawString returns the raw format "(real, imag)".
func (g *Gint) ToRawString() string {
	return fmt.Sprintf("(%s, %s)", g.real.String(), g.imag.String())
}

// GoString implements fmt.GoStringer for %#v formatting.
func (g *Gint) GoString() string {
	if g.IsReal() {
		return fmt.Sprintf("Gint(%s)", g.real.String())
	}
	return fmt.Sprintf("Gint(%s, %s)", g.real.String(), g.imag.String())
}

// --- Pramana Identity ---

// PramanaKey returns the canonical key: "real,1,imag,1".
func (g *Gint) PramanaKey() string {
	return fmt.Sprintf("%s,1,%s,1", g.real.String(), g.imag.String())
}

// PramanaID returns the UUID v5 for this value.
func (g *Gint) PramanaID() uuid.UUID {
	return uuid.NewSHA1(PramanaNamespace, []byte(g.PramanaKey()))
}

// PramanaLabel returns the Pramana label: "pra:num:real,1,imag,1".
func (g *Gint) PramanaLabel() string {
	return "pra:num:" + g.PramanaKey()
}

// PramanaURL returns the entity URL with the UUID.
func (g *Gint) PramanaURL() string {
	return "https://pramana.dev/entity/" + g.PramanaID().String()
}

// --- Parsing ---

// GintParse parses a Gint from "real,imag" format.
func GintParse(s string) (*Gint, error) {
	parts := strings.Split(strings.TrimSpace(s), ",")
	if len(parts) != 2 {
		return nil, fmt.Errorf("pramana: expected 2 comma-separated values, got %d", len(parts))
	}
	real, ok := new(big.Int).SetString(strings.TrimSpace(parts[0]), 10)
	if !ok {
		return nil, fmt.Errorf("pramana: invalid integer %q", parts[0])
	}
	imag, ok := new(big.Int).SetString(strings.TrimSpace(parts[1]), 10)
	if !ok {
		return nil, fmt.Errorf("pramana: invalid integer %q", parts[1])
	}
	return NewGintBig(real, imag), nil
}
