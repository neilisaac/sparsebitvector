// This file is distributed under the
// University of Illinois Open Source License.
// See LICENSE.TXT for details.

package sparsebitvector

// KeyType defines SparseBitVector's keyspace.
type KeyType uint64

// SparseBitVector implementation based on that from LLVM:
// https://github.com/llvm-mirror/llvm/blob/master/include/llvm/ADT/SparseBitVector.h
type SparseBitVector struct {
	start   *element
	current *element
}

type element struct {
	FiniteBitVector
	index KeyType
	prev  *element
	next  *element
}

func (sbv *SparseBitVector) insert(index KeyType, prev, next *element) *element {
	element := &element{index: index, next: next, prev: prev}

	if prev == nil {
		sbv.start = element
	} else {
		prev.next = element
	}

	if next != nil {
		next.prev = element
	}

	sbv.current = element
	return element
}

func (sbv *SparseBitVector) search(index KeyType) *element {
	if sbv.current == nil {
		if sbv.start == nil {
			return nil
		}
		sbv.current = sbv.start
	}

	if sbv.current.index > index {
		for element := sbv.current; element != nil; element = element.prev {
			sbv.current = element
			if element.index == index {
				break
			}
		}
	} else if sbv.current.index < index {
		for element := sbv.current; element != nil; element = element.next {
			sbv.current = element
			if element.index == index {
				break
			}
		}
	}

	return sbv.current
}

// New creates and instance of a SparseBitVector.
func New() *SparseBitVector {
	return new(SparseBitVector)
}

// Set sets a particular bit to true in a SparseBitVector.
func (sbv *SparseBitVector) Set(key KeyType) {
	index := key / elementsize
	nearest := sbv.search(index)

	if nearest == nil {
		element := sbv.insert(index, nil, nil)
		element.Set(uint(key % elementsize))
	} else if nearest.index < index {
		element := sbv.insert(index, nearest, nearest.next)
		element.Set(uint(key % elementsize))
		nearest.next = element
		if element.next != nil {
			element.next.prev = element
		}
	} else if nearest.index > index {
		element := sbv.insert(index, nearest.prev, nearest)
		element.Set(uint(key % elementsize))
	} else {
		nearest.Set(uint(key % elementsize))
	}
}

// Unset sets a particular bit to false.
func (sbv *SparseBitVector) Unset(key KeyType) {
	index := key / elementsize
	element := sbv.search(index)
	if element != nil && element.index == index {
		element.Unset(uint(key % elementsize))
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
