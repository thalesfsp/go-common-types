package safeslice

import (
	"reflect"
	"testing"
)

//nolint:goconst
func TestSafeSliceString(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)

	expected := "[1 2 3]"
	actual := s.String()

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceAdd(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)

	expected := "[1 2 3]"
	actual := s.String()

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceGet(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)

	expected := 2
	actual := s.Get(1)

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceDelete(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)

	s.Delete(1)

	expected := "[1 3]"
	actual := s.String()

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceContains(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)

	expected := true
	actual := s.Contains(2)

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceSize(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)

	expected := 3
	actual := s.Size()

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceEmpty(t *testing.T) {
	s := New[int]()

	expected := true
	actual := s.Empty()

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceClone(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)

	expected := "[1 2 3]"
	actual := s.Clone().String()

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceIndex(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)

	expected := 1
	actual, _ := s.Index(2)

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceUnique(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)
	s.Add(3)

	expected := "[1 2 3]"
	actual := s.Unique().String()

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceAll(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)

	expected := true
	actual := s.All(func(i int) bool { return i > 0 })

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceMap(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)

	expected := "[2 4 6]"
	actual := s.Map(func(i int) int { return i * 2 }).String()

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceFilter(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)

	expected := "[2 3]"
	actual := s.Filter(func(i int) bool { return i > 1 }).String()

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceEach(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)

	expected := "[1 2 3]"
	actual := s.Each(func(i int) {}).String()

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceReduce(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)

	expected := 6
	actual := s.Reduce(func(a, b int) int { return a + b }, 0)

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceFind(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)

	expected := 2
	actual := s.Find(func(i int) bool { return i > 1 })

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceAny(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)

	expected := true
	actual := s.Any(func(i int) bool { return i > 1 })

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceTakeWhile(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)

	expected := "[1 2]"
	actual := s.TakeWhile(func(i int) bool { return i < 3 }).String()

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceDropWhile(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)

	expected := "[3]"
	actual := s.DropWhile(func(i int) bool { return i < 3 }).String()

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceUnion(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)

	o := New[int]()

	o.Add(3)
	o.Add(4)
	o.Add(5)

	expected := "[1 2 3 4 5]"
	actual := s.Union(o).String()

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceDifference(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)

	o := New[int]()

	o.Add(3)
	o.Add(4)
	o.Add(5)

	expected := "[4 5]"
	actual := s.Difference(o).String()

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceSubset(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)

	o := New[int]()

	o.Add(3)
	o.Add(4)
	o.Add(5)

	expected := false
	actual := s.Subset(o)

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceSuperset(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)

	o := New[int]()

	o.Add(3)
	o.Add(4)
	o.Add(5)

	expected := false
	actual := s.Superset(o)

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceIntersection(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)

	o := New[int]()

	o.Add(3)
	o.Add(4)
	o.Add(5)

	expected := "[3]"
	actual := s.Intersection(o).String()

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceFrequency(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)
	s.Add(3)

	expected := map[int]int{
		1: 1,
		2: 1,
		3: 2,
	}
	actual := s.Frequency()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceMode(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)
	s.Add(3)

	expected := []int{3}
	actual := s.Mode()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

//nolint:unconvert
func TestSafeSliceMarshalJSON(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)

	expected := "[1,2,3]"
	actual, err := s.MarshalJSON()
	if err != nil {
		t.Errorf("Expected %v, got %v", expected, actual)
	}

	if string(expected) != string(actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceUnmarshalJSON(t *testing.T) {
	s := New[int]()

	expected := "[1, 2, 3]"
	err := s.UnmarshalJSON([]byte(expected))
	if err != nil {
		t.Errorf("Expected %v, got %v", expected, err)
	}

	actual := s.String()

	expected2 := "[1 2 3]"

	if expected2 != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceFirst(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)

	expected := 1
	actual, ok := s.First()
	if !ok {
		t.Errorf("Expected %v, got %v", expected, actual)
	}

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestSafeSliceLast(t *testing.T) {
	s := New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)

	expected := 3
	actual, ok := s.Last()
	if !ok {
		t.Errorf("Expected %v, got %v", expected, actual)
	}

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
