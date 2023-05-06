package safeorderedmap

import (
	"encoding/json"
	"sync"
)

//////
// Const, vars, and types.
//////

// SafeOrderedMap is a map that preserves the order of keys powered by generics.
type SafeOrderedMap[T any] struct {
	sync.RWMutex

	data map[string]T

	order []string
}

//////
// Methods.
//////

// String is the stringer implementation.
func (m *SafeOrderedMap[T]) String() string {
	m.RLock()
	defer m.RUnlock()

	json, err := json.Marshal(m.data)
	if err != nil {
		return ""
	}

	return string(json)
}

//////
// CRUD operations.

// Add a value in the map.
func (m *SafeOrderedMap[T]) Add(key string, value T) *SafeOrderedMap[T] {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.data[key]; !ok {
		m.order = append(m.order, key)
	}

	m.data[key] = value

	return m
}

// Get a value from the map.
func (m *SafeOrderedMap[T]) Get(key string) (T, bool) {
	m.RLock()
	defer m.RUnlock()

	value, ok := m.data[key]

	return value, ok
}

// Delete a value from the map.
func (m *SafeOrderedMap[T]) Delete(key string) *SafeOrderedMap[T] {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.data[key]; ok {
		delete(m.data, key)

		for i, k := range m.order {
			if k == key {
				m.order = append(m.order[:i], m.order[i+1:]...)

				break
			}
		}
	}

	return m
}

//////
// Key and Values operations.

// Keys returns a list of all keys.
func (m *SafeOrderedMap[T]) Keys() []string {
	m.RLock()
	defer m.RUnlock()

	keys := make([]string, len(m.order))

	copy(keys, m.order)

	return keys
}

// Values returns a list of all values.
func (m *SafeOrderedMap[T]) Values() []T {
	m.RLock()
	defer m.RUnlock()

	values := make([]T, len(m.order))

	for i, key := range m.order {
		values[i] = m.data[key]
	}

	return values
}

//////
// Meta operations.

// Contains checks if the set contains a given element.
func (m *SafeOrderedMap[T]) Contains(key string) bool {
	m.RLock()
	defer m.RUnlock()

	if _, ok := m.Get(key); ok {
		return true
	}

	return false
}

// Size returns the number of elements in the map.
func (m *SafeOrderedMap[T]) Size() int {
	m.RLock()
	defer m.RUnlock()

	return len(m.order)
}

// Empty checks if the map is empty and returns a boolean value.
func (m *SafeOrderedMap[T]) Empty() bool {
	m.RLock()
	defer m.RUnlock()

	return len(m.order) == 0
}

// Clone creates a deep copy of the map and returns it.
func (m *SafeOrderedMap[T]) Clone() *SafeOrderedMap[T] {
	m.RLock()
	defer m.RUnlock()

	clone := New[T]()

	for _, key := range m.order {
		clone.Add(key, m.data[key])
	}

	return clone
}

// Index returns the index and value of the given key.
func (m *SafeOrderedMap[T]) Index(key string) (int, T, bool) {
	m.RLock()
	defer m.RUnlock()

	for i, k := range m.order {
		if k == key {
			return i, m.data[key], true
		}
	}

	return -1, *new(T), false
}

//////
// Collection Operations (Higher-Order Functions).

// All checks if all elements in the map satisfy the predicate.
//
// This method checks if all elements in the map satisfy a given condition
// (predicate). It returns a boolean value, which is true if all elements meet
// the condition, and false otherwise. The All method stops processing as soon
// as it finds an element that does not satisfy the condition.
func (m *SafeOrderedMap[T]) All(predicate func(key string, value T) bool) bool {
	m.RLock()
	defer m.RUnlock()

	for _, key := range m.order {
		if !predicate(key, m.data[key]) {
			return false
		}
	}

	return true
}

// Map applies the function f to all elements in the map and returns a new
// ordered map with the results.
//
// This method applies a given function to all elements in the map and creates
// a new map containing the results. The original map remains unchanged. The new
// map maintains the insertion order of the original map.
func (m *SafeOrderedMap[T]) Map(f func(key string, value T) T) *SafeOrderedMap[T] {
	m.RLock()
	defer m.RUnlock()

	newMap := New[T]()

	for _, key := range m.order {
		newMap.Add(key, f(key, m.data[key]))
	}

	return newMap
}

// Filter filters elements in the map based on the predicate and returns a new
// ordered map containing only elements that satisfy the predicate.
//
// This method creates a new map containing only the elements that satisfy a
// given condition (predicate). The original map remains unchanged. The new map
// maintains the insertion order of the original map.
func (m *SafeOrderedMap[T]) Filter(predicate func(key string, value T) bool) *SafeOrderedMap[T] {
	m.RLock()
	defer m.RUnlock()

	filteredMap := New[T]()

	for _, key := range m.order {
		if predicate(key, m.data[key]) {
			filteredMap.Add(key, m.data[key])
		}
	}

	return filteredMap
}

// Each iterates over the map and calls the given function for each key-value
// pair.
//
// This method iterates over all elements in the map and applies a given
// function to each element. The function can perform any operation, such as
// printing or modifying the elements. However, the Each method itself does not
// return any result.
func (m *SafeOrderedMap[T]) Each(f func(key string, value T)) *SafeOrderedMap[T] {
	m.RLock()
	defer m.RUnlock()

	for _, key := range m.order {
		f(key, m.data[key])
	}

	return m
}

// Reduce accumulates the elements in the map using the given binary function.
//
// Iterates over the map and accumulates the elements using a given binary
// function. It takes a binary function (a function that accepts two arguments)
// and an initial value as input. The method applies the function to the initial
// value and the first element, then to the result and the next element, and so
// on, until all elements in the map have been processed. The final result is a
// single accumulated value.
func (m *SafeOrderedMap[T]) Reduce(reducer func(accum T, key string, value T) T, initial T) T {
	m.RLock()
	defer m.RUnlock()

	accum := initial

	for _, key := range m.order {
		accum = reducer(accum, key, m.data[key])
	}

	return accum
}

// Find returns the first element that satisfies the given predicate.
//
// Iterates over the map and returns the first element that satisfies a given
// predicate. It takes a predicate (a function that returns a boolean) as input.
// If there is an element that satisfies the predicate, it returns that element
// along with the corresponding key and a boolean value true. If no element
// satisfies the predicate, it returns a zero value for the type, an empty
// string for the key, and false for the boolean value.
func (m *SafeOrderedMap[T]) Find(predicate func(key string, value T) bool) (string, T, bool) {
	m.RLock()
	defer m.RUnlock()

	for _, key := range m.order {
		if predicate(key, m.data[key]) {
			return key, m.data[key], true
		}
	}

	return "", *new(T), false
}

// Any checks if any element in the map satisfies the given predicate.
//
// Iterates over the map and checks if any element satisfies a given predicate.
// It takes a predicate (a function that returns a boolean) as input. If any
// element satisfies the predicate, it returns true. If no element satisfies the
// predicate, it returns false.
func (m *SafeOrderedMap[T]) Any(predicate func(key string, value T) bool) bool {
	m.RLock()
	defer m.RUnlock()

	for _, key := range m.order {
		if predicate(key, m.data[key]) {
			return true
		}
	}

	return false
}

// TakeWhile returns a new ordered map containing the longest prefix of elements
// that satisfy the given predicate.
//
// Iterates over the map and returns a new ordered map containing the longest
// prefix of elements that satisfy a given predicate. It takes a predicate (a
// function that returns a boolean) as input. If an element satisfies the
// predicate, it is added to the resulting map. The process stops once an
// element that does not satisfy the predicate is encountered.
func (m *SafeOrderedMap[T]) TakeWhile(predicate func(key string, value T) bool) *SafeOrderedMap[T] {
	m.RLock()
	defer m.RUnlock()

	newMap := New[T]()

	for _, key := range m.order {
		if predicate(key, m.data[key]) {
			newMap.Add(key, m.data[key])
		} else {
			break
		}
	}

	return newMap
}

// DropWhile returns a new ordered map with all elements after (and not
// including) the first element that does not satisfy the given predicate.
//
// Iterates over the map and returns a new ordered map with all elements after
// (and not including) the first element that does not satisfy a given
// predicate. It takes a predicate (a function that returns a boolean) as input.
// The method iterates over the elements in the map and starts adding elements
// to the resulting map once an element that does not satisfy the predicate is
// encountered.
func (m *SafeOrderedMap[T]) DropWhile(predicate func(key string, value T) bool) *SafeOrderedMap[T] {
	m.RLock()
	defer m.RUnlock()

	newMap := New[T]()

	dropping := true
	for _, key := range m.order {
		if dropping && !predicate(key, m.data[key]) {
			dropping = false
		}

		if !dropping {
			newMap.Add(key, m.data[key])
		}
	}

	return newMap
}

//////
// Set operations

// Union returns a new ordered map containing all unique elements from both
// maps. The order of elements in the resulting map will be based on the order
// of elements in the original maps.
func (m *SafeOrderedMap[T]) Union(other *SafeOrderedMap[T]) *SafeOrderedMap[T] {
	m.RLock()
	defer m.RUnlock()

	result := New[T]()
	for _, key := range m.order {
		result.Add(key, m.data[key])
	}

	for _, key := range other.order {
		if _, ok := m.data[key]; !ok {
			result.Add(key, other.data[key])
		}
	}

	return result
}

// Difference returns a new ordered map containing elements present in the
// original map but not in the other map.
func (m *SafeOrderedMap[T]) Difference(other *SafeOrderedMap[T]) *SafeOrderedMap[T] {
	m.RLock()
	defer m.RUnlock()

	result := New[T]()

	for _, key := range m.order {
		if _, ok := other.data[key]; !ok {
			result.Add(key, m.data[key])
		}
	}

	return result
}

// Subset checks if all elements of the original map are present in the other
// map.
func (m *SafeOrderedMap[T]) Subset(other *SafeOrderedMap[T]) bool {
	m.RLock()
	defer m.RUnlock()

	for _, key := range m.order {
		if _, ok := other.data[key]; !ok {
			return false
		}
	}

	return true
}

// Superset checks if all elements of the other map are present in the original map.
func (m *SafeOrderedMap[T]) Superset(other *SafeOrderedMap[T]) bool {
	return other.Subset(m)
}

// Intersection returns a new ordered map containing elements present in both
// maps.
func (m *SafeOrderedMap[T]) Intersection(other *SafeOrderedMap[T]) *SafeOrderedMap[T] {
	m.RLock()
	defer m.RUnlock()

	result := New[T]()

	for _, key := range m.order {
		if _, ok := other.data[key]; ok {
			result.Add(key, m.data[key])
		}
	}

	return result
}

//////
// Conversion Operations.
//////

// MarshalJSON implements json.Marshaler interface for SafeOrderedMap.
func (m *SafeOrderedMap[T]) MarshalJSON() ([]byte, error) {
	m.RLock()
	defer m.RUnlock()

	jsonMap := make(map[string]T)

	for _, key := range m.order {
		jsonMap[key] = m.data[key]
	}

	return json.Marshal(jsonMap)
}

// UnmarshalJSON implements json.Unmarshaler interface for SafeOrderedMap.
func (m *SafeOrderedMap[T]) UnmarshalJSON(data []byte) error {
	m.Lock()
	defer m.Unlock()

	var temp map[string]T
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	m.order = []string{}

	for key := range temp {
		m.order = append(m.order, key)

		m.data[key] = temp[key]
	}

	return nil
}

//////
// Factory.
//////

// New creates a new Safe Ordered Map.
func New[T any]() *SafeOrderedMap[T] {
	return &SafeOrderedMap[T]{
		data:  make(map[string]T),
		order: []string{},
	}
}
