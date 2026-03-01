package pramana

import (
	"math/big"
	"testing"
)

func TestGintConstruction(t *testing.T) {
	g := NewGint(3, 4)
	if g.Real().Int64() != 3 || g.Imag().Int64() != 4 {
		t.Errorf("NewGint(3,4) = %s, want 3+4i", g)
	}
}

func TestGintProperties(t *testing.T) {
	zero := GintZero()
	one := GintOne()
	i := GintI()
	real := NewGintReal(5)

	if !zero.IsZero() {
		t.Error("Zero should be zero")
	}
	if !one.IsOne() {
		t.Error("One should be one")
	}
	if !one.IsUnit() {
		t.Error("One should be a unit")
	}
	if !i.IsUnit() {
		t.Error("i should be a unit")
	}
	if !i.IsPurelyImaginary() {
		t.Error("i should be purely imaginary")
	}
	if !real.IsReal() {
		t.Error("5 should be real")
	}
	if !real.IsPositive() {
		t.Error("5 should be positive")
	}
	if !GintMinusOne().IsNegative() {
		t.Error("-1 should be negative")
	}
}

func TestGintArithmetic(t *testing.T) {
	a := NewGint(3, 4)
	b := NewGint(1, 2)

	// Addition
	sum := a.Add(b)
	if !sum.Equal(NewGint(4, 6)) {
		t.Errorf("(3+4i) + (1+2i) = %s, want 4+6i", sum)
	}

	// Subtraction
	diff := a.Sub(b)
	if !diff.Equal(NewGint(2, 2)) {
		t.Errorf("(3+4i) - (1+2i) = %s, want 2+2i", diff)
	}

	// Multiplication: (3+4i)(1+2i) = -5+10i
	prod := a.Mul(b)
	if !prod.Equal(NewGint(-5, 10)) {
		t.Errorf("(3+4i) * (1+2i) = %s, want -5+10i", prod)
	}

	// Negation
	neg := a.Neg()
	if !neg.Equal(NewGint(-3, -4)) {
		t.Errorf("-(3+4i) = %s, want -3-4i", neg)
	}
}

func TestGintConjugate(t *testing.T) {
	g := NewGint(3, 4)
	conj := g.Conjugate()
	if !conj.Equal(NewGint(3, -4)) {
		t.Errorf("conj(3+4i) = %s, want 3-4i", conj)
	}
}

func TestGintNorm(t *testing.T) {
	g := NewGint(3, 4)
	norm := g.Norm()
	if norm.Int64() != 25 {
		t.Errorf("|3+4i|² = %s, want 25", norm)
	}
}

func TestGintPow(t *testing.T) {
	i := GintI()
	// i² = -1
	result := GintPow(i, 2)
	if !result.Equal(GintMinusOne()) {
		t.Errorf("i² = %s, want -1", result)
	}

	// i⁴ = 1
	result = GintPow(i, 4)
	if !result.Equal(GintOne()) {
		t.Errorf("i⁴ = %s, want 1", result)
	}

	// (2+i)² = 4+4i+i² = 3+4i
	g := NewGint(2, 1)
	sq := GintPow(g, 2)
	if !sq.Equal(NewGint(3, 4)) {
		t.Errorf("(2+i)² = %s, want 3+4i", sq)
	}
}

func TestGintUnits(t *testing.T) {
	units := GintUnits()
	if len(units) != 4 {
		t.Errorf("GintUnits() has %d elements, want 4", len(units))
	}
	for _, u := range units {
		if !u.IsUnit() {
			t.Errorf("%s should be a unit", u)
		}
	}
}

func TestGintAssociates(t *testing.T) {
	g := NewGint(2, 3)
	assocs := g.Associates()
	if len(assocs) != 3 {
		t.Errorf("Associates() has %d elements, want 3", len(assocs))
	}
	for _, a := range assocs {
		if !g.IsAssociate(a) {
			t.Errorf("%s should be associate of %s", a, g)
		}
	}
}

func TestGintModifiedDivmod(t *testing.T) {
	a := NewGint(7, 3)
	b := NewGint(2, 1)
	q, r := GintModifiedDivmod(a, b)

	// Verify: a = b*q + r
	check := b.Mul(q).Add(r)
	if !check.Equal(a) {
		t.Errorf("b*q + r = %s, want %s", check, a)
	}

	// Verify: |r|² < |b|²
	rNorm := r.Norm()
	bNorm := b.Norm()
	if rNorm.Cmp(bNorm) >= 0 {
		t.Errorf("|r|² = %s should be < |b|² = %s", rNorm, bNorm)
	}
}

func TestGintGCD(t *testing.T) {
	a := NewGint(3, 1)
	b := NewGint(1, 2)
	g := GintGCD(a, b)

	// GCD should divide both a and b
	_, rA := GintModifiedDivmod(a, g)
	_, rB := GintModifiedDivmod(b, g)
	if !rA.IsZero() {
		t.Errorf("GCD(%s, %s) = %s does not divide %s", a, b, g, a)
	}
	if !rB.IsZero() {
		t.Errorf("GCD(%s, %s) = %s does not divide %s", a, b, g, b)
	}
}

func TestGintXGCD(t *testing.T) {
	a := NewGint(3, 1)
	b := NewGint(1, 2)
	g, x, y := GintXGCD(a, b)

	// Verify: a*x + b*y = gcd
	check := a.Mul(x).Add(b.Mul(y))
	if !check.Equal(g) {
		t.Errorf("a*x + b*y = %s, want gcd = %s", check, g)
	}
}

func TestGintIsRelativelyPrime(t *testing.T) {
	// 3+2i (norm 13) and 1+i (norm 2) share no common Gaussian prime factor
	a := NewGint(3, 2)
	b := NewGint(1, 1)
	if !GintIsRelativelyPrime(a, b) {
		t.Errorf("%s and %s should be relatively prime", a, b)
	}
}

func TestGintIsGaussianPrime(t *testing.T) {
	// 3 is a Gaussian prime (3 ≡ 3 mod 4 and is prime)
	if !GintIsGaussianPrime(NewGintReal(3)) {
		t.Error("3 should be a Gaussian prime")
	}

	// 2 is NOT a Gaussian prime (2 = -i(1+i)²)
	if GintIsGaussianPrime(NewGintReal(2)) {
		t.Error("2 should not be a Gaussian prime")
	}

	// 1+i has norm 2 (prime), so it is a Gaussian prime
	if !GintIsGaussianPrime(NewGint(1, 1)) {
		t.Error("1+i should be a Gaussian prime")
	}

	// 5 is NOT a Gaussian prime (5 = (2+i)(2-i))
	if GintIsGaussianPrime(NewGintReal(5)) {
		t.Error("5 should not be a Gaussian prime")
	}
}

func TestGintToGauss(t *testing.T) {
	g := NewGint(3, 4)
	gauss := g.ToGauss()
	if !gauss.IsGaussianInteger() {
		t.Error("ToGauss() should be a Gaussian integer")
	}
	if gauss.A().Int64() != 3 || gauss.C().Int64() != 4 {
		t.Errorf("ToGauss() = %s, want 3+4i", gauss)
	}
}

func TestGintFromGauss(t *testing.T) {
	gauss := NewGaussInt(3, 4)
	g := GintFromGauss(gauss)
	if g.Real().Int64() != 3 || g.Imag().Int64() != 4 {
		t.Errorf("GintFromGauss(3+4i) = %s", g)
	}
}

func TestGintDivExact(t *testing.T) {
	a := NewGint(1, 0)
	b := NewGint(0, 1)
	result := a.DivExact(b)
	// 1/i = -i
	expected := NewGaussInt(0, -1)
	if !result.Equal(expected) {
		t.Errorf("1/i = %s, want -i", result)
	}
}

func TestGintFloorDiv(t *testing.T) {
	a := NewGint(7, 3)
	b := NewGint(2, 1)
	q := a.FloorDiv(b)

	// Verify: a - b*q has small norm
	r := a.Sub(b.Mul(q))
	rNorm := r.Norm()
	bNorm := b.Norm()
	if rNorm.Cmp(bNorm) >= 0 {
		t.Errorf("FloorDiv remainder norm %s >= divisor norm %s", rNorm, bNorm)
	}
}

func TestGintString(t *testing.T) {
	tests := []struct {
		g    *Gint
		want string
	}{
		{GintZero(), "0"},
		{GintOne(), "1"},
		{GintMinusOne(), "-1"},
		{GintI(), "i"},
		{NewGint(0, -1), "-i"},
		{NewGint(3, 4), "3 + 4i"},
		{NewGint(3, -4), "3 - 4i"},
		{NewGint(0, 5), "5i"},
	}
	for _, tt := range tests {
		got := tt.g.String()
		if got != tt.want {
			t.Errorf("Gint(%d,%d).String() = %q, want %q",
				tt.g.Real().Int64(), tt.g.Imag().Int64(), got, tt.want)
		}
	}
}

func TestGintPramanaIdentity(t *testing.T) {
	g := NewGint(3, 4)
	key := g.PramanaKey()
	if key != "3,1,4,1" {
		t.Errorf("PramanaKey() = %q, want \"3,1,4,1\"", key)
	}
	label := g.PramanaLabel()
	if label != "pra:num:3,1,4,1" {
		t.Errorf("PramanaLabel() = %q, want \"pra:num:3,1,4,1\"", label)
	}
	// UUID determinism
	id1 := g.PramanaID()
	id2 := g.PramanaID()
	if id1 != id2 {
		t.Error("PramanaID should be deterministic")
	}
	// Same value via Gint and Gauss should produce same ID
	gaussID := g.ToGauss().PramanaID()
	if id1 != gaussID {
		t.Error("Gint and Gauss PramanaID should match for same value")
	}
}

func TestGintNormsDivide(t *testing.T) {
	a := NewGint(2, 1) // norm = 5
	b := NewGint(3, 4) // norm = 25
	result := GintNormsDivide(a, b)
	if result == nil || result.Int64() != 5 {
		t.Errorf("NormsDivide(2+i, 3+4i) = %v, want 5", result)
	}
}

func TestGintCongruentModulo(t *testing.T) {
	a := NewGint(5, 3)
	b := NewGint(1, 1)
	c := NewGint(2, 1)
	isCong, _ := GintCongruentModulo(a, b, c)
	// (5+3i) - (1+i) = 4+2i; (4+2i)/(2+i) = (4+2i)(2-i)/5 = (8-4i+4i-2i²)/5 = 10/5 = 2
	if !isCong {
		t.Error("5+3i and 1+i should be congruent modulo 2+i")
	}
}

func TestGintParse(t *testing.T) {
	g, err := GintParse("3, 4")
	if err != nil {
		t.Fatal(err)
	}
	if g.Real().Cmp(big.NewInt(3)) != 0 || g.Imag().Cmp(big.NewInt(4)) != 0 {
		t.Errorf("GintParse(\"3, 4\") = %s", g)
	}
}

func TestGintTwo(t *testing.T) {
	two := GintTwo()
	if !two.Equal(NewGint(1, 1)) {
		t.Errorf("GintTwo() = %s, want 1+i", two)
	}
}

func TestGintComparison(t *testing.T) {
	a := NewGint(1, 0)
	b := NewGint(2, 0)
	if !a.Lt(b) {
		t.Error("1 should be < 2")
	}
	if !b.Gt(a) {
		t.Error("2 should be > 1")
	}
	if !a.Lte(a) {
		t.Error("1 should be <= 1")
	}
	if !a.Gte(a) {
		t.Error("1 should be >= 1")
	}
}
