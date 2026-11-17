// Package safeset provides a thread-safe implementation of a Set data structure
// that preserves the order of elements and is powered by generics.
package safeset

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"sync"
)

//////
// Const, vars, and types.
//////

// SafeSet is a set that preserves the order of keys powered by generics.
type SafeSet[T any] struct {
	data *setData[T]
}

// setData is shared by value-receiver methods without copying the mutex.
type setData[T any] struct {
	sync.RWMutex

	values []T
}

// equal uses Go equality when safe, and structural equality for noncomparable values.
func equal[T any](a, b T) bool {
	if reflect.ValueOf(a).Comparable() && reflect.ValueOf(b).Comparable() {
		return any(a) == any(b)
	}

	return reflect.DeepEqual(a, b)
}

//////
// Methods.
//////

// String is the stringer implementation.
func (s SafeSet[T]) String() string {
	// Will hold the string representations of the values in the set.
	var stringSlice []string

	// Iterate over the values in the set and convert each value to its string
	// representation.
	for _, value := range s.Values() {
		stringRepresentation := fmt.Sprintf("%v", value)

		if stringRepresentation != "" {
			stringSlice = append(stringSlice, stringRepresentation)
		}
	}

	// Return the joined string slice.
	return strings.Join(stringSlice, ",")
}

// ToSlice converts the SafeSet to a slice of its values.
func (s SafeSet[T]) ToSlice() []T {
	return s.Values()
}

//////
// CRUD operations.

// Add an element to the set.
func (s *SafeSet[T]) Add(value T) *SafeSet[T] {
	s.data.Lock()
	defer s.data.Unlock()

	if !slices.ContainsFunc(s.data.values, func(current T) bool { return equal(current, value) }) {
		s.data.values = append(s.data.values, value)
	}

	return s
}

// Get retrieves an element from the slice at the specified index.
func (s *SafeSet[T]) Get(index int) (T, bool) {
	s.data.RLock()
	defer s.data.RUnlock()

	if index < 0 || index >= len(s.data.values) {
		return *new(T), false
	}

	return s.data.values[index], true
}

// Delete removes an element from the slice at the specified index.
func (s *SafeSet[T]) Delete(index int) *SafeSet[T] {
	s.data.Lock()
	defer s.data.Unlock()

	if index < 0 || index >= len(s.data.values) {
		return s
	}

	s.data.values = slices.Delete(s.data.values, index, index+1)

	return s
}

// First returns the first element in the set.
func (s *SafeSet[T]) First() (T, bool) {
	return s.Get(0)
}

// Last returns the last element in the set.
func (s *SafeSet[T]) Last() (T, bool) {
	s.data.RLock()
	defer s.data.RUnlock()

	if len(s.data.values) == 0 {
		return *new(T), false
	}

	return s.data.values[len(s.data.values)-1], true
}

//////
// Values operations.

// Values returns a list of all values in the set.
func (s *SafeSet[T]) Values() []T {
	s.data.RLock()
	defer s.data.RUnlock()

	return slices.Clone(s.data.values)
}

//////
// Meta operations.

// Contains checks if the set contains a given element.
func (s *SafeSet[T]) Contains(value T) bool {
	s.data.RLock()
	defer s.data.RUnlock()

	return slices.ContainsFunc(s.data.values, func(current T) bool { return equal(current, value) })
}

// Size returns the number of elements in the set.
func (s *SafeSet[T]) Size() int {
	s.data.RLock()
	defer s.data.RUnlock()

	return len(s.data.values)
}

// Empty checks if the set is empty and returns a boolean value.
func (s *SafeSet[T]) Empty() bool {
	return s.Size() == 0
}

// Clone creates a deep copy of the set and returns it.
func (s *SafeSet[T]) Clone() *SafeSet[T] {
	return New[T](s.Values()...)
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
	newSet := New[T]()

	for _, value := range s.Values() {
		newSet.Add(f(value))
	}

	return newSet
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
	return slices.ContainsFunc(s.Values(), predicate)
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
	if s.data == other.data {
		return s.Clone()
	}

	result := s.Clone()

	for _, value := range other.Values() {
		result.Add(value)
	}

	return result
}

// Difference returns a new set containing elements present in the original set but not in the other set.
func (s *SafeSet[T]) Difference(other *SafeSet[T]) *SafeSet[T] {
	result := New[T]()
	if s.data == other.data {
		return result
	}

	for _, value := range s.Values() {
		if !other.Contains(value) {
			result.Add(value)
		}
	}

	return result
}

// Subset checks if all elements of the original set are present in the other set.
func (s *SafeSet[T]) Subset(other *SafeSet[T]) bool {
	if s.data == other.data {
		return true
	}

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
	if s.data == other.data {
		return s.Clone()
	}

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
// It marshals the set as a JSON array of its unique values.
func (s *SafeSet[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values())
}

// UnmarshalJSON implements json.Unmarshaler interface for SafeSet.
// It reconstructs the set from a JSON array of values.
func (s *SafeSet[T]) UnmarshalJSON(data []byte) error {
	var values []T

	err := json.Unmarshal(data, &values)
	if err != nil {
		return err
	}

	decoded := New[T](values...)
	if s.data == nil {
		s.data = decoded.data

		return nil
	}

	s.data.Lock()
	defer s.data.Unlock()

	s.data.values = decoded.data.values

	return nil
}

//////
// Factory.
//////

// New creates a new SafeSet.
func New[T any](v ...T) *SafeSet[T] {
	set := &SafeSet[T]{data: &setData[T]{RWMutex: sync.RWMutex{}, values: make([]T, 0, len(v))}}
	for _, value := range v {
		set.Add(value)
	}

	return set
}

//////
// Exported Functionalities.
//////

// Pluck returns a new slice with the result of applying the given predicate
// to each element of the slice.
func Pluck[T, R comparable](s *SafeSet[T], predicate func(T) R) []R {
	result := []R{}

	for _, item := range s.Values() {
		r := predicate(item)

		if r != *new(R) {
			result = append(result, predicate(item))
		}
	}

	return result
}
