package pramana

import "math/big"

// IsPrime returns true if n is a prime number.
// Uses trial division with the 6k±1 optimization.
func IsPrime(n *big.Int) bool {
	zero := big.NewInt(0)
	one := big.NewInt(1)
	two := big.NewInt(2)
	three := big.NewInt(3)

	if n.Cmp(two) < 0 {
		return false
	}
	if n.Cmp(two) == 0 || n.Cmp(three) == 0 {
		return true
	}

	mod2 := new(big.Int).Mod(n, two)
	mod3 := new(big.Int).Mod(n, three)
	if mod2.Cmp(zero) == 0 || mod3.Cmp(zero) == 0 {
		return false
	}

	i := big.NewInt(5)
	iSquared := new(big.Int)
	for {
		iSquared.Mul(i, i)
		if iSquared.Cmp(n) > 0 {
			break
		}
		mod := new(big.Int).Mod(n, i)
		if mod.Cmp(zero) == 0 {
			return false
		}
		iPlus2 := new(big.Int).Add(i, two)
		mod = new(big.Int).Mod(n, iPlus2)
		if mod.Cmp(zero) == 0 {
			return false
		}
		i.Add(i, big.NewInt(6))
	}

	return true
}

// IsPrimeInt is a convenience wrapper for IsPrime that accepts an int64.
func IsPrimeInt(n int64) bool {
	return IsPrime(big.NewInt(n))
}

// gcd returns the greatest common divisor of a and b.
// Both a and b should be non-negative.
func gcd(a, b *big.Int) *big.Int {
	x := new(big.Int).Abs(a)
	y := new(big.Int).Abs(b)
	return new(big.Int).GCD(nil, nil, x, y)
}

// abs returns the absolute value of x as a new big.Int.
func absBigInt(x *big.Int) *big.Int {
	return new(big.Int).Abs(x)
}

// bigIntSign returns -1, 0, or 1 for the sign of x.
func bigIntSign(x *big.Int) int {
	return x.Sign()
}

// copyBigInt returns a deep copy of x.
func copyBigInt(x *big.Int) *big.Int {
	if x == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Set(x)
}
