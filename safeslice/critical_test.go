package safeslice_test

import (
	"encoding/json"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/thalesfsp/go-common-types/safeslice"
)

func finish(t *testing.T, call func()) {
	t.Helper()

	done := make(chan struct{})
	go func() {
		defer close(done)

		call()
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("slice operation deadlocked")
	}
}

type reentrantValue struct{ mutate func() }

type reentrantNumber int

var duringDecode atomic.Pointer[func()]

func (v *reentrantNumber) UnmarshalJSON(data []byte) error {
	callback := duringDecode.Load()
	if callback != nil {
		(*callback)()
	}

	var number int

	err := json.Unmarshal(data, &number)
	if err != nil {
		return err
	}

	*v = reentrantNumber(number)

	return nil
}

func TestDecoderCallbackReentry(t *testing.T) {
	for _, access := range []func(*safeslice.SafeSlice[reentrantNumber]){
		func(s *safeslice.SafeSlice[reentrantNumber]) { s.Size() },
		func(s *safeslice.SafeSlice[reentrantNumber]) { s.Add(9) },
	} {
		s := safeslice.New[reentrantNumber]()
		callback := func() { access(s) }
		duringDecode.Store(&callback)

		t.Cleanup(func() { duringDecode.Store(nil) })

		var err error

		finish(t, func() { err = json.Unmarshal([]byte(`[7]`), s) })

		if err != nil || s.Size() != 1 || s.Get(0) != 7 {
			t.Fatalf("decoder reentry failed or replacement was not atomic: %v", err)
		}
	}
}

func (v *reentrantValue) String() string {
	v.mutate()
	return "value"
}

func (v *reentrantValue) MarshalJSON() ([]byte, error) {
	v.mutate()
	return []byte(`1`), nil
}

func TestFormattingAndJSONReentry(t *testing.T) {
	for _, encode := range []func(*safeslice.SafeSlice[*reentrantValue]) error{
		func(s *safeslice.SafeSlice[*reentrantValue]) error { _ = s.String(); return nil },
		func(s *safeslice.SafeSlice[*reentrantValue]) error { _, err := json.Marshal(s); return err },
	} {
		s := safeslice.New[*reentrantValue]()
		s.Add(&reentrantValue{mutate: func() { s.Add(nil) }})

		var err error

		finish(t, func() { err = encode(s) })

		if err != nil || s.Size() != 2 {
			t.Fatalf("formatter/JSON callback did not finish safely: %v", err)
		}
	}
}

func TestLastNIndependent(t *testing.T) {
	s := safeslice.New(1, 2, 3)
	tail := s.LastN(2)
	s.Delete(0)

	if !reflect.DeepEqual(tail.ToSlice(), []int{2, 3}) {
		t.Fatalf("original mutation corrupted tail: %v", tail.ToSlice())
	}

	tail.Delete(0)

	if !reflect.DeepEqual(s.ToSlice(), []int{2, 3}) {
		t.Fatal("tail mutation corrupted original")
	}

	if safeslice.New[int]().LastN(1) != nil || s.LastN(0).Size() != 0 || s.LastN(99).Size() != 2 {
		t.Fatal("LastN empty/zero/oversized conventions changed")
	}
}

func TestCallbackReentrySnapshot(t *testing.T) {
	calls := map[string]func(*safeslice.SafeSlice[int], func(int)){
		"all": func(s *safeslice.SafeSlice[int], visit func(int)) {
			s.All(func(v int) bool { visit(v); return true })
		},
		"map": func(s *safeslice.SafeSlice[int], visit func(int)) {
			s.Map(func(v int) int { visit(v); return v })
		},
		"filter": func(s *safeslice.SafeSlice[int], visit func(int)) {
			s.Filter(func(v int) bool { visit(v); return true })
		},
		"each": func(s *safeslice.SafeSlice[int], visit func(int)) { s.Each(visit) },
		"reduce": func(s *safeslice.SafeSlice[int], visit func(int)) {
			s.Reduce(func(a, v int) int { visit(v); return a + v }, 0)
		},
		"find": func(s *safeslice.SafeSlice[int], visit func(int)) {
			s.Find(func(v int) bool { visit(v); return true })
		},
		"any": func(s *safeslice.SafeSlice[int], visit func(int)) {
			s.Any(func(v int) bool { visit(v); return true })
		},
		"take": func(s *safeslice.SafeSlice[int], visit func(int)) {
			s.TakeWhile(func(v int) bool { visit(v); return true })
		},
		"drop": func(s *safeslice.SafeSlice[int], visit func(int)) {
			s.DropWhile(func(v int) bool { visit(v); return false })
		},
		"pluck": func(s *safeslice.SafeSlice[int], visit func(int)) {
			safeslice.Pluck(s, func(v int) int { visit(v); return v })
		},
	}
	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			s := safeslice.New(1)

			finish(t, func() { call(s, func(int) { s.Add(2) }) })

			if s.Size() < 2 {
				t.Fatal("callback mutation lost")
			}
		})
	}
}

func TestNestedReadersSelfAndBinaryProgress(t *testing.T) {
	s := safeslice.New(1, 2, 3)
	other := safeslice.New(1, 2, 3)

	finish(t, func() {
		var workers sync.WaitGroup

		workers.Go(func() {
			for range 2000 {
				s.Add(4).Delete(3)
				other.Add(4).Delete(3)
			}
		})

		for range 2000 {
			s.Mode()
			s.LastN(2)
			s.Union(other)
			s.Difference(other)
			s.Subset(other)
			s.Superset(other)
			s.Intersection(other)

			if !s.Subset(s) || !s.Superset(s) || !s.Difference(s).Empty() {
				t.Error("self-operation identities violated")
			}

			s.Intersection(s)
		}

		workers.Wait()
	})
}
