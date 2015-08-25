// This file is distributed under the
// University of Illinois Open Source License.
// See LICENSE.TXT for details.

package sparsebitvector

import "fmt"

// KeyType defines SparseBitVector's keyspace.
type KeyType uint64

// SparseBitVector implementation based on that from LLVM:
// https://github.com/llvm-mirror/llvm/blob/master/include/llvm/ADT/SparseBitVector.h
type SparseBitVector struct {
	start   *element
	current *element
}

// New creates and instance of a SparseBitVector, optionally initialized by set.
func New(set ...KeyType) *SparseBitVector {
	result := new(SparseBitVector)
	for _, i := range set {
		result.Set(i)
	}
	return result
}

// Set sets a particular bit to true in a SparseBitVector.
func (sbv *SparseBitVector) Set(key KeyType) {
	index := key / elementsize
	nearest := sbv.search(index)

	if nearest == nil {
		e := sbv.create(index, nil, nil)
		e.Set(uint(key % elementsize))
	} else if nearest.index < index {
		e := sbv.create(index, nearest, nearest.next)
		e.Set(uint(key % elementsize))
		nearest.next = e
		if e.next != nil {
			e.next.prev = e
		}
	} else if nearest.index > index {
		e := sbv.create(index, nearest.prev, nearest)
		e.Set(uint(key % elementsize))
	} else {
		nearest.Set(uint(key % elementsize))
	}
}

// Unset sets a particular bit to false.
func (sbv *SparseBitVector) Unset(key KeyType) {
	index := key / elementsize
	e := sbv.search(index)
	if e == nil || e.index != index {
		return
	}

	e.Unset(uint(key % elementsize))
	if e.Count() == 0 {
		sbv.delete(e)
	}
}

// Clear sets all bits to false.
func (sbv *SparseBitVector) Clear() {
	sbv.start = nil
	sbv.current = nil
}

// Count returns the number of distinct bits that are true.
func (sbv *SparseBitVector) Count() int {
	total := 0
	for element := sbv.start; element != nil; element = element.next {
		total += element.Count()
	}
	return total
}

// Test checks whether a particular bit is true.
func (sbv *SparseBitVector) Test(key KeyType) bool {
	index := key / elementsize
	element := sbv.search(index)
	if element == nil || element.index != index {
		return false
	}
	return element.Test(uint(key % elementsize))
}

// TestAndSet checks whether a bit was previously true before setting it to true.
func (sbv *SparseBitVector) TestAndSet(key KeyType) bool {
	if sbv.Test(key) {
		return false
	}
	sbv.Set(key)
	return true
}

// Equals returns true iff sbv and sbv2 contain equivalent true bits.
func (sbv *SparseBitVector) Equals(sbv2 *SparseBitVector) bool {
	for e1, e2 := sbv.start, sbv2.start; e1 != nil || e2 != nil; e1, e2 = e1.next, e2.next {
		if e1 == nil || e2 == nil || e1.index != e2.index || !e1.Equals(&e2.FiniteBitVector) {
			return false
		}
	}
	return true
}

// Contains returns true iff sbv contains all of sbv2's true bits.
func (sbv *SparseBitVector) Contains(sbv2 *SparseBitVector) bool {
	for e1, e2 := sbv.start, sbv2.start; e2 != nil; e1, e2 = e1.next, e2.next {
		for e1 != nil && e1.index < e2.index {
			e1 = e1.next
		}
		if e1 == nil || e1.index != e2.index || !e1.Contains(&e2.FiniteBitVector) {
			return false
		}
	}
	return true
}

// UnionAndIntersectionSize returns the number of true bits of the union and intersection with sbv2.
func (sbv *SparseBitVector) UnionAndIntersectionSize(sbv2 *SparseBitVector) (int, int) {
	union := 0
	intersection := 0
	for e1, e2 := sbv.start, sbv2.start; e1 != nil || e2 != nil; {
		// sbv catch-up
		for e1 != nil && (e2 == nil || e1.index < e2.index) {
			union += e1.Count()
			e1 = e1.next
		}
		// sbv2 catch-up
		for e2 != nil && (e1 == nil || e2.index < e1.index) {
			union += e2.Count()
			e2 = e2.next
		}
		// same index
		if e1 != nil && e2 != nil && e1.index == e2.index {
			u, i := e1.UnionAndIntersectionSize(&e2.FiniteBitVector)
			union += u
			intersection += i
			e1 = e1.next
			e2 = e2.next
		}
	}
	return union, intersection
}

// UnionWith returns the number of true bits of the union and intersection with sbv2.
func (sbv *SparseBitVector) UnionWith(sbv2 *SparseBitVector) {
	for e1, e2 := sbv.start, sbv2.start; e1 != nil || e2 != nil; {
		// sbv catch-up
		for e1 != nil && (e2 == nil || e1.index < e2.index) {
			e1 = e1.next
		}
		// sbv2 catch-up
		for e2 != nil && (e1 == nil || e2.index < e1.index) {
			// insert element and copy data
			sbv.Set(e2.index * elementsize)
			e1 = sbv.search(e2.index)
			e1.FiniteBitVector = e2.FiniteBitVector
			e1 = e1.next
			e2 = e2.next
		}
		// same index
		if e1 != nil && e2 != nil && e1.index == e2.index {
			e1.UnionWith(&e2.FiniteBitVector)
			e1 = e1.next
			e2 = e2.next
		}
	}
}

// IntersectWith sets sbv to the intersection of itself and sbv2.
func (sbv *SparseBitVector) IntersectWith(sbv2 *SparseBitVector) {
	for e1, e2 := sbv.start, sbv2.start; e1 != nil; {
		// remove sbv elements not in sbv2
		for e1 != nil && (e2 == nil || e1.index < e2.index) {
			sbv.delete(e1)
			e1 = e1.next
		}
		// skip sbv2 elements not in sbv
		for e2 != nil && e1 != nil && e2.index < e1.index {
			e2 = e2.next
		}
		// same index
		if e1 != nil && e2 != nil && e1.index == e2.index {
			e1.IntersectWith(&e2.FiniteBitVector)
			e1 = e1.next
			e2 = e2.next
		}
	}
}

// IntersectWithComplement sets sbv to the intersection of itself and the inverse of sbv2.
func (sbv *SparseBitVector) IntersectWithComplement(sbv2 *SparseBitVector) {
	for e1, e2 := sbv.start, sbv2.start; e1 != nil; {
		// skip sbv elements not in sbv2
		for e1 != nil && (e2 == nil || e1.index < e2.index) {
			e1 = e1.next
		}
		// skip sbv2 elements not in sbv
		for e2 != nil && e1 != nil && e2.index < e1.index {
			e2 = e2.next
		}
		// same index
		if e1 != nil && e2 != nil && e1.index == e2.index {
			e1.IntersectWithComplement(&e2.FiniteBitVector)
			e1 = e1.next
			e2 = e2.next
		}
	}
}

// Iterate returns a channel which publishes all true bits in ascending order.
// The behaviour is undefined for bits modified while iterating.
func (sbv *SparseBitVector) Iterate() <-chan KeyType {
	c := make(chan KeyType)
	go func(c chan<- KeyType) {
		for e := sbv.start; e != nil; e = e.next {
			for i := e.FindNext(0); i != -1; i = e.FindNext(i + 1) {
				c <- e.index*elementsize + KeyType(i)
			}
		}
		close(c)
	}(c)
	return c
}

func (sbv *SparseBitVector) String() string {
	result := []KeyType{}
	for i := range sbv.Iterate() {
		result = append(result, i)
	}
	return fmt.Sprint(result)
}
