go-exptable
===========

Compute tables to speed up Go's `math/big.Int` modular exponentiation
for fixed base.

```go
var table exptable.Table
table.Compute(&base, &m, 4)   // precomputes table for the given base
table.Exp(&result, &exponent) // sets result to base^exponent modulo m
```

[Documentation at godoc.](https://godoc.org/github.com/bwesterb/go-exptable)
