// This file is distributed under the
// University of Illinois Open Source License.
// See LICENSE.TXT for details.

package sparsebitvector

import "fmt"

// ElementWordType is the word type used for bit vector elements.
type ElementWordType uint64

// BitsPerWord is equivalent to sizeof(ElementWordType)
const BitsPerWord = 64

// WordsPerElement is the number of ElementWordTypes per Element.
const WordsPerElement = 2

// ElementSize is the number of bits stored per Element.
const ElementSize = WordsPerElement * BitsPerWord

// Element contains a discrete slice of elements
// [index*ElementSize, (index+1) * ElementSize -1] within a bit vector.
type Element struct {
	index KeyType
	words [WordsPerElement]ElementWordType
	prev  *Element
	next  *Element
}

func (e *Element) getWordBit(key KeyType) (KeyType, KeyType) {
	bit := key - e.index*ElementSize
	if bit < 0 || bit >= ElementSize {
		panic("key out of range for element")
	}

	word := bit / BitsPerWord
	return word, (bit % BitsPerWord)
}

// Set sets a bit to true.
func (e *Element) Set(key KeyType) {
	word, bit := e.getWordBit(key)
	e.words[word] |= 1 << bit
}

// TestAndSet sets a bit to true and returns true iff it was previously true.
func (e *Element) TestAndSet(key KeyType) bool {
	if e.Test(key) {
		return false
	}
	e.Set(key)
	return true
}

// Unset sets a bit to false.
func (e *Element) Unset(key KeyType) {
	word, bit := e.getWordBit(key)
	e.words[word] &^= 1 << bit
}

// Clear sets all bits in the Element to false.
func (e *Element) Clear() {
	for i := range e.words {
		e.words[i] = 0
	}
}

// Test returns true iff the given bit is true.
func (e *Element) Test(key KeyType) bool {
	word, bit := e.getWordBit(key)
	return (e.words[word] & (1 << bit)) != 0
}

// Count returns the number of true bits within the ELement.
func (e *Element) Count() (count int) {
	for _, word := range e.words {
		if word == 0 {
			continue
		}
		for i := ElementWordType(0); i < BitsPerWord; i++ {
			if word&(1<<i) != 0 {
				count++
			}
		}
	}
	return
}

func (e *Element) String() string {
	result := []byte{}
	for w := range e.words {
		for _, c := range fmt.Sprintf("%08x", e.words[len(e.words)-w-1]) {
			result = append(result, byte(c))
		}
	}
	return string(result)
}
