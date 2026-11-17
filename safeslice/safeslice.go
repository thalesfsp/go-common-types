// package safeslice

package safeslice

import (
	"encoding/json"
	"fmt"
	"slices"
	"sync"
)

//////
// Const, vars, and types.
//////

// SafeSlice is a slice that is safe for concurrent use powered by generics.
type SafeSlice[T comparable] struct {
	sync.RWMutex

	data []T
}

//////
// Methods.
//////

// String is the stringer implementation.
func (s *SafeSlice[T]) String() string {
	return fmt.Sprintf("%v", s.snapshot())
}

//////
// CRUD operations.

// Add appends a new element to the end of the slice.
func (s *SafeSlice[T]) Add(item T) *SafeSlice[T] {
	s.Lock()
	defer s.Unlock()

	s.data = append(s.data, item)

	return s
}

// Get retrieves an element from the slice at the specified index.
func (s *SafeSlice[T]) Get(index int) T {
	s.RLock()
	defer s.RUnlock()

	if index < 0 || index >= len(s.data) {
		return *new(T)
	}

	return s.data[index]
}

// Delete removes an element from the slice at the specified index.
func (s *SafeSlice[T]) Delete(index int) *SafeSlice[T] {
	s.Lock()
	defer s.Unlock()

	if index < 0 || index >= len(s.data) {
		return s
	}

	s.data = append(s.data[:index], s.data[index+1:]...)

	return s
}

// First return the first element.
func (s *SafeSlice[T]) First() (T, bool) {
	s.RLock()
	defer s.RUnlock()

	if len(s.data) == 0 {
		return *new(T), false
	}

	return s.data[0], true
}

// Last return the last element.
func (s *SafeSlice[T]) Last() (T, bool) {
	s.RLock()
	defer s.RUnlock()

	if len(s.data) == 0 {
		return *new(T), false
	}

	return s.data[len(s.data)-1], true
}

// ToSlice returns the underlying slice.
func (s *SafeSlice[T]) ToSlice() []T {
	s.RLock()
	defer s.RUnlock()

	return s.data
}

// LastN return the last N elements as a new slice.
func (s *SafeSlice[T]) LastN(n int) *SafeSlice[T] {
	s.RLock()
	defer s.RUnlock()

	if len(s.data) == 0 {
		return nil
	}

	if n > len(s.data) {
		n = len(s.data)
	}

	tail := make([]T, n)
	copy(tail, s.data[len(s.data)-n:])

	return New[T](tail...)
}

//////
// Meta operations.

// Contains checks if the given element is present in the slice.
func (s *SafeSlice[T]) Contains(item T) bool {
	s.RLock()
	defer s.RUnlock()

	return slices.Contains(s.data, item)
}

// Size returns the number of elements in the slice.
func (s *SafeSlice[T]) Size() int {
	s.RLock()
	defer s.RUnlock()

	return len(s.data)
}

// Empty checks if the slice is empty.
func (s *SafeSlice[T]) Empty() bool {
	s.RLock()
	defer s.RUnlock()

	return len(s.data) == 0
}

// Clone returns a new copy of the slice.
func (s *SafeSlice[T]) Clone() *SafeSlice[T] {
	return New[T](s.snapshot()...)
}

// snapshot copies container storage before invoking caller code or another collection.
func (s *SafeSlice[T]) snapshot() []T {
	s.RLock()
	defer s.RUnlock()

	return slices.Clone(s.data)
}

// Index returns the index of the first occurrence of the given element in the slice.
// If the element is not found, it returns -1 and false.
func (s *SafeSlice[T]) Index(element T) (int, bool) {
	s.RLock()
	defer s.RUnlock()

	for i, item := range s.data {
		if item == element {
			return i, true
		}
	}

	return -1, false
}

// Unique returns a new SafeSlice with all duplicates removed.
func (s *SafeSlice[T]) Unique() *SafeSlice[T] {
	s.RLock()
	defer s.RUnlock()

	uniqueMap := make(map[T]bool)

	uniqueSlice := New[T]()

	for _, item := range s.data {
		if _, ok := uniqueMap[item]; !ok {
			uniqueMap[item] = true

			uniqueSlice.Add(item)
		}
	}

	return uniqueSlice
}

//////
// Collection Operations (Higher-Order Functions).

// All checks if all elements in the slice satisfy a given condition (predicate) and returns a boolean value.
func (s *SafeSlice[T]) All(predicate func(T) bool) bool {
	for _, item := range s.snapshot() {
		if !predicate(item) {
			return false
		}
	}

	return true
}

// Map applies a given function to all elements in the slice and creates a new slice containing the results.
func (s *SafeSlice[T]) Map(mapper func(T) T) *SafeSlice[T] {
	result := New[T]()

	for _, item := range s.snapshot() {
		result.Add(mapper(item))
	}

	return result
}

// Filter creates a new slice containing only the elements that satisfy a given condition (predicate).
func (s *SafeSlice[T]) Filter(predicate func(T) bool) *SafeSlice[T] {
	result := New[T]()

	for _, item := range s.snapshot() {
		if predicate(item) {
			result.Add(item)
		}
	}

	return result
}

// Each iterates over the slice and calls the given function for each element.
func (s *SafeSlice[T]) Each(f func(T)) *SafeSlice[T] {
	for _, item := range s.snapshot() {
		f(item)
	}

	return s
}

// Reduce applies a given function to all elements in the slice and returns a single result.
func (s *SafeSlice[T]) Reduce(reducer func(T, T) T, initialValue T) T {
	result := initialValue

	for _, item := range s.snapshot() {
		result = reducer(result, item)
	}

	return result
}

// Find returns the first element in the slice that satisfies the given predicate.
// If no element satisfies the predicate, it returns the zero value of the type.
func (s *SafeSlice[T]) Find(predicate func(T) bool) T {
	for _, item := range s.snapshot() {
		if predicate(item) {
			return item
		}
	}

	return *new(T)
}

// Any checks if at least one element in the slice satisfies a given condition (predicate).
func (s *SafeSlice[T]) Any(predicate func(T) bool) bool {
	return slices.ContainsFunc(s.snapshot(), predicate)
}

// TakeWhile creates a new slice containing elements from the original slice
// until the predicate function returns false.
func (s *SafeSlice[T]) TakeWhile(predicate func(T) bool) *SafeSlice[T] {
	result := New[T]()

	for _, item := range s.snapshot() {
		if predicate(item) {
			result.Add(item)
		} else {
			break
		}
	}

	return result
}

// DropWhile creates a new slice without the elements from the original slice
// until the predicate function returns false.
func (s *SafeSlice[T]) DropWhile(predicate func(T) bool) *SafeSlice[T] {
	result := New[T]()

	dropping := true

	for _, item := range s.snapshot() {
		if dropping && predicate(item) {
			continue
		}

		dropping = false

		result.Add(item)
	}

	return result
}

//////
// Set operations.

// Union returns a new slice containing all unique elements from both slices.
func (s *SafeSlice[T]) Union(other *SafeSlice[T]) *SafeSlice[T] {
	result := s.Clone()
	if s == other {
		return result
	}

	for _, item := range other.snapshot() {
		if !result.Contains(item) {
			result.Add(item)
		}
	}

	return result
}

// Difference returns a new slice containing elements present in the
// original slice but not in the other slice.
func (s *SafeSlice[T]) Difference(other *SafeSlice[T]) *SafeSlice[T] {
	result := New[T]()
	if s == other {
		return result
	}

	compared := s.Clone()
	for _, item := range other.snapshot() {
		if !compared.Contains(item) {
			result.Add(item)
		}
	}

	return result
}

// Subset checks if all elements in the slice are present in the other slice.
func (s *SafeSlice[T]) Subset(other *SafeSlice[T]) bool {
	if s == other {
		return true
	}

	compared := other.Clone()
	for _, item := range s.snapshot() {
		if !compared.Contains(item) {
			return false
		}
	}

	return true
}

// Superset checks if all elements in the other slice are present in the slice.
func (s *SafeSlice[T]) Superset(other *SafeSlice[T]) bool {
	return other.Subset(s)
}

// Intersection returns a new slice containing elements present in both slices.
func (s *SafeSlice[T]) Intersection(other *SafeSlice[T]) *SafeSlice[T] {
	if s == other {
		return s.Clone()
	}

	result := New[T]()

	compared := other.Clone()
	for _, item := range s.snapshot() {
		if compared.Contains(item) {
			result.Add(item)
		}
	}

	return result
}

//////
// Statistical operations.

// Frequency returns a map with the frequency of each element in the slice.
func (s *SafeSlice[T]) Frequency() map[T]int {
	s.RLock()
	defer s.RUnlock()

	freq := make(map[T]int)

	for _, item := range s.data {
		freq[item]++
	}

	return freq
}

// Mode returns a slice with the mode(s) of the SafeSlice. The mode is the
// element(s) that appears the most frequently. If all elements appear with the
// same frequency, it returns a slice with all elements.
func (s *SafeSlice[T]) Mode() []T {
	data := s.snapshot()

	if len(data) == 0 {
		return nil
	}

	freqMap := make(map[T]int)
	for _, item := range data {
		freqMap[item]++
	}

	maxFreq := 0
	for _, freq := range freqMap {
		if freq > maxFreq {
			maxFreq = freq
		}
	}

	if maxFreq == 1 {
		uniqueSlice := New[T](data...).Unique()

		return uniqueSlice.data
	}

	modes := make([]T, 0)

	for item, freq := range freqMap {
		if freq == maxFreq {
			modes = append(modes, item)
		}
	}

	return modes
}

//////
// Conversion Operations.
//////

// MarshalJSON marshals the slice to JSON.
func (s *SafeSlice[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.snapshot())
}

// UnmarshalJSON unmarshals the slice from JSON.
func (s *SafeSlice[T]) UnmarshalJSON(data []byte) error {
	var temp []T

	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()

	s.data = temp

	return nil
}

//////
// Factory.
//////

// New creates a new Safe Slice.
func New[T comparable](v ...T) *SafeSlice[T] {
	return &SafeSlice[T]{
		data:    v,
		RWMutex: sync.RWMutex{},
	}
}

//////
// Exported Functionalities.
//////

// Pluck returns a new slice with the result of applying the given predicate
// to each element of the slice.
func Pluck[T, R comparable](s *SafeSlice[T], predicate func(T) R) []R {
	result := []R{}

	for _, item := range s.snapshot() {
		r := predicate(item)

		if r != *new(R) {
			result = append(result, predicate(item))
		}
	}

	return result
}
