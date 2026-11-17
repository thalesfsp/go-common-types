package safeset_test

import (
	"encoding/json"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/thalesfsp/go-common-types/safeset"
)

type formatted int

func (formatted) String() string { return "same" }

func (formatted) GoString() string { return "same" }

func TestEqualityIgnoresFormatting(t *testing.T) {
	type record struct{ First, Last string }

	first, second := record{"Mary Ann", "Lee"}, record{"Mary", "Ann Lee"}

	s := safeset.New(first)
	if s.Contains(second) {
		t.Fatal("distinct records have a false-positive membership result")
	}

	s.Add(second).Add(first)

	if s.Size() != 2 || !reflect.DeepEqual(s.Values(), []record{first, second}) {
		t.Fatalf("formatting collision merged records: %v", s.Values())
	}

	s.Delete(0)

	if s.Contains(first) || !s.Contains(second) {
		t.Fatal("indexed deletion removed the wrong colliding value")
	}

	if mixed := safeset.New[any](1, "1", 1); mixed.Size() != 2 {
		t.Fatalf("mixed dynamic types merged: %#v", mixed.Values())
	}

	if custom := safeset.New(formatted(1), formatted(2), formatted(1)); custom.Size() != 2 {
		t.Fatal("String/GoString hooks defined identity")
	}

	a, b := 1, 1
	pointers := safeset.New(&a, &b, &a)

	a = 2
	if pointers.Size() != 2 || !pointers.Contains(&a) || !pointers.Contains(&b) {
		t.Fatal("pointer identity was replaced by pointee equality")
	}

	slices := safeset.New([]int{1}, []int{2}, []int{1}, []int(nil), []int{})
	if slices.Size() != 4 || !slices.Contains([]int{2}) || slices.Contains([]int{3}) {
		t.Fatalf("noncomparable values lack structural equality: %v", slices.Values())
	}

	wrapped := safeset.New[any](struct{ Value any }{[]int{1}}, struct{ Value any }{[]int{1}})
	if wrapped.Size() != 1 {
		t.Fatal("interface containing a slice must compare without panic")
	}

	nils := safeset.New[any](nil, nil, []int(nil), (*int)(nil))
	if nils.Size() != 3 || !nils.Contains(nil) || !nils.Contains([]int(nil)) {
		t.Fatal("nil values and typed nils must retain their defined equality")
	}
}

func TestConcurrentIndexedOperations(t *testing.T) {
	s := safeset.New(0)

	done := make(chan struct{})
	go func() {
		defer close(done)

		var workers sync.WaitGroup

		workers.Add(2)
		go func() {
			defer workers.Done()

			for i := range 20000 {
				s.Delete(0).Add(i % 8)
			}
		}()
		go func() {
			defer workers.Done()

			for range 20000 {
				s.Get(0)
			}
		}()

		workers.Wait()
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("concurrent indexed operations deadlocked")
	}

	if _, ok := s.Get(-1); ok {
		t.Fatal("negative index accepted")
	}

	before := s.Values()
	s.Delete(100)

	if !reflect.DeepEqual(before, s.Values()) {
		t.Fatal("out-of-bounds deletion changed the set")
	}
}

func TestConcurrentIndexedDeletion(t *testing.T) {
	s := safeset.New(1)

	done := make(chan struct{})
	go func() {
		defer close(done)

		var workers sync.WaitGroup
		for range 2 {
			workers.Go(func() {
				for range 20000 {
					s.Add(1).Delete(0)
				}
			})
		}

		workers.Wait()
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("concurrent indexed deletion deadlocked")
	}

	if !s.Empty() {
		t.Fatal("completed indexed deletions left an entry")
	}
}

func TestReadersAndWriterProgress(t *testing.T) {
	reads := map[string]func(*safeset.SafeSet[int]){
		"string": func(s *safeset.SafeSet[int]) { _ = s.String() },
		"slice":  func(s *safeset.SafeSet[int]) { s.ToSlice() },
		"first":  func(s *safeset.SafeSet[int]) { s.First() },
		"last":   func(s *safeset.SafeSet[int]) { s.Last() },
		"clone":  func(s *safeset.SafeSet[int]) { s.Clone() },
		"map":    func(s *safeset.SafeSet[int]) { s.Map(func(v int) int { return v }) },
		"self": func(s *safeset.SafeSet[int]) {
			s.Union(s)
			s.Intersection(s)
			s.Difference(s)
			s.Subset(s)
			s.Superset(s)
		},
	}
	for name, read := range reads {
		t.Run(name, func(t *testing.T) {
			s := safeset.New(1, 2, 3)

			done := make(chan struct{})
			go func() {
				defer close(done)

				var workers sync.WaitGroup

				workers.Go(func() {
					for i := range 2000 {
						s.Add(i % 8)
					}
				})

				for range 2000 {
					read(s)
				}

				workers.Wait()
			}()

			select {
			case <-done:
			case <-time.After(2 * time.Second):
				t.Fatal("reader and writer deadlocked")
			}

			if !s.Difference(s).Empty() || !s.Subset(s) || !s.Superset(s) {
				t.Fatal("self-operation identities violated")
			}
		})
	}
}

func TestMapCallbackCanMutateSet(t *testing.T) {
	s := safeset.New(1)

	done := make(chan struct{})
	go func() {
		defer close(done)

		s.Map(func(v int) int { s.Add(2); return v })
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("set callback deadlocked")
	}

	if s.Size() != 2 {
		t.Fatal("reentrant insertion lost")
	}
}

func TestJSONConcurrentReplacement(t *testing.T) {
	s := safeset.New(1)

	var workers sync.WaitGroup

	workers.Go(func() {
		for range 1000 {
			s.Add(2)
			s.Values()
		}
	})

	for range 1000 {
		err := json.Unmarshal([]byte(`[1,1,3]`), s)
		if err != nil {
			t.Error(err)
		}

		_, err = json.Marshal(s)
		if err != nil {
			t.Error(err)
		}
	}

	workers.Wait()

	var empty safeset.SafeSet[int]

	err := json.Unmarshal([]byte(`[1,1,3]`), &empty)
	if err != nil || !reflect.DeepEqual(empty.Values(), []int{1, 3}) {
		t.Fatalf("zero-value JSON initialization failed: %v", err)
	}
}
