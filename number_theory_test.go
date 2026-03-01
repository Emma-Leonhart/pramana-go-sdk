package pramana

import (
	"math/big"
	"testing"
)

func TestIsPrime(t *testing.T) {
	tests := []struct {
		n    int64
		want bool
	}{
		{0, false},
		{1, false},
		{2, true},
		{3, true},
		{4, false},
		{5, true},
		{6, false},
		{7, true},
		{11, true},
		{13, true},
		{15, false},
		{17, true},
		{19, true},
		{23, true},
		{25, false},
		{29, true},
		{97, true},
		{100, false},
		{101, true},
	}

	for _, tt := range tests {
		got := IsPrime(big.NewInt(tt.n))
		if got != tt.want {
			t.Errorf("IsPrime(%d) = %v, want %v", tt.n, got, tt.want)
		}
	}
}

func TestIsPrimeInt(t *testing.T) {
	if !IsPrimeInt(7) {
		t.Error("IsPrimeInt(7) should be true")
	}
	if IsPrimeInt(8) {
		t.Error("IsPrimeInt(8) should be false")
	}
}

func TestIsPrimeNegative(t *testing.T) {
	if IsPrime(big.NewInt(-5)) {
		t.Error("IsPrime(-5) should be false")
	}
}
