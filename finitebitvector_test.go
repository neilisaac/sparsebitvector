// This file is distributed under the
// University of Illinois Open Source License.
// See LICENSE.TXT for details.

package sparsebitvector

import "testing"

func TestGetWordBit(t *testing.T) {
	vec := &FiniteBitVector{}
	if w, b := vec.getWordBit(0); w != 0 || b != 0 {
		t.Error("incorrect word or bit: word=", w, "bit=", b)
	}
	if w, b := vec.getWordBit(1); w != 0 || b != 1 {
		t.Error("incorrect word or bit: word=", w, "bit=", b)
	}
	if w, b := vec.getWordBit(bitsperword - 1); w != 0 || b != 63 {
		t.Error("incorrect word or bit: word=", w, "bit=", b)
	}
	if w, b := vec.getWordBit(bitsperword + 1); w != 1 || b != 1 {
		t.Error("incorrect word or bit: word=", w, "bit=", b)
	}
}

func TestTrivialBitVectorOperation(t *testing.T) {
	vec := &element{}

	if vec.Test(0) {
		t.Error("0 unexpected", vec)
	}

	if vec.Count() != 0 {
		t.Error("not empty", vec)
	}
	if vec.Test(17) {
		t.Error("17 unexpected", vec)
	}

	vec.Set(5)
	if !vec.Test(5) {
		t.Error("expected 5", vec)
	}
	if vec.Test(17) {
		t.Error("17 unexpected", vec)
	}
	if vec.String() != "[5]" {
		t.Error("incorrect string", vec)
	}

	vec.Unset(6)
	if !vec.Test(5) {
		t.Error("expected 5", vec)
	}
	if vec.Test(6) {
		t.Error("6 unexpected", vec)
	}

	vec.Unset(5)
	if vec.Test(5) {
		t.Error("5 unexpected", vec)
	}

	if !vec.TestAndSet(100) {
		t.Error("100 unexpected", vec)
	}
	if vec.TestAndSet(100) {
		t.Error("expected 100", vec)
	}
	if !vec.Test(100) {
		t.Error("expected 100", vec)
	}

	if vec.Count() != 1 {
		t.Error("incorrect count ", vec.Count(), vec)
	}

	vec.Set(18)
	if !vec.TestAndUnset(18) {
		t.Error("expected 18")
	}
	if vec.Test(18) {
		t.Error("unexpected 18")
	}

	vec.Clear()
	if vec.Test(17) {
		t.Error("17 unexpected", vec)
	}
}

func TestFindNext(t *testing.T) {
	vec := &FiniteBitVector{}

	if i := vec.FindNext(0); i != -1 {
		t.Error("unexpected result", i, vec)
	}

	vec.Set(0)
	vec.Set(5)
	vec.Set(63)
	if i := vec.FindNext(0); i != 0 {
		t.Error("unexpected result", i, vec)
	}
	if i := vec.FindNext(1); i != 5 {
		t.Error("unexpected result", i, vec)
	}
	if i := vec.FindNext(4); i != 5 {
		t.Error("unexpected result", i, vec)
	}
	if i := vec.FindNext(5); i != 5 {
		t.Error("unexpected result", i, vec)
	}
	if i := vec.FindNext(6); i != 63 {
		t.Error("unexpected result", i, vec)
	}
	if i := vec.FindNext(64); i != -1 {
		t.Error("unexpected result", i, vec)
	}

	vec.Unset(0)
	vec.Unset(63)
	vec.Set(67)
	vec.Set(127)
	if i := vec.FindNext(0); i != 5 {
		t.Error("unexpected result", i, vec)
	}
	if i := vec.FindNext(6); i != 67 {
		t.Error("unexpected result", i, vec)
	}
	if i := vec.FindNext(65); i != 67 {
		t.Error("unexpected result", i, vec)
	}
	if i := vec.FindNext(68); i != 127 {
		t.Error("unexpected result", i, vec)
	}
	if i := vec.FindNext(127); i != 127 {
		t.Error("unexpected result", i, vec)
	}
	if i := vec.FindNext(128); i != -1 {
		t.Error("unexpected result", i, vec)
	}
}

func TestFiniteBitVectorBinaryOperations(t *testing.T) {
	vec1 := NewFiniteBitVector()
	vec2 := NewFiniteBitVector()
	if !vec1.Equals(vec2) || !vec2.Equals(vec1) {
		t.Error("vec1 and vec2 be equal", vec1, vec2)
	}
	if u, i := vec1.UnionAndIntersectionSize(vec2); u != 0 || i != 0 {
		t.Error("incorrect union or intersection size", u, i, vec1, vec2)
	}
	if vec1.UnionWith(vec2); vec1.Count() != 0 {
		t.Error("incorrect union", vec1)
	}
	if vec1.IntersectWith(vec2); vec1.Count() != 0 {
		t.Error("incorrect intersection", vec1)
	}
	if vec1.IntersectWithComplement(vec2); vec1.Count() != 0 {
		t.Error("incorrect complement intersection", vec1)
	}
	if !vec1.Contains(vec2) {
		t.Error("vec1 should contain vec2")
	}

	vec1.Set(3)
	vec2.Set(3)
	if !vec1.Equals(vec2) || !vec2.Equals(vec1) {
		t.Error("vec1 and vec2 be equal", vec1, vec2)
	}
	if !vec1.Contains(vec2) {
		t.Error("vec1 should contain vec2")
	}

	vec1 = NewFiniteBitVector(0, 3, 5, 100, 101)
	vec2 = NewFiniteBitVector(1, 2, 3, 101, 127)
	if vec1.Equals(vec2) || vec2.Equals(vec1) {
		t.Error("vec1 and vec2 should not be equal", vec1, vec2)
	}
	if u, i := vec1.UnionAndIntersectionSize(vec2); u != 8 || i != 2 {
		t.Error("incorrect union or intersection size", u, i, vec1, vec2)
	}
	if vec1.Contains(vec2) {
		t.Error("vec1 should not contain vec2")
	}

	vec1 = NewFiniteBitVector(0, 3, 5, 100, 101)
	vec2 = NewFiniteBitVector(1, 2, 3, 101, 127)
	if vec1.UnionWith(vec2); vec1.Count() != 8 {
		t.Error("incorrect union", vec1)
	}

	vec1 = NewFiniteBitVector(0, 3, 5, 100, 101)
	vec2 = NewFiniteBitVector(1, 2, 3, 101, 127)
	if vec1.IntersectWith(vec2); vec1.Count() != 2 {
		t.Error("incorrect intersection", vec1)
	}

	vec1 = NewFiniteBitVector(0, 3, 5, 100, 101)
	vec2 = NewFiniteBitVector(1, 2, 3, 101, 127)
	if vec1.IntersectWithComplement(vec2); !vec1.Equals(NewFiniteBitVector(0, 5, 100)) {
		t.Error("incorrect complement intersection", vec1)
	}
}
