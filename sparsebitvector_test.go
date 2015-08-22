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

func TestIteration(t *testing.T) {
	vec := New()

	test := func() []int {
		result := []int{}
		for i := range vec.Iterate() {
			result = append(result, i)
		}
		return result
	}

	if result := test(); !reflect.DeepEqual(result, []int{}) {
		t.Error("incorrect result", result, vec)
	}

	vec.Set(5)
	if result := test(); !reflect.DeepEqual(result, []int{5}) {
		t.Error("incorrect result", result, vec)
	}

	vec.Set(0)
	vec.Set(65)
	if result := test(); !reflect.DeepEqual(result, []int{0, 5, 65}) {
		t.Error("incorrect result", result, vec)
	}
}
