# pramana-go-sdk

Go SDK for the [Pramana](https://pramana-data.ca) knowledge graph. Provides exact-arithmetic value types, item model mapping, and data source connectors for working with Pramana data in Go.

## Status

**Pre-implementation** - Project structure and implementation plan documented. See [IMPLEMENTATION.md](IMPLEMENTATION.md) for the full design.

## Key Features (Planned)

- **GaussianRational** - Exact complex rational arithmetic (`a/b + (c/d)i`) with stdlib `math/big.Int`
- **Deterministic Pramana IDs** - UUID v5 generation matching the canonical Pramana web app
- **Minimal dependencies** - Only `google/uuid`; everything else is stdlib
- **Implicit interfaces** - Clean expression of Pramana's type hierarchy
- **Struct tag mapping** - `pramana:"property_name"` tags for ORM-style mapping
- **Multiple data sources** - `.pra` files, SPARQL, REST API, SQLite

## Installation (Future)

```bash
go get github.com/Emma-Leonhart/pramana-go-sdk
```

## Quick Example (Planned API)

```go
package main

import (
    "fmt"
    pramana "github.com/Emma-Leonhart/pramana-go-sdk"
)

func main() {
    half, _ := pramana.NewFromInt64(1, 2, 0, 1)   // 1/2
    third, _ := pramana.NewFromInt64(1, 3, 0, 1)  // 1/3
    result := half.Add(third)                       // 5/6

    fmt.Println(result.PramanaID())  // deterministic UUID v5
}
```

## Documentation

- [General SDK Specification](08_SDK_LIBRARY_SPECIFICATION.md) - Cross-language design spec
- [Go Implementation Guide](IMPLEMENTATION.md) - Go-specific implementation details

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
