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

// FindNext retruns the next true bit starting from index, or -1 if none exist.
// The initial call should pass index 0.
// Successive calls should pass previous+1.
func (vec *FiniteBitVector) FindNext(index int) int {
	word, bit := vec.getWordBit(uint(index))
	for w := word; w < uint(len(vec)); w++ {
		bits := vec[w] >> bit
		for bits != 0 {
			if bits&1 == 1 {
				return int(w*bitsperword + bit)
			}
			bit++
			bits >>= 1
		}
		bit = 0
	}
	return -1
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
	result := []int{}
	for i := vec.FindNext(0); i != -1; i = vec.FindNext(i + 1) {
		result = append(result, i)
	}
	return fmt.Sprint(result)
}
