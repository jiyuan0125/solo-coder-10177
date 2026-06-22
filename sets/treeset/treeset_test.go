// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package treeset

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/emirpasic/gods/v2/testutils"
)

func TestSetNew(t *testing.T) {
	set := New[int](2, 1)
	if actualValue := set.Size(); actualValue != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	values := set.Values()
	if actualValue := values[0]; actualValue != 1 {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
	if actualValue := values[1]; actualValue != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
}

func TestSetAdd(t *testing.T) {
	set := New[int]()
	set.Add()
	set.Add(1)
	set.Add(2)
	set.Add(2, 3)
	set.Add()
	if actualValue := set.Empty(); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
	if actualValue := set.Size(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	testutils.SameElements(t, set.Values(), []int{1, 2, 3})
}

func TestSetContains(t *testing.T) {
	set := New[int]()
	set.Add(3, 1, 2)
	if actualValue := set.Contains(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := set.Contains(1); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := set.Contains(1, 2, 3); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := set.Contains(1, 2, 3, 4); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
}

func TestSetRemove(t *testing.T) {
	set := New[int]()
	set.Add(3, 1, 2)
	set.Remove()
	if actualValue := set.Size(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	set.Remove(1)
	if actualValue := set.Size(); actualValue != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	set.Remove(3)
	set.Remove(3)
	set.Remove()
	set.Remove(2)
	if actualValue := set.Size(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}
}

func TestSetEach(t *testing.T) {
	set := New[string]()
	set.Add("c", "a", "b")
	set.Each(func(index int, value string) {
		switch index {
		case 0:
			if actualValue, expectedValue := value, "a"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 1:
			if actualValue, expectedValue := value, "b"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 2:
			if actualValue, expectedValue := value, "c"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			t.Errorf("Too many")
		}
	})
}

func TestSetMap(t *testing.T) {
	set := New[string]()
	set.Add("c", "a", "b")
	mappedSet := set.Map(func(index int, value string) string {
		return "mapped: " + value
	})
	if actualValue, expectedValue := mappedSet.Contains("mapped: a", "mapped: b", "mapped: c"), true; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := mappedSet.Contains("mapped: a", "mapped: b", "mapped: x"), false; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if mappedSet.Size() != 3 {
		t.Errorf("Got %v expected %v", mappedSet.Size(), 3)
	}
}

func TestSetSelect(t *testing.T) {
	set := New[string]()
	set.Add("c", "a", "b")
	selectedSet := set.Select(func(index int, value string) bool {
		return value >= "a" && value <= "b"
	})
	if actualValue, expectedValue := selectedSet.Contains("a", "b"), true; actualValue != expectedValue {
		fmt.Println("A: ", selectedSet.Contains("b"))
		t.Errorf("Got %v (%v) expected %v (%v)", actualValue, selectedSet.Values(), expectedValue, "[a b]")
	}
	if actualValue, expectedValue := selectedSet.Contains("a", "b", "c"), false; actualValue != expectedValue {
		t.Errorf("Got %v (%v) expected %v (%v)", actualValue, selectedSet.Values(), expectedValue, "[a b]")
	}
	if selectedSet.Size() != 2 {
		t.Errorf("Got %v expected %v", selectedSet.Size(), 3)
	}
}

func TestSetAny(t *testing.T) {
	set := New[string]()
	set.Add("c", "a", "b")
	any := set.Any(func(index int, value string) bool {
		return value == "c"
	})
	if any != true {
		t.Errorf("Got %v expected %v", any, true)
	}
	any = set.Any(func(index int, value string) bool {
		return value == "x"
	})
	if any != false {
		t.Errorf("Got %v expected %v", any, false)
	}
}

func TestSetAll(t *testing.T) {
	set := New[string]()
	set.Add("c", "a", "b")
	all := set.All(func(index int, value string) bool {
		return value >= "a" && value <= "c"
	})
	if all != true {
		t.Errorf("Got %v expected %v", all, true)
	}
	all = set.All(func(index int, value string) bool {
		return value >= "a" && value <= "b"
	})
	if all != false {
		t.Errorf("Got %v expected %v", all, false)
	}
}

func TestSetFind(t *testing.T) {
	set := New[string]()
	set.Add("c", "a", "b")
	foundIndex, foundValue := set.Find(func(index int, value string) bool {
		return value == "c"
	})
	if foundValue != "c" || foundIndex != 2 {
		t.Errorf("Got %v at %v expected %v at %v", foundValue, foundIndex, "c", 2)
	}
	foundIndex, foundValue = set.Find(func(index int, value string) bool {
		return value == "x"
	})
	if foundValue != "" || foundIndex != -1 {
		t.Errorf("Got %v at %v expected %v at %v", foundValue, foundIndex, nil, nil)
	}
}

func TestSetChaining(t *testing.T) {
	set := New[string]()
	set.Add("c", "a", "b")
}

func TestSetIteratorNextOnEmpty(t *testing.T) {
	set := New[string]()
	it := set.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty set")
	}
}

func TestSetIteratorPrevOnEmpty(t *testing.T) {
	set := New[string]()
	it := set.Iterator()
	for it.Prev() {
		t.Errorf("Shouldn't iterate on empty set")
	}
}

func TestSetIteratorNext(t *testing.T) {
	set := New[string]()
	set.Add("c", "a", "b")
	it := set.Iterator()
	count := 0
	for it.Next() {
		count++
		index := it.Index()
		value := it.Value()
		switch index {
		case 0:
			if actualValue, expectedValue := value, "a"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 1:
			if actualValue, expectedValue := value, "b"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 2:
			if actualValue, expectedValue := value, "c"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			t.Errorf("Too many")
		}
		if actualValue, expectedValue := index, count-1; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	}
	if actualValue, expectedValue := count, 3; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestSetIteratorPrev(t *testing.T) {
	set := New[string]()
	set.Add("c", "a", "b")
	it := set.Iterator()
	for it.Prev() {
	}
	count := 0
	for it.Next() {
		count++
		index := it.Index()
		value := it.Value()
		switch index {
		case 0:
			if actualValue, expectedValue := value, "a"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 1:
			if actualValue, expectedValue := value, "b"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 2:
			if actualValue, expectedValue := value, "c"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			t.Errorf("Too many")
		}
		if actualValue, expectedValue := index, count-1; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	}
	if actualValue, expectedValue := count, 3; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestSetIteratorBegin(t *testing.T) {
	set := New[string]()
	it := set.Iterator()
	it.Begin()
	set.Add("a", "b", "c")
	for it.Next() {
	}
	it.Begin()
	it.Next()
	if index, value := it.Index(), it.Value(); index != 0 || value != "a" {
		t.Errorf("Got %v,%v expected %v,%v", index, value, 0, "a")
	}
}

func TestSetIteratorEnd(t *testing.T) {
	set := New[string]()
	it := set.Iterator()

	if index := it.Index(); index != -1 {
		t.Errorf("Got %v expected %v", index, -1)
	}

	it.End()
	if index := it.Index(); index != 0 {
		t.Errorf("Got %v expected %v", index, 0)
	}

	set.Add("a", "b", "c")
	it.End()
	if index := it.Index(); index != set.Size() {
		t.Errorf("Got %v expected %v", index, set.Size())
	}

	it.Prev()
	if index, value := it.Index(), it.Value(); index != set.Size()-1 || value != "c" {
		t.Errorf("Got %v,%v expected %v,%v", index, value, set.Size()-1, "c")
	}
}

func TestSetIteratorFirst(t *testing.T) {
	set := New[string]()
	set.Add("a", "b", "c")
	it := set.Iterator()
	if actualValue, expectedValue := it.First(), true; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if index, value := it.Index(), it.Value(); index != 0 || value != "a" {
		t.Errorf("Got %v,%v expected %v,%v", index, value, 0, "a")
	}
}

func TestSetIteratorLast(t *testing.T) {
	set := New[string]()
	set.Add("a", "b", "c")
	it := set.Iterator()
	if actualValue, expectedValue := it.Last(), true; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if index, value := it.Index(), it.Value(); index != 2 || value != "c" {
		t.Errorf("Got %v,%v expected %v,%v", index, value, 2, "c")
	}
}

func TestSetIteratorNextTo(t *testing.T) {
	// Sample seek function, i.e. string starting with "b"
	seek := func(index int, value string) bool {
		return strings.HasSuffix(value, "b")
	}

	// NextTo (empty)
	{
		set := New[string]()
		it := set.Iterator()
		for it.NextTo(seek) {
			t.Errorf("Shouldn't iterate on empty set")
		}
	}

	// NextTo (not found)
	{
		set := New[string]()
		set.Add("xx", "yy")
		it := set.Iterator()
		for it.NextTo(seek) {
			t.Errorf("Shouldn't iterate on empty set")
		}
	}

	// NextTo (found)
	{
		set := New[string]()
		set.Add("aa", "bb", "cc")
		it := set.Iterator()
		it.Begin()
		if !it.NextTo(seek) {
			t.Errorf("Shouldn't iterate on empty set")
		}
		if index, value := it.Index(), it.Value(); index != 1 || value != "bb" {
			t.Errorf("Got %v,%v expected %v,%v", index, value, 1, "bb")
		}
		if !it.Next() {
			t.Errorf("Should go to first element")
		}
		if index, value := it.Index(), it.Value(); index != 2 || value != "cc" {
			t.Errorf("Got %v,%v expected %v,%v", index, value, 2, "cc")
		}
		if it.Next() {
			t.Errorf("Should not go past last element")
		}
	}
}

func TestSetIteratorPrevTo(t *testing.T) {
	// Sample seek function, i.e. string starting with "b"
	seek := func(index int, value string) bool {
		return strings.HasSuffix(value, "b")
	}

	// PrevTo (empty)
	{
		set := New[string]()
		it := set.Iterator()
		it.End()
		for it.PrevTo(seek) {
			t.Errorf("Shouldn't iterate on empty set")
		}
	}

	// PrevTo (not found)
	{
		set := New[string]()
		set.Add("xx", "yy")
		it := set.Iterator()
		it.End()
		for it.PrevTo(seek) {
			t.Errorf("Shouldn't iterate on empty set")
		}
	}

	// PrevTo (found)
	{
		set := New[string]()
		set.Add("aa", "bb", "cc")
		it := set.Iterator()
		it.End()
		if !it.PrevTo(seek) {
			t.Errorf("Shouldn't iterate on empty set")
		}
		if index, value := it.Index(), it.Value(); index != 1 || value != "bb" {
			t.Errorf("Got %v,%v expected %v,%v", index, value, 1, "bb")
		}
		if !it.Prev() {
			t.Errorf("Should go to first element")
		}
		if index, value := it.Index(), it.Value(); index != 0 || value != "aa" {
			t.Errorf("Got %v,%v expected %v,%v", index, value, 0, "aa")
		}
		if it.Prev() {
			t.Errorf("Should not go before first element")
		}
	}
}

func TestSetSerialization(t *testing.T) {
	set := New[string]()
	set.Add("a", "b", "c")

	var err error
	assert := func() {
		if actualValue, expectedValue := set.Size(), 3; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
		if actualValue := set.Contains("a", "b", "c"); actualValue != true {
			t.Errorf("Got %v expected %v", actualValue, true)
		}
		if err != nil {
			t.Errorf("Got error %v", err)
		}
	}

	assert()

	bytes, err := set.ToJSON()
	assert()

	err = set.FromJSON(bytes)
	assert()

	bytes, err = json.Marshal([]interface{}{"a", "b", "c", set})
	if err != nil {
		t.Errorf("Got error %v", err)
	}

	err = json.Unmarshal([]byte(`["1","2","3"]`), &set)
	if err != nil {
		t.Errorf("Got error %v", err)
	}
}

func TestSetString(t *testing.T) {
	c := New[int]()
	c.Add(1)
	if !strings.HasPrefix(c.String(), "TreeSet") {
		t.Errorf("String should start with container name")
	}
}

func TestSetIntersection(t *testing.T) {
	set := New[string]()
	another := New[string]()

	intersection := set.Intersection(another)
	if actualValue, expectedValue := intersection.Size(), 0; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	set.Add("a", "b", "c", "d")
	another.Add("c", "d", "e", "f")

	intersection = set.Intersection(another)

	if actualValue, expectedValue := intersection.Size(), 2; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue := intersection.Contains("c", "d"); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
}

func TestSetUnion(t *testing.T) {
	set := New[string]()
	another := New[string]()

	union := set.Union(another)
	if actualValue, expectedValue := union.Size(), 0; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	set.Add("a", "b", "c", "d")
	another.Add("c", "d", "e", "f")

	union = set.Union(another)

	if actualValue, expectedValue := union.Size(), 6; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue := union.Contains("a", "b", "c", "d", "e", "f"); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
}

func TestSetDifference(t *testing.T) {
	set := New[string]()
	another := New[string]()

	difference := set.Difference(another)
	if actualValue, expectedValue := difference.Size(), 0; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	set.Add("a", "b", "c", "d")
	another.Add("c", "d", "e", "f")

	difference = set.Difference(another)

	if actualValue, expectedValue := difference.Size(), 2; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue := difference.Contains("a", "b"); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
}

func benchmarkContains(b *testing.B, set *Set[int], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			set.Contains(n)
		}
	}
}

func benchmarkAdd(b *testing.B, set *Set[int], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			set.Add(n)
		}
	}
}

func benchmarkRemove(b *testing.B, set *Set[int], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			set.Remove(n)
		}
	}
}

func BenchmarkTreeSetContains100(b *testing.B) {
	b.StopTimer()
	size := 100
	set := New[int]()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkContains(b, set, size)
}

func BenchmarkTreeSetContains1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	set := New[int]()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkContains(b, set, size)
}

func BenchmarkTreeSetContains10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	set := New[int]()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkContains(b, set, size)
}

func BenchmarkTreeSetContains100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	set := New[int]()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkContains(b, set, size)
}

func BenchmarkTreeSetAdd100(b *testing.B) {
	b.StopTimer()
	size := 100
	set := New[int]()
	b.StartTimer()
	benchmarkAdd(b, set, size)
}

func BenchmarkTreeSetAdd1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	set := New[int]()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkAdd(b, set, size)
}

func BenchmarkTreeSetAdd10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	set := New[int]()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkAdd(b, set, size)
}

func BenchmarkTreeSetAdd100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	set := New[int]()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkAdd(b, set, size)
}

func BenchmarkTreeSetRemove100(b *testing.B) {
	b.StopTimer()
	size := 100
	set := New[int]()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkRemove(b, set, size)
}

func BenchmarkTreeSetRemove1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	set := New[int]()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkRemove(b, set, size)
}

func BenchmarkTreeSetRemove10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	set := New[int]()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkRemove(b, set, size)
}

func BenchmarkTreeSetRemove100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	set := New[int]()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkRemove(b, set, size)
}

func newIntComparatorWithMultiplier(multiplier int) func(a, b int) int {
	return func(a, b int) int {
		return (a - b) * multiplier
	}
}

func TestSetIntersectionClosureDifferentMultiplier(t *testing.T) {
	set := NewWith(newIntComparatorWithMultiplier(1))
	another := NewWith(newIntComparatorWithMultiplier(2))
	set.Add(1, 2, 3, 4)
	another.Add(3, 4, 5, 6)

	intersection := set.Intersection(another)
	if actualValue := intersection.Size(); actualValue != 2 {
		t.Errorf("Got %v expected %v (same-order multipliers should be semantically equal)", actualValue, 2)
	}
	if actualValue := intersection.Contains(3, 4); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
}

func TestSetUnionClosureDifferentMultiplier(t *testing.T) {
	set := NewWith(newIntComparatorWithMultiplier(1))
	another := NewWith(newIntComparatorWithMultiplier(2))
	set.Add(1, 2, 3, 4)
	another.Add(3, 4, 5, 6)

	union := set.Union(another)
	if actualValue := union.Size(); actualValue != 6 {
		t.Errorf("Got %v expected %v (same-order multipliers should be semantically equal)", actualValue, 6)
	}
}

func TestSetDifferenceClosureDifferentMultiplier(t *testing.T) {
	set := NewWith(newIntComparatorWithMultiplier(1))
	another := NewWith(newIntComparatorWithMultiplier(2))
	set.Add(1, 2, 3, 4)
	another.Add(3, 4, 5, 6)

	difference := set.Difference(another)
	if actualValue := difference.Size(); actualValue != 2 {
		t.Errorf("Got %v expected %v (same-order multipliers should be semantically equal)", actualValue, 2)
	}
	if actualValue := difference.Contains(1, 2); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
}

func TestSetIntersectionClosureSameMultiplier(t *testing.T) {
	set := NewWith(newIntComparatorWithMultiplier(2))
	another := NewWith(newIntComparatorWithMultiplier(2))
	set.Add(1, 2, 3, 4)
	another.Add(3, 4, 5, 6)

	intersection := set.Intersection(another)
	if actualValue := intersection.Size(); actualValue != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	if actualValue := intersection.Contains(3, 4); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
}

func TestSetIntersectionReversedComparator(t *testing.T) {
	normalCmp := func(a, b int) int {
		return a - b
	}
	reversedCmp := func(a, b int) int {
		return b - a
	}

	set := NewWith(normalCmp)
	another := NewWith(reversedCmp)
	set.Add(1, 2, 3, 4)
	another.Add(3, 4, 5, 6)

	intersection := set.Intersection(another)
	if actualValue := intersection.Size(); actualValue != 0 {
		t.Errorf("Got %v expected %v (reversed comparator should yield empty set)", actualValue, 0)
	}
}

func TestSetUnionReversedComparator(t *testing.T) {
	normalCmp := func(a, b int) int {
		return a - b
	}
	reversedCmp := func(a, b int) int {
		return b - a
	}

	set := NewWith(normalCmp)
	another := NewWith(reversedCmp)
	set.Add(1, 2, 3)
	another.Add(4, 5, 6)

	union := set.Union(another)
	if actualValue := union.Size(); actualValue != 0 {
		t.Errorf("Got %v expected %v (reversed comparator should yield empty set)", actualValue, 0)
	}
}

type person struct {
	name string
	age  int
}

func newPersonComparator(field string) func(a, b person) int {
	return func(a, b person) int {
		switch field {
		case "name":
			if a.name < b.name {
				return -1
			}
			if a.name > b.name {
				return 1
			}
			return 0
		case "age":
			return a.age - b.age
		default:
			return 0
		}
	}
}

func TestSetIntersectionStructDifferentField(t *testing.T) {
	set := NewWith(newPersonComparator("name"))
	another := NewWith(newPersonComparator("age"))
	set.Add(person{"alice", 30}, person{"bob", 25}, person{"carol", 35})
	another.Add(person{"dave", 25}, person{"eve", 30}, person{"frank", 35})

	intersection := set.Intersection(another)
	if actualValue := intersection.Size(); actualValue != 0 {
		t.Errorf("Got %v expected %v (different field comparators should yield empty set)", actualValue, 0)
	}
}

func TestSetIntersectionStructSameField(t *testing.T) {
	set := NewWith(newPersonComparator("age"))
	another := NewWith(newPersonComparator("age"))
	set.Add(person{"alice", 30}, person{"bob", 25})
	another.Add(person{"carol", 30}, person{"dave", 28})

	intersection := set.Intersection(another)
	if actualValue := intersection.Size(); actualValue != 1 {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
}

func TestSetFloat64SpecialValues(t *testing.T) {
	set := New[float64]()
	another := New[float64]()

	zero := 0.0
	nan := zero / zero
	posZero := 0.0
	negZero := -0.0

	set.Add(nan, posZero, negZero, 1.0, 2.0)
	another.Add(nan, posZero, negZero, 2.0, 3.0)

	intersection := set.Intersection(another)
	if intersection.Size() == 0 {
		t.Errorf("Intersection should not be empty for same comparator")
	}

	union := set.Union(another)
	if union.Size() == 0 {
		t.Errorf("Union should not be empty for same comparator")
	}

	difference := set.Difference(another)
	if difference.Size() == 0 {
		t.Errorf("Difference should not be empty for same comparator")
	}
}

func TestSetChainedOperations(t *testing.T) {
	a := New[int]()
	b := New[int]()
	c := New[int]()
	d := New[int]()

	a.Add(1, 2, 3, 4, 5)
	b.Add(4, 5, 6, 7, 8)
	c.Add(3, 4, 5, 9, 10)
	d.Add(5, 11, 12)

	result := a.Intersection(b).Union(c).Intersection(d)
	if actualValue := result.Size(); actualValue != 1 {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
	if actualValue := result.Contains(5); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
}

func TestSetResultOperations(t *testing.T) {
	set := New[string]()
	another := New[string]()
	set.Add("a", "b", "c")
	another.Add("b", "c", "d")

	result := set.Intersection(another)

	if actualValue := result.Contains("b"); actualValue != true {
		t.Errorf("Contains got %v expected %v", actualValue, true)
	}

	values := result.Values()
	if len(values) != 2 {
		t.Errorf("Values got %v expected %v", len(values), 2)
	}

	count := 0
	result.Each(func(index int, value string) {
		count++
	})
	if count != 2 {
		t.Errorf("Each got %v expected %v", count, 2)
	}

	str := result.String()
	if !strings.Contains(str, "TreeSet") {
		t.Errorf("String should contain TreeSet")
	}

	jsonBytes, err := result.ToJSON()
	if err != nil {
		t.Errorf("ToJSON error: %v", err)
	}
	if len(jsonBytes) == 0 {
		t.Errorf("ToJSON should return non-empty bytes")
	}

	it := result.Iterator()
	iterCount := 0
	for it.Next() {
		iterCount++
	}
	if iterCount != 2 {
		t.Errorf("Iterator got %v expected %v", iterCount, 2)
	}
}

func TestSetEmptySetOperations(t *testing.T) {
	set := New[int]()
	another := New[int]()
	another.Add(1, 2, 3)

	intersection := set.Intersection(another)
	if actualValue := intersection.Size(); actualValue != 0 {
		t.Errorf("Intersection got %v expected %v", actualValue, 0)
	}

	union := set.Union(another)
	if actualValue := union.Size(); actualValue != 3 {
		t.Errorf("Union got %v expected %v", actualValue, 3)
	}

	difference := set.Difference(another)
	if actualValue := difference.Size(); actualValue != 0 {
		t.Errorf("Difference got %v expected %v", actualValue, 0)
	}
}

func TestSetIntersectionSameCmpBuiltin(t *testing.T) {
	set := New[int]()
	another := New[int]()
	set.Add(1, 2, 3)
	another.Add(2, 3, 4)

	intersection := set.Intersection(another)
	if actualValue := intersection.Size(); actualValue != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	if actualValue := intersection.Contains(2, 3); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
}

func TestSetLargeSetIntersection(t *testing.T) {
	size := 10000
	set := New[int]()
	another := New[int]()

	for i := 0; i < size; i++ {
		set.Add(i)
	}
	for i := size / 2; i < size+size/2; i++ {
		another.Add(i)
	}

	intersection := set.Intersection(another)
	expected := size / 2
	if actualValue := intersection.Size(); actualValue != expected {
		t.Errorf("Got %v expected %v", actualValue, expected)
	}
}

func BenchmarkSetIntersection10000(b *testing.B) {
	size := 10000
	set := New[int]()
	another := New[int]()
	for i := 0; i < size; i++ {
		set.Add(i)
	}
	for i := size / 2; i < size+size/2; i++ {
		another.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set.Intersection(another)
	}
}

func BenchmarkSetUnion10000(b *testing.B) {
	size := 10000
	set := New[int]()
	another := New[int]()
	for i := 0; i < size; i++ {
		set.Add(i)
	}
	for i := size / 2; i < size+size/2; i++ {
		another.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set.Union(another)
	}
}

func BenchmarkSetDifference10000(b *testing.B) {
	size := 10000
	set := New[int]()
	another := New[int]()
	for i := 0; i < size; i++ {
		set.Add(i)
	}
	for i := size / 2; i < size+size/2; i++ {
		another.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set.Difference(another)
	}
}

type record struct {
	ID   int
	Name string
}

func newRecordSubComparator() func(a, b record) int {
	return func(a, b record) int {
		return a.ID - b.ID
	}
}

func newRecordSignComparator() func(a, b record) int {
	return func(a, b record) int {
		if a.ID < b.ID {
			return -1
		}
		if a.ID > b.ID {
			return 1
		}
		return 0
	}
}

func TestSetStructSubVsSignComparatorIntersection(t *testing.T) {
	set := NewWith(newRecordSubComparator())
	another := NewWith(newRecordSignComparator())

	set.Add(record{1, "a"}, record{2, "b"}, record{3, "c"})
	another.Add(record{2, "x"}, record{3, "y"}, record{4, "z"})

	intersection := set.Intersection(another)
	if actualValue := intersection.Size(); actualValue != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
}

func TestSetStructSubVsSignComparatorUnion(t *testing.T) {
	set := NewWith(newRecordSubComparator())
	another := NewWith(newRecordSignComparator())

	set.Add(record{1, "a"}, record{2, "b"})
	another.Add(record{3, "c"}, record{4, "d"})

	union := set.Union(another)
	if actualValue := union.Size(); actualValue != 4 {
		t.Errorf("Got %v expected %v", actualValue, 4)
	}
}

func TestSetStructSubVsSignComparatorDifference(t *testing.T) {
	set := NewWith(newRecordSubComparator())
	another := NewWith(newRecordSignComparator())

	set.Add(record{1, "a"}, record{2, "b"}, record{3, "c"})
	another.Add(record{2, "x"}, record{3, "y"})

	difference := set.Difference(another)
	if actualValue := difference.Size(); actualValue != 1 {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
}

func TestSetLargeSetThreeOpsWithinOneSecond(t *testing.T) {
	size := 10000
	set := New[int]()
	another := New[int]()

	for i := 0; i < size; i++ {
		set.Add(i)
	}
	for i := size/2; i < size+size/2; i++ {
		another.Add(i)
	}

	start := time.Now()

	_ = set.Intersection(another)
	_ = set.Union(another)
	_ = set.Difference(another)

	elapsed := time.Since(start)
	if elapsed > time.Second {
		t.Errorf("Three ops took %v, expected <= 1s", elapsed)
	}
}
