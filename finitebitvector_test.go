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
	if vec.String() != "0000000000000020" {
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

	vec.Clear()
	if vec.Test(17) {
		t.Error("17 unexpected", vec)
	}
}

func TestFindNext(t *testing.T) {
	vec := &FiniteBitVector{}

	if i := vec.FindNext(-1); i != -1 {
		t.Error("unexpected result", i, vec)
	}

	vec.Set(0)
	vec.Set(5)
	vec.Set(63)
	if i := vec.FindNext(-1); i != 0 {
		t.Error("unexpected result", i, vec)
	}
	if i := vec.FindNext(0); i != 5 {
		t.Error("unexpected result", i, vec)
	}
	if i := vec.FindNext(4); i != 5 {
		t.Error("unexpected result", i, vec)
	}
	if i := vec.FindNext(5); i != 63 {
		t.Error("unexpected result", i, vec)
	}
	if i := vec.FindNext(63); i != -1 {
		t.Error("unexpected result", i, vec)
	}
	if i := vec.FindNext(65); i != -1 {
		t.Error("unexpected result", i, vec)
	}

	vec.Unset(0)
	vec.Unset(63)
	vec.Set(67)
	vec.Set(127)
	if i := vec.FindNext(-1); i != 5 {
		t.Error("unexpected result", i, vec)
	}
	if i := vec.FindNext(5); i != 67 {
		t.Error("unexpected result", i, vec)
	}
	if i := vec.FindNext(65); i != 67 {
		t.Error("unexpected result", i, vec)
	}
	if i := vec.FindNext(67); i != 127 {
		t.Error("unexpected result", i, vec)
	}
	if i := vec.FindNext(127); i != -1 {
		t.Error("unexpected result", i, vec)
	}

}
