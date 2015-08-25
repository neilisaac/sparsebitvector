// This file is distributed under the
// University of Illinois Open Source License.
// See LICENSE.TXT for details.

package sparsebitvector

// element is used internally by SparseBitVector.
type element struct {
	FiniteBitVector
	index KeyType
	prev  *element
	next  *element
}

func (sbv *SparseBitVector) create(index KeyType, prev, next *element) *element {
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

func (sbv *SparseBitVector) delete(e *element) {
	if sbv.start == e {
		sbv.start = e.next
	}
	if sbv.current == e {
		sbv.current = e.next
	}
	if e.prev != nil {
		e.prev.next = e.next
	}
	if e.next != nil {
		e.next.prev = e.prev
	}
}

func (sbv *SparseBitVector) search(index KeyType) *element {
	if sbv.current == nil {
		if sbv.start == nil {
			return nil
		}
		sbv.current = sbv.start
	}

	if sbv.current.index > index {
		for e := sbv.current; e != nil; e = e.prev {
			sbv.current = e
			if e.index == index {
				break
			}
		}
	} else if sbv.current.index < index {
		for e := sbv.current; e != nil; e = e.next {
			sbv.current = e
			if e.index == index {
				break
			}
		}
	}

	return sbv.current
}
