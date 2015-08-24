sparsebitvector
===============

[![TravisCI](https://travis-ci.org/neilisaac/sparsebitvector.svg)](https://travis-ci.org/neilisaac/sparsebitvector)

This library provides a SparseBitVector implementation in go based on
the similarly named class provided by LLVM.

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
 * `UnionWith` union itself with another SparseBitVector

### In progress

 * `IntersectWith` intersect itself with another SparseBitVector
 * `IntersectWithComplement` intersect itself with the bitwise inverse of another SparseBitVector

### TODO

 * `Union` returning a new bit vector
 * `UnionSize` redundant; use `UnionAndIntersectionSize`
 * `Intersection` returning a new bit vector
 * `Intersects` returning true if any bit is present in the intersection
 * `IntersectionSize` could be faster than `UnionAndIntersectionSize`
