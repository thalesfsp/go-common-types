package safeorderedmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSafeOrderedMapString(t *testing.T) {
	s := New[int]()
	s.Add("1", 1).Add("2", 2).Add("3", 3)

	assert.Equal(t, `{"1":1,"2":2,"3":3}`, s.String())
}

func TestSafeOrderedMapAdd(t *testing.T) {
	s := New[int]()
	s.Add("1", 1).Add("2", 2).Add("3", 3)

	assert.Equal(t, 3, s.Size())
	assert.True(t, s.Contains("1"))
	assert.True(t, s.Contains("2"))
	assert.True(t, s.Contains("3"))
}

func TestSafeOrderedMapGet(t *testing.T) {
	s := New[int]()
	s.Add("1", 1).Add("2", 2).Add("3", 3)

	if v, ok := s.Get("1"); v != 1 && !ok {
		t.Fatal("missing value")
	}

	if v, ok := s.Get("2"); v != 2 && !ok {
		t.Fatal("missing value")
	}

	if v, ok := s.Get("3"); v != 3 && !ok {
		t.Fatal("missing value")
	}
}

func TestSafeOrderedMapDelete(t *testing.T) {
	s := New[int]()
	s.Add("1", 1).Add("2", 2).Add("3", 3)

	s.Delete("1").Delete("2").Delete("3")

	assert.Equal(t, 0, s.Size())
	assert.False(t, s.Contains("1"))
	assert.False(t, s.Contains("2"))
	assert.False(t, s.Contains("3"))
}

func TestSafeOrderedMapKeys(t *testing.T) {
	s := New[int]()
	s.Add("1", 1).Add("2", 2).Add("3", 3)

	assert.Equal(t, []string{"1", "2", "3"}, s.Keys())
}

func TestSafeOrderedMapValues(t *testing.T) {
	s := New[int]()
	s.Add("1", 1).Add("2", 2).Add("3", 3)

	assert.Equal(t, []int{1, 2, 3}, s.Values())
}

func TestSafeOrderedMapContains(t *testing.T) {
	s := New[int]()
	s.Add("1", 1).Add("2", 2).Add("3", 3)

	assert.True(t, s.Contains("1"))
	assert.True(t, s.Contains("2"))
	assert.True(t, s.Contains("3"))
	assert.False(t, s.Contains("4"))
}

func TestSafeOrderedMapSize(t *testing.T) {
	s := New[int]()
	s.Add("1", 1).Add("2", 2).Add("3", 3)

	assert.Equal(t, 3, s.Size())
}

func TestSafeOrderedMapEmpty(t *testing.T) {
	s := New[int]()

	assert.Equal(t, 0, s.Size())
	assert.True(t, s.Empty())
}

func TestSafeOrderedMapClone(t *testing.T) {
	s := New[int]()
	s.Add("1", 1).Add("2", 2).Add("3", 3)

	c := s.Clone()

	assert.Equal(t, s.Size(), c.Size())
	assert.Equal(t, s.Keys(), c.Keys())
	assert.Equal(t, s.Values(), c.Values())
}

func TestSafeOrderedMapIndex(t *testing.T) {
	s := New[int]()
	s.Add("1", 1).Add("2", 2).Add("3", 3)

	if i, v, ok := s.Index("1"); i != 1 && v != 1 && !ok {
		t.Fatal("missing value")
	}

	if i, v, ok := s.Index("2"); i != 2 && v != 2 && !ok {
		t.Fatal("missing value")
	}

	if i, v, ok := s.Index("3"); i != 3 && v != 3 && !ok {
		t.Fatal("missing value")
	}
}

func TestSafeOrderedMapAll(t *testing.T) {
	s := New[int]()
	s.Add("1", 1).Add("2", 2).Add("3", 3)

	assert.True(t, s.All(func(key string, value int) bool {
		return value > 0
	}))

	assert.False(t, s.All(func(key string, value int) bool {
		return value > 1
	}))
}

func TestSafeOrderedMapMap(t *testing.T) {
	s := New[int]()
	s.Add("1", 1).Add("2", 2).Add("3", 3)

	m := s.Map(func(key string, value int) int {
		return value * 2
	})

	assert.Equal(t, []int{2, 4, 6}, m.Values())
}

func TestSafeOrderedMapFilter(t *testing.T) {
	s := New[int]()
	s.Add("1", 1).Add("2", 2).Add("3", 3)

	f := s.Filter(func(key string, value int) bool {
		return value > 1
	})

	assert.Equal(t, []int{2, 3}, f.Values())
}

func TestSafeOrderedMapEach(t *testing.T) {
	s := New[int]()
	s.Add("1", 1).Add("2", 2).Add("3", 3)

	var sum int
	s.Each(func(key string, value int) {
		sum += value
	})

	assert.Equal(t, 6, sum)
}

func TestSafeOrderedMap_Reduce(t *testing.T) {
	som := New[int]()
	som.Add("a", 2).Add("b", 3).Add("c", 4)

	// Test that reduce returns the correct accumulated value.
	accum := som.Reduce(func(acc int, key string, value int) int {
		return acc + value
	}, 1)
	if accum != 10 {
		t.Errorf("Reduce() failed. Expected %v, but got %v", 10, accum)
	}

	// Test that reduce works when the map is empty.
	somEmpty := New[int]()
	accumEmpty := somEmpty.Reduce(func(acc int, key string, value int) int {
		return acc + value
	}, 0)
	if accumEmpty != 0 {
		t.Errorf("Reduce() failed for empty map. Expected %v, but got %v", 0, accumEmpty)
	}
}

func TestSafeOrderedMapFind(t *testing.T) {
	s := New[int]()
	s.Add("1", 1).Add("2", 2).Add("3", 3)

	if key, value, ok := s.Find(func(key string, value int) bool {
		return value == 1
	}); key != "1" && value != 1 && !ok {
		t.Fatal("missing value")
	}

	if key, value, ok := s.Find(func(key string, value int) bool {
		return value == 2
	}); key != "2" && value != 2 && !ok {
		t.Fatal("missing value")
	}

	if key, value, ok := s.Find(func(key string, value int) bool {
		return value == 3
	}); key != "3" && value != 3 && !ok {
		t.Fatal("missing value")
	}
}

// Any
// TakeWhile
// DropWhile
// Union
// Difference
// Subset
// Superset
// Intersection

func TestSafeOrderedMapAny(t *testing.T) {
	s := New[int]()
	s.Add("1", 1).Add("2", 2).Add("3", 3)

	assert.True(t, s.Any(func(key string, value int) bool {
		return value == 1
	}))

	assert.False(t, s.Any(func(key string, value int) bool {
		return value == 4
	}))
}

func TestSafeOrderedMapTakeWhile(t *testing.T) {
	s := New[int]()
	s.Add("1", 1).Add("2", 2).Add("3", 3).Add("4", 4)

	assert.Equal(t, []int{1, 2}, s.TakeWhile(func(key string, value int) bool {
		return value < 3
	}).Values())
}

func TestSafeOrderedMapDropWhile(t *testing.T) {
	s := New[int]()
	s.Add("1", 1).Add("2", 2).Add("3", 3).Add("4", 4)

	assert.Equal(t, []int{3, 4}, s.DropWhile(func(key string, value int) bool {
		return value < 3
	}).Values())
}

func TestSafeOrderedMapUnion(t *testing.T) {
	s1 := New[int]()
	s1.Add("1", 1).Add("2", 2).Add("3", 3)

	s2 := New[int]()
	s2.Add("4", 4).Add("5", 5).Add("6", 6)

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, s1.Union(s2).Values())
}

func TestSafeOrderedMapDifference(t *testing.T) {
	s1 := New[int]()
	s1.Add("1", 1).Add("2", 2).Add("3", 3)

	s2 := New[int]()
	s2.Add("2", 2).Add("3", 3).Add("4", 4)

	assert.Equal(t, []int{1}, s1.Difference(s2).Values())
}

func TestSafeOrderedMapSubset(t *testing.T) {
	s1 := New[int]()
	s1.Add("1", 1).Add("2", 2).Add("3", 3)

	s2 := New[int]()
	s2.Add("2", 2).Add("3", 3)

	assert.True(t, s2.Subset(s1))
	assert.False(t, s1.Subset(s2))
}

func TestSafeOrderedMapSuperset(t *testing.T) {
	s1 := New[int]()
	s1.Add("1", 1).Add("2", 2).Add("3", 3)

	s2 := New[int]()
	s2.Add("2", 2).Add("3", 3)

	assert.True(t, s1.Superset(s2))
	assert.False(t, s2.Superset(s1))
}

func TestSafeOrderedMapIntersection(t *testing.T) {
	s1 := New[int]()
	s1.Add("1", 1).Add("2", 2).Add("3", 3)

	s2 := New[int]()
	s2.Add("2", 2).Add("3", 3).Add("4", 4)

	assert.Equal(t, []int{2, 3}, s1.Intersection(s2).Values())
}
