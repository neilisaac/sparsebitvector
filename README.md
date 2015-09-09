sparsebitvector
===============

[![TravisCI](https://travis-ci.org/neilisaac/sparsebitvector.svg)](https://travis-ci.org/neilisaac/sparsebitvector)
[![GoDoc](https://img.shields.io/badge/godoc-documentation-blue.svg)](https://godoc.org/github.com/neilisaac/sparsebitvector)

This library provides a SparseBitVector implementation in go based on
the similarly named class provided by LLVM.

See:

 * Documentation on [GoDoc](https://godoc.org/github.com/neilisaac/sparsebitvector)
 * LLVM's [SparseBitVector](https://github.com/llvm-mirror/llvm/blob/master/include/llvm/ADT/SparseBitVector.h)


### Usage

`go get github.com/neilisaac/sparsebitvector`

```
import "github.com/neilisaac/sparsebitvector"

vec := sparsebitvector.New(1, 2, 1000000)

vec.Set(1000001)
vec.Unset(1)

vec.IntersectWith(sparsebitvector.New(1, 1000001))

for value := range vec.Iterate() {
    fmt.Println("vec contains", value)
}
```

### Supported operations

 * `Set` set a bit to true
 * `Unset` set a bit to false (called `reset` in LLVM)
 * `Test` check whether a bit is true
 * `TestAndSet` set a bit to true and return true if it was changed
 * `Clear` set all bits to false
 * `Iterate` returns a channel that publishes all true bits
 * `Equals` compare to another SparseBitVector
 * `Contains` returns true if another SparseBitVector's bits are all true
 * `UnionAndIntersectionSize` return the size of the union and intersection with another SparseBitVector
 * `UnionSize`
 * `IntersectionSize`
 * `UnionWith` union itself with another SparseBitVector
 * `IntersectWith` intersect itself with another SparseBitVector
 * `IntersectWithComplement` intersect itself with the bitwise inverse of another SparseBitVector

### TODO

 * `Union` returning a new bit vector
 * `Intersection` returning a new bit vector
 * `Intersects` returning true if any bit is present in the intersection
