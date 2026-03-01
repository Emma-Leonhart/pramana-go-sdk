# Pramana Go SDK - Implementation Guide

**Package name:** `github.com/Emma-Leonhart/pramana-go-sdk`
**Minimum Go:** 1.21+
**Reference implementation:** [PramanaLib (C#)](https://github.com/Emma-Leonhart/PramanaLib)

---

## 1. Project Structure

```
pramana-go-sdk/
├── go.mod
├── go.sum
├── gaussian_rational.go            # GaussianRational implementation
├── gaussian_rational_test.go       # GaussianRational tests
├── pramana_id.go                   # UUID v5 generation utilities
├── pramana_id_test.go
├── number_type.go                  # NumberType constants
├── item.go                         # PramanaItem interface + types
├── entity.go                       # PramanaEntity struct
├── property.go                     # PramanaProperty struct
├── proposition.go                  # PramanaProposition struct
├── sense.go                        # PramanaSense struct
├── graph.go                        # PramanaGraph (loading, serialization)
├── graph_test.go
├── config.go                       # PramanaConfig
├── errors.go                       # Error types
├── orm/
│   ├── mapper.go                   # Entity mapping via struct tags
│   ├── mapper_test.go
│   ├── query.go                    # Query builder
│   └── query_test.go
├── datasources/
│   ├── prafile.go                  # .pra JSON file reader
│   ├── sparql.go                   # GraphDB SPARQL connector
│   ├── restapi.go                  # Pramana REST API connector
│   └── sqlite.go                   # SQLite export reader
├── structs/
│   ├── date.go                     # date: pseudo-class
│   ├── time.go                     # time: pseudo-class
│   ├── interval.go                 # interval: pseudo-class
│   ├── coordinate.go               # coord: pseudo-class
│   └── chemical.go                 # chem: / element: pseudo-classes
├── testdata/
│   └── test_vectors.json           # Cross-language test vectors
└── docs/
    └── api.md
```

## 2. Build & Packaging

### go.mod

```go
module github.com/Emma-Leonhart/pramana-go-sdk

go 1.21

require (
    github.com/google/uuid v1.6.0
)
```

### Key decisions:
- **Flat package structure** — Go convention, no nested `src/` directory
- **Minimal dependencies** — only `google/uuid` for UUID v5
- **`math/big`** from stdlib for arbitrary precision integers (no external BigInt)
- **Struct tags** for ORM mapping (Go convention, like `json:` tags)
- **Interfaces** for data source abstraction (Go's implicit interface satisfaction)
- **No generics for item types** — use type assertions (idiomatic pre-1.18 Go style) or generics where cleaner

## 3. GaussianRational (Gauss) Implementation

> **Naming convention:** The standard short name for a Gaussian rational is **`Gauss`**. When referring specifically to a Gaussian integer (both denominators are 1), the standard short name is **`Gint`**.

### 3.1 Struct Design

Go has `math/big.Int` in the standard library for arbitrary precision integers. No operator overloading; named methods throughout.

```go
package pramana

import (
    "fmt"
    "math/big"
    "strings"

    "github.com/google/uuid"
)

// GaussianRational represents an exact complex rational number: a/b + (c/d)i.
// Immutable by convention — all methods return new values.
type GaussianRational struct {
    a *big.Int // real numerator
    b *big.Int // real denominator (positive, nonzero)
    c *big.Int // imaginary numerator
    d *big.Int // imaginary denominator (positive, nonzero)
}

// New creates a GaussianRational from four big.Int values and normalizes.
func New(a, b, c, d *big.Int) (*GaussianRational, error) {
    if b.Sign() <= 0 || d.Sign() <= 0 {
        return nil, fmt.Errorf("denominators must be positive integers")
    }
    // Normalize to canonical form
    gReal := new(big.Int).GCD(nil, nil, new(big.Int).Abs(a), b)
    gImag := new(big.Int).GCD(nil, nil, new(big.Int).Abs(c), d)
    return &GaussianRational{
        a: new(big.Int).Div(a, gReal),
        b: new(big.Int).Div(b, gReal),
        c: new(big.Int).Div(c, gImag),
        d: new(big.Int).Div(d, gImag),
    }, nil
}

// NewFromInt64 creates a GaussianRational from four int64 values.
func NewFromInt64(a, b, c, d int64) (*GaussianRational, error) {
    return New(
        big.NewInt(a), big.NewInt(b),
        big.NewInt(c), big.NewInt(d),
    )
}
```

### 3.2 Constructor Functions

```go
// FromInt creates a GaussianRational from a single integer (imaginary = 0).
func FromInt(value int64) *GaussianRational {
    g, _ := NewFromInt64(value, 1, 0, 1)
    return g
}

// FromBigInt creates a GaussianRational from a big.Int (imaginary = 0).
func FromBigInt(value *big.Int) *GaussianRational {
    g, _ := New(value, big.NewInt(1), big.NewInt(0), big.NewInt(1))
    return g
}

// FromComplex creates a GaussianRational from integer real and imaginary parts.
func FromComplex(real, imag int64) *GaussianRational {
    g, _ := NewFromInt64(real, 1, imag, 1)
    return g
}

// Parse parses a canonical "a,b,c,d" or "num:a,b,c,d" string.
func Parse(s string) (*GaussianRational, error) {
    s = strings.TrimPrefix(s, "num:")
    parts := strings.Split(s, ",")
    if len(parts) != 4 {
        return nil, fmt.Errorf("expected 4 comma-separated integers: %s", s)
    }
    a, ok1 := new(big.Int).SetString(strings.TrimSpace(parts[0]), 10)
    b, ok2 := new(big.Int).SetString(strings.TrimSpace(parts[1]), 10)
    c, ok3 := new(big.Int).SetString(strings.TrimSpace(parts[2]), 10)
    d, ok4 := new(big.Int).SetString(strings.TrimSpace(parts[3]), 10)
    if !ok1 || !ok2 || !ok3 || !ok4 {
        return nil, fmt.Errorf("invalid integer in: %s", s)
    }
    return New(a, b, c, d)
}
```

### 3.3 Arithmetic Methods (No Operator Overloading)

Go does not support operator overloading. Methods follow Go naming conventions (exported, PascalCase):

```go
// Add returns the sum of g and other.
func (g *GaussianRational) Add(other *GaussianRational) *GaussianRational {
    // a1/b1 + a2/b2 = (a1*b2 + a2*b1) / (b1*b2)
    realNum := new(big.Int).Add(
        new(big.Int).Mul(g.a, other.b),
        new(big.Int).Mul(other.a, g.b),
    )
    realDen := new(big.Int).Mul(g.b, other.b)
    imagNum := new(big.Int).Add(
        new(big.Int).Mul(g.c, other.d),
        new(big.Int).Mul(other.c, g.d),
    )
    imagDen := new(big.Int).Mul(g.d, other.d)
    result, _ := New(realNum, realDen, imagNum, imagDen)
    return result
}

// Sub returns the difference of g and other.
func (g *GaussianRational) Sub(other *GaussianRational) *GaussianRational { ... }

// Neg returns the negation of g.
func (g *GaussianRational) Neg() *GaussianRational { ... }

// Mul returns the product of g and other.
func (g *GaussianRational) Mul(other *GaussianRational) *GaussianRational {
    // (a+bi)(c+di) = (ac-bd) + (ad+bc)i
    ...
}

// Div returns the quotient of g and other.
func (g *GaussianRational) Div(other *GaussianRational) *GaussianRational { ... }

// Mod returns the modulo (real values only).
func (g *GaussianRational) Mod(other *GaussianRational) (*GaussianRational, error) {
    if !g.IsReal() || !other.IsReal() {
        return nil, fmt.Errorf("modulo only defined for real values")
    }
    ...
}

// Pow returns g raised to an integer exponent.
func (g *GaussianRational) Pow(exp int) *GaussianRational { ... }
```

### 3.4 Comparison Methods

```go
// Equal returns true if g and other represent the same value.
func (g *GaussianRational) Equal(other *GaussianRational) bool {
    return g.a.Cmp(other.a) == 0 && g.b.Cmp(other.b) == 0 &&
           g.c.Cmp(other.c) == 0 && g.d.Cmp(other.d) == 0
}

// Cmp compares g and other (real values only). Returns -1, 0, or 1.
func (g *GaussianRational) Cmp(other *GaussianRational) (int, error) {
    if !g.IsReal() || !other.IsReal() {
        return 0, fmt.Errorf("ordering only defined for real values")
    }
    // Compare a1/b1 vs a2/b2 via cross-multiplication
    lhs := new(big.Int).Mul(g.a, other.b)
    rhs := new(big.Int).Mul(other.a, g.b)
    return lhs.Cmp(rhs), nil
}

// Lt returns true if g < other (real values only).
func (g *GaussianRational) Lt(other *GaussianRational) (bool, error) {
    cmp, err := g.Cmp(other)
    return cmp < 0, err
}

// Gt, Lte, Gte — same pattern
```

### 3.5 Properties

```go
func (g *GaussianRational) IsReal() bool            { return g.c.Sign() == 0 }
func (g *GaussianRational) IsInteger() bool          { return g.IsReal() && g.b.Cmp(big.NewInt(1)) == 0 }
func (g *GaussianRational) IsGaussianInteger() bool  { return g.b.Cmp(big.NewInt(1)) == 0 && g.d.Cmp(big.NewInt(1)) == 0 }
func (g *GaussianRational) IsZero() bool             { return g.a.Sign() == 0 && g.c.Sign() == 0 }
func (g *GaussianRational) IsPositive() bool         { return g.IsReal() && g.a.Sign() > 0 }
func (g *GaussianRational) IsNegative() bool         { return g.IsReal() && g.a.Sign() < 0 }

func (g *GaussianRational) Conjugate() *GaussianRational {
    result, _ := New(g.a, g.b, new(big.Int).Neg(g.c), g.d)
    return result
}

func (g *GaussianRational) MagnitudeSquared() *GaussianRational { ... }
func (g *GaussianRational) RealPart() *GaussianRational         { ... }
func (g *GaussianRational) ImaginaryPart() *GaussianRational    { ... }
func (g *GaussianRational) Reciprocal() *GaussianRational       { ... }

func (g *GaussianRational) RealNumerator() *big.Int      { return new(big.Int).Set(g.a) }
func (g *GaussianRational) RealDenominator() *big.Int    { return new(big.Int).Set(g.b) }
func (g *GaussianRational) ImagNumerator() *big.Int      { return new(big.Int).Set(g.c) }
func (g *GaussianRational) ImagDenominator() *big.Int    { return new(big.Int).Set(g.d) }

func (g *GaussianRational) Classify() NumberType { ... }
```

### 3.6 NumberType

```go
type NumberType string

const (
    NaturalNumber    NumberType = "Natural Number"
    WholeNumber      NumberType = "Whole Number"
    Integer          NumberType = "Integer"
    RationalNumber   NumberType = "Rational Number"
    GaussianRational NumberType = "Gaussian Rational"
)
```

### 3.7 Pramana ID (UUID v5)

```go
import "github.com/google/uuid"

var NumNamespace = uuid.MustParse("a6613321-e9f6-4348-8f8b-29d2a3c86349")

// Canonical returns the canonical num: string representation.
func (g *GaussianRational) Canonical() string {
    return fmt.Sprintf("num:%s,%s,%s,%s", g.a, g.b, g.c, g.d)
}

// PramanaID returns the deterministic UUID v5 for this value.
func (g *GaussianRational) PramanaID() uuid.UUID {
    return uuid.NewSHA1(NumNamespace, []byte(g.Canonical()))
}

// PramanaURI returns the Pramana URI for this value.
func (g *GaussianRational) PramanaURI() string {
    return "pra:" + g.PramanaID().String()
}
```

The `google/uuid` package provides `uuid.NewSHA1()` which implements UUID v5.

### 3.8 Stringer Interface

```go
// String implements fmt.Stringer.
func (g *GaussianRational) String() string {
    return g.Canonical()
}

func (g *GaussianRational) ToMixed() string    { ... }
func (g *GaussianRational) ToImproper() string { ... }
func (g *GaussianRational) ToRaw() string {
    return fmt.Sprintf("<%s,%s,%s,%s>", g.a, g.b, g.c, g.d)
}
```

### 3.9 JSON Marshaling

```go
import "encoding/json"

func (g *GaussianRational) MarshalJSON() ([]byte, error) {
    return json.Marshal(g.Canonical())
}

func (g *GaussianRational) UnmarshalJSON(data []byte) error {
    var s string
    if err := json.Unmarshal(data, &s); err != nil {
        return err
    }
    parsed, err := Parse(s)
    if err != nil {
        return err
    }
    *g = *parsed
    return nil
}
```

### 3.10 Intentionally Unsupported

```go
// Magnitude panics because complex magnitude produces irrationals.
func (g *GaussianRational) Magnitude() float64 {
    panic("complex magnitude produces irrationals; use MagnitudeSquared() for exact result")
}
// Phase, ToPolar, Sqrt — same treatment
```

## 4. Item Model

### 4.1 Interface-Based Design

Go's implicit interface satisfaction is ideal for the Pramana item hierarchy:

```go
type ItemType string

const (
    ItemTypeEntity      ItemType = "Entity"
    ItemTypeProperty    ItemType = "Property"
    ItemTypeProposition ItemType = "Proposition"
    ItemTypeSense       ItemType = "Sense"
    ItemTypeEvidence    ItemType = "Evidence"
    ItemTypeStanceLink  ItemType = "StanceLink"
)

// PramanaItem is the base interface for all Pramana items.
type PramanaItem interface {
    UUID() uuid.UUID
    Type() ItemType
    Properties() map[string]interface{}
    Edges() map[string]uuid.UUID
}
```

### 4.2 Typed Structs

```go
type PramanaEntity struct {
    id         uuid.UUID
    label      string
    instanceOf *uuid.UUID
    subclassOf *uuid.UUID
    props      map[string]interface{}
    edgeMap    map[string]uuid.UUID
}

func (e *PramanaEntity) UUID() uuid.UUID                  { return e.id }
func (e *PramanaEntity) Type() ItemType                   { return ItemTypeEntity }
func (e *PramanaEntity) Label() string                    { return e.label }
func (e *PramanaEntity) Properties() map[string]interface{} { return e.props }
func (e *PramanaEntity) Edges() map[string]uuid.UUID      { return e.edgeMap }
```

## 5. ORM-Style Mapping

### 5.1 Struct Tags (Go Convention)

```go
type ShintoShrine struct {
    Coordinates *Coordinate `pramana:"coordinates"`
    WikidataID  *string     `pramana:"Wikidata ID"`
    PartOf      *uuid.UUID  `pramana:"part of,ref"`
}
```

### 5.2 Mapper Using Reflection

```go
import "reflect"

func MapEntity(item PramanaItem, target interface{}) error {
    v := reflect.ValueOf(target).Elem()
    t := v.Type()

    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        tag := field.Tag.Get("pramana")
        if tag == "" {
            continue
        }
        // Parse tag, look up property in item, set field value
        ...
    }
    return nil
}
```

### 5.3 Query Interface

```go
shrines, err := pramana.Query("ShintoShrine").
    Where("coordinates", NotNull).
    Limit(100).
    All()

water, err := pramana.GetByID(
    uuid.MustParse("00000007-0000-4000-8000-000000000007"),
)
```

### 5.4 Interfaces for Multiple Classification

Go interfaces are implicit — no `implements` keyword needed:

```go
type ChemicalCompound interface {
    MolecularFormula() string
}

type QuantumSubstance interface {
    QuantumState() string
}

// Water satisfies both interfaces automatically
type Water struct { ... }
func (w *Water) MolecularFormula() string { return "H2O" }
func (w *Water) QuantumState() string     { return "..." }
```

## 6. Data Sources

| Source | File | Dependency |
|--------|------|------------|
| `.pra` JSON file | `datasources/prafile.go` | None (stdlib `encoding/json`) |
| GraphDB SPARQL | `datasources/sparql.go` | None (stdlib `net/http`) |
| Pramana REST API | `datasources/restapi.go` | None (stdlib `net/http`) |
| SQLite export | `datasources/sqlite.go` | `github.com/mattn/go-sqlite3` |

Go's standard library includes `net/http` and `encoding/json`, so most connectors need zero external dependencies.

## 7. Go-Specific Considerations

### 7.1 big.Int Mutability
`math/big.Int` is mutable. Always create new instances to avoid aliasing:

```go
// WRONG — mutates original
result := g.a
result.Add(result, other.a) // Mutates g.a!

// CORRECT — create new
result := new(big.Int).Add(g.a, other.a)
```

### 7.2 Error Handling
Go uses explicit error returns instead of exceptions:

```go
// Operations that can fail return (value, error)
result, err := g.Mod(other)
if err != nil {
    // Handle: modulo not defined for complex
}

// Operations that can't fail return value only
result := g.Add(other)
```

### 7.3 No Generics for Pre-1.18 Compatibility
The query builder uses `interface{}` for type flexibility. With Go 1.21+ as minimum, generics are available but used sparingly to stay idiomatic.

### 7.4 No Operator Overloading
Similar to Java/TypeScript. Method chaining helps readability:

```go
// Instead of: result = a + b * c
result := a.Add(b.Mul(c))
```

### 7.5 Implicit Interfaces
Go's biggest advantage for the Pramana type system — any struct that has the right methods automatically satisfies an interface. No explicit `implements` declaration needed.

### 7.6 Table-Driven Tests
Go convention for comprehensive test coverage:

```go
func TestGaussianRationalAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     *GaussianRational
        expected *GaussianRational
    }{
        {"1/2 + 1/3", NewFromInt64(1,2,0,1), NewFromInt64(1,3,0,1), NewFromInt64(5,6,0,1)},
        // more test cases...
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := tt.a.Add(tt.b)
            if !result.Equal(tt.expected) {
                t.Errorf("got %s, want %s", result, tt.expected)
            }
        })
    }
}
```

## 8. Testing Strategy

```bash
go test ./...                 # Run all tests
go test -v ./...              # Verbose output
go test -race ./...           # Race detector
go test -bench=. ./...        # Run benchmarks
go vet ./...                  # Static analysis
```

## 9. Implementation Priority

### Phase 1 - GaussianRational (core)
1. Implement `GaussianRational` struct with `math/big.Int` components
2. Implement all arithmetic methods (`Add`, `Sub`, `Mul`, `Div`, `Mod`, `Pow`)
3. Implement `Equal`, `Cmp`, `Lt`, `Gt`, `Lte`, `Gte`
4. Implement UUID v5 via `google/uuid` package
5. Implement `Parse`, `String`, formatting methods, `Classify`
6. Implement JSON marshal/unmarshal
7. Write table-driven tests against cross-language test vectors

### Phase 2 - Base Item Model
1. Define `PramanaItem` interface and typed structs
2. Implement `PramanaGraph` with JSON serialization
3. Implement `.pra` file reader

### Phase 3 - ORM Mapping
1. Implement struct-tag-based entity mapping via reflection
2. Implement query builder
3. Implement `PramanaConfig`

### Phase 4 - Data Sources & Provenance
1. SPARQL connector (stdlib `net/http`)
2. REST API connector (stdlib `net/http`)
3. SQLite connector (`go-sqlite3`)
4. Provenance metadata

### Phase 5 - Pseudo-Classes
1. `PramanaDate`, `PramanaTime`, `PramanaInterval` (stdlib `time` package)
2. `Coordinate` struct
3. `ChemicalIdentifier` / `ChemicalElement`

## 10. Go-Specific Advantages

- **`math/big.Int` in stdlib** — no external BigInt dependency
- **`net/http` in stdlib** — no external HTTP library for SPARQL/REST
- **`encoding/json` in stdlib** — no external JSON library
- **Implicit interfaces** — cleanest expression of Pramana's type hierarchy
- **Goroutines** — efficient concurrent graph traversal
- **Single binary compilation** — easy deployment
- **`google/uuid`** — well-maintained UUID v5 implementation
- **Table-driven tests** — comprehensive, readable test suites
