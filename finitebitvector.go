// This file is distributed under the
// University of Illinois Open Source License.
// See LICENSE.TXT for details.

package sparsebitvector

import "fmt"

type elementwordtype uint64

const bitsperword = 64
const wordsperelement = 2
const elementsize = bitsperword * wordsperelement

// FiniteBitVector provides a bit vector of length elementsize.
type FiniteBitVector [wordsperelement]elementwordtype

func (vec *FiniteBitVector) getWordBit(key uint) (uint, uint) {
	if key < 0 || key >= elementsize {
		panic("key out of range for element")
	}

	word := key / bitsperword
	bit := key % bitsperword
	return word, bit
}

// Set sets a bit to true.
func (vec *FiniteBitVector) Set(key uint) {
	word, bit := vec.getWordBit(key)
	vec[word] |= 1 << bit
}

// TestAndSet sets a bit to true and returns true iff it was previously true.
func (vec *FiniteBitVector) TestAndSet(key uint) bool {
	if vec.Test(key) {
		return false
	}
	vec.Set(key)
	return true
}

// Unset sets a bit to false.
func (vec *FiniteBitVector) Unset(key uint) {
	word, bit := vec.getWordBit(key)
	vec[word] &^= 1 << bit
}

// Clear sets all bits in the Element to false.
func (vec *FiniteBitVector) Clear() {
	for i := range vec {
		vec[i] = 0
	}
}

// Test returns true iff the given bit is true.
func (vec *FiniteBitVector) Test(key uint) bool {
	word, bit := vec.getWordBit(key)
	return (vec[word] & (1 << bit)) != 0
}

// Count returns the number of true bits within the ELement.
func (vec *FiniteBitVector) Count() (count int) {
	for _, word := range vec {
		for word != 0 {
			count += int(word & 1)
			word >>= 1
		}
	}
	return
}

func (vec *FiniteBitVector) String() string {
	result := []byte{}
	for w := range vec {
		for _, c := range fmt.Sprintf("%08x", vec[len(vec)-w-1]) {
			result = append(result, byte(c))
		}
	}
	return string(result)
}
