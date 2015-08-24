// This file is distributed under the
// University of Illinois Open Source License.
// See LICENSE.TXT for details.

package sparsebitvector

import "testing"
import "reflect"

func TestTrivialOperation(t *testing.T) {
	vec := New()

	if vec.Test(0) {
		t.Error("0 unexpected")
	}

	if vec.Count() != 0 {
		t.Error("not empty")
	}
	if vec.Test(17) {
		t.Error("17 unexpected")
	}

	vec.Set(5)
	if !vec.Test(5) {
		t.Error("expect 5")
	}
	if vec.Test(17) {
		t.Error("17 unexpected")
	}

	vec.Unset(6)
	if !vec.Test(5) {
		t.Error("expected 5")
	}
	if vec.Test(6) {
		t.Error("6 unexpected")
	}

	vec.Unset(5)
	if vec.Test(5) {
		t.Error("5 unexpected")
	}

	if !vec.TestAndSet(1000000000) {
		t.Error("1000000000 unexpected")
	}
	if vec.TestAndSet(1000000000) {
		t.Error("expected 1000000000")
	}
	if !vec.Test(1000000000) {
		t.Error("expected 1000000000")
	}

	if vec.Count() != 1 {
		t.Error("incorrect count ", vec.Count())
	}

	vec.Clear()
	if vec.Test(17) {
		t.Error("17 unexpected")
	}
}

func TestDelete(t *testing.T) {
	vec := New(0, 128, 1000000000)
	if vec.start.next.next.next != nil {
		t.Error("expected 3 elements")
	}

	vec.Unset(0)
	vec.Unset(1000000000)
	if vec.start.next != nil {
		t.Error("expected 1 element")
	}

	vec.Set(0)
	vec.Unset(128)
	if vec.start.next != nil {
		t.Error("expected 1 element")
	}

	vec.Unset(0)
	if vec.start != nil {
		t.Error("expected 0 elements")
	}
}

func TestEquals(t *testing.T) {
	vec1 := New()
	vec2 := New()
	if !vec1.Equals(vec1) {
		t.Error("vec1 should equal itself", vec1)
	}
	if !vec1.Equals(vec2) || !vec2.Equals(vec1) {
		t.Error("vec1 and vec2 should be equal", vec1, vec2)
	}

	vec1 = New(1, 63, 64, 127, 1000000)
	vec2 = New(1, 63, 64, 127, 1000000)
	if !vec1.Equals(vec1) {
		t.Error("vec1 should equal itself", vec1)
	}
	if !vec1.Equals(vec2) || !vec2.Equals(vec1) {
		t.Error("vec1 and vec2 should be equal", vec1, vec2)
	}

	vec1 = New(1, 1000000)
	vec2 = New(1, 1000001)
	if vec1.Equals(vec2) || vec2.Equals(vec1) {
		t.Error("vec1 and vec2 should not be equal", vec1, vec2)
	}

	vec1 = New(0)
	vec2 = New(elementsize)
	if vec1.Equals(vec2) || vec2.Equals(vec1) {
		t.Error("vec1 and vec2 should not be equal", vec1, vec2)
	}

	vec1 = New(1)
	vec2 = New(1, 1000001)
	if vec1.Equals(vec2) || vec2.Equals(vec1) {
		t.Error("vec1 and vec2 should not be equal", vec1, vec2)
	}

	vec1 = New()
	vec2 = New(1, 1000001)
	if vec1.Equals(vec2) || vec2.Equals(vec1) {
		t.Error("vec1 and vec2 should not be equal", vec1, vec2)
	}
}

func TestContains(t *testing.T) {
	vec1 := New()
	vec2 := New()
	if !vec1.Contains(vec1) {
		t.Error("vec1 should contain itself", vec1)
	}
	if !vec1.Contains(vec2) || !vec2.Contains(vec1) {
		t.Error("vec1 and vec2 should contain each other", vec1, vec2)
	}

	vec1 = New(1, 63, 64, 127, 1000000)
	vec2 = New(1, 63, 64, 127, 1000000)
	if !vec1.Contains(vec1) {
		t.Error("vec1 should contain itself", vec1)
	}
	if !vec1.Contains(vec2) || !vec2.Contains(vec1) {
		t.Error("vec1 and vec2 should contain each other", vec1, vec2)
	}

	vec1 = New(1, 1000000)
	vec2 = New(1, 1000001)
	if vec1.Contains(vec2) || vec2.Contains(vec1) {
		t.Error("vec1 and vec2 should not contain each other", vec1, vec2)
	}

	vec1 = New(0)
	vec2 = New(elementsize)
	if vec1.Contains(vec2) || vec2.Contains(vec1) {
		t.Error("vec1 and vec2 should not contain each other", vec1, vec2)
	}

	vec1 = New(1)
	vec2 = New(1, 1000001)
	if vec1.Equals(vec2) {
		t.Error("vec1 should not contain vec2", vec1, vec2)
	}
	if vec2.Equals(vec1) {
		t.Error("vec2 should contain vec1", vec1, vec2)
	}

	vec1 = New()
	vec2 = New(1, 1000001)
	if vec1.Equals(vec2) {
		t.Error("vec1 should not contain vec2", vec1, vec2)
	}
	if vec2.Equals(vec1) {
		t.Error("vec2 should contain vec1", vec1, vec2)
	}
}

func TestSparseBitVectorString(t *testing.T) {
	vec := New()
	if s := vec.String(); s != "[]" {
		t.Error("unexpected string", s)
	}

	vec.Set(0)
	vec.Set(5)
	vec.Set(100)
	vec.Set(1000000000)
	if s := vec.String(); s != "[0 5 100 1000000000]" {
		t.Error("unexpected string", s)
	}
}

func TestIteration(t *testing.T) {
	vec := New()

	test := func() []KeyType {
		result := []KeyType{}
		for i := range vec.Iterate() {
			result = append(result, i)
		}
		return result
	}

	if result := test(); !reflect.DeepEqual(result, []KeyType{}) {
		t.Error("incorrect result", result, vec)
	}

	vec.Set(5)
	if result := test(); !reflect.DeepEqual(result, []KeyType{5}) {
		t.Error("incorrect result", result, vec)
	}

	vec.Set(0)
	vec.Set(65)
	vec.Set(1000000000)
	if result := test(); !reflect.DeepEqual(result, []KeyType{0, 5, 65, 1000000000}) {
		t.Error("incorrect result", result, vec)
	}
}
