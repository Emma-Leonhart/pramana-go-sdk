package pramana

import (
	"math"
	"testing"
)

func TestGaussConstruction(t *testing.T) {
	g := NewGauss(3, 2, 1, 4)
	if g.A().Int64() != 3 || g.B().Int64() != 2 || g.C().Int64() != 1 || g.D().Int64() != 4 {
		t.Errorf("NewGauss(3,2,1,4) = %s, want 3/2 + 1/4i", g.ToRawString())
	}
}

func TestGaussNormalization(t *testing.T) {
	// 2/4 should normalize to 1/2
	g := NewGauss(2, 4, 0, 1)
	if g.A().Int64() != 1 || g.B().Int64() != 2 {
		t.Errorf("NewGauss(2,4,0,1) should normalize to 1/2, got %s", g.ToRawString())
	}

	// Negative denominator should be flipped
	g2 := NewGauss(3, -2, 0, 1)
	if g2.A().Int64() != -3 || g2.B().Int64() != 2 {
		t.Errorf("NewGauss(3,-2,0,1) should normalize to -3/2, got %s", g2.ToRawString())
	}

	// Zero numerator normalizes to 0/1
	g3 := NewGauss(0, 5, 0, 3)
	if g3.A().Int64() != 0 || g3.B().Int64() != 1 {
		t.Errorf("NewGauss(0,5,0,3) should normalize to 0/1, got %s", g3.ToRawString())
	}
}

func TestGaussProperties(t *testing.T) {
	zero := GaussZero()
	one := GaussOne()
	i := GaussI()
	half := NewGauss(1, 2, 0, 1)

	if !zero.IsZero() {
		t.Error("Zero should be zero")
	}
	if !one.IsOne() {
		t.Error("One should be one")
	}
	if !one.IsReal() {
		t.Error("One should be real")
	}
	if !i.IsPurelyImaginary() {
		t.Error("i should be purely imaginary")
	}
	if !one.IsInteger() {
		t.Error("1 should be an integer")
	}
	if half.IsInteger() {
		t.Error("1/2 should not be an integer")
	}
	if !i.IsGaussianInteger() {
		t.Error("i should be a Gaussian integer")
	}
	if !one.IsPositive() {
		t.Error("1 should be positive")
	}
	if !GaussMinusOne().IsNegative() {
		t.Error("-1 should be negative")
	}
}

func TestGaussArithmetic(t *testing.T) {
	a := NewGaussInt(3, 4)
	b := NewGaussInt(1, 2)

	// Addition
	sum := a.Add(b)
	expected := NewGaussInt(4, 6)
	if !sum.Equal(expected) {
		t.Errorf("(3+4i) + (1+2i) = %s, want 4+6i", sum)
	}

	// Subtraction
	diff := a.Sub(b)
	expected = NewGaussInt(2, 2)
	if !diff.Equal(expected) {
		t.Errorf("(3+4i) - (1+2i) = %s, want 2+2i", diff)
	}

	// Multiplication: (3+4i)(1+2i) = 3+6i+4i+8i² = 3+10i-8 = -5+10i
	prod := a.Mul(b)
	expected = NewGaussInt(-5, 10)
	if !prod.Equal(expected) {
		t.Errorf("(3+4i) * (1+2i) = %s, want -5+10i", prod)
	}

	// Negation
	neg := a.Neg()
	expected = NewGaussInt(-3, -4)
	if !neg.Equal(expected) {
		t.Errorf("-(3+4i) = %s, want -3-4i", neg)
	}
}

func TestGaussDivision(t *testing.T) {
	a := NewGaussInt(1, 0)
	b := NewGaussInt(0, 1)

	// 1/i = -i
	result := a.Div(b)
	expected := NewGaussInt(0, -1)
	if !result.Equal(expected) {
		t.Errorf("1/i = %s, want -i", result)
	}
}

func TestGaussConjugate(t *testing.T) {
	g := NewGaussInt(3, 4)
	conj := g.Conjugate()
	expected := NewGaussInt(3, -4)
	if !conj.Equal(expected) {
		t.Errorf("conj(3+4i) = %s, want 3-4i", conj)
	}
}

func TestGaussNorm(t *testing.T) {
	g := NewGaussInt(3, 4)
	norm := g.Norm()
	expected := NewGaussReal(25)
	if !norm.Equal(expected) {
		t.Errorf("|3+4i|² = %s, want 25", norm)
	}
}

func TestGaussMagnitude(t *testing.T) {
	g := NewGaussInt(3, 4)
	mag := g.Magnitude()
	if math.Abs(mag-5.0) > 1e-10 {
		t.Errorf("|3+4i| = %f, want 5.0", mag)
	}
}

func TestGaussPow(t *testing.T) {
	i := GaussI()
	// i² = -1
	result := GaussPow(i, 2)
	if !result.Equal(GaussMinusOne()) {
		t.Errorf("i² = %s, want -1", result)
	}

	// i⁴ = 1
	result = GaussPow(i, 4)
	if !result.Equal(GaussOne()) {
		t.Errorf("i⁴ = %s, want 1", result)
	}
}

func TestGaussUnits(t *testing.T) {
	units := GaussUnits()
	if len(units) != 4 {
		t.Errorf("GaussUnits() has %d elements, want 4", len(units))
	}
}

func TestGaussAssociates(t *testing.T) {
	g := NewGaussInt(2, 3)
	assocs := g.Associates()
	if len(assocs) != 3 {
		t.Errorf("Associates() has %d elements, want 3", len(assocs))
	}
	// Check that g is an associate of each
	for _, a := range assocs {
		if !g.IsAssociate(a) {
			t.Errorf("%s should be associate of %s", a, g)
		}
	}
}

func TestGaussInverse(t *testing.T) {
	g := NewGaussInt(1, 1)
	inv := g.Inverse()
	// 1/(1+i) = (1-i)/2 = 1/2 - 1/2*i
	expected := NewGauss(1, 2, -1, 2)
	if !inv.Equal(expected) {
		t.Errorf("1/(1+i) = %s, want 1/2-1/2i", inv)
	}
}

func TestGaussString(t *testing.T) {
	tests := []struct {
		g    *Gauss
		want string
	}{
		{GaussZero(), "0"},
		{GaussOne(), "1"},
		{GaussMinusOne(), "-1"},
		{GaussI(), "i"},
		{NewGaussInt(0, -1), "-i"},
		{NewGaussInt(3, 4), "3 + 4i"},
		{NewGaussInt(3, -4), "3 - 4i"},
		{NewGauss(1, 2, 3, 4), "1/2 + 3/4i"},
	}
	for _, tt := range tests {
		got := tt.g.String()
		if got != tt.want {
			t.Errorf("%s.String() = %q, want %q", tt.g.ToRawString(), got, tt.want)
		}
	}
}

func TestGaussParse(t *testing.T) {
	g, err := GaussParse("3,2,1,4")
	if err != nil {
		t.Fatal(err)
	}
	if g.A().Int64() != 3 || g.B().Int64() != 2 || g.C().Int64() != 1 || g.D().Int64() != 4 {
		t.Errorf("GaussParse(\"3,2,1,4\") = %s", g.ToRawString())
	}
}

func TestGaussFromPramana(t *testing.T) {
	g, err := GaussFromPramana("pra:num:3,2,1,4")
	if err != nil {
		t.Fatal(err)
	}
	if g.A().Int64() != 3 || g.B().Int64() != 2 {
		t.Errorf("GaussFromPramana failed: %s", g.ToRawString())
	}
}

func TestGaussPramanaIdentity(t *testing.T) {
	g := NewGaussInt(3, 4)
	key := g.PramanaKey()
	if key != "3,1,4,1" {
		t.Errorf("PramanaKey() = %q, want \"3,1,4,1\"", key)
	}
	label := g.PramanaLabel()
	if label != "pra:num:3,1,4,1" {
		t.Errorf("PramanaLabel() = %q, want \"pra:num:3,1,4,1\"", label)
	}
	// UUID should be deterministic
	id1 := g.PramanaID()
	id2 := g.PramanaID()
	if id1 != id2 {
		t.Error("PramanaID should be deterministic")
	}
}

func TestGaussComparison(t *testing.T) {
	a := NewGaussReal(3)
	b := NewGaussReal(5)
	if !a.Lt(b) {
		t.Error("3 should be < 5")
	}
	if !b.Gt(a) {
		t.Error("5 should be > 3")
	}
	if !a.Lte(a) {
		t.Error("3 should be <= 3")
	}
}

func TestGaussFloorCeiling(t *testing.T) {
	g := NewGauss(7, 2, 5, 3)
	floor := GaussFloor(g)
	// 7/2 = 3.5 -> floor = 3, 5/3 = 1.666 -> floor = 1
	if !floor.Equal(NewGaussInt(3, 1)) {
		t.Errorf("Floor(7/2 + 5/3i) = %s, want 3+i", floor)
	}

	ceil := GaussCeiling(g)
	// ceil(3.5) = 4, ceil(1.666) = 2
	if !ceil.Equal(NewGaussInt(4, 2)) {
		t.Errorf("Ceiling(7/2 + 5/3i) = %s, want 4+2i", ceil)
	}
}

func TestGaussMinMax(t *testing.T) {
	a := NewGaussReal(3)
	b := NewGaussReal(7)
	if !GaussMin(a, b).Equal(a) {
		t.Error("Min(3,7) should be 3")
	}
	if !GaussMax(a, b).Equal(b) {
		t.Error("Max(3,7) should be 7")
	}
}

func TestGaussDecimalString(t *testing.T) {
	g := NewGauss(1, 3, 0, 1)
	s := g.ToDecimalString(4)
	if s != "0.3333" {
		t.Errorf("ToDecimalString(4) = %q, want \"0.3333\"", s)
	}
}
