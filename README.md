# pramana-go-sdk

Go SDK for the [Pramana](https://pramana.dev) knowledge graph. Provides exact-arithmetic value types, item model mapping, and data source connectors for working with Pramana data in Go.

## Status

**Learning project** - This is the author's first Go project, being built as a way to learn the language through vibecoding. If you have experience with Go and want to help improve this SDK, contributions and feedback are very welcome!

The core Gaussian rational and integer types are implemented. See [IMPLEMENTATION.md](IMPLEMENTATION.md) for the full design.

## Key Features

- **Gauss** - Exact Gaussian rational arithmetic (`a/b + (c/d)i`) with stdlib `math/big.Int`
- **Gint** - Gaussian integer arithmetic with GCD, XGCD, primality testing, and modified division
- **Number Theory** - Primality testing with 6k+/-1 trial division
- **Deterministic Pramana IDs** - UUID v5 generation matching the canonical Pramana web app
- **Minimal dependencies** - Only `google/uuid`; everything else is stdlib

## Installation

```bash
go get github.com/Emma-Leonhart/pramana-go-sdk
```

## Quick Example

```go
package main

import (
    "fmt"
    pramana "github.com/Emma-Leonhart/pramana-go-sdk"
)

func main() {
    // Create Gaussian rationals
    half := pramana.NewGauss(1, 2, 0, 1)    // 1/2
    third := pramana.NewGauss(1, 3, 0, 1)   // 1/3
    sum := half.Add(third)                    // 5/6
    fmt.Println(sum)                          // "5/6"
    fmt.Println(sum.PramanaID())              // deterministic UUID v5

    // Create Gaussian integers
    z := pramana.NewGint(3, 4)               // 3 + 4i
    fmt.Println(z.Norm())                     // 25
    fmt.Println(pramana.GintIsGaussianPrime(z)) // false

    // GCD in Z[i]
    a := pramana.NewGint(3, 1)
    b := pramana.NewGint(1, 2)
    g := pramana.GintGCD(a, b)
    fmt.Println(g)
}
```

## Documentation

- [General SDK Specification](08_SDK_LIBRARY_SPECIFICATION.md) - Cross-language design spec
- [Go Implementation Guide](IMPLEMENTATION.md) - Go-specific implementation details

## Acknowledgments

The Gauss and Gint implementations across all Pramana SDKs were heavily inspired by [gaussian_integers](https://github.com/alreich/gaussian_integers) by **Alfred J. Reich, Ph.D.**, which provides exact arithmetic for Gaussian integers and Gaussian rationals in Python.

## Pramana SDK Family

| Language | Repository | Package |
|----------|-----------|---------|
| C# / .NET | [pramana-dotnet-sdk](https://github.com/Emma-Leonhart/pramana-dotnet-sdk) | `Pramana.SDK` (NuGet) |
| Python | [pramana-python-sdk](https://github.com/Emma-Leonhart/pramana-python-sdk) | `pramana-sdk` (PyPI) |
| TypeScript | [pramana-ts-sdk](https://github.com/Emma-Leonhart/pramana-ts-sdk) | `@pramana/sdk` (npm) |
| JavaScript | [pramana-js-sdk](https://github.com/Emma-Leonhart/pramana-js-sdk) | `@pramana/sdk` (npm) |
| Java | [pramana-java-sdk](https://github.com/Emma-Leonhart/pramana-java-sdk) | `org.pramana:pramana-sdk` (Maven) |
| Rust | [pramana-rust-sdk](https://github.com/Emma-Leonhart/pramana-rust-sdk) | `pramana-sdk` (crates.io) |
| Go | **pramana-go-sdk** (this repo) | `github.com/Emma-Leonhart/pramana-go-sdk` |

All SDKs implement the same core specification and must produce identical results for UUID v5 generation, canonical string normalization, and arithmetic operations.
