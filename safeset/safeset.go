package safeset

import (
	"fmt"
	"strings"

	"github.com/thalesfsp/go-common-types/safeorderedmap"
	"github.com/thalesfsp/go-common-types/shared"
)

//////
// Const, vars, and types.
//////

// SafeSet is a set that preserves the order of keys powered by generics.
type SafeSet[T any] struct {
	data *safeorderedmap.SafeOrderedMap[T]
}

//////
// Methods.
//////

// String is the stringer implementation.
func (s SafeSet[T]) String() string {
	s.data.RLock()
	defer s.data.RUnlock()

	// Shoud print only the values. Should use string builder.
	var sb strings.Builder

	sb.WriteString("[")

	for i, value := range s.data.Values() {
		sb.WriteString(fmt.Sprintf("%v", value))

		if i < s.data.Size()-1 {
			sb.WriteString(", ")
		}
	}

	sb.WriteString("]")

	return sb.String()
}

//////
// CRUD operations.

// Add an element to the set.
func (s *SafeSet[T]) Add(value T) *SafeSet[T] {
	s.data.Add(shared.GenerateHash(value), value)

	return s
}

// Get retrieves an element from the slice at the specified index.
func (s *SafeSet[T]) Get(index int) (T, bool) {
	s.data.RLock()
	defer s.data.RUnlock()

	if index < 0 || index >= s.data.Size() {
		return *new(T), false
	}

	return s.data.Values()[index], true
}

// Delete removes an element from the slice at the specified index.
func (s *SafeSet[T]) Delete(index int) *SafeSet[T] {
	s.data.Lock()
	defer s.data.Unlock()

	if index < 0 || index >= s.data.Size() {
		return s
	}

	s.data.Delete(shared.GenerateHash(s.data.Values()[index]))

	return s
}

// First returns the first element in the set.
func (s *SafeSet[T]) First() (T, bool) {
	s.data.RLock()
	defer s.data.RUnlock()

	if s.data.Empty() {
		return *new(T), false
	}

	return s.data.Values()[0], true
}

// Last returns the last element in the set.
func (s *SafeSet[T]) Last() (T, bool) {
	s.data.RLock()
	defer s.data.RUnlock()

	if s.data.Empty() {
		return *new(T), false
	}

	return s.data.Values()[s.data.Size()-1], true
}

//////
// Values operations.

// Values returns a list of all values in the set.
func (s *SafeSet[T]) Values() []T {
	return s.data.Values()
}

//////
// Meta operations.

// Contains checks if the set contains a given element.
func (s *SafeSet[T]) Contains(value T) bool {
	_, ok := s.data.Get(shared.GenerateHash(value))

	return ok
}

// Size returns the number of elements in the set.
func (s *SafeSet[T]) Size() int {
	return s.data.Size()
}

// Empty checks if the set is empty and returns a boolean value.
func (s *SafeSet[T]) Empty() bool {
	return s.data.Empty()
}

// Clone creates a deep copy of the set and returns it.
func (s *SafeSet[T]) Clone() *SafeSet[T] {
	clone := New[T]()

	for _, value := range s.data.Values() {
		clone.Add(value)
	}

	return clone
}

//////
// Collection Operations (Higher-Order Functions).

// All returns true if all elements in the set satisfy the given predicate.
func (s *SafeSet[T]) All(predicate func(value T) bool) bool {
	for _, value := range s.Values() {
		if !predicate(value) {
			return false
		}
	}

	return true
}

// Map returns a new set containing the results of applying the given function
// to each element.
func (s *SafeSet[T]) Map(f func(value T) T) *SafeSet[T] {
	for _, value := range s.Values() {
		s.Add(f(value))
	}

	return s
}

// Filter returns a new set containing only the elements that satisfy the given
// predicate.
func (s *SafeSet[T]) Filter(predicate func(value T) bool) *SafeSet[T] {
	result := New[T]()

	for _, value := range s.Values() {
		if predicate(value) {
			result.Add(value)
		}
	}

	return result
}

// Each iterates over the set and calls the given function for each element.
func (s *SafeSet[T]) Each(f func(value T)) *SafeSet[T] {
	for _, value := range s.Values() {
		f(value)
	}

	return s
}

// Reduce reduces the set to a single value by iteratively calling the reducer
// function for each element and passing along an accumulator.
func (s *SafeSet[T]) Reduce(reducer func(acc T, value T) T, initialValue T) T {
	acc := initialValue

	for _, value := range s.Values() {
		acc = reducer(acc, value)
	}

	return acc
}

// Find returns the first element that satisfies the given predicate.
func (s *SafeSet[T]) Find(predicate func(value T) bool) (T, bool) {
	for _, value := range s.Values() {
		if predicate(value) {
			return value, true
		}
	}

	return *new(T), false
}

// Any returns true if any element in the set satisfies the given predicate.
func (s *SafeSet[T]) Any(predicate func(value T) bool) bool {
	for _, value := range s.Values() {
		if predicate(value) {
			return true
		}
	}

	return false
}

// TakeWhile returns a new set containing the first n elements that satisfy the
// given predicate.
func (s *SafeSet[T]) TakeWhile(predicate func(value T) bool) *SafeSet[T] {
	result := New[T]()

	for _, value := range s.Values() {
		if predicate(value) {
			result.Add(value)
		} else {
			break
		}
	}

	return result
}

// DropWhile returns a new set containing all elements except the first n
// elements that satisfy the given predicate.
func (s *SafeSet[T]) DropWhile(predicate func(value T) bool) *SafeSet[T] {
	result := New[T]()

	for _, value := range s.Values() {
		if predicate(value) {
			continue
		}

		result.Add(value)
	}

	return result
}

//////
// Set operations.

// Union returns a new set containing all unique elements from both sets.
func (s *SafeSet[T]) Union(other *SafeSet[T]) *SafeSet[T] {
	result := s.Clone()

	for _, value := range other.Values() {
		result.Add(value)
	}

	return result
}

// Difference returns a new set containing elements present in the original set but not in the other set.
func (s *SafeSet[T]) Difference(other *SafeSet[T]) *SafeSet[T] {
	result := New[T]()

	for _, value := range s.Values() {
		if !other.Contains(value) {
			result.Add(value)
		}
	}

	return result
}

// Subset checks if all elements of the original set are present in the other set.
func (s *SafeSet[T]) Subset(other *SafeSet[T]) bool {
	for _, value := range s.Values() {
		if !other.Contains(value) {
			return false
		}
	}

	return true
}

// Superset checks if all elements of the other set are present in the original set.
func (s *SafeSet[T]) Superset(other *SafeSet[T]) bool {
	return other.Subset(s)
}

// Intersection returns a new set containing elements present in both sets.
func (s *SafeSet[T]) Intersection(other *SafeSet[T]) *SafeSet[T] {
	result := New[T]()

	for _, value := range s.Values() {
		if other.Contains(value) {
			result.Add(value)
		}
	}

	return result
}

//////
// Conversion Operations.
//////

// MarshalJSON implements json.Marshaler interface for SafeSet.
func (s *SafeSet[T]) MarshalJSON() ([]byte, error) {
	return s.data.MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler interface for SafeSet.
func (s *SafeSet[T]) UnmarshalJSON(data []byte) error {
	return s.data.UnmarshalJSON(data)
}

//////
// Factory.
//////

// New creates a new SafeSet.
func New[T any](v ...T) *SafeSet[T] {
	set := &SafeSet[T]{
		data: safeorderedmap.New[T](),
	}

	for _, value := range v {
		set.Add(value)
	}

	return set
}
