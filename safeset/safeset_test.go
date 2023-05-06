package safeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSafeSetAdd(t *testing.T) {
	s := New[int]()
	s.Add(1).Add(2).Add(3)

	assert.Equal(t, 3, s.Size())
	assert.True(t, s.Contains(1))
	assert.True(t, s.Contains(2))
	assert.True(t, s.Contains(3))
}

func TestSafeSetGet(t *testing.T) {
	s := New(1, 2, 3)

	value, ok := s.Get(2)
	assert.True(t, ok)
	assert.Equal(t, 2, value)

	value, ok = s.Get(4)
	assert.False(t, ok)
	assert.Equal(t, 0, value) // zero value for int
}

func TestSafeSetDelete(t *testing.T) {
	s := New(1, 2, 3)

	s.Delete(2)
	assert.Equal(t, 2, s.Size())
	assert.False(t, s.Contains(2))
}

func TestSafeSetClone(t *testing.T) {
	s := New(1, 2, 3)
	clone := s.Clone()

	assert.Equal(t, s.Size(), clone.Size())
	assert.True(t, clone.Contains(1))
	assert.True(t, clone.Contains(2))
	assert.True(t, clone.Contains(3))
}

func TestSafeSetMap(t *testing.T) {
	s := New(1, 2, 3)
	s.Map(func(value int) int {
		return value * 2
	})

	assert.True(t, s.Contains(2))
	assert.True(t, s.Contains(4))
	assert.True(t, s.Contains(6))
}

func TestSafeSetFilter(t *testing.T) {
	s := New(1, 2, 3, 4, 5)
	filtered := s.Filter(func(value int) bool {
		return value%2 == 0
	})

	assert.Equal(t, 2, filtered.Size())
	assert.True(t, filtered.Contains(2))
	assert.True(t, filtered.Contains(4))
}

func TestSafeSetReduce(t *testing.T) {
	s := New(1, 2, 3, 4)
	sum := s.Reduce(func(acc int, value int) int {
		return acc + value
	}, 0)

	assert.Equal(t, 10, sum)
}

func TestSafeSetUnion(t *testing.T) {
	s1 := New(1, 2, 3)
	s2 := New(3, 4, 5)
	union := s1.Union(s2)

	assert.Equal(t, 5, union.Size())
	assert.True(t, union.Contains(1))
	assert.True(t, union.Contains(2))
	assert.True(t, union.Contains(3))
	assert.True(t, union.Contains(4))
	assert.True(t, union.Contains(5))
}

func TestSafeSetIntersection(t *testing.T) {
	s1 := New(1, 2, 3)
	s2 := New(3, 4, 5)
	intersection := s1.Intersection(s2)

	assert.Equal(t, 1, intersection.Size())
	assert.True(t, intersection.Contains(3))
}

func TestSafeSetSubsetSuperset(t *testing.T) {
	s1 := New(1, 2, 3)
	s2 := New(1, 2)
	s3 := New(3, 4, 5)

	assert.True(t, s2.Subset(s1))
	assert.False(t, s3.Subset(s1))

	assert.True(t, s1.Superset(s2))
	assert.False(t, s1.Superset(s3))
}

func TestSafeSetDifference(t *testing.T) {
	s1 := New(1, 2, 3)
	s2 := New(3, 4, 5)
	difference := s1.Difference(s2)

	assert.Equal(t, 2, difference.Size())
	assert.True(t, difference.Contains(1))
	assert.True(t, difference.Contains(2))
}

func TestSafeSetMarshalUnmarshalJSON(t *testing.T) {
	s := New(1, 2, 3)
	marshaled, err := s.MarshalJSON()

	assert.NoError(t, err)
	assert.NotNil(t, marshaled)

	unmarshaled := New[int]()
	err = unmarshaled.UnmarshalJSON(marshaled)

	assert.NoError(t, err)
	assert.Equal(t, 3, unmarshaled.Size())
	assert.True(t, unmarshaled.Contains(1))
	assert.True(t, unmarshaled.Contains(2))
	assert.True(t, unmarshaled.Contains(3))
}

func TestSafeSetString(t *testing.T) {
	s := New(1, 2, 3)
	expected := "[1, 2, 3]"

	assert.Equal(t, expected, s.String())
}

func TestSafeSetAll(t *testing.T) {
	s := New(2, 4, 6)

	assert.True(t, s.All(func(value int) bool { return value%2 == 0 }))
	assert.False(t, s.All(func(value int) bool { return value > 4 }))
}

func TestSafeSetEach(t *testing.T) {
	s := New(1, 2, 3)
	sum := 0

	s.Each(func(value int) { sum += value })

	assert.Equal(t, 6, sum)
}

func TestSafeSetFind(t *testing.T) {
	s := New(1, 2, 3, 4, 5)

	value, ok := s.Find(func(value int) bool { return value%2 == 0 })

	assert.True(t, ok)
	assert.Equal(t, 2, value)
}

func TestSafeSetAny(t *testing.T) {
	s := New(1, 3, 5)

	// Update the predicate function to check for odd numbers
	assert.True(t, s.Any(func(value int) bool { return value%2 == 1 }))
	assert.False(t, s.Any(func(value int) bool { return value > 6 }))
}

func TestSafeSetTakeWhile(t *testing.T) {
	s := New(1, 2, 3, 4, 5)
	result := s.TakeWhile(func(value int) bool { return value < 4 })

	assert.Equal(t, 3, result.Size())
	assert.True(t, result.Contains(1))
	assert.True(t, result.Contains(2))
	assert.True(t, result.Contains(3))
}

func TestSafeSetDropWhile(t *testing.T) {
	s := New(1, 2, 3, 4, 5)
	result := s.DropWhile(func(value int) bool { return value < 4 })

	assert.Equal(t, 2, result.Size())
	assert.True(t, result.Contains(4))
	assert.True(t, result.Contains(5))
}

func TestSafeSetEmpty(t *testing.T) {
	s := New[int]()

	assert.True(t, s.Empty())
}
